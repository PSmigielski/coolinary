# IDN-04-04 Membership checks and richer `GetHousehold`

| | |
|---|---|
| Service | identity |
| Epic | [IDN-04 Households and members core](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | IDN-04-01 |
| Unblocks | GW-04, IDN-04-05, IDN-06, IDN-07, IDN-08 |

## Goal

Identity never shows or changes a household for someone who isn't in it, even if the gateway passes a stale
household id from a token issued before they left.

## Business rules

1. Every user-facing identity RPC carries `household_id` and `actor_user_id`, both filled by the gateway from the
   JWT (derived decision D3). Identity checks that `actor_user_id` is an account member of `household_id`; if not,
   `PERMISSION_DENIED` `not_a_member`.
2. `GetHousehold` adds `actor_user_id` to its request. Its response adds, per member: `joined_at`, `is_owner`; per
   household: `created_at`, `account_count`, `max_accounts` (5), `profile_count`.
3. Members are ordered: owner first, then account members by `joined_at`, then profiles by `joined_at`.
4. No member's e-mail address is ever returned to other members; display names only.
5. Display names that are empty are returned empty; the app shows a translated placeholder ("Me" for the caller,
   "Member" for others).
6. Service-to-service callers other than the gateway don't call these RPCs (they use events).

## Contract

`contracts/coolinary/identity/v1/identity.proto`:

```proto
message GetHouseholdRequest {
  string household_id = 1;
  string actor_user_id = 2;
}
message Household {
  string id = 1;
  string name = 2;
  string owner_user_id = 3;
  repeated Member members = 4;
  Plan plan = 5;
  string created_at = 6;       // RFC 3339
  int32 account_count = 7;
  int32 max_accounts = 8;      // 5
  int32 profile_count = 9;
}
message Member {
  string id = 1;
  string display_name = 2;
  Role role = 3;
  bool has_account = 4;
  string joined_at = 5;
  bool is_owner = 6;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Actor not a member of the household | `PERMISSION_DENIED` `not_a_member` |
| `household_id` not a UUID | `INVALID_ARGUMENT` |
| Household doesn't exist | `NOT_FOUND` |

## Acceptance criteria

- [ ] Calling `GetHousehold` with a household the actor left returns `PERMISSION_DENIED`.
- [ ] The response lists the owner first and includes `account_count`, `profile_count` and `max_accounts = 5`.
- [ ] No response field contains an e-mail address.
