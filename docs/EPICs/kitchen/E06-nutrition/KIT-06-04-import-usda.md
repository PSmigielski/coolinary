# KIT-06-04 Import USDA FoodData Central

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-06-01 |
| Unblocks | KIT-06-05 |

## Goal

Load USDA FoodData Central (public domain, CC0; credit requested) to fill gaps such as Asian sauces and spices
(features.md §5 table).

## Business rules

1. Command `server import-nutrition --source usda --dir <path> --version <label> [--dry-run]`; same behaviour as
   KIT-06-02 rules 1, 5, 6, 7, 8.
2. Only the Foundation Foods and SR Legacy datasets (CSV). Branded foods are not imported (scanned products use their
   pack values, KIT-06-07).
3. Mapping by nutrient id: energy kcal 1008, else 2048 (Atwater specific), else 2047 (Atwater general); protein 1003;
   carbohydrate by difference 1005; sugars total 2000; fat 1004; saturated fatty acids 1258; fibre 1079; sodium 1093
   (mg → g, salt = sodium × 2.5). kJ computed from kcal.
4. Name: the USDA description (English) in both `name` and `name_en`.
5. `state` from the description as in KIT-06-02 rule 4 ("raw", "dry", "cooked", "boiled" …).
6. Credit text: "U.S. Department of Agriculture, Agricultural Research Service. FoodData Central."

## Contract

No RPC (CLI).

## Data

No schema change.

## Errors

As KIT-06-02.

## Acceptance criteria

- [ ] Foundation and SR Legacy entries are imported; branded ones are skipped.
- [ ] A Foundation food with only 2047/2048 energy gets kcal from those.
- [ ] 10 sample foods match FoodData Central's published values.
