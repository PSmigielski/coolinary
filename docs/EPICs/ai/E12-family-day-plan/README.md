# AI-12 Family day plan

| | |
|---|---|
| Service | ai |
| Release | R2 |
| Goal | Compose a family day with one shared dish per meal that fits every eater's hard limits and suits their different calorie targets, from stock first; a second dish only when an allergy or diet forces it (features.md §2 "Calories and goals", "Planning for each person"; Plus). Portion sizes are calculated by planner without AI. |

## Flow

```
planner PLN-13: per eater (opaque ref, no name, no age): personal diets/allergens/dislikes, meal targets in kcal, away flags
  -> ai.ComposeFamilyDay (Sonnet 5.5 medium)
     validation: every present eater's hard limits; toppings per eater; second dish only when forced
  <- per meal: dish, who eats it, toppings per eater, optional second dish
planner PLN-12 sizes each portion (½-2, quarter steps), the pot, no-cook extras; if a dish can't fit someone, it asks for a swap
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-12-01](AI-12-01-family-contract.md) | ComposeFamilyDay contract | R2 | S | AI-10-01 |
| [AI-12-02](AI-12-02-family-compose.md) | Family composition prompt and call | R2 | L | AI-12-01, AI-10-02 |
| [AI-12-03](AI-12-03-family-validation.md) | Validation for eaters, toppings and second dishes | R2 | M | AI-12-02, AI-10-03 |
| [AI-12-04](AI-12-04-family-golden-set.md) | Golden set: families including the worked example | R2 | M | AI-12-03, AI-03-05 |
