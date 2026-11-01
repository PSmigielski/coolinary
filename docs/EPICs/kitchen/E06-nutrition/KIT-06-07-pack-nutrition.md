# KIT-06-07 Pack nutrition values of scanned products

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-04-07, KIT-06-01 |
| Unblocks | KIT-06-08 |

## Goal

When a household cooks with a product it scanned and the pack's values are known, those values replace the table's
(features.md §5 step 4). EU packs must state energy, fat, saturates, carbohydrate, sugars, protein and salt per
100 g; only values that pass checks are used.

## Business rules

1. Source: Open Food Facts `nutriments` stored by KIT-04-07, only when `nutrition_data_per = '100g'`. Fields:
   `energy-kcal_100g`, `energy-kj_100g`, `fat_100g`, `saturated-fat_100g`, `carbohydrates_100g`, `sugars_100g`,
   `proteins_100g`, `salt_100g` (else `sodium_100g` × 2.5), `fiber_100g` (optional).
2. Complete = all seven mandatory values present (fibre optional).
3. Decided (Q083): plausibility rules: each value 0–100 g; fat + carbohydrate + protein + fibre +
   salt ≤ 101 g; saturates ≤ fat; sugars ≤ carbohydrate; stated kcal within ±20% or ±20 kcal (whichever is larger)
   of 4 × protein + 4 × carbohydrate + 9 × fat + 2 × fibre. Failing values are stored with `plausible = false` and
   not used.
4. Liquids: OFF reports drinks per 100 g or per 100 ml in the same fields; values are used as per 100 g only when the
   ingredient has a density, else only for lines given in millilitres.
5. Stored per GTIN in `product_nutrition`, refreshed with the OFF row (KIT-04-07 rule 7).
6. Pack values are used only by calculations that name the product's GTIN on a line (KIT-06-08). Library dishes
   (planner) use table values.

## Contract

No RPC.

## Data

`kitchen/migrations/000NN_product_nutrition.sql`

```sql
CREATE TABLE product_nutrition (
    gtin            char(14) PRIMARY KEY,
    source          text NOT NULL DEFAULT 'off_pack',
    energy_kcal numeric(7,2), energy_kj numeric(8,2), protein_g numeric(6,2), carbohydrate_g numeric(6,2),
    sugars_g numeric(6,2), fat_g numeric(6,2), saturates_g numeric(6,2), fibre_g numeric(6,2), salt_g numeric(7,3),
    complete        boolean NOT NULL,
    plausible       boolean NOT NULL,
    per_basis       text NOT NULL CHECK (per_basis IN ('100g','100ml')),
    fetched_at      timestamptz NOT NULL DEFAULT now()
);
```

## Errors

None returned to users.

## Acceptance criteria

- [ ] A recorded OFF yoghurt with all seven values is stored complete and plausible.
- [ ] Values with sugars above carbohydrate are stored with `plausible = false` and never used.
- [ ] Per-serving-only data is not stored.
