# GW-07-07 Cooking mode and "cooked it"

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-01-03, PLN-10 (`GetCookingSteps`, `MarkCooked`, `DismissRatingPrompt`) |
| Unblocks | MOB-03, MOB-04 (cooked-it works offline) |

## Goal

Step-by-step cooking with timers scaled to the pot, and a "cooked it" button that reduces the stock and records the
history (features.md §2 "Dishes and the day plan").

## Business rules

1. `GET /v1/cooking/{dishId}?version=&planMealId=&servings=` → planner `GetCookingSteps(household_id, dish_id,
   version, servings, plan_meal_id, language)`: the dish with steps and timers, ingredients scaled to the pot (from
   the plan meal when `planMealId` is given, else `servings`, default household portions), and in R2 how to divide
   the pot per member (GW-07-08).
2. `POST /v1/cooked` body `{ "dishId", "dishVersion", "planMealId"?, "servings", "cookedAt", "chosenToppings"?,
   "leftOutIngredientIds"?, "optionalUsedIngredientIds"? }` → planner `MarkCooked(household_id, user_id, cooked_id,
   ...)`, which calls kitchen `ConsumeStock`. `servings` > 0 and ≤ 48; `cookedAt` RFC 3339, not in the future, not
   older than 7 days (offline queue window, MOB-04).
3. Idempotency **required**; the gateway passes the idempotency key as planner's `cooked_id`, so a replay from the
   offline queue never reduces stock twice (planner answers `replayed: true`).
4. Response: `entry` (the history entry), `shortfalls` (ingredients the stock did not have enough of; not an error),
   `replayed`, `ratingPrompt` (true when the app should ask for a rating).
5. `POST /v1/dishes/{dishId}/rating-prompt/dismiss` → planner `DismissRatingPrompt(household_id, user_id, dish_id)`:
   the user said "not now"; response 204.
6. "Cooked it" is free, no quota, works for AI and library dishes.
7. R3 leftovers of a double cook are shown in the week plan (GW-07-11); `MarkCooked` returns no leftover data.

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/cooking/{dishId}?version=&planMealId=&servings=` | 200 `{ "dish": { "...", "steps": [ { "n", "text", "timerSeconds" } ] }, "servings", "ingredients": [...], "divide": null }` |
| POST | `/v1/cooked` | 201 `{ "entry": { "id", "dishId", "cookedAt", "servings" }, "shortfalls": [ { "ingredientId", "missing" } ], "replayed": false, "ratingPrompt": true }` |
| POST | `/v1/dishes/{dishId}/rating-prompt/dismiss` | 204 |

RPCs: planner `GetCookingSteps`, `MarkCooked`, `DismissRatingPrompt` (PLN-10). Tags `cooking`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| `cookedAt` in the future or > 7 days old | 400 | `validation` |
| Unknown dish or plan | 404 | `not_found` |
| Kitchen unavailable (planner can't reduce stock) | 503 | `service_unavailable` (nothing recorded) |

## Acceptance criteria

- [ ] Cooking a 4-portion dish with 500 g rice in stock and 300 g used leaves 200 g.
- [ ] The same request replayed with the same key leaves 200 g (no second reduction) and says `replayed: true`.
- [ ] Missing stock is reported in `shortfalls`, and the history entry is still created.
- [ ] `GET /v1/cooking/...?servings=8` for a 4-serving dish doubles quantities and keeps timers unchanged.
- [ ] Dismissing the rating prompt returns 204 and the next `cooked` for that dish says `ratingPrompt: false`.
