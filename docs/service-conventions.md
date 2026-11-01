# Service conventions

Every service is its own project (own lockfile, Dockerfile and generated gRPC code), so the small amount of
start-up code each one needs is written per service. These rules keep the copies identical. When you add a
service, copy the closest existing one (`kitchen` for Go, `notifier` for Bun, `ai` for Python) and keep to them.

## Ports inside the container

| Port | Use |
|---|---|
| `50051` | gRPC (every service except `gateway`) |
| `8081` | Admin HTTP: `GET /metrics` (Prometheus), `GET /healthz` (Docker healthcheck) |
| `3000` | Public HTTP (`gateway`), Better Auth HTTP (`identity`, reached only through the gateway), back-office HTTP (`admin`, reached only through Caddy at `admin.<domain>`) |
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
| `REDIS_STATE_URL` | `redis-state`: quotas, rate limits, idempotency keys, event streams, live notices, AI spend counters; never evicts |
| `<SVC>_GRPC_ADDR` | Address of a service this one calls, e.g. `KITCHEN_GRPC_ADDR=kitchen:50051` |
| `ANTHROPIC_API_KEY` | `ai` only: the Claude API key (prod: `ANTHROPIC_API_KEY_FILE`) |

Any variable `X` can instead be given as `X_FILE=/path`, and the file's content is used. Prod uses this for Docker
secrets (`DATABASE_URL_FILE=/run/secrets/kitchen_database_url`).

Each service has three env files: `.env.example` (documents every variable), `.env.dev` (dev values, committed),
`.env.prod` (prod values without secrets, committed).

## gRPC

- Contracts live only in `contracts/` (packages `coolinary.<service>.v1`, buf STANDARD lint). `make proto`
  regenerates every service's code; generated code is committed.
- Every service serves `grpc.health.v1.Health`. Reflection is on only when `GRPC_REFLECTION=true`.
- Every call sets a deadline. Callers forward these metadata keys when they have them:
  `x-request-id`, `traceparent`, `x-user-id`, `x-household-id`, `x-plan` (`free` / `plus`), `x-language`
  (BCP 47, set by the gateway from `Accept-Language`). Calls from the admin service add `x-admin-user-id`.
- User-scoped requests carry `household_id` (and `member_id` where needed) as message fields. The gateway fills them
  from the JWT, never from client input; services trust them and pass them on in service-to-service calls.
- Plus rules: a service that needs a household's plan asks identity (`IdentityService.GetEntitlements`) on that
  request, memoised for the request only; background jobs ask per household per run. Never keep a local plan copy.
  If identity can't be reached, a request falls back to the token's `x-plan` (Free when missing) and a job skips
  the household until its next run.
- Each service has one public gRPC service (`KitchenService`) and one admin service (`KitchenAdminService` in
  `contracts/coolinary/<svc>/v1/admin.proto`) that only `admin` calls. Messages may be split into per-domain files
  in the same package.
- Types used by more than one service (`Cuisine`, `KitchenTool`, `Diet`, `Allergen`, `EffortTag`, `Unit`, `Money`,
  `Quantity`, paging) live in `contracts/coolinary/common/v1/`. Never redefine them per service.
- Errors use gRPC status codes:

  | Code | Use |
  |---|---|
  | `INVALID_ARGUMENT` | bad input |
  | `NOT_FOUND` | the thing doesn't exist (or isn't visible to this household) |
  | `PERMISSION_DENIED` | wrong household or role |
  | `FAILED_PRECONDITION` | a business rule says no, e.g. Plus required |
  | `ALREADY_EXISTS` | duplicate create |
  | `ABORTED` | version conflict on shared data (stock, shopping list); the client refetches and retries |
  | `RESOURCE_EXHAUSTED` | a quota, rate limit or the AI spend limit |
  | `UNAVAILABLE` | a failing dependency |
  | `INTERNAL` | bugs |

  The machine-readable reason goes in trailing metadata `x-error-code` as a stable snake_case value
  (`plus_required`, `quota_reached`, `ai_refused`, `version_mismatch`, ...). The gateway maps the status and the
  reason into its HTTP error body. Never put personal data in error messages.
- Bun services use unary RPCs only for now (Bun's HTTP/2 has open bugs with server streams; see
  `docs/infrastructure.md`).

## Events

Changes other services must learn about travel as events on Redis Streams in `redis-state`, never as extra calls.

- **Producing.** Write the event to the service's `outbox` table in the same transaction as the change. A relay
  publishes outbox rows in order to the service's stream `events.<service>` with
  `XADD events.<service> MINID ~ <now − 7 days> * type <type> envelope <json>`, then marks them sent. A crash between
  the two publishes twice; that is allowed. (Decision Q015 in `docs/EPICs/DECISIONS.md`.)
