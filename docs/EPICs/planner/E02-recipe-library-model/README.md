# PLN-02 Recipe library model and storage

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | Store every dish — library dishes and a household's own AI dishes — in one versioned model with translations, ingredients by kitchen ingredient id, steps with timers, tools, effort tags and derived diet and allergen flags. Replaces the skeleton `dishes` table. |

## Flow

```
import / AI draft / household AI  ->  draft version  ->  review  ->  published version (immutable)
                                                              \->  retired (hidden, still readable by plans and history)
dish_ingredients.ingredient_id -> kitchen catalogue (GetIngredients) -> derived diet/allergen flags
```

Plans, history and ratings point at `(dish_id, version)`, so a later edit never changes a plan already made.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-02-01](PLN-02-01-dish-schema.md) | Dish schema (replaces skeleton `dishes`) | R1 | M | — |
| [PLN-02-02](PLN-02-02-dish-validation.md) | Dish validation rules | R1 | M | PLN-02-01 |
| [PLN-02-03](PLN-02-03-derived-flags.md) | Derived diet and allergen flags | R1 | M | PLN-02-01, PLN-01-02, KIT-02 (GetIngredients) |
| [PLN-02-04](PLN-02-04-dish-lifecycle.md) | Dish lifecycle and versions | R1 | M | PLN-02-01, PLN-02-02 |
| [PLN-02-05](PLN-02-05-portion-scaling.md) | Portion scaling rules | R1 | S | PLN-02-01, KIT-02 (GetIngredients) |
| [PLN-02-06](PLN-02-06-get-dish.md) | GetDish and BatchGetDishes | R1 | M | PLN-02-01, PLN-02-04, PLN-02-05 |
