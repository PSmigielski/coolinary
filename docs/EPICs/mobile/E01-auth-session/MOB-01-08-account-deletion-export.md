# MOB-01-08 Account deletion and data export calls

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-01 Auth session](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | MOB-01-07, MOB-02-01, GW-02 (`POST /v1/me/deletion-code`, `DELETE /v1/me`, `POST /v1/me/export`), IDN-09, MOB-03-03 (`useTransferOwnership`, R2: shared households start in R2) |
| Unblocks | later settings screens |

## Goal

The user can delete their account and ask for a copy of their data from the app, as GDPR requires and as the
stores require for apps with accounts.

## Business rules

1. Registered accounts confirm deletion with an e-mail code (GW-02-05, IDN-09-01): `useRequestDeletionCode()`
   calls `POST /v1/me/deletion-code` and returns `expiresAt` (10 minutes); `429` carries `Retry-After`. Anonymous
   accounts skip this step (the server answers `422 code_not_needed` if asked).
2. `useDeleteAccount()` calls `DELETE /v1/me` with an `Idempotency-Key` and body `{ confirm: "DELETE", code? }`
   (`code` for registered accounts only). It is never queued offline. The hook takes `{ confirmed: true }` from
   the UI and refuses otherwise.
3. `400 invalid_code` lets the user retype; `422 code_expired` (third wrong try) sends the user back to step 1.
   Decided (Q137): the owner of a household with another account member gets `422 owner_must_transfer` on both
   calls. The hook returns that as its own result `{ kind: "ownerMustTransfer" }`; the flow then asks the owner to
   pick the new owner (`useTransferOwnership()`, MOB-03-03, `POST /v1/household/owner`) and starts again at step 1.
   An owner who is the only account member deletes normally, and the household goes with the account.
4. On `200`/`202`: run `wipeLocalData("account_deleted")` (MOB-01-07), which also deletes `health.profile`
   (Q136).
5. `useRequestExport()` calls `POST /v1/me/export`. The gateway assembles the export (README D9) and answers `200`
   with the JSON document as an attachment (GW-02-06; decided, Q143: a synchronous download, no e-mailed link). The app writes it to a temporary file, opens the system share sheet, and deletes the file afterwards;
   the document never enters the query cache or AsyncStorage.
6. Decided (Q142): anonymous users can delete (wipes the anonymous household) but cannot export; the export hook
   returns `sign_in_required` for them without a request (the gateway would answer `403 registration_required`).

## Contract

| Call | Request | Responses |
|---|---|---|
| `POST /v1/me/deletion-code` | `{}` | `202 { expiresAt }`, `422 code_not_needed`, `422 owner_must_transfer`, `429 rate_limited` |
| `DELETE /v1/me` | `{ confirm: "DELETE", code? }` | `202 { erasure: "in_progress" }`, `400 invalid_code`, `422 code_expired`, `422 owner_must_transfer`, `401` |
| `POST /v1/me/export` | `{}` | `200` JSON attachment, `429 rate_limited`, `503 export_unavailable` |

## Data

No storage.

## Errors

| Case | Behaviour |
|---|---|
| Offline | `OfflineError`; nothing queued |
| `429` on export | Show "you can request 3 exports a day" with `Retry-After` (GW-02-06 rule 6) |

## Acceptance criteria

- [ ] After deletion the phone is in the clean state of MOB-01-07 and `GET /v1/me` with the old session returns 401.
- [ ] A registered user requests a code, deletes with it, and a wrong code first returns 400 without deleting.
- [ ] The owner of a two-account household gets `ownerMustTransfer`; after naming the other member as owner, the
      deletion goes through and the other member keeps the household.
- [ ] Export as a signed-in user opens the share sheet; afterwards no export file remains on the device and nothing
      from it is in the query cache.
