# PLN-06-03 Extend SuggestDishes

| | |
|---|---|
| Service | planner |
| Epic | [PLN-06 What can I make now](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-06-02, PLN-05-02, PLN-05-03 |
| Unblocks | GW-07 (`GET /v1/suggestions`), PLN-08-05 |

## Goal

The existing `SuggestDishes` RPC returns ranked library dishes with what's missing and what's expiring, applying the
household's hard limits and the same filters as browse — fully without AI.

## Business rules

1. Request keeps `household_id`, `meal`, `limit`; adds `language`, `servings` (0 = household portions), `filters`
   (`DishFilters`, PLN-05-03), `eater_member_ids` (PLN-05-02), `include_household_ai` (default false).
2. `meal = MEAL_UNSPECIFIED` returns the best dishes over all meals (existing behaviour kept).
3. `limit` 1–50, default 10 (existing rule kept).
4. Hard limits (PLN-05-02) are applied before ranking; filters narrow further.
5. Each suggestion: dish id, name, share on hand, kcal per portion (0 when incomplete) plus: group, score, total minutes,
   missing lines (ingredient id, name, missing quantity, partly at home), expiring stock items used (stock item id,
   ingredient name, days left), `leave_out_ingredient_ids`.
6. Result cached in redis-cache by (household, profile_version, stock_version, request hash) for 2 minutes.
7. No ai call; works with ai down (Verification 1).
8. Existing field numbers stay; only new fields are added (`make proto-breaking` passes).

## Contract

```proto
message SuggestDishesRequest {
  string household_id = 1;
  Meal meal = 2;
  int32 limit = 3;
  string language = 4;
  double servings = 5;
  DishFilters filters = 6;
  repeated string eater_member_ids = 7;
  bool include_household_ai = 8;
}
message DishSuggestion {
  string dish_id = 1;
  string name = 2;
  double share_on_hand = 3;
  int32 kcal_per_portion = 4;
  SuggestionGroup group = 5;
  double score = 6;
  int32 total_minutes = 7;
  repeated MissingLine missing = 8;
  repeated ExpiringUse expiring = 9;
  repeated string leave_out_ingredient_ids = 10;
  bool check_first = 11;                      // uses an item past its date: app shows "check it first" (Q103)
}
message MissingLine { string ingredient_id = 1; string name = 2; coolinary.common.v1.Quantity missing = 3; bool partly_at_home = 4; }
message ExpiringUse { string stock_item_id = 1; string name = 2; int32 days_left = 3; bool check_first = 4; }  // days_left < 0 = past its date (Q103)
enum SuggestionGroup { SUGGESTION_GROUP_UNSPECIFIED = 0; SUGGESTION_GROUP_READY = 1;
                       SUGGESTION_GROUP_ALMOST = 2; SUGGESTION_GROUP_NEEDS_SHOPPING = 3; }
```

## Data

No schema change (the skeleton query `ListDishesForMeal` is removed).

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Unknown meal (existing) | INVALID_ARGUMENT | `unknown meal` |
| Profile or stock unavailable | UNAVAILABLE | from PLN-01-03/04 |

## Acceptance criteria

- [ ] `make proto-breaking` reports no breaking change.
- [ ] With 10 test items and AI off, breakfast, dinner and supper each return ≥ 3 suggestions filtered by cuisine and
      tools (Verification 1).
- [ ] A suggestion's missing lines never include a staple or an item fully at home.
- [ ] Two identical calls within 2 minutes make one stock and one profile read.
- [ ] A dish using yoghurt one day past its date is suggested with `check_first = true`; its `expiring` entry has
      `days_left = -1` and `check_first = true` (Q103).
