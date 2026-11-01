# MOB-01-04 JWT provider and refresh

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-01 Auth session](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | MOB-01-01, MOB-01-02, GW-02 (`GET /api/auth/token`) |
| Unblocks | MOB-01-06, MOB-02-01, MOB-05-01 |

## Goal

Every `/v1/*` request carries a valid 15-minute JWT without the user ever seeing a token expire, and the app reacts
when the token's household or plan changes.

## Business rules

1. The JWT is fetched with `authClient.token()` (`GET /api/auth/token`, authenticated by the session cookie) and kept
   **in memory only**. It is never written to SecureStore, AsyncStorage or logs. On cold start a new one is fetched.
2. The payload is decoded on the phone (base64url, no signature check) into:

   ```ts
   interface JwtClaims { sub: string; householdId?: string; plan: "free" | "plus"; exp: number; iat: number }
   ```
3. Expiry is computed against the phone's own clock at receipt, to survive clock skew:
   `expiresAt = receivedAt + (exp − iat) × 1000`. If `iat` is missing, `expiresAt = exp × 1000`.
4. **Refresh margin: 120 s.** `getToken()` returns the cached token while `now < expiresAt − 120 s`; otherwise it
   refreshes first.
5. Refresh is single-flight: all callers during a refresh await the same promise.
6. Proactive refresh: while the app is in the foreground, a timer fires at `expiresAt − 120 s`. When the app returns
   to the foreground (`AppState` → `active`), the provider checks the margin immediately. No timers run in the
   background.
7. Refresh failure:
   - network error or 5xx: keep using the current token until `expiresAt`; retry at 5 s, 15 s, 30 s, then every
     60 s while online. After `expiresAt`, `getToken()` waits up to 10 s for a refresh, then fails with
     `OfflineError` or `ServerDownError` (MOB-02-02).
   - 401: the session is gone → MOB-01-06.
8. After every new token, compare claims with `useAuth`: a changed `householdId` fires `householdChanged`, a changed
   `sub` fires `userChanged` (MOB-01-02), a changed `plan` invalidates entitlements and Plus-gated queries
   (MOB-02-05).
9. `forceRefresh()` exists for flows that change claims on the server: invite accepted, household left, plan changed,
   anonymous upgrade. It bypasses the margin but is still single-flight.

## Contract

```ts
export interface TokenProvider {
  getToken(): Promise<string>;          // valid for at least 120 s, or throws
  forceRefresh(): Promise<string>;
  current(): JwtClaims | null;
}
```

HTTP: `GET /api/auth/token` → `200 { "token": "<jwt>" }`, `401` when there is no session.

## Data

Memory only.

## Errors

| Case | Behaviour |
|---|---|
| Malformed token (can't decode) | Treat as refresh failure; report once to the dev log without the token |
| `householdId` missing in claims | Allowed only briefly after anonymous sign-up; retry the token once after 1 s, then `ServerDownError` |

## Acceptance criteria

- [ ] With a device clock 10 minutes wrong, tokens are still refreshed 120 s before their real expiry.
- [ ] 20 concurrent `getToken()` calls during expiry produce one `/api/auth/token` call.
- [ ] App backgrounded for 30 minutes, then foregrounded: the first request after resume uses a fresh token and
      gets no 401.
- [ ] A grep of the AsyncStorage dump and SecureStore keys finds no JWT.
