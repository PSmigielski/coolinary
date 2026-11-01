# KIT-13-04 household.deleted: erase the household

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-07, KIT-08-05, KIT-09-02, KIT-10, KIT-11, KIT-12 (their tables), KIT-04/KIT-05 (household products erase hook) |
| Unblocks | IDN-08, IDN-09 (deletion completes in kitchen) |

## Goal

When identity deletes a household (owner deleted the account, last person left, anonymous clean-up, merge), every
row kitchen holds for it is deleted, so nothing personal outlives it (GDPR; D9).

## Business rules

1. Handler for `household.deleted` (any `reason`), in one transaction:
   1. stock: `stock_items`, `stock_movements`, `household_stock_state`, `stock_add_refs`, `cook_events`;
   2. profile: `kitchen_profiles` and, by cascade, cuisines, tools, restrictions, dislikes, staples, member rows;
   3. reminders and savings: `reminder_runs`, `savings_monthly`;
   4. receipts: `receipt_name_links`, `receipt_fingerprints`;
   5. the household-scoped tables of KIT-04 and KIT-05 (household products, pending namings) through the erase hook
      those tasks register; shared-catalogue products stay but lose any link to this household.
2. Redis: `kitchen:planprof:<household_id>:*` and any `kitchen:receipt:*` copies of this household are deleted
   (the receipt copies are found through a per-household Redis set kept by KIT-08-04).
3. No `stock.changed` or `profile.changed` is written for a deleted household.
4. Large households: deletion in one transaction must finish inside the 10 s handler deadline; a household holds at
   most 1,000 items and 24 months of movements, which a test measures (target: < 2 s).
5. A household that does not exist in kitchen is a no-op.

## Contract

No RPC. Consumes `household.deleted`.

## Data

No schema change. KIT-08-04 adds a Redis set `kitchen:receipts:<household_id>` holding its receipt keys (TTL 2 h).

## Errors

| Case | Behaviour |
|---|---|
| Transaction fails | retried by the consumer; dead letter after 10 tries alerts through KIT-01-05 |

## Acceptance criteria

- [ ] After the event, a query over every kitchen table for the household id returns 0 rows (test lists the tables
      from the schema, so a new table without deletion fails the test).
- [ ] A shared product the household named stays in the catalogue without the household link.
- [ ] Deleting a household with 1,000 items and 50,000 movements takes < 2 s on the dev machine.
