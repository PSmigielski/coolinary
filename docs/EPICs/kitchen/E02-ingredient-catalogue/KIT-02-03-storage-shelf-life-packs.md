# KIT-02-03 Storage place, shelf life and pack sizes

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-02-01 |
| Unblocks | KIT-04-10 (scan suggestions), KIT-07 (default use-by), PLN-15 (round up to packs) |

## Goal

Know where an ingredient usually goes and how long it keeps there, so a scanned or typed item gets a sensible
storage place and use-by date, and know usual pack sizes so the shopping list can round up to whole packs
(features.md §2 "Fridge and pantry"; §2 "Shopping list" step 3).

## Business rules

1. `default_storage` is `fridge`, `freezer` or `pantry`.
2. Shelf life in days per storage place: `shelf_life_fridge_days`, `shelf_life_freezer_days`,
   `shelf_life_pantry_days`; `null` means "not suitable or unknown". Ranges: fridge 1–60, freezer 1–365,
   pantry 1–1,095.
3. Suggested use-by for an item stored in place P on day D (household's local date) = D + shelf life for P. If P has
   no value, the default storage's value is used and the suggestion carries `storage_unusual = true`. If neither
   exists, no date is suggested.
4. Conservative by design: shelf lives are for the unopened item as bought; the user can always change the date.
5. Pack sizes are listed per country (ISO 3166-1 alpha-2) as quantities in the ingredient's base unit, with one
   marked `is_default`. Country code `ZZ` holds sizes used when a country has none.
6. Pack lookup for country C returns C's sizes ascending; else `ZZ`'s; else none (the list then shows the exact
   amount, features.md §2 "where the catalogue knows the pack size").

## Contract

`contracts/coolinary/kitchen/v1/catalogue.proto`

```proto
enum StoragePlace { STORAGE_PLACE_UNSPECIFIED = 0; STORAGE_PLACE_FRIDGE = 1; STORAGE_PLACE_FREEZER = 2;
  STORAGE_PLACE_PANTRY = 3; }
message Storage { StoragePlace default_place = 1; optional int32 fridge_days = 2; optional int32 freezer_days = 3;
  optional int32 pantry_days = 4; }
message PackSize { coolinary.common.v1.Quantity quantity = 1; bool is_default = 2; }
```

## Data

`kitchen/migrations/000NN_storage_packs.sql`

```sql
ALTER TABLE ingredients
    ADD COLUMN default_storage text CHECK (default_storage IN ('fridge','freezer','pantry')),
    ADD COLUMN shelf_life_fridge_days  integer CHECK (shelf_life_fridge_days  BETWEEN 1 AND 60),
    ADD COLUMN shelf_life_freezer_days integer CHECK (shelf_life_freezer_days BETWEEN 1 AND 365),
    ADD COLUMN shelf_life_pantry_days  integer CHECK (shelf_life_pantry_days  BETWEEN 1 AND 1095);
CREATE TABLE ingredient_pack_sizes (
    ingredient_id  uuid NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    country        char(2) NOT NULL,
    value          double precision NOT NULL CHECK (value > 0),
    is_default     boolean NOT NULL DEFAULT false,
    PRIMARY KEY (ingredient_id, country, value)
);
CREATE UNIQUE INDEX ingredient_pack_default ON ingredient_pack_sizes (ingredient_id, country) WHERE is_default;
```

## Errors

Validation errors come with the admin RPCs (KIT-02-10).

## Acceptance criteria

- [ ] Chicken thighs (fridge 2 days, freezer 180) added to the fridge on 3 Nov get 5 Nov; to the freezer, 2 May.
- [ ] An ingredient with only `pantry_days` placed in the fridge gets the pantry date and `storage_unusual`.
- [ ] Pack lookup for `PL` falls back to `ZZ` and then to none.
