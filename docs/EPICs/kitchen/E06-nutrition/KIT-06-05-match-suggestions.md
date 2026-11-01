# KIT-06-05 Ingredient → nutrition entry suggestions

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-06-01, KIT-06-02 (at least one importer), KIT-01-08, AI-09 (`match_nutrition_entries`) |
| Unblocks | KIT-06-06 |

## Goal

Each catalogue ingredient gets a proposed nutrition table entry, prepared by AI once for everyone at batch price and
checked by a person before use (features.md §2 "Calories and goals, With AI … checked by hand before use").

## Business rules

1. A daily job at 04:30 UTC queues active ingredients with no approved link and no open suggestion, highest
   `popularity` first, at most 2,000 per batch.
2. Candidates per ingredient: the 5 best entries per source by trigram similarity of the ingredient's `en` name with
   `name_en`, plus, for BLS, its `de` name with `name`; only non-retired entries with energy present.
3. The batch item sends the ingredient's names, category and base unit and the candidates (source, code, names,
   state, energy and macros). The answer is one candidate id (or none), a confidence 0–1 and a short reason code.
4. Checks without AI:
   - the chosen id must be one of the candidates, else no suggestion;
   - recipe weights are raw or dry (features.md §5 step 3): for categories `meat`, `poultry`, `fish_seafood`,
     `grains_pasta_rice` and `legumes`, a chosen entry with `state = 'cooked'` or `'prepared'` gets
     `warning = 'cooked_state'`;
   - a chosen entry whose energy differs by more than 50% from the median of all candidates gets
     `warning = 'energy_outlier'`.
5. Source preference when the answer is unclear (confidence below 0.6): the best candidate from Ciqual, then BLS,
   then USDA, with `method = 'ai'` and `warning = 'low_confidence'`.
6. The result is stored as `status = 'suggested'`. Nothing is used for calculation until approved (KIT-06-06).
7. Items that fail three times are parked and listed in admin as "needs manual match".

## Contract

Uses ai batch kind `match_nutrition_entries` (AI-04, AI-09). No kitchen RPC.

## Data

`kitchen/migrations/000NN_nutrition_queue.sql`

```sql
CREATE TABLE nutrition_match_queue (
    ingredient_id  uuid PRIMARY KEY REFERENCES ingredients (id) ON DELETE CASCADE,
    attempts       integer NOT NULL DEFAULT 0,
    parked         boolean NOT NULL DEFAULT false,
    queued_at      timestamptz NOT NULL DEFAULT now()
);
```

## Errors

None returned to users.

## Acceptance criteria

- [ ] With a stub ai, every queued ingredient ends with one suggestion or a parked row.
- [ ] "Pasta" matched to a cooked-pasta entry carries `cooked_state`.
- [ ] An answer naming a non-candidate creates no suggestion.
