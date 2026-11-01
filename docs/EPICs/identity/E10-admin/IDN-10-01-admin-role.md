# IDN-10-01 Admin role, `IsAdmin` and `IdentityAdminService`

| | |
|---|---|
| Service | identity |
| Epic | [IDN-10 Admin RPCs and admin role](README.md) |
| Release | R1 |
| Size | M |
| Depends on | IDN-04-01 |
| Unblocks | IDN-01-04, IDN-10-02 … 04, ADM-02 |

## Goal

Only people the founder names can use the back office, and naming an admin can't be done over the network.

## Business rules

1. Table `admin_users` lists admin accounts. Only accounts with an e-mail (not anonymous) can be admins.
2. Admins are granted and revoked only from a command inside the identity container:
   `bun run src/cli/admin.ts grant <email> [--role admin|dietitian]` / `revoke <email>`, with a Makefile shortcut
   `make admin-grant email=… [role=dietitian]` / `make admin-revoke email=…`. There is no RPC to grant admin rights.
   The default role is `admin`. `revoke` refuses to remove the last `admin` (a `dietitian` doesn't count).
   Decided (Q028): two roles, `admin` (everything) and `dietitian` (dish reading and dietitian checks only).
3. Not Better Auth's `admin` plugin: it adds admin HTTP endpoints that would be reachable through the gateway's
   public `/api/auth/*` proxy.
   Decided (Q027): own `admin_users` table + CLI.
4. `IsAdmin(user_id)` returns `is_admin` and `roles` (`admin`, `dietitian`). Anonymous or unknown users → false
   with no roles (not an error). `is_admin` is true for either role; the admin service checks the role per route.
5. Every `IdentityAdminService` RPC except `IsAdmin` requires metadata `x-admin-user-id` naming a current admin;
   otherwise `PERMISSION_DENIED`. This is a second check behind the admin service's own (defence in depth); the
   network already limits who can reach port 50051.
6. Every admin RPC call is logged with the admin user id, the RPC and the target id (never an e-mail).
7. Deleting an admin's account removes them from `admin_users` (IDN-09-02).
8. No JWT claim for admin rights (IDN-04-03 rule 3).
9. `ListAdmins` returns every row of `admin_users` with user id, e-mail, role, granted at and granted by. It is
   read-only; granting stays in the CLI (rule 2).
10. `AdminDeleteAccount` runs the account deletion of IDN-09-02 for support requests (IDN-09-02 rule 8).

## Contract

`contracts/coolinary/identity/v1/admin.proto`:

```proto
service IdentityAdminService {
  rpc IsAdmin(IsAdminRequest) returns (IsAdminResponse);
  rpc FindAccount(FindAccountRequest) returns (FindAccountResponse);                        // IDN-10-02
  rpc SetPlan(SetPlanRequest) returns (SetPlanResponse);                                    // IDN-10-03
  rpc GetAccountStats(GetAccountStatsRequest) returns (AccountStats);                       // IDN-10-04
  rpc RepublishHouseholdSnapshots(RepublishHouseholdSnapshotsRequest) returns (SnapshotJob); // IDN-01-04
  rpc GetSnapshotJob(GetSnapshotJobRequest) returns (SnapshotJob);                          // IDN-01-04
  rpc ListAdmins(ListAdminsRequest) returns (ListAdminsResponse);                           // rule 9
  rpc AdminDeleteAccount(AdminDeleteAccountRequest) returns (AdminDeleteAccountResponse);   // IDN-09-02
}
message IsAdminRequest { string user_id = 1; }
message IsAdminResponse { bool is_admin = 1; repeated string roles = 2; }
message ListAdminsRequest {}
message AdminUser { string user_id = 1; string email = 2; string role = 3; string granted_at = 4; string granted_by = 5; }
message ListAdminsResponse { repeated AdminUser admins = 1; }
```

Registered on the same gRPC server as `IdentityService`, with health status and reflection (dev) like it.

## Data

```ts
export const adminUsers = pgTable("admin_users", {
  userId: text("user_id").primaryKey().references(() => user.id, { onDelete: "cascade" }),
  grantedAt: timestamp("granted_at", { withTimezone: true }).notNull().defaultNow(),
  role: text("role", { enum: ["admin", "dietitian"] }).notNull().default("admin"),
  grantedBy: text("granted_by").notNull().default("cli"),
});
```

## Errors

| Case | gRPC status |
|---|---|
| Admin RPC without a valid `x-admin-user-id` | `PERMISSION_DENIED` |
| CLI grant for an unknown or anonymous account | CLI exits 1 with a message |

## Acceptance criteria

- [ ] `make admin-grant email=founder@coolinary.test` makes `IsAdmin` return true for that user.
- [ ] `FindAccount` without `x-admin-user-id`, or with a non-admin id, returns `PERMISSION_DENIED`.
- [ ] No `/api/auth/admin/*` route exists.
- [ ] `make admin-grant email=d@coolinary.test role=dietitian` makes `IsAdmin` return `is_admin = true`,
      `roles = ["dietitian"]`, and `ListAdmins` lists the account; `make admin-revoke` on the only `admin` exits 1.
