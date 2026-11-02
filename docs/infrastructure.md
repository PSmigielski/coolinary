# Coolinary — infrastructure

Scope: how the back end, the app and their tooling are set up and run. Product decisions are in `business-model.md`
and `features.md`; per-service rules are in `service-conventions.md`. Day-to-day commands are in the root `README.md`.

## Decisions

- **Microservices in Go and Bun (Elysia), talking gRPC**; one Python service (uv + FastAPI) for every LLM call.
- **PostgreSQL 18** (one instance, one database per service) and **Redis 8** (cache and state kept apart).
- **Prometheus + Grafana** for resource and service metrics, **Langfuse** (self-hosted) for LLM traces and cost.
- **Mailpit** catches dev e-mail (MailHog has had no release since 2020; Mailpit is a drop-in, same ports).
- **Docker Compose** from three files: `docker-compose.abstract.yml` (shared, static IPs, no published ports),
  `docker-compose.dev.yml` (hot reload, debuggers) and `docker-compose.prod.yml` (hardened, Caddy in front).
  Every image we build is multi-stage.
- **Expo** app (TanStack Query, Zustand), outside Docker.
- **Folder layout: self-contained services.** Each service has its own lockfile, Dockerfile and generated gRPC code;
  only `contracts/` is shared, so any service can move to its own repository.
- **Zed** is the debugger front end.
- **Events between services: Redis Streams** on `redis-state`, written through a transactional outbox in each producer
  (identity, kitchen, offers). No message broker container.
- **LLM provider: the Claude API** through the official `anthropic` Python SDK, only from `ai`. Claude Haiku 5.5 reads
  and matches (text, photos, labels, receipts, flyers, catalogue matching, translation); Claude Sonnet 5.5 composes
  for Plus households when it beats Haiku on the golden test set; Free day plans compose on Haiku. Bulk jobs use
  Message Batches. The Claude API cannot pin inference to the EU, so the privacy policy says inference may run
  outside the EU.
- **A separate admin (back-office) service** in Bun + Elysia, reached at `admin.<domain>` by admin accounts only.

The planned build of all of this is in `docs/EPICs/` (decisions in `docs/EPICs/README.md`, order in
`docs/EPICs/ROADMAP.md`). Items below marked *(planned)* are designed there but not built yet.

Constraints from the business model that shaped this: one founder, fixed costs ≤ €200/month before revenue (few,
coarse services; one VPS running Compose); the app must work with AI off; AI cost caps per user; EU hosting and
health data kept on the phone.

## Services

| Service | Runtime | Owns | Calls |
|---|---|---|---|
| `gateway` | Bun + Elysia | Public HTTP API for the app, JWT check, OpenAPI, `/api/auth/*` proxy to identity, Free-tier AI quota and rate limits, live shopping-list WebSocket *(planned)* | identity, notifier, kitchen, planner, offers (health only to ai) |
| `identity` | Bun + Elysia + Better Auth | Accounts (sign-in by e-mail code, anonymous accounts), sessions, JWTs (JWKS), households, members and profiles without accounts, invites, plan and entitlements, admin role | notifier |
| `notifier` | Bun | E-mail (SMTP), push (Expo Push API) *(planned)*, devices and notification preferences, founder alerts, delivery log | — |
| `kitchen` | Go | Ingredient catalogue, name matching, barcode lookup, nutrition tables and calculation, fridge stock and AI intake, expiry reminders, the kitchen profile with per-person kcal targets, savings ledger | identity, ai, notifier |
| `planner` | Go | Recipe library (dishes, import and review, calories per dish), dish suggestions, rule checks on AI output, day/family/week plans, portions, cooking mode, shopping list, offer badges, trip price matrix | identity, kitchen, offers, ai, notifier |
| `offers` | Go | Shops (OpenStreetMap), postcode areas and flyer editions, flyer sync, offers, prices | identity, kitchen, ai, notifier |
| `ai` | Python 3.14 + uv + FastAPI | Every LLM call (Claude), structured output, batch jobs, Langfuse tracing, result cache, spend limit | — |
| `admin` *(planned)* | Bun + Elysia | Back office: review queues, imports, manual plan switch, audit log | identity, notifier, kitchen, planner, offers, ai |

- Service to service is gRPC only. The app talks HTTP to the gateway; the gateway proxies Better Auth's HTTP routes.
- The call graph has no cycles; `ai` and `notifier` call no one. kitchen, planner and offers read a household's plan from identity
  (`GetEntitlements`) when they need it, instead of keeping a copy. Each service also exposes a `<Svc>AdminService` that
  only `admin` calls.
