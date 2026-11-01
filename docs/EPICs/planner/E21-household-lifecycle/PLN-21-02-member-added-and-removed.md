# PLN-21-02 Member added and removed

| | |
|---|---|
| Service | planner |
| Epic | [PLN-21 Household lifecycle](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | PLN-01, PLN-12-06, PLN-12-01 |
| Unblocks | — |

## Goal

When someone joins or leaves the household, future meals include or drop them, and portions and the list follow.

## Business rules

1. `household.member_added`: invalidate the profile snapshot cache; recompute portions for the household's meals from
   today on (PLN-12-06). The new member eats every meal by default. Their own fridge items arrive through kitchen
   (`brought_fridge_from`), not planner.
2. `household.member_removed` (`left`, `removed`, `profile_deleted`): delete their attendance marks and rules,
   remove them from future meals' portions and plates, recompute those meals (PLN-12-06), which updates the list
   (PLN-15-10).
3. For `left` / `removed` with an account: list lines they added stay; `added_by` keeps their member id but the app
   shows "former member" because the member is gone from `GetHousehold`.
4. The person who left starts with an empty plan and list in their new household (`new_household_id`): nothing is
   copied.
5. Past and cooked meals don't change.

## Contract

Event consumers only.

## Data

No schema change.

## Errors

Handler errors roll back and are retried by the consumer framework; after 5 failures the event goes to `dead_events` (PLN-01-01). A member already added or already removed is a success (upsert).

## Acceptance criteria

- [ ] Kuba removed: tomorrow's dinner pot drops from 6 to 4 and the curry lines shrink.
- [ ] A new member added today appears in tonight's supper if it is not cooked yet.
