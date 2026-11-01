# PLN-07-04 Stock and allowed-ingredient check

| | |
|---|---|
| Service | planner |
| Epic | [PLN-07 Rule checker](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-07-01, PLN-06-01 |
| Unblocks | PLN-07-05, PLN-08-03, PLN-08-04 |

## Goal

An AI dish that claims to cook "from what's in the fridge" really uses only what's there (plus staples and the items
the user agreed to buy), in the amounts that are there — across all meals of the plan, not each meal alone.

## Business rules

1. Mode `FROM_STOCK` (day plans, swaps, wishes): every non-optional ingredient must be at home (stock or staple) or in
   the allowed purchases; otherwise → `NOT_ALLOWED_INGREDIENT`.
2. Allowed purchases come from the request: the 1–3 items of the cuisine-mismatch choice (PLN-08-06), offers in plan
   from deals (PLN-18), or the user's explicit "I can buy X".
3. Ingredient ids must be in the candidate set sent to AI (stock ingredients, staples, allowed purchases, and the
   pantry list sent with the request); an id outside that set but in the catalogue → `NOT_ALLOWED_INGREDIENT`;
   an id unknown to the catalogue → `UNKNOWN_INGREDIENT` (PLN-07-05).
4. Quantities are checked over the whole plan: the plan-level allocation (PLN-08-03) sums each ingredient's required
   amount over all meals; if the sum exceeds the available amount by more than 10%, the meals involved get
   `STOCK_SHORT` (ingredient, needed, available). Between 0% and 10% over → warning `LOW_STOCK_MARGIN`.
   A meal that uses stock past its date gets the warning `CHECK_FIRST` with those stock item ids; it never blocks the
   dish (Q103).
   Decided (Q107).
5. Ingredients whose amount can't be converted (PLN-06-01 rule 3) pass the quantity check when present.
6. Mode `ANY` (library browse, rewrites of library dishes): only rule 3's catalogue part applies.

## Contract

Internal.

## Data

No schema change.

## Errors

None of its own; internal check run by the checker core (PLN-07-01). Stock it needs comes from the snapshot already loaded by the caller.

## Acceptance criteria

- [ ] Stock 6 eggs; breakfast uses 4, supper uses 4 → both meals get `STOCK_SHORT eggs` (8 > 6 × 1.1).
- [ ] Day plan using rice that isn't at home but is an allowed purchase passes.
- [ ] Day plan using pine nuts, not at home and not allowed → `NOT_ALLOWED_INGREDIENT`.
