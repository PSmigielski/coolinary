# PLN-14-06 Week plan without AI

| | |
|---|---|
| Service | planner |
| Epic | [PLN-14 Week plan](README.md) |
| Release | R3 |
| Size | L |
| Depends on | PLN-14-02, PLN-14-03, PLN-14-04, PLN-14-05, PLN-13-06 |
| Unblocks | PLN-14-07 |

## Goal

Build a full week from the library with the fewest cooking sessions. It is the fallback for the AI week and works
with AI switched off.

## Business rules

1. Inputs: start date, rhythm and pins (PLN-14-02), profile snapshot, attendance for the 7 days, stock.
2. Main meals first: for each cooking day pick a main dish from library candidates that pass the shared limits for
   both its own eaters and the leftover meal's eaters (PLN-13-01), preferring dishes tagged "keeps 3 days" or
   "freezes well" so they can be doubled (PLN-14-03).
3. Stock is allocated across the whole week in date order, oldest items first: an item used on Monday is not on
   hand for Thursday. Items expiring before a meal's date are not used for it.
4. Search: for each cooking day keep the top 6 main candidates by stock score; choose days greedily in date order
   with a one-step look-ahead (try each of the 6 for day d with the best choice for the next cooking day), so the
   search stays under 2 seconds.
5. Breakfast and the other meal: "no cooking" dishes, chosen per day by PLN-13-06's day search with the main meal
   fixed; no breakfast dish more than twice in a week unless fewer than 4 candidates exist.
6. No dish as the main meal on more than 2 cooking days in one week (leftovers don't count).
7. Every week candidate is balanced (PLN-14-05), counted (PLN-14-04) and ranked by the comparator (PLN-13-03) with
   the rhythm allowance from PLN-14-02.
8. Saved with `ai_used = false` via PLN-14-01.

## Contract

Internal composer, used by `CreateWeekPlan` (PLN-14-07).

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| kitchen down | `UNAVAILABLE` |
| No main candidate for a cooking day | the day gets a no-cook main meal if one passes, else `no_candidate` |

## Acceptance criteria

- [ ] Test family of four, fixed fridge, "every other day", AI off: at most 4 cooking sessions, one dish per meal
      (unless an allergy or diet forces a second), every day within 10% and the week within 5% (Verification 18).
- [ ] A 300 g pack of mince is used by exactly one meal in the week.
- [ ] Runs in under 2 seconds with a 2,000-dish library.
