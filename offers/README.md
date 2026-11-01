# offers

Go. Shops, flyer sync, offers and prices, flyer editions by area. Distances are computed on the phone; the server only sees a postcode.

- gRPC `50051` (connect-go), admin `8081`, Delve `2345` in dev.
- RPC: `OffersService.ListOffers`.
- Migrations: `migrations/*.sql` (goose). Queries: `internal/offers/queries.sql` → `make sqlc`.

Runs with the rest of the stack: `make dev` from the repo root (hot reload, debugger). Settings: `.env.example`
(all variables), `.env.dev` (dev), `.env.prod` (prod, secrets via `*_FILE`). Shared rules:
`../docs/service-conventions.md`. Generated gRPC code is committed; regenerate with `make proto`.
