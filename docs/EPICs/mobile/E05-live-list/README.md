# MOB-05 Live shopping list client

Goal: phones of a shared household (R2, Plus) see each other's list changes within seconds through one WebSocket,
and never show a wrong list after missing a message.

## Flow

```
foreground + online + ≥ 2 accounts ──> wss://api…/v1/live (JWT) ──> hello
       every 25 s ping / pong;  re-auth 120 s before JWT expiry
       list.changed (version) ──> newer than v: GET /v1/list/changes?since=v │ own echo: ignore
       stock.changed (version) ──> invalidate stock
background 10 s ──> close;  foreground ──> reconnect + refetch list
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [MOB-05-01](MOB-05-01-live-connection.md) | Live connection manager | R2 | M | MOB-01-04, MOB-02-04, GW-09 |
| [MOB-05-02](MOB-05-02-live-list-updates.md) | Applying live list changes | R2 | M | MOB-05-01, MOB-03-07, MOB-04-03, PLN-15, GW-09 |
| [MOB-05-03](MOB-05-03-live-household-notices.md) | Live stock notices | R2 | S | MOB-05-01, MOB-03-02, GW-09 |
