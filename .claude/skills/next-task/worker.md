# Build one roadmap task

How one task from `docs/EPICs/` is built. A `/next-task` batch gives each task to a worker (a subagent) in its own
git worktree; a batch of one is built by the main session in the main checkout. Follow `CLAUDE.md`.

The worker's prompt names: the task ID, the worktree **W** (inside the main checkout, under `.claude/worktrees/`)
and the main checkout **ROOT**. Other workers are building other tasks of the same milestone at the same time, some
in the same service.

## 1. Workers: stay in your worktree

The main session building a task in ROOT skips this section.

- Start every Bash command with `cd W && `. Read, Edit and Write only absolute paths under W. Never change a file
  in ROOT outside W: the main session lands every task there, and the dev stack runs from there.
- Read the task file from `ROOT/<task file>`: the main checkout has its current status and any decision taken
  since the batch started. Everything else (code, other task files, docs) read from W.
- First `git status` and `git log --oneline -3` in W. If W already has changes or commits of this task (an
  interrupted run), read them and continue from there.
- Leave task status to the main session: don't run `tasks.py start|done|block|set` or `scripts/epics/roadmap.sh`, and
  don't edit anything in `docs/EPICs/`.
- Your branch `task/<ID>` is checked out in W: commit there. No other branch, and no push, merge, rebase, reset or
  worktree commands. Don't start subagents to build other tasks.
- Several workers may take the same next migration number or change the same generated files. Do what the task
  needs in W; the main session renumbers and regenerates when it lands the tasks.

## 2. Read before writing

1. The whole task file, then its epic `README.md`.
2. `docs/EPICs/README.md`: "Conventions for every task", "Free and Plus", "Cross-service calls", "Events".
3. `docs/service-conventions.md` and the service's `.claude/rules/*.md` (for `mobile`: `mobile/AGENTS.md`).
4. The task files listed in "Depends on" when you need their contracts, and the existing code you will change.

If anything is unclear or contradicts the code, another task or the docs, and no `Decided (Qnnn)` rule settles it:
don't guess. A worker stops at once and reports `BLOCKED <ID>: <what needs deciding, and the options>` as its final
line, with its changes left uncommitted in W. It may be continued with the answer. The main session building in ROOT
does what `SKILL.md` says for one task.

## 3. Build

Test first, as `.claude/rules/testing.md` describes.

- Contracts first (`contracts/…`, then `make proto`), then migrations and queries (`make sqlc`), so tests compile
  against the real types.
- Then one business rule or acceptance criterion at a time: write its test (named as `docs/service-conventions.md`
  "Tests" describes), run it and see it fail for the expected reason, write the code until it passes, refactor.
  Note which tests you saw fail first, for the commit message.
- Gateway routes: `make openapi`; the app's API types come from it. In a fresh worktree it skips the app types
  until the app's packages are installed, so a worker runs `make lint s=mobile` once before it.
- Only this task. A sketch detail that doesn't fit the code may be corrected; note it for the commit message.

## 4. Verify

- `make lint s="<services you touched>"` and `make test s="<services you touched>"` (add `contracts` to lint when you
  changed protos). Fix until both pass. In a worktree the tests get their own `<svc>_test_<id>` databases
  automatically; they still need the dev stack (`make dev`) running.
- Go through the acceptance criteria one by one. Prove each with a test. Checks against the running stack
  (`make dev`, `make grpc`, `make smoke`) only work in ROOT: the stack runs the main checkout's code, not a
  worktree's. A worker lists the criteria that need the stack as stack checks, for the main session to run after
  landing. List any criterion that needs a person (a device, an outside account) as a manual check.

## 5. Commit

`git add` only the files this task changed (a worker never adds `docs/EPICs/`; the main session adds the status and
roadmap updates), then commit:

```
feat(<service>): <id> <task title>

<two or three lines on what was built>

Acceptance: <criteria proven by tests or checks>
Tests first: <tests seen failing before their code was written; parts exempt from test-first, or "none">
Manual checks: <criteria left for a person, or "none">
Sketch changes: <deviations from the task's sketches, or "none">
```

A worker that fixes something after committing commits the fix too: the main session squashes a worker's commits into
one, with the first commit's message.

## 6. Workers: report

End with a short summary, then these lines, the last one exactly as shown:

```
Stack checks: <criteria that need the running stack, or "none">
Manual checks: <criteria left for a person, or "none">
DONE <ID>
```

or, when you stopped before committing, `BLOCKED <ID>: <reason>` as the last line.
