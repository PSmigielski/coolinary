# GW-09-02 Pub/sub relay and versioned change notices

| | |
|---|---|
| Service | gateway |
| Epic | [GW-09 Live shopping list](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | GW-09-01, PLN-15 (`list.changed`, `GetListChanges`), KIT-07 (`stock.changed`) |
| Unblocks | MOB-05 |

## Goal

Change notices from planner (shopping list) and kitchen (stock) reach every open socket of the household, in order,
and a client that missed any of them knows it must refetch.

## Business rules

1. One dedicated Redis connection to `redis-state` per gateway process in subscriber mode (Bun `RedisClient`
   subscribe). The gateway `SUBSCRIBE`s `live.household.<id>` when the first socket of that household is ready and
   `UNSUBSCRIBE`s when the last one closes.
2. Notices on the channel (docs/EPICs/README.md "Live notices"): JSON with `kind` and a version. Two kinds exist:
   `list.changed` (planner, PLN-15-09) and `stock.changed` (kitchen, KIT-07-08). Any other `kind`, or a message
   without a numeric version, is dropped with a debug log (forward-compatible).
3. The gateway sends every socket of the household, the author's included, exactly the README shape:
   `{"kind": "list.changed" | "stock.changed", "household_id": "<id>", "version": N}`. `household_id` is taken from
   the channel name; `version` from the notice's `version` field (planner and kitchen publish the README shape).
   Other fields are not forwarded.
   Control messages (`ready`, `pong`, `degraded`, GW-09-01) keep their `type` key; a message with `kind` is always a
   change notice.
4. **Gap rule for clients (documented here, implemented in MOB-05):** versions are counted per `kind`. After
   `ready`, and whenever a notice's `version` is not exactly last + 1 for its kind, the client catches up:
   `list.changed` → `GET /v1/list/changes?since=<last version>` (planner `GetListChanges`, which may answer with a
   full resync, GW-08-01); `stock.changed` → `GET /v1/stock`. The gateway keeps no history and does not replay.
5. **Redis lost:** every socket gets `{"type":"degraded"}`; the gateway retries the subscription every 5 s; on
   success it sends `{"type":"ready", ...}` again, which makes clients catch up.
6. Several gateway processes each subscribe on their own; Redis pub/sub fans out to all of them.
7. Notices carry ids and versions only (no names, no quantities), so they never leak personal data into logs.
8. Throughput limit: at most 20 notices per second per household and kind are forwarded; extra ones are collapsed
   into one notice with the latest `version` (clients then see a gap and catch up).

## Contract

Channel `live.household.<householdId>` on `redis-state` (published by planner and kitchen). Socket message:

```json
{ "kind": "list.changed", "household_id": "0192...", "version": 43 }
{ "kind": "stock.changed", "household_id": "0192...", "version": 118 }
```

## Data

No storage (pub/sub only).

## Errors

| Case | Behaviour |
|---|---|
| Malformed message, unknown `kind`, or no version | dropped, metric `gateway_ws_dropped_total{reason}` |
| Redis subscription lost | `degraded` to clients, retry every 5 s |

## Acceptance criteria

- [ ] Two sockets of one household: a `PATCH /v1/list/lines/{id}` from phone A produces
      `{"kind":"list.changed","household_id":…,"version":N}` with the new version on both sockets within 1 s
      (`make dev`, planner publishing).
- [ ] A stock edit produces `{"kind":"stock.changed","household_id":…,"version":N}` on both sockets.
- [ ] A socket of another household receives nothing.
- [ ] A notice with `kind: "household.changed"` (or any unknown kind) is dropped, not forwarded.
- [ ] Restarting `redis-state` sends `degraded`, then `ready` once it is back.
- [ ] 100 list notices in one second arrive as at most 20 messages, the last one with the latest version.
