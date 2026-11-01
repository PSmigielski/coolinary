# PLN-05-01 SearchDishes RPC (text search, scopes, sorting, paging)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-05 Library browse](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-02-06, PLN-04-01 |
| Unblocks | PLN-05-02, PLN-05-03, GW-07 (library routes) |

## Goal

One RPC behind the library screen: free-text search over dish names and ingredients, three scopes (library, my AI
dishes, favourites), stable sorting and paging.

## Business rules

1. Scopes: `LIBRARY` (published library dishes, default), `MINE` (this household's `household_ai` dishes),
   `FAVOURITES` (PLN-05-04). Retired dishes appear only in `FAVOURITES`, marked unavailable.
2. Text `q` (≤ 100 chars): matched against dish names in the request language and `en` with pg_trgm similarity
   ≥ 0.3, and against ingredient display names (a dish matches if any non-optional ingredient's name matches). Name
   matches rank above ingredient matches.
3. Sort orders: `RELEVANCE` (default when `q` is set), `POPULAR` (default otherwise), `QUICKEST` (total minutes asc),
   `KCAL_ASC`, `KCAL_DESC`, `NEWEST`. Dishes with incomplete calories sort last under the kcal orders.
4. `POPULAR` score = Bayesian average of ratings across all households: (C × m + Σ ratings) / (C + n), with m = 3.5
   and C = 5. Ties → name in the request language.
5. Paging: `page_size` 1–50 (default 20), opaque `page_token` encoding the sort key and dish id; results are stable while
   paging even if new dishes are published.
6. Each item has the list shape of `BatchGetDishes` (no steps) plus `is_favourite`, `my_rating`, `cooked_count` for this
   household.
7. Hard limits (PLN-05-02) are always applied; soft filters (PLN-05-03) when given.
8. Works with ai switched off; makes no ai call.

## Contract

```proto
rpc SearchDishes(SearchDishesRequest) returns (SearchDishesResponse);

message SearchDishesRequest {
  string household_id = 1;
  string language = 2;
  string q = 3;
  DishScope scope = 4;
  DishSort sort = 5;
  DishFilters filters = 6;                      // PLN-05-03
  repeated string eater_member_ids = 7;         // PLN-05-02; empty = everyone
  coolinary.common.v1.PageRequest page = 8;
}
message DishListItem {
  Dish dish = 1;                                // without steps
  bool is_favourite = 2;
  int32 my_rating = 3;
  int32 cooked_count = 4;
  bool available = 5;                           // false when retired
  repeated string leave_out_ingredient_ids = 6; // PLN-05-02 rule 5
}
message SearchDishesResponse {
  repeated DishListItem items = 1;
  DishFacets facets = 2;                        // PLN-05-03
  coolinary.common.v1.PageResponse page = 3;
  int32 total = 4;
}
enum DishScope { DISH_SCOPE_UNSPECIFIED = 0; DISH_SCOPE_LIBRARY = 1; DISH_SCOPE_MINE = 2; DISH_SCOPE_FAVOURITES = 3; }
enum DishSort { DISH_SORT_UNSPECIFIED = 0; DISH_SORT_RELEVANCE = 1; DISH_SORT_POPULAR = 2; DISH_SORT_QUICKEST = 3;
                DISH_SORT_KCAL_ASC = 4; DISH_SORT_KCAL_DESC = 5; DISH_SORT_NEWEST = 6; }
```

## Data

`planner/migrations/000NN_search_indexes.sql` — a materialised `library_search` view or table refreshed on publish /
retire / text change, holding per (dish, language) the current version's name, ingredient names, tools, meals,
cuisine, diets_capable, allergens, total minutes, kcal and protein, so search is one indexed query.

```sql
CREATE TABLE library_search (
    dish_id        uuid NOT NULL,
    language       text NOT NULL,
    version        int  NOT NULL,
    name           text NOT NULL,
    ingredient_text text NOT NULL,
    PRIMARY KEY (dish_id, language)
);
CREATE INDEX library_search_name_trgm ON library_search USING gin (name gin_trgm_ops);
CREATE INDEX library_search_ingr_trgm ON library_search USING gin (ingredient_text gin_trgm_ops);
```

## Errors

| Case | gRPC status | Message |
|---|---|---|
| `q` > 100 chars, page size > 50 | INVALID_ARGUMENT | limits |
| Profile unavailable | UNAVAILABLE | from PLN-01-03 (never unfiltered results) |

## Acceptance criteria

- [ ] "curry" finds "Chicken curry with rice"; "kurczak" (pl) finds it by its Polish name.
- [ ] "chickpeas" finds dishes containing chickpeas below dishes named "Chickpea …".
- [ ] Paging through 60 dishes with page size 20 returns each dish exactly once even if one is published mid-way.
- [ ] With the ai container stopped, search answers normally (Verification 1).
