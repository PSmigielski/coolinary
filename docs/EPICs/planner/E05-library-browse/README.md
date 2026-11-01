# PLN-05 Library browse, favourites, ratings, history

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | The recipe library works fully without AI (features.md rule 1): search and filter dishes by meal, cuisine, tools, diet, time, calories and protein, with the household's tools, allergens, diets and dislikes always applied as hard limits; keep favourites, ratings and a cooked history. |

## Flow

```
SearchDishes(q, filters) ──> profile snapshot (PLN-01-03) ──> hard limits (PLN-05-02) ──> soft filters + facets (PLN-05-03)
                                                                                     ──> page of dishes (BatchGetDishes shape)
favourites / ratings / history ──> per household tables ──> flags in every dish response
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-05-01](PLN-05-01-search-dishes-rpc.md) | SearchDishes RPC (text search, scopes, sorting, paging) | R1 | M | PLN-02-06, PLN-04-01 |
| [PLN-05-02](PLN-05-02-hard-limits.md) | Hard limits from the kitchen profile | R1 | M | PLN-05-01, PLN-01-03, PLN-02-03 |
| [PLN-05-03](PLN-05-03-filters-and-facets.md) | Filters and facet counts | R1 | M | PLN-05-01, PLN-04-01 |
| [PLN-05-04](PLN-05-04-favourites.md) | Favourites | R1 | S | PLN-02-01 |
| [PLN-05-05](PLN-05-05-ratings.md) | Ratings | R1 | S | PLN-02-01 |
| [PLN-05-06](PLN-05-06-cooked-history.md) | Cooked history | R1 | S | PLN-02-01 |
