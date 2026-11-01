# PLN-03-05 Dish translation batch (`translate_dishes`)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-03-04 (planner_batches, poll job), AI-04, AI-09 (`translate_dishes`) |
| Unblocks | PLN-03-07 (publish needs pl, de, en) |

## Goal

Every library dish exists in Polish, German and English. Missing languages are translated once for everyone in a
batch and marked as machine translated until a person checks them.

## Business rules

1. `TranslateDishes(dish_ids | all_missing=true, target_languages)`: picks review and published versions that lack any
   target language; ≤ 500 versions per batch.
2. Source language = a text with `machine_translated=false` (prefer `en`, then `pl`, then `de`).
3. Each item sends name, intro, step texts and ingredient notes; ingredient names are not translated here (they come
   from the catalogue).
4. Checks on each result, without AI: same number of steps; every number in a source step appears in the translated
   step (temperatures, minutes); no calorie number (PLN-02-02 rule 10); name ≤ 120 characters. A failing translation
   is dropped and reported.
5. Accepted results are written as `dish_texts` rows with `machine_translated=true` into the same version; for a
   published version this changes only texts, so no new version is created.
6. A reviewer can edit a text and clear `machine_translated` (PLN-03-07).
7. Decided (Q122): dishes may be published with machine-translated texts after a native speaker spot-checks 10% of
   dishes per language, before each language goes live (business model "one-off per language").
8. `dry_run = true` selects exactly as rule 1 and returns `items` without submitting a batch or writing
   anything. Admin turns the count into a cost estimate for the founder to confirm (ADM-04-04).

## Contract

```proto
// PlannerAdminService
rpc TranslateDishes(TranslateDishesRequest) returns (TranslateDishesResponse);
message TranslateDishesRequest { repeated string dish_ids = 1; bool all_missing = 2; repeated string target_languages = 3; bool dry_run = 4; }
message TranslateDishesResponse { string batch_id = 1; int32 items = 2; }
```

Progress and report through `GetDraftBatch` / `ListDraftBatches` (shared batch views).

## Data

No new table (uses `planner_batches`, `dish_texts`).

## Errors

| Case | gRPC status / behaviour |
|---|---|
| Neither `dish_ids` nor `all_missing`, or both | `INVALID_ARGUMENT` |
| A target language outside the launch languages (pl, de, en) | `INVALID_ARGUMENT` |
| More than 500 versions selected | `INVALID_ARGUMENT` (`batch_too_large`); the caller splits the run |
| An unknown dish id | `NOT_FOUND` naming the id; nothing is submitted |
| Nothing to translate | success with `items = 0` and no batch |
| ai down or over its spend limit at submission | `UNAVAILABLE` / `RESOURCE_EXHAUSTED` passed through; nothing stored |
| A result fails rule 4 | dropped and listed in the batch report; other results are kept |

## Acceptance criteria

- [ ] A dish with only `en` gets `pl` and `de` rows with `machine_translated=true`.
- [ ] A translation that drops "180 °C" from a step is rejected.
- [ ] Running it twice in a row submits nothing the second time (no missing languages).
- [ ] `dry_run = true` returns the same `items` the real call then submits, and creates no `planner_batches` row.
