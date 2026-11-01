# IDN-06-02 Delete a profile

| | |
|---|---|
| Service | identity |
| Epic | [IDN-06 Profiles without accounts](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | IDN-06-01 |
| Unblocks | GW-04, KIT-13 (deletes the person's diet and target) |

## Goal

A profile that is no longer needed is removed, and its diet and calorie target are erased everywhere.

## Business rules

1. `DeleteProfile(household_id, actor_user_id, member_id)`: any account member may delete a profile of their
   household. Works on Free and Plus.
2. Only rows with `role = profile` can be deleted this way; for account members see IDN-08.
3. The member row is deleted and `household.member_removed` is emitted with reason `profile_deleted`,
   `user_id: null`, `new_household_id: null`.
4. Deleting an already deleted profile returns `NOT_FOUND`.

## Contract

```proto
rpc DeleteProfile(DeleteProfileRequest) returns (DeleteProfileResponse);

message DeleteProfileRequest { string household_id = 1; string actor_user_id = 2; string member_id = 3; }
message DeleteProfileResponse {}
```

## Data

No schema change.

## Events

Publishes `household.member_removed` (reason `profile_deleted`).

## Errors

| Case | gRPC status / reason |
|---|---|
| Member is an account | `FAILED_PRECONDITION` `not_a_profile` |
| Unknown member | `NOT_FOUND` |
| Actor not an account member | `PERMISSION_DENIED` |

## Acceptance criteria

- [ ] Deleting a profile removes it from `GetHousehold` and emits one event.
- [ ] Trying to delete the owner through this RPC returns `not_a_profile`.
