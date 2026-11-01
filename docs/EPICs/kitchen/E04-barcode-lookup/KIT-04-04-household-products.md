# KIT-04-04 Household products (level 1)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-04-01 |
| Unblocks | KIT-04-09, KIT-05-01, KIT-07 (unpack confirm writes here), MOB-03 (offline cache) |

## Goal

Anything a household has scanned before resolves instantly with the household's own name for it, also offline on
the phone (features.md §4 step 1 "This household's own products … Works offline").

## Business rules

1. A row per (household, GTIN) is written when the household confirms a scanned product into stock (KIT-07) or names
   one (KIT-05-01): upsert, `scan_count + 1`, `last_scanned_at = now()`.
2. The household's values (name, brand, pack size, ingredient, storage place) override the shared ones for this
   household only. Editing them never changes the shared catalogue.
3. Level 1 of the chain: a row for (household, GTIN) answers with source `HOUSEHOLD`.
4. `ListHouseholdProducts(household_id, updated_since, page)` returns the household's rows changed since a time, 500
   per page, so the phone keeps an offline copy and resolves known codes without the network.
5. Decided (Q070): at most 5,000 products per household; beyond that the least recently scanned are deleted by a
   nightly job (a typical household scans a few hundred).
6. Rows are household data: part of the data export and deleted with the household (KIT-13).
7. `household_id` comes from the request, checked with the household guard (KIT-01-01).

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto`

```proto
rpc ListHouseholdProducts(ListHouseholdProductsRequest) returns (ListHouseholdProductsResponse);
message ListHouseholdProductsRequest { string household_id = 1; string updated_since = 2; coolinary.common.v1.PageRequest page = 3; }
message HouseholdProduct { string gtin = 1; string display_code = 2; string name = 3; string brand = 4;
  coolinary.common.v1.Quantity pack_size = 5; string ingredient_id = 6; StoragePlace storage = 7;
  string updated_at = 8; }
message ListHouseholdProductsResponse { repeated HouseholdProduct products = 1; coolinary.common.v1.PageResponse page = 2;
  string server_time = 3; }
```

## Data

`kitchen/migrations/000NN_household_products.sql`

```sql
CREATE TABLE household_products (
    household_id     uuid NOT NULL,
    gtin             char(14) NOT NULL,
    name             text NOT NULL,
    brand            text NOT NULL DEFAULT '',
    pack_value       double precision,
    pack_unit        smallint NOT NULL DEFAULT 0,
    ingredient_id    uuid REFERENCES ingredients (id),
    storage          text CHECK (storage IN ('fridge','freezer','pantry')),
    source           text NOT NULL CHECK (source IN ('shared','off','gs1','named','label')),
    scan_count       integer NOT NULL DEFAULT 1,
    first_scanned_at timestamptz NOT NULL DEFAULT now(),
    last_scanned_at  timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, gtin)
);
CREATE INDEX household_products_updated ON household_products (household_id, updated_at);
```

## Errors

| Case | Status | Message |
|---|---|---|
| Household mismatch | `PERMISSION_DENIED` | `household mismatch` |
| Bad `updated_since` | `INVALID_ARGUMENT` | `updated_since must be RFC 3339` |

## Acceptance criteria

- [ ] A product scanned and confirmed once resolves from level 1 on the next scan, with the household's own name.
- [ ] Household A's rename is invisible to household B.
- [ ] `ListHouseholdProducts` with the previous `server_time` returns only changes since then.
