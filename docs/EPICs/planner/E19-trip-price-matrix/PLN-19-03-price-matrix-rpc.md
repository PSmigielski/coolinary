# PLN-19-03 Price matrix RPC and shop totals

| | |
|---|---|
| Service | planner |
| Epic | [PLN-19 Trip price matrix](README.md) |
| Release | R2 |
| Size | M |
| Depends on | PLN-19-02, PLN-15-01 |
| Unblocks | GW-10, PLN-20-03 |

## Goal

Return the full matrix and per-shop totals so the phone can find the cheapest set of 1–3 shops within the walking
limit, apply the minimum saving per extra shop, and show up to two alternatives (features.md "How the plan is found").

## Business rules

1. Response: shops (id, chain, name, the position the phone sent), lines (id, ingredient id, name, amount bought),
   cells (line × shop: cost, source, about, offer id, condition flags, or unpriced), and per shop the total of its
   priced cells ("whole list in this shop"), the number of unpriced cells, and `about` when any cell is about.
2. Totals are sums of the rounded cells, so the phone's per-set sums match the server's numbers to the cent.
3. The matrix is for the list as it is now; it carries the list version so the phone can refetch after a change.
4. Planner does not pick shops, compute sets, apply the minimum saving or rank alternatives; the phone does
   (decision 7).
5. Plus only; validated shops per PLN-19-01.
6. Responses are not stored; each call recomputes (cache in PLN-19-02).
7. A list with no open ingredient lines returns empty lines and zero totals.

## Contract

`contracts/coolinary/planner/v1/prices.proto`

```proto
rpc GetTripPriceMatrix(GetTripPriceMatrixRequest) returns (TripPriceMatrix);

message GetTripPriceMatrixRequest { string household_id = 1; string viewer_member_id = 2; repeated string shop_ids = 3; }
message TripPriceMatrix {
  int64 list_version = 1;
  string currency = 2;
  repeated MatrixShop shops = 3;
  repeated MatrixLine lines = 4;
  repeated MatrixCell cells = 5;          // row-major: lines × shops
  repeated string ignored_shop_ids = 6;
  repeated string unpriced_line_ids = 7;  // free text, or no price anywhere
}
message MatrixShop { string shop_id = 1; string chain_name = 2; string shop_name = 3; int32 position = 4; coolinary.common.v1.Money total = 5; int32 unpriced_cells = 6; bool about = 7; }
message MatrixLine { string line_id = 1; string ingredient_id = 2; string name = 3; coolinary.common.v1.Quantity amount_bought = 4; }
message MatrixCell {
  string line_id = 1; string shop_id = 2;
  optional coolinary.common.v1.Money cost = 3;   // absent = unpriced
  PriceSource source = 4; bool about = 5;
  string offer_id = 6; repeated string condition_flags = 7;
}
enum PriceSource { PRICE_SOURCE_UNSPECIFIED = 0; PRICE_SOURCE_FLYER = 1; PRICE_SOURCE_RECEIPT = 2; PRICE_SOURCE_REFERENCE = 3; }
```

## Data

No schema change.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Free household | `FAILED_PRECONDITION` `plus_required` |
| No valid shop | `FAILED_PRECONDITION` `no_shops` |
| More than 10 shops | `INVALID_ARGUMENT` |
| offers down | `UNAVAILABLE` |

## Acceptance criteria

- [ ] Worked example (features.md "Planning the shopping trip"; prices given as flyer offers in a test double of
      `GetPrices`): totals Lidl 10.65, Kaufland 11.25, Aldi 12.11; taking each line's cheaper cell of Lidl and
      Kaufland gives 3.89 + 1.79 + 0.99 + 1.59 + 0.89 = 9.15 (Verification 13). The phone-side search, fed this
      matrix, returns Lidl only at a €2 minimum and Lidl then Kaufland at €1.
- [ ] A reference-priced line has the same cost in every shop, so it never changes the winner.
- [ ] A line priced nowhere is in `unpriced_line_ids` and in no total.
