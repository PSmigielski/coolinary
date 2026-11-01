# KIT-07-03 AddStockItems: add, merge, idempotent adds

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-07-01, KIT-07-02, KIT-02 (`GetIngredients`), KIT-04 (product → ingredient link, pack size) |
| Unblocks | GW-06 (add items), PLN-15 (tick a list line), KIT-07-06, KIT-08, KIT-12 |

## Goal

One RPC adds one item or a whole confirmed batch (unpack, photo, receipt, ticked list lines) to the household's
stock, in one transaction, merging identical items and never adding the same thing twice when a call is retried.

## Business rules

1. A call holds 1–100 items and one `source`. All items are stored or none (one transaction).
2. Each item names what it is with at least one of `ingredient_id`, `product_ean`, `name`:
   1. `product_ean` given: the product must exist (KIT-04/KIT-05); its linked ingredient becomes `ingredient_id`
      unless the caller sent a different one (the caller wins: the user corrected the match);
   2. only `name`: the item is stored "not matched" (no ingredient) and queued for KIT-07-09. Decided (Q079):
      free-text items are allowed and matched later.
3. Quantity: the given quantity; else the product's pack size; else 1 piece. Limits from KIT-07-01 rule 2.
4. Storage place and use-by date come from KIT-07-02 when not given; `no_expiry = true` stores no date
   (`expiry_source = none`).
5. `name`: product name, else ingredient name in the household's language (`x-language` falls back to the profile
   language), else the given text trimmed to 120 characters.
6. **Merge rule.** An item merges into an existing item of the same household when all of these are equal:
   ingredient (both set), product EAN (or both empty), unit, storage place, use-by date (or both without a date).
   Merging adds the quantity to `quantity` and `initial_quantity`, adds price and offer saving when the currency is
   the same (different currency → no merge), raises `version`, and returns `merged = true`. Items in the same call
   merge with each other by the same rule.
   A caller that later undoes an add (planner's untick, PLN-15-07) must remove only the added quantity
   (`RemoveEntry.quantity`, KIT-07-04) when the add was `merged`, never the whole item.
7. **Idempotent adds.** When `source_ref` is set and `(household_id, source, source_ref)` was already added in the
   last 90 days, nothing is added; the response returns the item it went to with `duplicate = true`. Planner sends
   one `source_ref` per tick (PLN-15); gateway idempotency keys (GW-01) cover the other sources.
8. **Limit.** New (not merged) items may bring the household to at most 1,000 items; otherwise the whole call fails
   with `RESOURCE_EXHAUSTED` and the number of free places in the message.
9. Each stored or merged item writes one `added` movement; the call bumps `stock_version` once and emits one
   `stock.changed` (KIT-07-08).
10. `intake_ref` (from KIT-08 proposals) is used only for metrics (KIT-08-01) and receipts (KIT-08-05/06).
11. Prices: `price_paid` and `offer_saving` must be ≥ 0 and share one currency. From a list tick they come from
    planner (offer or price-matrix price); from receipts from the line total; from the app only if the user typed
    a price.
12. Adding a staple ingredient (KIT-10-05) is allowed; such stock is used first by KIT-07-07.
13. Any account member of the household may add, including anonymous accounts. Planner calls on behalf of the
    household that owns the list.

## Contract

`contracts/coolinary/kitchen/v1/stock.proto`

```proto
rpc AddStockItems(AddStockItemsRequest) returns (AddStockItemsResponse);

message NewStockItem {
  string ingredient_id = 1;
  string product_ean = 2;
  string name = 3;                              // required if both ids are empty
  coolinary.common.v1.Quantity quantity = 4;    // g, ml or pieces; empty = pack size or 1 piece
  StoragePlace storage = 5;                     // empty = suggested
  string expires_on = 6;                        // empty = default
  bool no_expiry = 7;
  string pack_code_expires_on = 8;              // from a GS1 2D code
  coolinary.common.v1.Money price_paid = 9;
  coolinary.common.v1.Money offer_saving = 10;
  string source_ref = 11;                       // idempotency per item, e.g. "<list line id>:<tick n>"
}

message AddStockItemsRequest {
  string household_id = 1;
  StockSource source = 2;
  repeated NewStockItem items = 3;              // 1–100
  string intake_ref = 4;                        // KIT-08 proposal or receipt reference
}

message AddedStockItem {
  StockItem item = 1;
  bool merged = 2;
  bool duplicate = 3;
}

message AddStockItemsResponse {
  repeated AddedStockItem items = 1;            // same order as the request
  int64 stock_version = 2;
}
```

## Data

`kitchen/migrations/000NN_stock_add_refs.sql`

```sql
-- +goose Up
CREATE TABLE stock_add_refs (
    household_id  uuid NOT NULL,
    source        smallint NOT NULL,
    source_ref    text NOT NULL CHECK (char_length(source_ref) <= 120),
    stock_item_id uuid NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, source, source_ref)
);
-- Rows older than 90 days are deleted by KIT-13-08.

-- +goose Down
DROP TABLE stock_add_refs;
```

## Events

Publishes `stock.changed` through KIT-07-08.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Behaviour |
|---|---|---|
| 0 or more than 100 items, unknown source | `INVALID_ARGUMENT` | |
| An item with no ingredient, EAN or name | `INVALID_ARGUMENT` | names the item's index |
| Quantity ≤ 0 or above the unit limit; bad date | `INVALID_ARGUMENT` | |
| Prices in two currencies in one item, negative price | `INVALID_ARGUMENT` | |
| Unknown ingredient or EAN | `NOT_FOUND` | names the item's index |
| Household would exceed 1,000 items | `RESOURCE_EXHAUSTED` | "N places left" |

## Acceptance criteria

- [ ] Two cartons of milk 1 l, same date, added one after the other become one item of 2,000 ml; a third with
      another date stays separate.
- [ ] A batch of 30 unpack items where item 17 has an unknown EAN stores nothing and names index 17.
- [ ] A planner retry with the same `source_ref` returns `duplicate = true` and the stock is unchanged; a new tick
      with another `source_ref` adds again.
- [ ] Features Verification 12, last step: a ticked list line arrives as an item with the default use-by date.
- [ ] A free-text item "grandma's jam" is stored with `matched = false` and appears in KIT-07-09's queue.
- [ ] One call with 50 items bumps `stock_version` by exactly 1 and writes 50 `added` movements (fewer if merged
      inside the call).
