# PLN-02-05 Portion scaling rules

| | |
|---|---|
| Service | planner |
| Epic | [PLN-02 Recipe library model](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-02-01, KIT-02 (GetIngredients) |
| Unblocks | PLN-02-06, PLN-05, PLN-10, PLN-12, PLN-15 |

## Goal

One function scales a dish from its base servings to any number of portions with amounts a cook can measure. Dish
detail, cooking mode, the pot in family plans and the shopping list all use it.

## Business rules

1. `scaled = amount × servings / base_servings`, computed with exact decimals, then rounded per unit:
   - grams and millilitres: below 10 → nearest 1; 10–100 → nearest 5; 100–1,000 → nearest 10; above 1,000 → nearest 50.
   - pieces: below 2 → nearest ¼; 2–6 → nearest ½; above 6 → nearest 1. Never below ¼.
   - spoons: nearest ½; a scaled amount of 6 tablespoons or more is shown as tablespoons still (no conversion to
     millilitres in recipes).
2. Ingredients that kitchen marks as whole-only (eggs, for example) round to whole pieces, rounding half up, and never
   to 0: minimum 1. Decided (Q089): kitchen's `GetIngredients` returns a `whole_units_only` flag (KIT-02) that
   drives this rule.
3. Allowed servings: ¼ steps from 1 to 24 (the pot may be up to 12 portions and doubled on cooking days, PLN-14).
   Outside → INVALID_ARGUMENT.
4. Step timers and `total_minutes` do not scale. Decided (Q088): when servings > 2 × base, the response carries the
   hint "cooking times may be longer".
5. Scaling the shopping list and stock consumption uses the unrounded amount (PLN-15, PLN-10-02); only display and
   cooking mode use the rounded value.
6. Nutrition per portion does not change with scaling; the total does (PLN-04-03).

## Contract

Go package `portions`: `Scale(lines, base, servings) []ScaledLine` with `Exact` and `Display` quantities. Exposed
through `GetDish` (`servings` field, PLN-02-06).

## Data

No schema change.

## Errors

None of its own; a pure function called by PLN-02-06, PLN-05, PLN-10, PLN-12 and PLN-15. `servings` ≤ 0 or a base of 0 is a programming error caught by tests; callers validate servings first.

## Acceptance criteria

- [ ] 400 g pasta for 4 scaled to 6 → 600 g; 3 garlic cloves for 4 scaled to 6 → exact 4.5, display 4½.
- [ ] 2 eggs for 4 scaled to 3 → exact 1.5, display 2 (whole-only, half up).
- [ ] 7 g salt for 4 scaled to 6 → exact 10.5, display 10 (10–100 g rounds to the nearest 5).
- [ ] Servings 0.5 and 25 are rejected.
