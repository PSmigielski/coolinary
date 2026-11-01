# PLN-02-01 Dish schema (replaces skeleton `dishes`)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-02 Recipe library model](README.md) |
| Release | R1 |
| Size | M |
| Depends on | — |
| Unblocks | PLN-02-02 … PLN-02-06, PLN-03, PLN-04, PLN-05 |

## Goal

One schema for every dish planner knows: the shared library and each household's private AI dishes, with per-version
content and per-language text.

## Business rules

1. A **dish** has a stable `id` (UUID) and, for library dishes, a unique `slug` (lowercase, `a-z0-9-`, ≤ 60 chars).
2. `source` is one of: `import` (hand-written file), `ai_draft` (AI drafting batch), `ai_promoted` (copied from a
   household AI dish, PLN-11), `household_ai` (composed or rewritten for one household). `household_ai` dishes have a
   `household_id` and are visible only to that household; all others have none.
3. Content lives in **versions** `(dish_id, version)`, version starting at 1. `current_version` points at the version
   shown to users; `draft_version` at the version being edited, if any (PLN-02-04).
4. A version holds: meals (one or more of breakfast, dinner, supper — the existing `coolinary.planner.v1.Meal`),
   cuisine and optional sub-cuisine, required tools, total and active minutes, base servings, effort tags, ingredients,
   steps, texts per language.
5. Codes are shared with kitchen's profile so a dish and a profile can be compared without mapping. Cuisine:
   `polish, italian, asian, mediterranean, french, spanish, german, mexican, middle_eastern, american`; sub-cuisine
   (only under `asian`): `chinese, japanese, thai, indian, vietnamese, korean`. Tools: `hob, oven, microwave, kettle,
   air_fryer, blender, hand_blender, food_processor, multicooker, cooking_robot, toaster, grill, steamer, rice_cooker,
   mixer, wok, baking_tins` (features.md §1, plus `kettle` used by the "student room" preset). Effort tags:
   `no_cook, keeps_3_days, freezes_well, good_cold_next_day, base_plus_toppings`.
6. Each ingredient line: position, kitchen `ingredient_id`, amount (> 0), unit, `optional` flag, `role`
   (`main`, `base`, `topping`), `topping_group` (only for `topping`, e.g. "protein", "sauce"), and an optional
   per-language note ("finely chopped").
7. Each step: position, optional timer in seconds, optional tool code used in that step, text per language.
8. Texts per language: name, short intro, step texts. Every text row has `machine_translated` (true when written by
   AI translation, PLN-03-05).
9. Derived columns filled by PLN-02-03 (diet flags, allergens) and PLN-04-01 (nutrition) live on the version row.
10. The skeleton table `dishes` and its 5 seed rows are dropped; the seeds come back as import-file examples in
    PLN-03-01. `SuggestDishes` keeps working (PLN-06-03 rewrites its query).
11. Shared enums (README decision 16) live in `coolinary/common/v1/food.proto`, one definition each: `Cuisine`
    (KIT-10-02), `KitchenTool` (KIT-10-03), `Diet` (KIT-10-04), `Allergen` (KIT-02-01) and `EffortTag` (this task).
    Planner imports them and defines no copies; whichever task lands first creates the file.
12. Units: the common `Unit` enum is defined by KIT-02-02 (`UNIT_TABLESPOON = 4`, `UNIT_TEASPOON = 5`, `UNIT_CUP = 6`,
    `UNIT_PINCH = 7`, `UNIT_TO_TASTE = 8`); recipes use spoons and kitchen turns them into grams (features.md §5
    step 2). Planner reuses those values.

## Contract

`contracts/coolinary/common/v1/food.proto` — shared enum owned by this task (README decision 16; `Cuisine`,
`KitchenTool`, `Diet` and `Allergen` are defined in the same file by KIT-10-02/03/04 and KIT-02-01):

