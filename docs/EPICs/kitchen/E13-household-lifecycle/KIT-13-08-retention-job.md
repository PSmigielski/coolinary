# KIT-13-08 Retention clean-up job

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-07, KIT-08-05, KIT-08-06, KIT-09-02 |
| Unblocks | — |

## Goal

Rows that only exist for idempotency or short-term history are deleted on time, so kitchen keeps the minimum
(business model §10 "collect the minimum").

## Business rules

1. A daily job at 03:30 UTC (advisory lock `kitchen.retention`) deletes, in batches of 5,000 rows:

   | Table | Rows older than |
   |---|---|
   | `stock_movements` | 24 months (`occurred_at`, Decided Q141); monthly totals in `savings_monthly` stay |
   | `cook_events` | 30 days |
   | `stock_add_refs` | 90 days |
   | `receipt_fingerprints` | 90 days |
   | `reminder_runs` | 60 days |
   | `receipt_name_links` | 365 days since `last_seen_on` |
   | `unmatched_name_attempts` | 90 days after the last attempt with no unmatched item left |

2. Each batch is its own transaction; the job stops after 10 minutes and continues the next night.
3. Metrics: `kitchen_retention_deleted_rows_total{table}`.

## Contract

No RPC.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Lock held | skipped this night |

## Acceptance criteria

- [ ] Test data with rows on both sides of each limit: only the older ones are deleted.
- [ ] Deleting old movements leaves `savings_monthly` totals unchanged.
