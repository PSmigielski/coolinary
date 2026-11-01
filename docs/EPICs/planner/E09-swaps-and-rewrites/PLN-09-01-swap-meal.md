# PLN-09-01 Swap one meal

| | |
|---|---|
| Service | planner |
| Epic | [PLN-09 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-08-07, PLN-08-03, PLN-07-05, AI-11 (ModifyMeal) |
| Unblocks | PLN-09-02, PLN-09-03, GW-07 |

## Goal

"Not that one": replace one meal of a day plan with a different dish that still fits the stock left by the other meals,
the cuisine, the tools and everyone's limits.

## Business rules

1. `SwapMeal(plan_id, meal)` on a proposed or accepted plan whose meal is `planned` (not cooked).
2. Stock available for the swap = stock minus the allocation of the plan's other meals (PLN-08-03), so the new dish
   can't take what the others use.
3. Dishes to avoid: the current dish and every dish swapped out of this meal in the last 24 hours (kept in
   `plan_meal_swaps`).
4. **Plus**: call `ModifyMeal(kind=SWAP)` with the same request shape as the day plan for one meal (PLN-08-02), the
   remaining stock, the avoid list and library candidates. Check the answer (PLN-07, `source` metric `swap`). No repair
   call: a failing answer falls back to rule 5.
5. **Free, or AI not available**: the next best library dish for the meal against the remaining stock (PLN-08-05
   logic), not in the avoid list. No quota is used.
6. The meal's dish, plates, leave-outs and allowed purchases are replaced in place; the plan keeps its status.
7. If nothing fits (library exhausted), FAILED_PRECONDITION `no other dish fits`.
8. Response says `ai_used` (Plus) so the app can show the right label; swaps never consume Free quota.

## Contract

```proto
rpc SwapMeal(SwapMealRequest) returns (SwapMealResponse);
message SwapMealRequest { string household_id = 1; string plan_id = 2; Meal meal = 3; string request_id = 4; string language = 5; }
message SwapMealResponse { PlanMeal meal = 1; bool ai_used = 2; FallbackReason fallback_reason = 3; }
```

## Data

`planner/migrations/000NN_plan_meal_swaps.sql`

```sql
CREATE TABLE plan_meal_swaps (
    plan_meal_id uuid NOT NULL REFERENCES plan_meals (id) ON DELETE CASCADE,
    dish_id      uuid NOT NULL,
    swapped_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (plan_meal_id, dish_id)
);
```

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Meal already cooked | FAILED_PRECONDITION | `meal already cooked` |
| Nothing else fits | FAILED_PRECONDITION | `no other dish fits` |

## Acceptance criteria

- [ ] Swapping supper never proposes a dish needing the eggs allocated to dinner.
- [ ] Three swaps in a row give three different dishes, none of them the original.
- [ ] Free household: swap answers from the library with `ai_used=false` and no ai call.
- [ ] Repeating the same `request_id` returns the same new dish.
