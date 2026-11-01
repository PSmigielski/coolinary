# IDN-08-01 Move-out function and leaving

| | |
|---|---|
| Service | identity |
| Epic | [IDN-08 Leaving and removing members](README.md) |
| Release | R2 |
| Size | M |
| Depends on | IDN-04-02, IDN-04-04, IDN-01-01 |
| Unblocks | IDN-08-02, IDN-05-04, GW-04, KIT-13, PLN-21 |

## Goal

Anyone can leave a shared household and start again with their own empty fridge, without taking anything from the
people who stay.

## Business rules

1. `moveOut(tx, householdId, memberId, reason)` is the one function that takes an account member out of a household
   (used by leaving, removal, plan downgrade and account deletion paths). It:
   1. deletes the member row in the old household;
   2. creates a new personal household for the user with `createPersonalHousehold` (IDN-04-02), Free, the same
      display name, `previous_member_id` = the old member id (so kitchen can carry over the person's own diet and
      target, but not the fridge);
   3. emits `household.member_removed` (old household, old member id, `reason`, `new_household_id`) before the
      creation events;
   4. is where IDN-08-04 adds the token revocation key;
   5. leaves invites the member created valid (IDN-07-05 rule 4).
2. `LeaveHousehold(household_id, actor_user_id, new_owner_member_id)`:
   - A non-owner account member leaves with `moveOut(…, "left")`.
   - The owner must name a new owner (an account member who isn't anonymous). Identity hands over ownership
     (IDN-08-03) and then moves the old owner out, in one transaction.
     Decided (Q186): the owner names the new owner; there is no automatic hand-over.
   - An owner with no other account member can't leave (`already_alone`); they can delete profiles or their account.
3. Decided (Q185): the household's plan stays with the household when the owner leaves, until the payment route is
   chosen (IDN-05-05).
4. Profiles stay with the household.
5. The leaver never takes stock along; leaving always starts with an empty fridge.
6. The answer carries the new household and `token_stale: true`.

## Contract

```proto
rpc LeaveHousehold(LeaveHouseholdRequest) returns (LeaveHouseholdResponse);

message LeaveHouseholdRequest {
  string household_id = 1;
  string actor_user_id = 2;
  string new_owner_member_id = 3;   // required when the actor is the owner
}
message LeaveHouseholdResponse { Household new_household = 1; bool token_stale = 2; }
```

## Data

No schema change.

## Events

`household.member_removed` (reason `left`), then `household.created` and `household.member_added` for the new
household.

## Errors

| Case | gRPC status / reason |
|---|---|
| Owner without `new_owner_member_id` | `FAILED_PRECONDITION` `owner_must_transfer` |
| Owner alone | `FAILED_PRECONDITION` `already_alone` |
| New owner not an account member, or anonymous | `INVALID_ARGUMENT` `invalid_new_owner` |
| Actor not a member | `PERMISSION_DENIED` |

## Acceptance criteria

- [ ] A member who leaves owns a new household with no profiles; the old household still lists everyone else.
- [ ] Events arrive in the order: `member_removed`, `household.created`, `member_added`.
- [ ] An owner leaving with `new_owner_member_id` makes that member owner and then moves out.
