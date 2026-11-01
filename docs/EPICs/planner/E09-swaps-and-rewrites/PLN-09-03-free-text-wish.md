# PLN-09-03 Free-text wish for a meal

| | |
|---|---|
| Service | planner |
| Epic | [PLN-09 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-09-01, PLN-07-05, AI-11 (ModifyMeal) |
| Unblocks | GW-07 |

## Goal

"Something warm in 20 minutes": the user's own words steer one meal, while every hard limit still applies.

## Business rules

1. `WishMeal(plan_id | date, meal, text)`; text 1–200 characters after trimming. Plus only (Free →
   FAILED_PRECONDITION `plus_required`).
2. The text is user data: planner sends it to `ModifyMeal(kind=WISH)` in its own field; it never goes into logs,
   metrics, plans or dishes. Planner stores only a SHA-256 of the text for idempotency.
3. A wish can't loosen a hard limit: the answer runs through the full checker. "Ignore my nut allergy" still gets a
   nut-free dish or fails.
4. If AI returns a time limit read from the wish (`time_limit_minutes`, AI-11), the checker uses it instead of the
   meal's budget (PLN-07-05 rule 5); the limit can only be shorter than the profile's budget + 60 minutes.
5. Stock mode `FROM_STOCK` with the remaining stock (PLN-09-01 rule 2). AI may say the wish can't be met from stock
   and propose a dish needing up to 3 purchases; that answer is shown as a `BUY_FEW` option (PLN-08-06).
6. Without a plan (`date` given, no plan yet), the wish creates a proposed plan with only that meal.
7. Failing answer → FAILED_PRECONDITION `wish not possible`, with the main violation code (e.g. `TIME_OVER`), so the
   app can say why.
8. The response says `ai_used = true` when the meal came from AI, so the gateway counts Plus fair use (GW-03-04)
   and AI metrics (GW-03-05).

## Contract

```proto
rpc WishMeal(WishMealRequest) returns (WishMealResponse);
message WishMealRequest { string household_id = 1; string plan_id = 2; string date = 3; Meal meal = 4; string text = 5;
                          string request_id = 6; string language = 7; }
message WishMealResponse { PlanMeal meal = 1; string plan_id = 2; bool ai_used = 3; }
```

## Data

No schema change (the text hash lives with the idempotency key in `plans.request_id`/a request table).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Message |
|---|---|---|
| Text empty or > 200 chars | INVALID_ARGUMENT | `wish must be 1-200 characters` |
| Free household | FAILED_PRECONDITION | `plus_required` |
| No dish passes | FAILED_PRECONDITION | `wish not possible: <code>` |

## Acceptance criteria

- [ ] "something warm in 20 minutes" yields a dish with total time ≤ 20 + 10 minutes, or `wish not possible: TIME_OVER`.
- [ ] A wish asking to ignore an allergy still returns an allergen-free dish or fails.
- [ ] Grepping planner logs after a wish finds no part of the wish text.
- [ ] A successful wish returns `ai_used = true` (rule 8).
