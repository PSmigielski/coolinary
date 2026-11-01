---
paths:
  - "src/**/*.py"
  - "tests/**/*.py"
  - "pyproject.toml"
---

# AI service: FastAPI + grpc.aio patterns

Every LLM call in Coolinary goes through this service. Its callers are the Go services that own each feature,
over gRPC; they check the output with rules before anything is shown. The app never calls this service.

## FastAPI is the admin server only

- `app.py` serves `/healthz` and `/metrics` on 8081, and its `lifespan` starts and stops the gRPC server, so both
  share one event loop. **Features are gRPC methods, never FastAPI routes.**
- Keep one uvicorn worker (`__main__.py`): the process owns the gRPC port.
- New admin routes are `async def`, return `JSONResponse` with an explicit status, and stay out of OpenAPI
  (`docs_url` / `openapi_url` are off). `/healthz` checks only this service's own dependencies (Redis), never
  other services or the LLM provider.

## Adding an AI feature

1. Contract: add the RPC to `contracts/coolinary/ai/v1/ai.proto`, run `make proto` (generated code lands in
   `src/coolinary/`; never edit it).
2. Provider interface: add a method to `Provider` in `llm/base.py` that returns a frozen dataclass with the
   result, `model` and `usage` (token counts). Provider-specific types and SDKs never leave `llm/`.
3. Implement it in `llm/stub.py` with simple rules, so the whole pipeline runs with `LLM_PROVIDER=stub`, and in
   each real provider. Providers are picked in `llm/__init__.py:get_provider`.
4. Servicer method on `AiServicer` in `grpc_server.py`, in this order (copy `ExtractItemsFromText`):
   validate and cap input size → `cache.key(feature, provider.name, ...)` → `cache.get_json` →
   `with tracing.generation(...)` around the provider call → `gen.update(model, output, usage_details)` →
   `cache.set_json` → proto response.
5. Tests in `tests/` for the stub and the mapping.

## gRPC handlers

- Handlers are `async def`. Anything blocking or CPU-heavy goes through `asyncio.to_thread`; use async clients
  only (`redis.asyncio`, async HTTP / provider SDKs). A blocking call stalls both servers.
- Errors: `await context.abort(grpc.StatusCode.X, "message")`: `INVALID_ARGUMENT` for bad input,
  `UNAVAILABLE` when the provider fails or times out (callers then fall back to non-AI results), `INTERNAL` for
  bugs. Give provider calls a timeout shorter than the caller's deadline.
- Caller context: `dict(context.invocation_metadata() or ())` gives `x-plan`, `x-household-id`, `x-request-id`.
- Metrics come from `MetricsInterceptor` with the shared names (`docs/service-conventions.md`); don't add
  per-handler RPC metrics.

## Privacy and cost

- **Never put user text, item names, images or health data in traces or logs.** `input` / `output` on a
  generation hold sizes and counts only (`chars`, `items`). Add a key to `SAFE_KEYS` in `tracing.py` only for
  non-personal values. The household id is hashed (`tracing.household_hash`); tags are the feature and the plan.
- Always report token usage, so Langfuse can show cost per feature and plan.
- The cache key must include everything that changes the answer: provider, model or prompt version, language,
  input. The cache is best-effort; Redis errors never fail a request.

## Everything else

- Settings: add a field to `Settings` and read it in `load()` with `env()` (supports `X_FILE`), and add it to
  `.env.example`, `.env.dev` and `.env.prod`.
- Logs: `log.info("msg", extra={...})` from `log.py`.
- Style: Python 3.14, type hints everywhere, `X | None`, frozen dataclasses for internal types, line length 110.
- Check: `make lint s=ai` from the repo root (ruff check, ruff format, mypy strict; `make lint-fix s=ai` fixes
  formatting), and `uv run pytest` from `ai/` (`asyncio_mode = auto`: tests are plain `async def test_...`).
