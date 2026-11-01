# AI-02-03 Langfuse fields for every job

| | |
|---|---|
| Service | ai |
| Epic | [AI-02 Cost tracking and spend limit](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-02-01 |
| Unblocks | AI-03-03 |

## Goal

Every AI call is one Langfuse generation with job, plan, model, prompt version, token usage and cost, and never
any user content, so cost per household can be analysed safely.

## Business rules

1. Extend `tracing.generation` (already used by `text_to_items`) to every job via `run_job` (AI-01-02).
2. Fields: trace name = job; tags = job, `plan:<plan>`, `model:<model>`, `prompt:<version>`; user id = salted hash of
   the household id (existing `household_hash`); `usage_details` = input, output, cache_read, cache_write;
   `cost_details` = total in USD; metadata = `effort`, `repaired`, `cached`, `mode`.
3. Input and output are sizes only: counts of images, pages, characters, items, candidates, dishes, meals. Extend
   `SAFE_KEYS` with exactly these: `job`, `effort`, `images`, `pages`, `chars`, `items`, `candidates`, `dishes`,
   `meals`, `offers`, `lines`, `violations`, `repaired`, `mode`. Any other string is still masked to `[redacted]`.
4. A repair call is a second generation under the same trace.
5. Batch results are traced when fetched (AI-04-03): one generation per item, `mode=batch`, no household id.
6. Tracing stays optional: off when Langfuse keys are missing or the host doesn't resolve (existing behaviour).

## Contract

No gRPC change.

## Data

No schema change.

## Errors

None new; tracing failures never fail a call.

## Acceptance criteria

- [ ] With a fake Langfuse client, a `compose_day_plan` call records one generation with the tags above and usage; no string field other than the allowed keys.
- [ ] A text containing an e-mail address and a name never appears in the captured trace payload.
- [ ] With Langfuse down, calls succeed and log one `info` line that tracing is off.