- **Envelope** (JSON, snake_case): `id` (UUIDv7, the idempotency key), `type` (e.g. `household.member_added`),
  `version` (starts at 1), `occurred_at` (RFC 3339 UTC), `producer`, `household_id` (empty for catalogue-wide
  events), `payload`. Payloads hold ids, codes, dates and counts only: no names, e-mail addresses, free text or health
  data. At most 64 KB; split larger fan-outs.
- **Versioning.** Adding an optional field keeps the version. Removing, renaming or changing the meaning of a field
  raises it; the producer emits the new version only after every consumer accepts both.
- **Consuming.** One consumer group per consuming service (named after it). Handlers are idempotent: record each
  processed `id` and skip duplicates; treat "already exists" upserts as success. Ignore unknown types and unknown
  fields. Acknowledge only after the handler's transaction commits. Failed events are retried, then parked for a
  person to look at (counted in metrics).
- **Falling behind.** Streams keep about 7 days. A consumer that misses events asks the owner to republish snapshots
  (`household.snapshot` from identity) instead of guessing.
- **Live notices** are not events: Redis pub/sub on `live.household.<id>` with
  `{"kind": "...", "household_id": "...", "version": n}`, published after commit, best effort. The gateway relays
  them over its WebSocket; clients refetch when a version is skipped.

The event types, payloads and consumers are listed in `docs/EPICs/README.md` ("Events").

## Metrics

Same names in Go, Bun and Python:

| Metric | Labels |
|---|---|
| `rpc_server_requests_total` (counter) | `rpc_service`, `rpc_method`, `rpc_code` |
| `rpc_server_duration_seconds` (histogram) | `rpc_service`, `rpc_method` |
| `rpc_client_requests_total` (counter) | `rpc_service`, `rpc_method`, `rpc_code` |
| `http_server_requests_total`, `http_server_duration_seconds` | HTTP services (gateway, admin): `method`, `route`, `status` |

`rpc_code` is the lowercase status name: `ok`, `not_found`, `invalid_argument`, `deadline_exceeded`, ...
`rpc_service` is the full gRPC service name (`coolinary.kitchen.v1.KitchenService`). Prometheus adds a `service`
label (`kitchen`) per scrape target; don't set one yourself.

## Logs

JSON lines on stdout: `time`, `level`, `service`, `msg`, plus fields (`request_id`, `method`, `code`,
`duration_ms`, `error`). No personal data in logs: no e-mail addresses, no free text from users, no health data.

## Linting

`make lint s=<svc>` must pass before a push (the `pre-push` hook from `make hooks` runs it). Each service keeps its
lint config next to its code, identical across services of the same language: `biome.json` (Bun),
`.golangci.yml` (Go), `[tool.ruff]` / `[tool.mypy]` in `pyproject.toml` (Python). Change all copies together. A
`//nolint` or `# noqa` names the rule and says why on the same line or the one above.

## Tests

`make test s=<svc>` runs a service's tests (`scripts/test-service.sh`): `go test ./...` in a `golang` container for
Go, `bun test` for Bun, `pytest` for `ai`. A task is not done until `make lint s=<svc>` and `make test s=<svc>` pass.
Tests are written first (red, green, refactor); `.claude/rules/testing.md` has the workflow and what a good test
looks like.

- Every business rule and acceptance criterion in a task gets a test where it can be tested in code. Name the test
  after the rule (`TestConsumeStock_OldestExpiryFirst`, `"rejects a target under 1,200 kcal"`).
- Unit tests need nothing running. Tests that need Postgres or Redis read `TEST_DATABASE_URL`,
  `TEST_REDIS_CACHE_URL` and `TEST_REDIS_STATE_URL` and **skip** when they are unset (`t.Skip`, `test.skipIf`,
  `pytest.mark.skipif`). `make test` sets them only while `make dev` is up, pointing at a `<svc>_test` database
  (created on first use) and Redis database 15, never at dev data. Each test cleans up what it creates.
- Tests never call the real Claude API, Open Food Facts, retailer sites or Expo Push; use the `stub` LLM provider
  and recorded fixtures.
- Go: table tests, `_test.go` next to the code. Bun: `*.test.ts` next to the code. Python: `ai/tests/test_*.py`.

## Database migrations

Each service owns its migrations and runs them at start-up under a Postgres advisory lock:
Go uses `goose` (SQL files in `migrations/`, embedded) with `sqlc` for queries (`make sqlc`); Bun uses Drizzle
(`bun run db:generate` writes `drizzle/`). The Postgres init script (`platform/postgres/initdb/`) creates each
service's database and role only when the Postgres volume is first created.
