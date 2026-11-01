# PLN-01-04 Stock snapshot cache

| | |
|---|---|
| Service | planner |
| Epic | [PLN-01 Platform](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | PLN-01-01, PLN-01-02, KIT-07 (ListStock, stock.changed) |
| Unblocks | PLN-06, PLN-07-04, PLN-08, PLN-15 |

## Goal

Ranking, plans and the shopping list read the household's stock many times per minute. A short-lived snapshot keeps
kitchen calls low while staying fresh.

## Business rules

1. `stock.Get(household_id)` returns the items from `KIT-07 ListStock`: item id, ingredient id, product EAN if any,
   quantity and unit, storage place, expiry date, added date.
2. Cache key `planner:stock:<household_id>` in redis-cache with kitchen's `stock_version`; TTL 2 minutes.
3. Deleted on `stock.changed` and `household.deleted`. Planner deletes it itself right after its own successful
   `ConsumeStock` / `AddStockItems` calls, so the user sees the change at once.
4. Items are pre-sorted for planning: expiry date ascending (no date last), then added date ascending ("oldest items
   first", business model). Items without an ingredient id (not yet matched) are kept but marked `unmatched` and never
   count as on hand.
5. Single-flight per household as in PLN-01-03.
6. Kitchen down and no cache → UNAVAILABLE; ranking and plans can't run without stock.

## Contract

No planner contract change; consumes `KitchenService.ListStock`.

## Data

No schema change.

## Events

Consumes `stock.changed`, `household.deleted`.

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Kitchen down, no cache | UNAVAILABLE | `stock unavailable` |

## Acceptance criteria

- [ ] Items come back ordered by expiry then added date; items with no expiry date come last.
- [ ] An unmatched item is present in the snapshot but `OnHand(ingredient)` ignores it.
- [ ] After planner's own `ConsumeStock`, the next suggestion call reflects the reduced stock without waiting for the event.
