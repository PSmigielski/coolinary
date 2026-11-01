# KIT-04-01 GTIN parsing, check digits and prefix classes

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | — |
| Unblocks | every other KIT-04 task |

## Goal

Every scanned linear code is validated and turned into one canonical key, and codes that can't be looked up in a
register (shop-internal, books, coupons) are recognised before any lookup (features.md §4 "Codes that need special
handling").

## Business rules

1. Input: the scanned digits and the symbology reported by the phone (`EAN_13`, `EAN_8`, `UPC_A`, `UPC_E`, `ITF_14`,
   `UNKNOWN`). Spaces and hyphens are removed; any other non-digit is `INVALID_ARGUMENT`.
2. Accepted lengths: 8 (EAN-8, or UPC-E when the symbology says so), 12 (UPC-A), 13 (EAN-13), 14 (GTIN-14).
3. UPC-E is expanded to UPC-A with the standard rules (by its last data digit: 0–2, 3, 4, 5–9) before validation.
4. Check digit: GS1 mod 10 (from the right, data digits weighted 3, 1, 3, 1 …; check = (10 − sum mod 10) mod 10).
   A wrong check digit is `INVALID_ARGUMENT`.
5. Decided (Q067): canonical key = the GTIN-14 form (left-padded with zeros, the GS1 standard), stored in every
   table. Display form = 8, 12 or 13 digits, or 14 when the first digit isn't 0.
6. Classes, by the GTIN-13 form (GTIN-8 and GTIN-14 by their own rules):
   - GTIN-14 indicator 1–8 → `trade_unit_group` (multipack or case, KIT-04-11); indicator 9 → `variable_trade_item`
     (not supported; answered as no hit);
   - prefixes 020–029 and 040–049 (UPC number systems 2 and 4) and 200–299 → `restricted_internal` (KIT-04-03);
   - GTIN-8 starting with 0 or 2 → `restricted_internal`;
   - 977 → `serial_publication`, 978–979 → `book`, 980 → `refund_receipt`, 981–984 and 99 → `coupon`: all answered
     as "not food" without any lookup;
   - 952 → `demo`: resolvable only when `ENVIRONMENT` is not `prod` (the dev example product uses it);
   - everything else → `trade_item` (normal lookup).
7. Library `internal/kitchen/gtin`, pure functions, no I/O.
8. The existing `products.ean` values are migrated to the canonical form (KIT-04-05).

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto`

```proto
enum Symbology { SYMBOLOGY_UNSPECIFIED = 0; SYMBOLOGY_EAN_13 = 1; SYMBOLOGY_EAN_8 = 2; SYMBOLOGY_UPC_A = 3;
  SYMBOLOGY_UPC_E = 4; SYMBOLOGY_ITF_14 = 5; SYMBOLOGY_GS1_DATAMATRIX = 6; SYMBOLOGY_QR = 7; }
enum CodeClass { CODE_CLASS_UNSPECIFIED = 0; CODE_CLASS_TRADE_ITEM = 1; CODE_CLASS_TRADE_UNIT_GROUP = 2;
  CODE_CLASS_RESTRICTED_INTERNAL = 3; CODE_CLASS_NOT_FOOD = 4; CODE_CLASS_DEMO = 5; CODE_CLASS_UNSUPPORTED = 6; }
```

## Data

No schema change.

## Errors

| Case | Status | Message |
|---|---|---|
| Non-digit characters | `INVALID_ARGUMENT` | `code must contain digits only` |
| Wrong length | `INVALID_ARGUMENT` | `code must have 8, 12, 13 or 14 digits` |
| Wrong check digit | `INVALID_ARGUMENT` | `check digit does not match` |

## Acceptance criteria

- [ ] Table tests with real-format examples for each length and class, including a UPC-E expansion and a GS1 952
      demo code.
- [ ] `9520000000017` (the seeded example) canonicalises to `09520000000017` and is `demo`.
- [ ] A code with one digit changed fails the check digit (property test over 10,000 random valid codes).
- [ ] An ISBN (978…) returns "not food" and never reaches any lookup (Verification 14 edge cases).
