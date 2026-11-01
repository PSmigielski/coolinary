# GW-02 Auth and account routes

Release R1.

## Goal

The app can start without an account (anonymous), sign in with an e-mail code, upgrade the anonymous account,
read who it is and what its plan allows, delete the account, and download its data.

## Flow

1. First start: app calls `POST /api/auth/sign-in/anonymous` (proxied to identity, IDN-03) and gets a session; then
   `GET /api/auth/token` for the 15-minute JWT the `/v1` routes accept.
2. Sign-in: `POST /api/auth/email-otp/send-verification-otp`, then `POST /api/auth/sign-in/email-otp` (IDN-02). If
   the device holds an anonymous session, identity links it (IDN-03) and the household is kept.
3. The app fetches a fresh JWT; `GET /v1/me` and `GET /v1/entitlements` describe the account.
4. `DELETE /v1/me` and `POST /v1/me/export` serve GDPR rights.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-02-01](GW-02-01-jwt-claims.md) | JWT claims: member id, anonymous flag, registered-only routes | R1 | S | GW-01-02, IDN-03, IDN-04 |
| [GW-02-02](GW-02-02-auth-proxy-allowlist.md) | Better Auth proxy allowlist for anonymous and e-mail code sign-in | R1 | M | GW-01-01, IDN-02, IDN-03 |
| [GW-02-03](GW-02-03-me.md) | `GET /v1/me` | R1 | S | GW-02-01 |
| [GW-02-04](GW-02-04-entitlements.md) | `GET /v1/entitlements` | R1 | S | GW-02-01, GW-03-02, IDN-05 |
| [GW-02-05](GW-02-05-delete-account.md) | Account deletion with an e-mail code | R1 | M | GW-02-01, GW-01-03, IDN-09 |
| [GW-02-06](GW-02-06-data-export.md) | `POST /v1/me/export` (GDPR export) | R1 | M | GW-02-01, IDN-09, KIT-13, PLN-21, NTF-04 |
