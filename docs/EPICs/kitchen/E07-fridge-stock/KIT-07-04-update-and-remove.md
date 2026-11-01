# KIT-07-04 Edit and remove items (used, thrown away, mistake)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-07-01, KIT-07-02 |
| Unblocks | GW-06 (edit, remove), MOB-04 (offline edits and conflicts), KIT-12-02 |

## Goal

The user can correct an item (amount, place, date, what it is) and take it out, saying whether it was used, thrown
away or added by mistake, so stock stays true and savings are counted honestly. Two phones editing the same item
never overwrite each other silently.

## Business rules

1. **Optimistic concurrency.** Every edit and removal carries the item `version` the client saw. If it no longer
   matches, nothing changes and the call fails with `ABORTED`; the error details carry the current item so the app
   can show it (MOB-04). Decided (README decision 18): `ABORTED`, with
   `x-error-code: version_mismatch`.
2. **Edit** (`UpdateStockItem`) may change, each optional: quantity, storage place, use-by date or "no date",
   ingredient (re-link), name (only for unmatched items or to override the display name).
   1. Quantity up or down writes an `adjusted` movement with the difference; it has no effect on savings.
   2. Storage place applies KIT-07-02 rule 5 to the date.
   3. A typed date sets `expiry_source = user`.
   4. A new ingredient sets `matched = true`, and the name becomes the ingredient's name unless the item has a
      product or the call also sets a name.
   5. After an edit the merge rule (KIT-07-03 rule 6) is not applied; identical items may exist side by side.
3. **Remove** (`RemoveStockItems`) takes 1–100 entries, each with item id, version, reason and an optional
   quantity:
   1. reasons: `used` (eaten or used without the cooking mode), `thrown_away`, `mistake`;
   2. no quantity, or a quantity ≥ the item's quantity: the item is deleted;
   3. a smaller quantity: the item's quantity drops by that amount and the item stays (version + 1);
   4. one movement per entry with the reason (`used`, `thrown_away`, `mistake`) and the amount taken;
   5. all entries succeed or none (one transaction); a version conflict on any entry fails the whole call.
4. Movements for `used` and `thrown_away` feed the savings ledger (KIT-12-02); `mistake` and `adjusted` never do.
5. Any account member of the household may edit and remove any item of the household.
6. Every successful call bumps `stock_version` once and emits one `stock.changed` (KIT-07-08).
7. On a catalogue merge (KIT-02-09), items pointing to the old ingredient are re-pointed inside the merge
   transaction through the `ingredientRefs` registry; KIT-13-09 registers the stock tables there.

## Contract

`contracts/coolinary/kitchen/v1/stock.proto`

```proto
rpc UpdateStockItem(UpdateStockItemRequest) returns (UpdateStockItemResponse);
rpc RemoveStockItems(RemoveStockItemsRequest) returns (RemoveStockItemsResponse);

message UpdateStockItemRequest {
  string household_id = 1;
  string item_id = 2;
  int32 version = 3;
  optional coolinary.common.v1.Quantity quantity = 4;   // same unit as the item, or a convertible one (KIT-02 unit conversion)
  optional StoragePlace storage = 5;
  optional string expires_on = 6;                       // YYYY-MM-DD
  bool clear_expiry = 7;
  optional string ingredient_id = 8;
  optional string name = 9;
}
message UpdateStockItemResponse { StockItem item = 1; int64 stock_version = 2; }

enum RemoveReason {
  REMOVE_REASON_UNSPECIFIED = 0;
  REMOVE_REASON_USED = 1;
  REMOVE_REASON_THROWN_AWAY = 2;
  REMOVE_REASON_MISTAKE = 3;
}
message RemoveEntry {
  string item_id = 1;
  int32 version = 2;
  RemoveReason reason = 3;
  coolinary.common.v1.Quantity quantity = 4;            // empty = the whole item
}
message RemoveStockItemsRequest { string household_id = 1; repeated RemoveEntry entries = 2; }
message RemoveStockItemsResponse {
  repeated StockItem remaining = 1;                     // items that stayed with a smaller quantity
  repeated string deleted_item_ids = 2;
  int64 stock_version = 3;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| Item not in this household | `NOT_FOUND` | |
| Version differs | `ABORTED`, `x-error-code: version_mismatch` | current item in error details |
| Unspecified reason, 0 or > 100 entries, quantity ≤ 0 | `INVALID_ARGUMENT` | |
| Quantity in a unit that cannot be converted to the item's unit | `INVALID_ARGUMENT` | |
| Unknown ingredient | `NOT_FOUND` | |

## Acceptance criteria

- [ ] Two phones load the same yoghurt (version 3); phone A sets 2 pieces (version 4); phone B's edit with
      version 3 fails with `ABORTED` and receives version 4 in the details.
- [ ] Removing 200 g of a 500 g flour item as `used` leaves 300 g and writes a `used` movement of −200 g.
- [ ] Removing three expired items in one call with reason `thrown_away` deletes them and bumps `stock_version` by 1.
- [ ] A `mistake` removal writes a `mistake` movement and changes no savings totals (checked with KIT-12).
- [ ] Moving a `default`-dated item to the freezer changes its date per KIT-07-02; moving a `user`-dated one does not.
