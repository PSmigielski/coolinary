# AI-04-01 SubmitBatch contract and request building

| | |
|---|---|
| Service | ai |
| Epic | [AI-04 Batch jobs](README.md) |
| Release | R1 |
| Size | L |
| Depends on | AI-01-03, AI-01-05, AI-02-04, AI-01-08 |
| Unblocks | AI-04-02, AI-08-02, AI-09-*, AI-13-03; callers KIT-02, KIT-03, KIT-06, PLN-03, OFF-07 |

## Goal

A caller hands ai a list of same-kind items and gets back a batch id; ai turns each item into a Messages request
from the job's prompt and schema and submits them as one Message Batch.

## Business rules

1. One kind per batch. Kinds and their models/efforts are in the epic README; config lives in the job registry
   (AI-01-01) under `batch.<kind>`.
2. Limits per `SubmitBatch` call: 1–2,000 items and at most 30 MiB of request bytes (gRPC limit 32 MiB, AI-01-08).
   Callers split larger work into several batches. (Anthropic's own limit, 100,000 requests or 256 MB, is never
   reached.)
3. `custom_id` must match `^[a-zA-Z0-9_-]{1,64}$` and be unique within the call; the caller chooses it so it can join
   results to its own rows (e.g. `ing_01J…`, `pg_<flyer>_<page>`).
4. Each item becomes one batch request: `custom_id`, `params` = model, `max_tokens`, `thinking` adaptive,
   `output_config` with effort and the JSON-schema format, system prompt with a 1-hour cache breakpoint (AI-01-05
   rule 4), user content from the item. No `fallbacks` (rejected on the Batches API), no sampling parameters, no
   prefill.
5. Before submitting: input checks per kind (empty names, image types and sizes as in the job task) — any invalid
   item fails the whole call with `INVALID_ARGUMENT ai_bad_request` naming the first bad `custom_id`; then the
   monthly batch spend check (AI-02-04 rule 5).
6. Submission uses `client.messages.batches.create(requests=[...])`. The response returns Anthropic's batch id
   (`msgbatch_…`) unchanged, the prompt version, the model and the cost estimate.
7. ai keeps nothing about the batch (D10). The caller stores `batch_id`, `kind`, `prompt_version` and its inputs.
8. With `LLM_PROVIDER=stub`, `SubmitBatch` returns `stub_batch_<uuid>`, and results (AI-04-03) are the stub
   fixtures, so caller pipelines run end to end in dev without cost.

## Contract

`contracts/coolinary/ai/v1/batch.proto`:

```proto
rpc SubmitBatch(SubmitBatchRequest) returns (SubmitBatchResponse);

enum BatchKind {
  BATCH_KIND_UNSPECIFIED = 0;
  BATCH_KIND_TRANSLATE_INGREDIENTS = 1;
  BATCH_KIND_MATCH_NUTRITION_ENTRIES = 2;
  BATCH_KIND_PICK_CATALOGUE_MATCHES = 3;
  BATCH_KIND_DRAFT_DISHES = 4;
  BATCH_KIND_TRANSLATE_DISHES = 5;
  BATCH_KIND_SUGGEST_EFFORT_TAGS = 6;
  BATCH_KIND_EXTRACT_FLYER_OFFERS = 7;
}

message BatchItem {
  string custom_id = 1;
  oneof input {
    TranslateIngredientInput translate_ingredient = 10;   // AI-09-01
    NutritionMatchInput nutrition_match = 11;             // AI-09-02
    MatchItem pick_match = 12;                            // AI-08-01
    DishBrief dish_brief = 13;                            // AI-09-03
    DishTranslationInput dish_translation = 14;           // AI-09-04
    EffortTagInput effort_tags = 15;                      // AI-09-05
    FlyerPageInput flyer_page = 16;                       // AI-13-01
  }
}

message SubmitBatchRequest {
  BatchKind kind = 1;
  repeated BatchItem items = 2;
  string caller = 3;        // "kitchen" | "planner" | "offers", for metrics
}

message SubmitBatchResponse {
  string batch_id = 1;
  int32 items = 2;
  string prompt_version = 3;
  string model = 4;
  int64 estimated_cost_eur_micros = 5;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| 0 or > 2,000 items, bad or duplicate `custom_id`, item input of another kind | `INVALID_ARGUMENT`, `ai_bad_request` |
| Monthly batch limit reached | `RESOURCE_EXHAUSTED`, `ai_budget` |
| API errors on create | as AI-01-04 |

## Acceptance criteria

- [ ] A 3-item `translate_ingredients` submission produces three batch requests with identical system blocks and a 1-hour cache breakpoint (request body captured).
- [ ] No batch request contains `fallbacks`, `temperature`, `top_p`, `top_k` or forced `tool_choice`.
- [ ] Duplicate `custom_id` is refused before any API call.
- [ ] Stub mode returns a `stub_batch_` id and caller tests can fetch fixture results.
