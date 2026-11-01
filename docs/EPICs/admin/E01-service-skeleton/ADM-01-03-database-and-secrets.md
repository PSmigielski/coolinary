# ADM-01-03 Database, migrations and secrets

| | |
|---|---|
| Service | admin |
| Epic | [ADM-01 Service skeleton](README.md) |
| Release | R1 |
| Size | S |
| Depends on | ADM-01-01 |
| Unblocks | ADM-02-03 |

## Goal

admin has its own Postgres database and role (D2), used only for its audit log; no other service reads it.

## Business rules

1. `platform/postgres/initdb/01-databases.sh`: add `admin` to the database loop, so new volumes get database
   `admin` owned by role `admin` (dev password `admin-dev`, prod from `admin_db_password`).
2. Existing volumes are not re-initialised. `admin/README.md` gives the SQL to run once through
   `make psql s=postgres` (`CREATE ROLE admin LOGIN PASSWORD …; CREATE DATABASE admin OWNER admin;` plus
   `REVOKE ALL ON DATABASE admin FROM PUBLIC`).
3. `scripts/make-secrets.sh`: add `admin` to both loops (`admin_db_password`, `admin_database_url`).
4. `docker-compose.prod.yml`: declare secrets `admin_db_password` and `admin_database_url`; mount
   `admin_db_password` into `postgres`; admin reads `DATABASE_URL_FILE=/run/secrets/admin_database_url`.
5. Drizzle as in notifier: `src/db/schema.ts`, `drizzle.config.ts`, `bun run db:generate` writes `drizzle/`,
   migrations run at start-up under a Postgres advisory lock.
6. The postgres-exporter and backup job cover the new database without changes (verify).

## Contract

None.

## Data

Empty schema; first table in ADM-02-03.

## Errors

| Case | Behaviour |
|---|---|
| Database missing | Start-up fails with the connection error; `/healthz` 503 |

## Acceptance criteria

- [ ] Fresh `make dev` on an empty volume creates the `admin` database.
- [ ] `make psql s=admin` connects.
- [ ] `make secrets` creates the two new files and keeps existing ones.
- [ ] The nightly backup contains the `admin` database.
