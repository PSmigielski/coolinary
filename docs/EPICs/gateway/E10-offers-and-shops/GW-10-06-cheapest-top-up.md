# GW-10-06 Cheapest top-up

| | |
|---|---|
| Service | gateway |
| Epic | [GW-10 Shops, offers, flyers and price matrix routes](README.md) |
| Release | R3 |
| Size | S |
| Status | todo |
| Depends on | GW-10-05, GW-07-11, PLN-20 (`CreateCheapestTopUp`, `ApplyTopUp`) |
| Unblocks | MOB-03 |

## Goal

"Cheapest top-up": the smallest basket that completes the week plan using current offers, including swapping an
ingredient for a discounted equivalent (features.md §2 "Shopping list", AI per request, Plus).

## Business rules

1. **Plus only** (`plus: "cheapest_top_up"`). `POST /v1/list/cheapest-top-up` body `{ "weekPlanId", "shopIds":
   [1–10, nearest first], "useAi"?: true }` → planner `CreateCheapestTopUp(household_id, week_plan_id, shop_ids,
   use_ai)` with the `ai_long` deadline (120 s), idempotency required. It needs a week plan (GW-07-11).
2. Response: planner's `TopUp`: `id`, `expiresAt`, the lines to buy (ingredient, amount, shop, offer, cost,
   source), `swaps` (each with a `key`: an ingredient swapped for a discounted equivalent, the dishes it changes, the
   saving and the kcal change), estimated `total` and `saving` (both "about"), `aiUsed`, `fallbackReason`.
3. AI unavailable → planner answers with its non-AI top-up (PLN-20-03), `aiUsed: false`.
4. Swaps change dishes only after the user accepts them: `POST /v1/list/cheapest-top-up/{topUpId}/apply` body
   `{ "swapKeys": [...] }` → `ApplyTopUp(household_id, member_id, topup_id, swap_keys, client_mutation_id)`, which
   adds the lines to the list and applies the chosen swaps to the week plan. Response: applied and skipped swap
   keys, new list version and week plan version. Idempotency required.
5. An expired top-up → planner's `FAILED_PRECONDITION` → 410; the app asks for a new one.

## Contract

| Method | Path | 2xx |
|---|---|---|
| POST | `/v1/list/cheapest-top-up` | 200 top-up |
| POST | `/v1/list/cheapest-top-up/{topUpId}/apply` | 200 `{ "appliedSwapKeys", "skippedSwapKeys", "listVersion", "weekPlanVersion" }` |

```json
200 { "id": "tu_...", "expiresAt": "...",
      "lines": [ { "ingredientId": "ing_rice", "amount": { "value": 1000, "unit": "g" }, "shopId": "osm-lidl-1",
                   "offerId": "off_..", "cost": { "amountMinor": 179, "currency": "EUR" }, "source": "flyer" } ],
      "swaps": [ { "key": "s1", "fromIngredientId": "ing_basmati", "toIngredientId": "ing_jasmine",
                   "dishIds": ["..."], "saving": { "amountMinor": 120, "currency": "EUR" } } ],
      "total": { "amountMinor": 6420, "currency": "EUR" }, "saving": { "amountMinor": 310, "currency": "EUR" },
      "aiUsed": true, "fallbackReason": "" }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| Unknown week plan, or 0 / > 10 shops | 404 / 400 | `not_found` / `validation` |
| Top-up expired on apply | 410 | planner's code |

## Acceptance criteria

- [ ] For a test week and offers, every swap points to a valid offer at one of the given shops.
- [ ] Nothing changes on the list or in the week plan until `apply`; applying with one of two swap keys applies only that swap.
- [ ] With `ai` stopped, a non-AI proposal returns with `aiUsed: false`.
