# Coolinary — instructions for Claude

Fridge-to-meals planner: Go, Bun (Elysia) and Python services talking gRPC, an Expo app, Docker Compose.
Start with `README.md`. Product rules: `docs/features.md`, `docs/business-model.md`. How things are built (binding):
`docs/infrastructure.md`, `docs/service-conventions.md`.

## Building the planned work

- All planned work is in `docs/EPICs/`: one file per task, built in the order of `docs/EPICs/ROADMAP.md`.
- Work on one task at a time with `/next-task` (unattended: `scripts/epics/run-tasks.sh`). Never start a task whose
  dependencies (the roadmap's "After" column) aren't done.
- The task file is the spec: business rules, `Decided (Qnnn):` rules, contract and data sketches, errors and
  acceptance criteria. `docs/EPICs/README.md` holds the decisions and the cross-service names (RPCs, events, error
  reasons) that bind every task; `docs/EPICs/DECISIONS.md` indexes every decision.
- Don't make design decisions. If a task is unclear, contradicts the code or another task, or needs a choice no
  `Decided` rule covers: in an interactive session ask the user; unattended, block the task
  (`python3 scripts/epics/tasks.py block <ID> "<what needs deciding>"`) and stop.
- Stay inside the task: don't build later tasks or refactor unrelated code. Small corrections to a sketch (a field
  name, an index) are fine when the code needs them; say so in the commit message.
- Task status changes only through `scripts/epics/tasks.py`. After changing task files or dependencies, run
  `scripts/epics/roadmap.sh`.

## Rules for code

- Contracts: edit only `contracts/coolinary/<svc>/v1/*.proto` (shared types in `common/v1`), then `make proto`. Never
  edit generated code (`gen/`, `src/gen/`, `ai/src/coolinary/`) by hand; commit it with the change.
- Go: SQL in `queries.sql`, then `make sqlc`; goose migrations in `migrations/`. Bun: Drizzle schema, then
  `bun run db:generate`. Gateway route changes: `make openapi` (also regenerates the app's API types).
- Follow each service's `.claude/rules/*.md`, and `mobile/AGENTS.md` for the app (check the Expo SDK 58 docs, never
  rely on memory).
- No personal data, user free text or health data in logs, traces or events. Age, height and weight never leave the
  phone; the server sees only a postcode, never a position.
- Tests: write them first (`.claude/rules/testing.md`); what to test is "Tests" in `docs/service-conventions.md`. A
  task is done only when `make lint s=<svc>` and `make test s=<svc>` pass for every service it touched (a Stop
  hook checks this during `/next-task`).
- Never call real outside services in tests or by hand (Claude API, Open Food Facts, retailer sites, Expo Push).

## Commands

`make dev` · `make ps` · `make logs s=<svc>` · `make lint s="<svc> ..."` · `make test s="<svc> ..."` ·
`make proto` · `make sqlc` · `make openapi` · `make psql s=<svc>` · `make grpc s=<svc> args=list` · `make smoke`

## Git

- One branch per milestone (`m01-foundations`, … as `tasks.py next` names it), one commit per finished task:
  `feat(<svc>): <TASK-ID> <title>`. Never push, merge, rebase or reset; the owner reviews and merges each milestone.
- Commit only the files the task changed. Other work may be staged or modified in the tree: leave it alone.
