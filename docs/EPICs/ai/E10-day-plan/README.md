# AI-10 Day plan composition and cuisine choices

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | Compose a full day (breakfast, dinner, supper) from the household's stock, oldest items first, so the three meals don't compete for the same ingredients, within the kitchen's tools, diets, allergens and time; and when stock doesn't fit the chosen cuisine, offer three choices (features.md §2 "Dishes and the day plan", "Cuisine choice"). |

## Flow

```
gateway (quota: 1 AI day plan) -> planner PLN-08
  planner gathers: profile (kitchen GetPlanningProfile), stock (ListStock), ingredient facts (GetIngredients),
                   up to 30 ranked library candidates (PLN-06), buyable ids (cuisine base pantry, PLN-16)
  -> ai.ComposeDayPlan (Sonnet 5.5 medium; Free routing see AI-10-05)
     validation in ai (AI-10-03): ids, allocation, tools, allergens, diets, dislikes, time, calorie guard; 1 repair
  <- meals: library dish id or new dish, stock uses, items to buy, cuisine fit + alternatives
  planner rule checker (PLN-07) is the authority; calories by kitchen tables; fallback = library matches
```

## Shared types

AI-10-01 defines `NewDish`, `DishIngredient`, `Step`, `IngredientFact`, `StockLine`, `KitchenRules`,
`LibraryCandidate`, `ComposedMeal` and the enums `MealSlot`, `IngredientRole` (`Diet`, `EffortTag` and the other shared enums come from `coolinary.common.v1`, decision 16). AI-09, AI-11,
AI-12, AI-14, AI-15 and AI-17 reuse them.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-10-01](AI-10-01-composition-contract.md) | Composition contract and shared dish types | R1 | M | AI-02-01 |
| [AI-10-02](AI-10-02-compose-day-plan.md) | ComposeDayPlan prompt and call | R1 | L | AI-10-01, AI-01-03, AI-01-05, AI-01-09 |
| [AI-10-03](AI-10-03-plan-validation.md) | Validation and repair of composed meals | R1 | M | AI-10-02, AI-01-06 |
| [AI-10-04](AI-10-04-cuisine-choices.md) | Cuisine fit and the three choices | R1 | M | AI-10-03 |
| [AI-10-05](AI-10-05-model-routing-by-plan.md) | Model routing by plan and the cost decision | R1 | S | AI-10-06, AI-02-05 |
| [AI-10-06](AI-10-06-day-plan-golden-set.md) | Golden set and eval for day plans | R1 | M | AI-10-04, AI-03-05 |
