# PLN-12 Per-person portions (one dish, one pot)

Release R2. The family shares one dish per meal; each person gets a portion sized to their daily kcal target, and
the kitchen cooks one pot. Nobody gets a separate dish because of calories (features.md §2 "Planning for each
person", "How one meal is split").

## Flow

1. Who eats this meal: household members from the profile snapshot (KIT-11 `GetPlanningProfile`) minus people
   marked away (PLN-12-01).
2. Each person's share of their day target for this meal (PLN-12-02), plus what carried over from the previous meal
   (PLN-12-05).
3. Portion per person and the pot (PLN-12-03); people still short get no-cook extras (PLN-12-04).
4. The result is stored per planned meal and recomputed when inputs change (PLN-12-06), filtered by who is looking
   (PLN-12-07), and turned into scaled ingredients for cooking mode and the shopping list (PLN-12-08).

All arithmetic runs without AI. Only daily kcal targets are used; age, height, weight and goals never reach the
server.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-12-01](PLN-12-01-meal-attendance.md) | Meal attendance (marking people away) | R2 | M | PLN-01, PLN-08 |
| [PLN-12-02](PLN-12-02-meal-shares.md) | Meal shares per person | R2 | S | PLN-01, KIT-11 (GetPlanningProfile) |
| [PLN-12-03](PLN-12-03-portion-and-pot-calculator.md) | Portion and pot calculator | R2 | M | PLN-12-02, PLN-04 |
| [PLN-12-04](PLN-12-04-no-cook-extras.md) | No-cook extras | R2 | M | PLN-12-03, KIT-06 (ComputeNutrition), KIT-07 (ListStock), KIT-02 (GetIngredients) |
| [PLN-12-05](PLN-12-05-carry-over-and-day-balance.md) | Carry-over and day balance | R2 | M | PLN-12-03, PLN-12-04 |
| [PLN-12-06](PLN-12-06-portion-allocation-for-planned-meals.md) | Portion allocation for planned meals | R2 | M | PLN-12-01 … PLN-12-05, PLN-08 |
| [PLN-12-07](PLN-12-07-calorie-visibility.md) | Who sees which numbers | R2 | S | PLN-12-06 |
| [PLN-12-08](PLN-12-08-scaled-ingredients.md) | Scaled ingredients for a planned meal | R2 | S | PLN-12-06, PLN-02 |
