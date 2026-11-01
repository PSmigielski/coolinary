# AI-11 Swaps, wishes, tool rewrite, restyle

| | |
|---|---|
| Service | ai |
| Release | R1–R3 |
| Goal | Change one part of a plan without redoing the day: swap one meal, replace a missing ingredient, follow a free-text wish ("something warm in 20 minutes"), expand a chosen alternative (R1); rewrite a recipe for the tools at hand (R2); restyle a dish into another cuisine (R3). features.md §2 "Dishes and the day plan", "Kitchen tools", "Cuisine choice". Plus only, except expanding an alternative of a plan the user already got. |

## Flow

```
planner PLN-09 -> ai.ModifyMeal(mode, context, other meals, target)   SWAP | REPLACE_INGREDIENT | WISH | EXPAND_ALTERNATIVE
planner PLN-09 -> ai.RewriteDish(mode, dish, context)                 TOOLS (R2) | RESTYLE (R3)
validation: AI-10-03 (+ rules per mode); planner rule checker after; fallback = library match or "not possible"
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-11-01](AI-11-01-modify-meal-swap-expand.md) | ModifyMeal: contract, SWAP and EXPAND_ALTERNATIVE | R1 | M | AI-10-03, AI-10-04 |
| [AI-11-02](AI-11-02-replace-ingredient.md) | ModifyMeal: REPLACE_INGREDIENT | R1 | M | AI-11-01 |
| [AI-11-03](AI-11-03-free-text-wish.md) | ModifyMeal: WISH (free text) | R1 | M | AI-11-01, AI-01-07 |
| [AI-11-04](AI-11-04-rewrite-for-tools.md) | RewriteDish: TOOLS | R2 | M | AI-11-01 |
| [AI-11-05](AI-11-05-restyle.md) | RewriteDish: RESTYLE | R3 | S | AI-11-04 |
| [AI-11-06](AI-11-06-modify-golden-sets.md) | Golden sets for swaps, wishes and rewrites | R1–R3 | M | AI-11-01 … AI-11-05, AI-03-05 |
