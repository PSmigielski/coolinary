# PLN-08-03 Ingredient allocation across meals

| | |
|---|---|
| Service | planner |
| Epic | [PLN-08 AI day plan](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-06-01 |
| Unblocks | PLN-07-04, PLN-08-04, PLN-08-05, PLN-14, PLN-15 |

## Goal

Make sure breakfast doesn't use up what dinner needs (business model §1): stock is shared out between the day's meals
in a fixed, explainable way.

## Business rules

1. Allocation order of meals: the household's main hot meal first, then breakfast, then the remaining meal. (The main
   meal needs the most and has the fewest substitutes.)
2. For each meal's dish (scaled to its servings), each non-optional line takes stock items of that ingredient, oldest
   expiry first, until the required amount is met. Taken amounts are subtracted from a working copy of stock before the
   next meal.
3. Staples are never allocated (always available).
4. Each meal's allocation records, per line: stock item ids and amounts taken, shortfall. A shortfall > 10% of the
   required amount is reported to the checker as `STOCK_SHORT` (PLN-07-04).
5. The allocation is returned with the plan (not stored separately) and is used later by the shopping list (what's
   still missing for the plan, PLN-15) and by "cooked it" (which items to consume first, PLN-10-02).
6. Pure function; same input → same output (ties on expiry broken by added date, then stock item id).
7. Items past their date are taken first like any other old item (Q103). A meal whose allocation takes one gets
   `check_first_stock_item_ids`; the flag is stored on the plan meal (PLN-08-01) and shown in cooking mode (PLN-10-01).

## Contract

Internal Go package `allocation`.

## Data

No schema change.

## Errors

None of its own; internal logic called by PLN-07-04, PLN-08-04, PLN-08-05, PLN-14 and PLN-15.

## Acceptance criteria

- [ ] Stock 6 eggs (2 expire tomorrow, 4 next week); dinner (main) needs 4, breakfast needs 2: dinner gets the 2 old +
      2 new, breakfast the remaining 2; no shortfall.
- [ ] Same stock, breakfast needs 4 too → breakfast shortfall 2 (50%) → `STOCK_SHORT` on breakfast only.
- [ ] Changing the main meal to supper changes the allocation order.
