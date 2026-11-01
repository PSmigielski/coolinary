# GW-02-02 Better Auth proxy allowlist for anonymous and e-mail code sign-in

| | |
|---|---|
| Service | gateway |
| Epic | [GW-02 Auth and account routes](README.md) |
| Release | R1 |
| Size | M |
| Depends on | GW-01-01, IDN-02 (e-mail OTP), IDN-03 (anonymous) |
| Unblocks | MOB-01 |

## Goal

Only the Better Auth endpoints the app needs reach identity, with the client's IP and cookies intact, so sign-in
works through the gateway and password sign-up cannot be reached even by accident.

## Business rules

1. The catch-all `.all("/api/auth/*")` in `src/http/app.ts` is replaced by an allowlist (method + path, config file
   `src/http/auth-routes.ts`). Every other `/api/auth/*` path returns 404 `not_found` without calling identity.
2. Allowed (paths as served by Better Auth 1.7 with the plugins chosen in IDN-02/IDN-03; IDN-02 confirms the exact
   list when it ships):

   | Method | Path | Use |
   |---|---|---|
   | POST | `/api/auth/sign-in/anonymous` | first start (IDN-03) |
   | POST | `/api/auth/email-otp/send-verification-otp` | send the 6-digit code (IDN-02) |
   | POST | `/api/auth/sign-in/email-otp` | sign in / link anonymous with the code (IDN-02, IDN-03) |
   | GET | `/api/auth/get-session` | session check (Expo client) |
   | GET | `/api/auth/token` | 15-minute JWT for `/v1` |
   | GET | `/api/auth/jwks` | public keys (also used by admin) |
   | POST | `/api/auth/sign-out` | end the session |
   | GET | `/api/auth/list-sessions` | "signed-in devices" |
   | POST | `/api/auth/revoke-session` | sign a device out |

   Explicitly not allowed: `/api/auth/sign-up/email`, `/api/auth/sign-in/email`, password reset routes.
3. Forwarding keeps method, path, query, body and all headers except hop-by-hop ones (`Connection`, `Keep-Alive`,
   `Transfer-Encoding`, `Upgrade`, `TE`, `Trailer`, `Proxy-*`). It sets `X-Forwarded-For` to the client IP from
   GW-01-01 (replacing any incoming value) and `X-Forwarded-Proto`. Identity trusts `X-Forwarded-For` only from the
   gateway (IDN-02 config), so its own rate limits see the real client.
4. Responses pass back unchanged, `Set-Cookie` included; redirects are not followed (`redirect: "manual"`, exists).
5. Deadline 5 s; identity unreachable → 503 problem `service_unavailable`.
6. Bodies are never logged. The request log records the path template and status only.
7. Rate limits for the sign-in endpoints are applied before forwarding (GW-03-04, buckets `otp_send`,
   `otp_verify`, `anon_signup`).
8. The README's "Try the auth flow" (repo root) is updated to the anonymous + e-mail code flow in the same change.

## Contract

Pass-through HTTP; excluded from OpenAPI (exists). The allowlist table above is the contract MOB-01 relies on.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Path not on the allowlist | 404 | `not_found` |
| Identity down or slow | 503 | `service_unavailable` |
| Rate limit | 429 | `rate_limited` |

## Acceptance criteria

- [ ] `POST /api/auth/sign-up/email` returns 404 and identity receives nothing (fake identity counts calls).
- [ ] Anonymous sign-in → token → `GET /v1/me` works end to end in `make dev`.
- [ ] E-mail code sign-in: code arrives in Mailpit, `sign-in/email-otp` returns a session cookie, `/api/auth/token`
      returns a JWT accepted by `/v1/me`.
- [ ] A spoofed `X-Forwarded-For: 1.2.3.4` from the client reaches identity replaced by the real client IP.
- [ ] Hop-by-hop headers are stripped; `Set-Cookie` arrives at the client.
