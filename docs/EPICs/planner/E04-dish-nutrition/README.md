# PLN-04 Dish nutrition and calories

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | Every dish shows "about N kcal" per portion and, on tap, protein, carbohydrate, sugars, fat, saturates, fibre and salt — calculated without AI from kitchen's nutrition tables (features.md §5), kept current when the tables change, and honest when an ingredient has no values ("calories incomplete"). Free for everyone. |

## Flow

```
dish version written ──> ComputeNutrition (KIT-06) ──> dish_nutrition (per portion, rounded) ──> every dish response
catalogue.nutrition_changed / ingredient_merged ──> affected versions queued ──> recompute job
household views dish ──> scanned product with pack values in stock? ──> per-household override (PLN-04-04)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-04-01](PLN-04-01-compute-dish-nutrition.md) | Compute and store nutrition per dish version | R1 | M | PLN-02-01, PLN-01-02, KIT-06 (ComputeNutrition) |
| [PLN-04-02](PLN-04-02-recompute-on-catalogue-events.md) | Recompute on catalogue events | R1 | M | PLN-04-01, PLN-01-01 |
| [PLN-04-03](PLN-04-03-nutrition-in-responses.md) | Nutrition in dish responses (kcal/kJ, rounding, incomplete) | R1 | S | PLN-04-01, PLN-02-06 |
| [PLN-04-04](PLN-04-04-scanned-product-override.md) | Pack values of the household's scanned products | R1 | M | PLN-04-01, PLN-01-04 |
