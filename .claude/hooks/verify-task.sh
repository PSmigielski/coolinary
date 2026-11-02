#!/usr/bin/env bash
# Stop hook for /next-task. While roadmap tasks are marked "done" but not yet verified
# (.claude/state/current-task.json lists them: one task, or a parallel batch), Claude may stop only once lint and
# tests pass for every service the task or batch touched. A failure exits 2, which sends the output back to Claude to
# fix. After 3 failed attempts the done tasks whose service failed are marked blocked (all of them, if no task's own
# service failed) and Claude may stop. Outside a task (no marker) or while any task of the batch is still
# "in progress" (being built, or Claude pausing to ask something), it does nothing.
set -uo pipefail

input="$(cat)"
dir="$(jq -r '.cwd // empty' <<<"$input" 2>/dev/null)"
dir="${dir:-${CLAUDE_PROJECT_DIR:-$PWD}}"
root="$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null)" || exit 0
cd "$root" || exit 0

marker=.claude/state/current-task.json
[[ -f "$marker" ]] || exit 0
mapfile -t ids < <(jq -r '(.ids // [.id]) | .[] | select(. != null)' "$marker")
attempts="$(jq -r '.attempts // 0' "$marker")"

done_ids=()
for id in "${ids[@]}"; do
  status="$(python3 -I scripts/epics/tasks.py show "$id" 2>/dev/null | jq -r '.status // empty')"
  [[ "$status" == "in progress" ]] && exit 0
  [[ "$status" == "done" ]] && done_ids+=("$id")
done
if ((${#done_ids[@]} == 0)); then
  rm -f "$marker"
  exit 0
fi

services="$(python3 -I scripts/epics/tasks.py changed)"
if [[ -z "$services" ]]; then
  rm -f "$marker"
  exit 0
fi
testable="$(tr ' ' '\n' <<<"$services" | grep -v -x contracts | tr '\n' ' ' | xargs)"

log="$(mktemp)"
trap 'rm -f "$log"' EXIT
ok=1
make --no-print-directory lint s="$services" >"$log" 2>&1 || ok=""
if [[ -n "$testable" ]]; then
  make --no-print-directory test s="$testable" >>"$log" 2>&1 || ok=""
fi

if [[ -n "$ok" ]]; then
  rm -f "$marker"
  exit 0
fi

# make's own failure lines name the target ("*** [Makefile:172: lint-kitchen] Error 1", translated in some locales).
failed="$(grep -oE '\*\*\* \[[^]]*: (lint|test)-[a-z]+\]' "$log" | sed -E 's/.*(lint|test)-([a-z]+)\]$/\2/' | sort -u | xargs)"

attempts=$((attempts + 1))
if ((attempts >= 3)); then
  to_block=()
  for id in "${done_ids[@]}"; do
    svc="$(python3 -I scripts/epics/tasks.py show "$id" | jq -r '.service')"
    [[ " $failed " == *" $svc "* ]] && to_block+=("$id")
  done
  ((${#to_block[@]})) || to_block=("${done_ids[@]}")
  for id in "${to_block[@]}"; do
    python3 -I scripts/epics/tasks.py block "$id" "lint or tests still fail for: ${failed:-$services} after 3 fixes; see the last commits and make lint/test output" >/dev/null
  done
  rm -f "$marker"
  echo "Stop hook: ${to_block[*]} blocked after 3 failed checks of ${failed:-$services}." >&2
  exit 0
fi
jq --argjson a "$attempts" '.attempts = $a' "$marker" >"$marker.tmp" && mv "$marker.tmp" "$marker"

{
  echo "Tasks ${done_ids[*]} are marked done, but lint or tests fail for: ${failed:-$services} (attempt $attempts of 3)."
  echo "Fix the failures, commit each fix as fix(<service>): <ID> ... (the task that caused it), then stop again. Last 60 lines:"
  tail -n 60 "$log"
} >&2
exit 2
