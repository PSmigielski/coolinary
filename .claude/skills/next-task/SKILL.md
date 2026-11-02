---
name: next-task
description: Build the next tasks of docs/EPICs/ROADMAP.md end to end - a batch of ready tasks built in parallel by subagents, each in its own git worktree, or the task IDs given - implement, test, mark done and commit each on the milestone branch. Use only when the user runs /next-task or asks to work through the roadmap.
argument-hint: "[N | task-id ...]"
disable-model-invocation: true
---

# Build a batch of roadmap tasks

Build one batch of tasks from `docs/EPICs/`, then stop. Follow `CLAUDE.md`. All commands run from the repo root,
the main checkout (ROOT below). How a single task is built is in `worker.md`, next to this file.

`$ARGUMENTS`:

- nothing: a batch of up to 4 ready tasks; a number N: up to N.
- one or more task IDs: exactly those. Each must be ready, all in one milestone.

A batch of two or more is built in parallel: one subagent (worker) per task, each in its own git worktree (section
A). A batch of one is built here in ROOT, without a worker (section B).

## 1. Check the working tree

Run `git status --porcelain` and check for `.claude/state/current-task.json` (the marker).

- No marker and the tree has changes: stop. Final line: `BLOCKED -: working tree has changes that aren't from a task;
  commit or stash them first`.
- Marker present: a batch was interrupted. Its tasks are the marker's `ids` that are still `in progress`. If
  `.claude/state/landing.json` exists, a land was interrupted: finish or abort it first (A5). With worktrees
  (`scripts/epics/worktree.sh list`), continue at A1: `add` reuses each worktree and its worker continues from what's
  there. A single task without a worktree: continue at B, after reading `git diff`.

## 2. Pick the tasks

- Task IDs given: `python3 scripts/epics/tasks.py show <ID>` for each. If `ready` is false, stop:
  `BLOCKED <ID>: dependencies not done (<ids>)`. Tasks in different milestones: stop with
  `BLOCKED -: <IDs> are in different milestones; run them separately`.
- Otherwise `python3 scripts/epics/tasks.py ready <N>` (N = 4 unless given). If it returns `NO TASK LEFT` or
  `NOTHING READY`, print that as the final line and stop.

Note each task's `id`, `title`, `file` and `service`, and the batch's `branch`, `milestone` and `last_in_milestone`
(from the top level of `ready`, which counts the whole batch; a task's own `last_in_milestone` doesn't).

## 3. Branch and start

- `git switch <branch>` if it exists, otherwise `git switch -c <branch>` from the current HEAD.
- `python3 scripts/epics/tasks.py start <ID> <ID> ...` for the tasks not already `in progress`.

## A. Two or more tasks: build them in parallel

### A1. Worktrees

For each task: `scripts/epics/worktree.sh add <ID>`. It prints the worktree's path: `.claude/worktrees/<ID>` on branch
`task/<ID>`, made from the current HEAD. Its tests get their own databases.

### A2. Workers

Start one worker per task, all in one message so they run at the same time (Agent tool, `general-purpose`), and
wait until every worker has reported before going on. Prompt for each, with the values filled in:

> Build roadmap task `<ID>` (<title>) of the Coolinary repo. Your worktree W is `<worktree path>`; the main checkout
> ROOT is `<ROOT>`. Read `<ROOT>/.claude/skills/next-task/worker.md` and follow it, section 1 included. Task file:
> `<ROOT>/<file>`.

Each worker's last line is `DONE <ID>` (with its stack checks and manual checks above it) or
`BLOCKED <ID>: <reason>`.

### A3. Check the main checkout

`git status --porcelain` must show only the batch's task files, set to `in progress` by step 3. Anything else means
a worker wrote outside its worktree: land nothing and stop with `BLOCKED -: a worker changed <files> in the main
checkout; see git diff`.

### A4. Blocked workers

A worker that stopped on something to decide: in an interactive session, ask the user (batch the questions of all
workers). Record each answer as the task's `Decided (Qnnn)` rule and in `docs/EPICs/DECISIONS.md`, as other
decisions are recorded, then continue that worker with the answer (SendMessage to it) and wait for its report.
Unattended (`claude -p`), or when the user leaves it open: `python3 scripts/epics/tasks.py block <ID> "<reason>"`.
Keep a blocked task's worktree: its work so far is there for whoever unblocks it.

