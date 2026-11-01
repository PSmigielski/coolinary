# kitchen

Go. Ingredient catalogue, barcode lookup chain (household → shared catalogue → Open Food Facts → GS1), nutrition, fridge stock, recipe library.

- gRPC `50051` (connect-go: gRPC and Connect/HTTP on one port), admin `8081`, Delve `2345` in dev.
- RPC: `KitchenService.ResolveBarcode` (try `9520000000017`).
- Migrations: `migrations/*.sql` (goose, embedded, run at start-up). Queries: `internal/kitchen/queries.sql` → `make sqlc`.
- `go run ./cmd/smoketest [-idle 600s]`: Go → Bun gRPC smoke test (`make smoke`).

Runs with the rest of the stack: `make dev` from the repo root (hot reload, debugger). Settings: `.env.example`
(all variables), `.env.dev` (dev), `.env.prod` (prod, secrets via `*_FILE`). Shared rules:
`../docs/service-conventions.md`. Generated gRPC code is committed; regenerate with `make proto`.
