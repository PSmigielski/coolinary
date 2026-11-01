# PLN-09-02 Replace a missing ingredient

| | |
|---|---|
| Service | planner |
| Epic | [PLN-09 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-09-01, PLN-06-01, AI-11 (ModifyMeal) |
| Unblocks | GW-07 |

## Goal

"I have no cream": get a version of the planned dish that uses something at home instead of the missing ingredient.

## Business rules

1. `ReplaceIngredient(plan_id, meal, ingredient_id)` where the ingredient is a non-optional line of the meal's dish
   and is missing or short (PLN-06-01).
2. **Plus only.** Free → FAILED_PRECONDITION `plus_required`; the app offers "add to shopping list" instead.
   Decided (Q124): no non-AI substitution table for Free at launch.
3. Planner sends `ModifyMeal(kind=REPLACE_INGREDIENT)` with the current dish (ingredient ids, amounts, steps), the
   ingredient to replace, the remaining stock (as in PLN-09-01 rule 2) and staples. AI may change amounts of other lines
   and steps, but may add only ingredients that are at home or staples.
4. The answer is checked in `FROM_STOCK` mode; the replaced ingredient must be gone from non-optional lines, else
   `STRUCTURE`. A failing answer → FAILED_PRECONDITION `no replacement found` (no library fallback: the user asked
   about this dish).
5. The result is stored as a new household AI dish with `parent_dish_id` = the original dish, nutrition from tables;
   the plan meal points at it.
6. One replace per (plan meal, ingredient) per `request_id`; repeats return the stored result.
7. The response says `ai_used = true` when the new dish came from AI (fresh or from ai's result cache) and
   `replayed = true` when rule 6 returned a stored result, so the gateway counts Plus fair use (GW-03-04) and AI
   metrics (GW-03-05) once.

## Contract

```proto
rpc ReplaceIngredient(ReplaceIngredientRequest) returns (ReplaceIngredientResponse);
message ReplaceIngredientRequest { string household_id = 1; string plan_id = 2; Meal meal = 3; string ingredient_id = 4;
                                   string request_id = 5; string language = 6; }
message ReplaceIngredientResponse { PlanMeal meal = 1; bool ai_used = 2; bool replayed = 3; }
```

## Data

No schema change.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Message |
|---|---|---|
| Free household | FAILED_PRECONDITION | `plus_required` |
| Ingredient not in the dish, or not missing | INVALID_ARGUMENT | `nothing to replace` |
| AI result fails the checker or AI unavailable | FAILED_PRECONDITION / UNAVAILABLE | `no replacement found` / `ai unavailable` |

## Acceptance criteria

- [ ] Replacing cream in a pasta dish with yoghurt at home yields a dish without cream whose new lines are all in stock.
- [ ] The new dish's `parent_dish_id` is the original, and the original is unchanged.
- [ ] A successful replace returns `ai_used = true`; repeating it with the same `request_id` returns `replayed = true`.
