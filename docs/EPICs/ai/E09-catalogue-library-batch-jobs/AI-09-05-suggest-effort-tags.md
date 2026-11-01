# AI-09-05 suggest_effort_tags

| | |
|---|---|
| Service | ai |
| Epic | [AI-09 Catalogue and library batch jobs](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-04-03, AI-10-01 |
| Unblocks | PLN-03, PLN-14 (leftovers and cook-once plans need the tags) |

## Goal

Suggest effort tags for library dishes — no cooking, keeps 3 days in the fridge, freezes well, good cold the next
day, base plus toppings — for a human to confirm (features.md §2 "Calories and goals", effort tags; "checked by hand
before use").

## Business rules

1. Batch kind `suggest_effort_tags`, Haiku 5.5, effort `low`, prompt `suggest_effort_tags/v1`.
2. Item input: `dish_ref`, name, ingredients (English names with roles), steps, `total_minutes`, tool codes.
3. Output: for each of the five tags a boolean and a confidence 0–1; `keeps_days` (0–5, days safe in the fridge after
   cooking); `freezes` boolean.
4. Rules: `NO_COOKING` only if no step heats food (toaster and kettle count as cooking); `KEEPS_3_DAYS` only when
   `keeps_days ≥ 3`; dishes with raw fish, soft-boiled eggs or fresh salad leaves dressed never get `KEEPS_3_DAYS`;
   be conservative: food safety beats convenience (business model §10 "conservative expiry guidance").
5. Checks in ai: `KEEPS_3_DAYS` ⇔ `keeps_days ≥ 3`; `FREEZES_WELL` ⇒ `freezes`.
6. **Caller checks** (PLN-03): none beyond storing; every tag is confirmed by admin before planner uses it.

## Contract

```proto
message EffortIngredient {
  string name_en = 1;              // catalogue name; ids mean nothing to the model here
  IngredientRole role = 2;         // AI-10-01
}

message EffortTagInput {
  string dish_ref = 1;
  string name = 2;
  repeated EffortIngredient ingredients = 3;
  repeated Step steps = 4;
  int32 total_minutes = 5;
  repeated coolinary.common.v1.KitchenTool tools = 6;
}

message TagScore { coolinary.common.v1.EffortTag tag = 1; bool value = 2; float confidence = 3; }

message EffortTagSuggestion {
  string dish_ref = 1;
  repeated TagScore tags = 2;   // exactly 5
  int32 keeps_days = 3;
  bool freezes = 4;
}
```

## Data

No schema change.

## Errors

Per item, as AI-04-03: `ai_refused` (not retryable), `ai_truncated`, `ai_invalid_output`, `ai_item_errored`, `ai_canceled`, `ai_expired` (retryable). Whole-call errors as AI-04-01.

## Acceptance criteria

- [ ] Overnight oats → `NO_COOKING`; chili con carne → `KEEPS_3_DAYS`, `FREEZES_WELL`; dressed green salad → not `KEEPS_3_DAYS`.
- [ ] Golden set thresholds pass (AI-09-06).
