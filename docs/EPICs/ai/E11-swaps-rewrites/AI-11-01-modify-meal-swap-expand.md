# AI-11-01 ModifyMeal: contract, SWAP and EXPAND_ALTERNATIVE

| | |
|---|---|
| Service | ai |
| Epic | [AI-11 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-10-03, AI-10-04 |
| Unblocks | AI-11-02, AI-11-03, AI-11-04; PLN-09 |

## Goal

Replace one meal of an existing day plan, or write the steps of an alternative the user picked, while the other
meals keep the stock they already use.

## Business rules

1. Job `modify_meal`, Sonnet 5.5 at effort `low` (routing by plan: AI-10-05). Prompt `modify_meal/v1`.
2. Input: `CompositionContext`; `other_meals` (the day's other `ComposedMeal`s, with their `stock_uses`); `target`
   (`MealRequest`); `mode`; and per mode the fields below. Stock available to the target = stock − the other
   meals' `stock_uses`; ai computes it and sends it to the model as the stock for this call.
3. `SWAP`: a different dish for the target meal, following every rule of AI-10-02 rule 2. `current_dish_id` /
   `current_dish_name` and `exclude_dish_ids` are excluded. Output: one `ComposedMeal` (with cuisine fit and
   alternatives as AI-10-04).
4. `EXPAND_ALTERNATIVE`: input `alternative` (a `MealAlternative` from an earlier response). Output: the same dish
   with full steps; ingredients, `stock_uses` and `to_buy` must be unchanged (validated field by field); no
   alternatives. This mode is allowed on Free, because it completes a plan already counted.
5. Validation: AI-10-03 on the single meal against the reduced stock, plus mode rules (rule 4 for expand).
6. Cache TTL 1 hour; key includes mode, target, other meals' stock uses and excludes.

## Contract

```proto
rpc ModifyMeal(ModifyMealRequest) returns (ModifyMealResponse);

enum ModifyMode { MODIFY_MODE_UNSPECIFIED = 0; MODIFY_MODE_SWAP = 1; MODIFY_MODE_REPLACE_INGREDIENT = 2;
  MODIFY_MODE_WISH = 3; MODIFY_MODE_EXPAND_ALTERNATIVE = 4; }

message ModifyMealRequest {
  CompositionContext context = 1;
  repeated ComposedMeal other_meals = 2;
  MealRequest target = 3;
  ModifyMode mode = 4;
  // SWAP
  string current_dish_id = 5;
  string current_dish_name = 6;
  repeated string exclude_dish_ids = 7;
  // REPLACE_INGREDIENT (AI-11-02)
  ComposedMeal current_meal = 8;
  string missing_ingredient_id = 9;
  // WISH (AI-11-03)
  string wish = 10;
  // EXPAND_ALTERNATIVE
  MealAlternative alternative = 11;
  int32 variation = 12;
  bool bypass_cache = 13;
}

message ModifyMealResponse {
  ComposedMeal meal = 1;                           // SWAP, WISH, EXPAND_ALTERNATIVE
  repeated IngredientReplacement replacements = 2; // REPLACE_INGREDIENT (AI-11-02)
  CallInfo call = 3;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Mode fields missing (e.g. EXPAND without `alternative`) | `INVALID_ARGUMENT`, `ai_bad_request` |
| Others | AI-01-04 |

## Acceptance criteria

- [ ] SWAP of dinner never uses stock the breakfast and supper already use beyond what is left.
- [ ] EXPAND of a `BUY_ITEMS` alternative returns steps and the same ingredients, stock uses and items to buy; a fake answer that changes an ingredient fails validation.
- [ ] Golden set thresholds pass (AI-11-06).
