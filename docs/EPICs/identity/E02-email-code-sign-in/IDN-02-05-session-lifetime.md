# IDN-02-05 Session lifetime and JWT issuance

| | |
|---|---|
| Service | identity |
| Epic | [IDN-02 Sign-in with e-mail code](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | IDN-02-02 |
| Unblocks | MOB-01, GW-02 |

## Goal

People stay signed in on their phone for weeks without typing codes, while the gateway only ever sees
short-lived tokens.

## Business rules

1. Better Auth session: valid 60 days (`session.expiresIn`), extended on use once a day (`session.updateAge`), so a
   phone used at least every 60 days never signs out on its own.
   Decided (Q033): 60 days, extended on daily use.
2. JWT: 15 minutes (`jwt.expirationTime: "15m"`, exists), issued only by `GET /api/auth/token` for a valid session.
   The gateway accepts nothing else.
3. Sign-out (`POST /api/auth/sign-out`) deletes the session. A JWT already issued stays valid until it expires
   (at most 15 minutes); this is accepted.
4. A session belongs to one device. Signing in on a second phone creates a second session; both work.
5. Sessions store `ip_address` and `user_agent` as Better Auth does; they are personal data, included in the export
   (IDN-09-03) and deleted with the account.
6. No "remember me" switch, no password, no session "freshness" requirement: account deletion asks for its own
   code instead (IDN-09-01).

## Contract

No new endpoint. Settings in `src/auth.ts`.

## Data

No schema change.

## Errors

| Case | HTTP |
|---|---|
| `/api/auth/token` without a valid session | 401 |

## Acceptance criteria

- [ ] A session used on day 1 and day 50 is still valid on day 100 (test clock); one unused for 61 days is not.
- [ ] A JWT from `/api/auth/token` has `exp − iat = 900`.
- [ ] After sign-out, `/api/auth/token` returns 401 for that session; a second device's session still works.
