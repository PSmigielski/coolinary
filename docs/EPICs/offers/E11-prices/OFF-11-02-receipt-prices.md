# OFF-11-02 Receipt prices from scanned receipts

| | |
|---|---|
| Service | offers |
| Epic | [OFF-11 Prices](README.md) |
| Release | R3 |
| Size | M |
| Depends on | OFF-11-01, OFF-01-02, KIT-08 (receipt.scanned) |
| Unblocks | OFF-11-04 (receipt step) |

## Goal

Shelf prices reach the trip plan: the last price users paid at a chain is used when the item is not on offer
(features.md §2 "Otherwise the last price seen on receipts that users scanned at that chain (from release 3)";
Risks "Shelf prices are not in flyers").

## Business rules

1. Handler for `receipt.scanned` (stream `events.kitchen`). The chain is the payload's `chain_id` if it is a known
   chain; else the payload's chain name is compared case-insensitively with each chain's `receipt_aliases` in the
   payload's country. No match → the event is acknowledged and counted (`receipt_prices_unknown_chain_total`).
2. Each line with an ingredient id, a price above 0 and a purchase date within the last 30 days (not in the future)
   becomes an observation with source `receipt`, the chain, the purchase date, no edition, no shop, no household.
3. Unit price needs the item's amount. The README payload lists `ingredient_id`, `ean`, `unit_price`, `quantity`; this
   task reads `unit_price` as the price of one item and needs the amount of one item (grams, millilitres or
   pieces). Lines without that amount are used only when the ingredient is counted in pieces; others are skipped.
   (Needs agreement with KIT-08 on the payload; see Notes.)
4. Outliers: when the chain has at least 5 receipt observations of the ingredient in the last 90 days, a new one whose
   unit price is outside 0.2× to 5× their median is stored with `outlier = true` and never used.
5. Loyalty or coupon discounts printed on the receipt are already in the line price; nothing else is applied.
6. Nothing in the observation can identify the household: no ids, no store, no time of day.

## Contract

Consumes `receipt.scanned` (docs/EPICs/README.md "Events").

## Data

No schema change (uses `price_observations`).

## Errors

| Case | Behaviour |
|---|---|
| Unknown chain | Event acknowledged and counted |
| Line without an item amount | Skipped (rule 3) |

## Acceptance criteria

- [ ] A receipt event for "LIDL SP. Z O.O." resolves to `lidl-pl` through its alias and adds one observation per
      usable line.
- [ ] A line at 50× the chain's median for that ingredient is stored as outlier and ignored by `GetPrices`.
- [ ] The same event delivered twice adds the lines once (idempotent consumer, OFF-01-02).

## Notes

Resolved: KIT-08-06 now sends per line `quantity` (the whole line in base units), `item_quantity` (packs or pieces
on the line, 1 when unknown) and `unit_price_minor` per kg, litre or piece. Price per kg comes straight from
`unit_price_minor`; the pack size is `quantity` ÷ `item_quantity`.
