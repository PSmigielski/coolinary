# Coolinary

Fridge-to-meals planner. Product docs: `docs/business-model.md`, `docs/features.md`.
How the infrastructure works and why: `docs/infrastructure.md`. Rules every service follows:
`docs/service-conventions.md`.

## What's here

| Folder | What | Runtime |
|---|---|---|
| `gateway/` | Public HTTP API for the app | Bun + Elysia |
| `identity/` | Accounts, households, JWTs (Better Auth) | Bun + Elysia |
| `notifier/` | E-mail (later push) | Bun |
| `kitchen/` | Catalogue, barcodes, nutrition, fridge stock | Go |
| `planner/` | Dish suggestions, plans, shopping list | Go |
| `offers/` | Shops, flyers, offers | Go |
| `ai/` | Every LLM call, traced in Langfuse | Python + uv + FastAPI |
| `mobile/` | The app | Expo + TanStack Query + Zustand |
| `contracts/` | gRPC contracts shared by all services | buf |
| `platform/` | Config for Postgres, Redis, Prometheus, Grafana, Langfuse, Caddy | |

Services talk to each other over gRPC; the app talks HTTP to the gateway.

## Prerequisites

- Linux with Docker Engine and **Docker Compose ≥ 5.0** (`docker compose version`).
- `make`, `rsync`, `openssl`.
- For the app: Bun (`curl -fsSL https://bun.sh/install | bash`) and Node 26, the current LTS (`.nvmrc`; with nvm:
  `nvm install && nvm use`). Only the Expo tooling runs on Node; every service runs in Docker.
- For `make lint`: Bun and uv (`curl -LsSf https://astral.sh/uv/install.sh | sh`). The Go and buf linters run in
  containers.
- For editor support in Go files: Go (`sudo pacman -S go`). Nothing else: buf, grpcurl and all code generators
  run in containers.

## Start

```sh
make hooks      # once per clone: lint what you push (see Linting)
make dev        # services + Postgres + Redis + Mailpit, with hot reload
make dev-full   # the same plus Prometheus, Grafana, Langfuse
make ps         # status
make down       # stop (keeps data)
```

The first `make dev` builds the images (a few minutes). Then:

```sh
curl http://localhost:3080/v1/system/status          # every service, checked over gRPC
curl http://localhost:3080/v1/products/9520000000017 # gateway -> kitchen
make smoke                                           # gRPC smoke tests, Bun <-> Go both ways
```

| What | URL (dev) |
|---|---|
| API (gateway) | `http://localhost:3080` (port from `GATEWAY_HOST_PORT` in `.env`), also `http://10.77.0.20:3000` |
| OpenAPI | `http://localhost:3080/openapi` |
| Mailpit | `http://10.77.0.3:8025` |
| Grafana | `http://10.77.0.51:3000` (admin / admin; viewing needs no login) |
| Prometheus | `http://10.77.0.50:9090` |
| Langfuse | `http://10.77.0.60:3000` (dev@coolinary.test / coolinary-dev) |
| Postgres | `10.77.0.10:5432` (postgres / postgres-dev; each service: `<name>` / `<name>-dev`) |

Every container has a fixed address on `10.77.0.0/24` (table in `docs/infrastructure.md`) and is reachable from
this machine by IP; only the gateway port is published. Optional `/etc/hosts` lines:

```
10.77.0.51  grafana.coolinary.test
10.77.0.60  langfuse.coolinary.test
10.77.0.3   mail.coolinary.test
```

### Try the auth flow

```sh
curl -c /tmp/c -H 'Content-Type: application/json' -H 'Origin: http://localhost:3080' \
  -d '{"email":"ewa@coolinary.test","password":"correct-horse-battery","name":"Ewa"}' \
  http://localhost:3080/api/auth/sign-up/email
TOKEN=$(curl -s -b /tmp/c -H 'Origin: http://localhost:3080' http://localhost:3080/api/auth/token | jq -r .token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:3080/v1/household   # gateway -> identity over gRPC
```

## Everyday commands

Run `make help` for the full list.

