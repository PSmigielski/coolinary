# GW-04-05 Leave the household, remove a member

| | |
|---|---|
| Service | gateway |
| Epic | [GW-04 Household, profiles and invite routes](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | GW-01-03, IDN-08 (`LeaveHousehold`, `RemoveMember`, `TransferOwnership`, revocation key IDN-08-04) |
| Unblocks | MOB-01, MOB-03 |

## Goal

Anyone can leave and start again with their own empty fridge; the owner can remove people (features.md §2,
verification 15).

## Business rules

1. **Leave** (`POST /v1/household/leave`, body `{ "newOwnerMemberId"? }`): calls identity
   `LeaveHousehold(household_id, actor_user_id, new_owner_member_id)`. Identity creates the person's new personal
   household (empty fridge) and publishes the events. Response: new `householdId`, `memberId`,
   `tokenRefreshRequired: true` (identity's `token_stale`). An owner of a household with other accounts must name
   the new owner; without it identity refuses (`FAILED_PRECONDITION`, its `x-error-code`, → 422).
2. **Remove** (`DELETE /v1/household/members/{memberId}`): owner only (identity checks) → identity
   `RemoveMember(household_id, actor_user_id, member_id)`. Target must be a member with an account (profiles use
   GW-04-02). Response 204.
3. **Hand over ownership** (`POST /v1/household/owner`, body `{ "memberId" }`): owner only → identity
   `TransferOwnership(household_id, actor_user_id, new_owner_member_id)`. Response 200 household (GW-04-01 shape).
4. Both drop `gw:ent:*` for the affected households.
5. Idempotency required on leave (a retry must not create a second new household); optional on remove.
6. Access after removal: identity writes `auth:revoked:<user_id>:<household_id>` (TTL 900 s) in redis-state
   whenever an account stops being a member (IDN-08-04). The gateway checks it on every authenticated request for
   the token's `sub` and `householdId` and answers 401 `token_stale`; the app fetches a new token. If redis-state
   is down the check is skipped (the gap is then at most the 15-minute JWT lifetime).

## Contract

| Method | Path | Body | 2xx |
|---|---|---|---|
| POST | `/v1/household/leave` | `{ "newOwnerMemberId"? }` | 200 `{ "householdId", "memberId", "tokenRefreshRequired": true }` |
| POST | `/v1/household/owner` | `{ "memberId" }` | 200 household |
| DELETE | `/v1/household/members/{memberId}` | — | 204 |

Tag `household`, `auth: "registered"` (anonymous accounts are never in a shared household).

## Data

Reads `auth:revoked:<user_id>:<household_id>` in redis-state (written by identity, IDN-08-04).

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Owner with other accounts leaves without `newOwnerMemberId` | 422 | identity's code |
| Hand-over by a non-owner, or to a profile | 403 / 422 | `forbidden` / identity's code |
| Token of a removed member | 401 | `token_stale` |
| Remove by a non-owner | 403 | `forbidden` |
| Remove self, a profile, or a non-member | 422 / 404 | `use_leave` / `not_found` |

## Acceptance criteria

- [ ] After leave and a token refresh, `GET /v1/stock` returns an empty fridge in the new household.
- [ ] After the owner removes account B, B's next call with a refreshed token no longer sees the old fridge or list
      (verification 15).
- [ ] A retried leave with the same idempotency key returns the same new household id.
- [ ] The owner of a two-account household leaving without `newOwnerMemberId` gets 422; with it, the other account
      becomes owner and the leaver gets a new household.
- [ ] Right after removal, B's old token gets 401 `token_stale` (revocation key present), before the JWT expires.

## Notes

- Decided (Q191): access ends at once. The revocation marker `auth:revoked:<user_id>:<household_id>` is identity's
  (IDN-08-04); the gateway reads it on every authenticated request (one Redis read). Decided (Q186): an owner names the new owner before leaving; there is no automatic hand-over (IDN-08-01).
