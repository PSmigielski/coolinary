# GW-07-03 What can I make now

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-01-02, PLN-06 (`SuggestDishes`, extended) |
| Unblocks | MOB-03 |

## Goal

Without AI, rank library dishes by the share of ingredients at home and by what expires first (features.md §2), so
the app is useful with AI off (rule 1, verification 1).

## Business rules

1. `GET /v1/suggestions?meal=&limit=&cuisine=&onlyTools=` — `meal` required, `limit` 1–20 (default 10).
2. Calls planner `SuggestDishes(household_id, member_id, meal, limit, filters)` (exists; extended by PLN-06 with
   missing items and expiring items used).
3. Never calls AI and never counts against a quota.
4. Each suggestion: dish summary (as GW-07-01 list item), `shareOnHand` (0–1), `missing` (ingredients to buy, with
   quantities), `usesExpiring` (stock items expiring within 3 days that the dish uses).

## Contract

`GET /v1/suggestions?meal=dinner&limit=10` — `auth: true`, tag `recipes`.

```json
200 { "suggestions": [ { "dish": { "id": "pierogi", "name": "Pierogi ruskie", "kcalPerPortion": 610, "...": "" },
        "shareOnHand": 0.86, "missing": [ { "ingredientId": "ing_onion", "name": "cebula",
        "quantity": { "value": 1, "unit": "piece" } } ],
        "usesExpiring": [ { "stockItemId": "stk_...", "name": "Twaróg", "expiresOn": "2026-11-02" } ] } ] }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| `meal` missing or invalid | 400 | `validation` |
| Planner down | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] With `ai` stopped and 10 items in the fridge, each meal returns at least 3 suggestions (verification 1).
- [ ] A dish using an item that expires tomorrow lists it in `usesExpiring`.
- [ ] The route makes no AI call and leaves the quota unchanged.
