#!/bin/bash
# Prod only (pg-backup service): nightly pg_dump of every service database,
# kept for 7 days in the pg-backups volume. Run once by hand with:
#   docker compose -f docker-compose.prod.yml exec pg-backup /backup.sh once
set -euo pipefail
export PGHOST=postgres PGUSER=postgres
export PGPASSWORD="$(cat /run/secrets/postgres_password)"
DBS="identity notifier kitchen planner offers langfuse"
DIR=/backups

dump() {
  stamp="$(date -u +%Y%m%dT%H%M%SZ)"
  for db in $DBS; do
    pg_dump --format=custom --file="$DIR/$db-$stamp.dump" "$db"
    echo "backup: $db-$stamp.dump"
  done
  find "$DIR" -name '*.dump' -mtime +7 -delete
}

if [[ "${1:-}" == "once" ]]; then dump; exit 0; fi
while true; do
  # next run at 03:00 UTC
  now=$(date -u +%s); next=$(date -u -d "tomorrow 03:00" +%s)
  sleep $(( next - now ))
  dump || echo "backup failed" >&2
done
