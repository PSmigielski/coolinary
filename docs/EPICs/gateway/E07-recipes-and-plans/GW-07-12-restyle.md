# GW-07-12 Restyle a dish

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R3 |
| Size | S |
| Status | todo |
| Depends on | GW-03-03, PLN-09 (`RewriteDish`, kind `RESTYLE`) |
| Unblocks | MOB-03 |

## Goal

"Make this Italian-style": a dish rewritten in another cuisine (features.md §2 "Cuisine choice", Plus).

## Business rules

1. Plus-only (`plus: "restyle"`). Body `{ "cuisine" }`; `cuisine` must be a code from `GET /v1/profile/options`.
2. Planner `RewriteDish(household_id, dish_id, kind = RESTYLE, target_cuisine, request_id, language)` (PLN-09-05);
   same rules and response as GW-07-09, plus `toBuyIngredientIds` (ingredients the restyled dish needs that are not
   at home).
3. AI unavailable → 503 `ai_unavailable`.

## Contract

`POST /v1/dishes/{dishId}/restyle` — `auth: true`, `plus: "restyle"`, tag `recipes`, idempotency required.

```json
request { "cuisine": "italian" }
200     { "dish": { "...": "", "parentDishId": "pierogi" }, "toBuyIngredientIds": ["ing_parmesan"], "aiUsed": true, "cached": false }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| Unknown cuisine | 400 | `validation` |
| AI unavailable | 503 | `ai_unavailable` |

## Acceptance criteria

- [ ] Restyling pierogi as `italian` returns a household dish with cuisine `italian` that passes the rule checker.
- [ ] A Free token gets 403.
