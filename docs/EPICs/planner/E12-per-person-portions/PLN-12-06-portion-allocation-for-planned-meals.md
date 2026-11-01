# PLN-12-06 Portion allocation for planned meals

| | |
|---|---|
| Service | planner |
| Epic | [PLN-12 Per-person portions](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | PLN-12-01, PLN-12-02, PLN-12-03, PLN-12-04, PLN-12-05, PLN-08 |
| Unblocks | PLN-12-07, PLN-12-08, PLN-13-06, PLN-13-07, PLN-15-10, GW-07 |

## Goal

Store each planned meal's pot and per-person portions, recompute them when anything they depend on changes, and
return them with "divide the pot" data for the dish card and cooking mode (features.md step 7).

## Business rules

1. Applies to every planned meal (PLN-08's plan meal) that has a dish, for households whose snapshot has at least
   one target. Households without targets keep R1 behaviour: pot = household portions setting.
2. Allocation runs per day, all three meals together (carry-over needs the order), using PLN-12-02 … PLN-12-05.
3. Stored per plan meal: the pot is written to PLN-08's `plan_meals.servings` and the eaters to `plan_meals.eaters`; in
   addition: pot portions, spare portions, eaters count, `calories_incomplete`, per eater: portion,
   share, portion kcal, gap, status, carry-in; extras in `meal_extras` (PLN-12-04).
4. An `inputs_hash` (profile snapshot version, attendance version, dish id and its kcal, stock version) is stored;
   recompute is skipped when it is unchanged.
5. Recompute triggers: `profile.changed`, attendance change (PLN-12-01), dish change on the meal (swap, PLN-09),
   the dish's kcal change (`catalogue.nutrition_changed` handled by PLN-04), `household.member_added` /
   `member_removed` (PLN-21-02). Only meals dated today or later and not yet cooked are recomputed; cooked and past
   meals are frozen.
6. If the calculator returns `UNSUITABLE`, the meal is marked `needs_new_dish`. Inside a composer run the composer
   picks again; for a dish the user chose by hand, the allocation keeps ½ portion for that eater and returns a
   warning `portion_too_big_for` with the member id.
7. "Divide the pot" data: pot portions and, in member order from the snapshot, `{member_id, portions}`; unnamed
   eaters appear as "guest". The app renders fractions as ½, ¼, ¾ ("6 portions: Adam 1½, Ewa 1½, Ola 1, Kuba 2").
8. Every response passes through the visibility filter (PLN-12-07).
9. Plus only for more than one target (features.md §6); Free households run with their one personal goal
   (PLN-12-02 rule 6).
10. After a recompute that changes a pot, PLN-15-10 is notified so list lines follow.

## Contract

`contracts/coolinary/planner/v1/portions.proto`

```proto
rpc GetMealPortions(GetMealPortionsRequest) returns (MealPortions);

message GetMealPortionsRequest { string household_id = 1; string plan_meal_id = 2; string viewer_member_id = 3; }
message MealPortions {
  string plan_meal_id = 1;
  int32 pot_portions = 2;
  double spare_portions = 3;
  bool calories_incomplete = 4;
  repeated EaterPortion eaters = 5;  // member order
  repeated string warnings = 6;      // e.g. "portion_too_big_for:<member_id>"
}
message EaterPortion {
  string member_id = 1;              // empty for unnamed eaters
  bool guest = 2;
  double portions = 3;               // 0.5 … 2.0 in steps of 0.25
  optional int32 portion_kcal = 4;   // omitted by the visibility filter
  repeated MealExtra extras = 5;
  EaterStatus status = 6;            // OK, SHORT, STILL_SHORT, OVER_MINIMUM
}
message MealExtra { string ingredient_id = 1; string label_key = 2; coolinary.common.v1.Quantity serving = 3; optional int32 kcal = 4; bool from_stock = 5; }
```

PLN-08's `GetDayPlan` response embeds `MealPortions` per meal (field added there).

## Data

`planner/migrations/000NN_meal_portions.sql`

```sql
CREATE TABLE plan_meal_pots (
    plan_meal_id        uuid PRIMARY KEY REFERENCES plan_meals ON DELETE CASCADE, -- PLN-08
    pot_portions        integer NOT NULL CHECK (pot_portions >= 0),
    spare_portions      numeric(4,2) NOT NULL DEFAULT 0,
    calories_incomplete boolean NOT NULL DEFAULT false,
    needs_new_dish      boolean NOT NULL DEFAULT false,
    inputs_hash         text NOT NULL,
    computed_at         timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE meal_portions (
    plan_meal_id  uuid NOT NULL REFERENCES plan_meal_pots ON DELETE CASCADE,
    member_key    text NOT NULL,                     -- member id or "unnamed:<n>"
    position      smallint NOT NULL,
    portions      numeric(3,2) NOT NULL CHECK (portions BETWEEN 0.5 AND 2.0),
    share_kcal    integer,
    portion_kcal  integer,
    gap_kcal      integer,
    carry_in_kcal integer NOT NULL DEFAULT 0,
    status        smallint NOT NULL,
    PRIMARY KEY (plan_meal_id, member_key)
);
```

## Errors

| Case | gRPC status |
|---|---|
| Plan meal not in this household | `NOT_FOUND` |
| Viewer not a member | `PERMISSION_DENIED` |
| kitchen down during recompute | stored allocation kept, response flagged `stale`; recompute retried on next trigger |

## Acceptance criteria

- [ ] Worked example stored as pot 6 with Adam 1½, Ewa 1½, Ola 1, Kuba 2 and Kuba's extras (Verification 17).
- [ ] Marking Ola away recomputes only that dinner and the supper after it; the pot becomes 5.
- [ ] Re-running with unchanged inputs does no writes (hash match).
- [ ] A meal marked cooked is never recomputed, even after `profile.changed`.
- [ ] A household without any target gets pot = household portions and no per-eater rows.
