# planner

Go. Dish suggestions ("what can I make now"), day and week plans, per-person portions, rule checks on AI-composed dishes, shopping list, shop-trip plan.

- gRPC `50051` (connect-go), admin `8081`, Delve `2345` in dev.
- RPC: `PlannerService.SuggestDishes`.
- Admin RPC: `PlannerAdminService` (`ListAdminChanges`, `GetQueueCounts`), same port, called only by admin and
  refused without `x-admin-user-id`; the gateway has no client for it. Admin RPCs that change data call
  `recordAdminChange` in their transaction (`admin_changes`: admin, RPC, target such as `dish:<id>`, a JSON summary
  of ids and counts; kept 1 year). Try `make grpc s=planner d='{}' flags="-H x-admin-user-id:a1"
  args=coolinary.planner.v1.PlannerAdminService/GetQueueCounts`.
- Migrations: `migrations/*.sql` (goose). Queries: `internal/planner/queries.sql` → `make sqlc`.
- Events: `events.Consumer` (`internal/platform/events`) reads the streams in `EVENT_STREAMS` (group `planner`) and
  runs the handlers registered with `consumer.Handle(type, version, fn)` once per event (`processed_events`). Entries
  it gives up on (5 failed deliveries, a newer version, a malformed envelope) go to `dead_events`. Planner publishes
  no events.
- Calls: `clients.New` (`internal/clients`) builds the kitchen, offers, ai, notifier and identity clients. Each call
  gets a deadline (the per-service default, or the caller's deadline minus 200 ms), kitchen and offers reads retry
  once after 100 ms, and `clients.KindOf(err)` sorts errors into `Unavailable`, `Rejected` and `Bug`. Background
  jobs call with `clients.JobContext(ctx, plan)`.

Runs with the rest of the stack: `make dev` from the repo root (hot reload, debugger). Settings: `.env.example`
(all variables), `.env.dev` (dev), `.env.prod` (prod, secrets via `*_FILE`). Shared rules:
`../docs/service-conventions.md`. Generated gRPC code is committed; regenerate with `make proto`.