| Command | Does |
|---|---|
| `make logs s=kitchen` | Follow one service's logs (or all without `s=`) |
| `make sh s=kitchen` | Shell in a running container |
| `make psql s=kitchen` | psql into a service's database (`s=postgres` for the admin database) |
| `make grpc s=kitchen args=list` | grpcurl (reflection is on in dev); call a method with `d='{"ean":"9520000000017"}' args=coolinary.kitchen.v1.KitchenService/ResolveBarcode` |
| `make proto` | Lint `contracts/` and regenerate gRPC code in every service |
| `make proto-breaking` | Check contract changes against `main` |
| `make sqlc` | Regenerate Go query code after changing `queries.sql` |
| `make openapi` | Export the gateway's OpenAPI document and regenerate the app's API types |
| `make mobile` | Start the Expo dev server |
| `make lint` | Lint and type-check every service, the app and `contracts/` (or some: `s="kitchen ai"`) |
| `make lint-fix` | Apply formatters and safe lint fixes (same `s=`) |

Changing a contract: edit `contracts/coolinary/<svc>/v1/*.proto`, run `make proto`, implement, commit the generated
code with it. Running services reload on their own; restart the Bun ones (`make dev`) to refresh their reflection data.

Changing env files: running containers keep their old environment; run `make dev` again to recreate them.

## Linting

Every service has a linter, a type checker and a formatter check, configured in the service's own folder:

| Folder | Checks | Config |
|---|---|---|
| `gateway/`, `identity/`, `notifier/` | Biome (lint, format, import order), `tsc` | `biome.json`, `tsconfig.json` |
| `kitchen/`, `planner/`, `offers/` | golangci-lint (go vet, staticcheck, gosec, gofmt, import groups, ...) | `.golangci.yml` |
| `ai/` | ruff (lint + format), mypy (strict) | `pyproject.toml` |
| `mobile/` | ESLint (Expo config, no warnings allowed), `tsc` | `eslint.config.js`, `tsconfig.json` |
| `contracts/` | buf lint (STANDARD), buf format | `buf.yaml` |

`make hooks` sets `core.hooksPath` to `.githooks/`, whose `pre-push` hook runs `make lint` for the folders the
pushed commits touch and stops the push when something fails. It checks the working tree, so uncommitted changes
count too. `make lint-fix` fixes formatting and most simple issues; `git push --no-verify` skips the hook once.
Generated code (`gen/`, sqlc `store/`, `src/coolinary/`, `auth-schema.ts`, Drizzle migrations) is never linted.

## Planned work

Every planned feature is a task in `docs/EPICs/` (start with `docs/EPICs/README.md`), built in the order of
`docs/EPICs/ROADMAP.md`. Claude Code builds them one at a time: `/next-task` in an interactive session, or
`scripts/epics/run-tasks.sh` unattended (stops after each milestone). `make test s=<svc>` runs a service's tests.

## Debugging (Zed)

Start `make dev`, then in Zed open the debugger and pick a configuration from `.zed/debug.json`:

- **Go** (`kitchen`, `planner`, `offers`): Delve is always listening in the container. Attach, set breakpoints, call
  the service (`make grpc ...` or through the gateway). Each code change restarts the service, so attach again.
  If a Go service stops answering after a debug session, it is paused in Delve: press Continue, or
  `docker compose -f docker-compose.dev.yml restart <svc>`.
- **Python** (`ai`): same, debugpy on `10.77.0.40:5678`.
- **Bun** (`gateway`, `identity`, `notifier`): Zed can't debug Bun yet. Run the Zed task "debug Bun: identity"
  (or `make debug s=identity`), which restarts the service under Bun's inspector and opens
  `https://debug.bun.sh/#localhost:6500/identity`. `make undebug s=identity` turns hot reload back on.

`.zed/tasks.json` has tasks for start/stop, logs, `make proto` and the smoke tests.

## The app

```sh
make mobile     # then scan the QR code with Expo Go
```

The app reads the API address from `mobile/.env` (`EXPO_PUBLIC_API_URL`, this machine's LAN IP and
`GATEWAY_HOST_PORT`). Details in `mobile/README.md`.

## Production

On the server (Linux, Docker Compose ≥ 5.0), with DNS for `api.`, `grafana.` and `langfuse.` under your domain:

```sh
# set DOMAIN in .env
make secrets    # generates ./secrets (never overwrites); then edit secrets/smtp_url
make prod       # builds the images and starts everything behind Caddy (ports 80/443 only)
make backup     # run a database backup now (it also runs nightly at 03:00 UTC)
```

Grafana: `admin` / `secrets/grafana_admin_password`. Langfuse: `admin@<domain>` / `secrets/langfuse_admin_password`.
Langfuse needs about 4 cores / 16 GB RAM on its own; see `docs/infrastructure.md`.
