# IDN-09-02 Delete an account

| | |
|---|---|
| Service | identity |
| Epic | [IDN-09 Account deletion and data export](README.md) |
| Release | R1 |
| Size | L |
| Status | todo |
| Depends on | IDN-09-01, IDN-01-01, IDN-04-01 |
| Unblocks | IDN-03-04, GW-02, KIT-13, PLN-21, OFF-14, NTF-01, ADM-05-03 |

## Goal

Deleting an account removes the person's data everywhere, keeps a shared household working for the people who
stay, and can't be undone halfway.

## Business rules

1. `DeleteAccount(actor_user_id, code)`. Accounts with an e-mail must pass the current code from IDN-09-01
   (3 wrong tries kill it). Anonymous accounts pass no code. The cleanup job (IDN-03-04) calls the same function
   without a code.
2. One transaction, in this order:
   1. Decided (Q137): if the actor owns a household with other account members, deletion is refused with
      `FAILED_PRECONDITION` `owner_must_transfer` and nothing changes. The owner first names a new owner
      (`TransferOwnership`, IDN-08-03, or the leave flow, IDN-08-01). There is no automatic hand-over. Shared
      households exist from R2, so this check matters from R2.
   2. If the actor is the only account member of their household: the household is deleted with its profiles,
      invites and plan history; emit `household.deleted` with reason `owner_deleted` (or `anonymous_cleanup` from
      IDN-03-04).
   3. Otherwise (a non-owner member): their member row is deleted.
   4. Pending invites the actor created stay valid with no creator (IDN-07-05 rule 4); invites of a deleted
      household go with it.
   5. Remove the actor from `admin_users` (IDN-10-01) if present.
   6. Emit `account.deleted` with `user_id` and `member_ids` (the actor's member id).
   7. Delete the `user` row; Better Auth's sessions, accounts and verifications cascade.
   Because of step i, a deleted user never owns a household that survives, so the cascade of
   `households.owner_user_id` on user delete can't remove a shared household.
3. After commit, write the token revocation key (IDN-08-04) once that task exists.
4. Deletion is final: no soft delete, no grace period, no undo. Backups age out on their normal schedule.
5. Events are the only signal to other services; identity doesn't wait for them.
6. Logs record "account deleted" with no e-mail and no user id; the metric `identity_accounts_deleted_total{kind}`
   counts with `kind` in `email`, `anonymous`, `cleanup`.
7. The answer is empty; the app signs out and starts a new anonymous account if the person continues.
8. Support deletion: `IdentityAdminService.AdminDeleteAccount(user_id, confirm_email, request_reference)` runs the
   same function without a code, for a person who asked support to delete their account. `confirm_email` must equal
   the account's e-mail (case-insensitive) or the call fails with `INVALID_ARGUMENT` `email_mismatch`; anonymous
   accounts can't be deleted this way (`FAILED_PRECONDITION`), and an owner of a shared household gets
   `owner_must_transfer` here too (support asks the person to name a new owner in the app first).
   `x-admin-user-id` is required (IDN-10-01 rule 5);
   the log line carries the admin id and `request_reference` (≤ 100 characters), never the e-mail.

## Contract

```proto
rpc DeleteAccount(DeleteAccountRequest) returns (DeleteAccountResponse);

message DeleteAccountRequest { string actor_user_id = 1; string code = 2; }
message DeleteAccountResponse {}

// IdentityAdminService (IDN-10-01), rule 8
rpc AdminDeleteAccount(AdminDeleteAccountRequest) returns (AdminDeleteAccountResponse);
message AdminDeleteAccountRequest { string user_id = 1; string confirm_email = 2; string request_reference = 3; }
message AdminDeleteAccountResponse {}
```

## Data

No schema change.

## Events

Publishes `household.deleted` (when the household goes) and `account.deleted`, in that order.

## Errors

| Case | gRPC status / reason |
|---|---|
| Code missing or wrong (e-mail accounts) | `INVALID_ARGUMENT` `invalid_code` |
| Third wrong code | `FAILED_PRECONDITION` `code_expired` |
| Unknown user | `NOT_FOUND` |
| Support deletion with a different e-mail | `INVALID_ARGUMENT` `email_mismatch` |
| Owner of a household with other account members | `FAILED_PRECONDITION` `owner_must_transfer` |

## Acceptance criteria

- [ ] A single-person account: user, household and member rows are gone; `household.deleted` then
      `account.deleted` are in the stream.
- [ ] An owner of a household with two other accounts gets `owner_must_transfer`; nothing is deleted and no event
      is emitted. After `TransferOwnership` to another member, deletion succeeds: the household survives with the
      new owner and only `account.deleted` is emitted.
- [ ] Network traffic and logs of the flow contain no e-mail address after the code e-mail.
- [ ] A deleted account's e-mail can sign up again as a brand-new account.
- [ ] `AdminDeleteAccount` with the right `confirm_email` deletes the account exactly like the in-app flow; with a
      different e-mail it returns `email_mismatch` and deletes nothing.
