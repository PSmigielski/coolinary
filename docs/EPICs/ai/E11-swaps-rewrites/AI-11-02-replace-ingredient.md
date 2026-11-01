# AI-11-02 ModifyMeal: REPLACE_INGREDIENT

| | |
|---|---|
| Service | ai |
| Epic | [AI-11 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-11-01 |
| Unblocks | PLN-09 |

## Goal

When one ingredient of a planned dish is missing, suggest up to three replacements from what is at home, so the user
can still cook the dish.

## Business rules

1. Mode `REPLACE_INGREDIENT`, job `modify_meal`.
2. Input: `current_meal` (library id or new dish, with its stock uses) and `missing_ingredient_id` (must be an
   ingredient of that dish).
3. Output: 1–3 `IngredientReplacement`s, best first. Each: `ingredient_id` from stock, staples or `buyable_ids`;
   `amount`; `source` (`STOCK`, `STAPLE`, `BUY`); `note` ≤ 80 chars ("use half the amount, it's saltier");
   `step_changes` (index + new text) only where a step must change; `drop_only=true` when the dish works without the
   ingredient and nothing is needed.
4. Replacement rules: same culinary role (fat for fat, acid for acid, binder for binder); keeps every hard limit
   (allergens, diets, dislikes, tools); stock replacements must fit the remaining stock after the other meals; prefer
   `STOCK` over `STAPLE` over `BUY`; at most one `BUY` suggestion.
5. Validation: ids as rule 4; hard limits via facts (AI-10-03 rules 5–7); `step_changes` index < number of steps;
   calorie guard on notes and steps.
6. No replacement possible → empty list with `none_possible=true` (OK, not an error); planner shows "swap the meal
   instead".

## Contract

```proto
enum ReplacementSource { REPLACEMENT_SOURCE_UNSPECIFIED = 0; REPLACEMENT_SOURCE_STOCK = 1;
  REPLACEMENT_SOURCE_STAPLE = 2; REPLACEMENT_SOURCE_BUY = 3; }

message StepChange { int32 index = 1; string text = 2; }

message IngredientReplacement {
  string ingredient_id = 1;
  coolinary.common.v1.Quantity amount = 2;
  ReplacementSource source = 3;
  string note = 4;
  repeated StepChange step_changes = 5;
  bool drop_only = 6;
}
```

`ModifyMealResponse` gets `bool none_possible = 4;`.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Input limits broken, unknown ids or codes in the request | `INVALID_ARGUMENT`, `ai_bad_request` |
| Output still invalid after repair | `UNAVAILABLE`, `ai_invalid_output` |
| Others | AI-01-04; the caller falls back to its non-AI result |

## Acceptance criteria

- [ ] Missing butter in a lactose-free household never suggests a dairy product.
- [ ] Missing lemon with vinegar in stock suggests vinegar with an amount note.
- [ ] A replacement id not at home and not in `buyable_ids` fails validation.
