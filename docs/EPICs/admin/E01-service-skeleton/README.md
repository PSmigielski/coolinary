# ADM-01 Service skeleton

**Goal.** A new Bun + Elysia service `admin` (decision 9, D2) that boots like every other service: fixed address
`10.77.0.23`, HTTP `3000`, admin port `8081`, its own `admin` database, hot reload and debugger in dev, hardened in
prod, published at `admin.<domain>` through Caddy, scraped by Prometheus, with gRPC clients for every owning
service's admin RPCs and the same error shape as the app API.

**Flow.** admin UI (later) → Caddy `admin.<domain>` → admin `/admin/v1/*` → `<Svc>AdminService` RPCs.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [ADM-01-01](ADM-01-01-project-scaffold.md) | Project scaffold (Bun + Elysia) | R1 | M | — |
| [ADM-01-02](ADM-01-02-compose-and-makefile.md) | Compose, Makefile and editor wiring | R1 | M | ADM-01-01 |
| [ADM-01-03](ADM-01-03-database-and-secrets.md) | Database, migrations and secrets | R1 | S | ADM-01-01 |
| [ADM-01-04](ADM-01-04-caddy-and-prometheus.md) | Caddy route and Prometheus scrape | R1 | S | ADM-01-02 |
| [ADM-01-05](ADM-01-05-grpc-clients-errors-openapi.md) | gRPC clients, error model and OpenAPI | R1 | M | ADM-01-01, GW-01 |
