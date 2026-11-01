# PLN-04-04 Pack values of the household's scanned products

| | |
|---|---|
| Service | planner |
| Epic | [PLN-04 Dish nutrition and calories](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-04-01, PLN-01-04, KIT-06 (ComputeNutrition with product overrides) |
| Unblocks | PLN-10-01 |

## Goal

When the household cooks with a product it scanned and the pack's values are known, those values replace the table's
(features.md §5 step 4), so "about 520 kcal" reflects the oat drink actually in the fridge.

## Business rules

1. Applies to dish detail and cooking mode for one household (`GetDish`, `GetCookingSteps`), not to browse lists,
   filters or plans (those use the library values so results are comparable). Decided (Q098).
2. For each non-optional line, look in the stock snapshot for items with the same `ingredient_id` and a product EAN.
   If several, use the one that expires first (it will be used first; tie → added first).
3. Planner calls `ComputeNutrition` with that line's `ean` as an override; kitchen uses the pack values when it has
   all of energy, fat, saturates, carbohydrate, sugars, protein and salt per 100 g (the EU mandatory set), else the
   table entry. Fibre falls back to the table when the pack lacks it.
4. Result cached in redis-cache `planner:hhnut:<household>:<dish>:<version>:<stock_version>` for 1 hour.
5. The response sets `from_household_products=true` when at least one line used pack values; the app says "using your
   products".
6. Kitchen unavailable → fall back silently to the stored library values.

## Contract

Uses `KitchenService.ComputeNutrition` lines with an optional `ean`. Decided (Q097): KIT-06 accepts a product EAN
per line (features.md §5 step 4 is kitchen's rule).

## Data

No schema change.

## Errors

None of its own: kitchen unavailable → silent fallback to the stored library values (rule 6). Errors of the calling RPC (PLN-10-01) apply.

## Acceptance criteria

- [ ] A household with a scanned oat drink whose pack says 46 kcal/100 ml sees porridge kcal computed from 46, not
      from the table value, and `from_household_products=true`.
- [ ] Another household without that product sees the library value.
- [ ] Browse results show the same kcal for both households.
