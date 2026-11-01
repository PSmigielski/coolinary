# PLN-10-01 Cooking mode payload

| | |
|---|---|
| Service | planner |
| Epic | [PLN-10 Cooking mode and "cooked it"](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-02-06, PLN-04-04, PLN-08-01 |
| Unblocks | GW-07 (`GET /v1/dishes/{id}/cooking`), PLN-12 (pot split text) |

## Goal

Everything the cooking screen needs in one call, so it works offline once loaded: scaled amounts, steps, timers, tools,
and who leaves out what.

## Business rules

1. `GetCookingSteps(dish_id, version?, servings? | plan_meal_id)`: with a plan meal, dish, version, servings, plates and
   leave-outs come from the plan; otherwise from the request (servings default = household portions).
2. Returns: dish name and texts in the request language; scaled ingredient lines (display values, PLN-02-05), grouped
   main / base / toppings; steps in order with timer seconds and tool; total and active minutes; nutrition per portion
   with household product values where available (PLN-04-04); per-eater plates and leave-outs (member ids only — the app
   shows names).
3. R2: the pot split per person (`portions_by_member`, e.g. member A 1½, B 1) is filled by PLN-12; until then it is
   empty.
4. Timers are returned, not run: the app runs them on the phone.
5. Each scaled ingredient line covered by stock past its date carries `check_first = true` so the app shows "check it
   first" next to it (Q103); with a plan meal the ids come from the plan meal, otherwise from PLN-06-01.
6. The response is cacheable by the app for 24 h (`etag` = dish id, version, servings, plan meal id, language).

## Contract

```proto
rpc GetCookingSteps(GetCookingStepsRequest) returns (GetCookingStepsResponse);
message GetCookingStepsRequest { string household_id = 1; string dish_id = 2; int32 version = 3; double servings = 4;
                                 string plan_meal_id = 5; string language = 6; }
message MemberPortion { string member_id = 1; double portions = 2; }
message GetCookingStepsResponse {
  Dish dish = 1;                          // with steps
  double servings = 2;
  repeated ScaledIngredient scaled = 3;
  repeated Plate plates = 4;
  repeated LeaveOut leave_out = 5;
  repeated MemberPortion portions_by_member = 6;   // PLN-12
  string etag = 7;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status / behaviour |
|---|---|
| Neither `dish_id` nor `plan_meal_id` | `INVALID_ARGUMENT` |
| `servings` ≤ 0 or above 24 | `INVALID_ARGUMENT` |
| Unknown dish, version or plan meal, or a plan meal of another household | `NOT_FOUND` |
| Version retired | still served (history and plans may point at it) |
| kitchen down while reading household product values (PLN-04-04) | served with table values only; `nutrition.product_values_used = false` |

## Acceptance criteria

- [ ] For a plan meal with 6 servings, every amount is the 6-serving display value.
- [ ] A step with a 20-minute timer returns `timer_seconds = 1200`.
- [ ] Same inputs return the same `etag`.
- [ ] With milk one day past its date in the plan meal's allocation, the milk line has `check_first = true`; other
      lines have `false` (Q103).
