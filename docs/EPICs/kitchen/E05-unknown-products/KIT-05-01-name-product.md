# KIT-05-01 Name an unknown product once

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-05 Unknown products](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-04-04, KIT-04-05, KIT-04-09 |
| Unblocks | KIT-05-02, KIT-05-03, KIT-07 (unpack mode with unknown items), GW-06 |

## Goal

A household that scans an unknown code types its name once; from then on the code is known in this household, and a
proposal is ready for the next household.

## Business rules

1. `NameProduct(household_id, code, name, brand?, pack_size?, units_per_pack?, ingredient_id?, chain?, source)` with
   `source` `typed` or `label` (from KIT-05-03).
2. The code is parsed as in KIT-04-01. Shop-internal codes are stored in `household_internal_codes` (KIT-04-03) with
   the optional `chain` and never become submissions.
3. Validation: name 2–120 characters after trimming; brand ≤ 60; pack size > 0 and ≤ 100 kg / 100 l / 500 pieces;
   `units_per_pack` 1–100; `ingredient_id` must be active. The normaliser's `not_a_name` rule (KIT-03-01 rule 4)
   applies to name and brand.
4. If `ingredient_id` is missing, `MatchNames` (`INTERACTIVE`, context `product`) fills it when the answer is
   `matched`.
5. Effects in one transaction: upsert `household_products` (source `named` or `label`), and create or update this
   household's submission (KIT-05-02). The response is the household's product.
6. If the code is already `confirmed` in the shared catalogue, the household's name is saved only on its own row
   (no submission); the shared row is unchanged.
7. Decided (Q073): at most 50 namings per household per day (redis-state counter
   `kitchen:naming:<household>:<date>`), then `RESOURCE_EXHAUSTED` until the next day. This stops abuse of the
   shared catalogue.

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto`

```proto
rpc NameProduct(NameProductRequest) returns (NameProductResponse);
message NameProductRequest { string household_id = 1; string code = 2; Symbology symbology = 3; string name = 4;
  string brand = 5; coolinary.common.v1.Quantity pack_size = 6; int32 units_per_pack = 7; string ingredient_id = 8;
  string chain = 9; NamingSource source = 10; string language = 11; string country = 12; }
enum NamingSource { NAMING_SOURCE_UNSPECIFIED = 0; NAMING_SOURCE_TYPED = 1; NAMING_SOURCE_LABEL = 2; }
message NameProductResponse { HouseholdProduct product = 1; bool submitted_to_shared = 2; }
```

## Data

No new table (uses KIT-04-03, KIT-04-04, KIT-05-02 tables).

## Errors

| Case | Status | Message |
|---|---|---|
| Bad name, brand or size | `INVALID_ARGUMENT` | names the field |
| More than 50 namings today | `RESOURCE_EXHAUSTED` | `naming limit reached for today` |
| Household mismatch | `PERMISSION_DENIED` | `household mismatch` |

## Acceptance criteria

- [ ] A named product resolves from level 1 for the same household at once.
- [ ] Naming a confirmed product changes only the household's row.
- [ ] A name containing an e-mail address is refused.