### A5. Land the tasks, one at a time, in roadmap order

For each task whose worker reported `DONE`:

1. `scripts/epics/worktree.sh land <ID>` applies the task's commits to ROOT, staged and not committed. Clean: go
   to 3.
2. On a conflict the script lists each file as generated, lockfile, migration or code. Resolve only what is
   mechanical, then `git add` it:
   - generated: take either side (`git checkout --ours -- <file>`), then regenerate from the merged sources
     (`make proto`, `make sqlc`, `bun run db:generate` in the service, `make openapi`) and add the result.
   - lockfile: take HEAD's (`git checkout --ours -- <file>`), then add the incoming task's dependencies again
     (`uv add` / `bun add` in the service), or relock (`uv lock`, `bun install`) when the manifest merged cleanly.
   - migration: a goose number taken twice: rename the incoming file to the next free number. Drizzle: take HEAD's
     migrations and journal, then `bun run db:generate` for the incoming schema change.
   - code where each side only added its own lines (imports, registry entries, routes, settings, test cases,
     queries in `queries.sql`): keep both.

   Anything that needs a choice (both sides changed the same lines differently, one side renamed or removed what
   the other changed, both changed one contract): `scripts/epics/worktree.sh abort <ID>`, then
   `scripts/epics/worktree.sh remove <ID>` and `python3 scripts/epics/tasks.py set <ID> todo`. It is rebuilt on
   top of the new HEAD by the next batch; note why for the final lines.

   After resolving, run `make lint` and `make test` for the services involved and fix until they pass.
3. `scripts/epics/worktree.sh finish <ID> "<note>"`: marks the task done, updates the roadmap, commits with the
   worker's message and removes the worktree, its branch and its test databases. Pass a note only when you
   changed something while landing ("migration renamed to 00007; registry entries kept from AI-01-08"); it becomes
   a `Landed:` line in the commit.

### A6. Verify the batch

- Services: `python3 scripts/epics/tasks.py changed` (everything the batch touched). Run
  `make lint s="<services>"` and `make test s="<services>"`. Fix until both pass and commit each fix as
  `fix(<service>): <ID> <what>`, naming the task that caused it.
- If a task changed gateway routes, `make openapi` must leave the tree unchanged; commit what it changes as a fix.
- Run the workers' stack checks now that their code is in ROOT (`make dev` reloads it; `make grpc`, `make smoke`).
- When you stop, the Stop hook runs lint and tests on everything the batch touched. If it reports a failure, fix
  it and commit the fix before stopping again.

## B. One task: build it here

1. Build it as `worker.md` sections 2 to 4 describe, in ROOT (skip its worker sections 1 and 6).
2. If something needs deciding: in an interactive session ask the user and record the answer as a
   `Decided (Qnnn)` rule; otherwise run `python3 scripts/epics/tasks.py block <ID> "<what needs deciding, and the
   options>"`, leave your changes uncommitted, and stop with `BLOCKED <ID>: <reason>`.
3. Finish: `python3 scripts/epics/tasks.py done <ID>`, then `scripts/epics/roadmap.sh`. `git add` only the files the
   task changed (plus the `docs/EPICs/` status and roadmap updates) and commit with the message in `worker.md`
   section 5.
4. When you stop, the Stop hook runs lint and tests on every service the task touched. If it reports a failure, fix
   it and commit the fix (`fix(<service>): <ID> ...`) before stopping again.

## Final lines

End your reply with one line per task of the batch, in roadmap order, so `scripts/epics/run-tasks.sh` can read them:

- `DONE <ID>`
- `BLOCKED <ID>: <reason>`
- `REQUEUED <ID>: <why it couldn't land>` (back to `todo`; the next batch rebuilds it)

Then `MILESTONE <Mn> COMPLETE` when `last_in_milestone` was true and every task of the batch is done. When there was
no batch, the only line is `NO TASK LEFT`, `NOTHING READY` or `BLOCKED -: <reason>`.
