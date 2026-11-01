# KIT-10-04 Household diet, allergens and dislikes

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-10-01, KIT-02 (ingredient ids for dislikes) |
| Unblocks | KIT-10-08, KIT-11-04, PLN-07 (allergen and diet check) |

## Goal

The household states its diet, the allergens it must avoid and the ingredients it won't eat; these are hard limits
for every dish shown (features.md §1 "Diet and allergens", rule 2 "AI proposes, rules check").

## Business rules

1. **Diet flags** (shared enum `coolinary.common.v1.Diet`): vegetarian, vegan, gluten-free, lactose-free,
   nut-free. Vegan implies vegetarian: saving `vegan` stores both. Their meaning against ingredients is KIT-02-01
   rule 5 (gluten-free = no gluten allergen; nut-free = no peanuts and no tree nuts; lactose-free = no lactose).
2. **Allergens to avoid**: any of the 14 EU allergens (shared enum `Allergen`, codes as in KIT-02-01). DECISION
   NEEDED: features.md lists only the diet flags; offering all 14 allergens costs nothing and is safer (e.g. sesame,
   celery, mustard); alternative: only the five flags.
3. **Dislikes**: 0–50 ingredient ids. They are hard limits too (features.md §1 table). Unknown ids are rejected;
   retired ingredients stay valid; merged ones are re-pointed by the KIT-02-09 registry (KIT-13-09).
4. All three apply to every meal of the household. In R2, a person's own values add to these for the meals they eat
   (KIT-11-04).
5. The section is replaced as a whole on save (KIT-10-01 rule 4).
6. The allergen filter is free for everyone (features.md §6: safety).
7. Values are not health data under the docs' definition (only age, height, weight and goals are), so they are
   stored normally, but they are never logged.

## Contract

```proto
// coolinary/common/v1/food.proto
enum Diet {
  DIET_UNSPECIFIED = 0;
  DIET_VEGETARIAN = 1; DIET_VEGAN = 2; DIET_GLUTEN_FREE = 3; DIET_LACTOSE_FREE = 4; DIET_NUT_FREE = 5;
}
// Allergen: the 14 EU allergens, defined once in this file by KIT-02-01.

// kitchen/v1/profile.proto
message DietSection {
  repeated coolinary.common.v1.Diet diets = 1;
  repeated coolinary.common.v1.Allergen avoid_allergens = 2;
  repeated string dislike_ingredient_ids = 3;
}
```

## Data

`kitchen/migrations/000NN_profile_restrictions.sql`

```sql
-- +goose Up
-- member_id NULL = the whole household; KIT-11-04 adds rows per member.
CREATE TABLE profile_restrictions (
    household_id uuid NOT NULL REFERENCES kitchen_profiles (household_id) ON DELETE CASCADE,
    member_id    uuid,
    kind         text NOT NULL CHECK (kind IN ('diet', 'allergen')),
    code         text NOT NULL,
    UNIQUE NULLS NOT DISTINCT (household_id, member_id, kind, code)
);

CREATE TABLE profile_dislikes (
    household_id  uuid NOT NULL REFERENCES kitchen_profiles (household_id) ON DELETE CASCADE,
    member_id     uuid,
    ingredient_id uuid NOT NULL REFERENCES ingredients (id),
    UNIQUE NULLS NOT DISTINCT (household_id, member_id, ingredient_id)
);

-- +goose Down
DROP TABLE profile_dislikes;
DROP TABLE profile_restrictions;
```

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Unknown code, more than 50 dislikes | `INVALID_ARGUMENT` | `invalid_value` |
| Unknown ingredient id | `INVALID_ARGUMENT` | `unknown_ingredient` |

## Acceptance criteria

- [ ] Saving `[vegan]` returns `[vegetarian, vegan]`.
- [ ] Saving allergens `[gluten, tree_nuts, peanuts]` and 3 dislikes round-trips unchanged.
- [ ] features.md Verification 3 ("no gluten, no nuts": 50 dishes, zero violations) reads its limits from here.
- [ ] A 51st dislike is rejected.
