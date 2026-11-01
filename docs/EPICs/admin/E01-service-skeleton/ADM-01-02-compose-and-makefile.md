# ADM-01-02 Compose, Makefile and editor wiring

| | |
|---|---|
| Service | admin |
| Epic | [ADM-01 Service skeleton](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | ADM-01-01 |
| Unblocks | ADM-01-04 |

## Goal

`make dev` and `make prod` start admin like every other service: fixed address, hot reload in dev, hardened in prod.

## Business rules

1. `docker-compose.abstract.yml`: service `admin` with `<<: *app`, build context `./admin`, arg `BUN_VERSION`,
   `depends_on` Postgres healthy only (not the other services, which it calls lazily), address
   `${NET_PREFIX:-10.77.0}.23` (D2).
2. `docker-compose.dev.yml`: image `coolinary/admin:dev`, target `dev`, `env_file: [./admin/.env.dev]`,
   `BUN_INSPECT: ${ADMIN_INSPECT:-0}`, `command: *bun-command`, volumes `./admin:/app` and `/app/node_modules`,
   ports `127.0.0.1:6502:6499` (inspector). No public port: dev reaches it at `http://10.77.0.23:3000`.
3. `docker-compose.prod.yml`: `<<: *hardened`, image `coolinary/admin:${TAG:-latest}`, target `prod`,
   `env_file: [./admin/.env.prod]`, environment `AUTH_ISSUER: https://api.${DOMAIN}` and
   `ADMIN_PUBLIC_URL: https://admin.${DOMAIN}`, `secrets: [admin_database_url]`, limits `cpus: "0.5"`,
   `memory: 192M`.
4. Makefile: add `admin` to `BUN_SERVICES` (so lint, typecheck, install and `make proto` sync include it), add
   `admin) p=6502;;` to `debug`, and an `Admin` line to `make urls`. Header comment of the dev compose lists 6502.
5. `make proto`: admin gets generated client code in `admin/src/gen` for the admin services it calls
   (identity, notifier, kitchen, planner, offers, ai); committed like the other services.
6. `.zed/tasks.json`: task "debug Bun: admin" like the others.
7. Root `README.md` "What's here" and dev URL table get the admin row.

## Contract

No API change.

## Data

No schema change.

## Errors

None.

## Acceptance criteria

- [ ] `make dev` starts `admin`; `make ps` shows it healthy; `make logs s=admin` works.
- [ ] `make debug s=admin` prints `https://debug.bun.sh/#localhost:6502/admin`.
- [ ] `make prod` (local `.localhost` domain) starts it with read-only root filesystem and no capabilities.
- [ ] `make lint` covers admin.
