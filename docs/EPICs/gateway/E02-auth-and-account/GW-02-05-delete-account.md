# GW-02-05 Account deletion (`POST /v1/me/deletion-code`, `DELETE /v1/me`)

| | |
|---|---|
| Service | gateway |
| Epic | [GW-02 Auth and account routes](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-02-01, GW-01-03, IDN-09 (`RequestAccountDeletion`, `DeleteAccount`) |
| Unblocks | MOB-03 |

## Goal

A user can delete their account from the app (GDPR right to erasure). Identity deletes the account and publishes
`account.deleted`; every other service erases its data from that event (D9).

## Business rules

1. Any account may delete itself, anonymous included.
2. **Registered accounts confirm with an e-mail code.** `POST /v1/me/deletion-code` (body `{}`) calls identity
   `RequestAccountDeletion(actor_user_id, language)` (IDN-09-01): identity e-mails a 6-digit code, valid 10 minutes,
   3 wrong tries, throttled to 1 per 60 s and 3 per hour (identity answers `RESOURCE_EXHAUSTED` → 429 with
   `Retry-After`). For an anonymous account identity answers `FAILED_PRECONDITION` `code_not_needed` → 422; the app
   then deletes without a code.
3. `DELETE /v1/me` body `{ "confirm": "DELETE", "code"? }`: `confirm` is required (protects against an accidental
   call); `code` is required for registered accounts and absent for anonymous ones. Calls identity
   `DeleteAccount(actor_user_id, code)` (IDN-09-02). A missing or wrong code is identity's `INVALID_ARGUMENT`
   `invalid_code` (→ 400); the third wrong try is `FAILED_PRECONDITION` `code_expired` (→ 422) and the app asks for
   a new code (GW-01-02 passes identity's `x-error-code` through).
4. Decided (Q137): the owner of a household that has another account member must name the new owner first
   (`POST /v1/household/owner`, GW-04-05); until then identity refuses both the code request and the deletion with
   `FAILED_PRECONDITION` `owner_must_transfer` (→ 422, passed through). There is no automatic hand-over. An owner who
   is the only account member deletes normally; the household (and its profiles without accounts) is deleted too.
5. Success → 202 `{ "erasure": "in_progress" }`. Identity's answer is empty; erasure in the other services is
   asynchronous, which the response says.
6. After success the gateway deletes the user's idempotency entries (`SCAN` + `DEL gw:idem:<userId>:*`, best effort)
   and the cached entitlements of the household (`gw:ent:<householdId>`). Quota counters expire on their own.
7. Idempotency: required on `DELETE /v1/me`, so a retried deletion does not fail with "not found". Not used on
   `POST /v1/me/deletion-code` (a retry sends a new code, which identity throttles).

## Contract

`POST /v1/me/deletion-code` — `auth: true` (registered), tag `account`.

```json
request  {}
202      { "expiresAt": "2026-11-01T10:10:00Z" }
```

`DELETE /v1/me` — `auth: true`, tag `account`, idempotency required.

```json
request  { "confirm": "DELETE", "code": "482913" }
202      { "erasure": "in_progress" }
```

RPCs: identity `RequestAccountDeletion`, `DeleteAccount` (IDN-09).

## Data

Deletes `gw:idem:<userId>:*` and `gw:ent:<householdId>`.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| `confirm` missing or wrong | 400 | `validation` |
| Code missing or wrong (registered account) | 400 | `invalid_code` |
| Third wrong code | 422 | `code_expired` |
| Code requested for an anonymous account | 422 | `code_not_needed` |
| Too many code requests | 429 | `rate_limited` (with `Retry-After`) |
| Owner of a household with another account member | 422 | `owner_must_transfer` |
| Identity down | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] An anonymous account deletes itself with `{ "confirm": "DELETE" }` and no code; its old token then fails at
      the next token refresh (identity revoked the sessions).
- [ ] A registered account requests a code, then deletes with it; deleting with a wrong code returns 400
      `invalid_code` and the account still exists.
- [ ] A second code request within 60 s returns 429.
- [ ] After deletion, `account.deleted` appears on `events.identity` (integration test with identity).
- [ ] A retried `DELETE` with the same idempotency key returns the same 202.
- [ ] The owner of a household with a second account member gets 422 `owner_must_transfer` on both
      `POST /v1/me/deletion-code` and `DELETE /v1/me`; after `POST /v1/household/owner` names the other member, the
      deletion goes through.

## Notes

- The gateway logs neither the code nor the e-mail address.
- Re-authentication is the e-mail code from identity (IDN-09-01), so the gateway needs no session-age rule.
