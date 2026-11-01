# PLN-21-03 Plus → Free downgrade

| | |
|---|---|
| Service | planner |
| Epic | [PLN-21 Household lifecycle](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | PLN-01-01, PLN-01-03, PLN-01-05, PLN-12-06, KIT-11 (allergen merge on downgrade) |
| Unblocks | — |

## Goal

When a household drops from Plus to Free (`household.plan_changed`), Plus-only planning stops cleanly without losing
data.

## Business rules

1. Plus checks ask identity on every request (PLN-01-05, Q004), so a downgrade takes effect on the next request; this
   handler only recomputes future meals.
2. Future meals are recomputed with Free rules: one personal goal only (PLN-12-02 rule 6). Per-person vegetarian,
   vegan and dislikes stop applying; per-person allergens and the gluten-free, lactose-free and nut-free flags keep
   applying because kitchen merges them into the household list on the downgrade (KIT-11-04, Q195, Q230), and
   planner filters by the household list from the profile snapshot.
3. Members of a shared household stay (Q187): planner removes nobody from plans; eaters are unchanged, only goals
   and per-person diets fall back to Free rules.
4. Existing week plans, deals plans and top-ups stay readable; creating or editing them returns `plus_required`.
5. The price matrix (PLN-19) and cheapest top-up (PLN-20) return `plus_required`.
6. Offer badges keep working with the shops offers now returns (2 on Free).
7. Upgrading back to Plus recomputes future meals with all targets.

## Contract

Event consumer only.

## Data

No schema change.

## Errors

Handler errors roll back and are retried by the consumer framework; after 5 failures the event goes to `dead_events` (PLN-01-01). A plan change to the plan already held is a success.

## Acceptance criteria

- [ ] After a downgrade, tomorrow's family dinner uses only the account owner's target; the others get 1 portion.
- [ ] `CreateWeekPlan` returns `plus_required`; `GetWeekPlan` still works.
- [ ] After a downgrade, a member's own nut allergy (now in the household list) still blocks nut dishes; their
      vegetarian diet no longer filters dishes (Q195).
- [ ] After a downgrade no member is removed from future plans (Q187).
