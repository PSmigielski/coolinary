# GW-02-01 JWT claims: member id, anonymous flag, registered-only routes

| | |
|---|---|
| Service | gateway |
| Epic | [GW-02 Auth and account routes](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-01-02, IDN-03 (anonymous claim), IDN-04 (memberId claim) |
| Unblocks | GW-02-03 … GW-02-06, every route that needs `member_id` or refuses anonymous users |

## Goal

The gateway reads every claim it needs from the JWT (no call to identity per request) and can refuse anonymous
accounts on the few routes that need a real one.

## Business rules

1. `verifyBearer` (exists, `src/lib/auth.ts`) keeps checking signature (JWKS, cached), `iss`, `aud`, `exp`. It also
   requires `nbf`/`iat` not more than 60 s in the future (clock skew).
2. `AuthUser` gains:
   - `memberId` — claim `memberId` (identity's member id of this user in their household). Missing → the token is
     treated as invalid (401), because every R1 token from identity carries it once IDN-04 ships.
   - `isAnonymous` — claim `isAnonymous === true`; anything else is `false`.
   - `issuedAt` — from `iat`, used by GW-02-05.
3. Route macro `auth` takes `true` (any account, anonymous included) or `"registered"` (non-anonymous only).
   Anonymous on a `"registered"` route → 403 `registration_required`.
4. Registered-only routes: invite create/list/cancel/preview/accept (GW-04-03, GW-04-04), adding profiles
   (GW-04-02), GDPR export (GW-02-06). Everything else accepts anonymous users, including AI routes within the Free
   quota (business model: "first scan and first full day-plan are free with no account").
5. `plan` claim: `plus` only when the claim is exactly `plus`; anything else is `free` (exists).
6. A JWT whose `householdId` is missing returns 403 `no_household` on any route that needs a household (today
   `/v1/household` returns 404; changed to the problem model).

## Contract

```ts
interface AuthUser {
  id: string; email?: string; householdId?: string; memberId: string;
  plan: "free" | "plus"; isAnonymous: boolean; issuedAt: number;
}
// macro: { auth: true | "registered" }
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| No or invalid token, missing `memberId` | 401 | `unauthorized` |
| Anonymous on a registered-only route | 403 | `registration_required` |
| No household in token | 403 | `no_household` |

## Acceptance criteria

- [ ] A token with `isAnonymous: true` passes `auth: true` and gets 403 `registration_required` on `auth: "registered"`.
- [ ] A token without `memberId` gets 401.
- [ ] A token with `iat` 5 minutes in the future gets 401; 30 s in the future passes.
- [ ] `plan: "PLUS"` (wrong case) is treated as `free`.

## Notes

- Needs identity to add `memberId` and `isAnonymous` to the JWT payload (`householdClaims` in identity); listed as a
  cross-service need in the gateway README.
