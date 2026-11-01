# KIT-01-03 Outbox relay to Redis Streams

| | |
|---|---|
| Service | kitchen (shared platform code; copy to planner and offers) |
| Epic | [KIT-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-01-02 |
| Unblocks | KIT-01-05, every kitchen event (planner and offers consume them) |

## Goal

Committed outbox rows reach the `events.kitchen` stream within about a second, in order, at least once, even when
Redis or the service restarts.

## Business rules

1. The relay runs as one `g.Go(...)` loop in `main.go` and stops on shutdown after finishing the current batch.
2. It wakes on a Postgres `NOTIFY outbox` sent by an `AFTER INSERT` trigger, and at the latest every 1 second.
3. Each pass: `SELECT … FROM outbox WHERE published_at IS NULL ORDER BY id LIMIT 100 FOR UPDATE SKIP LOCKED`, then
   for each row in order `XADD events.kitchen MINID ~ <now − 7 days in ms> * type <type> envelope <json>`, then
   `UPDATE outbox SET published_at = now()` for the rows sent, all in one transaction.
4. If `XADD` fails, the pass stops at that row (later rows wait, keeping order), `attempts` and `last_error` are set,
   and the next pass starts after a backoff of 1 s, 2 s, 4 s … capped at 30 s.
5. A crash between `XADD` and the commit publishes the same event twice; this is allowed (at-least-once) because
   consumers deduplicate by envelope `id` (KIT-01-04).
6. Published rows older than 7 days are deleted by the same loop once an hour, at most 10,000 rows per run.
7. The stream is trimmed by time (7 days) on every `XADD` with approximate trimming (`MINID ~`).
8. `redis-state` is used for streams (never evicted, AOF). `.claude/rules/go.md` currently says redis-state is for
   "quotas, rate limits and pub/sub only"; add "event streams" in all three Go services.

## Contract

No proto change.

## Data

`kitchen/migrations/000NN_outbox_notify.sql`

```sql
-- +goose StatementBegin
CREATE FUNCTION outbox_notify() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN PERFORM pg_notify('outbox', ''); RETURN NULL; END $$;
-- +goose StatementEnd
CREATE TRIGGER outbox_notify AFTER INSERT ON outbox FOR EACH STATEMENT EXECUTE FUNCTION outbox_notify();
```

## Errors

None returned to callers; failures show in metrics and logs (`error` field, no payload content).

## Acceptance criteria

- [ ] 1,000 events written by concurrent transactions appear in `events.kitchen` in `id` order.
- [ ] With Redis stopped for 2 minutes, events written meanwhile are published in order after it returns; none lost.
- [ ] Killing the process between `XADD` and commit produces a duplicate entry with the same envelope `id` and no
      gap.
- [ ] Metrics `events_published_total{type}`, `outbox_pending` (gauge) and `outbox_oldest_age_seconds` (gauge) are
      exported on `/metrics`.
- [ ] Median latency from commit to `XADD` under 1 s in the dev stack.
