# PLN-21-01 Household and account deletion

| | |
|---|---|
| Service | planner |
| Epic | [PLN-21 Household lifecycle](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-01, PLN-15-01 |
| Unblocks | IDN-09 (deletion completes across services) |

## Goal

When a household or an account is deleted, planner erases what it holds about them (GDPR erasure; README events
`household.deleted`, `account.deleted`).

## Business rules

1. `household.deleted`: delete, for that household, all day and week plans, plan meals and their allocations,
   attendance marks and rules, the shopping list with lines, sources, badges, proposals and mutation records,
   household-private dishes (AI, top-up variants), favourites, ratings, cooking history, offer reminders, top-ups,
   Shared library dishes are not touched.
2. Ratings that fed the shared library's review (PLN-11) are kept only as anonymous aggregates (count and average per
   dish) without household or member ids.
3. `account.deleted`: for each `member_id` in the event, delete that member's attendance marks and rules, portion
   rows and plate toppings in future meals (and recompute those meals, PLN-12-06), and their personal favourites and
   ratings; replace `added_by` / `ticked_by` with null on list lines.
4. Deletion runs in one transaction per household; on failure the event is retried (at-least-once) and the handler
   is idempotent (already-deleted data is fine).
5. Completion is logged with the event id and row counts only.
6. Within 30 days nothing of the household remains in planner's database; backups age out on their own schedule.

## Contract

Event consumers only (PLN-01 consumer framework).

## Data

All household-scoped tables are checked for an index on `household_id` (add missing ones here).

## Errors

Handler errors roll back and are retried by the consumer framework; after 5 failures the event goes to `dead_events` and an alert is raised (PLN-01-01), because an unfinished erasure must be seen. Deleting data that is already gone is a success.

## Acceptance criteria

- [ ] After `household.deleted`, a query over every planner table for that household id returns 0 rows (test
      enumerates tables from the schema, so new tables are covered).
- [ ] Replaying the same event does nothing and succeeds.
- [ ] After `account.deleted` for Ewa, future meals no longer include her and their pots shrink; past list lines
      show no "added by".
