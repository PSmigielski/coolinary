# KIT-10-02 Cuisines per meal

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | KIT-10-01 |
| Unblocks | KIT-10-08, PLN-05 (library filter), PLN-08 (AI composition in the chosen cuisine) |

## Goal

The household picks favourite cuisines, separately for breakfast, dinner and supper ("Polish breakfast, Asian
dinner"), to filter the library and steer AI composition (features.md §1 "Cuisines", §2 "Cuisine choice").

## Business rules

1. Cuisine codes (shared enum `coolinary.common.v1.Cuisine`, also used for dish tags by PLN-02-01): Polish, Italian,
   Asian, Chinese, Japanese, Thai, Indian, Vietnamese, Korean, Mediterranean, French, Spanish, German, Mexican,
   Middle Eastern, American.
2. Chinese, Japanese, Thai, Indian, Vietnamese and Korean are sub-cuisines of Asian. Choosing a sub-cuisine narrows
   Asian: `[ASIAN, THAI]` means Thai only; `[ASIAN]` alone means any Asian. A sub-cuisine sent without `ASIAN` is
   stored with `ASIAN` added.
3. Each meal holds 0–12 codes, no duplicates. An empty list means "no preference" for that meal.
4. `no_preference = true` on the section clears every meal (the explicit "no preference" choice in the app).
5. Default: the home cuisine of the country (KIT-10-01 rule 2) for all three meals; no home cuisine → empty.
6. Cuisines are a preference, not a hard limit: planner may show other cuisines when nothing fits (features.md §2,
   "three choices"); tools, diet and allergens stay hard limits.
7. Order inside a meal is kept (the first is the favourite) and returned as stored.

## Contract

`contracts/coolinary/kitchen/v1/profile.proto` (shared enum in `coolinary/common/v1`)

```proto
// coolinary/common/v1/food.proto
enum Cuisine {
  CUISINE_UNSPECIFIED = 0;
  CUISINE_POLISH = 1; CUISINE_ITALIAN = 2; CUISINE_ASIAN = 3;
  CUISINE_CHINESE = 4; CUISINE_JAPANESE = 5; CUISINE_THAI = 6; CUISINE_INDIAN = 7;
  CUISINE_VIETNAMESE = 8; CUISINE_KOREAN = 9;
  CUISINE_MEDITERRANEAN = 10; CUISINE_FRENCH = 11; CUISINE_SPANISH = 12; CUISINE_GERMAN = 13;
  CUISINE_MEXICAN = 14; CUISINE_MIDDLE_EASTERN = 15; CUISINE_AMERICAN = 16;
}

// kitchen/v1/profile.proto
message CuisineSection {
  repeated coolinary.common.v1.Cuisine breakfast = 1;
  repeated coolinary.common.v1.Cuisine dinner = 2;
  repeated coolinary.common.v1.Cuisine supper = 3;
  bool no_preference = 4;
}
```

## Data

`kitchen/migrations/000NN_profile_cuisines.sql`

```sql
-- +goose Up
CREATE TABLE profile_cuisines (
    household_id uuid NOT NULL REFERENCES kitchen_profiles (household_id) ON DELETE CASCADE,
    meal         smallint NOT NULL CHECK (meal IN (1, 2, 3)),   -- breakfast, dinner, supper
    cuisine      smallint NOT NULL CHECK (cuisine BETWEEN 1 AND 16),
    position     smallint NOT NULL,
    PRIMARY KEY (household_id, meal, cuisine)
);

-- +goose Down
DROP TABLE profile_cuisines;
```

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Unknown code, duplicate, more than 12 per meal | `INVALID_ARGUMENT` | `invalid_value` |

## Acceptance criteria

- [ ] Saving breakfast [Polish], dinner [Thai], supper [] returns dinner [Asian, Thai] and supper "no preference".
- [ ] A German household's defaults are German for all three meals.
- [ ] `no_preference = true` empties all three meals.
- [ ] features.md Verification 4 depends on this: a day plan for each launch cuisine reads the cuisine from here.
