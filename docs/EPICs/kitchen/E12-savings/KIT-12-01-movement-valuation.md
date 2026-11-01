# KIT-12-01 Grams and value of each stock movement

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-12 Savings](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | KIT-07-01 (movement columns), KIT-02-02 (`ToGrams`) |
| Unblocks | KIT-12-02 |

## Goal

Every movement carries its weight in grams and its share of what the item cost, computed once when it happens, so
savings never depend on items that were deleted later.

## Business rules

1. When a movement is written (KIT-07-03, 04, 07, KIT-13-02), kitchen fills:
   1. `grams` = |delta| converted with KIT-02-02 `ToGrams` (ml by density, pieces by piece weight); empty when it
      cannot be converted;
   2. `value_minor` = the item's `price_paid_minor` × |delta| ÷ `initial_quantity`, rounded half up to a whole
      minor unit; empty when the item has no price;
   3. `currency` = the item's currency.
2. For an `added` movement whose item came with an offer saving, the saving is recorded once, in the new column
   `saving_minor` of that movement (the saving that came with the added quantity; merges add their own).
3. **Only known prices count.** Kitchen has no reference prices (prices are owned by offers, and kitchen does not
   call offers). Decided (Q209): items added by typing or barcode without a price add kilograms but no euros;
   euros are counted only where the price paid is known.
4. Merged items (KIT-07-03 rule 6) add their prices, so the share stays right.
5. Values are never recomputed later (a later price change does not rewrite the past).

## Contract

No RPC.

## Data

`kitchen/migrations/000NN_movement_saving.sql`

```sql
-- +goose Up
ALTER TABLE stock_movements ADD COLUMN saving_minor bigint CHECK (saving_minor >= 0);

-- +goose Down
ALTER TABLE stock_movements DROP COLUMN saving_minor;
```

## Errors

None (a failed conversion leaves `grams` empty).

## Acceptance criteria

- [ ] A 1,000 g chicken item paid 15.99 PLN, 400 g cooked → movement `grams = 400`, `value_minor = 640`.
- [ ] 2 of 6 eggs (piece weight 60 g, paid 9.00) used → `grams = 120`, `value_minor = 300`.
- [ ] A ticked list line with a 2.00 offer saving writes `saving_minor = 200` on its `added` movement only.
- [ ] An item without a price produces grams and no value.
