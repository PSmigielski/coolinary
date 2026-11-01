# PLN-07 Rule checker for AI and library dishes

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | "AI proposes, rules check" (features.md rule 2): every AI-composed or AI-rewritten dish is checked without AI against the household's tools, the allergens, diets and dislikes of everyone eating it, the stock it claims to use, the time budget and the catalogue, before anything is shown or stored. |

## Flow

```
ai response ─> CandidateDish (PLN-07-01) ─> structure + catalogue + text (PLN-07-05)
                                         ─> tools (PLN-07-02)
                                         ─> allergens / diets / dislikes per eater (PLN-07-03)
                                         ─> stock + allowed purchases, across the whole plan (PLN-07-04)
                                         ─> Verdict { pass, violations[], warnings[], plates[] }
pass ─> stored as household AI dish        fail ─> one repair (PLN-08-04) or library fallback
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-07-01](PLN-07-01-checker-core.md) | Checker core, candidate dish and verdict | R1 | M | PLN-02-02, PLN-01-03 |
| [PLN-07-02](PLN-07-02-tool-check.md) | Tool check | R1 | S | PLN-07-01 |
| [PLN-07-03](PLN-07-03-allergen-diet-dislike-check.md) | Allergen, diet and dislike check per eater | R1 | M | PLN-07-01, PLN-02-03 |
| [PLN-07-04](PLN-07-04-stock-check.md) | Stock and allowed-ingredient check | R1 | M | PLN-07-01, PLN-06-01 |
| [PLN-07-05](PLN-07-05-ai-output-checks.md) | AI output checks, time budget, metrics, 50-dish harness | R1 | M | PLN-07-01 … PLN-07-04 |
