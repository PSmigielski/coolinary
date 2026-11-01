# PLN-18-03 Plan from deals without AI

| | |
|---|---|
| Service | planner |
| Epic | [PLN-18 Plan from deals](README.md) |
| Release | R2 |
| Size | M |
| Depends on | PLN-18-01, PLN-13-06 |
| Unblocks | PLN-18-02 |

## Goal

The library-based version of "plan from deals": the fallback for PLN-18-02 and the result when AI is off.

## Business rules

1. Runs the family-day library composer (PLN-13-06) with two changes:
   - a dish's "on hand" share counts stock plus staples as before; ingredients covered by a deal candidate count as
     `on_offer`;
   - comparator key 6 (offer score) is active: the share of bought ingredients that are on offer, weighted by shop
     position (`1 / position`).
2. Candidates with less than 30% of ingredients on hand are dropped when stock has usable items (fridge first).
3. Each chosen meal records the offer items it uses (best candidate per ingredient: earliest shop position, then
   lowest unit price).
4. Same response shape as PLN-18-02 with `ai_used = false`.

## Contract

Internal composer variant.

## Data

No schema change.

## Errors

None of its own; the library fallback inside `CreatePlanFromDeals` (PLN-18-02), whose errors apply.

## Acceptance criteria

- [ ] AI off: a deals plan comes back, using at least one offer when one fits.
- [ ] Two dinners with equal stock share: the one using an offer at shop position 1 wins over one using position 2.
