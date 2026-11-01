# ADM-02 Admin sign-in, roles and audit log

**Goal.** Admins sign in with the normal e-mail code on the admin site; every request is checked against identity's
admin list (cached at most 60 s, so revocation is fast); every change and every look at personal data is written to
an append-only audit log.

**Flow.** `admin.<domain>/api/auth/*` (allow-listed e-mail code routes) → identity → session cookie → `/api/auth/token`
→ JWT → `/admin/v1/*` with Bearer → JWKS check → `IdentityAdminService.IsAdmin` (60 s cache) → role check → audit row
`pending` → downstream admin RPC → audit row `ok` / `error`.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [ADM-02-01](ADM-02-01-sign-in-proxy.md) | Sign-in through the admin origin | R1 | M | ADM-01-04, ADM-01-05, IDN-02 |
| [ADM-02-02](ADM-02-02-admin-check.md) | Admin check on every request | R1 | M | ADM-02-01, IDN-10-01 |
| [ADM-02-03](ADM-02-03-audit-log.md) | Audit log | R1 | M | ADM-01-03, ADM-02-02 |
| [ADM-02-04](ADM-02-04-admin-accounts.md) | Admin accounts list | R1 | S | ADM-02-02, ADM-02-03, IDN-10-01 |
