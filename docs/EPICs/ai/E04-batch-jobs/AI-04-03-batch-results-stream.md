# AI-04-03 GetBatchResults stream, validation and cost

| | |
|---|---|
| Service | ai |
| Epic | [AI-04 Batch jobs](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-04-02, AI-01-06, AI-02-03 |
| Unblocks | AI-04-04; callers KIT-02, KIT-03, KIT-06, PLN-03, OFF-07 |

## Goal

Callers receive every item's result as a typed, schema-checked message, streamed so large batches don't hit message
size limits, and the batch's cost is counted exactly once.

## Business rules

1. Server-streaming RPC (callers are Go services, which handle streams; Bun services never call it). One stream
   message per batch item, in the order the API returns them — **not** submission order; callers join by
   `custom_id`.
2. The caller sends `kind` and `prompt_version` (from `SubmitBatchResponse`), because ai keeps no batch state (D10).
   ai parses each result with that version's schema.
3. Per item, by `result.type`:
   - `succeeded` → check `stop_reason`: `refusal` → error `ai_refused` (not retryable); `max_tokens` → error
     `ai_truncated` (retryable); otherwise parse, apply the schema and the kind's input-independent checks (calorie
     guard, ranges, formats). Failures → error `ai_invalid_output` (retryable). No repair calls in batch mode.
   - `errored` → error `ai_item_errored`, retryable when the API error type is overloaded or server-side, not
     retryable for invalid request.
   - `canceled` → `ai_canceled` (retryable); `expired` → `ai_expired` (retryable).
4. **Input-dependent checks are the caller's job** (e.g. a picked `ingredient_id` must be one of that item's
   candidates; a draft's ingredient ids must be in the allowed list). ai can't do them because it doesn't keep the
   inputs. Each kind's task lists the checks the caller must run.
5. Cost: computed per succeeded item with the batch discount (AI-02-01 rule 4). The batch total is added to the
   monthly batch counter only when a stream reaches the end for the first time (`SET ai:batch:counted:<id> NX`,
   expires after 30 days). A stream cut midway counts nothing; the next complete one counts.
6. Tracing: one Langfuse generation per item, `mode=batch`, sizes only (AI-02-03 rule 5).
7. Calling before the batch has ended returns `FAILED_PRECONDITION ai_batch_not_ended`.

## Contract

```proto
rpc GetBatchResults(GetBatchResultsRequest) returns (stream BatchItemResult);

message GetBatchResultsRequest {
  string batch_id = 1;
  BatchKind kind = 2;
  string prompt_version = 3;
}

message BatchItemResult {
  string custom_id = 1;
  oneof outcome {
    IngredientTranslation translate_ingredient = 10;   // AI-09-01
    NutritionMatch nutrition_match = 11;               // AI-09-02
    MatchPick pick_match = 12;                         // AI-08-01
    DishDraft dish_draft = 13;                         // AI-09-03
    DishTranslation dish_translation = 14;             // AI-09-04
    EffortTagSuggestion effort_tags = 15;              // AI-09-05
    FlyerPageOffers flyer_page = 16;                   // AI-13-01
    BatchItemError error = 20;
  }
  CallInfo call = 30;
}

message BatchItemError {
  string reason = 1;    // ai_refused | ai_truncated | ai_invalid_output | ai_item_errored | ai_canceled | ai_expired
  bool retryable = 2;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Batch not ended | `FAILED_PRECONDITION`, `ai_batch_not_ended` |
| Unknown batch | `NOT_FOUND` |
| Unknown `prompt_version` for the kind | `INVALID_ARGUMENT`, `ai_bad_request` |

## Acceptance criteria

- [ ] A recorded 5-item result file with one success, one refusal, one invalid JSON, one errored and one expired produces the five outcomes above with the right `retryable` flags.
- [ ] Streaming the same batch twice adds its cost to the monthly counter once.
- [ ] A stream cut after 2 of 5 items adds nothing; a later full stream adds the full cost.
