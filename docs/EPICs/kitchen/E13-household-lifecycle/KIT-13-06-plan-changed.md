# KIT-13-06 Plan changes: planning version and effects

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-01-04 (consumer), KIT-01-06 (plan from identity), KIT-10-09, KIT-11-04 (R2: allergen merge) |
| Unblocks | PLN-01-03 (fresh snapshot after a plan change) |

## Goal

When a household goes from Free to Plus or back, every kitchen rule that depends on the plan switches at once, and
no data is lost on a downgrade.

## Business rules

1. Handler for `household.plan_changed` v1: bump `planning_version` and write `profile.changed` with section
   `plan`, so planner drops its cached snapshot. The event stores no plan (Decided, Q004: kitchen asks identity,
   KIT-01-06).
2. On a downgrade (`plan = free`, `previous_plan = plus`, R2) the same transaction runs the allergen merge of
   KIT-11-04 rule 4 (Decided, Q195).
3. Effects, each reading the plan from identity at the moment it runs:

   | Rule | Plus | Free |
   |---|---|---|
   | Expiry reminders (KIT-09) | sent | not sent |
   | Per-person targets of others, per-person diets (KIT-11-02) | editable | not editable |
   | Per-person allergens in planning (KIT-11-04) | applied per person | merged into the household allergens |
   | Per-person diets and dislikes in planning (KIT-11-04) | applied | not applied (kept, apply again with Plus) |
   | Savings reports, push, card (KIT-12) | yes | current month only |

4. A downgrade deletes nothing: targets, per-person data and history stay and come back with Plus. The allergen
   merge only adds to the household list.
5. Request-time checks ask identity too (KIT-01-06 rule 3), so they switch even before this event arrives.

## Contract

No RPC. Consumes `household.plan_changed`.

## Data

No schema change.

## Errors

None beyond the consumer's retry rules.

## Acceptance criteria

- [ ] Plus → Free: the next reminder tick sends nothing; `GetPlanningProfile` shows `plan = free` with a higher
      `profile_version`; member targets are still stored.
- [ ] Free → Plus: reminders start at the next tick after 09:00 local.
- [ ] Plus → Free with Kuba nut-free and Ola vegetarian (R2): the household allergens now include nuts; Ola's
      vegetarian no longer limits plans; both rows are still stored and editing them returns `plus_required`.
- [ ] Replaying the downgrade event changes nothing further.
