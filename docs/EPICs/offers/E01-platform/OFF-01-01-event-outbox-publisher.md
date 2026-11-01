# OFF-01-01 Event outbox and Redis Streams publisher

| | |
|---|---|
| Service | offers |
| Epic | [OFF-01 Platform](README.md) |
| Release | R2 |
| Size | M |
| Depends on | KIT-01 (outbox and publisher pattern in Go) |
| Unblocks | OFF-07-07, OFF-08-05, OFF-10-04 |

## Goal

offers can publish `offers.published` and `offers.expired` reliably: an event is written in the same database
transaction as the change it describes, so no change is lost and no event describes a change that was rolled back.

## Business rules

1. Every event is first inserted into `outbox` inside the transaction that makes the change. Code may not publish
   directly to Redis.
2. A relay loop reads unsent rows in `id` order, `XADD`s them to the stream `events.offers` on `redis-state` with two fields, `type` and `envelope` (KIT-01-02), then sets
   `sent_at`. It reads at most 500 rows per round and sleeps 1 s when there was nothing to send.
3. The envelope is exactly the one in docs/EPICs/README.md "Events": `id` (UUIDv7, generated at insert), `type`,
   `version` (starts at 1), `occurred_at` (UTC, RFC 3339), `producer` = `offers`, `household_id` (empty string for
   events that concern everyone), `payload` (JSON object).
4. `XADD` trims by time like kitchen and identity: `MINID ~ <now − 7 days in ms>` (docs/EPICs/README.md "Events").
5. Only one relay runs at a time: it holds the Postgres advisory lock `hashtext('offers.outbox')`. If the lock is held,
   the loop waits and retries every 5 s.
6. If Redis is down, rows stay unsent and the relay retries with backoff 1 s, 2 s, 4 s … capped at 30 s. Business
   writes never fail because Redis is down.
7. Sent rows are deleted after 7 days by a daily cleanup.
8. Payloads carry ids only: no product text typed by users, no postcode, no user names.

## Contract

No gRPC change. Events are listed in docs/EPICs/README.md. Payload of each offers event is defined where it is
produced (OFF-07-07, OFF-08-05, OFF-10-04).

## Data

`offers/migrations/000NN_outbox.sql`

```sql
CREATE TABLE outbox (
    id           uuid PRIMARY KEY,               -- UUIDv7 = event id
    type         text NOT NULL,
    version      smallint NOT NULL DEFAULT 1,
    household_id text NOT NULL DEFAULT '',
    payload      jsonb NOT NULL,
    occurred_at  timestamptz NOT NULL DEFAULT now(),
    sent_at      timestamptz
);
CREATE INDEX outbox_unsent ON outbox (id) WHERE sent_at IS NULL;
```

## Errors

| Case | Behaviour |
|---|---|
| Redis unreachable | Rows stay unsent; `events_outbox_unsent` gauge rises; no error to callers |
| Payload not valid JSON | Insert fails; the business transaction rolls back (bug, `INTERNAL`) |

## Acceptance criteria

- [ ] A transaction that inserts an offer and an outbox row, then rolls back, produces no stream entry.
- [ ] With Redis stopped, 100 events are written; after Redis restarts all 100 appear in `events.offers` in order,
      each once.
- [ ] Two offers containers started by hand (dev) never both relay: only the lock holder sends.
- [ ] Metrics `events_published_total{type}` and `events_outbox_unsent` exist with the same names as in kitchen.
- [ ] The relay code is a copy of kitchen's (KIT-01) with only the stream name changed.
