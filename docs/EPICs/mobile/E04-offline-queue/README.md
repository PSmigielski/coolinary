# MOB-04 Offline mutation queue

Goal: everyday edits (fridge, shopping list, cooked it, profile, favourites) work without signal, survive restarts,
reach the server once and in order, and resolve conflicts with a fixed rule set. AI calls are never queued.

## Flow

```
user change ──> allowed kind? ──no──> online call only (offline error)
                    │ yes
                    ▼
        outbox (AsyncStorage, ≤ 500 entries / 1 MB, coalesced)
                    │                         └─> overlay keeps it visible on screens
   online / foreground / every 30 s
                    ▼
        replay one by one, oldest first, same Idempotency-Key
                    ├─ 2xx → remove, invalidate, update cache
                    ├─ network / 5xx → backoff 2 s … 5 min
                    └─ 404 / 409 / 4xx → conflict rules → drop + tell the user when it matters
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [MOB-04-01](MOB-04-01-outbox-storage.md) | Outbox: storage and allowed changes | R1 | M | MOB-01-02, MOB-02-01 |
| [MOB-04-02](MOB-04-02-outbox-replay.md) | Outbox replay | R1 | M | MOB-04-01, MOB-02-02, MOB-02-04, MOB-01-03, MOB-01-04, MOB-03-02 |
| [MOB-04-03](MOB-04-03-conflict-rules.md) | Conflict rules | R1 | M | MOB-04-02, GW-01, GW-06, GW-08 |
| [MOB-04-04](MOB-04-04-pending-overlay.md) | Pending changes overlay | R1 | M | MOB-04-01, MOB-03-01 |
