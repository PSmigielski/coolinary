# PLN-08-06 Cuisine mismatch: three choices

| | |
|---|---|
| Service | planner |
| Epic | [PLN-08 AI day plan](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-08-04, PLN-08-05, PLN-08-01 |
| Unblocks | PLN-08-07, PLN-15 (allowed purchases to the list) |

## Goal

When what's at home doesn't fit the cuisine picked for a meal, the user gets three choices instead of a bad dish:
the closest dish possible now, the 1–3 items to buy, or a different cuisine for today (features.md "Cuisine choice").

## Business rules

1. A meal "doesn't fit" when AI returns `cuisine_fit=false` for it, or (library path) no library dish of that cuisine
   reaches share on hand ≥ 0.5.
2. The three options, each a full dish that passed the checker:
   - `CLOSEST_NOW`: in the chosen cuisine, from stock only.
   - `BUY_FEW`: in the chosen cuisine, needing 1–3 ingredients not at home; those ids are the option's
     `to_buy` and become the meal's allowed purchases (PLN-07-04 rule 2).
   - `OTHER_CUISINE`: the best dish from stock in another cuisine, named in the option.
3. With AI: all three come from the same `ComposeDayPlan` response (no extra call), are checked like any AI dish and,
   if one fails after the repair, that option is filled from the library; if no library dish fits it, the option is
   left out (2 options shown).
4. Library path: `CLOSEST_NOW` = highest share in the cuisine; `BUY_FEW` = the dish in the cuisine with the fewest
   missing lines (1–3), tie → higher share; `OTHER_CUISINE` = top-ranked dish of any other cuisine.
5. The meal is stored with `options` and no dish; `ChooseMealOption(plan_id, meal, option)` sets the dish, and for
   `BUY_FEW` keeps `to_buy` as allowed purchases so the list proposal (PLN-15) adds them.
6. A cuisine of "no preference" never triggers this.
7. Choosing costs no AI call and no quota.

## Contract

```proto
rpc ChooseMealOption(ChooseMealOptionRequest) returns (ChooseMealOptionResponse);

message MealOption {
  MealOptionKind kind = 1;
  string dish_id = 2;
  int32 dish_version = 3;
  repeated string to_buy_ingredient_ids = 4;
  coolinary.common.v1.Cuisine cuisine = 5;
  Dish dish = 6;
}
enum MealOptionKind { MEAL_OPTION_KIND_UNSPECIFIED = 0; MEAL_OPTION_KIND_CLOSEST_NOW = 1;
                      MEAL_OPTION_KIND_BUY_FEW = 2; MEAL_OPTION_KIND_OTHER_CUISINE = 3; }
message ChooseMealOptionRequest { string household_id = 1; string plan_id = 2; Meal meal = 3; MealOptionKind kind = 4; }
message ChooseMealOptionResponse { Plan plan = 1; }
```

## Data

`plan_meals.options` (PLN-08-01) holds the options JSON.

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Meal has no open choice | FAILED_PRECONDITION | `no choice to make` |
| Option kind not offered | INVALID_ARGUMENT | `option not available` |

## Acceptance criteria

- [ ] For each launch cuisine, a day plan from the fixed test fridge `planner/testdata/fridges/mismatch.yaml` either
      has the meal in that cuisine or offers `BUY_FEW` naming the missing items (Verification 4).
- [ ] Choosing `BUY_FEW` puts its `to_buy` items into the meal's allowed purchases.
- [ ] No `BUY_FEW` option ever lists more than 3 items.
