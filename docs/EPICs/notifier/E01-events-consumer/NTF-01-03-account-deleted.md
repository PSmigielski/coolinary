# NTF-01-03 Erase a deleted account's data

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-01 Events consumer](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | NTF-01-01, NTF-01-02, NTF-02-02, NTF-03-02, NTF-04-01, NTF-04-02; IDN-09 (`account.deleted`) |
| Unblocks | — |

## Goal

When someone deletes their account (GDPR erasure), notifier removes everything it holds about them within the same
event handling.

## Business rules

1. On `account.deleted` (`user_id`, `member_ids`), in one transaction:
   1. delete every device of `user_id` (push tokens);
   2. delete the user's notification preferences;
   3. mark every `household_users` pair of `user_id` removed (NTF-01-02 tombstone rules);
   4. set every queued push row of `user_id` that is not final to `cancelled` and clear its rendered title, body and
      data;
   5. in the delivery log, set `user_id` and `recipient_hash` to null on all of the user's rows (counts stay for
      statistics, the link to the person is gone).
2. An unknown `user_id` is a no-op (idempotent).
3. Founder alerts never contain user data, so they are untouched.
4. `member_ids` is not used by notifier.

## Contract

No gRPC change.

## Data

No schema change (uses tables from NTF-01-02, NTF-02-02, NTF-03-02, NTF-04-01, NTF-04-02).

## Events

Consumes `account.deleted` from `events.identity`.

## Errors

| Case | Behaviour |
|---|---|
| Postgres error | Rolled back, retried (NTF-01-01 rule 8) |

## Acceptance criteria

- [ ] A user with 2 devices, custom preferences, 3 queued pushes and 10 delivery rows: after the event, no device
      or preference row exists, the 3 pushes are `cancelled` with empty texts, the 10 delivery rows have null
      `user_id` and `recipient_hash`.
- [ ] `GetPreferences` for that user returns defaults; `ExportUserData` returns empty lists.
- [ ] Running the event twice changes nothing the second time.
