# NTF-03-04 Push receipts and cleanup

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-03 Push delivery](README.md) |
| Release | R1 |
| Size | S |
| Depends on | NTF-03-03 |
| Unblocks | — |

## Goal

notifier learns whether Apple and Google accepted each push, removes tokens that no longer work, and deletes message
texts as soon as they are no longer needed.

## Business rules

1. Every 15 minutes, a job takes rows with `status = 'sent'` and `sent_at <= now() - 15 min`, in groups of at most
   1,000 ticket ids, and calls `POST https://exp.host/--/api/v2/push/getReceipts`.
2. Receipt `ok` → row and delivery `delivered`.
3. Receipt error `DeviceNotRegistered` → the device is deleted, row `failed` / `expo_DeviceNotRegistered`.
4. Any final state (`delivered`, `failed`, `expired`, `cancelled`, `unknown`) clears `title`, `body` and `data` and
   sets `final_at`.
5. `InvalidCredentials` in a receipt → founder alert `notifier.expo_credentials`.
6. A row still `sent` 24 h after `sent_at` without a receipt → `unknown`.
7. Queue rows are deleted 48 h after `final_at`; delivery rows follow NTF-02-02 retention.
8. A failed receipts request is retried at the next run; it never marks rows failed.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Expo receipts endpoint down | Logged; retried in 15 minutes |

## Acceptance criteria

- [ ] With a mock: 2,500 sent rows older than 15 min → 3 receipt requests (1,000, 1,000, 500).
- [ ] After the job, final rows have empty `title`, `body` and `data`.
- [ ] Rows final for 49 h are gone.
- [ ] Metric `push_receipts_total{result}`.
