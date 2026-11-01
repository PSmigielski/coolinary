# AI-01-03 Claude provider

| | |
|---|---|
| Service | ai |
| Epic | [AI-01 Claude provider foundation](README.md) |
| Release | R1 |
| Size | L |
| Depends on | AI-01-02 |
| Unblocks | AI-01-04, AI-01-05, AI-04-01, every job epic |

## Goal

Call Claude through the official SDK with structured output, adaptive thinking and per-job effort, and return a
`Completion` the rest of the service can trust.

## Business rules

1. One `anthropic.AsyncAnthropic` client is created in the FastAPI lifespan (same place the gRPC server starts) and
   closed on shutdown. It uses the SDK's default async HTTP client (`DefaultAsyncHttpxClient`, or
   `DefaultAioHttpClient` if measured faster); never a hand-built HTTP client.
2. Requests use the Messages API. Structured output through `client.messages.parse(..., output_format=<pydantic
   model>)` or `output_config={"format": {"type": "json_schema", "schema": ...}}`; never forced `tool_choice`
   (`any`/`tool` returns 400 on Sonnet 5.5), never an assistant prefill (400 on both models).
3. Thinking: `thinking={"type": "adaptive"}` (or omitted) on both models; depth through
   `output_config.effort` from the job spec. Effort is always sent explicitly, because the model defaults differ
   (Haiku 5.5 `medium`, Sonnet 5.5 `high`). No `budget_tokens`.
4. No `temperature`, `top_p` or `top_k` on any request (400 since API version `2026-10-21`). Variety between "new
   plan" requests comes from the request content (AI-01-09 rule 4), not sampling.
5. `max_tokens` from the job spec. Requests are non-streaming while `max_tokens` ≤ 16000; above that the provider
   uses `client.messages.stream(...)` and `get_final_message()` so long outputs don't hit HTTP timeouts.
6. Refusals: before reading content, check `stop_reason`. On `"refusal"`, raise `ProviderRefusal(category)` with
   `stop_details.category` when present. Never return partial content from a refused response.
7. Server-side fallbacks for Sonnet jobs only: when the job's `fallbacks=true`, send the request on the beta
   messages endpoint with `betas=["server-side-fallback-2026-07-01"]` and `fallbacks="default"`. The model that
   actually answered (`response.model`) is recorded and priced (AI-02-01). If the final `stop_reason` is still
   `"refusal"`, rule 6 applies. Haiku 5.5 has no server-side fallback, so Haiku refusals go straight to rule 6.
   Fallbacks are never sent on Batches (rejected there; AI-04-01). Verify the beta binding in the installed SDK
   before coding; do not guess method names.
8. `stop_reason == "max_tokens"`: the output is incomplete; raise `ProviderTruncated`. The validation step
   (AI-01-06) may repair once with `max_tokens` doubled, capped at 64000.
9. Content blocks are read by type: `thinking` blocks are ignored (display stays the default `"omitted"`); the
   structured result is the parsed output, never string matching on raw text.
10. Images go in as base64 `image` blocks with the right `media_type` (`image/jpeg`, `image/png`, `image/webp`);
    PDFs as base64 `document` blocks with `media_type="application/pdf"`, placed before the text block of the
    same message.
11. Usage is copied from `response.usage`: `input_tokens`, `output_tokens`, `cache_read_input_tokens`,
    `cache_creation_input_tokens` (split into 5-minute and 1-hour where the SDK reports it).
12. `get_provider("anthropic")` returns this provider; `get_provider("stub")` the stub (AI-01-02).

## Contract

No gRPC change (each job adds its RPC).

## Data

No schema change.

## Errors

| Case | Raised (mapped in AI-01-04) |
|---|---|
| `stop_reason == "refusal"` | `ProviderRefusal` |
| `stop_reason == "max_tokens"` | `ProviderTruncated` |
| Parsed output fails pydantic validation | `ProviderInvalidOutput` |
| SDK `RateLimitError`, `APIStatusError` 5xx/529, `APIConnectionError`, `APITimeoutError` | passed through |
| SDK `BadRequestError` (our request is wrong) | passed through, logged at `error` |

## Acceptance criteria

- [ ] Recorded-response tests (SDK transport mocked) cover: success, refusal with category, max_tokens, invalid JSON, 429, 529, connection error.
- [ ] A request log assertion shows no `temperature`/`top_p`/`top_k`, no prefill, no forced `tool_choice` in any job's request body.
- [ ] A Sonnet job with `fallbacks=true` sends the beta header and `fallbacks: "default"`; a Haiku job never does.
- [ ] A live smoke test (`make ai-smoke`, manual, needs a key, costs < $0.01) runs `text_to_items` once on each model and prints model, tokens and cost.
- [ ] `CallInfo.model` equals `response.model`, including after a fallback.

## Notes

- Decided (Q017): server-side refusal fallbacks are on for every Sonnet job (`fallbacks=true` in the registry), so a
  misfired safety decline on a recipe still returns a plan. The answering model's real price is recorded. Haiku jobs
  have no server-side fallback; a refusal goes to the caller's non-AI fallback.
