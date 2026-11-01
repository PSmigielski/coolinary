# KIT-04-05 Shared catalogue (level 2)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-04-01 |
| Unblocks | KIT-04-09, KIT-04-11, KIT-05-02 |

## Goal

Every product a user has named and a second household has confirmed is known to everyone (features.md §4 step 2
"Coolinary's shared catalogue: every product any user has confirmed").

## Business rules

1. The existing `products` table becomes the shared catalogue. Its key moves to the canonical 14-digit GTIN
   (KIT-04-01); the example row `9520000000017` becomes `09520000000017`.
2. Only Coolinary's own data lives here: products named by users (KIT-05) or entered by admin. Data from Open Food
   Facts and GS1 is kept in their own tables (KIT-04-07, KIT-04-08).
3. `status`: `pending` (named by one household), `confirmed` (KIT-05-02 rule), `disputed` (in admin review),
   `retired`. The existing `confirmed boolean` is replaced by `status`.
4. Level 2 answers only `confirmed` rows. A `pending` row is shown to other households only as a proposal
   (KIT-05-02), never as a hit.
5. A row may carry `ingredient_id`, `category`, `units_per_pack` and `unit_gtin` (KIT-04-11).
6. Cache: redis-cache `kitchen:gtin:<gtin>` for 24 hours, deleted on any change to the row. The existing
   `kitchen:ean:` keys are dropped.
7. `report_count` counts "this is wrong" reports from distinct households (KIT-05-02).

## Contract

`Product` in `kitchen.proto` gains fields:

```proto
message Product { string ean = 1; string name = 2; string brand = 3; coolinary.common.v1.Quantity pack_size = 4;
  string gtin = 5; string ingredient_id = 6; string ingredient_name = 7; int32 units_per_pack = 8;
  string unit_gtin = 9; repeated coolinary.common.v1.Allergen allergens = 10; repeated coolinary.common.v1.Allergen traces = 11; }
```

(`ean` keeps the display form for existing callers.)

## Data

`kitchen/migrations/000NN_products_shared.sql`

```sql
ALTER TABLE products RENAME COLUMN ean TO gtin;
UPDATE products SET gtin = lpad(gtin, 14, '0');
ALTER TABLE products ALTER COLUMN gtin TYPE char(14);
ALTER TABLE products
    ADD COLUMN status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','confirmed','disputed','retired')),
    ADD COLUMN source text NOT NULL DEFAULT 'user' CHECK (source IN ('user','admin')),
    ADD COLUMN ingredient_id uuid REFERENCES ingredients (id),
    ADD COLUMN category text,
    ADD COLUMN units_per_pack integer NOT NULL DEFAULT 1 CHECK (units_per_pack >= 1),
    ADD COLUMN unit_gtin char(14),
    ADD COLUMN report_count integer NOT NULL DEFAULT 0,
    ADD COLUMN confirmed_at timestamptz,
    ADD COLUMN updated_at timestamptz NOT NULL DEFAULT now();
UPDATE products SET status = 'confirmed', confirmed_at = now() WHERE confirmed;
ALTER TABLE products DROP COLUMN confirmed;
```

## Errors

None beyond the chain's (KIT-04-09).

## Acceptance criteria

- [ ] After migration the example product resolves from level 2 by both `9520000000017` and `09520000000017`
      (in dev).
- [ ] A `pending` product is not a hit for another household.
- [ ] Changing a product deletes its cache key.
