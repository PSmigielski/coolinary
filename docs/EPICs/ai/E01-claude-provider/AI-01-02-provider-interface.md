# AI-01-02 Provider interface for structured jobs

| | |
|---|---|
| Service | ai |
| Epic | [AI-01 Claude provider foundation](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-01-01 |
| Unblocks | AI-01-03, AI-01-06, AI-01-07, AI-01-09, every job epic |

## Goal

Replace the single-purpose `Provider.extract_items` with one generic call that takes a job spec, a prompt and an
output schema and returns a validated object plus usage, so every job (and the stub) plugs in the same way.

## Business rules

1. New types in `llm/base.py`:
   - `JobSpec` — job name, model, effort, max_tokens, timeout, max_retries, fallbacks flag (resolved from the
     registry and the caller's plan, AI-01-01 rule 3).
   - `PromptParts` — `system_blocks` (stable, cacheable text), `context_blocks` (stable per catalogue or library
     version, cacheable), `user_content` (per-request blocks: text, image, document), all built by the job runner.
   - `Completion[T]` — `parsed: T` (a pydantic model), `model` (as reported by the API), `stop_reason`,
     `usage` (`input`, `output`, `cache_read`, `cache_write_5m`, `cache_write_1h`), `prompt_chars` (size only).
2. `Provider.complete(job: JobSpec, prompt: PromptParts, schema: type[T]) -> Completion[T]` is the only method
   providers implement. `Provider.name` stays.
3. Output schemas are pydantic models, one module per job under `prompts/<job>/schema.py` (AI-03-01). Every model
   sets `extra="forbid"`, so the JSON schema sent to the API has `additionalProperties: false`.
4. The stub provider implements `complete` for every job with deterministic, rule-based or fixture output:
   - `text_to_items` keeps the current regex parser.
   - Every other job returns a fixed, schema-valid fixture from `tests/fixtures/stub/<job>.json`, with ids copied
     from the request where the schema needs them (so caller integration tests get valid references).
   - Stub usage is reported as zero tokens and zero cost.
5. The existing `ExtractItemsFromText` handler moves onto the new interface without changing its gRPC behaviour;
   the existing stub tests keep passing.
6. A job runner helper `run_job(job_name, request_context, prompt, schema)` does, in order: resolve `JobSpec`,
   cache lookup (AI-01-09), spend check (AI-02-04), `complete`, validation and repair (AI-01-06), cost and tracing
   (AI-02), cache store. Every RPC handler calls `run_job`; no handler calls a provider directly.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Provider raises a typed provider error | Passed to the error mapper (AI-01-04) |
| Stub has no fixture for a job | `INTERNAL`, `stub fixture missing for <job>` (test-time bug) |

## Acceptance criteria

- [ ] `ExtractItemsFromText` returns the same items as before for the two existing stub tests.
- [ ] A unit test runs `run_job` with a fake provider and checks the call order: cache, spend check, provider, validation, cost, cache store.
- [ ] Every job schema module produces a JSON schema with `additionalProperties: false` on every object (test walks the schema).
- [ ] With `LLM_PROVIDER=stub`, each RPC added later returns a schema-valid response in its own unit test.
