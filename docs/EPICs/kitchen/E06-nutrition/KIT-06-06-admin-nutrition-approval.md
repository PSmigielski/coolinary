# KIT-06-06 Admin approval of nutrition links

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-06-05, KIT-01-07, KIT-06-10 |
| Unblocks | ADM-03, PLN-04 (dishes get calories once their ingredients are approved) |

## Goal

The founder (and later a dietitian) approves, corrects or rejects each proposed link, so only checked values reach
users.

## Business rules

1. `ListNutritionSuggestions(filter, page)`: suggestions ordered by ingredient `popularity` descending; filters:
   `warning` present, source, category. Each item shows the ingredient, the chosen entry with all values, the other
   candidates, confidence and warnings. 50 per page, maximum 200.
2. `ApproveNutritionLink(ingredient_id, entry_id)`: approves this entry; any other approved link for the ingredient
   becomes `rejected`. Requires the entry to have energy.
3. `RejectNutritionLink(ingredient_id, entry_id)`: rejects; the ingredient is queued again (KIT-06-05) unless
   `requeue = false`.
4. `SetNutritionLink(ingredient_id, entry_id)`: picks any entry by hand (search with `SearchNutritionEntries(query,
   source)`) and approves it, `method = 'admin'`.
5. `CreateAdminEntry(values, name)`: an entry with `source = 'admin'` (e.g. values a dietitian provides) that can
   then be linked.
6. Every approval, rejection or change writes `admin_changes`, raises the ingredient's `nutrition_version` and
   emits `catalogue.nutrition_changed` (KIT-06-10) when the approved link changed.
7. Coverage view `NutritionCoverage()`: count and share of active ingredients with an approved link, overall and
   for ingredients used by published dishes (planner passes the list in, or admin calls planner first; ADM-03).

## Contract

`contracts/coolinary/kitchen/v1/admin.proto`

```proto
rpc ListNutritionSuggestions(ListNutritionSuggestionsRequest) returns (ListNutritionSuggestionsResponse);
rpc ApproveNutritionLink(NutritionLinkRequest) returns (NutritionLink);
rpc RejectNutritionLink(RejectNutritionLinkRequest) returns (NutritionLink);
rpc SetNutritionLink(NutritionLinkRequest) returns (NutritionLink);
rpc SearchNutritionEntries(SearchNutritionEntriesRequest) returns (SearchNutritionEntriesResponse);
rpc CreateAdminEntry(CreateAdminEntryRequest) returns (NutritionEntry);
rpc NutritionCoverage(NutritionCoverageRequest) returns (NutritionCoverageResponse);
message NutritionCoverageRequest { repeated string ingredient_ids = 1; } // optional subset, e.g. ingredients of published dishes
```

## Data

No schema change.

## Events

Emits `catalogue.nutrition_changed` v1 through KIT-06-10.

## Errors

| Case | Status | Message |
|---|---|---|
| Entry without energy | `FAILED_PRECONDITION` | `entry has no energy value` |
| Unknown ingredient or entry | `NOT_FOUND` | `not found` |

## Acceptance criteria

- [ ] Approving a link makes `ComputeNutrition` use it on the next call.
- [ ] Changing an approved link emits one `catalogue.nutrition_changed` for that ingredient.
- [ ] Coverage for a given ingredient list is reported correctly.
