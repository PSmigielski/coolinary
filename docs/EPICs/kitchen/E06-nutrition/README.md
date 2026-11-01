# KIT-06 Nutrition tables and calculation

| | |
|---|---|
| Service | kitchen |
| Release | R1 |
| Unblocks | PLN-04 (calories per dish), PLN-12 (portions), GW-07, Verification 16 |

## Goal

Calories and nutrients are always calculated without AI, from official nutrition tables or the pack (features.md §5
"How a dish's calories are calculated"). Kitchen imports Ciqual, BLS 4.0 and USDA FoodData Central, links each
catalogue ingredient to one table entry (AI suggests, a person approves), keeps validated pack values of scanned
products, and offers one calculation RPC. A dish with any unmatched ingredient says "calories incomplete" and names
it; it never shows a number that is too low.

## Flow

```
Ciqual / BLS / USDA files ──import──> nutrition_entries
ingredient ──trigram candidates──> batch match_nutrition_entries (AI-09) ──> suggestion ──admin approves (ADM-03)──> approved link
OFF nutriments ──validate──> product_nutrition
planner ──ComputeNutrition(lines, portions)──> totals, per portion (rounded and exact), or "incomplete" + missing list
any change ──> nutrition_version + catalogue.nutrition_changed ──> planner recalculates dishes
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-06-01](KIT-06-01-nutrition-schema.md) | Nutrition entries, links and sources | R1 | M | KIT-02-01 |
| [KIT-06-02](KIT-06-02-import-ciqual.md) | Import Ciqual | R1 (GATE) | M | KIT-06-01 |
| [KIT-06-03](KIT-06-03-import-bls.md) | Import Bundeslebensmittelschlüssel 4.0 | R1 (GATE) | M | KIT-06-01 |
| [KIT-06-04](KIT-06-04-import-usda.md) | Import USDA FoodData Central | R1 | M | KIT-06-01 |
| [KIT-06-05](KIT-06-05-match-suggestions.md) | Ingredient → nutrition entry suggestions | R1 | M | KIT-06-01, KIT-06-02, KIT-01-08, AI-09 |
| [KIT-06-06](KIT-06-06-admin-nutrition-approval.md) | Admin approval of nutrition links | R1 | M | KIT-06-05, KIT-01-07, KIT-06-10 |
| [KIT-06-07](KIT-06-07-pack-nutrition.md) | Pack nutrition values of scanned products | R1 | M | KIT-04-07, KIT-06-01 |
| [KIT-06-08](KIT-06-08-calculation-rules.md) | Calculation rules (library) | R1 | M | KIT-06-01, KIT-02-02, KIT-06-07 |
| [KIT-06-09](KIT-06-09-compute-nutrition-rpc.md) | ComputeNutrition RPC | R1 | M | KIT-06-08 |
| [KIT-06-10](KIT-06-10-nutrition-changed-event.md) | Nutrition version and `catalogue.nutrition_changed` | R1 | S | KIT-06-01, KIT-01-02, KIT-02-02 |

## Gates

- Ciqual and BLS: licence check for use in a paid app and the exact credit wording (features.md §5 open points).
- Polish tables (NIZP PZH-PIB): not imported until a licence is agreed; Polish ingredients use the closest Ciqual or
  BLS entry.
- A dietitian checks the calorie values of the 200 starter dishes before launch (features.md §5; PLN-04).

## Decisions (taken 1 November 2026)

- Q082 (KIT-06-01): "Traces" = 0; "< x" = half of x; both flagged.
- Q083 (KIT-06-07): Each 0–100 g, macros ≤ 101 g, kcal within ±20% of computed.
- Q084 (KIT-06-08): Recipe must state absorbed oil; otherwise "calories incomplete" and not publishable.
- Q085 (KIT-06-08): Yes: EU label rounding; kJ to the nearest 10.
