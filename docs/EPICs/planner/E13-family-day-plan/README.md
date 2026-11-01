# PLN-13 Family day plan

Release R2, Plus only. One day (breakfast, dinner, supper) for the whole household: one dish per meal that meets
every eater's allergens and diet, a portion per person, as few cooking sessions as possible, stock first
(features.md §2 "Calories and goals" with AI, "One dish per meal, few cooking sessions").

## Flow

1. Combine the hard limits of everyone eating each meal (PLN-13-01).
2. Compose: with AI through `AI-12 (ComposeFamilyDay)` (PLN-13-07), or without AI from the library (PLN-13-06).
   Both use the same plan comparator (PLN-13-03) with "eat lighter" (PLN-13-04), base-plus-toppings plates
   (PLN-13-02) and, only when an allergy or diet forces it, a second dish (PLN-13-05).
3. Every AI dish passes the rule checker (PLN-07) before it is kept; portions come from PLN-12.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-13-01](PLN-13-01-shared-meal-limits.md) | Hard limits for a shared meal | R2 | S | PLN-07, KIT-11 (GetPlanningProfile) |
| [PLN-13-02](PLN-13-02-base-plus-toppings.md) | Base-plus-toppings plates | R2 | M | PLN-13-01, PLN-02, PLN-12-08 |
| [PLN-13-03](PLN-13-03-plan-comparator.md) | Plan comparator (priority order) | R2 | M | PLN-12-05 |
| [PLN-13-04](PLN-13-04-eat-lighter-as-a-family.md) | Eat lighter as a family | R2 | S | PLN-13-03, PLN-12-04, KIT-02 (GetIngredients) |
| [PLN-13-05](PLN-13-05-second-dish-when-forced.md) | Second dish only when an allergy or diet forces it | R2 | S | PLN-13-01 |
| [PLN-13-06](PLN-13-06-non-ai-family-day.md) | Family day without AI | R2 | L | PLN-13-01 … PLN-13-05, PLN-06, PLN-12-06 |
| [PLN-13-07](PLN-13-07-ai-family-day.md) | Family day with AI | R2 | L | PLN-13-06, AI-12 (ComposeFamilyDay), PLN-07, PLN-08 |
