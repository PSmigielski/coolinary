# PLN-08-08 SetPlanMeal: put a chosen dish into a plan meal

| | |
|---|---|
| Service | planner |
| Epic | [PLN-08 AI day plan](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-08-01, PLN-08-03, PLN-07-01, PLN-07-02, PLN-07-03, PLN-07-04, PLN-01-03, PLN-02-06 |
| Unblocks | GW-07 (manual meal route), PLN-15-10 (list follows plan) |

## Goal

Decided (Q117): users can put a library dish they picked (from browse, favourites or "what can I make now") into a
day-plan meal themselves. It is checked by rules, uses no AI and no Free quota.

## Business rules

1. `SetPlanMeal(plan_date, meal, dish_id, servings?)`. The dish is a published library dish or one of the household's
   own AI dishes; its latest version is used.
2. Target plan: the household's accepted plan for that date; else its newest proposed plan for that date; else a new
   `proposed` plan for that date holding only this meal (`ai_used = false`).
3. The meal is replaced in place. A meal already marked `cooked` → FAILED_PRECONDITION `meal_already_cooked`. A meal
   with an open cuisine choice (PLN-08-06) is resolved by the chosen dish; its options are cleared.
4. Eaters: the meal's existing eaters; for a new meal, the household's eaters for that meal (profile snapshot
   `EatersFor(meal)`). Servings: the request value, else the household's portions (R1); from R2 the pot from PLN-12.
5. Rule check (PLN-07) with the day's other meals, so the allocation (PLN-08-03) shares stock across meals:
   - hard limits block: `TOOL_MISSING`, `ALLERGEN`, `DIET` for any eater → FAILED_PRECONDITION
     `x-error-code: rule_violation` with the problem codes; nothing is saved;
   - the rest are returned as warnings and don't block, because the user chose the dish: `STOCK_SHORT`,
     `LOW_STOCK_MARGIN`, `TIME_OVER`, `TIME_TIGHT`, `DISLIKE`, `CUISINE_DIFFERS`, `CHECK_FIRST` (Q103);
   - missing ingredients are returned so the app can offer "add to shopping list".
6. The meal is stored with `source = manual`, the allocation's `check_first` stock item ids (Q103), status `planned`.
   The plan's `ai_used` and `fallback_reason` stay as they were. No call to ai; the response has `ai_used = false`, so
   the gateway counts no quota.
7. An accepted plan stays accepted; the shopping list follows the change (PLN-15-10).
8. Idempotent by `request_id`: the same request again returns the same result and changes nothing.

## Contract

`contracts/coolinary/planner/v1/plan.proto`

```proto
rpc SetPlanMeal(SetPlanMealRequest) returns (SetPlanMealResponse);

message SetPlanMealRequest {
  string household_id = 1;
  string plan_date = 2;          // YYYY-MM-DD, household local
  Meal meal = 3;
  string dish_id = 4;
  double servings = 5;           // 0 = household default
  string request_id = 6;
  string language = 7;
}
message SetPlanMealResponse {
  Plan plan = 1;                         // PLN-08-01
  repeated ValidationProblem warnings = 2;   // PLN-02-02 shape, rule codes of PLN-07-01
  repeated MissingLine missing = 3;      // PLN-06-03 shape
  bool ai_used = 4;                      // always false
}
```

Uses `PLAN_MEAL_SOURCE_MANUAL` from `PlanMealSource` (PLN-08-01).

## Data

No new table. `plan_meals.source` accepts `manual` (PLN-08-01).

## Errors

| Case | gRPC status / reason | Message or behaviour |
|---|---|---|
| Dish not found, retired, or another household's AI dish | NOT_FOUND | `dish not found` |
| Tool missing, allergen or diet of an eater | FAILED_PRECONDITION, `x-error-code: rule_violation` | problem codes listed; nothing saved |
| Meal already cooked | FAILED_PRECONDITION, `x-error-code: meal_already_cooked` | |
| Date more than 31 days ahead or in the past | INVALID_ARGUMENT | `plan_date out of range` |
| Servings outside ¼ steps 1–24 | INVALID_ARGUMENT | PLN-02-05 rule 3 |
| kitchen unavailable (profile or stock) | UNAVAILABLE | nothing saved |

## Acceptance criteria

- [ ] Setting a library curry as tomorrow's dinner with no plan for that date creates a `proposed` plan with one meal,
      `source = manual`, `ai_used = false`, and no quota counted by the gateway.
- [ ] Setting a dish that needs an oven in a "student room" household fails with `rule_violation` (`TOOL_MISSING`).
- [ ] A dish containing nuts for a household with a nut allergy fails with `rule_violation` (`ALLERGEN`).
- [ ] A dish missing 200 g chicken is saved, with a `STOCK_SHORT` warning and chicken in `missing`.
- [ ] Replacing the dinner of an accepted plan keeps the plan accepted and the shopping list proposes the new lines.
- [ ] A meal using milk past its date is saved with that stock item in `check_first`.
- [ ] Repeating the call with the same `request_id` changes nothing.
