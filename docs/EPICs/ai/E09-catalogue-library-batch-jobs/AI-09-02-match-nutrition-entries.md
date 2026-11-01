# AI-09-02 match_nutrition_entries

| | |
|---|---|
| Service | ai |
| Epic | [AI-09 Catalogue and library batch jobs](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-04-03 |
| Unblocks | KIT-06 (ingredient → nutrition entry) |

## Goal

Suggest which nutrition-table entry each catalogue ingredient should use, so calories per dish can be calculated
without AI (features.md §5: "Each catalogue ingredient is matched once to an entry in a nutrition table"; checked by
hand before use).

## Business rules

1. Batch kind `match_nutrition_entries`, Haiku 5.5, effort `medium`, prompt `match_nutrition_entries/v1`.
2. Item input: `ingredient_id`, names (`en` and local), `category`, `typical_state`, and 1–15 candidate entries found
   by kitchen (source `CIQUAL`, `BLS`, `USDA`; `entry_id`; original name; English name; food group).
3. Output: `entry_id` (or empty for none), `source`, `state_match` (`SAME`, `RAW_VS_COOKED`, `DRY_VS_COOKED`,
   `OTHER`), `alternative_entry_id` (second best or empty), `confidence` 0–1, `reason_code` (`EXACT_FOOD`,
   `CLOSEST_VARIANT`, `NO_GOOD_ENTRY`).
4. Choosing rules:
   - The entry must be the same food in the same state the recipe counts it in: recipes count ingredients as weighed
     (dry pasta, raw meat; features.md §5 rule 3), so dry pasta → "pasta, dry", not "pasta, cooked".
   - When two sources fit equally, prefer Ciqual for plain ingredients, BLS for German-specific foods, USDA for gaps
     such as Asian sauces and spices (features.md §5 "Sources").
   - A Polish ingredient uses the closest Ciqual or BLS entry until Polish tables are licensed.
   - No good entry → none; a wrong entry is worse than "calories incomplete".
5. Checks in ai: `reason_code=NO_GOOD_ENTRY` ⇔ empty `entry_id`; `alternative_entry_id` ≠ `entry_id`.
6. **Caller checks** (KIT-06): `entry_id` and `alternative_entry_id` are among the candidates sent; every suggestion
   waits for admin approval (ADM-03) before any calorie uses it, regardless of confidence.
7. The model never outputs nutrient values; kitchen reads them from the table entry.

## Contract

```proto
enum NutritionSource { NUTRITION_SOURCE_UNSPECIFIED = 0; NUTRITION_SOURCE_CIQUAL = 1; NUTRITION_SOURCE_BLS = 2; NUTRITION_SOURCE_USDA = 3; }

message NutritionCandidate {
  NutritionSource source = 1;
  string entry_id = 2;
  string name_original = 3;
  string name_en = 4;
  string food_group = 5;
}

message NutritionMatchInput {
  string ingredient_id = 1;
  string name_en = 2;
  map<string, string> local_names = 3;
  string category = 4;
  string typical_state = 5;
  repeated NutritionCandidate candidates = 6;   // 1-15
}

enum StateMatch { STATE_MATCH_UNSPECIFIED = 0; STATE_MATCH_SAME = 1; STATE_MATCH_RAW_VS_COOKED = 2; STATE_MATCH_DRY_VS_COOKED = 3; STATE_MATCH_OTHER = 4; }
enum NutritionReason { NUTRITION_REASON_UNSPECIFIED = 0; NUTRITION_REASON_EXACT_FOOD = 1; NUTRITION_REASON_CLOSEST_VARIANT = 2; NUTRITION_REASON_NO_GOOD_ENTRY = 3; }

message NutritionMatch {
  string ingredient_id = 1;
  NutritionSource source = 2;
  string entry_id = 3;
  StateMatch state_match = 4;
  string alternative_entry_id = 5;
  float confidence = 6;
  NutritionReason reason_code = 7;
}
```

## Data

No schema change.

## Errors

Per item, as AI-04-03: `ai_refused` (not retryable), `ai_truncated`, `ai_invalid_output`, `ai_item_errored`, `ai_canceled`, `ai_expired` (retryable). Whole-call errors as AI-04-01.

## Acceptance criteria

- [ ] "spaghetti (dry)" with candidates "pasta, dry, raw" and "pasta, cooked" picks the dry entry, `SAME`.
- [ ] An ingredient with no fitting candidate returns none with `NO_GOOD_ENTRY`.
- [ ] Golden set thresholds pass (AI-09-06).
