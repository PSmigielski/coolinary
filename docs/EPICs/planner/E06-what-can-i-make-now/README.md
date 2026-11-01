# PLN-06 What can I make now

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | Rank library dishes by how much of each is already at home and by what expires first (features.md "Dishes and the day plan"), without AI, so the app is useful with AI switched off (rule 1, Verification 1). Extends the existing `SuggestDishes` RPC. |

## Flow

```
SuggestDishes(meal, filters) ─> profile (PLN-01-03) + stock (PLN-01-04) + library cache
   ─> hard limits (PLN-05-02) ─> coverage per dish (PLN-06-01) ─> score + tie-breaks (PLN-06-02) ─> top N with missing items
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-06-01](PLN-06-01-stock-coverage.md) | Stock coverage per dish | R1 | M | PLN-01-03, PLN-01-04, PLN-02-05, KIT-02 (GetIngredients) |
| [PLN-06-02](PLN-06-02-ranking.md) | Ranking: share on hand, expiry, tie-breaks | R1 | M | PLN-06-01, PLN-05-04, PLN-05-06 |
| [PLN-06-03](PLN-06-03-suggest-dishes-rpc.md) | Extend SuggestDishes | R1 | M | PLN-06-02, PLN-05-02, PLN-05-03 |
