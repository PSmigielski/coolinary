#!/usr/bin/env bash
# Runs one service's tests; `make test s=<svc>` calls it for each service.
#
# Tests that need Postgres or Redis read TEST_DATABASE_URL, TEST_REDIS_CACHE_URL and TEST_REDIS_STATE_URL and must
# skip when they are unset. Those are set only when the dev stack is up (`make dev`): the service then gets its own
# `<svc>_test` database (created on first use, owned by the service role) and Redis database 15, never the dev data.
set -euo pipefail
cd "$(dirname "$0")/.."

svc="${1:?usage: test-service.sh <service>}"
NET_PREFIX="${NET_PREFIX:-10.77.0}"
GO_VERSION="${GO_VERSION:-1.27}"
BUN="${BUN:-$(command -v bun || echo "$HOME/.bun/bin/bun")}"
UV="${UV:-$(command -v uv || echo "$HOME/.local/bin/uv")}"
CACHE="${XDG_CACHE_HOME:-$HOME/.cache}/coolinary-lint"
DEV=(docker compose -f docker-compose.dev.yml)

echo "── test $svc"
[[ -d "$svc" ]] || { echo "no folder $svc"; exit 1; }

stack=""
if docker network inspect coolinary >/dev/null 2>&1 && [[ -n "$("${DEV[@]}" ps --status running -q postgres 2>/dev/null)" ]]; then
  stack=1
fi

# Create <svc>_test when the service has a database role.
db=""
if [[ -n "$stack" ]]; then
  psql=("${DEV[@]}" exec -T postgres psql -U postgres -qtA)
  if [[ "$("${psql[@]}" -c "SELECT 1 FROM pg_roles WHERE rolname = '$svc'")" == 1 ]]; then
    if [[ "$("${psql[@]}" -c "SELECT 1 FROM pg_database WHERE datname = '${svc}_test'")" != 1 ]]; then
      "${psql[@]}" -c "CREATE DATABASE ${svc}_test OWNER $svc" >/dev/null
    fi
    db=1
  fi
else
  echo "   dev stack not running: tests that need Postgres or Redis are skipped (start it with make dev)"
fi

# Environment for tests, by where they run: on the host (by IP) or in a container on the compose network (by name).
test_env() {
  local pg="$1" cache="$2" state="$3"
  [[ -n "$stack" ]] || return 0
  [[ -n "$db" ]] && echo "TEST_DATABASE_URL=postgres://$svc:$svc-dev@$pg:5432/${svc}_test?sslmode=disable"
  echo "TEST_REDIS_CACHE_URL=redis://$cache:6379/15"
  echo "TEST_REDIS_STATE_URL=redis://$state:6379/15"
}

case "$svc" in
  kitchen | planner | offers)
    mkdir -p "$CACHE"
    args=(--rm --user "$(id -u):$(id -g)" -e HOME=/tmp -e GOTOOLCHAIN=local -e GOFLAGS=-buildvcs=false
      -e CGO_ENABLED=0 -e GOCACHE=/cache/go-build -e GOMODCACHE=/cache/go-mod
      -v "$CACHE:/cache" -v "$PWD/$svc:/src" -w /src)
    if [[ -n "$stack" ]]; then
      args+=(--network coolinary)
      while read -r kv; do [[ -n "$kv" ]] && args+=(-e "$kv"); done < <(test_env postgres redis-cache redis-state)
    fi
    docker run "${args[@]}" "golang:${GO_VERSION}-alpine" go test ./...
    ;;
  gateway | identity | notifier | admin)
    cd "$svc"
    if [[ -z "$(find . -path ./node_modules -prune -o \( -name '*.test.ts' -o -name '*_test.ts' -o -name '*.spec.ts' \) -print -quit)" ]]; then
      echo "   no tests yet"
      exit 0
    fi
    "$BUN" install --frozen-lockfile --silent
    env $(test_env "$NET_PREFIX.10" "$NET_PREFIX.11" "$NET_PREFIX.12") "$BUN" test
    ;;
  ai)
    cd ai
    env $(test_env "$NET_PREFIX.10" "$NET_PREFIX.11" "$NET_PREFIX.12") "$UV" run --locked -q pytest -q
    ;;
  mobile)
    echo "   no test runner yet (type-checking runs in make lint s=mobile)"
    ;;
  *)
    echo "unknown service $svc"
    exit 1
    ;;
esac
