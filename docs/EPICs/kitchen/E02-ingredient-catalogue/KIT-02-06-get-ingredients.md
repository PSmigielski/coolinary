# KIT-02-06 GetIngredients and catalogue version

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-02-02, KIT-02-03, KIT-02-04 |
| Unblocks | PLN-02, PLN-07, PLN-15 (planner reads flags, weights, packs), OFF-07, GW-05 |

## Goal

Other services get the full record of many ingredients in one call, and can tell cheaply whether their cached copy
is still current.

## Business rules

1. Input: up to 500 `ids`, `language`, optional `country` (for pack sizes), `include_names` (all names, default
   false).
2. For each id: an active or retired ingredient is returned as is; a merged one returns its target with
   `requested_id` set to the asked id; unknown ids are listed in `not_found` (not an error).
3. The record holds: id, slug, display name (KIT-02-04 rule 6), category, base unit, allergens, diet flags
   (KIT-02-01 rule 5), density, measures, storage and shelf lives, pack sizes for `country`, status,
   `nutrition_version`, and optionally all names.
4. `catalogue_version` is one number for the whole catalogue, raised in the same transaction as any change to
   ingredients, names, measures, storage or packs. Every response carries it.
5. `GetCatalogueVersion` returns only the number, so planner can poll it every minute and refresh its cache when it
   changes.
6. Records are cached in process per (`catalogue_version`, language, country); a version change empties the cache.

## Contract

`contracts/coolinary/kitchen/v1/catalogue.proto`

```proto
rpc GetIngredients(GetIngredientsRequest) returns (GetIngredientsResponse);
rpc GetCatalogueVersion(GetCatalogueVersionRequest) returns (GetCatalogueVersionResponse);
message GetIngredientsRequest { repeated string ids = 1; string language = 2; string country = 3; bool include_names = 4; }
message GetIngredientsResponse { repeated Ingredient ingredients = 1; repeated string not_found = 2; int64 catalogue_version = 3; }
message Ingredient {
  string id = 1; string slug = 2; string name = 3; bool translated = 4; Category category = 5;
  coolinary.common.v1.Unit base_unit = 6; repeated Allergen allergens = 7; DietFlags diet = 8;
  optional double density_g_per_ml = 9; repeated Measure measures = 10; Storage storage = 11;
  repeated PackSize pack_sizes = 12; IngredientStatus status = 13; string requested_id = 14;
  int32 nutrition_version = 15; repeated IngredientName names = 16; bool staple_candidate = 17;
}
message IngredientName { string language = 1; string name = 2; string kind = 3; bool reviewed = 4; }
message GetCatalogueVersionResponse { int64 catalogue_version = 1; }
```

## Data

`kitchen/migrations/000NN_catalogue_version.sql`

```sql
CREATE TABLE catalogue_version (id boolean PRIMARY KEY DEFAULT true CHECK (id), version bigint NOT NULL);
INSERT INTO catalogue_version VALUES (true, 1);
```

## Errors

| Case | Status | Message |
|---|---|---|
| More than 500 ids | `INVALID_ARGUMENT` | `at most 500 ids` |
| An id that isn't a UUID | listed in `not_found` | — |

## Acceptance criteria

- [ ] 500 ids return in ≤ 40 ms p95 from a warm cache.
- [ ] A merged id returns the target with `requested_id`.
- [ ] Editing a measure raises `catalogue_version`; the next call reflects the change.
