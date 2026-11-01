# AI-04 Batch jobs (Message Batches)

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | Every once-for-everyone job (catalogue translation, nutrition matching, bulk name matching, dish drafts and translations, effort tags, chain flyers) runs through Claude's Message Batches at 50% of the normal price, while ai stays stateless (D10). |

## Flow

```
caller (kitchen KIT-02/03/06, planner PLN-03/04, offers OFF-07)
  SubmitBatch(kind, items[custom_id, input])  -> batch_id, prompt_version, estimate   (AI-04-01)
  stores batch_id + its own inputs
  GetBatch(batch_id) every 10 min              -> status, counts                    (AI-04-02)
  when ended: GetBatchResults(batch_id, kind, prompt_version)  (server stream)      (AI-04-03)
     -> per item: output (schema-checked) or error {reason, retryable}
  caller checks input-dependent rules (ids from its own candidate lists) and stores results
  caller resubmits retryable failures at most twice                                  (AI-04-04)
```

## Kinds

| Kind | Model | Effort | Item input | Defined in |
|---|---|---|---|---|
| `translate_ingredients` | Haiku 5.5 | low | one ingredient | AI-09-01 |
| `match_nutrition_entries` | Haiku 5.5 | medium | one ingredient + table candidates | AI-09-02 |
| `pick_catalogue_matches` | Haiku 5.5 | low | one name + candidates | AI-08-02 |
| `draft_dishes` | Sonnet 5.5 | high | one dish brief | AI-09-03 |
| `translate_dishes` | Haiku 5.5 | low | one dish | AI-09-04 |
| `suggest_effort_tags` | Haiku 5.5 | low | one dish | AI-09-05 |
| `extract_flyer_offers` | Haiku 5.5 | medium | one flyer page | AI-13-03 |

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-04-01](AI-04-01-submit-batch.md) | SubmitBatch contract and request building | R1 | L | AI-01-03, AI-01-05, AI-02-04, AI-01-08 |
| [AI-04-02](AI-04-02-get-and-cancel-batch.md) | GetBatch and CancelBatch | R1 | S | AI-04-01 |
| [AI-04-03](AI-04-03-batch-results-stream.md) | GetBatchResults stream, validation and cost | R1 | M | AI-04-02, AI-01-06, AI-02-03 |
| [AI-04-04](AI-04-04-batch-failures-and-metrics.md) | Failed items, retry rules and batch metrics | R1 | S | AI-04-03 |
