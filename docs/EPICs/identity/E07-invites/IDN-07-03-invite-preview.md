# IDN-07-03 Invite preview

| | |
|---|---|
| Service | identity |
| Epic | [IDN-07 Invites](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | IDN-07-02 |
| Unblocks | GW-04, IDN-07-04 |

## Goal

Before joining, the person sees whose household it is and what will happen to their own fridge, and the app knows
whether joining is possible at all.

## Business rules

1. `GetInvitePreview(code, actor_user_id)`, after normalising the code (IDN-07-01 rule 3).
2. Unknown and cancelled codes return `NOT_FOUND` `invite_not_found` (the same answer for both). Expired codes
   return `FAILED_PRECONDITION` `invite_expired`; accepted ones `FAILED_PRECONDITION` `invite_used`.
3. For a usable invite it returns: household name, owner display name, inviter display name, account count,
   `expires_at`.
4. It also returns, for this actor:
   - `can_join` and, if false, `reason`: `sign_in_required` (anonymous), `already_member`,
     `owns_shared_household` (they own a household with other account members and must hand it over first),
     `household_full`, `plus_required` (the inviting household lost Plus).
   - `can_bring_fridge`: true when the actor is the only account member of their current household.
   - `current_profiles_count`: profiles in the actor's current personal household, which will be deleted on joining.
   - `current_plan`: the actor's current household plan, so the app can warn that a manual Plus is lost.
5. Guessing protection: an actor may have at most 10 lookups of unknown codes per hour; then every preview and
   accept call from that actor returns `RESOURCE_EXHAUSTED` for the rest of the hour. The gateway rate-limits per
   IP in addition (GW-03).
6. The preview never shows e-mail addresses or the member list.

## Contract

```proto
rpc GetInvitePreview(GetInvitePreviewRequest) returns (InvitePreview);

message GetInvitePreviewRequest { string code = 1; string actor_user_id = 2; }
message InvitePreview {
  string household_name = 1;
  string owner_display_name = 2;
  string inviter_display_name = 3;
  int32 account_count = 4;
  string expires_at = 5;
  bool can_join = 6;
  string reason = 7;
  bool can_bring_fridge = 8;
  int32 current_profiles_count = 9;
  Plan current_plan = 10;
}
```

## Data

No schema change. redis-state key `identity:invite_miss:<user id>` (TTL 1 h) counts unknown-code lookups.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status / reason |
|---|---|
| Unknown or cancelled | `NOT_FOUND` `invite_not_found` |
| Expired / accepted | `FAILED_PRECONDITION` `invite_expired` / `invite_used` |
| Too many unknown codes | `RESOURCE_EXHAUSTED` `too_many_attempts` |
| Malformed code | `INVALID_ARGUMENT` `invalid_code` |

## Acceptance criteria

- [ ] An anonymous actor gets `can_join: false`, `reason: sign_in_required`.
- [ ] An actor alone in their household gets `can_bring_fridge: true`; a member of a shared household gets false.
- [ ] The eleventh unknown code in an hour returns `RESOURCE_EXHAUSTED`, and so does a valid code right after.
