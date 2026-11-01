# PLN-02-03 Derived diet and allergen flags

| | |
|---|---|
| Service | planner |
| Epic | [PLN-02 Recipe library model](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-02-01, PLN-01-02, KIT-02 (GetIngredients) |
| Unblocks | PLN-05-02, PLN-07-03, PLN-03-07 |

## Goal

A dish's diets (vegetarian, vegan, gluten-free, lactose-free, nut-free) and allergens are computed from its
ingredients, never typed by hand, so the hard limits in the profile can be trusted.

## Business rules

1. Kitchen's `GetIngredients` gives each ingredient its diet flags and its EU allergens (the 14 groups of Regulation
   1169/2011 Annex II: cereals with gluten, crustaceans, eggs, fish, peanuts, soybeans, milk, nuts, celery, mustard,
   sesame, sulphites, lupin, molluscs).
2. A dish version **is** a diet when every non-optional `main` and `base` line has that flag, and for every
   `topping_group` at least one option has it. Example: tacos with a beef and a bean topping are vegetarian-capable.
3. Planner stores two sets per version: `diets` (whole dish as listed, all non-optional lines including every topping)
   and `diets_capable` (rule 2). Browse and rule checks use `diets_capable` only for `base_plus_toppings` dishes,
   together with the per-plate topping choice (PLN-07-03).
4. Allergens = union over all non-optional lines (all toppings included). Optional lines' allergens are stored
   separately as `optional_allergens` so the app can say "contains milk if you add the cheese".
5. Flags are recomputed when a version is written, when `catalogue.ingredient_merged` rewrites an id (PLN-04-02), and
   on a daily job that re-reads ingredients changed in kitchen (see Notes).
6. If kitchen has no flags for an ingredient (`flags_known=false`), the dish gets no diet that depends on it and the
   allergen set is marked `allergens_complete=false`; such dishes are never shown to a profile with allergens
   (PLN-05-02, PLN-07-03).

## Contract

Fields 17–18 of `Dish` (PLN-02-01) plus:

```proto
message Dish {
  // ...
  repeated coolinary.common.v1.Diet diets_capable = 20;
  repeated coolinary.common.v1.Allergen optional_allergens = 21;
  bool allergens_complete = 22;
}
```

## Data

`planner/migrations/000NN_dish_flags.sql`

```sql
ALTER TABLE dish_versions
    ADD COLUMN diets_capable      text[]  NOT NULL DEFAULT '{}',
    ADD COLUMN optional_allergens text[]  NOT NULL DEFAULT '{}',
    ADD COLUMN allergens_complete boolean NOT NULL DEFAULT false;
CREATE INDEX dish_versions_allergens ON dish_versions USING gin (allergens);
CREATE INDEX dish_versions_diets ON dish_versions USING gin (diets_capable);
```

## Errors

| Case | Behaviour |
|---|---|
| Kitchen unavailable at write time | Write fails with UNAVAILABLE (flags are required) |
| Kitchen unavailable in the daily job | Job retries at the next run; old flags stay |

## Acceptance criteria

- [ ] Pasta with parmesan (milk) and a gluten ingredient: allergens = {gluten, milk}; diets = {vegetarian}.
- [ ] Tacos with beef and beans toppings: `diets` lacks vegetarian, `diets_capable` has it.
- [ ] An optional walnut garnish puts `nuts` in `optional_allergens` only.
- [ ] An ingredient with unknown flags makes `allergens_complete=false`.

## Notes

Kitchen publishes no event when an ingredient's allergen or diet flags change (only merges and nutrition changes are
in the README). Decided (Q054): kitchen publishes `catalogue.ingredient_flags_changed` with `ingredient_ids`
(KIT-02-10), and planner re-derives the flags of every dish using those ingredients at once; the daily job stays as a
safety net.
