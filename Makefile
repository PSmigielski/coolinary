# Coolinary — everyday commands. Run `make help` for the list.
include .env
export

SHELL := /bin/bash
UID := $(shell id -u)
GID := $(shell id -g)

# Bun's installer adds ~/.bun/bin to interactive shells only; find it either way. Same for uv.
BUN := $(shell command -v bun 2>/dev/null || echo $(HOME)/.bun/bin/bun)
UV  := $(shell command -v uv 2>/dev/null || echo $(HOME)/.local/bin/uv)

DEV  := docker compose -f docker-compose.dev.yml
PROD := docker compose -f docker-compose.prod.yml

GO_SERVICES  := kitchen planner offers
BUN_SERVICES := gateway identity notifier
SERVICES     := $(BUN_SERVICES) $(GO_SERVICES) ai
LINTABLE     := $(SERVICES) mobile contracts

# Code generators and grpcurl run from containers; nothing to install on the host.
PROTOGEN := docker run --rm --user $(UID):$(GID) -v $(CURDIR):/workspace

s ?=
args ?=
d ?=
flags ?=

.PHONY: help
help: ## Show this list
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# ---------------------------------------------------------------- dev

.PHONY: dev
dev: ## Start the dev stack (services + data + Mailpit), hot reload on
	$(DEV) up -d --build --renew-anon-volumes --remove-orphans
	@$(MAKE) --no-print-directory urls

.PHONY: dev-full
dev-full: ## Start the dev stack plus monitoring and Langfuse
	COMPOSE_PROFILES=observability $(DEV) up -d --build --renew-anon-volumes --remove-orphans
	@$(MAKE) --no-print-directory urls

.PHONY: down
down: ## Stop the dev stack (keeps volumes)
	COMPOSE_PROFILES=observability $(DEV) down --remove-orphans

.PHONY: ps
ps: ## Show container status
	COMPOSE_PROFILES=observability $(DEV) ps --format 'table {{.Service}}\t{{.Status}}'

.PHONY: logs
logs: ## Follow logs (all, or one: make logs s=kitchen)
	COMPOSE_PROFILES=observability $(DEV) logs -f --tail 100 $(s)

.PHONY: sh
sh: ## Shell in a running container: make sh s=kitchen
	@test -n "$(s)" || (echo "usage: make sh s=<service>" && exit 1)
	$(DEV) exec $(s) sh

.PHONY: psql
psql: ## psql into a service database: make psql s=kitchen (s=postgres for the admin db)
	@test -n "$(s)" || (echo "usage: make psql s=<service>" && exit 1)
	@if [ "$(s)" = "postgres" ]; then $(DEV) exec postgres psql -U postgres; \
	else $(DEV) exec -e PGPASSWORD=$(s)-dev postgres psql -h 127.0.0.1 -U $(s) $(s); fi

.PHONY: debug
debug: ## Restart a Bun service under the inspector: make debug s=identity (open the printed URL)
	@case "$(s)" in gateway) p=6499;; identity) p=6500;; notifier) p=6501;; *) echo "usage: make debug s=gateway|identity|notifier"; exit 1;; esac; \
	var=$$(echo $(s) | tr a-z A-Z)_INSPECT; env $$var=1 $(DEV) up -d $(s); \
	echo "Open https://debug.bun.sh/#localhost:$$p/$(s)  (make undebug s=$(s) to go back to --watch)"

.PHONY: undebug
undebug: ## Restart a Bun service with hot reload again: make undebug s=identity
	$(DEV) up -d $(s)

.PHONY: grpc
grpc: ## grpcurl: make grpc s=kitchen args=list | make grpc s=kitchen d='{"ean":"9520000000017"}' args=coolinary.kitchen.v1.KitchenService/ResolveBarcode
	@test -n "$(s)" || (echo 'usage: make grpc s=<service> args=list' && exit 1)
	@docker run --rm --network coolinary fullstorydev/grpcurl -plaintext $(if $(d),-d '$(d)',) $(flags) $(s):50051 $(args)

.PHONY: smoke
smoke: ## gRPC smoke tests Bun -> Go and Go -> Bun (idle=<seconds> adds an idle gap)
	$(DEV) exec -T gateway bun scripts/smoke.ts $(idle)
	$(DEV) exec -T kitchen go run ./cmd/smoketest $(if $(idle),-idle $(idle)s,)

