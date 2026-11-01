#!/usr/bin/env bash
# Works through docs/EPICs/ROADMAP.md unattended: one `claude -p "/next-task <id>"` run per task, in roadmap order.
# Stops after the last task of a milestone (review and merge its branch, then run again), on a blocked task, when a
# run fails, or when nothing is left.
#
# Usage: scripts/epics/run-tasks.sh [--dry-run] [--budget <usd per task>]
#   --dry-run   show the next task and the command, run nothing
#   --budget    spending cap per task run, passed to claude --max-budget-usd (default 5)
# Needs: the dev stack up (make dev) for tests against Postgres/Redis, jq, python3, claude.
set -euo pipefail
cd "$(dirname "$0")/../.."

dry=""
budget=5
while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run) dry=1 ;;
    --budget) budget="$2"; shift ;;
    *) echo "unknown option: $1" >&2; exit 2 ;;
  esac
  shift
done

# Commands a task may run without asking. Everything else is denied in print mode.
allowed=(
  Read Edit Write Glob Grep
  "Bash(make *)" "Bash(scripts/epics/*)" "Bash(python3 scripts/epics/*)" "Bash(python3 -I scripts/epics/*)"
  "Bash(scripts/test-service.sh *)" "Bash(bun *)" "Bash(uv *)"
  "Bash(git status *)" "Bash(git diff *)" "Bash(git log *)" "Bash(git show *)" "Bash(git add *)"
  "Bash(git commit *)" "Bash(git switch *)" "Bash(git branch *)" "Bash(git rev-parse *)"
)
denied=(
  "Bash(make prod *)" "Bash(make prod-down *)" "Bash(make secrets *)" "Bash(make backup *)"
  "Bash(git push *)" "Bash(git reset *)" "Bash(git rebase *)" "Bash(git merge *)" "Bash(rm -rf *)"
)

tasks() { python3 -I scripts/epics/tasks.py "$@"; }
mkdir -p .claude/state/runs
total=0

while true; do
  next="$(tasks next)"
  id="$(jq -r '.id // empty' <<<"$next")"
  if [[ -z "$id" ]]; then
    echo "$(jq -r '.message' <<<"$next")$(jq -r 'if .blocked then " (blocked: " + (.blocked | join(", ")) + ")" else "" end' <<<"$next")"
    break
  fi
  milestone="$(jq -r '.milestone' <<<"$next")"
  branch="$(jq -r '.branch' <<<"$next")"
  last="$(jq -r '.last_in_milestone' <<<"$next")"
  echo "▶ $id $(jq -r '.title' <<<"$next")  [$milestone, branch $branch]"

  cmd=(claude -p "/next-task $id" --permission-mode acceptEdits --max-budget-usd "$budget" --output-format json
    --allowedTools "${allowed[@]}" --disallowedTools "${denied[@]}")
  if [[ -n "$dry" ]]; then
    printf '  would run:'; printf ' %q' "${cmd[@]}"; echo
    break
  fi

  log=".claude/state/runs/$(date +%Y%m%d-%H%M%S)-$id.json"
  rc=0
  "${cmd[@]}" >"$log" || rc=$?
  if ((rc != 0)); then
    echo "  run failed (exit $rc); output in $log"
    break
  fi
  cost="$(jq -r '.total_cost_usd // 0' "$log")"
  total="$(python3 -c "print(round($total + $cost, 2))")"
  echo "  $(jq -r '.result // ""' "$log" | tail -n 1)   (\$$cost, total \$$total)"

  # The Stop hook verifies and clears the marker; if it didn't run, verify here.
  if [[ -f .claude/state/current-task.json && "$(tasks show "$id" | jq -r '.status')" == "done" ]]; then
    services="$(tasks changed)"
    testable="$(tr ' ' '\n' <<<"$services" | grep -v -x contracts | tr '\n' ' ' | xargs)"
    if make --no-print-directory lint s="$services" && { [[ -z "$testable" ]] || make --no-print-directory test s="$testable"; }; then
      rm -f .claude/state/current-task.json
    else
      tasks block "$id" "lint or tests fail after the run (checked by run-tasks.sh)"
    fi
  fi

  status="$(tasks show "$id" | jq -r '.status')"
  if [[ "$status" != "done" ]]; then
    echo "stopped: $id is $status. Read the task file and $log, fix or decide, then run again."
    break
  fi
  if [[ "$last" == "true" ]]; then
    echo "milestone $milestone complete. Review and merge branch $branch, then run again."
    break
  fi
done
