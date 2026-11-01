# IDN-08-02 Owner removes a member

| | |
|---|---|
| Service | identity |
| Epic | [IDN-08 Leaving and removing members](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | IDN-08-01 |
| Unblocks | GW-04 |

## Goal

The owner can take someone out of the household; afterwards that person no longer sees the shared fridge or list
(features.md Verification 15).

## Business rules

1. `RemoveMember(household_id, actor_user_id, member_id)`: only the owner.
2. The target must be another account member. Profiles are deleted with `DeleteProfile` (IDN-06-02); the owner
   can't remove themself (they leave, IDN-08-01).
3. The target is moved out with `moveOut(…, "removed")`: they get a new empty personal household.
4. The removed person gets no message from identity; the app shows "You are no longer in <household>" when its next
   call fails with `not_a_member` or `token_stale` (IDN-08-04).
   Decided (Q189): no push.
5. The answer carries the household as it is now.

## Contract

```proto
rpc RemoveMember(RemoveMemberRequest) returns (GetHouseholdResponse);

message RemoveMemberRequest { string household_id = 1; string actor_user_id = 2; string member_id = 3; }
```

## Data

No schema change.

## Events

As IDN-08-01, with reason `removed`.

## Errors

| Case | gRPC status / reason |
|---|---|
| Actor is not the owner | `PERMISSION_DENIED` `owner_only` |
| Target is a profile | `FAILED_PRECONDITION` `use_delete_profile` |
| Target is the owner | `FAILED_PRECONDITION` `cannot_remove_owner` |
| Unknown member | `NOT_FOUND` |

## Acceptance criteria

- [ ] After removal the removed account's next token has a different `householdId`, and its calls with the old
      token are refused within seconds (IDN-08-04).
- [ ] A member (not owner) trying to remove someone gets `owner_only`.