.PHONY: urls
urls: ## Print the dev URLs
	@echo ""
	@echo "  API (gateway)   http://localhost:$(GATEWAY_HOST_PORT)  ·  http://$(NET_PREFIX).20:3000"
	@echo "  Mailpit         http://$(NET_PREFIX).3:8025"
	@echo "  Grafana         http://$(NET_PREFIX).51:3000   (dev-full)"
	@echo "  Prometheus      http://$(NET_PREFIX).50:9090   (dev-full)"
	@echo "  Langfuse        http://$(NET_PREFIX).60:3000   (dev-full; dev@coolinary.test / coolinary-dev)"
	@echo "  Postgres        $(NET_PREFIX).10:5432  (user postgres / postgres-dev)"
	@echo ""

# ---------------------------------------------------------------- code generation

.PHONY: protogen
protogen: ## Build the pinned code-generator image (contracts/Dockerfile)
	@docker build -q -t coolinary/protogen:local --build-arg GO_VERSION=$(GO_VERSION) --build-arg BUF_VERSION=$(BUF_VERSION) \
		--build-arg BUN_VERSION=$(BUN_VERSION) --build-arg PYTHON_VERSION=$(PYTHON_VERSION) contracts >/dev/null

.PHONY: proto
proto: protogen ## Lint contracts and generate gRPC code into every service
	$(PROTOGEN) -w /workspace/contracts coolinary/protogen:local buf lint
	@# Generate into .gen/ first, then sync only changed files into each service, so
	@# running dev containers (bun --watch, air, watchfiles) never see a half-empty folder.
	@rm -rf .gen && mkdir -p .gen
	@for s in $(GO_SERVICES) $(BUN_SERVICES); do echo "generate $$s"; \
		$(PROTOGEN) -w /workspace/$$s coolinary/protogen:local buf generate -o /workspace/.gen/$$s || exit 1; done
	@for s in $(BUN_SERVICES); do \
		$(PROTOGEN) -w /workspace/contracts coolinary/protogen:local buf build -o /workspace/.gen/$$s/src/gen/descriptor.binpb || exit 1; done
	@echo "generate ai"; $(PROTOGEN) -w /workspace/ai -e OUT=/workspace/.gen/ai coolinary/protogen:local ./scripts/gen-proto.sh
	@for s in $(GO_SERVICES); do rsync -a --delete --checksum .gen/$$s/gen/ $$s/gen/; done
	@for s in $(BUN_SERVICES); do rsync -a --delete --checksum .gen/$$s/src/gen/ $$s/src/gen/; done
	@rsync -a --delete --checksum .gen/ai/coolinary/ ai/src/coolinary/
	@rm -rf .gen

.PHONY: proto-breaking
proto-breaking: protogen ## Check contracts for breaking changes against main
	@if git rev-parse --verify -q main >/dev/null; then \
		$(PROTOGEN) -w /workspace coolinary/protogen:local buf breaking contracts --against '.git#branch=main,subdir=contracts'; \
	else echo "no main branch yet; skipping"; fi

.PHONY: sqlc
sqlc: ## Regenerate sqlc query code in the Go services
	@for s in $(GO_SERVICES); do docker run --rm --user $(UID):$(GID) -v $(CURDIR)/$$s:/src -w /src sqlc/sqlc:1.30.0 generate || exit 1; done

.PHONY: openapi
openapi: ## Export the gateway OpenAPI document and regenerate the app's API types
	cd gateway && $(BUN) install --frozen-lockfile && $(BUN) scripts/export-openapi.ts ../contracts/openapi/gateway.json
	@if [ -d mobile/node_modules ]; then cd mobile && $(BUN) run api:types; fi

# ---------------------------------------------------------------- lint

# Go and buf run from pinned images (Go isn't needed on the host); Bun and uv run on the host.
# Go caches live outside the repo, so later runs only re-check what changed.
GOLANGCI_LINT_VERSION := v2.15.0
LINT_CACHE := $(or $(XDG_CACHE_HOME),$(HOME)/.cache)/coolinary-lint
GOLANGCI = docker run --rm --user $(UID):$(GID) -e HOME=/tmp -e GOTOOLCHAIN=local -e GOFLAGS=-buildvcs=false \
	-e GOCACHE=/cache/go-build -e GOMODCACHE=/cache/go-mod -e GOLANGCI_LINT_CACHE=/cache/golangci \
	-v $(LINT_CACHE):/cache -v $(CURDIR)/$*:/src -w /src golangci/golangci-lint:$(GOLANGCI_LINT_VERSION)-alpine golangci-lint
BUF = docker run --rm --user $(UID):$(GID) -e BUF_CACHE_DIR=/tmp/buf -v $(CURDIR)/contracts:/workspace -w /workspace \
	bufbuild/buf:$(BUF_VERSION)

