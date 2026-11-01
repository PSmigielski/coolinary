# KIT-01-04 Stream consumer with idempotent handlers

| | |
|---|---|
| Service | kitchen (shared platform code; copy to planner and offers) |
| Epic | [KIT-01 Platform](README.md) |
| Release | R1 |
| Size | L |
| Status | todo |
| Depends on | KIT-01-02 |
| Unblocks | KIT-01-05, KIT-13 (household lifecycle) |

## Goal

Kitchen handles events from other services (first `events.identity`) so that each event's effect is applied
exactly once in kitchen's database, in the order it was published.

## Business rules

1. Configuration lists the streams kitchen reads: `EVENT_STREAMS=events.identity`. The consumer group is the service
   name (`kitchen`); the consumer name is the container hostname.
2. On start the group is created with `XGROUP CREATE <stream> kitchen 0 MKSTREAM` if it doesn't exist, so a new
   service processes everything still retained (up to 7 days).
3. One loop per stream reads `XREADGROUP GROUP kitchen <consumer> COUNT 50 BLOCK 5000 STREAMS <stream> >` and
   processes entries one at a time in stream order. Per-household order therefore matches publish order. Each entry
   is parsed from its `envelope` field; the `type` field is read first so unhandled types are skipped without
   parsing (entry layout: KIT-01-02).
4. Handlers are registered per event `type` with the major `version` they understand:
   `consumer.Handle("household.plan_changed", 1, fn)`. An entry whose type has no handler is acknowledged and
   counted as `skipped`. An entry with a higher `version` than the handler knows is not acknowledged; it is retried
   and ends as a dead letter (KIT-01-05), so the incompatibility is visible.
5. Each handler runs inside one Postgres transaction that also inserts `(event_id, stream, type)` into
   `processed_events`. If the insert hits the primary key, the event was already applied: the transaction rolls back
   and the entry is acknowledged.
6. `XACK` is sent only after the transaction commits. A crash before `XACK` re-delivers the entry, which step 5
   turns into a no-op.
7. A handler gets at most 10 seconds (context deadline). A handler error or timeout leaves the entry pending.
8. Every 60 seconds the loop claims entries pending longer than 60 seconds with `XAUTOCLAIM` and retries them.
   After 10 deliveries an entry becomes a dead letter (KIT-01-05) and is acknowledged.
9. Handlers that call other services (e.g. notifier) must be idempotent themselves (pass the event id as the
   request's idempotency key where the callee supports one).
10. `processed_events` rows older than 14 days are deleted daily (the stream keeps only 7 days, so older ids can't
    come back).
11. Shutdown: stop reading, finish the current handler, then return.

## Contract

No proto change.

## Data

`kitchen/migrations/000NN_processed_events.sql`

```sql
CREATE TABLE processed_events (
    event_id      uuid PRIMARY KEY,
    stream        text NOT NULL,
    type          text NOT NULL,
    processed_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX processed_events_age ON processed_events (processed_at);
```

## Events

Consumes the streams in `EVENT_STREAMS`. Handlers are added by the tasks that need them (KIT-13, KIT-11-01).

## Errors

None returned to callers. Handler errors are logged with `event_id`, `type` and `error` (no payload content).

## Acceptance criteria

- [ ] The same entry delivered twice changes the database once.
- [ ] A handler that fails 3 times then succeeds leaves exactly one effect and one `processed_events` row.
- [ ] Entries for one household are applied in publish order across a restart.
- [ ] Unknown types are acknowledged and counted in `events_consumed_total{stream,type,result="skipped"}`.
- [ ] Metrics `events_consumed_total{stream,type,result=ok|duplicate|skipped|error|dead}` and
      `events_handler_duration_seconds{type}`.
- [ ] Integration test against the dev Redis using a fake producer.
