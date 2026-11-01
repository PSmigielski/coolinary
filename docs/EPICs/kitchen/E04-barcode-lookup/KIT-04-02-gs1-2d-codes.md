# KIT-04-02 GS1 2D codes and Digital Link

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-04-01 |
| Unblocks | KIT-04-09; KIT-07 (use-by date from the code) |

## Goal

A GS1 DataMatrix or QR code gives the product code, the use-by date and the batch in one scan, so the user doesn't
type the date (features.md §4 table "2D codes on the pack … Read from release 1").

## Business rules

1. Inputs: a GS1 element string (DataMatrix or GS1 QR; group separator ASCII 29 between variable-length fields; a
   leading symbology identifier such as `]d2` or `]Q3` is removed) or a Digital Link URI (`http`/`https`, at most
   2,000 characters).
2. Application identifiers read: `01` GTIN (14 digits, required); `17` expiry date and `15` best-before date
   (YYMMDD); `10` batch (up to 20 characters); `310n` net weight in kg with `n` decimals (6 digits). Other AIs are
   skipped. `21` (serial) is never stored or returned.
3. Fixed-length AIs (`01`, `11`, `13`, `15`, `16`, `17`, `310n`) need no separator; variable-length ones end at a
   group separator or the end of the string.
4. Dates: `DD = 00` means the last day of the month. The century is chosen so the year lies between 49 years before
   and 50 years after the current year (GS1 rule).
5. Use-by priority: `17` over `15`; `16` (sell-by) is not used. A date in the past is returned with
   `already_expired = true`.
6. Digital Link: the path holds `/01/<gtin>` and optionally `/10/<batch>`, `/21/<serial>`; data attributes are query
   parameters with numeric AI names (`?17=261231&3103=000500`). Any host is accepted. The URI is never fetched.
7. The GTIN from `01` goes through KIT-04-01 rules (check digit, classes). A malformed string is `INVALID_ARGUMENT`.
8. Output: GTIN, `use_by` date, `date_kind` (`expiry` or `best_before`), batch, net weight in grams.

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto`

```proto
message ParsedCode { string gtin = 1; CodeClass code_class = 2; string use_by = 3; DateKind date_kind = 4;
  bool already_expired = 5; string batch = 6; optional double net_weight_g = 7; optional int64 price_minor = 8; }
enum DateKind { DATE_KIND_UNSPECIFIED = 0; DATE_KIND_EXPIRY = 1; DATE_KIND_BEST_BEFORE = 2; }
```

## Data

No schema change.

## Errors

| Case | Status | Message |
|---|---|---|
| No `01` AI | `INVALID_ARGUMENT` | `no product code in 2D code` |
| Bad date or truncated field | `INVALID_ARGUMENT` | `malformed GS1 data` |

## Acceptance criteria

- [ ] `]d201095060001343521726123110ABC123` → GTIN `09506000134352`, use-by 2026-12-31 (expiry), batch `ABC123`.
- [ ] `https://id.gs1.org/01/09506000134352/10/ABC?17=260200` → use-by 2026-02-28.
- [ ] A string with a serial `21` returns no serial anywhere.
- [ ] Verification 14 extended: a sample 2D code fills the use-by date without typing.