```proto
enum EffortTag { EFFORT_TAG_UNSPECIFIED = 0; EFFORT_TAG_NO_COOKING = 1; EFFORT_TAG_KEEPS_3_DAYS = 2;
  EFFORT_TAG_FREEZES_WELL = 3; EFFORT_TAG_GOOD_COLD_NEXT_DAY = 4; EFFORT_TAG_BASE_PLUS_TOPPINGS = 5; }
```

`contracts/coolinary/planner/v1/dish.proto`

```proto
message Dish {
  string id = 1;
  string slug = 2;                       // empty for household_ai
  DishSource source = 3;
  DishStatus status = 4;
  int32 version = 5;
  repeated Meal meals = 6;
  coolinary.common.v1.Cuisine cuisine = 7;
  coolinary.common.v1.Cuisine sub_cuisine = 8;     // CUISINE_UNSPECIFIED when none
  repeated coolinary.common.v1.KitchenTool tools = 9;
  int32 total_minutes = 10;
  int32 active_minutes = 11;
  int32 base_servings = 12;
  repeated coolinary.common.v1.EffortTag effort_tags = 13;
  repeated DishIngredient ingredients = 14;
  repeated DishStep steps = 15;
  DishText text = 16;                    // in the requested language (PLN-02-06)
  repeated coolinary.common.v1.Diet diets = 17;          // derived, PLN-02-03
  repeated coolinary.common.v1.Allergen allergens = 18;  // derived, PLN-02-03
  DishNutrition nutrition = 19;          // PLN-04
}

message DishIngredient {
  string ingredient_id = 1;
  string display_name = 2;               // from kitchen, in the requested language
  coolinary.common.v1.Quantity quantity = 3;
  bool optional = 4;
  IngredientRole role = 5;
  string topping_group = 6;
  string note = 7;
}

message DishStep {
  int32 position = 1;
  string text = 2;
  int32 timer_seconds = 3;               // 0 = no timer
  coolinary.common.v1.KitchenTool tool = 4;
}

message DishText {
  string language = 1;
  string name = 2;
  string intro = 3;
  bool machine_translated = 4;
}

enum DishSource { DISH_SOURCE_UNSPECIFIED = 0; DISH_SOURCE_IMPORT = 1; DISH_SOURCE_AI_DRAFT = 2;
                  DISH_SOURCE_AI_PROMOTED = 3; DISH_SOURCE_HOUSEHOLD_AI = 4; }
enum DishStatus { DISH_STATUS_UNSPECIFIED = 0; DISH_STATUS_DRAFT = 1; DISH_STATUS_REVIEW = 2;
                  DISH_STATUS_PUBLISHED = 3; DISH_STATUS_RETIRED = 4; DISH_STATUS_HOUSEHOLD = 5; }
enum IngredientRole { INGREDIENT_ROLE_UNSPECIFIED = 0; INGREDIENT_ROLE_MAIN = 1;
                      INGREDIENT_ROLE_BASE = 2; INGREDIENT_ROLE_TOPPING = 3; }
```

## Data

`planner/migrations/000NN_dish_library.sql`

