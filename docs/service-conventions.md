# Service conventions

Every service is its own project (own lockfile, Dockerfile and generated gRPC code), so the small amount of
start-up code each one needs is written per service. These rules keep the copies identical. When you add a
service, copy the closest existing one (`kitchen` for Go, `notifier` for Bun, `ai` for Python) and keep to them.

## Ports inside the container

| Port | Use |
|---|---|
| `50051` | gRPC (every service except `gateway`) |
| `8081` | Admin HTTP: `GET /metrics` (Prometheus), `GET /healthz` (Docker healthcheck) |
| `3000` | Public HTTP (`gateway`), Better Auth HTTP (`identity`, reached only through the gateway) |
| `2345` / `6499` / `5678` | Debuggers in dev: Delve (Go), Bun inspector, debugpy (Python) |

`/healthz` returns `200 {"ok":true,"checks":{...}}` only when the service's own dependencies (Postgres, Redis)
answer, and `503` otherwise. It never checks other services: a service is healthy even when the services it calls
are down.

## Environment variables

| Variable | Meaning |
|---|---|
| `SERVICE_NAME` | Used in logs and metrics |
| `GRPC_PORT`, `ADMIN_PORT`, `HTTP_PORT` | See above |
| `LOG_LEVEL` | `debug`, `info`, `warn`, `error` |
| `GRPC_REFLECTION` | `true` only in dev (for `grpcurl`) |
| `DATABASE_URL` | The service's own database; no service reads another's |
| `REDIS_CACHE_URL` | `redis-cache`: may evict anything at any time |
| `REDIS_STATE_URL` | `redis-state`: quotas, rate limits, pub/sub; never evicts |
| `<SVC>_GRPC_ADDR` | Address of a service this one calls, e.g. `KITCHEN_GRPC_ADDR=kitchen:50051` |

Any variable `X` can instead be given as `X_FILE=/path`, and the file's content is used. Prod uses this for Docker
secrets (`DATABASE_URL_FILE=/run/secrets/kitchen_database_url`).

Each service has three env files: `.env.example` (documents every variable), `.env.dev` (dev values, committed),
`.env.prod` (prod values without secrets, committed).

## gRPC

- Contracts live only in `contracts/` (packages `coolinary.<service>.v1`, buf STANDARD lint). `make proto`
  regenerates every service's code; generated code is committed.
- Every service serves `grpc.health.v1.Health`. Reflection is on only when `GRPC_REFLECTION=true`.
- Every call sets a deadline. Callers forward these metadata keys when they have them:
  `x-request-id`, `traceparent`, `x-user-id`, `x-household-id`, `x-plan` (`free` / `plus`).
- Errors use gRPC status codes: `INVALID_ARGUMENT` for bad input, `NOT_FOUND`, `UNAVAILABLE` for a failing
  dependency, `INTERNAL` for bugs.
- Bun services use unary RPCs only for now (Bun's HTTP/2 has open bugs with server streams; see
  `docs/infrastructure.md`).

## Metrics

Same names in Go, Bun and Python:

| Metric | Labels |
|---|---|
| `rpc_server_requests_total` (counter) | `rpc_service`, `rpc_method`, `rpc_code` |
| `rpc_server_duration_seconds` (histogram) | `rpc_service`, `rpc_method` |
| `rpc_client_requests_total` (counter) | `rpc_service`, `rpc_method`, `rpc_code` |
| `http_server_requests_total`, `http_server_duration_seconds` | gateway only: `method`, `route`, `status` |

`rpc_code` is the lowercase status name: `ok`, `not_found`, `invalid_argument`, `deadline_exceeded`, ...
`rpc_service` is the full gRPC service name (`coolinary.kitchen.v1.KitchenService`). Prometheus adds a `service`
label (`kitchen`) per scrape target; don't set one yourself.

## Logs

JSON lines on stdout: `time`, `level`, `service`, `msg`, plus fields (`request_id`, `method`, `code`,
`duration_ms`, `error`). No personal data in logs: no e-mail addresses, no free text from users, no health data.

## Database migrations

Each service owns its migrations and runs them at start-up under a Postgres advisory lock:
Go uses `goose` (SQL files in `migrations/`, embedded) with `sqlc` for queries (`make sqlc`); Bun uses Drizzle
(`bun run db:generate` writes `drizzle/`). The Postgres init script (`platform/postgres/initdb/`) creates each
service's database and role only when the Postgres volume is first created.
