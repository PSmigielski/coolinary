# AI-09-04 translate_dishes

| | |
|---|---|
| Service | ai |
| Epic | [AI-09 Catalogue and library batch jobs](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-09-03 |
| Unblocks | PLN-03 (library in every launch language) |

## Goal

Every library dish reads naturally in Polish, German and English: name, ingredient notes and steps, with numbers and
timers unchanged.

## Business rules

1. Batch kind `translate_dishes`, Haiku 5.5, effort `low`, prompt `translate_dishes/v1`.
2. Item input: `dish_ref`, `source_language`, `name`, `original_name`, `ingredient_notes` (list, may be empty
   strings), `display_hints` (list), `steps` (texts), `target_languages` (1–5).
3. Output per target language: `name` (keep `original_name` when set; otherwise translate; ≤ 80 chars),
   `ingredient_notes` and `display_hints` in the same order and count, `steps` in the same order and count.
4. Rules: keep every number, temperature and duration exactly ("180 °C", "20 min"); use the target language's
   kitchen terms and metric units; convert nothing; no added advice or claims.
5. Checks in ai: calorie guard and claim words (AI-09-03 rule 6); lengths.
6. **Caller checks** (PLN-03): counts equal the source; the multiset of numbers in each step equals the source
   step's (a regex extracts numbers); otherwise the item is invalid and retried once, then goes to admin review.

## Contract

```proto
message DishTranslationInput {
  string dish_ref = 1;
  string source_language = 2;
  string name = 3;
  string original_name = 4;
  repeated string ingredient_notes = 5;
  repeated string display_hints = 6;
  repeated string steps = 7;
  repeated string target_languages = 8;
}

message DishText {
  string language = 1;
  string name = 2;
  repeated string ingredient_notes = 3;
  repeated string display_hints = 4;
  repeated string steps = 5;
}

message DishTranslation {
  string dish_ref = 1;
  repeated DishText texts = 2;
}
```

## Data

No schema change.

## Errors

Per item, as AI-04-03: `ai_refused` (not retryable), `ai_truncated`, `ai_invalid_output`, `ai_item_errored`, `ai_canceled`, `ai_expired` (retryable). Whole-call errors as AI-04-01.

## Acceptance criteria

- [ ] A 6-step English dish translates to pl and de with 6 steps each and identical numbers per step.
- [ ] "Pierogi ruskie" keeps its name in de and en.
- [ ] Golden set thresholds pass (AI-09-06).
