# PLN-04-01 Compute and store nutrition per dish version

| | |
|---|---|
| Service | planner |
| Epic | [PLN-04 Dish nutrition and calories](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-02-01, PLN-01-02, KIT-06 (ComputeNutrition) |
| Unblocks | PLN-02-02 (publication rule 12), PLN-03-07, PLN-03-08, PLN-04-02, PLN-05-03, PLN-12 |

## Goal

Store, per dish version, nutrition per standard portion computed by kitchen from the ingredient lines, so browse,
filters and plans read it without calling kitchen.

## Business rules

1. Input to `KIT-06 ComputeNutrition`: every non-optional line's `ingredient_id`, amount and unit, plus
   `portions = base_servings`. Kitchen converts pieces and spoons to grams with typical weights and millilitres with
   density (features.md §5 steps 1–2) and returns totals, per-portion values and the ids it could not value.
2. Amounts are counted as weighed in the recipe (dry pasta, raw meat; §5 step 3).
3. Deep-frying oil: a line may carry `nutrition_amount`, the amount counted for nutrition (the oil the food absorbs,
   §5 step 3), written by the recipe author. A line whose ingredient is an oil with amount > 200 ml and no
   `nutrition_amount` is a blocking publication problem `FRYING_OIL_UNCOUNTED`. Decided (Q084): the recipe must
   state the absorbed oil; without it the dish shows "calories incomplete" and cannot be published. No default share.
4. Base-plus-toppings dishes: the stored values are base + the first option of each topping group ("default plate").
   Planner also stores `kcal_min` and `kcal_max` per portion over all topping combinations, shown as a range.
   Decided (Q092): the default plate is base + the first option of each topping group, plus the min–max range.
5. Stored values per portion: energy kcal (unrounded and rounded), protein, carbohydrate, sugars, fat, saturates, fibre,
   salt in grams (unrounded); rounding happens at display (PLN-04-03) except `kcal_rounded`, used by filters.
6. `kcal_rounded` = nearest 10 kcal (§5 step 5).
7. If kitchen returns any missing ingredient, `complete=false` and the missing ids are stored; `kcal_rounded` is NULL
   so no filter or plan uses a number that is too low (§5 step 6).
8. Computed on every write of a version (import, draft, admin edit, household AI dish, rewrite) inside the write
   path; kitchen unavailable → the write fails with UNAVAILABLE for library dishes, but household AI dishes are stored
   with `complete=false, pending=true` and queued for PLN-04-02's job (the user's plan must not fail because of it).

## Contract

```proto
message DishNutrition {
  bool complete = 1;
  repeated string missing_ingredient_ids = 2;
  repeated string missing_names = 3;          // in the requested language
  int32 kcal_per_portion = 4;                  // rounded to 10; 0 when incomplete
  int32 kj_per_portion = 5;                    // PLN-04-03
  double protein_g = 6;
  double carbohydrate_g = 7;
  double sugars_g = 8;
  double fat_g = 9;
  double saturates_g = 10;
  double fibre_g = 11;
  double salt_g = 12;
  int32 kcal_min = 13;                         // base + toppings range; = kcal_per_portion otherwise
  int32 kcal_max = 14;
  bool from_household_products = 15;           // PLN-04-04
}
```

## Data

`planner/migrations/000NN_dish_nutrition.sql`

```sql
ALTER TABLE dish_ingredients ADD COLUMN nutrition_amount numeric(10,3) CHECK (nutrition_amount > 0);

CREATE TABLE dish_nutrition (
    dish_id         uuid NOT NULL,
    version         int  NOT NULL,
    complete        boolean NOT NULL,
    pending         boolean NOT NULL DEFAULT false,
    missing_ids     text[] NOT NULL DEFAULT '{}',
    kcal            numeric(8,2),                 -- unrounded, per portion
    kcal_rounded    int,                          -- NULL when incomplete
    kcal_min        int,
    kcal_max        int,
    protein_g       numeric(7,2),
    carbohydrate_g  numeric(7,2),
    sugars_g        numeric(7,2),
    fat_g           numeric(7,2),
    saturates_g     numeric(7,2),
    fibre_g         numeric(7,2),
    salt_g          numeric(7,3),
    computed_at     timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (dish_id, version),
    FOREIGN KEY (dish_id, version) REFERENCES dish_versions (dish_id, version) ON DELETE CASCADE
);
CREATE INDEX dish_nutrition_kcal ON dish_nutrition (kcal_rounded);
CREATE INDEX dish_nutrition_protein ON dish_nutrition (protein_g);
```

## Errors

| Case | Behaviour |
|---|---|
| Kitchen down, library write | UNAVAILABLE, nothing written |
| Kitchen down, household AI dish | Stored with `pending=true`, recomputed later |

## Acceptance criteria

- [ ] Chicken curry with rice (worked example ingredients) gives one `dish_nutrition` row with `complete=true` and
      `kcal_rounded` divisible by 10.
- [ ] A dish with one unmatched ingredient has `complete=false`, `kcal_rounded IS NULL`, and the missing id stored
      (Verification 16: "calories incomplete", never a number).
- [ ] 1 l frying oil without `nutrition_amount` blocks publishing with `FRYING_OIL_UNCOUNTED`.
- [ ] No value in `dish_nutrition` ever comes from an ai response (only `ComputeNutrition`).
