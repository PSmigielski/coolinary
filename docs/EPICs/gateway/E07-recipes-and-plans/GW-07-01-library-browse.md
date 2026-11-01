# GW-07-01 Library browse and dish detail

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-01-04, PLN-05 (`SearchDishes`, `GetDish`), PLN-04 (calories per dish) |
| Unblocks | MOB-03 |

## Goal

Search and filter the recipe library, and open a dish with ingredients, steps, tools, effort tags and calories per
portion (features.md §2 "Dishes", "Calories and goals", "Kitchen tools").

## Business rules

1. `GET /v1/dishes` query (all optional): `q` (2–60 chars), `meal` (`breakfast` | `dinner` | `supper`), `cuisine`
   (repeatable), `maxMinutes` (5–240), `kcal` (`lt400` | `400to600` | `gt600`), `protein25` (`true` = 25 g protein
   or more per portion), `effort` (repeatable: `no_cooking`, `keeps_3_days`, `freezes_well`, `good_cold`,
   `base_plus_toppings`), `onlyTools` (repeatable; the one-off "air fryer only" filter for tonight), `favourites`
   (`true`), `pageSize`, `pageToken`.
2. Planner always applies the household's hard limits (tools, diets, allergens from the kitchen profile) whatever
   the query says; `onlyTools` narrows further and can't widen. The gateway passes `household_id` and `member_id`.
3. Calorie bands use the per-standard-portion value; dishes with `caloriesIncomplete` are excluded from `kcal` and
   `protein25` filters (they have no number).
4. List item: id, name, meals, cuisine, total minutes, kcal per portion (rounded to 10, `about`), effort tags,
   `caloriesIncomplete`, `isFavourite`, `rating` (household average), `imageUrl` (nullable).
5. `GET /v1/dishes/{dishId}?portions=<n>` (n 1–24, default the household's portions): ingredients scaled to `n`,
   steps with timers, tools, diet flags, allergens, effort tags, nutrition per portion (kcal, protein,
   carbohydrate, sugars, fat, saturates, fibre, salt), `caloriesIncomplete` with the names of unmatched
   ingredients (features.md §5 rule 6), `onHand` per ingredient (true/false from stock).
6. kcal and kJ: the server returns kcal and kJ both (`kcal`, `kj`); the app shows the user's choice.
7. Dish names and steps come in `x-language`; planner falls back to English.
8. No label like "healthy" or "light" appears anywhere in the response (features.md §2).

## Contract

`GET /v1/dishes`, `GET /v1/dishes/{dishId}` — `auth: true`, tag `recipes`.

```json
list 200 { "dishes": [ { "id": "chicken-curry", "name": "Kurczak curry z ryżem", "meals": ["dinner"],
             "cuisine": "asian:indian", "totalMinutes": 40, "kcalPerPortion": 520, "kjPerPortion": 2180,
             "caloriesIncomplete": false, "effort": ["keeps_3_days"], "isFavourite": false, "rating": 4.5,
             "imageUrl": null } ], "nextPageToken": "..." }
detail 200 { "id": "...", "name": "...", "portions": 4,
             "ingredients": [ { "ingredientId": "ing_rice_basmati", "name": "ryż basmati",
                                "quantity": { "value": 300, "unit": "g" }, "optional": false, "onHand": true } ],
             "steps": [ { "n": 1, "text": "...", "timerSeconds": 600 } ],
             "tools": ["hob"], "dietFlags": ["gluten_free"], "allergens": [], "effort": ["keeps_3_days"],
             "nutritionPerPortion": { "kcal": 520, "kj": 2180, "proteinG": 31, "carbohydrateG": 62, "sugarsG": 6,
                                      "fatG": 14, "saturatesG": 4, "fibreG": 3, "saltG": 1.4, "about": true },
             "caloriesIncomplete": false, "unmatchedIngredients": [] }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bad filter value | 400 | `validation` |
| Dish unknown or not visible to this household | 404 | `not_found` |

## Acceptance criteria

- [ ] With the "microwave + kettle" preset, no dish in any page needs another tool (verification 2, planner rule).
- [ ] `kcal=lt400` returns only dishes under 400 kcal per portion and no `caloriesIncomplete` dish.
- [ ] `portions=8` doubles every ingredient of a 4-portion dish; kcal per portion is unchanged.
- [ ] A dish with an unmatched ingredient shows `caloriesIncomplete: true` and names it (verification 16).
