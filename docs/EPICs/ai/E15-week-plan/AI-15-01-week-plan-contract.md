# AI-15-01 ComposeWeekPlan contract

| | |
|---|---|
| Service | ai |
| Epic | [AI-15 Week plan](README.md) |
| Release | R3 |
| Size | S |
| Depends on | AI-12-01 |
| Unblocks | AI-15-02; PLN-14 |

## Goal

A request and response shape for a 7-day plan where every meal is either cooked, a leftover of an earlier cooked
meal, or a no-cook dish.

## Business rules

1. Input: `CompositionContext` (its `max_items_to_buy` is per week here, 0–40; the week plan feeds the shopping
   list), `days` (1–7 dates in order, each with its `MealRequest`s), `rhythm` (`EVERY_DAY`, `EVERY_OTHER_DAY`,
   `TWICE_A_WEEK`), `preferred_cooking_dates` (optional), `eaters` (optional, AI-12-01; absent for single-person or
   no-goal households), `eat_lighter`, `max_new_dishes` (default 7, limits output size and cost).
2. Output per day and meal a `WeekSlot`: kind `COOK` (dish + `double_batch` flag), `LEFTOVER` (points to the date
   and meal of the `COOK` slot it reuses; may add `freeze=true` when it was frozen), or `NO_COOK` (a no-cooking dish);
   plus stock uses and items to buy on the slot where they are spent.
3. No calorie numbers, prices, times totals or session counts in output text; planner computes "3 cooking sessions
   this week, about 2 hours at the stove" from the slots.

## Contract

```proto
rpc ComposeWeekPlan(ComposeWeekPlanRequest) returns (ComposeWeekPlanResponse);

enum CookingRhythm { COOKING_RHYTHM_UNSPECIFIED = 0; COOKING_RHYTHM_EVERY_DAY = 1;
  COOKING_RHYTHM_EVERY_OTHER_DAY = 2; COOKING_RHYTHM_TWICE_A_WEEK = 3; }

message WeekDay {
  string date = 1;                     // YYYY-MM-DD
  repeated MealRequest meals = 2;
}

message ComposeWeekPlanRequest {
  CompositionContext context = 1;
  repeated WeekDay days = 2;           // 1-7
  CookingRhythm rhythm = 3;
  repeated string preferred_cooking_dates = 4;
  repeated Eater eaters = 5;           // AI-12-01, optional
  bool eat_lighter = 6;
  int32 max_new_dishes = 7;            // default 7
  repeated string exclude_dish_ids = 8;
  int32 variation = 9;
  bool bypass_cache = 10;
}

enum SlotKind { SLOT_KIND_UNSPECIFIED = 0; SLOT_KIND_COOK = 1; SLOT_KIND_LEFTOVER = 2; SLOT_KIND_NO_COOK = 3; }

message SlotRef { string date = 1; MealSlot meal = 2; }

message WeekSlot {
  string date = 1;
  SlotKind kind = 2;
  FamilyMeal meal = 3;                 // dish, eaters, toppings (COOK and NO_COOK)
  bool double_batch = 4;               // COOK only
  SlotRef leftover_of = 5;             // LEFTOVER only
  bool freeze = 6;                     // LEFTOVER from the freezer
}

message ComposeWeekPlanResponse {
  repeated WeekSlot slots = 1;
  CallInfo call = 2;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Request breaks the limits above | `INVALID_ARGUMENT`, `ai_bad_request` |

## Acceptance criteria

- [ ] Proto compiles; a stub response for 7 days × 3 meals passes schema checks.
