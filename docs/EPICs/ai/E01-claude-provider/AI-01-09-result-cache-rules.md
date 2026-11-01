# AI-01-09 Result cache rules per job

| | |
|---|---|
| Service | ai |
| Epic | [AI-01 Claude provider foundation](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | AI-01-02 |
| Unblocks | every job epic |

## Goal

The same request never costs twice, while a user asking for "another plan" still gets a new one.

## Business rules

1. Cache key: `ai:<job>:<sha256-32>` over, in order: job name, provider name, model, effort, prompt version
   (AI-03-01), language, and the canonical request bytes (protobuf deterministic serialisation of the request
   without `household_id` where the job output doesn't depend on the household, e.g. extraction; with it for
   composition jobs).
2. TTLs (registry `cache_ttl_seconds`): text_to_items 7 days, photo_to_items 1 day, read_label 30 days,
   pick_matches 30 days, extract_flyer_offers 30 days, extract_receipt 1 day, composition jobs 1 hour.
3. Never cached: errors, refusals, outputs that needed item dropping above 0, and any response when the request
   sets `bypass_cache=true`.
4. Composition requests carry `variation` (int, the n-th "new plan" for the same inputs) and
   `exclude_dish_names` / `exclude_dish_ids`; both are part of the key, so "new plan" is a cache miss by design.
5. A cache hit returns `CallInfo.cached=true` and `cost_eur_micros=0`, and still emits the request metric with
   `outcome="cached"`.
6. The cache stays best-effort (`redis-cache` may evict at any time; errors are ignored, as today).
7. Photos and documents are hashed for the key, never stored (D6); only the job output is cached.

## Contract

Composition request messages get `int32 variation` and `bool bypass_cache` (defined in AI-10-01 and reused).

## Data

No schema change (Redis keys only).

## Errors

None new.

## Acceptance criteria

- [ ] Two identical `ExtractItemsFromPhoto` calls: second is `cached=true`, cost 0, one HTTP call total.
- [ ] Same compose request with `variation=1` then `variation=2`: two HTTP calls.
- [ ] A refused request is not cached (second identical call calls the API again).
- [ ] Changing the prompt version of a job changes every key of that job (test).
