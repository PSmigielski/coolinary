# ADM-02-04 Admin accounts list

| | |
|---|---|
| Service | admin |
| Epic | [ADM-02 Admin sign-in, roles and audit log](README.md) |
| Release | R1 |
| Size | S |
| Depends on | ADM-02-02, ADM-02-03; IDN-10-01 (`IdentityAdminService.ListAdmins`, `make admin-grant` / `make admin-revoke`) |
| Unblocks | — |

## Goal

The founder can see who has back-office access and with which role. Granting and revoking happen in identity's
command line, so no network request can create an admin.

## Business rules

1. Only role `admin` may use this route.
2. `GET /admin/v1/admins` calls `IdentityAdminService.ListAdmins` and shows user id, e-mail, role (`admin` or
   `dietitian`), granted at and granted by.
3. Granting and revoking are not offered here. The page shows the commands instead:
   `make admin-grant email=<address> [role=dietitian]` and `make admin-revoke email=<address>`, run on the server
   (IDN-10-01 rule 2). The person must have signed in once (in the app or on admin) so the account exists; identity
   refuses anonymous accounts and refuses to revoke the last `admin`.
4. A grant or revoke takes effect in admin within 60 s (the `IsAdmin` cache, ADM-02-02 rule 2).
5. Viewing the list is audited as a `read` (`identity.admins.list`); the e-mail addresses are not stored in the audit
   row (ADM-02-03 rule 5).
6. Decided (Q048): admin rights are never granted from the back office; only identity's CLI grants them
   (IDN-10-01). There is no `SetAdminRole` RPC.

## Contract

| Method | Path | Response |
|---|---|---|
| GET | `/admin/v1/admins` | `{ items: [{ user_id, email, role, granted_at, granted_by }], grant_command, revoke_command }` |

Calls `IdentityAdminService.ListAdmins` (IDN-10-01).

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Role `dietitian` | 403 `role_required` |
| identity unreachable | 503 `identity_unavailable` |

## Acceptance criteria

- [ ] After `make admin-grant email=d@coolinary.test role=dietitian`, the list shows the account with role
      `dietitian`, and within 60 s that account can read the dietitian queue.
- [ ] A `dietitian` gets 403 on `GET /admin/v1/admins`.
- [ ] The audit row for a list view contains no e-mail address.
