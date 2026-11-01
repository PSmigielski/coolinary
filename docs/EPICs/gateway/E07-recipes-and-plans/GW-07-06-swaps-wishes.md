# GW-07-06 Swap a meal, replace an ingredient, free-text wish

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R1 |
| Size | M |
| Depends on | GW-03-03, GW-01-03, PLN-09 (`SwapMeal`, `ReplaceIngredient`, `WishMeal`) |
| Unblocks | MOB-03 |

## Goal

Plus users change one meal of a plan: swap it, replace a missing ingredient, or ask for something specific
("something warm in 20 minutes") (features.md §2, AI per request).

## Business rules

1. **Swap** `POST /v1/plans/{planId}/meals/{meal}/swap` body `{}`: open to every household. Calls planner
   `SwapMeal(household_id, plan_id, meal, request_id, language)`. Plus gets an AI dish when AI is available; Free
   (and Plus when AI is down) gets the next best library dish (PLN-09-01 rules 4–5). Swaps never use Free quota.
2. **Replace an ingredient** `POST /v1/plans/{planId}/meals/{meal}/replace-ingredient` body `{ "ingredientId" }`
   (a non-optional ingredient of the meal's dish, e.g. one that is missing): Plus only (`plus: "swap_meal"`, Free gets
   403 `plus_required` before planner is called). Calls planner `ReplaceIngredient`. There is no library fallback:
   when no replacement is found planner answers `FAILED_PRECONDITION` `no replacement found` → 422.
3. **Free-text wish** `POST /v1/plans/wish` body `{ "planId"? | "date"?, "meal", "text" }` (one of `planId` or `date`;
   without a plan for that date planner creates a proposed plan with only that meal): Plus only (`plus: "swap_meal"`).
   `text` trimmed, 1–200 characters. It is user text: passed to planner as data, never logged, never put in metrics.
   Calls planner `WishMeal(household_id, plan_id, date, meal, text, request_id, language)`. A wish that can't be met
   is `FAILED_PRECONDITION` `wish not possible` with the main violation code (e.g. `TIME_OVER`) → 422.
4. All three use the `ai` deadline class, idempotency required (the key is passed as `request_id`). Planner keeps no
   plan version, so no `If-Match`.
5. Response: the changed meal (`PlanMeal` shape of GW-07-04) and `aiUsed` (plus `fallbackReason` on swap,
   `replayed` on replace). The app refetches the plan when it needs the whole day.
6. Plus fair use (GW-03-04) counts every answer with `aiUsed: true`.

## Contract

`POST /v1/plans/{planId}/meals/{meal}/swap`, `POST /v1/plans/{planId}/meals/{meal}/replace-ingredient`,
`POST /v1/plans/wish` — `auth: true`, tag `plans`.

```json
request  { "date": "2026-11-02", "meal": "supper", "text": "coś ciepłego w 20 minut" }
200      { "planId": "pln_...", "meal": { "...": "" }, "aiUsed": true }
```

RPCs: planner `SwapMeal`, `ReplaceIngredient`, `WishMeal` (PLN-09-01 … 03).

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free on replace-ingredient or wish | 403 | `plus_required` |
| Wish empty or > 200 chars, neither `planId` nor `date` | 400 | `validation` |
| Meal already cooked / ingredient not in dish / nothing else fits | 422 | planner's code |
| No replacement found / wish not possible | 422 | `no_replacement_found` / `wish_not_possible` (+ violation in `meta`) |
| Fair use exceeded | 429 | `fair_use_limit` |

## Acceptance criteria

- [ ] A Free token on replace-ingredient or wish gets 403 and planner is not called.
- [ ] A Free swap returns a library dish with `aiUsed: false` and no quota change.
- [ ] A wish returns a dish (or 422 `wish_not_possible`); the wish text is absent from all logs.
- [ ] With `ai` stopped, a Plus swap returns a library dish with `aiUsed: false`.
