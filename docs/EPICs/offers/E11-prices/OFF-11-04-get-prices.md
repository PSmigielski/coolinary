# OFF-11-04 GetPrices for the trip price matrix

| | |
|---|---|
| Service | offers |
| Epic | [OFF-11 Prices](README.md) |
| Release | R2 (receipt step R3) |
| Size | M |
| Depends on | OFF-11-01, OFF-11-03, OFF-04-01 |
| Unblocks | PLN-19 (trip price matrix) |

## Goal

planner gets, for each list line and each candidate shop, the best known price and where it came from, so the phone
can find the cheapest set of shops (features.md §2 "Planning the shopping trip", Verification 13).

## Business rules

1. Input: household id, shop ids (at most 40, one country), lines (line id, ingredient id, optional wanted amount),
   and the shopping date (default today).
2. For each line × shop, the first source that has a price wins:
   1. **Flyer**: visible offers of the shop's editions (and the household's own user-flyer offers for that shop) for
      the ingredient, valid on the date; the one with the lowest unit price. Loyalty-only offers are included and
      flagged "with app or card" (Decided, Q175).
   2. **Receipt** (R3): the latest non-outlier receipt observation of the shop's chain for the ingredient within the
      last 90 days.
   3. **Reference**: the country's reference price (same for every shop).
   4. **Unknown**: no price.
3. Each cell returns: source, price per pack, pack, unit price and basis, offer id and end date (flyer only),
   loyalty flag and multi-buy terms (not applied; planner decides when the wanted amount reaches them), and
   `estimated = true` for every source except flyer. Totals built from it are labelled "about" (features.md).
4. Cells are returned shop by shop in the request order and line by line in the request order.
5. Shops of removed chains still get receipt and reference prices.
6. The plan check (Plus for the trip plan) is planner's (PLN-19); `GetPrices` itself is not plan-gated, because
   planner also uses it for "about" totals.
7. p95 under 200 ms for 10 shops × 100 lines.

## Contract

`contracts/coolinary/offers/v1/prices.proto`

```proto
rpc GetPrices(GetPricesRequest) returns (GetPricesResponse);

message GetPricesRequest {
  string household_id = 1;
  repeated string shop_ids = 2;
  repeated PriceLine lines = 3;
  string date = 4;                  // YYYY-MM-DD
}
message PriceLine {
  string line_id = 1;
  string ingredient_id = 2;
  coolinary.common.v1.Quantity wanted = 3;
}
message GetPricesResponse {
  repeated ShopPrices shops = 1;
  string currency = 2;
}
message ShopPrices {
  string shop_id = 1;
  repeated PriceCell cells = 2;
}
message PriceCell {
  string line_id = 1;
  PriceSource source = 2;
  coolinary.common.v1.Money pack_price = 3;
  OfferPack pack = 4;
  UnitPrice unit_price = 5;
  string offer_id = 6;
  string valid_to = 7;
  bool loyalty_required = 8;
  OfferConditions conditions = 9;
  bool estimated = 10;
}
enum PriceSource {
  PRICE_SOURCE_UNSPECIFIED = 0; PRICE_SOURCE_FLYER = 1; PRICE_SOURCE_RECEIPT = 2;
  PRICE_SOURCE_REFERENCE = 3; PRICE_SOURCE_UNKNOWN = 4;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Shops in two countries | `INVALID_ARGUMENT` |
| More than 40 shops or 300 lines | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] With the worked example of features.md §2 loaded as flyer offers (Lidl, Kaufland, Aldi; chicken thighs, basmati
      rice, onions, soy sauce, ginger), every cell returns the example's price with source `flyer`, so planner and the
      phone reach Lidl only (10.65) at a €2 minimum and Lidl then Kaufland (9.15) at €1 (Verification 13, offers side).
- [ ] A line with no offer anywhere returns the same reference price for all three shops.
- [ ] From R3, a line with a receipt price at Kaufland only returns `receipt` for Kaufland and `reference` for the
      others.
