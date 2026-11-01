# KIT-01-05 Dead letters, lag metrics and alerts

| | |
|---|---|
| Service | kitchen (shared platform code; copy to planner and offers) |
| Epic | [KIT-01 Platform](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-01-03, KIT-01-04, NTF-05 (SendAlert) |
| Unblocks | — |

## Goal

When events stop flowing or keep failing, the founder hears about it and the failing events are kept for replay.

## Business rules

1. A dead letter is an entry that failed 10 deliveries (KIT-01-04 rule 8). It is copied to `dead_letters` with its
   stream id, envelope and last error, then acknowledged.
2. Admin RPC `ReplayDeadLetter(id)` runs the handler again inside the normal idempotent transaction; on success the
   row gets `replayed_at`. Admin RPC `ListDeadLetters` lists unreplayed rows, newest first, 50 per page.
3. Alerts go through notifier `SendAlert` with a dedup key, so one cause sends one alert per 6 hours:
   - `kitchen.dead_letter.<type>` when a dead letter is written;
   - `kitchen.outbox_stuck` when `outbox_oldest_age_seconds` > 300;
   - `kitchen.consumer_lag.<stream>` when the group's `lag` (from `XINFO GROUPS`) stays above 1,000 entries or the
     oldest pending entry is older than 15 minutes, checked every minute;
   - `kitchen.consumer_gap.<stream>` when the last processed entry id is older than the stream's first entry
     (kitchen fell behind the 7-day retention and missed events).
4. If notifier is down, the alert is logged at `error` and retried on the next check.
5. Metric `events_consumer_lag{stream}` (gauge) and `events_dead_letters_total{type}` (counter).

## Contract

`contracts/coolinary/kitchen/v1/admin.proto`

```proto
rpc ListDeadLetters(ListDeadLettersRequest) returns (ListDeadLettersResponse);
rpc ReplayDeadLetter(ReplayDeadLetterRequest) returns (ReplayDeadLetterResponse);
message DeadLetter { int64 id = 1; string stream = 2; string type = 3; string event_id = 4;
  string last_error = 5; string failed_at = 6; }
```

## Data

`kitchen/migrations/000NN_dead_letters.sql`

```sql
CREATE TABLE dead_letters (
    id           bigserial PRIMARY KEY,
    stream       text NOT NULL,
    entry_id     text NOT NULL,
    event_id     uuid NOT NULL,
    type         text NOT NULL,
    envelope     jsonb NOT NULL,
    last_error   text NOT NULL,
    failed_at    timestamptz NOT NULL DEFAULT now(),
    replayed_at  timestamptz
);
```

## Errors

| Case | Status | Message |
|---|---|---|
| Replay of an unknown id | `NOT_FOUND` | `dead letter not found` |
| Replay fails again | `FAILED_PRECONDITION` | `handler failed: <short reason>` |

## Acceptance criteria

- [ ] A handler that always fails produces one dead letter after 10 deliveries and one alert.
- [ ] Ten dead letters of the same type within an hour send one alert.
- [ ] Replaying after fixing the handler applies the effect once.
- [ ] Simulated gap (stream trimmed past the group's position) raises `kitchen.consumer_gap.events.identity`.

## Notes

Decided (Q042): no automatic resync in R1. The alert links to a runbook: ask identity to republish household
snapshots (IDN-01-04) by hand.
