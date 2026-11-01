# PLN-06-01 Stock coverage per dish

| | |
|---|---|
| Service | planner |
| Epic | [PLN-06 What can I make now](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-01-03, PLN-01-04, PLN-02-05, KIT-02 (GetIngredients: typical weights, density) |
| Unblocks | PLN-06-02, PLN-07-04, PLN-08-03, PLN-15 |

## Goal

For a dish at the household's number of portions, work out which ingredients are at home, how much of each, and what
is missing — the base of "what can I make now", plan checks and the shopping list.

## Business rules

1. Required amounts = the dish's non-optional lines scaled to the household's portions (profile, PLN-01-03) using the
   exact (unrounded) values of PLN-02-05.
2. Available amount of an ingredient = sum of all the household's stock items with that `ingredient_id`, including
   items past their date (household time zone). Decided (Q103): items past their date count as at home; every line
   they cover is flagged `check_first` with those stock item ids, and the flag travels to suggestions, plans and
   cooking mode so the app shows "check it first".
3. Amounts are compared in grams: grams as is; millilitres × density; pieces × typical piece weight; spoons × typical
   spoon weight (all from `GetIngredients`, cached 1 day). If a conversion is not possible (no typical weight) or the
   stock item has no quantity ("some cheese"), the ingredient counts as fully available when any item is present.
4. Staples from the profile ("always at home", features.md §1) count as fully available and are never missing.
5. Coverage per line = min(1, available / required). A topping group counts as one line with the coverage of its best
   option. `optional` lines are ignored.
6. `share_on_hand` = average line coverage, 0.0–1.0, rounded to 2 decimals in responses (unrounded for ranking).
7. Missing items = lines with coverage < 1, each with `missing = required − available` in the recipe's unit (display
   rounding PLN-02-05) and whether any of it is at home.
8. Stock items used = for each covered line, the stock items that would be used, oldest expiry first (features.md
   "oldest items first"), so items past their date are used first; their days to expiry (negative when past the
   date) feed PLN-06-02. Each line's result carries `check_first` (bool) and `past_date_stock_item_ids`.
9. Unmatched stock items (no ingredient id) never count (PLN-01-04 rule 4).
10. The calculation is a pure function over (dish lines, stock snapshot, ingredient data, staples, servings) so ranking,
    the rule checker and the list can share it and test it with fixtures.

## Contract

Internal Go package `coverage`. No RPC.

## Data

No schema change.

## Errors

None of its own; internal logic called by PLN-06-02, PLN-07-04, PLN-08-03 and PLN-15.

## Acceptance criteria

- [ ] Curry for 4 needs 600 g chicken; stock has 400 g → line coverage 0.67, missing 200 g.
- [ ] Stock "2 onions" against a recipe line of 150 g onion with typical weight 110 g/piece → 220 g available → covered.
- [ ] Salt (a staple) with no stock is covered and not missing.
- [ ] A yoghurt one day past its date counts as at home, and its line has `check_first = true` with that item's id.
- [ ] A line covered only by items within their date has `check_first = false`.
- [ ] Tacos with toppings beef (absent) or beans (present) count the topping group as covered.