- Changes that several services must learn about (members joining or leaving, plan changes, account deletion,
  catalogue merges, new offers) travel as events, not calls; see "Events" below.
- The app never calls `ai` directly. AI features go through the Go service that owns them, which checks the AI output
  with rules (tools, allergens, stock) before anything is shown, and falls back to non-AI results when `ai` is down.
- Each service runs as a single container (fixed IP). Scheduled jobs live in the owning service.

What exists today is a skeleton: each service boots, migrates its database, passes health checks, serves one real
RPC (`identity.GetHousehold`, `kitchen.ResolveBarcode`, `planner.SuggestDishes`, `offers.ListOffers`,
`notifier.SendEmail`, `ai.ExtractItemsFromText`) and is wired into monitoring. The gateway serves
`GET /v1/system/status`, `GET /v1/products/{ean}`, `GET /v1/me` and `GET /v1/household`. Sign-up through
`/api/auth/sign-up/email` creates the user's own household; `/api/auth/token` returns the JWT the gateway accepts.

## Repository layout

```
contracts/              shared .proto files (buf), OpenAPI export of the gateway, code-generator image
gateway/ identity/ notifier/     Bun services (src/gen = generated, committed)
kitchen/ planner/ offers/        Go services (gen/ = generated, committed)
ai/                     Python service (src/coolinary = generated, committed)
mobile/                 Expo app
platform/               config for third-party containers (postgres, redis, prometheus, grafana, langfuse, caddy)
docs/                   product and infrastructure docs
.zed/                   Zed debugger and task configs
scripts/make-secrets.sh prod secrets generator
docker-compose.{abstract,dev,prod}.yml, Makefile, .env (versions and network prefix)
```

## gRPC and code generation

- `contracts/` is one buf module, buf STANDARD lint, FILE-level breaking-change checks (`make proto-breaking`
  compares against `main`).
- `make proto` builds `coolinary/protogen` from `contracts/Dockerfile`: buf, protoc-gen-go, protoc-gen-connect-go,
  ts-proto (run by Bun) and grpcio-tools, all pinned. Generation needs no network and no Buf account (the hosted
  plugins rate-limit anonymous use). Output goes to `.gen/` first and is then synced into each service, so running
  dev containers never see a half-empty folder.
- **Go:** connect-go (v1 API) on plain `net/http` with h2c. One port serves standard gRPC to grpc-js and grpcio
  clients, and the Connect protocol over HTTP/1.1, so a Go service can be called with curl:
  `curl -H 'Content-Type: application/json' -d '{"ean":"9520000000017"}' http://10.77.0.30:50051/coolinary.kitchen.v1.KitchenService/ResolveBarcode`.
