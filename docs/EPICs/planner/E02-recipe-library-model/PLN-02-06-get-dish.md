# PLN-02-06 GetDish and BatchGetDishes

| | |
|---|---|
| Service | planner |
| Epic | [PLN-02 Recipe library model](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-02-01, PLN-02-04, PLN-02-05, KIT-02 (GetIngredients) |
| Unblocks | GW-07 (dish detail), PLN-05, PLN-08, PLN-10 |

## Goal

Read one dish (or many) in the user's language, at a chosen number of portions, with ingredient names from the
catalogue — the base read every screen and plan uses.

## Business rules

1. `GetDish(dish_id, version?, servings?, language)`:
   - no version → `current_version`; a version → that version (used by plans and history).
   - servings default = the household's portions from the profile (PLN-01-03); falls back to `base_servings`.
2. Visibility: published and retired library dishes are readable by any household; `household_ai` dishes only by their
   household (`household_id` must match) — otherwise NOT_FOUND (never PERMISSION_DENIED, to not leak ids). Draft and
   review versions only through `PlannerAdminService`.
3. Language: requested language → `en` → `pl` → any. The response says which language it used and whether it is
   machine translated.
4. Ingredient display names come from `KIT-02 GetIngredients` in the same language, cached in redis-cache
   `planner:ingname:<lang>:<id>` for 1 day.
5. `BatchGetDishes(ids, language)` returns up to 100 dishes at their current versions, without steps (list views).
6. Responses include derived flags (PLN-02-03), nutrition (PLN-04-03), scaled quantities (PLN-02-05) with both
   `exact` and `display` values, and, for the calling household, `is_favourite` and `my_rating` (PLN-05-04/05).

## Contract

```proto
rpc GetDish(GetDishRequest) returns (GetDishResponse);
rpc BatchGetDishes(BatchGetDishesRequest) returns (BatchGetDishesResponse);

message GetDishRequest {
  string household_id = 1;
  string dish_id = 2;
  int32 version = 3;          // 0 = current
  double servings = 4;        // 0 = household default
  string language = 5;
}
message GetDishResponse {
  Dish dish = 1;
  double servings = 2;
  repeated ScaledIngredient scaled = 3;
  bool is_favourite = 4;
  int32 my_rating = 5;        // 0 = not rated
}
message ScaledIngredient {
  int32 position = 1;
  coolinary.common.v1.Quantity exact = 2;
  coolinary.common.v1.Quantity display = 3;
  bool check_first = 4;            // set in cooking mode when stock past its date covers it (PLN-10-01, Q103)
}
message BatchGetDishesRequest { string household_id = 1; repeated string dish_ids = 2; string language = 3; }
message BatchGetDishesResponse { repeated Dish dishes = 1; }
```

## Data

No schema change.

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Unknown id, or another household's AI dish | NOT_FOUND | `dish not found` |
| More than 100 ids | INVALID_ARGUMENT | `at most 100 dishes` |
| Bad servings | INVALID_ARGUMENT | from PLN-02-05 |

## Acceptance criteria

- [ ] Household B asking for household A's AI dish gets NOT_FOUND.
- [ ] A dish with only `pl` and `en` texts asked in `de` returns `en` with `language="en"`.
- [ ] A plan referencing version 1 still gets version 1 after version 2 is published.
- [ ] p95 latency of `GetDish` with warm caches < 50 ms in the dev stack.