```sql
DROP TABLE dishes;

CREATE TABLE dishes (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug             text UNIQUE CHECK (slug ~ '^[a-z0-9-]{1,60}$'),
    source           text NOT NULL CHECK (source IN ('import','ai_draft','ai_promoted','household_ai')),
    status           text NOT NULL CHECK (status IN ('draft','review','published','retired','household')),
    household_id     uuid,
    parent_dish_id   uuid REFERENCES dishes (id),       -- rewrites (PLN-09) and promotions (PLN-11)
    current_version  int,
    draft_version    int,
    created_at       timestamptz NOT NULL DEFAULT now(),
    published_at     timestamptz,
    retired_at       timestamptz,
    CHECK ((source = 'household_ai') = (household_id IS NOT NULL)),
    CHECK ((source = 'household_ai') = (slug IS NULL))
);
CREATE INDEX dishes_household ON dishes (household_id) WHERE household_id IS NOT NULL;
CREATE INDEX dishes_status ON dishes (status);

CREATE TABLE dish_versions (
    dish_id          uuid NOT NULL REFERENCES dishes (id) ON DELETE CASCADE,
    version          int  NOT NULL CHECK (version >= 1),
    meals            smallint[] NOT NULL,
    cuisine          text NOT NULL,
    sub_cuisine      text,
    tools            text[] NOT NULL DEFAULT '{}',
    total_minutes    int  NOT NULL CHECK (total_minutes BETWEEN 1 AND 600),
    active_minutes   int  NOT NULL CHECK (active_minutes BETWEEN 0 AND 600),
    base_servings    int  NOT NULL CHECK (base_servings BETWEEN 1 AND 12),
    effort_tags      text[] NOT NULL DEFAULT '{}',
    diets            text[] NOT NULL DEFAULT '{}',     -- PLN-02-03
    allergens        text[] NOT NULL DEFAULT '{}',     -- PLN-02-03
    content_hash     text NOT NULL,                    -- sha256 of normalised content, PLN-03-03
    created_at       timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (dish_id, version)
);
CREATE INDEX dish_versions_tools ON dish_versions USING gin (tools);
CREATE INDEX dish_versions_meals ON dish_versions USING gin (meals);

CREATE TABLE dish_ingredients (
    dish_id        uuid NOT NULL,
    version        int  NOT NULL,
    position       int  NOT NULL,
    ingredient_id  text NOT NULL,                      -- kitchen catalogue id
    amount         numeric(10,3) NOT NULL CHECK (amount > 0),
    unit           smallint NOT NULL,                  -- coolinary.common.v1.Unit
    optional       boolean NOT NULL DEFAULT false,
    role           text NOT NULL DEFAULT 'main' CHECK (role IN ('main','base','topping')),
    topping_group  text,
    PRIMARY KEY (dish_id, version, position),
    FOREIGN KEY (dish_id, version) REFERENCES dish_versions (dish_id, version) ON DELETE CASCADE,
    CHECK ((role = 'topping') = (topping_group IS NOT NULL))
);
CREATE INDEX dish_ingredients_ingredient ON dish_ingredients (ingredient_id);

CREATE TABLE dish_steps (
    dish_id        uuid NOT NULL,
    version        int  NOT NULL,
    position       int  NOT NULL,
    timer_seconds  int  NOT NULL DEFAULT 0 CHECK (timer_seconds BETWEEN 0 AND 86400),
    tool           text,
    PRIMARY KEY (dish_id, version, position),
    FOREIGN KEY (dish_id, version) REFERENCES dish_versions (dish_id, version) ON DELETE CASCADE
);

CREATE TABLE dish_texts (
    dish_id             uuid NOT NULL,
    version             int  NOT NULL,
    language            text NOT NULL,                 -- 'pl', 'de', 'en', ...
    name                text NOT NULL CHECK (length(name) BETWEEN 1 AND 120),
    intro               text NOT NULL DEFAULT '',
    step_texts          text[] NOT NULL,               -- aligned with dish_steps.position
    ingredient_notes    jsonb NOT NULL DEFAULT '{}',   -- position -> note
    machine_translated  boolean NOT NULL DEFAULT false,
    PRIMARY KEY (dish_id, version, language),
    FOREIGN KEY (dish_id, version) REFERENCES dish_versions (dish_id, version) ON DELETE CASCADE
);
CREATE INDEX dish_texts_name_trgm ON dish_texts USING gin (name gin_trgm_ops);
```

(`CREATE EXTENSION IF NOT EXISTS pg_trgm;` at the top, as kitchen does.)

## Errors

No RPC in this task.

## Acceptance criteria

- [ ] Migration up drops the skeleton table and creates the new ones; down restores the skeleton table empty.
- [ ] A `household_ai` row without `household_id`, or a library row with one, is rejected by the database.
- [ ] A topping line without `topping_group` is rejected.
- [ ] `make sqlc` generates queries for insert/select of a full version (dish, version, ingredients, steps, texts).
- [ ] `SuggestDishes` still answers (empty list) after the migration.
