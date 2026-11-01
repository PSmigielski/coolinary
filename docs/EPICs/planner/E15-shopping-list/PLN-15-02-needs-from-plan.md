# PLN-15-02 Ingredient needs from the plan

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-15-01, PLN-08, PLN-02 |
| Unblocks | PLN-15-03, PLN-15-10, PLN-18-04, PLN-20-01 |

## Goal

Step 1 of making the list: take every ingredient of the dishes in the chosen plan dates, scaled to the household's
portions, or to the pot when goals are set (features.md "Shopping list" step 1).

## Business rules

1. Input: household, a date range (`from` … `to`, at most 14 days, not before today) and optionally a set of meals.
2. Included: planned meals in the range with a dish, not yet cooked, with at least one eater.
3. Scale factor per meal:
   - R1: the portions chosen on the planned meal (PLN-05 portion scaling) or else the household portions setting,
     ÷ the dish's base servings;
   - from R2, when a pot is stored (PLN-12-06): the pot ÷ base servings, via PLN-12-08;
   - leftover meals (R3, PLN-14-03) add nothing; their session's cooking meal carries the doubled pot.
4. Each need keeps: ingredient id, amount and unit as in the recipe × factor, meal date and meal, plan meal id,
   dish id, `optional`.
5. Base-plus-toppings: toppings scaled per PLN-12-08 rule 4.
6. Extras not at home (R2, PLN-12-04) become needs of source kind `extra` with the member key.
7. Amounts are converted to the ingredient's base unit (g, ml, piece) with `GetIngredients` typical weights and
   densities. If no conversion exists, the need keeps its unit and is flagged `unit_unconverted`.
8. Two dishes in a second-dish meal (PLN-13-05) both count.
9. No AI is used.

## Contract

Internal (`list.Needs(household, range, meals) []Need`), used by `GenerateList` (PLN-15-05) and PLN-15-10.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Range > 14 days or in the past | `INVALID_ARGUMENT` |
| No planned meals in range | empty proposal, not an error |

## Acceptance criteria

- [ ] A day plan with porridge (2 servings, household portions 2) and curry (4 servings, household portions 2):
      porridge × 1, curry × 0.5.
- [ ] With goals (R2) and pot 6 for the curry: curry × 1.5; Ola away → × 1.25.
- [ ] A cooked meal adds nothing.
- [ ] "1 onion" with a typical weight of 150 g becomes 150 g.
