# PLN-14-03 Cook once, eat twice

| | |
|---|---|
| Service | planner |
| Epic | [PLN-14 Week plan](README.md) |
| Release | R3 |
| Size | M |
| Depends on | PLN-14-02, PLN-12-08, PLN-02 |
| Unblocks | PLN-14-04, PLN-14-05, PLN-14-06, PLN-14-07, PLN-15-02, PLN-10 |

## Goal

On a cooking day the main dish is cooked in a double amount if it keeps 3 days or freezes well; the second half
becomes a later meal and is split again for whoever eats then (features.md "One dish per meal, few cooking sessions").

## Business rules

1. A dish may be doubled only with effort tag `keeps_3_days` ("keeps 3 days in the fridge") or `freezes_well` (PLN-02).
2. The leftover meal is the next main meal without its own cooking session:
   - with "keeps 3 days": allowed on the cooking day + 1 or + 2, stored in the fridge;
   - with "freezes well": any later day in the same week, stored in the freezer; when it is more than 2 days later
     the dish card and the day before show "take out of the freezer the evening before".
   - with both tags, the fridge is used when the gap is ≤ 2 days.
3. The session's pot = pot of the cooking meal + pot of the leftover meal, each computed with its own eaters and
   carry-over (PLN-12). Ingredients are scaled to the session pot (PLN-12-08); the leftover meal adds no ingredients.
4. "Cooked it" on the cooking meal (PLN-10) consumes stock for the whole session pot.
5. The leftover pot is fixed once cooked. If eaters of the leftover meal change after cooking, portions are re-split
   within the fixed pot: compute portions as usual, and if their sum exceeds the pot, scale everyone's portion down
   by the same factor, rounded down to quarters, minimum ½; any shortfall is filled with extras (PLN-12-04).
   Spare portions are reported ("1 portion left over").
6. Before cooking, a change to the leftover meal's eaters recomputes the session pot.
7. Base-plus-toppings dishes keep the same topping groups for the leftover meal unless eaters' plates change.
8. A leftover meal is shown as "Leftovers: chicken curry (cooked Monday)".

## Contract

Fields in `CookingSession` (PLN-14-01). Internal `compose.PlanLeftovers(days, sessions)`.

## Data

No schema change beyond PLN-14-01.

## Errors

None of its own; internal logic inside the week composers (PLN-14-06, PLN-14-07).

## Acceptance criteria

- [ ] Worked example with "every other day": Monday's curry pot = 6 + 6 = 12, Tuesday dinner is leftovers split
      again (Adam 1½, Ewa 1½, Ola 1, Kuba 2) (features.md worked example).
- [ ] A dish with neither tag is never doubled.
- [ ] "Freezes well" leftovers on Saturday from a Wednesday session show the defrost note on Friday.
- [ ] After cooking, Kuba away on Tuesday: the 6 leftover portions are split among three (1½ + 1½ + 1) and 2 portions are spare.
