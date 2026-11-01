#!/usr/bin/env bash
# Stop hook for /next-task. While a roadmap task is marked "done" but not yet verified
# (.claude/state/current-task.json exists), Claude may stop only once lint and tests pass for every service the
# task touched. A failure exits 2, which sends the output back to Claude to fix. After 3 failed attempts the task is
# marked blocked and Claude may stop. Outside a task (no marker) or while a task is still "in progress" (Claude
# pausing to ask something), it does nothing.
set -uo pipefail

input="$(cat)"
dir="$(jq -r '.cwd // empty' <<<"$input" 2>/dev/null)"
dir="${dir:-${CLAUDE_PROJECT_DIR:-$PWD}}"
root="$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null)" || exit 0
cd "$root" || exit 0

marker=.claude/state/current-task.json
[[ -f "$marker" ]] || exit 0
id="$(jq -r '.id' "$marker")"
attempts="$(jq -r '.attempts // 0' "$marker")"

status="$(python3 -I scripts/epics/tasks.py show "$id" 2>/dev/null | jq -r '.status // empty')"
[[ "$status" == "done" ]] || exit 0

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
if [[ -n "$ok" && -n "$testable" ]]; then
  make --no-print-directory test s="$testable" >>"$log" 2>&1 || ok=""
fi

if [[ -n "$ok" ]]; then
  rm -f "$marker"
  exit 0
fi

attempts=$((attempts + 1))
if ((attempts >= 3)); then
  python3 -I scripts/epics/tasks.py block "$id" "lint or tests still fail for: $services after 3 fixes; see the last commit and make lint/test output" >/dev/null
  echo "Stop hook: $id blocked after 3 failed checks of $services." >&2
  exit 0
fi
jq --argjson a "$attempts" '.attempts = $a' "$marker" >"$marker.tmp" && mv "$marker.tmp" "$marker"

{
  echo "Task $id is marked done, but lint or tests fail for: $services (attempt $attempts of 3)."
  echo "Fix the failures, commit the fix as fix(<service>): $id ..., then stop again. Last 60 lines:"
  tail -n 60 "$log"
} >&2
exit 2
