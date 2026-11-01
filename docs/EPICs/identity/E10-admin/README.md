# IDN-10 Admin RPCs and admin role

Decision 9: the admin service (ADM) is used only by accounts identity marks as admin. Identity owns that mark and
offers `IdentityAdminService` (in `contracts/coolinary/identity/v1/admin.proto`, same gRPC port), called only by
admin, for support work: find an account, set a plan by hand, see account numbers, republish snapshots (IDN-01-04).

Admins sign in to the admin UI with the normal e-mail code flow (through the gateway); the admin service then asks
`IsAdmin` for the JWT's `sub` on each request (cached briefly there, ADM-02).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-10-01 | Admin role, `IsAdmin` and `IdentityAdminService` | R1 | M | IDN-04-01 |
| IDN-10-02 | Find an account | R1 | S | IDN-10-01, IDN-05-01 |
| IDN-10-03 | Set a plan by hand | R1 | S | IDN-10-01, IDN-05-01 |
| IDN-10-04 | Account numbers | R1 | S | IDN-10-01 |
