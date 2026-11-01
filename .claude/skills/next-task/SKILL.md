---
name: next-task
description: Build the next task of docs/EPICs/ROADMAP.md (or the task ID given) end to end - implement, test, mark done and commit on the milestone branch. Use only when the user runs /next-task or asks to work through the roadmap.
argument-hint: "[task-id]"
disable-model-invocation: true
---

# Build one roadmap task

Build exactly one task from `docs/EPICs/`, then stop. Follow `CLAUDE.md`. All commands run from the repo root.

## 1. Check the working tree

Run `git status --porcelain` and check for `.claude/state/current-task.json`.

- No marker and the tree has changes: stop. Final line: `BLOCKED -: working tree has changes that aren't from a task;
  commit or stash them first`.
- Marker present: a task is in progress. Resume it (step 2 picks it) and read `git diff` to see what's done.

## 2. Pick the task

- With an argument (`$ARGUMENTS`), run `python3 scripts/epics/tasks.py show $ARGUMENTS`. If `ready` is false, stop:
  `BLOCKED $ARGUMENTS: dependencies not done (<ids>)`.
- Without one, run `python3 scripts/epics/tasks.py next`. If it returns `NO TASK LEFT` or `NOTHING READY`, print
  that as the final line and stop.

Note `id`, `file`, `service`, `branch` and `last_in_milestone`.

## 3. Branch and start

- `git switch <branch>` if it exists, otherwise `git switch -c <branch>` from the current HEAD.
- Unless the task is already `in progress`: `python3 scripts/epics/tasks.py start <id>`.

## 4. Read before writing

1. The whole task file, then its epic `README.md`.
2. `docs/EPICs/README.md`: "Conventions for every task", "Free and Plus", "Cross-service calls", "Events".
3. `docs/service-conventions.md` and the service's `.claude/rules/*.md` (for `mobile`: `mobile/AGENTS.md`).
4. The task files listed in "Depends on" when you need their contracts, and the existing code you will change.

If anything is unclear or contradicts the code, another task or the docs, and no `Decided (Qnnn)` rule settles it:
don't guess. Run `python3 scripts/epics/tasks.py block <id> "<what needs deciding, and the options>"`, leave your
changes uncommitted, and stop with `BLOCKED <id>: <reason>`.

## 5. Build

- Contracts first (`contracts/…`, then `make proto`), then migrations and queries (`make sqlc`), then code.
- Gateway routes: `make openapi`. App API types come from it.
- Write tests for the business rules and acceptance criteria as `docs/service-conventions.md` "Tests" describes.
- Only this task. A sketch detail that doesn't fit the code may be corrected; note it for the commit message.

## 6. Verify

- `make lint s="<services you touched>"` and `make test s="<services you touched>"` (add `contracts` to lint when you
  changed protos). Fix until both pass.
- Go through the acceptance criteria one by one. Prove each with a test, or check it against the running stack
  (`make dev`, `make grpc`, `make smoke`). List any that need a person (a device, an outside account) as manual
  checks.

## 7. Finish

1. `python3 scripts/epics/tasks.py done <id>`, then `scripts/epics/roadmap.sh`.
2. `git add` only the files this task changed (plus `docs/EPICs/` status and roadmap updates), then commit:

   ```
   feat(<service>): <id> <task title>

   <two or three lines on what was built>

   Acceptance: <criteria proven by tests or checks>
   Manual checks: <criteria left for a person, or "none">
   Sketch changes: <deviations from the task's sketches, or "none">
   ```

3. When you stop, the Stop hook runs lint and tests on every service the task touched. If it reports a failure, fix
   it and commit the fix (`fix(<service>): <id> ...`) before stopping again.

## Final line

End your reply with exactly one line, so `scripts/epics/run-tasks.sh` can read it:

- `DONE <id>`, or `DONE <id> MILESTONE <Mn> COMPLETE` when `last_in_milestone` was true
- `BLOCKED <id>: <reason>`
- `NO TASK LEFT` or `NOTHING READY`
