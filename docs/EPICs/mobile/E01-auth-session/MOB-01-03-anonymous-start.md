# MOB-01-03 Anonymous start

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-01 Auth session](README.md) |
| Release | R1 |
| Size | S |
| Depends on | MOB-01-02, IDN-03 (anonymous accounts), GW-02 (`POST /api/auth/sign-in/anonymous`) |
| Unblocks | MOB-01-04, MOB-04-02 |

## Goal

A new user can scan, stock the fridge and get a first day plan without giving an e-mail address (business model §3
"first scan and first full day-plan are free with no account"). The server sees an anonymous Free household.

## Business rules

1. The anonymous account is created **lazily**: the first time a request needs auth and `status === "none"`.
   Browsing static content (presets, the system status screen) needs no account. Decided (Q043): lazy creation, so
   installs that never use the app leave no server rows.
2. Creation is single-flight: concurrent callers await the same `authClient.signIn.anonymous()` promise.
3. After success: fetch the JWT (MOB-01-04), set `status = "anonymous"`, write the auth hint. identity creates the
   household (IDN-03); the app learns `householdId` from the JWT.
4. Offline with `status === "none"`: do not create the account. Mutations that may run offline are queued without a
   household (MOB-04-01 stores `householdId: null`) and the account is created first when the queue replays.
5. Calling `signIn.anonymous()` while an anonymous session exists returns an error from Better Auth; never call it
   when `status` is `anonymous` or `signedIn`.
6. Failure handling: retry with backoff 1 s, 3 s, 10 s, then stop and surface `ServerDownError` (MOB-02-02) to the
   caller; the next user action tries again.
7. Anonymous households have the Free limits (README "Free and Plus"); nothing in the app treats them differently
   except the "save your data" prompt that a later frontend task shows.

## Contract

`POST /api/auth/sign-in/anonymous` → session cookie (stored by `expoClient`), then `GET /api/auth/token`.

```ts
export function ensureSession(): Promise<void>; // resolves when a session and JWT exist
```

## Data

No new storage; session cookie in SecureStore via MOB-01-01.

## Errors

| Case | Behaviour |
|---|---|
| 429 from Better Auth rate limit | Wait `Retry-After` (default 30 s), then one retry |
| 5xx / network | Rule 6 |

## Acceptance criteria

- [ ] Fresh install, online: the first `GET /v1/stock` triggers exactly one anonymous sign-up, then succeeds.
- [ ] Ten parallel first requests produce one sign-up call (single-flight test).
- [ ] Fresh install, offline: adding two fridge items queues them; on reconnect the account is created first and both
      items appear in `GET /v1/stock`.
- [ ] After app restart, no second anonymous account is created.

## Notes

- Decided (Q044): identity deletes an anonymous account and its data after **30 days without use** (IDN-03-04).
  The app then gets a 401 on its old session and MOB-01-06 starts a new anonymous account with an empty fridge. Signing in
  with an e-mail code keeps the data.
