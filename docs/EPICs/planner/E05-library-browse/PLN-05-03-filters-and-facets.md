# PLN-05-03 Filters and facet counts

| | |
|---|---|
| Service | planner |
| Epic | [PLN-05 Library browse](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-05-01, PLN-04-01 |
| Unblocks | GW-07, PLN-06-03 |

## Goal

The library filters from features.md ("Dishes and the day plan", "Calories and goals", "Kitchen tools") with counts per
option, so users see how many dishes each choice leaves.

## Business rules

1. Filters (all optional, combined with AND; values inside one filter with OR):
   - `meals`: breakfast, dinner, supper.
   - `cuisines`: cuisine codes; `asian` matches every Asian sub-cuisine; a sub-cuisine matches only itself.
   - `tools_only` (one-off "air fryer only" for tonight): dish tools ⊆ the selected tools (and still ⊆ the household's).
     Dishes with no tool match.
   - `max_total_minutes`: 5–240.
   - `kcal_bands`: `UNDER_400` (< 400), `FROM_400_TO_600` (400–600 inclusive), `OVER_600` (> 600), on `kcal_rounded`
     per portion; dishes with incomplete calories never match a band.
   - `protein_25g`: protein per portion ≥ 25.0 g (unrounded).
   - `effort_tags`: dish has all selected tags.
   - `ingredient_ids`: dish contains all of them in non-optional lines.
2. Labels are numeric only; no filter or tag is called "healthy" or "light" (features.md "Calories and goals").
3. Facets: counts for meals, cuisines, kcal bands, protein_25g, effort tags and time bands (≤ 15, ≤ 30, ≤ 45, ≤ 60,
   > 60 min). Each facet's counts apply all other filters but not its own (standard "multi-select facet" rule).
4. Facets are computed only on the first page (no `page_token`) to keep paging cheap.

## Contract

```proto
message DishFilters {
  repeated Meal meals = 1;
  repeated coolinary.common.v1.Cuisine cuisines = 2;
  repeated coolinary.common.v1.KitchenTool tools_only = 3;
  int32 max_total_minutes = 4;
  repeated KcalBand kcal_bands = 5;
  bool protein_25g = 6;
  repeated coolinary.common.v1.EffortTag effort_tags = 7;
  repeated string ingredient_ids = 8;
}
enum KcalBand { KCAL_BAND_UNSPECIFIED = 0; KCAL_BAND_UNDER_400 = 1; KCAL_BAND_FROM_400_TO_600 = 2; KCAL_BAND_OVER_600 = 3; }
message FacetCount { string value = 1; int32 count = 2; }
message DishFacets {
  repeated FacetCount meals = 1;
  repeated FacetCount cuisines = 2;
  repeated FacetCount kcal_bands = 3;
  int32 protein_25g = 4;
  repeated FacetCount effort_tags = 5;
  repeated FacetCount time_bands = 6;
}
```

## Data

`library_search` gains `meals smallint[]`, `cuisine text`, `sub_cuisine text`, `total_minutes int`, `kcal_rounded int`,
`protein_g numeric`, `effort_tags text[]`, indexed.

## Errors

| Case | gRPC status | Message |
|---|---|---|
| `max_total_minutes` outside 5–240 | INVALID_ARGUMENT | limits |

## Acceptance criteria

- [ ] A 400 kcal dish matches `FROM_400_TO_600`, not `UNDER_400`; a 600 kcal dish matches `FROM_400_TO_600`.
- [ ] A dish with 24.9 g protein doesn't match `protein_25g`.
- [ ] Selecting cuisine `polish` leaves the cuisine facet counts unchanged and changes the meal facet counts.
- [ ] `tools_only=[air_fryer]` returns no oven dish even if the household has an oven.
