# IDN-03 Anonymous accounts and upgrade

Decision 11: on first start the app creates an anonymous account (Better Auth `anonymous` plugin) with its own
household, so the first scan and the first day plan work without an e-mail. Signing in with an e-mail code later
upgrades the account and keeps everything the person did.

## Flow

1. First start: app → `/api/auth/sign-in/anonymous` → anonymous user + personal household (Free), JWT with
   `anon: true`.
2. Later, with the anonymous session still active, the app runs the e-mail code flow (IDN-02).
   - E-mail has no account yet → the anonymous household becomes the new account's household (IDN-03-02).
   - E-mail already has an account → the anonymous fridge moves into that account's household, the anonymous
     household is deleted (IDN-03-03).
3. Anonymous accounts unused for 30 days are deleted (IDN-03-04).

Better Auth calls `onLinkAccount({ anonymousUser, newUser })` in step 2. Identity sets
`disableDeleteAnonymousUser: true` and deletes the anonymous user itself, after moving the household, because
`households.owner_user_id` cascades on user delete.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-03-01 | Anonymous sign-in | R1 | M | IDN-04-02, IDN-02-03 |
| IDN-03-02 | Upgrade to a new e-mail account | R1 | M | IDN-03-01, IDN-02-02, IDN-01-01 |
| IDN-03-03 | Upgrade into an existing account | R1 | M | IDN-03-02 |
| IDN-03-04 | Delete stale anonymous accounts | R1 | M | IDN-03-01, IDN-09-02 |
