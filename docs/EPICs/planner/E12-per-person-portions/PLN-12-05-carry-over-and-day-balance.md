# PLN-12-05 Carry-over and day balance

| | |
|---|---|
| Service | planner |
| Epic | [PLN-12 Per-person portions](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | PLN-12-03, PLN-12-04 |
| Unblocks | PLN-12-06, PLN-13-03, PLN-14-05 |

## Goal

Small gaps carry over: a bigger dinner makes the next meal smaller. Each person's planned day lands within 10% of
their target (features.md "How one meal is split", step 6).

## Business rules

1. Meals of a day are processed in time order: breakfast, dinner, supper.
2. For each eater with a target: `effective share = meal target + carry-in`. Carry-in to breakfast is 0 in a day
   plan; week plans pass the previous day's remainder (PLN-14-05).
3. After portion and extras: `gap = effective share − (portion kcal + extras kcal)`. A positive gap means short, a
   negative gap means over. The gap becomes the carry-in of the person's next meal that day.
4. Carry-in into one meal is capped at ±25% of that meal's target; any rest is dropped and the person is flagged.
5. A meal the person is away for is skipped: it receives no carry-in, and the carry passes on to the next meal
   they eat.
6. Carry only moves forward in time. Already planned earlier meals never change because of a later meal.
7. **Day check:** `planned day = Σ (portion kcal + extras kcal)` over the meals the person eats.
   `day target part = meal part of the day target (PLN-12-02) − meal targets of meals they are away for`.
   The day is in range when `|planned day − day target part| ≤ 10% × day target part`. Otherwise the person is
   flagged `day_out_of_range` with the deviation in percent (whole number).
8. A person away for every meal has no day check.
9. Eaters without a target are never checked.
10. Output per person per day: planned kcal, target part, deviation %, in range, and per meal the carry-in used.

## Contract

Internal package `planner/internal/portions` (`BalanceDay(day, eaters) DayBalance`). No RPC.

## Data

No schema change (stored by PLN-12-06).

## Errors

None of its own; invalid inputs are rejected by PLN-12-02 and PLN-12-03.

## Acceptance criteria

- [ ] Worked example: Ewa's dinner is 1½ portions = 780 against a 720 share; her supper share drops from 450 to
      390 ("supper is 60 kcal smaller").
- [ ] Adam's dinner is 20 kcal short; his supper share rises by 20.
- [ ] Kuba after milk and banana (about 265 kcal for a 240 gap) carries −25 into supper.
- [ ] A carry of +400 into a 300 kcal supper is capped at +75 and flagged.
- [ ] A test family day where all three meals fit lands within 10% for every person; a day built from a too-large
      breakfast for a child is flagged `day_out_of_range`.
- [ ] Ola away for dinner: her breakfast gap carries straight to supper.

## Notes

- Decided (Q202): the day check compares planned meals with 90% of the target (target minus the 10% snacks share),
  because snacks are not planned.
- Decided (Q201): the ±25% carry cap (rule 4).
