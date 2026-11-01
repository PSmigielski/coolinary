# OFF-01-02 Event consumer group and idempotent handlers

| | |
|---|---|
| Service | offers |
| Epic | [OFF-01 Platform](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-01-01, KIT-01 (consumer pattern) |
| Unblocks | OFF-01-03, OFF-04-04, OFF-11-02, OFF-14-01, OFF-14-02 |

## Goal

offers reacts to changes in other services (plan changes, deletions, merged ingredients, scanned receipts) exactly
once in effect, even though Redis Streams deliver at least once.

## Business rules

1. offers reads two streams with consumer group `offers`: `events.identity` and `events.kitchen`. Consumer name is
   the container hostname.
2. On first start the group is created at `$` (only new events): a fresh offers database has nothing to catch
   up on; household plans default to Free (OFF-01-03).
3. Each entry's `envelope` field is parsed (the `type` field lets unknown types be skipped unparsed; KIT-01-02).
   Each event is handled in one Postgres transaction that also inserts its `id` into `processed_events`. If the id is
   already there, the handler is skipped and the event is acknowledged.
4. Unknown event types are acknowledged and ignored (logged at `debug` with type only).
5. A handler error leaves the event unacknowledged. Pending events older than 60 s are re-claimed (`XAUTOCLAIM`).
   After 10 failed attempts the event is copied to `dead_events` with the error text, acknowledged, and a founder
   alert is sent through notifier `SendAlert` with key `offers.dead_event.<type>` (NTF-05).
6. Events of a newer `version` than the handler knows are treated like rule 5 errors (so nothing is half-handled).
7. `processed_events` rows older than 30 days are deleted daily.
8. Handlers registered by later tasks: `household.deleted`, `account.deleted` (OFF-14-01), `catalogue.ingredient_merged` (OFF-14-02), `receipt.scanned` (OFF-11-02, R3).

## Contract

No gRPC change.

## Data

`offers/migrations/000NN_consumer.sql`

```sql
CREATE TABLE processed_events (
    id           uuid PRIMARY KEY,
    type         text NOT NULL,
    processed_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE dead_events (
    id         uuid PRIMARY KEY,
    stream     text NOT NULL,
    type       text NOT NULL,
    envelope   jsonb NOT NULL,
    error      text NOT NULL,
    attempts   int NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
```

## Errors

| Case | Behaviour |
|---|---|
| Redis down | Consumer retries with backoff up to 30 s; health stays OK only while `redis-state` ping works (existing check) |
| Handler fails 10 times | Event moved to `dead_events`, alert sent, stream moves on |

## Acceptance criteria

- [ ] The same `catalogue.ingredient_merged` event delivered twice rewrites the offer references once and leaves one
      `processed_events` row.
- [ ] A handler that panics on one event does not stop later events; after 10 attempts the event is in
      `dead_events` and one alert reached Mailpit.
- [ ] `events_consumed_total{stream,type,result}` is exported with kitchen's names.
