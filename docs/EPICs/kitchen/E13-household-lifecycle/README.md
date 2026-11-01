# KIT-13 Household lifecycle: event consumers and data export

| | |
|---|---|
| Service | kitchen |
| Release | R1 (create, delete, merge stock, plan effects, export, retention); R2 (members leaving, person data carried over) |
| Depends on | KIT-01 (stream consumer, plan from identity), KIT-07, KIT-10, KIT-11, IDN-01-03 (event payloads) |

## Goal

Kitchen's household data follows what happens in identity — a household is created, someone joins with their fridge,
leaves, is removed, the plan changes, an account is deleted — without kitchen ever calling identity, and a user can
export everything kitchen holds about their household (GDPR; derived decisions D9).

## Flow

```
events.identity ──► KIT-01 consumer (idempotent, in stream order)
   household.created ─────────► default kitchen profile (KIT-13-01)
   household.member_added ────► member row (KIT-11-01); brought_fridge_from → move stock;
                                previous_member_id → carry the person's data (KIT-13-02)
   household.member_removed ──► delete that person's data (KIT-13-03)
   household.plan_changed ────► planning version, profile.changed; on a downgrade the allergen merge (KIT-13-06)
   household.deleted ─────────► erase everything of the household (KIT-13-04)
   account.deleted ───────────► erase the person's rows (KIT-13-05)
gateway ─► ExportHouseholdData (KIT-13-07)
daily ─► retention clean-up (KIT-13-08)      catalogue merge ─► household tables re-pointed (KIT-13-09)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-13-01](KIT-13-01-household-created.md) | household.created: default profile and stock state | R1 | S | KIT-01-04, KIT-10-01 |
| [KIT-13-02](KIT-13-02-member-joined-with-fridge.md) | Joining with a fridge: move stock, carry the person's data | R1 (person data R2) | M | KIT-13-01, KIT-07-03, KIT-11-01 |
| [KIT-13-03](KIT-13-03-member-removed.md) | member_removed: delete the person's data | R2 | S | KIT-11-01 |
| [KIT-13-04](KIT-13-04-household-deleted.md) | household.deleted: erase the household | R1 | M | KIT-07, KIT-10 |
| [KIT-13-05](KIT-13-05-account-deleted.md) | account.deleted: erase the person | R1 | S | KIT-11-01 |
| [KIT-13-06](KIT-13-06-plan-changed.md) | Plan changes: planning version and effects | R1 | S | KIT-01-06, KIT-10-09, KIT-11-04 (R2 merge) |
| [KIT-13-07](KIT-13-07-export-household-data.md) | ExportHouseholdData | R1 | M | KIT-07, KIT-10, KIT-11-06, KIT-12 |
| [KIT-13-08](KIT-13-08-retention-job.md) | Retention clean-up job | R1 | S | KIT-07, KIT-08, KIT-09 |
| [KIT-13-09](KIT-13-09-ingredient-merge-household-tables.md) | Register household tables for ingredient merges | R1 | S | KIT-02-09, KIT-07-01, KIT-10 |

## Decisions (taken 1 November 2026)

- Q077 (KIT-13-02): Yes, 1,000 items, as an abuse guard.
- Q139 (KIT-13-02): Yes, if they were that household's only account member.
- Q140 (KIT-13-02): Move the soonest-expiring items first, up to 1,000; the rest stay behind.
- Q141 (KIT-13-08): 24 months; monthly savings totals are kept longer.
