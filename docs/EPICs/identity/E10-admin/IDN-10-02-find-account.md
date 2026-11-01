# IDN-10-02 Find an account

| | |
|---|---|
| Service | identity |
| Epic | [IDN-10 Admin RPCs and admin role](README.md) |
| Release | R1 |
| Size | S |
| Depends on | IDN-10-01, IDN-05-01 |
| Unblocks | ADM-05 |

## Goal

When someone writes to support, the founder can find their account and household from the e-mail address.

## Business rules

1. `FindAccount` takes exactly one of: `email` (normalised as IDN-02-02 rule 7), `user_id`, `household_id`.
2. It returns: user id, e-mail, created at, anonymous flag, last seen (newest session `updated_at`), and the
   household summary: id, name, role of the user, plan, plan source and expiry, account count, profile count,
   pending invite count, `sharing_limited` (IDN-05-04).
3. For `household_id` it returns the household summary and its owner's account.
4. Nothing about health, fridge or plans (other services own that).
5. The call is logged with the admin id and the found user id (rule 6 of IDN-10-01); the searched e-mail isn't
   logged.

## Contract

```proto
message FindAccountRequest {
  oneof key { string email = 1; string user_id = 2; string household_id = 3; }
}
message FindAccountResponse {
  string user_id = 1;
  string email = 2;
  string created_at = 3;
  bool is_anonymous = 4;
  string last_seen_at = 5;
  HouseholdSummary household = 6;
}
message HouseholdSummary {
  string id = 1; string name = 2; Role role = 3; Plan plan = 4;
  string plan_source = 5; string plan_expires_at = 6;
  int32 account_count = 7; int32 profile_count = 8; int32 pending_invites = 9;
  bool sharing_limited = 10;   // Free with more than one account member (IDN-05-04)
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| None or more than one key | `INVALID_ARGUMENT` |
| Nothing found | `NOT_FOUND` |

## Acceptance criteria

- [ ] Searching `Ewa@Coolinary.test` finds the account stored as `ewa@coolinary.test`.
- [ ] Searching by household id returns its owner.
