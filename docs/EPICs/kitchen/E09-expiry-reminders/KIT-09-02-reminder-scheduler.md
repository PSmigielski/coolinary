# KIT-09-02 Daily reminder scheduler per household timezone

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-09 Expiry reminders](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-09-01, KIT-10-01 (timezone), KIT-01 (job runner lock) |
| Unblocks | KIT-09-03, KIT-12-04 (same pattern) |

## Goal

Run the reminder once per household per local day, at a sensible local time, surviving restarts and downtime
without sending twice.

## Business rules

1. The job ticks every 15 minutes inside kitchen. A Postgres advisory lock (`pg_try_advisory_lock`, key
   `kitchen.expiry_reminders`) makes sure only one runner works at a time, also across restarts.
2. **Due households:** local time ≥ 09:00 and < 21:00, at least one dated stock item, no run row for today's local
   date, and plan `plus` from identity (`PlanOf`, KIT-01-06, asked once per household per run). A household whose
   plan is `free` gets the run row `not_plus` for the day, so it is asked again only tomorrow; identity down → the
   household is skipped this tick.
   Decided (Q134): 09:00 local time, in time to plan the day's meals. Notifier
   delays anything that falls in a user's quiet hours (NTF-03-02).
3. Per tick at most 500 households, ordered by household id; the rest wait for the next tick.
4. Per household: KIT-09-01 selection → KIT-09-03 send. The run row (household, local date, outcome) is written when
   the push was accepted by notifier or nothing was due.
5. **Notifier unreachable:** no run row; the household is tried again next tick, until 21:00 local. At 21:00 the
   row is written with outcome `skipped` (no late-night reminders).
6. **Downtime:** after a restart the job continues the same day's households from the first tick; a day that passed
   completely is not caught up.
7. Households without a kitchen profile (no timezone) are skipped until KIT-13-01 creates it.
8. Metrics: `kitchen_expiry_reminders_total{outcome="sent|nothing_due|not_plus|failed|skipped"}`,
   `kitchen_expiry_reminder_tick_duration_seconds`.

## Contract

No RPC (background job).

## Data

`kitchen/migrations/000NN_reminder_runs.sql`

```sql
-- +goose Up
CREATE TABLE reminder_runs (
    household_id uuid NOT NULL,
    kind         text NOT NULL,                 -- 'expiry' (KIT-09), 'savings' (KIT-12-04)
    local_date   date NOT NULL,
    outcome      text NOT NULL CHECK (outcome IN ('sent', 'nothing_due', 'skipped')),
    created_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, kind, local_date)
);
-- Rows older than 60 days are deleted by KIT-13-08.

-- +goose Down
DROP TABLE reminder_runs;
```

## Errors

| Case | Behaviour |
|---|---|
| Lock held by another runner | tick skipped |
| One household fails | logged (household hash only), others continue |

## Acceptance criteria

- [ ] A Plus household in `Europe/Warsaw` with an item expiring today gets one push between 09:00 and 09:15 local
      and none at the next ticks.
- [ ] Two kitchen processes started together send one push (advisory lock).
- [ ] With notifier down from 09:00 to 10:20 local, the push goes out at the 10:30 tick; with notifier down all day
      the row says `skipped` at 21:00 and nothing is sent.
- [ ] A household that became Plus at 15:00 local gets its first reminder that day if items are due.
