# KIT-06-01 Nutrition entries, links and sources

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-02-01 |
| Unblocks | every other KIT-06 task |

## Goal

Store nutrition table entries from all sources in one shape, the one approved link per ingredient, and the credits
each source requires.

## Business rules

1. An entry holds values per 100 g: energy (kcal and kJ), protein, carbohydrate, sugars, fat, saturates, fibre,
   salt (features.md §5 step 1). Values are grams except energy.
2. Missing value = `NULL` (unknown), never 0.
3. Decided (Q082): "traces" is stored as 0 and "< x" as x / 2, each with a qualifier kept in `qualifiers`
   (common practice; keeps the numbers conservative).
4. Energy: if a source gives only kcal or only kJ, the other is computed (1 kcal = 4.184 kJ). An entry without
   energy can't be approved for an ingredient.
5. Salt: if a source gives sodium but not salt, salt = sodium × 2.5 (Regulation (EU) 1169/2011, Annex I).
6. `state` records how the food was measured: `raw`, `dry`, `cooked`, `prepared`, `other`, from the source's
   description; matching uses it (KIT-06-05 rule 4).
7. One approved link per ingredient at most (`ingredient_nutrition.status = 'approved'`); any number of suggested
   or rejected ones.
8. `nutrition_sources` lists each source with its licence, credit text, URL and version, for the app's methodology
   page (`ListNutritionSources`, needed for app store reviews of health accuracy claims).

## Contract

`contracts/coolinary/kitchen/v1/nutrition.proto`

```proto
message Nutrients { optional double energy_kcal = 1; optional double energy_kj = 2; optional double protein_g = 3;
  optional double carbohydrate_g = 4; optional double sugars_g = 5; optional double fat_g = 6;
  optional double saturates_g = 7; optional double fibre_g = 8; optional double salt_g = 9; }
rpc ListNutritionSources(ListNutritionSourcesRequest) returns (ListNutritionSourcesResponse);
message NutritionSource { string code = 1; string name = 2; string version = 3; string licence = 4; string credit = 5; string url = 6; }
```

## Data

`kitchen/migrations/000NN_nutrition.sql`

```sql
CREATE TABLE nutrition_sources (
    code      text PRIMARY KEY CHECK (code IN ('ciqual','bls','usda','off_pack','admin')),
    name      text NOT NULL,
    version   text NOT NULL DEFAULT '',
    licence   text NOT NULL,
    credit    text NOT NULL,
    url       text NOT NULL
);
CREATE TABLE nutrition_entries (
    id              bigserial PRIMARY KEY,
    source          text NOT NULL REFERENCES nutrition_sources (code),
    source_code     text NOT NULL,
    source_version  text NOT NULL,
    name            text NOT NULL,
    name_en         text,
    food_group      text,
    state           text NOT NULL DEFAULT 'other' CHECK (state IN ('raw','dry','cooked','prepared','other')),
    energy_kcal     numeric(7,2), energy_kj numeric(8,2),
    protein_g numeric(6,2), carbohydrate_g numeric(6,2), sugars_g numeric(6,2), fat_g numeric(6,2),
    saturates_g numeric(6,2), fibre_g numeric(6,2), salt_g numeric(7,3),
    qualifiers      jsonb NOT NULL DEFAULT '{}',
    retired         boolean NOT NULL DEFAULT false,
    imported_at     timestamptz NOT NULL DEFAULT now(),
    UNIQUE (source, source_code)
);
CREATE INDEX nutrition_entries_name_trgm ON nutrition_entries USING gin (name gin_trgm_ops);
CREATE INDEX nutrition_entries_name_en_trgm ON nutrition_entries USING gin (name_en gin_trgm_ops);
CREATE TABLE ingredient_nutrition (
    ingredient_id  uuid NOT NULL REFERENCES ingredients (id),
    entry_id       bigint NOT NULL REFERENCES nutrition_entries (id),
    status         text NOT NULL CHECK (status IN ('suggested','approved','rejected')),
    method         text NOT NULL CHECK (method IN ('ai','admin')),
    confidence     real,
    warning        text,
    decided_by     text,
    decided_at     timestamptz,
    created_at     timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (ingredient_id, entry_id)
);
CREATE UNIQUE INDEX ingredient_nutrition_one_approved ON ingredient_nutrition (ingredient_id) WHERE status = 'approved';
```

## Errors

`ListNutritionSources` has none.

## Acceptance criteria

- [ ] Inserting a second approved link for one ingredient fails.
- [ ] An entry with only kJ gets kcal computed on import; one with only sodium gets salt.
- [ ] `ListNutritionSources` returns each imported source with its credit text.
