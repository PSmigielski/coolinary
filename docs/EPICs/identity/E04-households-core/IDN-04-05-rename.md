# IDN-04-05 Rename household and display names

| | |
|---|---|
| Service | identity |
| Epic | [IDN-04 Households and members core](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | IDN-04-04 |
| Unblocks | GW-04, IDN-06-01, IDN-07-02 |

## Goal

People can name their household and themselves, so invites say whose household it is and portions read
"Adam 1½, Ewa 1½" (features.md §2 "How one meal is split", step 7).

## Business rules

1. `UpdateHousehold(household_id, actor_user_id, name)`: only the owner may rename.
   Decided (Q024): owner only.
2. `UpdateMember(household_id, actor_user_id, member_id, display_name)`:
   - an account member may rename only their own member;
   - any account member may rename a profile without an account (IDN-06);
   - nobody renames another account member.
3. Text rules for both: trim, collapse inner whitespace to one space, no control characters. Household name 1–40
   characters, display name 1–30 characters (counted in Unicode code points). Emoji allowed.
4. Display names are unique within a household, compared case-insensitively after the rules above, so portions
   are never "Ola 1, Ola 2". A clash returns `ALREADY_EXISTS` `name_taken`.
5. Names are personal data: never in events, logs or metric labels.
6. No events (other services don't store names).

## Contract

```proto
rpc UpdateHousehold(UpdateHouseholdRequest) returns (GetHouseholdResponse);
rpc UpdateMember(UpdateMemberRequest) returns (GetHouseholdResponse);

message UpdateHouseholdRequest { string household_id = 1; string actor_user_id = 2; string name = 3; }
message UpdateMemberRequest {
  string household_id = 1; string actor_user_id = 2; string member_id = 3; string display_name = 4;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status / reason |
|---|---|
| Not the owner (household rename) | `PERMISSION_DENIED` `owner_only` |
| Renaming another account member | `PERMISSION_DENIED` `not_your_member` |
| Empty or too long after trimming | `INVALID_ARGUMENT` `invalid_name` |
| Display name used by another member | `ALREADY_EXISTS` `name_taken` |

## Acceptance criteria

- [ ] `"  Kasia   W "` is stored as `"Kasia W"`.
- [ ] A 41-character household name is refused; 40 is accepted.
- [ ] Renaming a profile to "ola" when "Ola" exists returns `ALREADY_EXISTS`.
- [ ] A member renaming the household gets `PERMISSION_DENIED`.
