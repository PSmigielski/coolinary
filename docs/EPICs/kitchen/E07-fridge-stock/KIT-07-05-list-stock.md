# KIT-07-05 ListStock with expiry status and totals

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-07-01, KIT-02 (unit conversion library, KIT-02-02), KIT-10-01 (timezone) |
| Unblocks | GW-06 (fridge list), PLN-01 (stock snapshot), PLN-06, PLN-15 |

## Goal

The app shows the fridge with what to use first on top; planner gets the same stock with per-ingredient totals in
one call, and can skip the call cheaply when nothing changed.

## Business rules

1. Returns every item of the household (at most 1,000; no paging).
2. **Expiry status** per item, from the household's local date: `expired` (date < today), `today`, `tomorrow`,
   `soon` (2 or 3 days ahead), `ok` (later), `no_date`.
3. **Order:** expired, today, tomorrow, soon, ok, no date; inside a group by date ascending, then by name
   (case-insensitive, household language collation), then by `added_at`.
4. **Filters**, all optional and combined with AND: storage places; statuses; ingredient ids; `matched_only`
   (planner sets it, so unmatched items are left out).
5. **Cheap re-check:** when `if_stock_version` equals the current `stock_version`, the response has
   `not_modified = true` and no items.
6. **Totals** (when `with_totals = true`), one entry per ingredient over matched items:
   1. `total_grams` = sum of all items converted to grams (KIT-02-02 `ToGrams`: ml via density, pieces via typical weight);
      empty when any item of that ingredient cannot be converted;
   2. `by_unit` = the raw sums per unit (g, ml, pieces), always present;
   3. `earliest_expires_on` = the earliest date among the ingredient's items (no date if none has one);
   4. `in_freezer_only` = true when every item of that ingredient is in the freezer;
   5. `includes_expired` = true when any counted item has status `expired`.
7. The response always carries `stock_version` and the household's `today` (so callers use the same date).
8. Planner calls this with the household id it is planning for; the gateway with the caller's household.
9. Decided (Q103): items past their date stay in the list and in the totals. Planner counts them as at home and
   marks the dish or plan line "check it first" (PLN-06-01, PLN-08); the status `expired` and `includes_expired`
   tell it which. Kitchen never drops them on its own.

## Contract

`contracts/coolinary/kitchen/v1/stock.proto`

```proto
rpc ListStock(ListStockRequest) returns (ListStockResponse);

message ListStockRequest {
  string household_id = 1;
  repeated StoragePlace storage = 2;
  repeated ExpiryStatus statuses = 3;
  repeated string ingredient_ids = 4;
  bool matched_only = 5;
  bool with_totals = 6;
  int64 if_stock_version = 7;     // 0 = always return items
}

message IngredientStockTotal {
  string ingredient_id = 1;
  optional double total_grams = 2;
  repeated coolinary.common.v1.Quantity by_unit = 3;
  string earliest_expires_on = 4;
  bool in_freezer_only = 5;
  bool includes_expired = 6;
}

message ListStockResponse {
  repeated StockItem items = 1;
  repeated IngredientStockTotal totals = 2;
  int64 stock_version = 3;
  bool not_modified = 4;
  string today = 5;               // YYYY-MM-DD in the household's timezone
}
```

## Data

No schema change (uses `stock_items_household_expiry`).

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| Unknown enum values in filters | `INVALID_ARGUMENT` | |
| Household with no stock | — | empty list, `stock_version = 0` |

## Acceptance criteria

- [ ] With items dated yesterday, today, tomorrow, +3 days, +10 days and one without a date, the order is exactly
      that, and the statuses are `expired`, `today`, `tomorrow`, `soon`, `ok`, `no_date`.
- [ ] `if_stock_version` equal to the current version returns `not_modified = true`; after one add it returns items.
- [ ] Two eggs items (6 pieces and 4 pieces, typical weight 60 g) give `total_grams = 600` and `by_unit` 10 pieces.
- [ ] An item of an ingredient with no typical weight in pieces makes `total_grams` empty for that ingredient only.
- [ ] Milk dated yesterday is returned with status `expired`, counts in the milk total, and the total has
      `includes_expired = true` (Q103).
- [ ] p95 latency below 50 ms for a household with 1,000 items on the dev machine.
