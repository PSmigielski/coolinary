# OFF-11 Prices: flyer, receipt and reference

| | |
|---|---|
| Service | offers |
| Release | R2–R3 |
| Goal | Give planner one price per shopping-list line and shop for the trip plan's price matrix: the flyer price when on offer this week, else the last price seen on users' receipts at that chain (R3), else a reference price that is the same for every shop (features.md §2 "Where prices come from"). |

## Flow

1. Every published flyer offer with a matched ingredient is also recorded as a price observation.
2. From R3, kitchen's `receipt.scanned` events add chain-level receipt prices, without any household id.
3. A weekly job computes a reference price per country and ingredient.
4. `GetPrices` answers planner's matrix request (PLN-19); the phone then runs the trip plan (README decision 7).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-11-01](OFF-11-01-price-observations.md) | Price observations from flyers | R2 | M | OFF-07-07 |
| [OFF-11-02](OFF-11-02-receipt-prices.md) | Receipt prices from scanned receipts | R3 | M | OFF-11-01, OFF-01-02, KIT-08 (receipt.scanned) |
| [OFF-11-03](OFF-11-03-reference-prices.md) | Reference prices per country | R2 | M | OFF-11-01, OFF-01-05 |
| [OFF-11-04](OFF-11-04-get-prices.md) | GetPrices for the trip price matrix | R2 | M | OFF-11-01, OFF-11-03, OFF-04-01 |
