# GW-07-10 Plan from deals

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | GW-03-03, GW-10-02, PLN-18 (`CreatePlanFromDeals`), OFF-04 (`ListMyShops`, used by planner) |
| Unblocks | MOB-03 |

## Goal

Compose dishes from the fridge plus this week's offers, preferring the nearest shops (features.md §2 "Flyers and
offers", Plus). Coolinary stays fridge-first: offers fill the gaps.

## Business rules

1. Plus-only (`plus: "plan_from_deals"`). Body `{ "date", "replace"?: false, "useAi"?: true }`.
2. Shops: planner reads the household's shops from offers (`ListMyShops`), in the order the phone saved them, nearest
   first (OFF-04-01 rule 1, D13). The route takes no shop ids: the app keeps that order current through
   `PUT /v1/my-shops` (GW-10-02) when the user's position changes. The server never computes distance.
3. Planner `CreatePlanFromDeals(household_id, viewer_member_id, date, replace, use_ai)` with the `ai` deadline;
   idempotency required. A plan for that date needs `replace: true`, otherwise planner's `ALREADY_EXISTS` → 409
   `plan_exists` (same as the family route, GW-07-04 rule 9).
4. Response: planner's `DealsPlan`: the family day plan (GW-07-04 rule 9 shape, with `aiUsed`, `fallbackReason`),
   the deal `items` (meal, offer id, ingredient, shop id, price, pack, valid to) and `offerCostEstimate`.
5. AI unavailable → planner falls back to library dishes ranked with offers (PLN-18-03), `aiUsed: false`.

## Contract

`POST /v1/plans/from-deals` — `auth: true`, `plus: "plan_from_deals"`, tag `plans`.

```json
request { "date": "2026-11-03" }
200     { "plan": { "...": "", "aiUsed": true, "fallbackReason": "" },
          "items": [ { "meal": "dinner", "offerId": "...", "ingredientId": "ing_chicken_thigh", "shopId": "osm-n123",
                       "price": { "amountMinor": 1599, "currency": "PLN" }, "validTo": "2026-11-08" } ],
          "offerCostEstimate": { "amountMinor": 2398, "currency": "PLN" } }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| No shops saved for the household | 422 | planner's code |
| Plan exists for the date (no `replace`) | 409 | `plan_exists` |

## Acceptance criteria

- [ ] Every offer in `items` belongs to one of the household's saved shops and is valid on the plan date.
- [ ] With `ai` stopped the route returns library dishes with offers and `aiUsed: false`.
- [ ] The request body accepts no coordinates or distances (schema test).
