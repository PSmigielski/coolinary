# ai

Python 3.14 + uv + FastAPI. Every LLM call in the system. Callers are the Go services that own each feature; they check AI output with rules before anything is shown. Results are cached in redis-cache; every call is traced in Langfuse with feature, plan and token usage, never user text.

- gRPC `50051` (`grpc.aio`, started in FastAPI's lifespan), admin HTTP `8081` (FastAPI: `/healthz`, `/metrics`), debugpy `5678` in dev.
- RPC: `AiService.ExtractItemsFromText`.
- LLM providers live in `src/coolinary_ai/llm/`; only `stub` (rule-based) exists until a provider is chosen (`LLM_PROVIDER`).
- Generated gRPC code: `src/coolinary/` (`make proto`, via `scripts/gen-proto.sh`).
- Tests: `uv run pytest`; lint: `make lint s=ai` from the repo root (ruff check + format, mypy strict).

Runs with the rest of the stack: `make dev` from the repo root (hot reload, debugger). Settings: `.env.example`
(all variables), `.env.dev` (dev), `.env.prod` (prod, secrets via `*_FILE`). Shared rules:
`../docs/service-conventions.md`. Generated gRPC code is committed; regenerate with `make proto`.
