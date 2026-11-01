# PLN-10 Cooking mode and "cooked it"

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | Step-by-step cooking with timers, amounts scaled to the pot, and a "cooked it" button that reduces the fridge stock, records the history and asks for a rating (features.md "Dishes and the day plan", "Stock is reduced automatically when a dish is marked cooked"). |

## Flow

```
GetCookingSteps(dish, version, servings | plan_meal) ─> scaled ingredients, steps + timers, plates, leave-outs
MarkCooked(cooked_id, dish, servings, plan_meal?) ─> exact amounts (non-optional, chosen toppings, no staples)
   ─> KIT-07 ConsumeStock (idempotent by cooked_id, oldest first) ─> cooked_history ─> plan meal "cooked"
   ─> rating prompt?
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-10-01](PLN-10-01-cooking-steps.md) | Cooking mode payload | R1 | S | PLN-02-06, PLN-04-04, PLN-08-01 |
| [PLN-10-02](PLN-10-02-mark-cooked.md) | "Cooked it": consume stock, history, plan status | R1 | M | PLN-05-06, PLN-08-03, KIT-07 (ConsumeStock) |
| [PLN-10-03](PLN-10-03-rating-prompt.md) | Rating prompt after cooking | R1 | S | PLN-10-02, PLN-05-05 |
