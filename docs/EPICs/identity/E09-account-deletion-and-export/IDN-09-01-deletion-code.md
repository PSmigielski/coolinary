# IDN-09-01 Deletion confirmation code

| | |
|---|---|
| Service | identity |
| Epic | [IDN-09 Account deletion and data export](README.md) |
| Release | R1 |
| Size | S |
| Depends on | IDN-02-01, NTF-02 (template `account_deletion_code`) |
| Unblocks | IDN-09-02, GW-02 |

## Goal

An account is deleted only by its owner holding the e-mail inbox, not by anyone who picks up an unlocked phone.

## Business rules

1. `RequestAccountDeletion(actor_user_id, language)` for an account with an e-mail: create a 6-digit code, valid
   10 minutes, 3 wrong tries, stored hashed in redis-state `identity:delcode:<user id>`; a new request replaces the
   old code.
2. Sent with notifier `SendEmail`, template `account_deletion_code`, variables `{code, valid_minutes: 10}`, in the
   given language (the gateway passes `x-language`).
3. Throttle: 1 request per 60 s and 3 per hour per account → `RESOURCE_EXHAUSTED` with retry time.
4. Anonymous accounts have no e-mail: the request returns `FAILED_PRECONDITION` `code_not_needed`, and
   `DeleteAccount` works without a code for them.
5. The code and e-mail are never logged.
6. Decided (Q137): the owner of a household with other account members gets `FAILED_PRECONDITION`
   `owner_must_transfer` and no code is sent; the app asks them to name a new owner first (IDN-08-03, IDN-08-01).
   Shared households exist from R2, so this check matters from R2.

## Contract

```proto
rpc RequestAccountDeletion(RequestAccountDeletionRequest) returns (RequestAccountDeletionResponse);

message RequestAccountDeletionRequest { string actor_user_id = 1; string language = 2; }
message RequestAccountDeletionResponse { string expires_at = 1; }
```

## Data

No schema change. redis-state key `identity:delcode:<user id>` (hash, attempts; TTL 600 s).

## Errors

| Case | gRPC status / reason |
|---|---|
| Anonymous account | `FAILED_PRECONDITION` `code_not_needed` |
| Owner of a shared household | `FAILED_PRECONDITION` `owner_must_transfer` |
| Throttled | `RESOURCE_EXHAUSTED` `rate_limited` |
| Notifier down | `UNAVAILABLE` `code_not_sent` |

## Acceptance criteria

- [ ] A request sends one e-mail visible in Mailpit; a second within 60 s is refused.
- [ ] The stored value is not the plain code.
- [ ] The owner of a household with another account member gets `owner_must_transfer` and no e-mail.
