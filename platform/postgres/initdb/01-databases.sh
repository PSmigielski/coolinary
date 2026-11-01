#!/bin/bash
# Runs once, when the Postgres volume is first created.
# Creates one database + owner role per service, and a read-only monitoring role.
# Passwords: prod reads $DB_PASSWORD_DIR/<name>_db_password (Docker secrets, required);
#            dev uses "<name>-dev".
# Adding a service later: create its role and database by hand (make psql s=postgres).
set -euo pipefail

password_for() {
  if [[ -n "${DB_PASSWORD_DIR:-}" ]]; then
    cat "$DB_PASSWORD_DIR/$1_db_password"
  else
    echo "$1-dev"
  fi
}

for name in identity notifier kitchen planner offers langfuse; do
  pw="$(password_for "$name")"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres -v name="$name" -v pw="$pw" <<'SQL'
SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', :'name', :'pw') \gexec
SELECT format('CREATE DATABASE %I OWNER %I', :'name', :'name') \gexec
SELECT format('REVOKE ALL ON DATABASE %I FROM PUBLIC', :'name') \gexec
SQL
done

pw="$(password_for postgres_exporter)"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres -v pw="$pw" <<'SQL'
SELECT format('CREATE ROLE postgres_exporter LOGIN PASSWORD %L IN ROLE pg_monitor', :'pw') \gexec
SQL
echo "coolinary: databases created"
