# PLN-20-01 Top-up needs and allowed swaps

| | |
|---|---|
| Service | planner |
| Epic | [PLN-20 Cheapest top-up](README.md) |
| Release | R3 |
| Size | M |
| Depends on | PLN-15-02, PLN-15-03, PLN-18-01, PLN-07 |
| Unblocks | PLN-20-02, PLN-20-03 |

## Goal

Know exactly what the week plan still needs after the fridge, and which ingredient swaps would be acceptable, before
any AI is asked.

## Business rules

1. Needs: PLN-15-02 and PLN-15-03 for the active week plan's remaining (uncooked) meals, merged per ingredient
   (PLN-15-04), each with the dishes and meals that use it. Ingredients already on the list as open lines count as
   needed; ticked lines count as bought.
2. Offers: current offers at the household's shops (PLN-18-01 candidates, but asked for every needed ingredient and
   for equivalents, cap 500 ids).
3. Equivalent candidates for a needed ingredient: catalogue ingredients in the same catalogue category (KIT-02
   `GetIngredients`) that are on offer this week.
4. A swap is allowed only if:
   - the dish with the swap still passes PLN-07 for every eater of every meal using it (allergens, diet, dislikes,
     tools);
   - the dish's kcal per standard portion changes by at most 15% (PLN-04 calculation on the modified ingredient
     list), so portions stay sensible;
   - the swap is cheaper than the cheapest price of the original at the household's shops by at least 10% or €0.50.
5. At most 3 swaps per week plan; at most 1 per dish.
6. Output: needs (ingredient, amount, using meals), allowed swaps (from, to, dish ids, kcal change %, price
   difference), candidate offers.

## Contract

Internal (`topup.Prepare(household, week) Prepared`).

## Data

No schema change.

## Errors

None of its own; internal preparation called by `CreateCheapestTopUp` (PLN-20-02), whose errors apply.

## Acceptance criteria

- [ ] Chicken breast needed, chicken thighs on offer 35% cheaper, same category: allowed for a curry when no eater
      has a limit against thighs.
- [ ] A swap from feta to a cow's-milk cheese for a lactose-free eater's meal is not allowed.
- [ ] A swap changing the dish's kcal per portion by 20% is not allowed.

## Notes

- Decided (Q224): calories change ≤ 15%, a swap saves ≥ 10% or €0.50, at most 3 swaps a week.
