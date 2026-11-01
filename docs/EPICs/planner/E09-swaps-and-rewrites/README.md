# PLN-09 Swaps, wishes, tool rewrite, restyle

| | |
|---|---|
| Service | planner |
| Release | R1–R3 |
| Goal | Let people change a plan without starting over: swap one meal, replace a missing ingredient, ask for "something warm in 20 minutes" (R1, D12), rewrite a recipe for the tools at hand (R2, D12) and restyle a dish into another cuisine (R3). AI changes are Plus-only and checked like every AI dish; Free swaps use the library. |

## Flow

```
SwapMeal / ReplaceIngredient / WishMeal ──> Plus? ──yes──> AI-11 ModifyMeal (other meals' stock fixed) ──> PLN-07 checker
                                          └─no───> library: next best dish (swap) / FAILED_PRECONDITION plus_required
RewriteDish(tools | restyle) ──> Plus ──> AI-11 RewriteDish ──> checker ──> new household AI dish (parent = original)
```

The gateway blocks Plus-only calls for Free households (GW-03); planner checks the plan again with identity (PLN-01-05, Q004).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-09-01](PLN-09-01-swap-meal.md) | Swap one meal | R1 | M | PLN-08-07, AI-11 (ModifyMeal) |
| [PLN-09-02](PLN-09-02-replace-ingredient.md) | Replace a missing ingredient | R1 | M | PLN-09-01 |
| [PLN-09-03](PLN-09-03-free-text-wish.md) | Free-text wish for a meal | R1 | M | PLN-09-01 |
| [PLN-09-04](PLN-09-04-tool-rewrite.md) | Rewrite a recipe for my tools | R2 | M | PLN-07-05, PLN-02-04, AI-11 (RewriteDish) |
| [PLN-09-05](PLN-09-05-restyle.md) | Restyle a dish ("make this Italian-style") | R3 | M | PLN-09-04 |
