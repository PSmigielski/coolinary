# planner

Go. Dish suggestions ("what can I make now"), day and week plans, per-person portions, rule checks on AI-composed dishes, shopping list, shop-trip plan.

- gRPC `50051` (connect-go), admin `8081`, Delve `2345` in dev.
- RPC: `PlannerService.SuggestDishes`.
- Migrations: `migrations/*.sql` (goose). Queries: `internal/planner/queries.sql` → `make sqlc`.
- Events: `events.Consumer` (`internal/platform/events`) reads the streams in `EVENT_STREAMS` (group `planner`) and
  runs the handlers registered with `consumer.Handle(type, version, fn)` once per event (`processed_events`). Entries
  it gives up on (5 failed deliveries, a newer version, a malformed envelope) go to `dead_events`. Planner publishes
  no events.

Runs with the rest of the stack: `make dev` from the repo root (hot reload, debugger). Settings: `.env.example`
(all variables), `.env.dev` (dev), `.env.prod` (prod, secrets via `*_FILE`). Shared rules:
`../docs/service-conventions.md`. Generated gRPC code is committed; regenerate with `make proto`.
