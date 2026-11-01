# KIT-13-05 account.deleted: erase the person

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R1 (handler); data exists from R2 |
| Size | S |
| Status | todo |
| Depends on | KIT-11-01 |
| Unblocks | IDN-09 |

## Goal

When a person deletes their account, their personal rows in kitchen are erased everywhere, even in a household that
lives on without them.

## Business rules

1. Handler for `account.deleted` (`user_id`, `member_ids`):
   1. delete the member rows for `member_ids` (targets, consent, personal diet and dislikes go with them);
   2. set `user_id` to null in any remaining kitchen row that holds this `user_id` (none expected after step 1;
      the step guards future columns);
   3. bump `planning_version` and write `profile.changed` (`members`) for each household that lost a row.
2. If the person's households are deleted too, identity sends `household.deleted`, which KIT-13-04 handles; the order
   does not matter.
3. Stock and savings are household data, not the person's; they stay with a household that lives on.

## Contract

No RPC. Consumes `account.deleted`.

## Data

No schema change.

## Errors

None beyond the consumer's retry rules.

## Acceptance criteria

- [ ] Ewa deletes her account; her member row in Adam's household is gone, Adam's stock is unchanged.
- [ ] Replaying the event is a no-op.
