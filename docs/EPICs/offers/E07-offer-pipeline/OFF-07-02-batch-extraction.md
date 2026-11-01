# OFF-07-02 Submit flyers for extraction and collect results

| | |
|---|---|
| Service | offers |
| Epic | [OFF-07 Offer extraction, validation and publishing](README.md) |
| Release | R2 |
| Size | L |
| Depends on | OFF-07-01, OFF-05-05, OFF-01-04, AI-04 (SubmitBatch, GetBatch, GetBatchResults), AI-13 (ExtractFlyerOffers) |
| Unblocks | OFF-07-03 |

## Goal

Each covered flyer is read by AI exactly once for everyone, through Message Batches at half price (README decision 5),
and the raw results are kept until code has checked them.

## Business rules

1. A flyer version with fetch done and extraction pending is submitted as one or more batch jobs of kind
   `extract_flyer_offers`. Each item is one page (custom id `<version_id>:<page_no>`), or for a PDF one item per
   chunk of pages (custom id `<version_id>:p<first>-<last>`).
2. Each item carries: page image or PDF bytes, the flyer's language (`pl`, `de`), chain name, currency, flyer start and
   end date, and the page number. Nothing about users.
3. Items are grouped into batch jobs so that each `SubmitBatch` call stays under the request size limit set by
   AI-04; the chunk size is a setting (`EXTRACT_PAGES_PER_CALL`, default 8 pages).
4. Urgent flyers skip the batch (Decided, Q154): when a version is fetched less than 36 hours before its start date
   (or after it started), pages are sent to the per-request `ExtractFlyerOffers` instead, at full price, so offers
   are not late.
5. The version's `extract_state` becomes `submitted`; batch ids are stored. A job polls every 10 minutes
   (`GetBatch`); when a batch has ended, all results are read with `GetBatchResults` (paged) and stored as raw rows
   in `extracted_offers`, keyed by custom id. Results are read within 24 h of the batch ending (they expire on the
   provider side after 29 days).
6. Items that errored or expired are resubmitted once in a new batch. A second failure marks the version
   `extract_state = failed`, the flyer `failed`, and sends an alert `offers.extract.<source_id>`.
7. A refusal or an empty result for a page is recorded as zero offers for that page, not as an error (some pages hold
   no products).
8. When all items of a version have results, the version is `done` and OFF-07-03 runs for it.
9. The token cost returned for each job (AI-02) is stored per version, for the metric
   `offers_extraction_cost_eur_total{chain}` (the fixed weekly cost of the flyer release, business model §6).
10. `ai` down or over its daily spend limit (`RESOURCE_EXHAUSTED`): submission is retried every 30 minutes; after 12 h
    an alert is sent. Users see the flyer status rules of OFF-05-06 meanwhile.

## Contract

Uses `AiService.SubmitBatch`, `GetBatch`, `GetBatchResults` (AI-04) and `ExtractFlyerOffers` (AI-13). The result
fields offers needs from each extracted offer: product name, brand, pack text, price text, price, regular price text
and value, conditions text, valid from/to as printed, page number, kind (`food`, `drink`, `non_food`), confidence.
The authoritative message is defined in AI-13.

## Data

`offers/migrations/000NN_extraction.sql`

```sql
CREATE TABLE extraction_jobs (
    batch_id     text PRIMARY KEY,
    version_id   uuid REFERENCES flyer_versions (id) ON DELETE CASCADE,
    user_flyer_id uuid,
    items        int NOT NULL,
    state        text NOT NULL DEFAULT 'submitted' CHECK (state IN ('submitted', 'ended', 'collected', 'failed')),
    submitted_at timestamptz NOT NULL DEFAULT now(),
    ended_at     timestamptz,
    cost_eur     numeric(10,4)
);
CREATE TABLE extracted_offers (
    id          bigserial PRIMARY KEY,
    version_id  uuid REFERENCES flyer_versions (id) ON DELETE CASCADE,
    user_flyer_id uuid,
    custom_id   text NOT NULL,
    item_index  int NOT NULL,
    raw         jsonb NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (custom_id, item_index)
);
```

`extracted_offers` rows are deleted 14 days after the version's offers are published (they hold no personal data).

## Errors

| Case | Behaviour |
|---|---|
| `SubmitBatch` `INVALID_ARGUMENT` (too large) | Chunk halved and resubmitted; at 1 page still too large → version failed, alert |
| `RESOURCE_EXHAUSTED` | Retry every 30 min, alert after 12 h |

## Acceptance criteria

- [ ] With the stub provider, a 24-page fixture version is submitted as 3 jobs of 8 pages and collected into
      `extracted_offers` with one row per extracted offer.
- [ ] Polling the same ended batch twice does not duplicate rows.
- [ ] A version fetched 20 h before its start date uses `ExtractFlyerOffers`, not a batch (if the recommended rule
      is accepted).
- [ ] Cost per version is stored and the metric increases.
