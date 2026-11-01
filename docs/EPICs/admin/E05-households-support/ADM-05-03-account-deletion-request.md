# ADM-05-03 Delete an account on request

| | |
|---|---|
| Service | admin |
| Epic | [ADM-05 Households and plans support](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | ADM-05-01; IDN-09-02 (`IdentityAdminService.AdminDeleteAccount`) |
| Unblocks | — |

## Goal

A person who asks support to delete their account (GDPR Article 17), for example because they lost access to the
app, can have it deleted the same way as when they do it in the app.

## Business rules

1. Input: user id, `confirm_email` (must equal the account's e-mail; identity compares), and a request reference
   (≤ 100 characters, e.g. the support ticket id). Admin calls `AdminDeleteAccount(user_id, confirm_email,
   request_reference)` (IDN-09-02 rule 8).
2. identity runs the same deletion as the in-app deletion (IDN-09-02): ownership handed over or household deleted,
   `account.deleted` published, the other services erase their part. Anonymous accounts can't be deleted this way
   (identity's cleanup job removes them).
3. Back-office accounts are not deleted here: if `IsAdmin(user_id)` is true, admin refuses with 409 `is_admin`
   (revoke first with `make admin-revoke`, ADM-02-04).
4. Audited (`identity.account.delete`, user id, request reference); the e-mail is not stored.
5. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/accounts/{userId}/deletion` | `{ confirm_email, request_reference }` → `{ deleted: true }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| `confirm_email` does not match | 409 `email_mismatch` (identity's `INVALID_ARGUMENT`) |
| Unknown user | 404 |
| Anonymous account | 409 `anonymous_account` (identity's `FAILED_PRECONDITION`) |
| Account has back-office access | 409 `is_admin` |

## Acceptance criteria

- [ ] After deletion, `ADM-05-01` lookup by the old e-mail returns 404 and notifier holds no device for the user.
