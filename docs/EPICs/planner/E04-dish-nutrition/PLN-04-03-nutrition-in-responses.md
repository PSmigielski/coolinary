# PLN-04-03 Nutrition in dish responses (kcal/kJ, rounding, incomplete)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-04 Dish nutrition and calories](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-04-01, PLN-02-06 |
| Unblocks | GW-07 (dish detail), PLN-05, PLN-10-01, PLN-12 |

## Goal

Every dish response carries nutrition rounded the way EU pack labels are, in kcal and kJ, per portion and for the whole
pot, with "calories incomplete" instead of a number when values are missing.

## Business rules

1. Per portion values come from `dish_nutrition`; they do not change when servings change (PLN-02-05 rule 6). The
   response also carries `total_kcal = kcal × servings`, rounded to 10.
2. `kj_per_portion` = unrounded kcal × 4.184, rounded to the nearest 10. Both units are always sent; the app shows kcal
   by default and kJ when the user picks it (features.md "Calories and goals"; the setting lives in the app).
3. Gram values are rounded as on EU labels (European Commission guidance on tolerances and rounding, 2012): fat,
   saturates, carbohydrate, sugars, protein, fibre ≥ 10 g → nearest 1 g; 0.5–10 g → nearest 0.1 g; < 0.5 g → "0"
   (sent as 0 and listed in `zero_nutrients`); salt ≥ 1 g → nearest 0.1 g; 0.0125–1 g → nearest 0.01 g; below → 0.
4. Incomplete: `complete=false`, `kcal_per_portion=0`, `missing_names` in the requested language; the app shows
   "calories incomplete" and the names (§5 step 6).
5. Base plus toppings: `kcal_min`/`kcal_max` sent; the app shows "about 450–610 kcal".
6. No text field ever says "healthy" or "light" (features.md "Calories and goals").
7. Calories are never hidden by the server; hiding per person (children by default) is decided by the app from the
   profile flag (KIT-11). Decided (Q096): the app hides them; planner never strips values.

## Contract

`DishNutrition` (PLN-04-01) plus:

```proto
message DishNutrition {
  // ...
  int32 total_kcal = 16;                 // kcal x servings, rounded to 10
  repeated string zero_nutrients = 17;   // nutrients sent as 0 because they are below the rounding threshold
}
```

## Data

No schema change.

## Errors

None of its own; formatting applied inside dish responses (PLN-04-01), whose errors apply. Incomplete calories are a value (`calories_incomplete` with the missing names), never an error.

## Acceptance criteria

- [ ] 523.4 kcal unrounded → 520 kcal, 2,190 kJ (523.4 × 4.184 = 2,189.9).
- [ ] Fat 12.46 g → 12 g; sugars 4.36 g → 4.4 g; salt 0.837 g → 0.84 g; salt 1.26 g → 1.3 g.
- [ ] An incomplete dish returns `kcal_per_portion=0` and its missing names in Polish when asked in `pl`.
- [ ] Scaling servings from 4 to 6 changes `total_kcal` and not `kcal_per_portion`.
