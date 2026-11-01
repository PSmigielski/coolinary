# PLN-19-04 Receipt prices in the matrix

| | |
|---|---|
| Service | planner |
| Epic | [PLN-19 Trip price matrix](README.md) |
| Release | R3 |
| Size | S |
| Depends on | PLN-19-02, OFF-11 (GetPrices, receipt source) |
| Unblocks | — |

## Goal

From release 3, when an item is not on offer at a shop, use the last price seen on receipts that users scanned at
that chain, before falling back to the reference price (features.md "Where prices come from").

## Business rules

1. offers gathers receipt prices from `receipt.scanned` events (OFF-11); planner only reads them via `GetPrices`.
2. A receipt price is used only when it is at most 60 days old; older ones are skipped.
3. Receipt prices are per chain (and country), not per store; every shop of the chain gets the same receipt price.
4. Receipt prices are unit prices (per kg, l or piece) and follow PLN-19-02 rule 5's unit-price formula.
5. Cells from receipt prices are always `about = true` and carry the observation date (`observed_on`) so the phone
   can say "price from 12 Oct".

## Contract

`MatrixCell` gains `string observed_on = 8;`.

## Data

No schema change.

## Errors

None of its own: read inside `GetTripPriceMatrix` (PLN-19-03), whose errors apply. Receipt prices older than 60 days are skipped (rule 2) and the cell falls back to the reference price, not an error.

## Acceptance criteria

- [ ] Onions not on offer at Kaufland, a receipt price of 1.09/kg 10 days old → Kaufland cell from the receipt.
- [ ] The same receipt price 90 days old → reference price is used instead.

## Notes

- Decided (Q220): receipt prices older than 60 days are ignored.
