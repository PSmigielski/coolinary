# NTF-01-01 Redis Streams consumer (Bun)

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-01 Events consumer](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | IDN-01 (outbox publisher, envelope and stream entry format on `events.identity`) |
| Unblocks | NTF-01-02, NTF-01-03, NTF-01-04 |

## Goal

notifier reads identity's events reliably, exactly once in effect, so it can keep its household mirror and delete a
user's data when the account is deleted. This is the first Bun consumer; other Bun services copy it.

## Business rules

1. notifier reads one stream, `events.identity`, with consumer group `notifier`. The group is created at start-up
   with `XGROUP CREATE events.identity notifier 0 MKSTREAM` when it doesn't exist, so the first start replays what
   the stream still holds (about 7 days).
2. Consumer name = `notifier-<hostname>`. One container runs today; the code must still be safe with two.
3. The read loop uses its own Redis connection (a blocking read must not stall the shared `redisState` client):
   `XREADGROUP GROUP notifier <consumer> COUNT 50 BLOCK 5000 STREAMS events.identity >`.
4. Each entry has two fields, `type` and `envelope` (docs/EPICs/README.md "Events"). The `envelope` field is parsed
   as the envelope: `id`, `type`, `version`, `occurred_at`, `producer`, `household_id`, `payload`. An entry without
   an `envelope` field, or whose envelope is not valid JSON or lacks `id` or `type`, is
   moved to `dead_events` with error code `malformed` and acknowledged.
5. Exactly-once effect: in one Postgres transaction, `INSERT INTO processed_events (event_id, type) ... ON CONFLICT
   DO NOTHING`; only when the row was inserted does the handler run, inside the same transaction. After commit the
   entry is acknowledged (`XACK`). A duplicate (row already there) is acknowledged without running the handler.
6. Handlers are registered by `type`. A type with no handler is recorded in `processed_events` with
   `handled = false` and acknowledged (metric result `ignored`).
7. A known type whose `version` is higher than the handler supports is not acknowledged; it is logged at `warn`
   with `type` and `version` and retried like an error (rule 8). Lower or equal versions are handled.
8. When a handler throws, the transaction rolls back and the entry is not acknowledged. Every 30 s the consumer runs
   `XAUTOCLAIM events.identity notifier <consumer> 60000 0 COUNT 50` to retry entries idle for at least 60 s.
9. An entry delivered 10 times (delivery count from `XPENDING`) is copied to `dead_events` (envelope, error code,
   delivery count) and acknowledged, so it stops blocking retries. A founder alert is raised with key
   `notifier.dead_event.<type>` (NTF-05-01; until NTF-05 exists, log at `error`).
10. Gap check at start-up and every hour: if the group's `last-delivered-id` is older than the first entry still in
    the stream (`XINFO STREAM`), events were trimmed before notifier read them. Set gauge `events_consumer_gap` to 1,
    log at `error`, raise alert `notifier.events_gap`. Processing continues.
11. `processed_events` rows older than 14 days are deleted by a daily job at 04:00 UTC (the stream keeps about 7
    days, so older duplicates cannot arrive).
12. On SIGTERM the loop stops reading, waits up to 10 s for running handlers, then exits; unacknowledged entries are
    retried after restart.
13. The consumer does not change `/healthz` (it already checks `redis-state` and Postgres).
14. Logs carry `event_id`, `type`, `version`, `result`, `duration_ms`; never the payload.

## Contract

No gRPC change.

## Data

`notifier/src/db/schema.ts` → `bun run db:generate`.

```sql
CREATE TABLE processed_events (
  event_id     uuid PRIMARY KEY,
  type         text NOT NULL,
  handled      boolean NOT NULL DEFAULT true,
  processed_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX processed_events_at ON processed_events (processed_at);

CREATE TABLE dead_events (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  stream      text NOT NULL,
  entry_id    text NOT NULL,
  event_id    uuid,
  type        text,
  envelope    jsonb,          -- ids only, per the envelope rules
  error_code  text NOT NULL,  -- malformed | handler_error | unsupported_version
  deliveries  integer NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now()
);
```

## Events

Consumes `events.identity`. Handlers come in NTF-01-02, NTF-01-03 and NTF-01-04.

## Errors

| Case | Behaviour |
|---|---|
| Redis down | Loop waits 1 s, 2 s, 4 s … up to 30 s between attempts; logs once per minute |
| Postgres down | Handler fails, entry stays pending, retried (rule 8) |
| Malformed entry | `dead_events` `malformed`, acknowledged |

## Acceptance criteria

- [ ] Publishing the same envelope twice (same `id`) runs the handler once; both entries are acknowledged.
- [ ] A handler that throws 9 times then succeeds is applied once; one that always throws ends in `dead_events`
      after 10 deliveries and the alert is raised.
- [ ] An unknown type is acknowledged and counted as `ignored`.
- [ ] After `XTRIM events.identity MAXLEN 0` with unread entries, the gap gauge reads 1.
- [ ] Metrics: `events_consumed_total{stream,type,result}` with result `ok|duplicate|ignored|error|dead`,
      `events_consumer_pending{stream}`, `events_consumer_lag_seconds{stream}` (now minus `occurred_at` of the last
      handled event).
- [ ] No payload field appears in any log line.

## Notes

- Decided (Q042): no automatic resync after a gap. The risk (notifier down for more than 7 days) is accepted and the
  lag alert is the safeguard.
- Bun's built-in `RedisClient` supports raw commands through `send`; use a second client instance for the
  blocking read.
