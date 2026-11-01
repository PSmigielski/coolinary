# KIT-04-03 Shop-internal codes (prefixes 20–29)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-04-01, KIT-10 (household country) |
| Unblocks | KIT-04-09, Verification 14 ("a weighed-counter code … handled") |

## Goal

Weighed cheese, meat and bakery codes are handled: the weight (or price) is read where the chain's pattern is known,
and otherwise the household names the item once (features.md §4 "Codes starting with 20 to 29").

## Business rules

1. Applies to codes classed `restricted_internal` (KIT-04-01). They never go to the shared catalogue, Open Food Facts
   or GS1.
2. Patterns are data, not code: per country (and optionally per chain), a prefix, the positions of the item digits,
   the positions of the value digits, the value kind (`weight_g` or `price_minor`), decimals, and whether a price
   check digit is present. Admin maintains them (ADM-03); the seed starts empty.
3. For a code, the household's country patterns with a matching prefix are tried. Exactly one match → item key and
   value are read. Several matches (different chains use the same prefix differently) and no chain known → only the
   item key is read, the value is ignored.
4. Decided (Q068): with no pattern, the item key is the first 7 digits (2-digit prefix + 5 item digits), the common
   European layout 2 + 5 + 5 + 1, so a new weight of a known item doesn't ask for a name again.
5. Decided (Q069): names for internal codes are kept per household, keyed by (household, country, item key),
   optionally with the chain the user picked (the same digits mean different things at different chains). They are
   never shared with other households. When two names exist for one key with different chains, both are returned
   as choices.
6. Known key → found, source `INTERNAL_CODE`, with the household's name and ingredient, plus the weight (as the
   quantity) or the price (as price paid) when a pattern gave it.
7. Unknown key → not found with `needs_naming = true`; naming goes through KIT-05-01 with the internal key.
8. Rows are deleted with the household (KIT-13).
9. Admin RPC `DecodeInternalCode(pattern, code)` decodes a sample code with a pattern that may not be saved
   yet, with the same decoder as `ResolveBarcode`, and returns the item key, value and check-digit results without
   writing anything (the admin test box, ADM-03-07).

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto` — `ResolveSource` gains `RESOLVE_SOURCE_INTERNAL_CODE = 5`.
`contracts/coolinary/kitchen/v1/admin.proto`

```proto
rpc ListInternalCodePatterns(ListInternalCodePatternsRequest) returns (ListInternalCodePatternsResponse);
rpc SetInternalCodePattern(InternalCodePattern) returns (InternalCodePattern);
rpc DecodeInternalCode(DecodeInternalCodeRequest) returns (DecodeInternalCodeResponse);
message DecodeInternalCodeRequest { InternalCodePattern pattern = 1; string code = 2; }
message DecodeInternalCodeResponse { string item_key = 1; int64 value = 2; string value_kind = 3; bool valid_check_digit = 4; bool price_check_digit_ok = 5; }
message InternalCodePattern { int32 id = 1; string country = 2; string chain = 3; string prefix = 4;
  int32 item_start = 5; int32 item_length = 6; int32 value_start = 7; int32 value_length = 8;
  string value_kind = 9; int32 value_decimals = 10; bool price_check_digit = 11; bool enabled = 12; string note = 13; }
```

## Data

`kitchen/migrations/000NN_internal_codes.sql`

```sql
CREATE TABLE internal_code_patterns (
    id             serial PRIMARY KEY,
    country        char(2) NOT NULL,
    chain          text,
    prefix         text NOT NULL CHECK (prefix ~ '^2[0-9]$'),
    item_start     integer NOT NULL, item_length integer NOT NULL,
    value_start    integer NOT NULL, value_length integer NOT NULL,
    value_kind     text NOT NULL CHECK (value_kind IN ('weight_g','price_minor')),
    value_decimals integer NOT NULL DEFAULT 0,
    price_check_digit boolean NOT NULL DEFAULT false,
    enabled        boolean NOT NULL DEFAULT true,
    note           text NOT NULL DEFAULT ''
);
CREATE TABLE household_internal_codes (
    household_id   uuid NOT NULL,
    country        char(2) NOT NULL,
    item_key       text NOT NULL,
    chain          text NOT NULL DEFAULT '',
    name           text NOT NULL,
    ingredient_id  uuid REFERENCES ingredients (id),
    created_at     timestamptz NOT NULL DEFAULT now(),
    last_used_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, country, item_key, chain)
);
```

## Errors

| Case | Status | Message |
|---|---|---|
| Pattern positions outside 13 digits | `INVALID_ARGUMENT` | `pattern does not fit a 13-digit code` |
| `DecodeInternalCode` with a code that isn't 13 digits or doesn't start with the pattern's prefix | `INVALID_ARGUMENT` | `code does not match the pattern` |

## Acceptance criteria

- [ ] With a pattern "PL, prefix 29, item 3–7, weight 8–12 in g", code `2912345004569` returns 456 g of the named
      item.
- [ ] With no pattern, two scans of the same cheese at different weights ask for a name once.
- [ ] An internal code never triggers an Open Food Facts request (checked with a stub that fails on any call).
- [ ] `DecodeInternalCode` with the pattern above and `2912345004569` returns the same item key `ResolveBarcode`
      derives for that code, value 456 and `value_kind = weight_g`, and writes nothing.
