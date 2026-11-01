# KIT-04-11 Multipacks

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | KIT-04-01, KIT-04-05, KIT-03-01 |
| Unblocks | KIT-07 (stock amounts), Verification 14 ("… and a multipack are handled") |

## Goal

A multipack has its own code; Coolinary knows how many units it holds, so stock and recipes count the right amount
(features.md §4 "Multipacks have their own code, different from the single item; the catalogue stores how many units
a pack holds").

## Business rules

1. A product has `units_per_pack` (≥ 1) and optionally `unit_gtin` (the single item's code) and the single unit's
   size.
2. Sources of the count: OFF or label text "N x size" (KIT-03-01), a GTIN-14 with indicator 1–8 whose inner GTIN
   (digits 2–13 with a new check digit) is known, or the user when naming (KIT-05-01).
3. For an indicator-1–8 code with an unknown count, the response returns the inner product as `unit_product` and
   `needs_unit_count = true`; the app asks "how many in the pack?".
4. Decided (Q072): a multipack is stocked as one line: for a `PIECE` ingredient the quantity is N pieces;
   otherwise N × unit size in the base unit (6 × 0.5 l water → 3,000 ml).
5. `units_per_pack` is between 1 and 100.

## Contract

Fields already listed in KIT-04-05 (`units_per_pack`, `unit_gtin`); `ResolveBarcodeResponse` gains
`Product unit_product = 14; bool needs_unit_count = 15;`.

## Data

No schema change beyond KIT-04-05.

## Errors

| Case | Status | Message |
|---|---|---|
| `units_per_pack` outside 1–100 | `INVALID_ARGUMENT` | `units per pack must be 1-100` |

## Acceptance criteria

- [ ] A 6 × 0.5 l water multipack stocks 3,000 ml.
- [ ] A 4 × 125 g yoghurt pack of a `PIECE` ingredient stocks 4 pieces.
- [ ] An ITF-14 code with indicator 1 and a known inner product asks for the count once; the answer is saved.
