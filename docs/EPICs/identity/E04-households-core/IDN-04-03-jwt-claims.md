# IDN-04-03 JWT claims

| | |
|---|---|
| Service | identity |
| Epic | [IDN-04 Households and members core](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | IDN-04-01 |
| Unblocks | GW-02 (reads the claims), MOB-01 |

## Goal

The gateway can fill `household_id`, `member_id`, `x-user-id` and `x-plan` for every call from the token alone,
and knows whether the user is anonymous.

## Business rules

1. Claims (in `jwt({ definePayload })`): `sub` (user id), `householdId`, `memberId` (the user's member id in that
   household), `plan` (`free` / `plus`), `anon` (boolean). `email` only when the account is not anonymous.
2. No other personal data: no name, no household name, no role. Identity checks roles itself on each RPC.
3. No admin claim; admin rights are checked per request with `IsAdmin` (IDN-10-01).
4. Values are read when the token is issued. A change (plan, joining, leaving) shows in the next token, at most
   15 minutes later. RPCs that change the caller's household or plan return `token_stale: true` so the gateway can
   tell the app to fetch a new token right away (used in IDN-07-04, IDN-08-01, IDN-08-02, IDN-10-03).
5. Issuer and audience stay `BETTER_AUTH_URL` (the public API URL), as the gateway expects.
6. The gateway's `toUser` (`gateway/src/lib/auth.ts`) changes in the same commit (GW-02), per `.claude/rules/elysia.md`.
7. The encoded token stays under 1 KB.

## Contract

JWT payload example:

```json
{ "sub": "u_8f2…", "householdId": "6b1e…", "memberId": "a93c…", "plan": "free", "anon": false,
  "email": "ewa@coolinary.test", "iss": "https://api.example.com", "aud": "https://api.example.com",
  "iat": 1793000000, "exp": 1793000900 }
```

## Data

No schema change (`householdClaims` gains `memberId`).

## Errors

None.

## Acceptance criteria

- [ ] Tokens for anonymous users have `anon: true` and no `email`; for e-mail accounts `anon: false` and `email`.
- [ ] `memberId` matches the user's row in `household_members`.
- [ ] Gateway `GET /v1/me` returns the new fields (GW-02 test).
