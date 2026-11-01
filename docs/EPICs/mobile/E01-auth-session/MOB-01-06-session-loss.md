# MOB-01-06 401 recovery and session loss

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-01 Auth session](README.md) |
| Release | R1 |
| Size | S |
| Depends on | MOB-01-04, MOB-02-01 |
| Unblocks | MOB-04-02 |

## Goal

An expired token never shows up as an error, and a lost session is handled the same way everywhere.

## Business rules

1. A `401` from any `/v1/*` call triggers one `forceRefresh()` and one retry of that request with the new token
   (same `Idempotency-Key` and `x-request-id`). A second `401` is not retried.
2. While a refresh is running, other `401`s wait for it instead of starting their own.
3. If `GET /api/auth/token` itself returns `401`, the session is gone (expired, revoked, account deleted elsewhere):
   - `status === "signedIn"` → local wipe as in sign-out (MOB-01-07, without the server call) and the app shows the
     signed-out state with the message "Please sign in again". Pending outbox entries are kept for 24 h, bound to
     the old `userId`, and replayed only if the same user signs in again; otherwise dropped.
   - `status === "anonymous"` → the anonymous data can't be recovered. The app wipes locally and starts a new
     anonymous account on the next request (MOB-01-03). Pending outbox entries are dropped and the user is told
     that unsent changes were lost.
4. A `403` with error code `not_member` (the user was removed from the household, IDN-08 / GW-04) triggers
   `forceRefresh()`; the new JWT's household replaces the old one through `householdChanged`.
5. A `401` from the live socket (close code 4401) follows rule 1 for the socket (MOB-05-01).

## Contract

Uses `GET /api/auth/token` and the error codes from MOB-02-02 (`unauthorized`, `not_member`).

## Data

No new storage.

## Errors

| Case | Behaviour |
|---|---|
| Retry after refresh fails with another status | Normal error handling for that status (MOB-02-02) |

## Acceptance criteria

- [ ] With a token revoked on the server, the next query retries once and succeeds without user action.
- [ ] Deleting the session in identity's database while signed in leads to the signed-out state on the next call.
- [ ] Deleting an anonymous session leads to a new anonymous household, and the user sees the data-loss message.
- [ ] Removing a member in identity moves that phone to its new household within one request (verification 15).
