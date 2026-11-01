# KIT-07-08 stock.changed event and live notice

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 (event); R2 (live notice for shared households) |
| Size | S |
| Status | todo |
| Depends on | KIT-01 (outbox, Redis Streams publisher), KIT-07-03, KIT-07-04, KIT-07-07 |
| Unblocks | PLN-01 (stock snapshot cache), GW-09 (live updates), MOB-05 |

## Goal

Planner learns that a household's stock changed without polling, and in a shared household every phone sees the
other's change within a second.

## Business rules

1. Every transaction that bumps `stock_version` writes exactly one `stock.changed` row to the outbox (KIT-01) in the
   same transaction. Payload: `household_id`, `stock_version`.
2. The outbox relay publishes it to `events.kitchen`. Consumers treat it as "your cache for this household is stale";
   they must not expect one event per item.
3. **Live notice (R2):** after the transaction commits, kitchen publishes to Redis pub/sub on `redis-state`, channel
   `live.household.<household_id>`, the JSON `{"kind":"stock.changed","household_id":"...","version":N}` (the README "Live notices" shape; `version` is the household's stock version). It is best-effort: a
   failed publish is logged (no personal data) and not retried; clients resync from `ListStock` on reconnect.
4. Live notices are published for every household (the cost of a publish to a channel nobody listens to is
   negligible), so no plan check is needed.
5. Background changes (KIT-07-09 linking, KIT-13 moves and merges) publish both the same way.

## Contract

No RPC. Event in docs/EPICs/README.md (`stock.changed`), payload extended with `stock_version`.

## Data

No schema change (outbox table from KIT-01).

## Events

| Event | Stream / channel | Payload |
|---|---|---|
| `stock.changed` | `events.kitchen` | `household_id`, `stock_version` |
| live notice | pub/sub `live.household.<id>` | `{"kind":"stock.changed","household_id":"...","version":N}` |

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | Behaviour |
|---|---|
| Redis pub/sub down | stock change still commits; notice skipped; metric `kitchen_live_publish_failures_total` + 1 |
| Outbox relay down | events wait in the outbox and go out in order when it recovers (KIT-01) |

## Acceptance criteria

- [ ] Adding 50 items in one call produces one `stock.changed` with the new `stock_version`.
- [ ] A rolled-back transaction produces no event and no live notice.
- [ ] Two subscribers on `live.household.<id>` receive the notice within 1 s of an edit.
- [ ] Stopping Redis pub/sub does not make the edit fail.
