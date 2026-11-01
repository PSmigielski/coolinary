# PLN-03-02 Ingredient name resolution

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-01-02, KIT-03 (MatchNames) |
| Unblocks | PLN-03-03, PLN-03-04, PLN-08-04 |

## Goal

Turn ingredient names written by a person or by AI ("chicken thigh", "Hähnchenschenkel") into kitchen catalogue ids,
using kitchen's matching (decision 13), with clear thresholds for what is accepted automatically.

## Business rules

1. Names are normalised (trim, lower case, collapse spaces) and de-duplicated per request, then sent to
   `KIT-03 MatchNames` in one call per 100 names with `language` and context `recipe`.
2. Thresholds on the returned confidence: ≥ 0.90 accepted; 0.60–0.89 accepted and the line is flagged
   `low_confidence_match` for the reviewer; < 0.60 or "none" → unresolved.
3. An unresolved name blocks the dish (no ingredient id, no dish). The import or batch report lists each unresolved
   name once with the dishes that use it, so the admin can add it to the catalogue (ADM-03) or write an `id` instead.
4. A local map `recipe_name_matches` (language, normalised name → ingredient id, confidence, matched_at) avoids asking
   kitchen twice for the same name. Entries older than 30 days, or touched by `catalogue.ingredient_merged`, are
   refreshed on next use.
5. Decided (Q091): accept ≥ 0.90, review 0.60–0.90, reject < 0.60; the thresholds are config values, tuned after
   the first 200 dishes.

## Contract

Calls `KitchenService.MatchNames(names, language, context="recipe")`. Decided (Q090): kitchen (KIT-03) accepts the
context value `recipe` next to product, receipt and offer.

## Data

`planner/migrations/000NN_recipe_name_matches.sql`

```sql
CREATE TABLE recipe_name_matches (
    language       text NOT NULL,
    name           text NOT NULL,
    ingredient_id  text,                -- NULL = no match
    confidence     real NOT NULL,
    matched_at     timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (language, name)
);
CREATE INDEX recipe_name_matches_ingredient ON recipe_name_matches (ingredient_id);
```

## Events

`catalogue.ingredient_merged`: rows with `from_ingredient_id` are rewritten to `into_ingredient_id` (handler shared with
PLN-04-02).

## Errors

| Case | Behaviour |
|---|---|
| Kitchen unavailable | The whole import/batch step fails with UNAVAILABLE and can be re-run; nothing is half-written |

## Acceptance criteria

- [ ] 300 names with 120 duplicates make 2 `MatchNames` calls (180 unique names, 100 per call).
- [ ] A name at confidence 0.75 is accepted and flagged; at 0.5 it is unresolved and blocks its dish.
- [ ] Re-importing the same file makes no `MatchNames` calls (all names cached).
