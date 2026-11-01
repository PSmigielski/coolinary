# PLN-02-02 Dish validation rules

| | |
|---|---|
| Service | planner |
| Epic | [PLN-02 Recipe library model](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-02-01 |
| Unblocks | PLN-02-04, PLN-03-01, PLN-03-04, PLN-07-05 |

## Goal

One validator every write path uses (import, AI draft, admin edit, household AI dish), so no broken dish reaches the
database whatever its source.

## Business rules

Structural rules (all sources; a failure rejects the write):

1. 1–40 ingredient lines; 1–30 steps; 1–3 meals; base servings 1–12.
2. Amount > 0 and ≤ 10,000 for grams and millilitres, ≤ 100 for pieces, ≤ 20 for spoons.
3. Every `ingredient_id` exists in kitchen (`KIT-02 GetIngredients`). Unknown ids fail with their positions.
4. The same ingredient id may appear twice only with different roles or topping groups.
5. Every tool used in a step is listed in the dish's tools. A dish lists no tool its steps don't use, except `kettle`,
   `hob` and `oven` when a step text needs them without a timer (the validator then only warns).
6. `active_minutes` ≤ `total_minutes`. The sum of step timers ≤ `total_minutes × 60`.
7. Cuisine is a known code; sub-cuisine only with `asian`.
8. Effort tag `no_cook`: tools ⊆ {`kettle`, `toaster`, `blender`, `hand_blender`, `microwave`}, no step timer over
   600 s. Decided (Q087): a dish needing only a microwave, kettle, toaster or blender counts as "no cooking"
   (features.md: breakfast and the non-main meal come from "no cooking" dishes; the student preset is microwave + kettle).
9. Effort tag `base_plus_toppings`: at least one `base` line and at least two distinct `topping_group`s, each with at
   least two options.
10. Texts: name ≤ 120 characters, step text count = step count, each step 1–600 characters. No number followed by
    `kcal`, `kJ` or `calories` (in any launch language) in any text: AI never states calories (features.md rule 2).

Publication rules (checked on the move to `published`, PLN-02-04):

11. Texts exist in all launch languages `pl`, `de`, `en`.
12. Nutrition is complete (PLN-04-01), or the dish is explicitly published with "calories incomplete" by an admin
    who confirms it. Decided (Q086): the 200 starter dishes must have complete calories to be published; later
    dishes may be published with "calories incomplete" only after an admin confirms it.
13. Derived flags are computed (PLN-02-03).

The validator returns all problems at once, each with `code`, `path` (e.g. `ingredients[3].amount`) and a short
English message; admin shows them next to the field.

## Contract

`contracts/coolinary/planner/v1/dish.proto`

```proto
message ValidationProblem {
  string code = 1;      // e.g. "UNKNOWN_INGREDIENT", "STEP_TOOL_NOT_LISTED", "CALORIES_IN_TEXT"
  string path = 2;
  string message = 3;
  bool blocking = 4;    // false = warning
}
```

## Data

No schema change.

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Any blocking problem on write | INVALID_ARGUMENT | `dish invalid`; details carry the `ValidationProblem` list |
| Kitchen unavailable during check | UNAVAILABLE | `catalogue unavailable` (never skip rule 3) |

## Acceptance criteria

- [ ] Table-driven tests cover every rule above with one passing and one failing case.
- [ ] A dish text "about 520 kcal" (en), "ok. 520 kcal" (pl), "ca. 520 Kalorien" (de) fails with `CALORIES_IN_TEXT`.
- [ ] A `base_plus_toppings` dish with one topping group fails; with two groups of two options passes.
- [ ] A dish with three problems returns all three in one response.
