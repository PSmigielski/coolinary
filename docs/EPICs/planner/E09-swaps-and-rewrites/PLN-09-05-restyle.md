# PLN-09-05 Restyle a dish ("make this Italian-style")

| | |
|---|---|
| Service | planner |
| Epic | [PLN-09 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | PLN-09-04 |
| Unblocks | GW-07 |

## Goal

Turn a dish into another cuisine's style with what's at home: swap spices, herbs and sauces, keep the base (features.md
"Cuisine choice").

## Business rules

1. `RewriteDish(kind=RESTYLE, target_cuisine)`; Plus only; target cuisine ≠ the dish's cuisine.
2. AI may replace or add seasoning, herb, sauce and fat lines (by kitchen ingredient category) and change steps; main
   and base ingredients stay. New lines must be at home or staples; up to 2 may be purchases, returned as `to_buy`.
   Decided (Q223): the allowed lines are those in KIT-02's seasoning, herb, sauce and fat category codes; the code
   list is fixed in this task at implementation.
3. Checked with the full checker for the household's eaters (`FROM_STOCK` with the `to_buy` items allowed). The result
   gets `cuisine = target_cuisine`.
4. Stored as household AI dish with `parent_dish_id`; when used in a plan meal, `to_buy` becomes the meal's allowed
   purchases.
5. Failing answer → FAILED_PRECONDITION `restyle not possible`.

## Contract

`RewriteDish` (PLN-09-04) with `kind=REWRITE_KIND_RESTYLE` and `target_cuisine`; response gains
`repeated string to_buy_ingredient_ids = 2;`.

## Data

No schema change.

## Errors

Reason codes travel as trailing gRPC metadata `x-error-code` (README decision 18).

| Case | gRPC status / behaviour |
|---|---|
| Free household | `FAILED_PRECONDITION`, `x-error-code: plus_required` |
| `target_cuisine` unspecified or equal to the dish's cuisine | `INVALID_ARGUMENT` |
| Unknown dish | `NOT_FOUND` |
| AI answer fails the checks | `FAILED_PRECONDITION`, `x-error-code: restyle_not_possible` (rule 5) |
| ai down, over budget or refusing | `UNAVAILABLE`, `x-error-code: ai_unavailable`; restyle has no non-AI fallback, the dish stays as it was |

## Acceptance criteria

- [ ] Restyling a Polish pork stew to Italian keeps pork and potatoes, changes herbs/sauce, and lists at most 2 items to buy.
- [ ] Restyling into the dish's own cuisine is rejected with INVALID_ARGUMENT.
