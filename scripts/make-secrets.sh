#!/usr/bin/env bash
# Generates production secrets into ./secrets (gitignored). Never overwrites an existing file,
# so it is safe to re-run after adding a new secret. Edit secrets/smtp_url and secrets/anthropic_api_key by hand.
#
# Files are 0644 because Compose (without Swarm) bind-mounts them as they are and our
# images run as uid 65532. The folder itself is 0700.
set -euo pipefail
cd "$(dirname "$0")/.."
set -a; source .env; set +a

dir=secrets
mkdir -p "$dir" && chmod 700 "$dir"

rand() { openssl rand -base64 48 | tr -dc 'A-Za-z0-9' | head -c "${1:-40}"; }

write() { # write <name> <value>
  local f="$dir/$1"
  if [[ -e "$f" ]]; then echo "keep     $f"; return; fi
  printf '%s' "$2" > "$f" && chmod 644 "$f" && echo "created  $f"
}

write postgres_password "$(rand)"
for svc in identity notifier kitchen planner offers langfuse postgres_exporter; do
  write "${svc}_db_password" "$(rand)"
done
for svc in identity notifier kitchen planner offers; do
  write "${svc}_database_url" "postgres://${svc}:$(cat "$dir/${svc}_db_password")@postgres:5432/${svc}?sslmode=disable"
done

write better_auth_secret "$(rand 64)"
write smtp_url "smtp://user:password@smtp.example.com:587"
write anthropic_api_key "sk-ant-replace-me"
write trace_hash_salt "$(rand 32)"
write grafana_admin_password "$(rand 24)"
write langfuse_public_key "pk-lf-$(rand 32)"
write langfuse_secret_key "sk-lf-$(rand 32)"
write langfuse_admin_password "$(rand 24)"
write clickhouse_password "$(rand)"
write s3_secret_key "$(rand)"

# Langfuse and ClickHouse don't read *_FILE variables, so they get env files.
write clickhouse.env "CLICKHOUSE_USER=clickhouse
CLICKHOUSE_PASSWORD=$(cat "$dir/clickhouse_password")
"
s3_secret="$(cat "$dir/s3_secret_key")"
write langfuse.env "NEXTAUTH_URL=https://langfuse.${DOMAIN}
NEXTAUTH_SECRET=$(rand 48)
SALT=$(rand 32)
ENCRYPTION_KEY=$(openssl rand -hex 32)
TELEMETRY_ENABLED=false
DATABASE_URL=postgresql://langfuse:$(cat "$dir/langfuse_db_password")@postgres:5432/langfuse
CLICKHOUSE_MIGRATION_URL=clickhouse://clickhouse:9000
CLICKHOUSE_URL=http://clickhouse:8123
CLICKHOUSE_USER=clickhouse
CLICKHOUSE_PASSWORD=$(cat "$dir/clickhouse_password")
CLICKHOUSE_CLUSTER_ENABLED=false
REDIS_CONNECTION_STRING=redis://redis-langfuse:6379
S3_ACCESS_KEY=langfuse
S3_SECRET_KEY=${s3_secret}
LANGFUSE_S3_EVENT_UPLOAD_BUCKET=langfuse
LANGFUSE_S3_EVENT_UPLOAD_REGION=us-east-1
LANGFUSE_S3_EVENT_UPLOAD_ACCESS_KEY_ID=langfuse
LANGFUSE_S3_EVENT_UPLOAD_SECRET_ACCESS_KEY=${s3_secret}
LANGFUSE_S3_EVENT_UPLOAD_ENDPOINT=http://seaweedfs:8333
LANGFUSE_S3_EVENT_UPLOAD_FORCE_PATH_STYLE=true
LANGFUSE_S3_EVENT_UPLOAD_PREFIX=events/
LANGFUSE_S3_MEDIA_UPLOAD_BUCKET=langfuse
LANGFUSE_S3_MEDIA_UPLOAD_REGION=us-east-1
LANGFUSE_S3_MEDIA_UPLOAD_ACCESS_KEY_ID=langfuse
LANGFUSE_S3_MEDIA_UPLOAD_SECRET_ACCESS_KEY=${s3_secret}
LANGFUSE_S3_MEDIA_UPLOAD_ENDPOINT=http://seaweedfs:8333
LANGFUSE_S3_MEDIA_UPLOAD_FORCE_PATH_STYLE=true
LANGFUSE_S3_MEDIA_UPLOAD_PREFIX=media/
LANGFUSE_INIT_ORG_ID=coolinary
LANGFUSE_INIT_ORG_NAME=Coolinary
LANGFUSE_INIT_PROJECT_ID=coolinary-prod
LANGFUSE_INIT_PROJECT_NAME=coolinary-prod
LANGFUSE_INIT_PROJECT_PUBLIC_KEY=$(cat "$dir/langfuse_public_key")
LANGFUSE_INIT_PROJECT_SECRET_KEY=$(cat "$dir/langfuse_secret_key")
LANGFUSE_INIT_USER_EMAIL=admin@${DOMAIN}
LANGFUSE_INIT_USER_NAME=Admin
LANGFUSE_INIT_USER_PASSWORD=$(cat "$dir/langfuse_admin_password")
"
echo
echo "Done. Edit $dir/smtp_url with your mail provider and $dir/anthropic_api_key with your Claude API key."
echo "Grafana:  admin / $dir/grafana_admin_password"
echo "Langfuse: admin@${DOMAIN} / $dir/langfuse_admin_password"
