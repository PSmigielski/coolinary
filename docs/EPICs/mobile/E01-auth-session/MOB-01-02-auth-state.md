# MOB-01-02 Auth state and app start

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-01 Auth session](README.md) |
| Release | R1 |
| Size | M |
| Depends on | MOB-01-01 |
| Unblocks | MOB-01-03, MOB-01-04, MOB-01-05, MOB-01-07, MOB-03-01 |

## Goal

The app always knows, in one place, whether there is no session, an anonymous session or a signed-in session, and
which user and household it is acting for. Everything that caches or queues data keys it by that household.

## Business rules

1. State lives in a non-persisted Zustand store `useAuth` (`src/stores/auth.ts`):

   ```ts
   type AuthStatus = "booting" | "none" | "anonymous" | "signedIn";
   interface AuthState {
     status: AuthStatus;
     userId: string | null;
     householdId: string | null;   // from the JWT claim `householdId`
     plan: "free" | "plus";        // from the JWT claim `plan`
     email: string | null;         // memory only, never persisted
   }
   ```
2. Boot order at app start (before any `/v1/*` query runs):
   1. `status = "booting"`; queries that need auth are disabled while booting (`enabled: status !== "booting"`).
   2. Read the session with `authClient.getSession()`. Offline: use the cached session from SecureStore that the
      Expo client keeps; if there is none, `status = "none"`.
   3. Session present → `status = "anonymous"` when the user is anonymous (Better Auth's anonymous flag on the user),
      otherwise `"signedIn"`; then fetch a JWT (MOB-01-04) when online.
   4. Boot must finish within 3 s; after 3 s without a server answer the app continues with the cached session and
      the last auth hint (rule 4), and retries in the background.
3. `userId`, `householdId` and `plan` are set only from a decoded JWT (MOB-01-04), never from client guesses. Until a
   JWT exists, they come from the auth hint (rule 4) so the persisted cache of the last household can render offline.
4. An auth hint is persisted in AsyncStorage key `coolinary.auth-hint`:
   `{ v: 1, userId, householdId, status: "anonymous" | "signedIn" }`. It contains ids only, no e-mail, no token.
   It is rewritten whenever a JWT changes one of these values and deleted on sign-out.
5. **Household switch.** When a new JWT carries a `householdId` different from the current one (invite accepted,
   member left or was removed, anonymous upgrade into an existing account), the store emits `householdChanged(old,
   new)`. Subscribers: query cache (MOB-03-01 removes every key under the old household), outbox (MOB-04-02 flushes
   or drops entries for the old household), live socket (MOB-05-01 reconnects), push registration (MOB-06-03).
6. **User switch.** When `userId` changes, the store emits `userChanged(old, new)`; the push registration is
   redone (MOB-06-03) and user-scoped queries (`me`, notification preferences) are removed.
7. Expose hooks `useAuthStatus()`, `useHouseholdId()`, `usePlan()`; components never read SecureStore or the JWT.

## Contract

`GET /api/auth/get-session` via `authClient.getSession()` (GW-02 proxy).

## Data

- Memory: `useAuth` store.
- AsyncStorage: `coolinary.auth-hint` (ids only).

## Errors

| Case | Behaviour |
|---|---|
| `getSession` network error at boot | Continue with cached session and hint; `status` from the hint; retry on reconnect (MOB-02-04) |
| Session cookie present but server says no session (expired, revoked) | Handled by MOB-01-06 |
| Hint exists but no session cookie (e.g. Android reinstall with backup) | Delete the hint, `status = "none"` |

## Acceptance criteria

- [ ] Cold start online with an anonymous session: status goes `booting → anonymous` and `householdId` matches
      `GET /v1/me`.
- [ ] Cold start in airplane mode with a session: status leaves `booting` within 3 s and the last fridge list renders
      from the persisted cache.
- [ ] Changing `householdId` in a test JWT fires `householdChanged` exactly once and the old household's queries are
      gone from the cache.
- [ ] `coolinary.auth-hint` never contains an e-mail or a token (unit test on the serializer).
