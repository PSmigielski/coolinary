# PLN-03 Library import, AI drafting and review

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | Fill the library with about 200 checked starter dishes in Polish, German, Italian, Asian and Mediterranean cuisine (features.md §7, release 1), from two sources that end in one review queue: a hand-written import file and an AI drafting batch job (decision 10). Translations and effort tags are AI-suggested once for everyone and checked by a person. |

## Flow

```
YAML import file ──> parse + validate ──> name matching (KIT-03) ──┐
                                                                   ├─> dish (status review) ──> admin review ──> published
admin briefs ──> ai batch draft_dishes (AI-04) ──> poll ──> same ──┘          ^        ^
                                                          translate_dishes ───┘        │
                                                          suggest_effort_tags ─────────┘
dietitian sign-off (per version) ──> launch gate report;  coverage report ──> "do we have enough dishes?"
```

All RPCs here are on `PlannerAdminService` and are called by admin (ADM-04).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-03-01](PLN-03-01-import-file-format.md) | Import file format and parser | R1 | M | PLN-02-02 |
| [PLN-03-02](PLN-03-02-ingredient-name-resolution.md) | Ingredient name resolution | R1 | M | PLN-01-02, KIT-03 (MatchNames) |
| [PLN-03-03](PLN-03-03-import-dishes-rpc.md) | ImportDishes admin RPC | R1 | M | PLN-03-01, PLN-03-02, PLN-02-04, PLN-01-06 |
| [PLN-03-04](PLN-03-04-ai-drafting-batch.md) | AI drafting batch (`draft_dishes`) | R1 | L | PLN-03-02, PLN-02-02, AI-04 (SubmitBatch), AI-09 |
| [PLN-03-05](PLN-03-05-dish-translation-batch.md) | Dish translation batch (`translate_dishes`) | R1 | M | PLN-03-04 (batch tracking), AI-04, AI-09 |
| [PLN-03-06](PLN-03-06-effort-tag-suggestions.md) | Effort tag suggestions (`suggest_effort_tags`) | R1 | S | PLN-03-04 (batch tracking), AI-04, AI-09 |
| [PLN-03-07](PLN-03-07-review-queue.md) | Review queue and decisions | R1 | L | PLN-02-04, PLN-02-03, PLN-04-01, PLN-01-06 |
| [PLN-03-08](PLN-03-08-dietitian-sign-off.md) | Dietitian sign-off | R1 | S | PLN-03-07, PLN-04-01 |
| [PLN-03-09](PLN-03-09-library-coverage-report.md) | Library coverage report | R1 | M | PLN-02-03, PLN-05-02 |
