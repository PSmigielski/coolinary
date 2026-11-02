# PLN-01-01 Event consumer (Redis Streams, Go)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Status | done |
| Depends on | KIT-01 (Go consumer pattern) |
| Unblocks | PLN-01-03, PLN-01-04, PLN-04-02, PLN-21 |

## Goal

Planner learns about changes it doesn't own (household plan, members, catalogue merges, nutrition changes, profile
and stock changes, offers) through the Redis Streams events in docs/EPICs/README.md, processed exactly once in effect.

## Business rules

1. Copy the Go consumer from kitchen (KIT-01) into `planner/internal/platform/events`, keeping the same API, envelope
   type and metrics so the copies stay identical (docs/service-conventions.md, "copy the closest service").
   Decided (Q232): where the rules below differ from kitchen's code (group start, read size, 5 deliveries, dead at
   once, Redis retry, cleanup time), these rules win; the API, envelope type and metric names stay kitchen's.
2. Consumer group name `planner`; consumer name = container hostname. Streams read: `events.identity`,
   `events.kitchen`, `events.offers`. The group is created at start-up with `XGROUP CREATE ... $ MKSTREAM` if missing.
3. Read with `XREADGROUP COUNT 100 BLOCK 5000`. Each entry has the fields `type` and `envelope` (KIT-01-02); the
   envelope JSON is parsed from `envelope`. Each message is handled in its own Postgres transaction:
   insert `event_id` into `processed_events`; if the insert conflicts, the event was already handled → `XACK` and skip.
4. Handlers are registered per event `type`. Unknown types are acknowledged and counted, never fail the loop.
   Decided (Q235): counted with kitchen's result label `skipped`.
5. A handler error rolls back the transaction and leaves the message pending. Pending messages older than 60 s are
   reclaimed with `XAUTOCLAIM` and retried. After 5 failed deliveries (from `XPENDING` delivery count) the message is
   copied to `dead_events` with the last error, acknowledged, and an error log line is written.
6. Envelope `version` greater than the handler supports → the message goes to `dead_events` at once (no retries).
7. `processed_events` rows older than 14 days are deleted by a daily job at 03:30 UTC (streams keep about 7 days, so
   older ids can't come back).
8. The consumer stops cleanly on shutdown: it finishes the current message and does not take new ones.
9. Handlers must not call other services inside the transaction; if they need to, they enqueue work (e.g. a row in a
   job table) and return.
10. `/healthz` does not depend on the consumer (Redis-state being down already fails `redis_state`).

## Contract

No gRPC change.

## Data

`planner/migrations/000NN_events.sql`

```sql
CREATE TABLE processed_events (           -- kitchen's columns, so the copied query is the same
    event_id     uuid PRIMARY KEY,
    stream       text NOT NULL,
    type         text NOT NULL,
    processed_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX processed_events_age ON processed_events (processed_at);

-- Decided (Q234): keyed by stream entry, so malformed entries (no JSON, no event id) fit too.
CREATE TABLE dead_events (
    stream      text NOT NULL,
    entry_id    text NOT NULL,
    event_id    uuid,                      -- NULL when the envelope has no valid id
    type        text NOT NULL,
    envelope    text NOT NULL,             -- the entry's envelope field exactly as received
    last_error  text NOT NULL,
    deliveries  int NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (stream, entry_id)
);
```

## Events

Consumes the envelope `id`, `type`, `version`, `occurred_at`, `producer`, `household_id`, `payload`. Handlers are added by
the epics that need them.

## Errors

| Case | Behaviour |
|---|---|
| Redis unreachable | Retry every 2 s with jitter, log once per minute, keep serving gRPC |
| Malformed JSON | Move to `dead_events`, acknowledge |
| Handler error | Leave pending, retry via `XAUTOCLAIM`, dead-letter after 5 deliveries |

## Acceptance criteria

- [ ] Publishing the same envelope twice to `events.kitchen` runs the handler once (`processed_events` has one row).
- [ ] A handler that fails 5 times leaves one row in `dead_events` and the stream has no pending entry for it.
- [ ] Killing planner mid-handler and restarting processes the message again and only once in effect.
- [ ] Metrics `events_consumed_total{stream,type,result}` (`ok`, `duplicate`, `skipped`, `error`, `dead`, Q235) and
      `events_consumer_lag{stream,group}` (the group's lag in entries from `XINFO GROUPS`, refreshed every minute,
      Q233) are exposed on `/metrics`, with the same names as kitchen.
- [ ] Logs carry `event_id` and `type`, never payload contents.

## Notes

Decided (Q232): planner's numbered rules win over kitchen's code where they differ; API, envelope type and metric
names stay kitchen's.
Decided (Q233): the lag metric is Q041's `events_consumer_lag{stream,group}` in entries (as KIT-01-05 defines it), not
a seconds metric.
Decided (Q234): `dead_events` is keyed by `(stream, entry_id)` with a nullable `event_id` and the raw envelope as text.
Decided (Q235): unknown types are counted as `skipped`, kitchen's label.
