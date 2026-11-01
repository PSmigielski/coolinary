# PLN-12-04 No-cook extras

| | |
|---|---|
| Service | planner |
| Epic | [PLN-12 Per-person portions](README.md) |
| Release | R2 |
| Size | M |
| Depends on | PLN-12-03, KIT-06 (ComputeNutrition), KIT-07 (ListStock), KIT-02 (GetIngredients) |
| Unblocks | PLN-12-05, PLN-13-04, PLN-15-02 |

## Goal

When the 2-portion limit leaves someone more than 15% short, add no-cook extras (bread, fruit, yoghurt, milk, cheese,
nuts, raw vegetables), taken from what is at home first (features.md "How one meal is split", step 5).

## Business rules

1. Extras come from the planner table `no_cook_extras`: ingredient id, one serving (amount and unit), a label key
   ("a glass of milk"), a default rank, a lighter rank (null = not used when eating lighter), and `child_ok`.
2. Seed rows (default rank): bread 1 slice 35 g (1), fruit: banana 1 piece / apple 1 piece (2), natural yoghurt
   150 g (3), milk 250 ml (4), hard cheese 30 g (5), mixed nuts 30 g (6), raw carrot or pepper sticks 100 g (7).
   Lighter ranks: fruit 1, raw vegetables 2, milk 3, yoghurt 4; bread, cheese and nuts null (features.md "eat
   lighter as a family").
3. Serving kcal comes from KIT-06 `ComputeNutrition` and is cached for 24 hours in redis-cache. An extra whose
   kcal is incomplete is skipped.
4. Candidates are filtered by the eater's hard limits (allergens, diet, dislikes) using ingredient flags from
   KIT-02 `GetIngredients`. For child profiles, rows with `child_ok = false` (nuts) are skipped.
5. Order: candidates at home first (stock quantity, after reservations by earlier meals of the same plan, covers one
   serving), then the rest; within each group by default rank, or by lighter rank when the household eats lighter.
6. Greedy fill: while the remaining gap is more than 5% of the eater's share and fewer than 3 extras were chosen,
   take the next candidate whose serving kcal ≤ remaining gap + 10% of the share. Each candidate at most once per
   eater per meal.
7. If the remaining gap is still more than 15% of the share, the eater is flagged `still_short`; the rest carries
   over (PLN-12-05).
8. Extras used from stock reserve that stock for this plan; extras not at home become shopping-list contributions of
   kind `extra` (PLN-15-02).
9. Extras are shown per person on the dish card ("Kuba: a glass of milk and a banana, about 260 kcal"); the kcal
   number follows the visibility rules (PLN-12-07).

## Contract

Internal package `planner/internal/portions` (`ChooseExtras(gap, share, eater, stock, lighter) []Extra`).
Admin edits through `PlannerAdminService`:

```proto
rpc ListNoCookExtras(ListNoCookExtrasRequest) returns (ListNoCookExtrasResponse);
rpc UpsertNoCookExtra(NoCookExtra) returns (NoCookExtra);
message NoCookExtra {
  string ingredient_id = 1; coolinary.common.v1.Quantity serving = 2; string label_key = 3;
  int32 rank_default = 4; optional int32 rank_lighter = 5; bool child_ok = 6; bool active = 7;
}
```

## Data

`planner/migrations/000NN_no_cook_extras.sql`

```sql
CREATE TABLE no_cook_extras (
    ingredient_id text PRIMARY KEY,           -- kitchen catalogue id
    serving_value double precision NOT NULL CHECK (serving_value > 0),
    serving_unit  smallint NOT NULL,          -- common.v1.Unit
    label_key     text NOT NULL,
    rank_default  smallint NOT NULL,
    rank_lighter  smallint,
    child_ok      boolean NOT NULL DEFAULT true,
    active        boolean NOT NULL DEFAULT true
);
CREATE TABLE meal_extras (
    plan_meal_id  uuid NOT NULL,
    member_key    text NOT NULL,              -- member id, or "unnamed:<n>"
    ingredient_id text NOT NULL,
    serving_value double precision NOT NULL,
    serving_unit  smallint NOT NULL,
    kcal          integer NOT NULL,
    from_stock    boolean NOT NULL,
    PRIMARY KEY (plan_meal_id, member_key, ingredient_id)
);
```

## Errors

| Case | Behaviour |
|---|---|
| kitchen down | no extras chosen, eater flagged `still_short`, result marked `extras_unavailable` |
| No candidate passes the eater's limits | `still_short` |

## Acceptance criteria

- [ ] Kuba (share 1,280, gap 240) with only milk and bananas at home: banana (~105 kcal), then milk (~160): total
      about 265 kcal, displayed "about 260" after rounding; nothing goes on the list (worked example).
- [ ] Same with an empty fridge: bread then fruit are chosen and both appear on the list as `extra` lines.
- [ ] Eating lighter: bread, cheese and nuts are never chosen.
- [ ] Nut allergy or child profile: nuts never chosen.
- [ ] Never more than 3 extras per person per meal.

## Notes

- Decided (Q199): the default order (rule 2) is bread, fruit, yoghurt, milk, cheese, nuts, raw vegetables, with the
  lighter order for "eat lighter as a family".
- Decided (Q200): nuts are never offered to child profiles (choking risk for small children; age is not on the
  server).
