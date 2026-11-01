# PLN-19 Trip price matrix

Release R2 (receipt prices R3), Plus only. The trip plan (which 1–3 shops to buy at within the walking limit) runs
on the phone, because only the phone knows distances (decision 7, D13). Planner gives it the numbers it needs: a price
for each open list line at each candidate shop, with the price source, and each shop's total (features.md §2
"Planning the shopping trip", Verification 13).

Price order per cell (features.md "Where prices come from"): the flyer price when the item is on offer at that shop
this week; otherwise the last price seen on receipts at that chain (R3); otherwise a reference price that is the same
for every shop. Totals are labelled "about".

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-19-01](PLN-19-01-candidate-shops.md) | Candidate shops for a trip | R2 | S | OFF-04 (ListMyShops) |
| [PLN-19-02](PLN-19-02-price-per-line-and-shop.md) | Price per line and shop | R2 | M | PLN-19-01, PLN-15-04, OFF-11 (GetPrices) |
| [PLN-19-03](PLN-19-03-price-matrix-rpc.md) | Price matrix RPC and shop totals | R2 | M | PLN-19-02, PLN-15-01 |
| [PLN-19-04](PLN-19-04-receipt-prices.md) | Receipt prices in the matrix | R3 | S | PLN-19-02, OFF-11 (GetPrices, receipt source) |
