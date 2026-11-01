# PLN-11 Library growth from well-rated AI dishes

| | |
|---|---|
| Service | planner |
| Release | R2 |
| Goal | Well-rated AI dishes are reviewed and added to the library, so the same dish later costs nothing (features.md "Dishes and the day plan", once for everyone). A household's private dish is copied without any link to the household and goes through the normal review (PLN-03-07). |

## Flow

```
weekly job ─> household AI dishes rated ≥ 4 and cooked ≥ 2 times ─> cluster similar dishes across households
          ─> skip if the library already has it ─> candidate queue ─> admin promotes ─> anonymised library draft
          ─> translations, effort tags, nutrition, review ─> published
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-11-01](PLN-11-01-candidate-selection.md) | Promotion candidate job | R2 | M | PLN-05-05, PLN-05-06, PLN-01-05 |
| [PLN-11-02](PLN-11-02-anonymised-copy.md) | Anonymised copy into review | R2 | M | PLN-11-01, PLN-03-07, PLN-03-05, PLN-03-06 |
| [PLN-11-03](PLN-11-03-promotion-admin-and-metrics.md) | Promotion admin RPCs and savings metric | R2 | S | PLN-11-01, PLN-11-02, PLN-01-06 |
