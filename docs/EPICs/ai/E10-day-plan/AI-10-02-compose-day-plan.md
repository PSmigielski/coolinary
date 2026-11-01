# AI-10-02 ComposeDayPlan prompt and call

| | |
|---|---|
| Service | ai |
| Epic | [AI-10 Day plan composition and cuisine choices](README.md) |
| Release | R1 |
| Size | L |
| Status | todo |
| Depends on | AI-10-01, AI-01-03, AI-01-05, AI-01-09 |
| Unblocks | AI-10-03, AI-10-04; PLN-08 |

## Goal

Turn the household's stock and kitchen rules into a day of meals that a person would actually cook, using what
expires first and preferring checked library dishes.

## Business rules

1. Job `compose_day_plan` at effort `medium`: Haiku 5.5 for Free (decision 15); for Plus, Haiku 5.5 unless AI-10-05
   moves it to Sonnet 5.5. Prompt `compose_day_plan/v1`.
2. Composition rules given to the model, in priority order (features.md §2 "One dish per meal", conflict order):
   1. Hard limits: only the kitchen's tools; every diet; no excluded allergen; no disliked ingredient. A meal that
      can't be made within them is left out with reason in `cuisine_fit` (AI-10-04), never bent.
   2. Each requested meal once. The `main_hot_meal` meal is the cooked one; meals with `prefer_no_cooking` come
      from no-cooking dishes (library candidates tagged `NO_COOKING` first).
   3. Stock: use only stock, staples and, when allowed, `buyable_ids` (at most `max_items_to_buy` for the whole day,
      0 by default). Across the day, the sum of `stock_uses` per stock line never exceeds `available`.
   4. Oldest first: `past_date` lines first, then lines with the smallest `days_to_expiry`; opened items before
      closed ones of the same food; unknown expiry (−1) last. Past-date lines are used like other stock and come back
      marked `check_first` (Q103, AI-10-01 rule 7).
   5. Library first: when a library candidate fits as well as a new dish (same hard limits, uses the expiring items
      about as well), choose the library dish — it is already checked and its calories are known.
   6. Time: `total_minutes` ≤ the meal's `max_total_minutes`.
   7. Cuisine: within the meal's cuisines; otherwise the cuisine-choice path (AI-10-04).
   8. Variety: no dish or main protein twice in one day; nothing in `exclude_dish_ids` / `exclude_dish_names`.
3. New dishes: written in the request language; servings = `portions`; steps 1–20, short and imperative, timers
   where food waits; `tools` lists every tool used; effort tags as suggestions; no calorie or health claims.
4. `stock_uses` list every stock line the meal uses, in the stock line's unit, for library and new dishes alike.
5. Request layout (AI-01-05): the system prompt and the tool/cuisine/allergen vocabulary are the cached prefix;
   facts, stock, candidates and rules follow in the user message as sorted JSON.
6. Variation: `variation` > 0 asks for a different plan than earlier ones (the planner also sends the previous dish
   ids in `exclude_dish_ids`). Cache TTL 1 hour (AI-01-09).
7. Tracing: counts of `meals`, stock lines, `candidates`, new vs library dishes; no names.
8. Expected size: ≈ 7,000 input and ≈ 3,500 output tokens (AI-02 estimate); `max_tokens` 16000.

## Contract

```proto
rpc ComposeDayPlan(ComposeDayPlanRequest) returns (ComposeDayPlanResponse);

message ComposeDayPlanRequest {
  CompositionContext context = 1;           // AI-10-01
  repeated MealRequest meals = 2;           // 1-3
  repeated string exclude_dish_ids = 3;
  repeated string exclude_dish_names = 4;
  int32 variation = 5;
  bool bypass_cache = 6;
}

message ComposeDayPlanResponse {
  repeated ComposedMeal meals = 1;
  CallInfo call = 2;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Limits of AI-10-01 rule 5 broken; an id in stock/candidates without a fact; unknown tool or allergen code | `INVALID_ARGUMENT`, `ai_bad_request` |
| Others | AI-01-04; planner falls back to library matches (PLN-08) |

## Acceptance criteria

- [ ] Fixed test fridge (12 lines, eggs and spinach expiring in 1 day), tools hob+oven, Polish breakfast / Asian dinner: eggs and spinach are used, breakfast is no-cooking or quick, dinner is Asian, no tool outside hob/oven.
- [ ] With a library candidate that fits (share on hand 1.0), it is chosen over a new dish in ≥ 8 of 10 runs (eval tag `library-first`).
- [ ] `variation=1` with the previous dish ids excluded returns different dishes.
- [ ] Golden set thresholds pass (AI-10-06).
