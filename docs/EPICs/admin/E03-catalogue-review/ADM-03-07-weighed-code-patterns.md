# ADM-03-07 Weighed-item barcode patterns

| | |
|---|---|
| Service | admin |
| Epic | [ADM-03 Catalogue, matching and nutrition review](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | ADM-02-03; KIT-04-03 (`ListInternalCodePatterns`, `SetInternalCodePattern`, `DecodeInternalCode`) |
| Unblocks | — |

## Goal

Shop-internal barcodes starting with 20–29 (weighed cheese, meat, bakery) carry a weight or a price whose layout
differs by country and chain (features.md §4 "Codes that need special handling"). The founder maintains these
layouts so the app can read the weight.

## Business rules

1. A pattern has kitchen's fields (`InternalCodePattern`, KIT-04-03): country, chain (empty = any), prefix (two
   digits 20–29), item digit positions (`item_start`, `item_length`), value digit positions (`value_start`,
   `value_length`), value kind (`weight_g` or `price_minor`), decimals, whether a price check digit is present,
   enabled, and a note.
2. A test box decodes a sample code with the pattern being edited (`DecodeInternalCode(pattern, code)`) and shows
   item key, value and check-digit results before saving.
3. Saving a pattern while another enabled one has the same country, chain and prefix → 409; the admin disables the
   old one first. Admin checks this with `ListInternalCodePatterns(country)` before calling
   `SetInternalCodePattern`.
4. Audited (`kitchen.internal_code_pattern.set`).
5. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/kitchen/internal-code-patterns?country=` | `{ items }` |
| PUT | `/admin/v1/kitchen/internal-code-patterns/{id}` (or POST to create) | pattern → pattern |
| POST | `/admin/v1/kitchen/internal-code-patterns/decode` | `{ pattern, code }` → `{ item_key, value, value_kind, valid_check_digit, price_check_digit_ok }` |

Calls `KitchenAdminService.ListInternalCodePatterns`, `SetInternalCodePattern`, `DecodeInternalCode` (KIT-04-03).

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Overlapping active pattern | 409 |
| Code does not fit the layout | 400 |

## Acceptance criteria

- [ ] A sample Polish code with prefix 29 decodes to the expected grams with the test box before saving.
