# gateway

Bun + Elysia. Public HTTP API for the mobile app. Verifies JWTs from identity (JWKS, cached), turns HTTP calls into gRPC calls on the services that own the data, proxies `/api/auth/*` to identity, and serves OpenAPI at `/openapi`.

- HTTP `3000` (public, published on the host as `GATEWAY_HOST_PORT`), admin `8081`.
- Routes: `GET /v1/system/status`, `GET /v1/products/{ean}`, `GET /v1/me`, `GET /v1/household`.
- `bun scripts/smoke.ts [idleSeconds]`: Bun → Go gRPC smoke test (`make smoke`).
- `make openapi` exports `contracts/openapi/gateway.json` and refreshes the app's API types.

Runs with the rest of the stack: `make dev` from the repo root (hot reload, debugger). Settings: `.env.example`
(all variables), `.env.dev` (dev), `.env.prod` (prod, secrets via `*_FILE`). Shared rules:
`../docs/service-conventions.md`. Generated gRPC code is committed; regenerate with `make proto`.
