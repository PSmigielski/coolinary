#!/usr/bin/env bash
# Works through docs/EPICs/ROADMAP.md unattended, in batches: one `claude -p "/next-task <ID> <ID> ..."` run per
# batch of up to --jobs ready tasks, each built by its own subagent in its own git worktree (a batch of one is built
# directly). Stops after the last task of a milestone (review and merge its branch, then run again), on a blocked
# task, when a run fails, or when nothing is left.
#
# Usage: scripts/epics/run-tasks.sh [--dry-run] [--jobs <n>] [--budget <usd per run>]
#   --dry-run   show the next batch and the command, run nothing
#   --jobs      tasks per batch (default 4; 1 builds one task at a time)
#   --budget    spending cap per batch run, passed to claude --max-budget-usd (default: no cap)
# Needs: the dev stack up (make dev) for tests against Postgres/Redis, jq, python3, claude.
set -euo pipefail
cd "$(dirname "$0")/../.."

dry=""
budget=""
jobs=4
while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run) dry=1 ;;
    --jobs) jobs="$2"; shift ;;
    --budget) budget="$2"; shift ;;
    *) echo "unknown option: $1" >&2; exit 2 ;;
  esac
  shift
done
[[ "$jobs" =~ ^[1-9][0-9]*$ ]] || { echo "--jobs needs a positive number" >&2; exit 2; }

# Commands a task may run without asking. Everything else is denied in print mode. Subagents (Agent) build the
# tasks of a batch; they work in .claude/worktrees/<ID> (cd there) and the main session lands their commits with
# scripts/epics/worktree.sh.
allowed=(
  Read Edit Write Glob Grep Agent SendMessage
  "Bash(make *)" "Bash(scripts/epics/*)" "Bash(python3 scripts/epics/*)" "Bash(python3 -I scripts/epics/*)"
  "Bash(scripts/test-service.sh *)" "Bash(bun *)" "Bash(uv *)" "Bash(cd *)"
  "Bash(git status *)" "Bash(git diff *)" "Bash(git log *)" "Bash(git show *)" "Bash(git add *)"
  "Bash(git commit *)" "Bash(git switch *)" "Bash(git branch *)" "Bash(git rev-parse *)"
  "Bash(git checkout --ours *)" "Bash(git checkout --theirs *)"
)
denied=(
  "Bash(make prod *)" "Bash(make prod-down *)" "Bash(make secrets *)" "Bash(make backup *)"
  "Bash(git push *)" "Bash(git reset *)" "Bash(git rebase *)" "Bash(git merge *)" "Bash(rm -rf *)"
)

tasks() { python3 -I scripts/epics/tasks.py "$@"; }
marker=.claude/state/current-task.json
mkdir -p .claude/state/runs
total=0
declare -A requeued=()

while true; do
  batch="$(tasks ready "$jobs")"
  mapfile -t ids < <(jq -r '.tasks[]?.id' <<<"$batch")
  if ((${#ids[@]} == 0)); then
    echo "$(jq -r '.message' <<<"$batch")$(jq -r 'if .blocked then " (blocked: " + (.blocked | join(", ")) + ")" else "" end' <<<"$batch")"
    break
  fi
  milestone="$(jq -r '.milestone' <<<"$batch")"
  branch="$(jq -r '.branch' <<<"$batch")"
  echo "▶ ${#ids[@]} task(s) [$milestone, branch $branch]"
  jq -r '.tasks[] | "  \(.id) \(.title)"' <<<"$batch"

  cmd=(claude -p "/next-task ${ids[*]}" --permission-mode acceptEdits --output-format json
    --allowedTools "${allowed[@]}" --disallowedTools "${denied[@]}")
  [[ -n "$budget" ]] && cmd+=(--max-budget-usd "$budget")
  if [[ -n "$dry" ]]; then
    printf '  would run:'; printf ' %q' "${cmd[@]}"; echo
    break
  fi

  name="${ids[0]}"
  ((${#ids[@]} > 1)) && name+="+$((${#ids[@]} - 1))"
  log=".claude/state/runs/$(date +%Y%m%d-%H%M%S)-$name.json"
  rc=0
  "${cmd[@]}" >"$log" || rc=$?
  if ((rc != 0)); then
    echo "  run failed (exit $rc); output in $log"
    break
  fi
  cost="$(jq -r '.total_cost_usd // 0' "$log")"
  total="$(python3 -c "print(round($total + $cost, 2))")"
  jq -r '.result // ""' "$log" | grep -E '^(DONE|BLOCKED|REQUEUED|MILESTONE|NO TASK LEFT|NOTHING READY)' | sed 's/^/  /' || true
  echo "  (\$$cost, total \$$total)"

  # The Stop hook verifies and clears the marker; if it didn't run, verify here.
  if [[ -f "$marker" ]]; then
    mapfile -t marked < <(jq -r '(.ids // [.id]) | .[] | select(. != null)' "$marker")
    open=""
    for id in "${marked[@]}"; do
      [[ "$(tasks show "$id" | jq -r '.status')" == "in progress" ]] && open=1
    done
    if [[ -z "$open" ]]; then
      services="$(tasks changed)"
      testable="$(tr ' ' '\n' <<<"$services" | grep -v -x contracts | tr '\n' ' ' | xargs)"
      if [[ -z "$services" ]] || { make --no-print-directory lint s="$services" && { [[ -z "$testable" ]] || make --no-print-directory test s="$testable"; }; }; then
        rm -f "$marker"
      else
        for id in "${marked[@]}"; do
          [[ "$(tasks show "$id" | jq -r '.status')" == "done" ]] && tasks block "$id" "lint or tests fail after the run (checked by run-tasks.sh)"
        done
        rm -f "$marker"
      fi
    fi
  fi

  stop=""
  for id in "${ids[@]}"; do
    status="$(tasks show "$id" | jq -r '.status')"
    case "$status" in
      done) ;;
      todo)
        if [[ -n "${requeued[$id]:-}" ]]; then
          echo "stopped: $id could not be landed twice. Read $log, then run again."
          stop=1
        fi
        requeued[$id]=1
        ;;
      *)
        echo "stopped: $id is $status. Read the task file and $log, fix or decide, then run again."
        stop=1
        ;;
    esac
  done
  [[ -n "$stop" ]] && break
  # Every task of the batch is done here (or requeued, which is not done): the milestone is complete when no
  # other task of it is left.
  if tasks show "${ids[-1]}" | jq -e '.status == "done" and .last_in_milestone' >/dev/null; then
    echo "milestone $milestone complete. Review and merge branch $branch, then run again."
    break
  fi
done
