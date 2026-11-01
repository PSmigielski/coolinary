# IDN-08-03 Transfer ownership

| | |
|---|---|
| Service | identity |
| Epic | [IDN-08 Leaving and removing members](README.md) |
| Release | R2 |
| Size | S |
| Depends on | IDN-04-01, IDN-04-04 |
| Unblocks | IDN-08-01, GW-04 |

## Goal

The owner can hand the household to another member, for example before leaving.

## Business rules

1. `TransferOwnership(household_id, actor_user_id, new_owner_member_id)`: only the owner.
2. The new owner is an account member of the household and not anonymous.
3. Uses `setOwner` (IDN-04-01): the old owner becomes `member`, the new one `owner`, and
   `households.owner_user_id` changes, in one transaction.
4. The plan stays with the household (see IDN-08-01 rule 3).
5. No event: no other service stores who the owner is.
   Decided (Q190): no event.

## Contract

```proto
rpc TransferOwnership(TransferOwnershipRequest) returns (GetHouseholdResponse);

message TransferOwnershipRequest {
  string household_id = 1; string actor_user_id = 2; string new_owner_member_id = 3;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status / reason |
|---|---|
| Actor not the owner | `PERMISSION_DENIED` `owner_only` |
| Target is a profile, anonymous, or not in the household | `INVALID_ARGUMENT` `invalid_new_owner` |

## Acceptance criteria

- [ ] After transfer, `GetHousehold` lists the new owner first with `is_owner: true`, and the database still has
      exactly one owner.
- [ ] Only the new owner can now remove members.
