# AI-13-03 Batch kind extract_flyer_offers for chain flyers

| | |
|---|---|
| Service | ai |
| Epic | [AI-13 Flyer offer extraction](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | AI-13-01, AI-04-03 |
| Unblocks | OFF-07 |

## Goal

Read every page of every covered chain flyer once a week at batch prices (about $0.0005 a page).

## Business rules

1. Batch kind `BATCH_KIND_EXTRACT_FLYER_OFFERS`; one item per page.
2. Item input `FlyerPageInput`: `page_ref` (used as `custom_id`, e.g. `pg_<flyer>_<n>`), the page as an image (or a
   single-page PDF), `page_number`, and the flyer's `FlyerMeta` (chain, country, currency, dates, language known from
   the list page).
3. Output `FlyerPageOffers` (AI-13-01). In batch, ai runs the input-independent checks only; offers runs every check
   that needs the meta (AI-13-01 rule 5).
4. One `SubmitBatch` per flyer (a flyer of 20–60 pages at ≈ 400 KB per page fits the 30 MiB limit); larger flyers are
   split by the caller.
5. Store-edition flyers (Kaufland per store) are separate flyers with their own batches, submitted only when a user
   has that store in range (features.md §3 step 8; decided in offers).

## Contract

```proto
message FlyerPageInput {
  string page_ref = 1;
  InputImage image = 2;     // or
  bytes pdf_page = 3;       // single-page PDF
  int32 page_number = 4;
  FlyerMeta meta = 5;
}
```

## Data

No schema change.

## Errors

Per item, as AI-04-03: `ai_refused` (not retryable), `ai_truncated`, `ai_invalid_output`, `ai_item_errored`, `ai_canceled`, `ai_expired` (retryable). Whole-call errors as AI-04-01.

## Acceptance criteria

- [ ] A 30-page stub batch round-trips with 30 `FlyerPageOffers` keyed by `page_ref`.
- [ ] A live 30-page batch costs about half of the same pages sent sync (eval report with `--batch`).
