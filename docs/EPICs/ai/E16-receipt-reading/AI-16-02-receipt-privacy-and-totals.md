# AI-16-02 Receipt privacy filter and totals check

| | |
|---|---|
| Service | ai |
| Epic | [AI-16 Receipt reading](README.md) |
| Release | R3 |
| Size | S |
| Status | todo |
| Depends on | AI-16-01 |
| Unblocks | AI-16-03 |

## Goal

Receipts carry personal data; none of it leaves ai, and a misread receipt is flagged rather than trusted.

## Business rules

1. The prompt asks for no personal or payment data. ai also filters every output string: any run of 8 or more
   digits (card numbers, loyalty ids, terminal and transaction ids) is removed unless it is the line's `ean` field
   with a valid GTIN-8/13 check digit; e-mail addresses and phone numbers are removed; lines whose text is only such
   data are dropped.
2. `store_city_printed` keeps the city name only (no street, no postcode).
3. Totals: sum of `line_total − discount` over non-deposit lines, plus deposit lines, compared with `total`; a
   difference above 1% or 0.50 of the currency unit (whichever is larger) sets `total_mismatch=true`. The result is
   still returned; kitchen shows "check the items" and does not send prices from a mismatched receipt to offers.
4. Each line: `line_total ≈ quantity × unit_price` within 0.01 per line, else the line's prices are cleared (text
   and quantity kept) and counted in `ai_output_dropped_total{rule="line_math"}`.

## Contract

No change (fields in AI-16-01).

## Data

No schema change.

## Errors

None new; a total mismatch is a flag, not an error.

## Acceptance criteria

- [ ] A receipt photo showing a card number "**** 1234" and a loyalty number "1234 5678 9012" returns neither.
- [ ] A receipt with one misread line total sets `total_mismatch`.
