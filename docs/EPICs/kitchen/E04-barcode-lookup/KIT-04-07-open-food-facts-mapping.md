# KIT-04-07 Open Food Facts mapping and storage

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 (GATE: licence check before production) |
| Size | M |
| Status | todo |
| Depends on | KIT-04-06, KIT-03-01 |
| Unblocks | KIT-04-09, KIT-05-04, KIT-06-07 |

## Goal

An Open Food Facts answer becomes a Coolinary product (name, brand, pack size, category, allergens, pack nutrition
values) while OFF data stays separate from Coolinary's own catalogue and is credited.

## Business rules

1. Decided (Q071): OFF-derived products are stored in their own table `off_products`, never mixed into `products`,
   which keeps the share-alike ODbL database separate from Coolinary's catalogue and from household data
   (features.md §8).
2. Name, in order: `product_name_<household language>`, `product_name`, `generic_name`, `product_name_en`. If all are
   empty the answer counts as a miss (the user names it, KIT-05).
3. Brand: the first entry of `brands` (comma-separated), trimmed, at most 60 characters.
4. Pack size: `product_quantity` with `product_quantity_unit` (`g`, `kg`, `ml`, `l`, `cl` converted to g or ml);
   otherwise the `quantity` text through the normaliser's size parser (KIT-03-01), which also detects multipacks
   ("6 x 0,5 l", KIT-04-11).
5. Allergens and traces: OFF tags (`en:gluten`, `en:milk`, `en:nuts`, `en:peanuts`, `en:soybeans`, `en:eggs`,
   `en:fish`, `en:crustaceans`, `en:molluscs`, `en:celery`, `en:mustard`, `en:sesame-seeds`,
   `en:sulphur-dioxide-and-sulphites`, `en:lupin`) map to the 14 codes; unknown tags are ignored.
6. `nutriments` and `nutrition_data_per` are stored raw for KIT-06-07, which validates them.
7. Refresh: a stored row older than 90 days is fetched again on the next real scan (KIT-04-06 rule 2).
8. Every response that uses OFF data carries `attribution = "Open Food Facts (ODbL)"` and the product page URL
   `https://world.openfoodfacts.org/product/<code>`; the app shows it (credit required by the licence).
9. When a household confirms an OFF product into stock, its values are copied to `household_products`
   (KIT-04-04) with `source = 'off'`.

## Contract

`ResolveBarcodeResponse` gains `attribution` and `attribution_url` (KIT-04-09).

## Data

`kitchen/migrations/000NN_off_products.sql`

```sql
CREATE TABLE off_products (
    gtin            char(14) PRIMARY KEY,
    name            text NOT NULL,
    names           jsonb NOT NULL DEFAULT '{}',     -- language -> name
    brand           text NOT NULL DEFAULT '',
    quantity_text   text NOT NULL DEFAULT '',
    pack_value      double precision,
    pack_unit       smallint NOT NULL DEFAULT 0,
    units_per_pack  integer NOT NULL DEFAULT 1,
    categories      text[] NOT NULL DEFAULT '{}',
    countries       text[] NOT NULL DEFAULT '{}',
    allergens       text[] NOT NULL DEFAULT '{}',
    traces          text[] NOT NULL DEFAULT '{}',
    nutriments      jsonb,
    nutrition_per   text,
    ingredient_id   uuid REFERENCES ingredients (id),
    fetched_at      timestamptz NOT NULL DEFAULT now()
);
```

## Errors

A malformed OFF answer counts as a miss and is logged (code only, no content).

## Acceptance criteria

- [ ] A recorded OFF response for a Polish yoghurt maps to name, brand, 150 g and allergen `milk`.
- [ ] "6 x 0,5 l" sets 6 units of 500 ml.
- [ ] The attribution appears in every response that used OFF data and in no other.
