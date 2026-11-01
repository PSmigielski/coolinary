# PLN-12-08 Scaled ingredients for a planned meal

| | |
|---|---|
| Service | planner |
| Epic | [PLN-12 Per-person portions](README.md) |
| Release | R2 |
| Size | S |
| Depends on | PLN-12-06, PLN-02 |
| Unblocks | PLN-15-02, PLN-13-02, PLN-14-03, PLN-10 |

## Goal

Ingredients of a planned meal are scaled to its pot, so cooking mode and the shopping list use the right amounts
(features.md step 4, "Shopping list" step 1).

## Business rules

1. Scale factor = pot portions ÷ the dish's base servings (PLN-02). Without goals (no pot stored) the factor
   is household portions ÷ base servings, as in R1.
2. Every ingredient amount is multiplied by the factor; units are kept; amounts are not rounded here (the list
   rounds to packs in PLN-15-04, cooking mode rounds for display).
3. Optional ingredients are included and marked optional.
4. Base-plus-toppings dishes: base ingredients scale with the pot; each topping scales with the sum of portions of
   the eaters who take it (PLN-13-02).
5. Leftover meals (second half of a double cook, PLN-14-03) return no ingredients; their session's cooking meal
   carries the doubled pot.
6. Extras of the meal are returned separately, per person, not mixed into the dish's ingredients.
7. A pot of 0 (everyone away) returns an empty list.

## Contract

`contracts/coolinary/planner/v1/portions.proto`

```proto
rpc GetScaledIngredients(GetScaledIngredientsRequest) returns (ScaledIngredients);
message GetScaledIngredientsRequest { string household_id = 1; string plan_meal_id = 2; }
message ScaledIngredients {
  string plan_meal_id = 1;
  double factor = 2;
  int32 pot_portions = 3;
  repeated ScaledIngredient ingredients = 4;
  repeated MealExtra extras = 5;           // from PLN-12-06
  bool leftover = 6;
}
message ScaledIngredient {
  string ingredient_id = 1; coolinary.common.v1.Quantity amount = 2;
  bool optional = 3; string group = 4;     // "base", "topping", "" (PLN-02)
}
```

Also exposed as an internal Go function used by PLN-15-02 and cooking mode (PLN-10).

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Plan meal not in household | `NOT_FOUND` |

## Acceptance criteria

- [ ] Chicken curry with 4 base servings and pot 6: every amount × 1.5.
- [ ] Ola away (pot 5): × 1.25, and the shopping list follows (Verification 17).
- [ ] Leftover meal returns `leftover = true` and no ingredients.
- [ ] Tacos where only two of four eaters take the meat topping: meat scales with those two eaters' portions.