.PHONY: lint
lint: ## Lint + type-check everything (the pre-push hook runs this), or some: make lint s="kitchen ai"
	@$(MAKE) --no-print-directory -k -j4 --output-sync=target $(addprefix lint-,$(or $(s),$(LINTABLE)))

.PHONY: lint-fix
lint-fix: ## Apply formatters and safe lint fixes: make lint-fix (or s=gateway)
	@$(MAKE) --no-print-directory -k $(addprefix lint-fix-,$(or $(s),$(LINTABLE)))

.PHONY: hooks
hooks: ## Install the git hooks in .githooks (pre-push runs make lint on what you push)
	git config core.hooksPath .githooks
	@echo "pre-push hook on; skip it once with: git push --no-verify"

# Bun services and the app: Biome / ESLint + tsc, from each package.json ("lint", "lint:fix").
.PHONY: $(addprefix lint-,$(BUN_SERVICES) mobile) $(addprefix lint-fix-,$(BUN_SERVICES) mobile)
$(addprefix lint-,$(BUN_SERVICES) mobile): lint-%:
	@echo "── lint $*"
	@cd $* && $(BUN) install --frozen-lockfile --silent && $(BUN) run --silent lint
$(addprefix lint-fix-,$(BUN_SERVICES) mobile): lint-fix-%:
	@cd $* && $(BUN) install --frozen-lockfile --silent && $(BUN) run --silent lint:fix

# Go services: golangci-lint (.golangci.yml, identical in all three), which includes go vet and gofmt.
.PHONY: $(addprefix lint-,$(GO_SERVICES)) $(addprefix lint-fix-,$(GO_SERVICES))
$(addprefix lint-,$(GO_SERVICES)): lint-%: | $(LINT_CACHE)
	@echo "── lint $*"
	@$(GOLANGCI) run
$(addprefix lint-fix-,$(GO_SERVICES)): lint-fix-%: | $(LINT_CACHE)
	@$(GOLANGCI) run --fix

$(LINT_CACHE):
	@mkdir -p $@

# Python: ruff (lint + format) and mypy, configured in ai/pyproject.toml.
.PHONY: lint-ai lint-fix-ai
lint-ai:
	@echo "── lint ai"
	@cd ai && $(UV) run --locked -q ruff check -q && $(UV) run --locked -q ruff format -q --check && $(UV) run --locked -q mypy
lint-fix-ai:
	@cd ai && $(UV) run --locked -q ruff check --fix && $(UV) run --locked -q ruff format

# Contracts: buf lint (STANDARD) and buf format.
.PHONY: lint-contracts lint-fix-contracts
lint-contracts:
	@echo "── lint contracts"
	@$(BUF) lint && $(BUF) format --diff --exit-code
lint-fix-contracts:
	@$(BUF) format -w

# ---------------------------------------------------------------- tests

# scripts/test-service.sh per service: go test (in a container), bun test, pytest. Tests that need Postgres or Redis
# run only while `make dev` is up, against a <svc>_test database and Redis database 15; otherwise they skip.
TESTABLE := $(SERVICES) mobile

.PHONY: test
test: ## Run tests, all or some: make test s="kitchen ai" (needs `make dev` for Postgres/Redis tests)
	@$(MAKE) --no-print-directory -k $(addprefix test-,$(or $(s),$(TESTABLE)))

.PHONY: $(addprefix test-,$(TESTABLE))
$(addprefix test-,$(TESTABLE)): test-%:
	@./scripts/test-service.sh $*

# ---------------------------------------------------------------- mobile

.PHONY: mobile
mobile: ## Start the Expo dev server (scan the QR code with Expo Go)
	cd mobile && $(BUN) install && $(BUN) x expo start

# ---------------------------------------------------------------- prod

.PHONY: secrets
secrets: ## Generate prod secrets into ./secrets (never overwrites existing files)
	./scripts/make-secrets.sh

.PHONY: prod
prod: ## Build and start the prod stack (needs `make secrets` first)
	@test -d secrets || (echo "run 'make secrets' first" && exit 1)
	COMPOSE_PROFILES=observability $(PROD) up -d --build --remove-orphans

.PHONY: prod-down
prod-down: ## Stop the prod stack (keeps volumes)
	COMPOSE_PROFILES=observability $(PROD) down --remove-orphans

.PHONY: backup
backup: ## Run a database backup now (prod)
	$(PROD) exec pg-backup /backup.sh once