- **Bun:** ts-proto + `@grpc/grpc-js` (the gRPC library Bun tests against), a small built-in `grpc.health.v1` service
  (the `grpc-health-check` package reads `.proto` files at runtime, which the bundled prod image doesn't have), and
  `@grpc/reflection` in dev only, fed by a descriptor set that `make proto` writes (it includes the standard
  `contracts/grpc/health/v1/health.proto`, which is used for nothing else).
- **Python:** grpcio-tools at the same version as the `grpcio` runtime (1.85.1), `grpc.aio`.

**Bun HTTP/2 risk.** Bun has open HTTP/2 bugs (bun#45728 server streams keep accepting writes after the client
leaves; bun#31881 client hangs after idle; bun#39796 client stalls under load). So: Bun is pinned, Bun services use
unary RPCs only, clients send keepalive pings and set deadlines, and `make smoke` checks both directions (unary call,
deadline, cancel, and with `idle=600` a call after 10 idle minutes). Fallback if it ever fails: Bun → Go calls switch
to the Connect protocol over Bun's `fetch`, which the Go servers already accept, with the same contracts.

## Data

**Postgres 18** (`postgres:18`): databases and roles `identity`, `notifier`, `kitchen`, `planner`, `offers`,
`langfuse`, plus `postgres_exporter` (role `pg_monitor`). Created by `platform/postgres/initdb/01-databases.sh` on the
first start of an empty volume (dev passwords `<name>-dev`, prod from secrets). Data lives under
`/var/lib/postgresql/18/docker`; the volume mounts at `/var/lib/postgresql`. The admin service adds an `admin`
database *(planned, ADM-01)*; an existing volume needs it created by hand (`make psql s=postgres`). PostGIS isn't
needed (distances are computed on the phone). pgvector isn't needed either: names are matched to the catalogue with
`pg_trgm` candidates and a Claude Haiku pick, stored once for everyone.

**Redis 8**, three instances because caching and state need opposite eviction rules:

| Instance | Policy | Used for |
|---|---|---|
| `redis-cache` | `allkeys-lru`, no persistence | Barcode/catalogue caches, AI result cache, key prefix per service |
| `redis-state` | `noeviction`, AOF | Free-tier AI quota, rate limits, idempotency keys, event streams (`events.<producer>`), live notices (pub/sub `live.household.<id>`), AI spend counters |
| `redis-langfuse` | `noeviction` (required by Langfuse) | Langfuse queues only (observability profile) |

## Events *(planned)*

- Producers: identity, kitchen and offers. Each writes events to an outbox table in the same transaction as the
  change; a relay publishes them to its stream (`events.identity`, `events.kitchen`, `events.offers`). Planner
  publishes no events.
- Consumers read with Redis consumer groups, at least once, and drop duplicates by the event id. Streams keep about
  7 days; a consumer that falls further behind resyncs from snapshots the owner republishes.
- Events carry ids, codes and dates only, never names, e-mail addresses, free text or health data.
- Live notices for the shared shopping list and fridge are separate: plain pub/sub on `live.household.<id>`, relayed
  by the gateway WebSocket. Clients refetch when they miss one.
- The event types, payloads and consumers are listed in `docs/EPICs/README.md` ("Events"); the conventions are in
  `docs/service-conventions.md`.

## Network

One bridge network named `coolinary`, subnet `10.77.0.0/24` (`NET_PREFIX` in `.env`; Docker's own pool and this
machine's other projects use 172.16/12). Every container has a fixed address; one-off containers get addresses from
`10.77.0.128/25`. The base file publishes no ports.

| Address | Container | Address | Container |
|---|---|---|---|
| `.2` | caddy (prod) | `.40` | ai |
| `.3` | mailpit (dev) | `.50` | prometheus |
| `.10` | postgres | `.51` | grafana |
| `.11` | redis-cache | `.52` | cadvisor |
| `.12` | redis-state | `.53` | node-exporter (prod) |
| `.13` | pg-backup (prod) | `.54` | postgres-exporter |
| `.20` | gateway | `.55` | redis-exporter |
| `.21` | identity | `.60` | langfuse-web |
| `.22` | notifier | `.61` | langfuse-worker |
| `.23` | admin *(planned)* | | |
| `.30` | kitchen | `.62` | clickhouse |
| `.31` | planner | `.63` | seaweedfs |
| `.32` | offers | `.64` | redis-langfuse |

- **Dev:** the Linux host reaches every container by IP (native rootful Docker). Published: the gateway on
  `GATEWAY_HOST_PORT` (3080; port 3000 is taken on this machine) for phones on the LAN, and the Bun inspector ports
  on 127.0.0.1 only.
- **Prod:** only Caddy publishes `80`, `443/tcp`, `443/udp` with automatic TLS, routing `api.<domain>` → gateway,
  `grafana.<domain>` → Grafana, `langfuse.<domain>` → Langfuse, and `admin.<domain>` → admin *(planned)*.
- Configs use DNS names (`kitchen:50051`); the fixed IPs are for host access, Zed and Prometheus.
- With fixed IPs, `docker compose run <svc>` would clash with the running container; the Makefile uses `exec`.

## Compose files

`docker-compose.dev.yml` and `docker-compose.prod.yml` `include` the base file and redefine the services they change;
Compose merges the two and the local definition wins. This needs **Docker Compose ≥ 5.0** (compose-go PR #835; the
Docker docs page still describes the old error). If it ever misbehaves, `docker compose -f docker-compose.abstract.yml
-f docker-compose.dev.yml` gives the same result. `make dev` / `make prod` are the only entry points.

- **Base:** network, volumes, every service with build context, fixed IP and `depends_on`. Healthcheck commands live
  in the images (dev images have a shell, prod images don't); the base sets only timing. Services wait only for
  Postgres and Redis, never for each other or for `ai`. Monitoring and Langfuse are in the `observability` profile
  and nothing outside it depends on it.
- **Dev:** `target: dev`, source bind-mounted at `/app`, a watcher as the command, cache volumes, debug ports,
  Mailpit, Grafana with anonymous viewing. Containers write nothing into the repo (air builds to `/tmp/air`, the Python
  venv is in `/opt/venv`, `node_modules` is an anonymous volume refreshed by `make dev`).
- **Prod:** `target: prod`, `restart: unless-stopped`, resource limits, Docker secrets from `./secrets`
  (`make secrets`), Caddy, node-exporter, `pg-backup`. Our seven images run non-root with a read-only root
  filesystem, `cap_drop: ALL` and `no-new-privileges`; third-party images keep the capabilities their entrypoints
  need (cAdvisor is privileged, node-exporter uses the host PID namespace). Images are built on the VPS at first
  (`make prod`); a registry and CI come later.

## Hot reload and debugging

| Runtime | Reload | Debugger |
|---|---|---|
| Go | `air` rebuilds with `-gcflags "all=-N -l"` and restarts under headless Delve on `:2345` | Zed: `.zed/debug.json` (attach, `mode: remote`, `substitutePath` to `/app`) |
| Python | `watchfiles` restarts `python -m debugpy --listen 0.0.0.0:5678 -m coolinary_ai` | Zed: `.zed/debug.json` (attach, `pathMappings` to `/app`) |
| Bun | `bun --watch` | Zed can't debug Bun yet (zed#49033, bun#22914). `make debug s=<svc>` restarts the service under the inspector (without `--watch`, which leaves the inspector port busy: bun#45200) and you open `https://debug.bun.sh/#localhost:<6499/6500/6501>/<svc>`. |

Every reload restarts the debug server, so re-attach after a change. Both the Delve and the debugpy attach were
checked with a DAP client sending what Zed sends.

## Images

| Image | Stages | Final |
|---|---|---|
| Go services | base → dev (air, delve) / build → prod | `gcr.io/distroless/static:nonroot`, ~31 MB; healthcheck = the binary's `healthcheck` subcommand |
| Bun services | base → deps → dev / build (typecheck + bundle) → prod | `oven/bun:<v>-distroless`, uid 65532, ~195 MB |
| ai | uv → base → deps → dev / build → prod | `python:3.14-slim`, uid 65532, venv only, ~270 MB |
| protogen | go-plugins, ts-plugin, buf, bun → final | code generators for `make proto` |

Third-party containers use pinned upstream images (`.env`).

## Monitoring

- **Prometheus** scrapes every service's `/metrics` (same metric names in all languages), cAdvisor (per container),
  postgres-exporter, redis-exporter (all three Redis instances) and, in prod, node-exporter. 15 days of data in dev,
  30 in prod.
- **Grafana** is provisioned from `platform/grafana`: dashboards *Coolinary Containers*, *Coolinary Services* and
  *Coolinary Data Stores*.
- **Langfuse v4** (web, worker, ClickHouse, SeaweedFS as the S3 store, its own Redis, a database on the shared
  Postgres). MinIO's community edition was archived in Feb 2026, hence SeaweedFS. First boot creates the org,
  project and API keys from `LANGFUSE_INIT_*`. The `ai` service traces each LLM call as a generation with the
  feature and plan as tags, a salted hash of the household id as user, token usage, and only sizes as input/output;
  a mask function redacts any string that slips through. Tracing switches itself off when Langfuse isn't running.
  Langfuse asks for 4 cores / 16 GB / ~100 GB disk (mostly ClickHouse); size the prod VPS for it, or point the SDK at
  Langfuse Cloud (EU) and drop five containers.

## Versions (checked 31 Oct 2026)

Go 1.27, Bun 1.4.4, Elysia 1.4.30, Better Auth 1.7.14, Python 3.14, uv 0.13, FastAPI 0.162, grpcio 1.85.1,
connect-go 1.21 (v1 API), ts-proto 2.18, Expo SDK 58 (React Native 0.88) on Node 26 LTS (`.nvmrc`; app tooling
only, no service uses Node), TanStack Query 5.105, Zustand 5.0,
PostgreSQL 18, Redis 8, Prometheus 3.15, Grafana 13.3, cAdvisor 0.60, Langfuse 4, ClickHouse 26.3, SeaweedFS 4.00,
Mailpit 1.32, Docker Compose ≥ 5.0. Image versions are set in `.env`.

## Not built yet

Everything below is planned task by task in `docs/EPICs/` (order in `docs/EPICs/ROADMAP.md`).

- Product features beyond the one RPC per service.
- Events: outboxes, Redis Streams publishers and consumers (IDN-01, KIT-01, PLN-01, OFF-01, NTF-01).
- identity: sign-in by e-mail code instead of passwords, anonymous accounts, invites, entitlements, account deletion.
- Gateway: Free-tier AI quota, rate limits, idempotency keys, the live shopping-list WebSocket.
- notifier: push (Expo Push API), devices and preferences, founder alerts.
- ai: the Claude provider (settings, `ANTHROPIC_API_KEY` and the job registry exist; the provider is AI-01-03),
  cost tracking and the spend limit (needs `REDIS_STATE_URL`), batch jobs.
- admin: the whole service (ADM-01), its database, Caddy route and Prometheus target.
- Remote push in the app needs a development build: Expo Go on SDK 58 does not receive remote push notifications.
- Off-site copies of the nightly database backups; a container registry and CI.
- Better Auth client in the app (MOB-01).
