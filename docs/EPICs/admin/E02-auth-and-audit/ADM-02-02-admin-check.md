# ADM-02-02 Admin check on every request

| | |
|---|---|
| Service | admin |
| Epic | [ADM-02 Admin sign-in, roles and audit log](README.md) |
| Release | R1 |
| Size | M |
| Depends on | ADM-02-01, ADM-01-05; IDN-10-01 (`IdentityAdminService.IsAdmin`) |
| Unblocks | every `/admin/v1/*` route |

## Goal

Only people identity marks as admins can use the back office, and removing someone's admin role takes effect within
a minute.

## Business rules

1. Every `/admin/v1/*` request needs `Authorization: Bearer <JWT>`. The JWT is verified with `jose` against
   identity's JWKS (`<IDENTITY_HTTP_URL>/api/auth/jwks`, cooldown 30 s), issuer and audience `AUTH_ISSUER` (the
   same checks as the gateway). Missing, invalid or expired → 401.
2. Then `IdentityAdminService.IsAdmin(user_id)` returns `is_admin` and `roles`. The answer (yes or no) is cached per
   user for 60 s in memory.
3. When identity cannot be reached and the cache entry is older than 60 s → 503 `identity_unavailable`. A missing or
   stale answer is never treated as "yes".
4. Not an admin → 403 `not_admin`, and an audit row `auth.denied` (ADM-02-03), at most one per user per minute.
5. Roles: `admin` may use every route. `dietitian` may only read dishes and record dietitian checks (ADM-04-03 read
   routes, ADM-04-05). Identity stores the role per admin account and grants it only from its CLI (IDN-10-01 rules
   2 and 4). Decided (Q028): these two roles.
6. Each route declares its required role (Elysia macro `role: "admin" | "dietitian"`); a route without a declaration
   requires `admin`.
7. The request context carries `adminUserId`, `roles`, and the e-mail from the JWT (used only for
   `SendTestEmail`, ADM-07-04); `adminUserId` is forwarded as `x-admin-user-id`.
8. CSRF guard: on `POST`, `PUT`, `PATCH`, `DELETE`, a present `Origin` header must equal `ADMIN_PUBLIC_URL`;
   otherwise 403.
9. Anonymous accounts can never pass (identity answers `is_admin = false` for them).

## Contract

`IdentityAdminService.IsAdmin(IsAdminRequest{user_id}) → IsAdminResponse{is_admin, roles[]}` (IDN-10-01).
`IsAdmin` needs no `x-admin-user-id`; every other identity admin RPC does (IDN-10-01 rule 5).

## Data

None.

## Errors

| Case | HTTP |
|---|---|
| No or bad JWT | 401 `unauthorized` |
| Not admin, or role too low | 403 `not_admin` / `role_required` |
| identity unreachable, no fresh cache | 503 `identity_unavailable` |
| Foreign `Origin` on a write | 403 `bad_origin` |

## Acceptance criteria

- [ ] A normal user's JWT gets 403 on `/admin/v1/overview`.
- [ ] Revoking the admin role in identity (`make admin-revoke`): the next request after at most 60 s gets 403.
- [ ] With identity stopped: requests succeed for up to 60 s after the last check, then 503.
- [ ] A `dietitian` gets 403 on `POST /admin/v1/kitchen/ingredient-merges`.
