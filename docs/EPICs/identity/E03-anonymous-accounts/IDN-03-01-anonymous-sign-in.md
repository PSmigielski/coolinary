# IDN-03-01 Anonymous sign-in

| | |
|---|---|
| Service | identity |
| Epic | [IDN-03 Anonymous accounts and upgrade](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | IDN-04-02, IDN-02-03 |
| Unblocks | IDN-03-02, MOB-01, GW-02 |

## Goal

A new user can use the Free features (one scan, one day plan, a fridge list) right after installing, without
giving an e-mail address (business-model §3 "Paywall").

## Business rules

1. Add Better Auth's `anonymous` plugin with `disableDeleteAnonymousUser: true` (identity deletes anonymous users
   itself, IDN-03-02, IDN-09-02). This also turns off the public `/delete-anonymous-user` route, so every deletion
   goes through the path that emits events.
2. `POST /api/auth/sign-in/anonymous` creates a user with `is_anonymous = true`, name `""`, and a placeholder e-mail
   `anon-<user id>@anonymous.coolinary.invalid` (`generateRandomEmail`). The placeholder is never shown, never
   e-mailed and never put in a JWT.
3. The user-create hook creates the personal household as for any new account (IDN-04-02): plan Free, one owner
   member with display name `""`. Events `household.created` and `household.member_added` are emitted.
4. Calling anonymous sign-in while already holding an anonymous session returns Better Auth's error (400); the app
   keeps the session it has.
5. Per-IP limit: 5 anonymous sign-ins per hour (IDN-02-03).
   Decided (Q034): 5 per hour per IP; the gateway's anonymous sign-up limit uses the same number (GW-03-04).
   Device attestation (App Attest / Play Integrity) may tighten this later.
6. What anonymous accounts may not do: create or accept invites (IDN-07-02, IDN-07-04), be an admin (IDN-10-01),
   receive any e-mail. Everything else follows the household's plan.
7. JWT claims for anonymous users: `anon: true`, no `email` claim (IDN-04-03).
8. Metric `identity_anonymous_sign_ins_total`.

## Contract

| Method and path | Answer |
|---|---|
| `POST /api/auth/sign-in/anonymous` (empty body) | `200` session; then `/api/auth/token` gives the JWT |

## Data

`bun run auth:schema` adds the plugin's user column (`is_anonymous`, boolean, default false), then
`bun run db:generate`. Add an index for the cleanup job:

```ts
index("user_anonymous_idx").on(user.isAnonymous).where(sql`is_anonymous`)
```

## Errors

| Case | HTTP |
|---|---|
| Already anonymous | 400 (Better Auth) |
| Over 5 per hour from one IP | 429 `rate_limited` |

## Acceptance criteria

- [ ] A fresh anonymous sign-in yields a JWT with `anon: true`, a `householdId`, a `memberId`, `plan: "free"`, and
      no `email`.
- [ ] `events.identity` receives `household.created` and `household.member_added` for that household.
- [ ] The sixth anonymous sign-in from one IP within an hour returns 429.
- [ ] `POST /api/auth/delete-anonymous-user` is not available.
