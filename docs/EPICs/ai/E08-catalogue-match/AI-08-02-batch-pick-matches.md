# AI-08-02 Batch kind pick_catalogue_matches

| | |
|---|---|
| Service | ai |
| Epic | [AI-08 Catalogue match pick](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-08-01, AI-04-03 |
| Unblocks | KIT-03 bulk matching, OFF-07 offer matching |

## Goal

Match thousands of names (a week of flyer offers, a catalogue import) at half price.

## Business rules

1. Batch kind `BATCH_KIND_PICK_CATALOGUE_MATCHES`; each `BatchItem.pick_match` is one `MatchItem` (AI-08-01) and
   gives one `MatchPick`.
2. Same prompt version and schema as the sync job; model and effort from `batch.pick_catalogue_matches`.
3. In batch mode ai checks only `reason=NONE` ⇔ empty id and the confidence range. **The caller must check that the
   picked id is one of the candidates it sent** (AI-04-03 rule 4) and treat a wrong id as an `ai_invalid_output`
   item (retryable once).
4. `custom_id` should be the caller's match row id.

## Contract

Uses `BatchItem.pick_match` / `BatchItemResult.pick_match` from AI-04.

## Data

No schema change.

## Errors

Per item, as AI-04-03: `ai_refused` (not retryable), `ai_truncated`, `ai_invalid_output`, `ai_item_errored`, `ai_canceled`, `ai_expired` (retryable). Whole-call errors as AI-04-01.

## Acceptance criteria

- [ ] A 3-item batch (stub provider) round-trips through SubmitBatch → GetBatchResults with the three picks keyed by `custom_id`.
- [ ] A live 100-item batch costs about half of the same 100 items sent sync (report from the eval runner with `--batch`).
