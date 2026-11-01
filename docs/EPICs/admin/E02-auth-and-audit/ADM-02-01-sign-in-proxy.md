# ADM-02-01 Sign-in through the admin origin

| | |
|---|---|
| Service | admin |
| Epic | [ADM-02 Admin sign-in, roles and audit log](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | ADM-01-04, ADM-01-05; IDN-02 (e-mail code routes, `TRUSTED_ORIGINS` includes the admin origin) |
| Unblocks | ADM-02-02 |

## Goal

Admins sign in with the same e-mail code as everyone else (decision 12), on `admin.<domain>`, so the session cookie
belongs to the admin site and the admin UI gets the same JWT the gateway accepts.

## Business rules

1. admin forwards a fixed list of Better Auth paths to identity's HTTP server (`IDENTITY_HTTP_URL`), unchanged,
   the way the gateway does today (method, headers, body, `redirect: "manual"`):
   - send a sign-in code (Better Auth e-mail OTP "send verification OTP" route, type `sign-in`),
   - sign in with the code (e-mail OTP sign-in route),
   - `get-session`, `token` (JWT for the API), `jwks`, `sign-out`.
   Exact paths are the ones IDN-02 configures.
2. Every other `/api/auth/*` path answers 404. In particular anonymous sign-in and any sign-up route are not
   reachable through admin.
3. identity must list `https://admin.<domain>` (prod) and `http://10.77.0.23:3000` (dev) in `TRUSTED_ORIGINS`;
   this is a change in identity (IDN-02 or IDN-10) and a dependency of this task.
4. Code sends through admin are limited to 5 per IP per 15 minutes (in memory; admin runs one container), on top of
   identity's own limits. Over the limit → 429.
5. Anyone can request a code here; only accounts identity marks as admin get past ADM-02-02. Responses for
   non-admin addresses are identical to admin ones (no hint which addresses are admins).
6. Cookies are set by identity's responses; with no `Domain` attribute they belong to the admin host. Prod is
   HTTPS-only, so `Secure` applies.
7. The admin UI then calls `GET /api/auth/token` for a 15-minute JWT and sends it as `Authorization: Bearer` to
   `/admin/v1/*`; it refreshes the JWT before expiry.

## Contract

| Method | Path | Behaviour |
|---|---|---|
| any | `/api/auth/<allowed path>` | Forwarded to identity |
| any | `/api/auth/<other>` | 404 |

## Data

None.

## Errors

| Case | HTTP |
|---|---|
| Path not on the list | 404 |
| Over 5 code requests per IP per 15 min | 429 `rate_limited` |
| identity down | 503 `service_unavailable` |

## Acceptance criteria

- [ ] In dev: request a code for an admin account at `10.77.0.23:3000`, read it in Mailpit, sign in, get a JWT.
- [ ] `POST /api/auth/sign-in/anonymous` on admin → 404.
- [ ] 6th code request from one IP in 15 minutes → 429.

## Notes

Decided (Q046): sign-in goes through the admin service's proxy (first-party cookies); the public API gets no CORS
for the admin origin.
