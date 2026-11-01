# IDN-01-02 Publisher to `events.identity`

| | |
|---|---|
| Service | identity |
| Epic | [IDN-01 Events: outbox and publisher](README.md) |
| Release | R1 |
| Size | M |
| Depends on | IDN-01-01, NTF-05 (SendAlert) for the stuck-outbox alert |
| Unblocks | KIT-13, PLN-21, OFF-14, NTF-01 (all consume `events.identity`) |

## Goal

Committed outbox rows reach the Redis stream `events.identity` in commit order, at least once, even when Redis or
identity restarts in between.

## Business rules

1. A loop runs inside the identity process. It wakes every 250 ms, and right after any transaction that called
   `emit` (an in-process signal), whichever comes first.
2. Each round, in one Postgres transaction: select up to 100 rows `WHERE published_at IS NULL ORDER BY id
   FOR UPDATE SKIP LOCKED`, publish them one by one, set `published_at = now()` on each published row, commit.
3. Publish = `XADD events.identity MINID ~ <now − 7 days in ms> * type <type> envelope <envelope JSON>`: two
   fields, as in docs/EPICs/README.md "Events" (layout from KIT-01-02, so Bun and Go consumers read the same shape). The stream entry id is
   Redis's own (`*`); the envelope `id` is the idempotency key.
4. Envelope JSON, exactly: `{"id","type","version","occurred_at","producer":"identity","household_id","payload"}`.
   `occurred_at` is RFC 3339 in UTC with milliseconds. `household_id` is a string or null.
5. Order: if publishing a row fails, the round stops at that row (later rows are not published ahead of it),
   `attempts` is incremented, `last_error` stores the error class (no payload), and the next round waits with
   backoff: 1 s, 2 s, 4 s … capped at 30 s, reset after a success.
6. At-least-once: a crash after `XADD` and before commit publishes the same rows again. Consumers drop duplicates by
   envelope `id` (documented in IDN-01-03).
7. Redis down never blocks RPCs: business transactions only write the outbox, which keeps growing until Redis is
   back.
8. Two identity containers at once (rolling deploy) are safe: `SKIP LOCKED` gives each row to one publisher. Order
   across the two may interleave by at most one round; consumers must not rely on order across households.
9. Cleanup: once a day at 04:00 UTC delete rows with `published_at` older than 7 days. Unpublished rows are never
   deleted.
10. Stuck outbox: if the oldest unpublished row is older than 5 minutes, log `error` once a minute and call notifier
    `SendAlert` with dedup key `identity.outbox_stuck` (notifier sends at most one alert per key per hour, NTF-05).
11. Shutdown (SIGTERM): finish the current round, then stop before the gRPC server closes.
12. Metrics (same names in every service; see Notes): `events_published_total{stream,type}` counter,
    `events_outbox_pending` gauge, `events_outbox_oldest_seconds` gauge, `events_publish_errors_total{stream}`.

## Contract

No RPC. Stream `events.identity` on redis-state; each entry has field `type` (the event type) and field `envelope`
(the envelope JSON).

## Data

No schema change (uses `outbox_events` from IDN-01-01).

## Errors

| Case | Behaviour |
|---|---|
| Redis unreachable | Round stops, backoff, `events_publish_errors_total` +1, RPCs unaffected |
| Postgres unreachable | Round skipped; `/healthz` already reports Postgres down |
| Oldest pending > 5 min | `SendAlert` with key `identity.outbox_stuck` |

## Acceptance criteria

- [ ] With redis-state stopped, sign-ups still succeed and `events_outbox_pending` grows; after Redis restarts,
      every event appears in `XRANGE events.identity - +` in id order within 5 seconds.
- [ ] Killing identity between `XADD` and commit (fault-injection test) produces a duplicate entry with the same
      envelope `id`, never a missing one.
- [ ] Entries older than 7 days are trimmed (`XADD … MINID ~`), checked with a test clock.
- [ ] Published rows older than 7 days are deleted by the daily cleanup; unpublished rows are not.
- [ ] No e-mail address or name appears in any stream entry (test scans the payloads of a full sign-up → link →
      delete run).

## Notes

Decided (Q041): every producer uses the four names in rule 12 and every consumer
`events_consumed_total{stream,type,result}` and `events_consumer_lag{stream,group}`; all six are written into
`docs/service-conventions.md` so Go and Bun services match.
