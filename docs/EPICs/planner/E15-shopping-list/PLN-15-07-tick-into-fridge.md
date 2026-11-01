# PLN-15-07 Tick into the fridge (and undo)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-15-01, PLN-15-08, KIT-07 (AddStockItems, RemoveStockItems) |
| Unblocks | GW-08, KIT-12 (savings from offer prices) |

## Goal

When an item is ticked in the shop, it moves into the fridge with a default use-by date (features.md "Shopping list"
step 7; §4 "Tick in the shop" is the fastest way to stock the fridge).

## Business rules

1. Ticking an ingredient line calls KIT-07 `AddStockItems` with: ingredient id, quantity = pack size × pack count
   when packs are known (what was actually bought), else the line amount; source `shopping_list`; source ref = line
   id + tick number (`<line id>:<tick n>`, kitchen's idempotency key). kitchen sets the storage place and use-by date from its defaults.
2. If the tick request names the shop (optional `shop_id`, R2) and the line has an offer badge for that shop
   (PLN-17-01), the offer price is passed as price paid and the saving against the reference price is passed too,
   for the savings ledger (KIT-12).
3. The stock item ids returned are stored on the line.
4. Free-text lines are ticked on the list only; nothing is added to the fridge, and the response says
   `not_added_to_fridge`.
5. kitchen down: the tick is saved on the list with flag `fridge_pending`; a retry job (every minute, up to 24 hours)
   adds the stock with the same idempotency key, so nothing is added twice.
6. **Undo:** unticking within 24 hours removes the stock items it created through KIT-07 `RemoveStockItems` with reason
   `mistake`, those item ids and their versions, then returns the line to active. If kitchen reports the items were already used
   or changed, the line is unticked but stock is left as it is (`stock_kept`). After 24 hours, unticking only
   reopens the line.
7. Ticking is idempotent: ticking a ticked line returns it unchanged.
8. Ticking a line removed by someone else in the last 7 days restores it as ticked (the item was bought).
9. One tick is one list change and one live notice.

## Contract

`contracts/coolinary/planner/v1/list.proto`

```proto
rpc TickListLine(TickListLineRequest) returns (TickListLineResponse);
rpc UntickListLine(UntickListLineRequest) returns (ListLine);

message TickListLineRequest {
  string household_id = 1; string member_id = 2; string client_mutation_id = 3;
  string line_id = 4; string shop_id = 5;          // optional, R2
  string ticked_at = 6;                              // device time, for offline ticks
}
message TickListLineResponse { ListLine line = 1; bool added_to_fridge = 2; bool fridge_pending = 3; }
message UntickListLineRequest { string household_id = 1; string member_id = 2; string client_mutation_id = 3; string line_id = 4; }
```

## Data

`planner/migrations/000NN_list_ticks.sql`

```sql
ALTER TABLE list_lines ADD COLUMN tick_count integer NOT NULL DEFAULT 0;
CREATE TABLE list_fridge_outbox (
    line_id      uuid NOT NULL,
    tick_count   integer NOT NULL,
    payload      jsonb NOT NULL,            -- AddStockItems request
    attempts     integer NOT NULL DEFAULT 0,
    next_try_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (line_id, tick_count)
);
```

## Errors

| Case | gRPC status / behaviour |
|---|---|
| Line not found (and not a recent tombstone) | `NOT_FOUND` |
| kitchen down | success with `fridge_pending = true` |

## Acceptance criteria

- [ ] Ticking "rice 1 × 1 kg" adds 1,000 g rice to the pantry with kitchen's default use-by (Verification 12:
      "ticking an item in the shop adds it to the fridge").
- [ ] Replaying the same tick offline twice adds stock once.
- [ ] Unticking within 24 hours removes the rice again; after 24 hours it does not.
- [ ] With kitchen stopped, the tick succeeds; after kitchen returns, the stock appears exactly once.

## Notes

- Uses KIT-07-03 `AddStockItems` (`source_ref` idempotency, `price_paid`, `offer_saving`) and KIT-07-04
  `RemoveStockItems` (reason `mistake`). `RemoveStockItems` is not in the README's RPC table yet.
- Decided (Q129): rule 4 (free-text lines are ticked on the list only; they don't go into the fridge).
