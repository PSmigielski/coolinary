# KIT-02-01 Ingredient core: categories, allergens, diet flags

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | — |
| Unblocks | every other KIT-02 task, KIT-03-02, KIT-06-01 |

## Goal

Store each ingredient once with the facts that hard limits depend on: its allergens and whether it fits vegetarian,
vegan and lactose-free diets (features.md §1 "Diet and allergens: Hard limit").

## Business rules

1. An ingredient has a UUID `id` (what other services store) and a unique `slug` (`chicken-thigh`, lowercase ASCII,
   digits and hyphens, 2–60 characters) used by seed files and admin.
2. `category` is one of: `vegetables`, `fruit`, `herbs_spices`, `grains_pasta_rice`, `bread_bakery`, `legumes`,
   `nuts_seeds`, `meat`, `poultry`, `fish_seafood`, `eggs`, `dairy`, `cheese`, `plant_alternatives`, `oils_fats`,
   `sauces_condiments`, `baking_sweeteners`, `drinks`, `canned_jarred`, `frozen_prepared`, `snacks_sweets`, `other`.
   The shopping list groups lines by this category as the shop section (PLN-15).
3. `base_unit` is how stock and recipes keep amounts: `GRAM` for solids, `MILLILITRE` for liquids, `PIECE` for
   countable items (eggs). It uses `coolinary.common.v1.Unit`.
4. `allergens` is a set of the 14 EU allergens (Regulation (EU) 1169/2011, Annex II), coded: `gluten`
   (cereals containing gluten), `crustaceans`, `eggs`, `fish`, `peanuts`, `soy`, `milk`, `tree_nuts`, `celery`,
   `mustard`, `sesame`, `sulphites`, `lupin`, `molluscs`. It lists what the ingredient contains by nature, not
   "may contain" traces (those come with products, KIT-04-07).
5. Diet flags: `is_vegetarian`, `is_vegan`, `contains_lactose`. Derived, never stored: `gluten_free` = no `gluten`;
   `nut_free` = neither `peanuts` nor `tree_nuts`; `lactose_free` = not `contains_lactose`.
6. Consistency checks enforced by the database: vegan implies vegetarian; `contains_lactose` implies `milk` in
   allergens; `fish`, `crustaceans` or `molluscs` imply not vegetarian; categories `meat`, `poultry`,
   `fish_seafood` imply not vegetarian.
7. `status` is `active`, `retired` (hidden from search and new matches, still resolvable) or `merged`
   (`merged_into` points to an active ingredient; KIT-02-09).
8. `is_staple_candidate` marks items the profile suggests as "always at home" (salt, oil, flour, sugar, pepper,
   water). It changes no other rule.
9. `nutrition_version` (integer, starts at 1) and `popularity` (integer, starts at 0) are columns used by KIT-06-10
   and KIT-02-05.

## Contract

`contracts/coolinary/common/v1/food.proto` — the one definition of `Allergen` for every service (README decision 16;
`Cuisine`, `KitchenTool`, `Diet` and `EffortTag` live in the same file, see KIT-10-02/03/04 and PLN-02-01):

```proto
enum Allergen { ALLERGEN_UNSPECIFIED = 0; ALLERGEN_GLUTEN = 1; ALLERGEN_CRUSTACEANS = 2; ALLERGEN_EGGS = 3;
  ALLERGEN_FISH = 4; ALLERGEN_PEANUTS = 5; ALLERGEN_SOY = 6; ALLERGEN_MILK = 7; ALLERGEN_TREE_NUTS = 8;
  ALLERGEN_CELERY = 9; ALLERGEN_MUSTARD = 10; ALLERGEN_SESAME = 11; ALLERGEN_SULPHITES = 12; ALLERGEN_LUPIN = 13;
  ALLERGEN_MOLLUSCS = 14; }
```

`contracts/coolinary/kitchen/v1/catalogue.proto` (message types used by KIT-02-06 and later RPCs; `Allergen` is
imported as `coolinary.common.v1.Allergen`)

```proto
enum Category { CATEGORY_UNSPECIFIED = 0; CATEGORY_VEGETABLES = 1; CATEGORY_FRUIT = 2; /* … one per rule 2 */ }
message DietFlags { bool vegetarian = 1; bool vegan = 2; bool gluten_free = 3; bool lactose_free = 4; bool nut_free = 5; }
enum IngredientStatus { INGREDIENT_STATUS_UNSPECIFIED = 0; INGREDIENT_STATUS_ACTIVE = 1;
  INGREDIENT_STATUS_RETIRED = 2; INGREDIENT_STATUS_MERGED = 3; }
```

## Data

`kitchen/migrations/000NN_ingredients.sql`

```sql
CREATE TABLE ingredients (
    id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug                 text NOT NULL UNIQUE CHECK (slug ~ '^[a-z0-9][a-z0-9-]{1,59}$'),
    category             text NOT NULL,
    base_unit            smallint NOT NULL CHECK (base_unit IN (1, 2, 3)),   -- common.v1.Unit
    allergens            text[] NOT NULL DEFAULT '{}',
    is_vegetarian        boolean NOT NULL,
    is_vegan             boolean NOT NULL,
    contains_lactose     boolean NOT NULL DEFAULT false,
    is_staple_candidate  boolean NOT NULL DEFAULT false,
    status               text NOT NULL DEFAULT 'active' CHECK (status IN ('active','retired','merged')),
    merged_into          uuid REFERENCES ingredients (id),
    nutrition_version    integer NOT NULL DEFAULT 1,
    popularity           integer NOT NULL DEFAULT 0,
    created_at           timestamptz NOT NULL DEFAULT now(),
    updated_at           timestamptz NOT NULL DEFAULT now(),
    CHECK (allergens <@ ARRAY['gluten','crustaceans','eggs','fish','peanuts','soy','milk','tree_nuts','celery',
                              'mustard','sesame','sulphites','lupin','molluscs']),
    CHECK (NOT is_vegan OR is_vegetarian),
    CHECK (NOT contains_lactose OR 'milk' = ANY (allergens)),
    CHECK (NOT is_vegetarian OR NOT (allergens && ARRAY['fish','crustaceans','molluscs'])),
    CHECK (NOT is_vegetarian OR category NOT IN ('meat','poultry','fish_seafood')),
    CHECK ((status = 'merged') = (merged_into IS NOT NULL))
);
CREATE INDEX ingredients_category ON ingredients (category) WHERE status = 'active';
```

## Errors

Schema task; RPC errors come with KIT-02-06 and KIT-02-10.

## Acceptance criteria

- [ ] Inserting a vegan ingredient that isn't vegetarian, a lactose ingredient without `milk`, or a vegetarian
      `poultry` ingredient fails.
- [ ] A Go function `DietFlagsOf(ingredient)` returns all five flags; table tests cover each rule in 5–6.
- [ ] The proto enums map 1:1 to the stored codes (test fails if a code is added on one side only).
