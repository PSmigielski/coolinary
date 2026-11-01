# PLN-19-02 Price per line and shop

| | |
|---|---|
| Service | planner |
| Epic | [PLN-19 Trip price matrix](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | PLN-19-01, PLN-15-04, OFF-11 (GetPrices) |
| Unblocks | PLN-19-03, PLN-19-04, PLN-20-03 |

## Goal

Work out what each open list line would cost at each candidate shop, and where that price comes from.

## Business rules

1. Lines: open (unticked) ingredient lines. Free-text lines are not priced and are reported as `unpriced_lines`.
2. Prices: one `GetPrices(ingredient ids, shop ids, today)` call. For each (ingredient, shop) offers returns the
   available observations: flyer offer (price, pack, conditions, valid dates), receipt price at the chain (R3,
   PLN-19-04), reference price (unit price per kg, l or piece, same for every shop in the country).
3. Source order per cell: flyer → receipt (R3) → reference. The first available is used.
4. Amount bought: pack size × pack count from the line when known (what will actually be bought), else the line's
   amount.
5. Cost:
   - flyer price with a pack: `packs = ceil(amount bought / offer pack)`, `cost = packs × offer price`;
   - unit price (receipt or reference): `cost = amount bought in kg, l or pieces × unit price`;
   - each cell rounded to whole cents, halves up.
6. Conditions: `app_or_card` and `multi_buy` flags are copied to the cell; the price is taken as stated (no multi-buy
   arithmetic). The phone shows the condition.
7. A cell with no source at all is `unpriced`; a line unpriced at every shop is reported once.
8. A cell is `about = true` unless its source is a flyer offer valid today.
9. All cells use the household's currency (the country's); an observation in another currency is ignored.

## Contract

Internal (`trip.PriceCells(lines, shops, prices) [][]Cell`).

## Data

No schema change. `GetPrices` results cached in redis-cache for 15 minutes per (ingredient set, shop set, date).

## Errors

| Case | Behaviour |
|---|---|
| offers down | `UNAVAILABLE` from `GetTripPriceMatrix` |

## Acceptance criteria

- [ ] Chicken thighs 1 kg on the list, Lidl flyer 3.89 per 1 kg pack → Lidl cell 3.89, source flyer, not "about".
- [ ] Need 600 g rounded to 1 × 1 kg, no flyer, reference 2.19 per kg → 2.19, source reference, "about".
- [ ] Need 1.5 kg, offer pack 1 kg at 3.89 → 2 packs = 7.78.
- [ ] A Lidl "with app" offer carries the `app_or_card` flag.
