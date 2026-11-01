# IDN-02-04 Remove password sign-up

| | |
|---|---|
| Service | identity |
| Epic | [IDN-02 Sign-in with e-mail code](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | IDN-02-02 |
| Unblocks | GW-02, MOB-01 |

## Goal

There is one way in (an e-mail code), so there are no passwords to store, reset or leak (decision 12).

## Business rules

1. In `src/auth.ts` set `emailAndPassword: { enabled: false }` and remove `minPasswordLength`.
2. `POST /api/auth/sign-up/email` and `POST /api/auth/sign-in/email` then answer 404 (Better Auth drops the routes).
3. Existing accounts keep their e-mail address and sign in with a code to the same address; their user id and
   household stay the same.
4. Migration deletes every `account` row with `provider_id = 'credential'` (password hashes). Before launch these
   exist only in dev.
   `GATE:` if the prod database has users when this ships, confirm with the founder before running the migration.
5. The `password` column stays in Better Auth's generated schema (owned by Better Auth); it stays empty.
6. Docs and tests that use password sign-up switch to the code flow: the "Try the auth flow" section in the root
   `README.md` (sign-in code read from Mailpit's API) and any identity or gateway smoke test.

## Contract

Removes `POST /api/auth/sign-up/email` and `POST /api/auth/sign-in/email`.

## Data

`identity/drizzle/000NN_remove_credentials.sql`:

```sql
DELETE FROM account WHERE provider_id = 'credential';
```

## Errors

| Case | HTTP |
|---|---|
| Call to the removed routes | 404 |

## Acceptance criteria

- [ ] `POST /api/auth/sign-up/email` returns 404 through the gateway.
- [ ] A dev user created earlier with a password signs in with a code and gets the same user id and household id.
- [ ] `SELECT count(*) FROM account WHERE provider_id = 'credential'` is 0 after migration.
- [ ] The root README's auth example works as written against `make dev`.
