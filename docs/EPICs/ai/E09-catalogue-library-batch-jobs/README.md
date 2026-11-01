# AI-09 Catalogue and library batch jobs

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | The once-for-everyone AI work behind the catalogue and the recipe library (features.md §2 "once for everyone" items, §5 nutrition matching, decision 10 starter dishes) runs as batch kinds at 50% cost, and every result is a proposal a human approves in admin. |

## Flow

```
kitchen KIT-02  -> translate_ingredients     -> names + synonyms + search forms per language -> admin spot-check
kitchen KIT-06  -> match_nutrition_entries   -> Ciqual / BLS / USDA entry per ingredient      -> admin approves each
planner PLN-03  -> draft_dishes (Sonnet)     -> dish drafts                                  -> catalogue check, calories (KIT-06), admin review
planner PLN-03  -> translate_dishes          -> names and steps in pl / de / en               -> admin review
planner PLN-03  -> suggest_effort_tags       -> effort tags + keeps/freezes                   -> admin checks
```

All via SubmitBatch / GetBatch / GetBatchResults (AI-04). In batch mode ai checks the schema and input-independent
rules; each task lists the checks the caller must run because it holds the inputs (AI-04-03 rule 4).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-09-01](AI-09-01-translate-ingredients.md) | translate_ingredients | R1 | M | AI-04-03 |
| [AI-09-02](AI-09-02-match-nutrition-entries.md) | match_nutrition_entries | R1 | M | AI-04-03 |
| [AI-09-03](AI-09-03-draft-dishes.md) | draft_dishes | R1 | L | AI-04-03, AI-10-01 |
| [AI-09-04](AI-09-04-translate-dishes.md) | translate_dishes | R1 | M | AI-09-03 |
| [AI-09-05](AI-09-05-suggest-effort-tags.md) | suggest_effort_tags | R1 | S | AI-04-03, AI-10-01 |
| [AI-09-06](AI-09-06-batch-golden-sets.md) | Golden sets for catalogue and library jobs | R1 | M | AI-09-01 … AI-09-05, AI-03-03 |
