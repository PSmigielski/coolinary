# ADM-01-01 Project scaffold (Bun + Elysia)

| | |
|---|---|
| Service | admin |
| Epic | [ADM-01 Service skeleton](README.md) |
| Release | R1 |
| Size | M |
| Depends on | — |
| Unblocks | ADM-01-02, ADM-01-03, ADM-01-05 |

## Goal

A new self-contained `admin/` service that boots, serves health and metrics, and follows the same conventions as
the other Bun services, so the back-office routes can be added on top (decision 9).

## Business rules

1. Copy the gateway's layout (closest Bun HTTP service): `admin/src/{index,config,admin,healthcheck}.ts`,
   `src/http/app.ts`, `src/lib/{log,metrics}.ts`, `src/grpc/`, `src/db/`, own `package.json`, `bun.lock`,
   `tsconfig.json`, `biome.json`, `Dockerfile`, `.dockerignore`, `README.md`.
2. Dockerfile stages as the gateway: base → deps → dev / build (typecheck + bundle) → prod on
   `oven/bun:<v>-distroless`, uid 65532, healthcheck = the `healthcheck` entry point.
3. Ports: HTTP `3000` (admin API and auth proxy), admin port `8081` (`GET /metrics`, `GET /healthz`). No gRPC server:
   nobody calls admin.
4. `/healthz` checks only admin's own Postgres (ADM-01-03); `200 {"ok":true,"checks":{...}}` or `503`. It never
   checks other services.
5. Metrics: `http_server_requests_total{method,route,status}`, `http_server_duration_seconds{method,route}`,
   `rpc_client_requests_total{rpc_service,rpc_method,rpc_code}` (same names as the gateway).
6. Logs: JSON lines as in `docs/service-conventions.md`; additionally never log e-mail addresses typed into lookups.
7. Config (any `X` also as `X_FILE`):

| Variable | Dev value |
|---|---|
| `SERVICE_NAME` | `admin` |
| `HTTP_PORT`, `ADMIN_PORT` | `3000`, `8081` |
| `LOG_LEVEL` | `debug` |
| `DATABASE_URL` | `postgres://admin:admin-dev@postgres:5432/admin` |
| `IDENTITY_HTTP_URL` | `http://identity:3000` |
| `AUTH_ISSUER` | `http://localhost:3080` (the public API URL, same as the gateway) |
| `ADMIN_PUBLIC_URL` | `http://10.77.0.23:3000` (prod `https://admin.<domain>`) |
| `IDENTITY_GRPC_ADDR` … `AI_GRPC_ADDR` | `identity:50051`, `notifier:50051`, `kitchen:50051`, `planner:50051`, `offers:50051`, `ai:50051` |
| `PROMETHEUS_URL` | `http://prometheus:9090` |
| `GRAFANA_PUBLIC_URL`, `LANGFUSE_PUBLIC_URL` | `http://10.77.0.51:3000`, `http://10.77.0.60:3000` |

8. Env files: `.env.example` (all variables), `.env.dev` (committed), `.env.prod` (no secrets).

## Contract

`GET /healthz` and `GET /metrics` on 8081 only.

## Data

None yet (ADM-01-03).

## Errors

| Case | Behaviour |
|---|---|
| Required variable missing | Start-up fails naming the variable |

## Acceptance criteria

- [ ] `bun run typecheck` and `bun run lint` pass; the prod image builds and runs as uid 65532 with a read-only root
      filesystem.
- [ ] `curl 10.77.0.23:8081/healthz` → 200 with Postgres `ok`; with Postgres stopped → 503.
- [ ] `/metrics` exposes the metric names above.
