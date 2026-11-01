# MOB-03 Query keys, hooks and cache rules

Goal: one typed hook per gateway route, grouped by feature area and release, with household-scoped keys, a single
invalidation map, optimistic updates for stock and list, and an allow-list deciding what may be stored on the phone.

## Flow

```
screen ──> useX() hook ──> keys.x(householdId) ──> $api (MOB-02) ──> gateway
                 │                                   │
     offline-capable mutation ──> outbox (MOB-04)    └─> response ──> cache (persist only if allow-listed)
change ──> invalidate(kind) ──> stale queries refetch (now, or on reconnect)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [MOB-03-01](MOB-03-01-query-keys-cache-rules.md) | Query keys and cache rules | R1 | M | MOB-01-02, MOB-02-02 |
| [MOB-03-02](MOB-03-02-invalidation-map.md) | Invalidation map | R1 | S | MOB-03-01 |
| [MOB-03-03](MOB-03-03-account-household-hooks.md) | Account and household hooks | R1–R2 | M | MOB-03-01, MOB-03-02, GW-02, GW-04 |
| [MOB-03-04](MOB-03-04-profile-catalogue-hooks.md) | Kitchen profile and catalogue hooks | R1–R2 | M | MOB-03-01, MOB-03-02, MOB-04-01, GW-05 |
| [MOB-03-05](MOB-03-05-fridge-barcode-intake-hooks.md) | Fridge, barcode and intake hooks | R1–R3 | L | MOB-03-01, MOB-03-02, MOB-04-01, MOB-04-04, MOB-07-02, GW-06 |
| [MOB-03-06](MOB-03-06-library-plans-cooking-hooks.md) | Library, plans and cooking hooks | R1–R3 | L | MOB-03-01, MOB-03-02, MOB-02-05, MOB-04-01, GW-07 |
| [MOB-03-07](MOB-03-07-shopping-list-hooks.md) | Shopping list hooks | R1–R2 | M | MOB-03-01, MOB-03-02, MOB-04-01, MOB-04-04, GW-08 |
| [MOB-03-08](MOB-03-08-offers-shops-price-matrix-hooks.md) | Shops, offers, flyers and price matrix hooks | R2–R3 | M | MOB-03-01, MOB-03-02, MOB-07-03, GW-10 |
| [MOB-03-09](MOB-03-09-notifications-savings-hooks.md) | Notification preference and savings hooks | R1–R3 | S | MOB-03-01, MOB-03-02, GW-11, GW-12 |
