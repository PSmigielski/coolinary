# GW-09 Live shopping list (WebSocket)

Release R2 (comes with the shared household).

## Goal

When one household member changes the list, the others see it within a second or two (features.md §2: "The list is
shared live between household members"). Planner publishes change notices on Redis pub/sub channel
`live.household.<id>` (PLN-15); the gateway relays them to every open socket of that household. Notices are small
("list version 43", "stock version 118"); clients catch up with `GET /v1/list/changes?since=` or `GET /v1/stock`
when they miss one.

## Protocol (client ↔ gateway)

```
client: connect wss://api.<domain>/v1/live
client: {"type":"auth","token":"<JWT>"}                       within 5 s
server: {"type":"ready","householdId":"...","heartbeatSeconds":25}
server: {"kind":"list.changed","household_id":"...","version":43}
server: {"kind":"stock.changed","household_id":"...","version":118}
client: {"type":"auth","token":"<new JWT>"}                   before the old one expires
server: {"type":"degraded"}                                    pub/sub lost; poll every 30 s
server: close 4401 (auth), 4403 (household changed), 4429 (too many sockets), 1001 (gateway shutting down)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-09-01](GW-09-01-socket-auth.md) | `/v1/live` socket: auth, token refresh, limits, heartbeat | R2 | M | GW-02-01, GW-03-04 |
| [GW-09-02](GW-09-02-pubsub-relay.md) | Pub/sub relay and versioned change notices | R2 | M | GW-09-01, PLN-15, KIT-07 (live notices) |
