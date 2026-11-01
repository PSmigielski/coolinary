# KIT-02-08 AI translation batch job

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-02-04, KIT-01-08, AI-09 (`translate_ingredients` batch kind) |
| Unblocks | new languages (wave 2: French, Spanish, Italian), ADM-03 (review of machine names) |

## Goal

Fill missing names in every catalogue language once for everyone at batch price, so adding an ingredient or a
language does not mean typing translations by hand (features.md §2 "Ingredient catalogue translated … (once for
everyone)").

## Business rules

1. Languages come from config `CATALOGUE_LANGUAGES` (launch: `en,pl,de`).
2. A daily job at 03:30 UTC selects active ingredients missing a primary name in any configured language, ordered by
   `popularity` descending, at most 1,000 ingredients per batch.
3. Each item sends: slug, category, every existing name with its language and kind, and the target languages. It
   asks for a primary name, up to 5 synonyms and a plural per target language. No user data is sent.
4. Results are checked without AI before insert:
   - 1–60 characters, letters of the target language's script plus spaces, hyphens and apostrophes;
   - not equal to the primary name of another ingredient in the same category and language (such a name is
     skipped and listed for review);
   - at most 5 synonyms; duplicates after normalisation dropped.
5. Inserted names get `origin = 'machine'`, `reviewed = false`. They are used at once for display, search and
   matching.
6. Decided (Q066): machine names go live before review; the native-speaker check per language (business model)
   follows and marks them `reviewed`.
7. Items that fail three times are parked and shown in the admin review list (ADM-03).

## Contract

Uses ai `SubmitBatch` / `GetBatchResults` with kind `translate_ingredients` (AI-04, AI-09). No kitchen RPC.

## Data

`kitchen/migrations/000NN_translation_queue.sql`

```sql
CREATE TABLE translation_queue (
    ingredient_id  uuid PRIMARY KEY REFERENCES ingredients (id) ON DELETE CASCADE,
    attempts       integer NOT NULL DEFAULT 0,
    parked         boolean NOT NULL DEFAULT false,
    last_error     text,
    queued_at      timestamptz NOT NULL DEFAULT now()
);
```

## Errors

None returned to users; failures are counted in `ai_batch_items_total{kind="translate_ingredients"}`.

## Acceptance criteria

- [ ] An ingredient with only an `en` name gets `pl` and `de` primaries after one batch (stub ai).
- [ ] A suggested name that collides with another ingredient's primary in the same category is not inserted and
      appears in the review list.
- [ ] Adding `fr` to `CATALOGUE_LANGUAGES` queues every active ingredient for French.
- [ ] No seed or admin name is changed by the job.
