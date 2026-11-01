# MOB-01-07 Sign-out and local wipe

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-01 Auth session](README.md) |
| Release | R1 |
| Size | M |
| Depends on | MOB-01-02, MOB-03-01, MOB-04-01, MOB-06-03, GW-11 (`DELETE /v1/devices/{installationId}`) |
| Unblocks | MOB-01-05, MOB-01-06, MOB-01-08 |

## Goal

After sign-out nothing of the previous user remains on the phone except device preferences, and the phone stops
receiving that user's notifications.

## Business rules

1. Sign-out is offered only to `signedIn` users; anonymous users have nothing to sign out to.
2. Order of steps:
   1. If the outbox has entries and the app is online: flush it, waiting up to 10 s. If entries remain, ask the user
      to confirm losing them.
   2. `DELETE /v1/devices/{installationId}` (MOB-06-03), best effort, 5 s timeout.
   3. `authClient.signOut()` (`POST /api/auth/sign-out`), best effort, 5 s timeout. Offline sign-out still wipes
      locally; the server session then expires by itself.
   4. Local wipe (rule 3).
   5. `status = "none"`; the next request starts a new anonymous account (MOB-01-03).
3. Local wipe removes: the session cookie (Better Auth's SecureStore entries), the in-memory JWT, the `useAuth` state,
   `coolinary.auth-hint`, the whole query cache (`queryClient.clear()` and the persisted copy
   `coolinary.query-cache`), the outbox `coolinary.outbox.v1`, the device-registration fingerprint, the live socket,
   and the health profile in SecureStore (`health.profile`). Decided (Q136): `health.profile` is always deleted on
   sign-out, so the next person on the phone never inherits someone's age, height and weight.
4. Kept: device preferences in `coolinary.settings` (energy unit, later language) and the installation id
   (`device.installation`).
5. The wipe is also used without step 2–3 by MOB-01-06 (session lost) and MOB-01-08 (account deleted).

## Contract

```ts
export function signOut(opts?: { force?: boolean }): Promise<void>;
export function wipeLocalData(reason: "sign_out" | "session_lost" | "account_deleted"): Promise<void>;
```

## Data

Removes the keys listed in rule 3.

## Errors

| Case | Behaviour |
|---|---|
| Server calls fail or time out | Ignored; the wipe always runs |

## Acceptance criteria

- [ ] After sign-out, AsyncStorage contains only `coolinary.settings`; SecureStore contains only
      `device.installation`.
- [ ] A push sent to the old user after sign-out does not arrive on the phone (device was unregistered).
- [ ] Sign-out in airplane mode still leaves the phone in the clean state of the first check.
