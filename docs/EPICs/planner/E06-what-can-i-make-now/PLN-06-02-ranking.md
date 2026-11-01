# PLN-06-02 Ranking: share on hand, expiry, tie-breaks

| | |
|---|---|
| Service | planner |
| Epic | [PLN-06 What can I make now](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-06-01, PLN-05-04, PLN-05-06 |
| Unblocks | PLN-06-03, PLN-08-02, PLN-08-05 |

## Goal

Order dishes so the ones that use most of what is at home, and especially what expires first, come first; never repeat
yesterday's dinner at the top.

## Business rules

1. Expiry urgency of one stock item by days left (household local date): past its date or 0 (today) → 1.0; 1 → 0.8; 2 → 0.6; 3 → 0.4;
   4–5 → 0.2; more, or no date → 0.
2. Dish urgency = the highest urgency among the stock items it would use (PLN-06-01 rule 8), plus 0.05 for each other
   used item with urgency ≥ 0.6, capped at 1.0.
3. Score = 0.7 × share_on_hand + 0.3 × dish urgency
   + 0.05 if the dish is a household favourite
   − 0.20 if cooked in the last 3 days, − 0.10 if cooked 4–7 days ago (PLN-05-06).
4. Tie-breaks in order: fewer missing lines; shorter total time; higher POPULAR score (PLN-05-01 rule 4); name.
5. Groups for the app: `READY` (share_on_hand = 1.0), `ALMOST` (share ≥ 0.75 or ≤ 2 missing lines), `NEEDS_SHOPPING`
   (the rest).
6. Dishes with share < 0.30 are left out, unless that leaves fewer than 3 dishes for the meal; then the best dishes
   below 0.30 fill up to 3, in group `NEEDS_SHOPPING`.
7. Decided (Q104): the weights above (0.7/0.3, +0.05 favourite, −0.20/−0.10 recently cooked, 0.30 cut-off) are
   config values, tuned with the concierge-test data.
8. Library dishes are held in an in-process cache of current published versions with their lines (refreshed on publish,
   retire and every 5 minutes), so a request computes ranking in memory without a database scan per dish.

## Contract

Internal; exposed through PLN-06-03.

## Data

No schema change.

## Errors

None of its own; internal logic called by `SuggestDishes` (PLN-06-03) and the day-plan composer (PLN-08-02, PLN-08-05).

## Acceptance criteria

- [ ] Two dishes both fully covered: the one using milk that expires tomorrow ranks first.
- [ ] A dish cooked yesterday drops below an equal dish not cooked this week.
- [ ] With 10 items in the fridge from `planner/testdata/fridges/*.yaml`, every meal returns at least 3 dishes
      (Verification 1).
- [ ] Ranking 2,000 dishes for one request takes < 30 ms on the dev machine.
