# IDN-08-04 Cut off removed members' live tokens

| | |
|---|---|
| Service | identity |
| Epic | [IDN-08 Leaving and removing members](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | IDN-08-01, GW-02 (checks the key) |
| Unblocks | — |

## Goal

Someone who left or was removed can't keep reading the household's fridge and list with a token issued before,
for up to 15 minutes.

## Business rules

1. Whenever an account stops being a member of a household (left, removed, account deleted, joined
   another household), identity writes redis-state key `auth:revoked:<user_id>:<household_id>` = `1` with a TTL of
   900 s (the JWT lifetime), after the transaction commits.
2. The gateway checks this key on every authenticated request for the token's `sub` and `householdId`; when it
   exists, it answers 401 `token_stale` and the app fetches a new token (GW-02).
3. If writing the key fails (redis-state down), identity logs an error and carries on; the gap is then at most
   15 minutes, as without this task.
4. The key holds no personal data beyond the ids.

## Contract

Shared redis-state key `auth:revoked:<user_id>:<household_id>`, written by identity, read by the gateway.
Decided (Q191): access ends at once; the gateway reads this key on every request (one Redis read).

## Data

No schema change.

## Errors

None.

## Acceptance criteria

- [ ] A member removed at time T gets 401 `token_stale` on their next request at T + 5 s with their old token.
- [ ] The key is gone after 15 minutes.
