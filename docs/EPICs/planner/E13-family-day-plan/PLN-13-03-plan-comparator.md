# PLN-13-03 Plan comparator (priority order)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-13 Family day plan](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | PLN-12-05 |
| Unblocks | PLN-13-04, PLN-13-06, PLN-13-07, PLN-14-06, PLN-14-07, PLN-18-03 |

## Goal

When aims conflict, one fixed order decides (features.md "One dish per meal": hard limits; then the household's
cooking rhythm; then each person's day within 10%; then oldest stock first; then this week's offers). Every
composer uses this one comparator to pick between candidate plans.

## Business rules

1. A candidate is a day (or a week, PLN-14) with a dish per meal, its allocations (PLN-12-06 logic) and the stock it
   uses.
2. Candidates are compared lexicographically by this key, best first:
   1. hard-limit violations (PLN-07): 0 beats any; candidates with violations are never returned;
   2. cooking rhythm: sessions over the rhythm's allowance (PLN-14-02; a day plan allows 1 session for the main
      meal plus 0 for the others unless the household cooks other meals), fewer is better;
   3. day balance: number of people out of the 10% range (PLN-12-05), then the sum of their deviations in percent;
   4. eat-lighter score when the household eats lighter (PLN-13-04), higher is better;
   5. stock score: Σ over used stock items of `quantity share used × urgency`, `urgency = 1 / (1 + days to expiry)`
      (items without expiry count 0.05; items past their date count as 0 days, urgency 1.0, Q103), higher is better;
   6. offer score (PLN-18, R2): share of bought ingredients on offer at the household's shops, higher is better;
   7. fewer distinct ingredients to buy;
   8. stable tie-break: lower sum of dish ids' hash, so the same inputs always give the same plan.
3. Extra cooking from toppings needing their own pan (PLN-13-02) and second dishes (PLN-13-05) count as sessions in
   step 2.
4. The comparator explains its choice: for the winner it returns the first key where it beat the runner-up, used in
   logs and tests (never shown with personal numbers).

## Contract

Internal package `planner/internal/compose` (`Compare(a, b Candidate) int`, `Score(c) Key`). No RPC.

## Data

No schema change.

## Errors

None of its own; internal comparator called by PLN-13-04, PLN-13-06, PLN-13-07, PLN-14-06, PLN-14-07 and PLN-18-03.

## Acceptance criteria

- [ ] A candidate with an allergen violation never wins, however good its other scores.
- [ ] Between two valid days, the one with 1 cooking session beats the one with 2 even if it uses less old stock.
- [ ] Between two days with equal sessions, the one where nobody is out of the 10% range wins.
- [ ] Same inputs twice → same winner (deterministic).
- [ ] Unit tests build candidates for each key and check the order.

## Notes

- Decided (Q205): "eat lighter" sits at step 4, after each person's day target and before using up stock.
- Decided (Q206): the urgency formula in step 5 (share used × 1/(1 + days to expiry); undated items 0.05).
