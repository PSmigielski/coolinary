# KIT-12-04 Monthly savings push

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-12 Savings](README.md) |
| Release | R2 |
| Size | S |
| Depends on | KIT-12-03, KIT-09-02 (scheduler, `reminder_runs`), NTF-03 (`SendPush`, template `savings_monthly`) |
| Unblocks | — |

## Goal

At the start of each month, Plus households get one push with last month's result, which keeps the value of the
subscription visible.

## Business rules

1. Runs with the KIT-09-02 scheduler pattern (15-minute ticks, advisory lock, batches of 500), kind `savings`.
2. Due: on the 1st of the month, local time ≥ 10:00 and < 21:00, no `savings` run row for that date, and plan
   `plus` from identity (`PlanOf`, KIT-01-06, once per household per run; `free` → run row `not_plus`). Decided (Q213): the 1st of the month at 10:00 local time.
3. Sent only when last month's `money_saved` > 0 or `kg_saved` ≥ 0.1; otherwise run row `nothing_due`.
4. `SendPush`: target `household_id`, category `CATEGORY_SAVINGS`, template `savings_monthly`, variables `month`
   (`YYYY-MM`), `amount` (the household currency's `money_saved`), `kg` (one decimal), `path = "/savings"`,
   `dedup_key = "savings:<YYYY-MM>"`, `ttl_seconds = 259200` (3 days).
5. Free households get no push (the report is Plus).
6. Notifier unavailable → retried at the next tick until 21:00, then `skipped`; there is no second attempt on the 2nd.

## Contract

Calls `NotifierService/SendPush` (NTF-03-02). No new kitchen RPC.

## Data

No schema change (uses `reminder_runs` with `kind = 'savings'`).

## Errors

| Case | Behaviour |
|---|---|
| notifier unavailable | retry next tick, `skipped` at 21:00 |

## Acceptance criteria

- [ ] A Plus household with 38 € saved in October gets one push on 1 November between 10:00 and 10:15 local.
- [ ] A household with nothing saved gets none (`nothing_due`).
- [ ] A Free household gets none.
