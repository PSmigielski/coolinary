---
paths:
  - "**/*.go"
  - "migrations/*.sql"
  - "internal/**/queries.sql"
  - "sqlc.yaml"
---

# Go service patterns

Applies to kitchen, planner and offers. This file is identical in all three; change all three together. Stack:
connect-go (gRPC + Connect on one h2c port), pgx v5, sqlc, goose, go-redis, slog. `<svc>` below means the
service's name.

## Layout

| Path | Holds |
|---|---|
| `cmd/server/main.go` | Wiring only: config, pool, migrations, Redis, clients, `<svc>.New(...)`, servers, shutdown |
| `internal/<svc>/service.go` | The RPC methods (split into more files in the package as it grows) |
| `internal/<svc>/queries.sql` | Every SQL query; `make sqlc` generates `internal/<svc>/store/` |
| `migrations/NNNNN_name.sql` | goose migrations, embedded and run at start-up under an advisory lock |
| `internal/platform/*`, `internal/config` | Shared start-up code, **identical in all three services** (only defaults and addresses differ in `config`) |
| `gen/`, `internal/<svc>/store/` | Generated (`make proto`, `make sqlc`). Never edit |

## RPC methods

- `Service` gets its dependencies through `New(...)`, and the file asserts
  `var _ <svc>v1connect.<Svc>ServiceHandler = (*Service)(nil)`.
- Keep the exported method thin: read `req.Msg` with the getters (`GetEan()`, nil-safe), call an unexported
  method that takes plain values and returns the proto message, and wrap it with `connect.NewResponse`. Test the
  unexported method.
- Errors are always `*connect.Error`; a bare `error` reaches the caller as `unknown`:
  - bad input: `connect.NewError(connect.CodeInvalidArgument, errors.New("ean must have 8, 12, 13 or 14 digits"))`
  - missing: `CodeNotFound`, or `found: false` if the contract has that field (as `ResolveBarcode` does).
    Check `errors.Is(err, pgx.ErrNoRows)`.
  - failing dependency (another service, an external API): `CodeUnavailable`
  - anything else: `connect.NewError(connect.CodeInternal, err)`
  - messages are lowercase and contain no personal data.
- Lists: clamp `page_size` (default and max as in `ListOffers`), fetch `limit+1` rows, and return the last
  returned id as `next_page_token`.
- Pass `ctx` everywhere. Never use `context.Background()` inside a request.
- A new RPC starts in `contracts/coolinary/<svc>/v1/*.proto`, then `make proto`, then the method. Mount new
  services in `main.go` with `rpc.Options(logger)`.

## Calling other services

- Build each client once in `main.go` and pass it to `New`:
  `aiv1connect.NewAiServiceClient(rpc.NewHTTPClient(), rpc.BaseURL(cfg.AIAddr), rpc.ClientOptions()...)`.
- Every call has a deadline: `ctx, cancel := context.WithTimeout(ctx, 2*time.Second); defer cancel()`. The client
  interceptor forwards `x-request-id`, `traceparent`, `x-user-id`, `x-household-id` and `x-plan` from `ctx`.
- The call graph has no cycles (`docs/infrastructure.md`): `ai` and `notifier` call no one, and `planner` calls
  `kitchen` and `offers`, never the other way round.
- **AI output is never trusted.** Check it with rules (allergens, tools at home, stock) before returning it, and
  when `ai` fails or times out, return the non-AI result instead of an error.

## Database

- Schema change = a new file `migrations/0000N_name.sql` with `-- +goose Up` and `-- +goose Down`. Never edit a
  migration that has been applied. Money is integer minor units plus currency; times are `timestamptz`.
- Queries go in `queries.sql` as `-- name: GetX :one` / `:many` / `:exec`, then `make sqlc`. No SQL strings in Go
  code. Writes that must succeed together use `pool.Begin` and `q.WithTx(tx)`.
- Each service reads only its own database. Data from another service comes over gRPC.

## Redis

- `redis-cache` can lose any key at any time: keys are prefixed `<svc>:`, always have a TTL, and cache errors are
  ignored (fall through to the database).
- `redis-state` never evicts and is for quotas, rate limits and pub/sub only.

## Everything else

- Settings: add to `internal/config/config.go` with `get()` / `getInt()` (supports `X_FILE`), and to
  `.env.example`, `.env.dev` and `.env.prod`.
- Logs: `s.log` (slog) with key/value pairs. The interceptor already logs every RPC, so don't log routine success.
  Never log e-mail addresses, free text from users or health data.
- Scheduled jobs (flyer sync, expiry checks) live in the owning service: another `g.Go(...)` in `main.go` that
  stops on `gctx.Done()`.
- Imports in three groups: standard library, third-party, `coolinary/<svc>/...`. Code is `gofmt`ed.
- Go isn't installed on the host: with `make dev` running, check with
  `docker compose -f docker-compose.dev.yml exec <svc> sh -c 'go vet ./... && go test ./...'` from the repo root.
  Tests are table-driven `_test.go` files next to the code.
