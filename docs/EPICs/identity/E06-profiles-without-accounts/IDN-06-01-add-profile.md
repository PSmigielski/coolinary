# IDN-06-01 Add a profile

| | |
|---|---|
| Service | identity |
| Epic | [IDN-06 Profiles without accounts](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | IDN-04-04, IDN-04-05, IDN-05-01, IDN-01-01 |
| Unblocks | GW-04, KIT-11 (per-person profile rows), PLN-12 (portions) |

## Goal

A parent can add "Ola" and "Kuba" to the household, so each gets a portion of the shared dish.

## Business rules

1. `CreateProfile(household_id, actor_user_id, display_name)`: the actor is any account member of the household
   (owner or member).
2. The household must be Plus: portions and goals for every member, children included, are a Plus feature
   (features.md §6). Free → `FAILED_PRECONDITION` `plus_required`.
3. At most 10 profiles per household. Profiles don't count toward the 5 accounts.
   Decided (Q180): at most 10.
4. Display name rules and uniqueness as IDN-04-05 rules 3–4.
5. The profile is a `household_members` row with `role = profile`, `user_id = null`.
6. Emit `household.member_added` with `has_account: false`, `role: profile`, `brought_fridge_from: null`,
   `previous_member_id: null`.
7. A household that drops to Free keeps its profiles (IDN-05-04 rule 6); new ones can't be added until Plus again.

## Contract

```proto
rpc CreateProfile(CreateProfileRequest) returns (CreateProfileResponse);

message CreateProfileRequest { string household_id = 1; string actor_user_id = 2; string display_name = 3; }
message CreateProfileResponse { Member member = 1; }
```

## Data

No schema change (`role = 'profile'` already exists).

## Events

Publishes `household.member_added`.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status / reason |
|---|---|
| Household is Free | `FAILED_PRECONDITION` `plus_required` |
| 10 profiles already | `FAILED_PRECONDITION` `profile_limit` |
| Name invalid or taken | `INVALID_ARGUMENT` `invalid_name` / `ALREADY_EXISTS` `name_taken` |
| Actor not an account member | `PERMISSION_DENIED` |

## Acceptance criteria

- [ ] A Plus household with 5 accounts can still add profiles (they don't count).
- [ ] The eleventh profile is refused.
- [ ] A Free household gets `plus_required`.
- [ ] The event carries `has_account: false` and no name.
