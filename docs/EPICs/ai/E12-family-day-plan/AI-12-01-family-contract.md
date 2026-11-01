# AI-12-01 ComposeFamilyDay contract

| | |
|---|---|
| Service | ai |
| Epic | [AI-12 Family day plan](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | AI-10-01 |
| Unblocks | AI-12-02; PLN-13 |

## Goal

Give AI exactly what it needs about each eater to choose one dish for all, and nothing more.

## Business rules

1. Eaters are opaque refs `e1` … `e8`, assigned by planner per request and mapped back by planner. No member id,
   name, age, child flag, weight, height or goal ever reaches ai (features.md §2 "Privacy"; AI receives only each
   person's target).
2. Per eater: personal `KitchenRules` (diets, allergens, dislikes; `tools` ignored) and, per requested meal, a
   `kcal` target (the meal's share of the daily target, computed by planner) or `away=true`.
3. `context.rules` holds the household's tools and any household-wide rules; each eater's rules add to them for the
   meals that eater eats.
4. `eat_lighter=true` when the household turned on "eat lighter as a family" (features.md §2 "Children and
   teenagers").
5. Library candidates in family requests carry `kcal_per_portion` (computed from nutrition tables), so the model can
   prefer dishes that portion well; the number is never echoed in text (calorie guard).
6. Limits: 1–8 eaters; at least one present eater per requested meal.

## Contract

```proto
rpc ComposeFamilyDay(ComposeFamilyDayRequest) returns (ComposeFamilyDayResponse);

message MealTarget {
  MealSlot meal = 1;
  int32 kcal = 2;      // 0 when away
  bool away = 3;
}

message Eater {
  string ref = 1;                        // "e1".."e8"
  KitchenRules rules = 2;
  repeated MealTarget meal_targets = 3;
}

message ComposeFamilyDayRequest {
  CompositionContext context = 1;
  repeated MealRequest meals = 2;
  repeated Eater eaters = 3;
  bool eat_lighter = 4;
  repeated string exclude_dish_ids = 5;
  int32 variation = 6;
  bool bypass_cache = 7;
}

message EaterToppings {
  string eater_ref = 1;
  repeated string ingredient_ids = 2;    // toppings this eater gets
}

message FamilyMeal {
  ComposedMeal meal = 1;
  repeated string eater_refs = 2;        // who eats the main dish
  repeated EaterToppings toppings = 3;   // only for base-plus-toppings dishes
  ComposedMeal second_dish = 4;          // only when a hard limit forces it
  repeated string second_dish_eater_refs = 5;
}

message ComposeFamilyDayResponse {
  repeated FamilyMeal meals = 1;
  CallInfo call = 2;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Request breaks the limits above | `INVALID_ARGUMENT`, `ai_bad_request` |
| Eater ref outside `e1`–`e8`, more than 8 eaters, a meal with no present eater | `INVALID_ARGUMENT`, `ai_bad_request` |

## Acceptance criteria

- [ ] A request built from the features.md worked example (4 eaters, dinner targets 800/720/560/1280) passes input checks; a request with an eater ref other than `e1`–`e8` is refused.
