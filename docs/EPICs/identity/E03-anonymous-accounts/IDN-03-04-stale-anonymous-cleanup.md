# IDN-03-04 Delete stale anonymous accounts

| | |
|---|---|
| Service | identity |
| Epic | [IDN-03 Anonymous accounts and upgrade](README.md) |
| Release | R1 |
| Size | M |
| Depends on | IDN-03-01, IDN-09-02 |
| Unblocks | — |

## Goal

Anonymous accounts from people who tried the app once and left don't pile up, and their data (fridge photos'
results, plans) doesn't stay on the servers ("collect the minimum").

## Business rules

1. A job runs daily at 03:30 UTC under a Postgres advisory lock (one run at a time).
2. An anonymous account is stale when its newest session's `updated_at` (or the user's `created_at` if it has no
   session) is older than 30 days.
   Decided (Q044): 30 days without use.
3. Each stale account is deleted through the account deletion function of IDN-09-02, without a code, with
   `household.deleted` reason `anonymous_cleanup`, followed by `account.deleted`.
4. An anonymous account whose household has other account members is skipped and logged as a warning (it should
   not happen: anonymous accounts can't invite or join).
5. At most 500 accounts per transaction and 10,000 per run; the rest wait for the next day.
6. Metric `identity_anonymous_deleted_total`; a log line per run with the count (no ids).

## Contract

No RPC.

## Data

No schema change (uses `user_anonymous_idx` from IDN-03-01).

## Events

Publishes `household.deleted` (reason `anonymous_cleanup`) and `account.deleted` per account.

## Errors

| Case | Behaviour |
|---|---|
| Lock held by another run | Skip this run |
| A deletion fails | Log, continue with the next account |

## Acceptance criteria

- [ ] With test accounts last active 29 and 31 days ago, only the 31-day one is deleted.
- [ ] Its household and member are gone and both events are in the stream.
- [ ] 12,000 stale accounts are cleared over two runs (10,000 + 2,000).
