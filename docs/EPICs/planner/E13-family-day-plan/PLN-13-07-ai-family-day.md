# PLN-13-07 Family day with AI

| | |
|---|---|
| Service | planner |
| Epic | [PLN-13 Family day plan](README.md) |
| Release | R2 |
| Size | L |
| Depends on | PLN-13-06, AI-12 (ComposeFamilyDay), PLN-07, PLN-08, PLN-12-06, PLN-12-07 |
| Unblocks | GW-07, PLN-14-07 |

## Goal

Compose a family day with one dish per meal that fits every person's target, from stock first (features.md §2
"Calories and goals" with AI). AI proposes; rules check; portions are calculated without AI.

## Business rules

1. Plus only: the plan from identity (PLN-01-05) must be `plus`, else `FAILED_PRECONDITION` (`plus_required`). The gateway's quota does not
   apply (Plus is unlimited).
2. Request to `AI-12 (ComposeFamilyDay)`:
   - stock: ingredient ids, quantities, days to expiry (negative and `past_date` when past the date, Q103), sorted
     oldest first; meals using such items carry `check_first_stock_item_ids` as in PLN-08-03;
   - per meal: eaters as anonymous keys (`p1`, `p2`, …), each with their kcal budget for that meal (PLN-12-02) and
     hard limits; no names, ages, weights or goal types;
   - household tools, cuisines per meal, time per meal, "no cooking" rule for non-main meals, `eat_lighter_as_family`;
   - up to 30 library candidates per meal that pass the shared limits, so AI can pick one instead of inventing.
3. AI returns per meal: a library dish id, or a new dish (ingredient ids from the catalogue with amounts per standard
   portion, base servings, steps, tools, effort tags), plus topping groups for base-plus-toppings dishes. AI
   returns no kcal numbers; any it returns are ignored.
4. Every new dish is checked by PLN-07 against the shared limits, tools, time and stock; kcal comes from PLN-04;
   portions from PLN-12 (with plates from PLN-13-02). A dish that is `UNSUITABLE` for an eater fails the check.
5. If any meal fails, one retry is sent with the failing meals and their violations (no personal data). If it fails
   again, that meal is filled by the non-AI composer (PLN-13-06) and the plan says `partial_fallback`.
6. If `ai` returns `UNAVAILABLE`, `RESOURCE_EXHAUSTED` (spend limit) or a refusal, the whole day comes from
   PLN-13-06 with `ai_used = false` and `fallback_reason`.
7. The AI plan and the non-AI plan are compared with PLN-13-03; the AI plan is kept unless it is worse on keys 1–3,
   in which case the non-AI plan is returned (`ai_plan_rejected`).
8. New AI dishes are stored as household-private dishes (`source = ai`, PLN-02/PLN-08) so swaps and "cooked it" work;
   they never enter the shared library without review (PLN-11).
9. Saving replaces any existing plan for that date only after the user confirms (`replace = true`); otherwise
   `ALREADY_EXISTS`.
10. The response passes through the visibility filter (PLN-12-07).

## Contract

`contracts/coolinary/planner/v1/family.proto`

```proto
rpc CreateFamilyDayPlan(CreateFamilyDayPlanRequest) returns (FamilyDayPlan);

message CreateFamilyDayPlanRequest {
  string household_id = 1;
  string viewer_member_id = 2;
  string date = 3;                 // YYYY-MM-DD, today … today + 13
  bool replace = 4;
  bool use_ai = 5;                 // false = library only
}
message FamilyDayPlan {
  string day_plan_id = 1;
  string date = 2;
  repeated FamilyMeal meals = 3;
  repeated PersonDay people = 4;   // day totals, filtered by PLN-12-07
  int32 cooking_sessions = 5;
  bool ai_used = 6;
  string fallback_reason = 7;      // "", "ai_unavailable", "ai_budget", "ai_refused", "rule_check_failed", "partial_fallback", "ai_plan_rejected"
}
message FamilyMeal {
  Meal meal = 1;
  repeated MealDish dishes = 2;    // PLN-13-05
  MealPortions portions = 3;       // PLN-12-06
  repeated PlateToppings plates = 4;
  bool cuisine_relaxed = 5;
}
message PersonDay { string member_id = 1; optional int32 planned_kcal = 2; optional int32 target_kcal = 3; optional bool in_range = 4; }
```

## Data

No new tables (PLN-08 plan, PLN-12-06 allocations, PLN-13-02/05 tables).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Free household | `FAILED_PRECONDITION` `plus_required` |
| Date outside today … today + 13 | `INVALID_ARGUMENT` |
| Plan exists and `replace = false` | `ALREADY_EXISTS` |
| kitchen down | `UNAVAILABLE` |

## Acceptance criteria

- [ ] Worked example household: one dish per meal, dinner pot 6, Kuba's extras, each day total within 10%.
- [ ] AI returning a dish with a nut ingredient for a nut-free household: the dish is rejected and replaced; zero
      violations shown (Verification 3).
- [ ] With `ai` stopped, the call still returns a plan with `ai_used = false` (Verification 1).
- [ ] The AI request payload contains no names, ages or weights (checked in a test against the serialized request;
      Verification 20).
- [ ] Langfuse trace tag `family_day` with plan `plus`; cost recorded by `ai`.
