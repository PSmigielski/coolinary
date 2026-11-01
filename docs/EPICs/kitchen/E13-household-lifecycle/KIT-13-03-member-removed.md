# KIT-13-03 member_removed: delete the person's data

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R2 |
| Size | S |
| Depends on | KIT-11-01 |
| Unblocks | IDN-08 (leave and remove) |

## Goal

When someone leaves or is removed, their personal data leaves with them, and the household's shared fridge stays
(features.md §2: "Anyone can leave, and then starts again with their own empty fridge").

## Business rules

1. Handler for `household.member_removed` (`reason`: `left`, `removed`, `profile_deleted`, `plan_ended`):
   1. delete the member row; its personal restrictions and dislikes go with it (foreign keys, KIT-11-01);
   2. stock, profile sections and savings stay with the household;
   3. bump `planning_version`, write `profile.changed` (`members`).
2. The person's new household (`new_household_id`) gets its own `household.created` and default profile; nothing is
   copied there: they start with an empty fridge.
3. A missing row is not an error.

## Contract

No RPC. Consumes `household.member_removed`.

## Data

No schema change.

## Errors

None beyond the consumer's retry rules.

## Acceptance criteria

- [ ] Removing Ewa deletes her target and allergens; Adam's fridge is unchanged (Verification 15: after removal she
      no longer sees it — identity's part).
- [ ] A profile without an account deleted by a parent (`profile_deleted`) loses its row.
- [ ] Replaying the event is a no-op.
