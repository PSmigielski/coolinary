# AI-12-03 Validation for eaters, toppings and second dishes

| | |
|---|---|
| Service | ai |
| Epic | [AI-12 Family day plan](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | AI-12-02, AI-10-03 |
| Unblocks | AI-12-04 |

## Goal

No eater is ever given food that breaks their own hard limits, and second dishes appear only when they must.

## Business rules

1. AI-10-03 runs on every main and second dish with the household rules.
2. `EATERS` — `eater_refs` of a meal = exactly the present eaters (not away) for that meal, minus
   `second_dish_eater_refs`; every ref exists.
3. `EATER_LIMITS` — for each eater of the main dish: no ingredient of the base and of that eater's toppings breaks
   the eater's allergens, diets or dislikes. Toppings given to other eaters are not checked against this eater.
4. `TOPPINGS` — `toppings` only on dishes tagged or built as base plus toppings (ingredients with role `TOPPING`);
   every topping id is a `TOPPING` ingredient of the dish; every eater of the dish has an entry.
5. `SECOND_DISH` — allowed only if the main dish breaks a hard limit for each eater in `second_dish_eater_refs`
   (checked: removing that eater's personal rules would make the main dish valid for them); the second dish meets
   those eaters' limits.
6. Stock: the sum of `stock_uses` of main and second dishes over the day ≤ available.

## Contract

No change.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Still invalid after the one repair | `UNAVAILABLE`, `ai_invalid_output` |

## Acceptance criteria

- [ ] A fake answer giving a nut topping to a nut-allergic eater fails `EATER_LIMITS`.
- [ ] A second dish for an eater the main dish already suits fails `SECOND_DISH`.
- [ ] An away eater listed in `eater_refs` fails `EATERS`.
