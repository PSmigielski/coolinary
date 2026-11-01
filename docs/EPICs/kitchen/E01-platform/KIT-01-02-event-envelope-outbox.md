# KIT-01-02 Event envelope and outbox writer

| | |
|---|---|
| Service | kitchen (shared platform code; copy to planner and offers) |
| Epic | [KIT-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | — |
| Unblocks | KIT-01-03, KIT-01-04, KIT-02-09, KIT-06-10 |

## Goal

A handler can record an event in the same database transaction as the change it describes, so an event is
published if and only if the change commits.

## Business rules

1. Envelope (README "Events"): `id` (UUIDv7 string), `type` (e.g. `catalogue.ingredient_merged`), `version`
   (integer, starts at 1, raised only for breaking payload changes), `occurred_at` (RFC 3339 UTC, millisecond
   precision), `producer` (`kitchen`), `household_id` (UUID or empty for catalogue-wide events), `payload` (object).
2. Package `internal/platform/events` exposes `events.Add(ctx, q *store.Queries, e Event) error`, which inserts one
   outbox row using the caller's transaction (`q.WithTx(tx)`). Calling it outside a transaction is a programming
   error and is caught by a test helper, not at runtime.
3. `type` must match `^[a-z]+(\.[a-z_]+)+$` and be registered in a Go constant list (`events.Types`); unknown types
   fail with an error before insert.
4. Payloads contain ids, codes, dates and counts only. No names, e-mail addresses, free text or health data. Each
   event type has a Go struct; `payload` is its JSON encoding with stable field names in `snake_case`.
5. The payload of one event is at most 64 KB. Larger fan-out is split into several events by the caller (e.g.
   ≤ 500 ingredient ids per `catalogue.nutrition_changed`).
6. `id` is generated in Go (UUIDv7) so it sorts by creation time; the relay publishes in `id` order.

## Contract

No proto change.

## Data

`kitchen/migrations/000NN_outbox.sql`

```sql
CREATE TABLE outbox (
    id            uuid PRIMARY KEY,               -- UUIDv7, also the envelope id
    type          text NOT NULL,
    version       integer NOT NULL DEFAULT 1,
    household_id  uuid,
    payload       jsonb NOT NULL,
    occurred_at   timestamptz NOT NULL,
    published_at  timestamptz,
    attempts      integer NOT NULL DEFAULT 0,
    last_error    text
);
CREATE INDEX outbox_unpublished ON outbox (id) WHERE published_at IS NULL;
```

## Events

Defines the envelope used by every kitchen event in README "Events".

## Errors

Internal only: `events.Add` returns an error that the handler turns into `INTERNAL` and the transaction rolls back.

## Acceptance criteria

- [ ] A test handler that writes a business row and an event, then rolls back, leaves neither row.
- [ ] Ids generated in one second sort in creation order.
- [ ] An unregistered type or a payload over 64 KB is rejected before insert.
- [ ] Envelope JSON produced for a sample event matches a golden file.

## Notes

Decided (Q015): each stream entry has two fields, `type` (for metrics and skipping) and `envelope` (the full JSON).
The identity (Bun) publisher uses the same layout (README "Events").
