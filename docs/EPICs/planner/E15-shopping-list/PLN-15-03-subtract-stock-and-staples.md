# PLN-15-03 Subtract stock and staples

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-15-02, KIT-07 (ListStock), KIT-10 (GetPlanningProfile) |
| Unblocks | PLN-15-04, PLN-15-10 |

## Goal

Step 2: remove what is already in the fridge, freezer and pantry, and the "always at home" staples, so nothing
already at home appears on the list (features.md "Shopping list" step 2, Verification 12).

## Business rules

1. Staples from the profile snapshot (salt, oil, flour, …) are never on the list from plan needs. They can still be
   added by hand (PLN-15-06).
2. Stock comes from KIT-07 `ListStock` (all storage places), converted to base units like the needs.
3. Allocation runs per ingredient in meal date order (then breakfast, dinner, supper): each need takes from stock
   items with the earliest expiry first; what is taken is not available to later needs.
4. A stock item counts for a need only if it has no expiry date or its expiry date is on or after the meal's date.
5. A need is covered fully, partly (the remainder goes on), or not at all.
6. If stock of the ingredient exists only in a unit that cannot be converted (e.g. pieces with no typical weight),
   the need stays and the line gets the flag `check_at_home`.
7. Stock already reserved by earlier meals of the same plan (PLN-08's "from stock" choice) is respected: the
   allocation is the same one the plan used.
8. Optional ingredients are subtracted the same way.

## Contract

Internal (`list.SubtractStock(needs, stock, staples) []Need`).

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| kitchen down | `GenerateList` returns `UNAVAILABLE`; no list from unknown stock |

## Acceptance criteria

- [ ] Fixed test plan and test fridge (Verification 12): nothing already at home appears; a 300 g need against 200 g
      at home leaves 100 g.
- [ ] Yoghurt expiring the day before the meal does not count.
- [ ] Two meals needing 200 g onion each with 300 g at home: 100 g on the list, attributed to the later meal.
- [ ] Salt is never proposed even if a dish needs 30 g.

## Notes

- Decided (Q127): rule 4 (stock expiring before the meal date does not count, so the item goes on the list). Items
  already past their date are still offered in plans with "check it first" (Q103); the list buys a replacement.
