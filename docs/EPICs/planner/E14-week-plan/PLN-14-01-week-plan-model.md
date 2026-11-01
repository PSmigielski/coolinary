# PLN-14-01 Week plan model, read and edit

| | |
|---|---|
| Service | planner |
| Epic | [PLN-14 Week plan](README.md) |
| Release | R3 |
| Size | M |
| Depends on | PLN-08, PLN-12-06 |
| Unblocks | PLN-14-02, PLN-14-06, PLN-14-07, PLN-15-02, GW-07 |

## Goal

Store a week as 7 linked day plans with cooking sessions, so the app can show it, the shopping list can be built
from it, and single meals can be changed without breaking leftovers.

## Business rules

1. A week plan has a start date (any day, chosen by the user; default next Monday) and covers 7 consecutive dates.
   One active week plan per household per start date; weeks may not overlap.
2. Each date is a PLN-08 `plans` row (kind `week`); each main-meal slot points to a cooking session, either as the cooking meal or as
   its leftover meal.
3. A cooking session: date, meal, dish, amount factor (1 or 2), pot portions (sum of both meals' pots), leftover
   storage (`fridge` or `freezer`), leftover date and meal.
4. Editing one meal:
   - replacing a cooking meal's dish replaces the leftover meal's dish too (they are the same food) and both are
     re-allocated;
   - replacing only the leftover meal turns the session into a single amount (factor 1) and the leftover slot gets
     the new dish;
   - marking someone away re-allocates that meal; a leftover meal's pot is fixed by the session, so portions are
     re-split within it (PLN-14-03 rule 5).
5. Each edit bumps the week plan version; edits on a stale version return `ABORTED` with the current version.
6. Deleting a week deletes its day plans for dates not yet cooked; cooked meals stay in history.
7. Reads return days, meals, portions (visibility-filtered, PLN-12-07), sessions and the effort summary
   (PLN-14-04).
8. Plus only. A household that drops to Free keeps read access to existing weeks; new weeks and edits need Plus.

## Contract

`contracts/coolinary/planner/v1/week.proto`

```proto
rpc GetWeekPlan(GetWeekPlanRequest) returns (WeekPlan);
rpc ReplaceWeekMeal(ReplaceWeekMealRequest) returns (WeekPlan);
rpc DeleteWeekPlan(DeleteWeekPlanRequest) returns (DeleteWeekPlanResponse);

message GetWeekPlanRequest { string household_id = 1; string viewer_member_id = 2; string start_date = 3; }
message WeekPlan {
  string id = 1; string start_date = 2; int64 version = 3;
  repeated FamilyDayPlan days = 4;          // PLN-13-07 shape
  repeated CookingSession sessions = 5;
  EffortSummary effort = 6;                 // PLN-14-04
  repeated PersonWeek people = 7;           // PLN-14-05, visibility-filtered
  bool ai_used = 8; string fallback_reason = 9;
}
message CookingSession {
  string id = 1; string date = 2; Meal meal = 3; string dish_id = 4;
  int32 amount_factor = 5; int32 pot_portions = 6;
  string leftover_date = 7; Meal leftover_meal = 8; LeftoverStorage storage = 9;
}
enum LeftoverStorage { LEFTOVER_STORAGE_UNSPECIFIED = 0; LEFTOVER_STORAGE_FRIDGE = 1; LEFTOVER_STORAGE_FREEZER = 2; }
message ReplaceWeekMealRequest { string household_id = 1; string week_plan_id = 2; int64 base_version = 3; string date = 4; Meal meal = 5; string dish_id = 6; string viewer_member_id = 7; }
message DeleteWeekPlanRequest { string household_id = 1; string week_plan_id = 2; }
message DeleteWeekPlanResponse {}
```

## Data

`planner/migrations/000NN_week_plans.sql`

```sql
CREATE TABLE week_plans (
    id              uuid PRIMARY KEY,
    household_id    uuid NOT NULL,
    start_date      date NOT NULL,
    version         bigint NOT NULL DEFAULT 1,
    ai_used         boolean NOT NULL,
    fallback_reason text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT now(),
    EXCLUDE USING gist (household_id WITH =, daterange(start_date, start_date + 7) WITH &&)
);
CREATE TABLE cooking_sessions (
    id             uuid PRIMARY KEY,
    week_plan_id   uuid NOT NULL REFERENCES week_plans ON DELETE CASCADE,
    date           date NOT NULL,
    meal           smallint NOT NULL,
    dish_id        uuid NOT NULL,
    amount_factor  smallint NOT NULL CHECK (amount_factor IN (1, 2)),
    pot_portions   integer NOT NULL,
    leftover_date  date,
    leftover_meal  smallint,
    storage        smallint NOT NULL DEFAULT 0,
    active_minutes integer NOT NULL
);
ALTER TABLE plans      ADD COLUMN week_plan_id uuid REFERENCES week_plans ON DELETE CASCADE;   -- PLN-08; day rows of a week use kind 'week'
ALTER TABLE plan_meals ADD COLUMN cooking_session_id uuid REFERENCES cooking_sessions ON DELETE SET NULL,
                       ADD COLUMN is_leftover boolean NOT NULL DEFAULT false;
```

(`btree_gist` extension needed for the exclusion constraint.)

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Overlapping week | `ALREADY_EXISTS` |
| Stale `base_version` | `ABORTED` |
| Free household editing | `FAILED_PRECONDITION` `plus_required` |
| Editing a cooked meal | `FAILED_PRECONDITION` `meal_cooked` |

## Acceptance criteria

- [ ] Replacing Monday's curry (cooked double) also changes Tuesday's leftover dinner.
- [ ] Replacing only Tuesday's leftover makes Monday's session factor 1 and its pot shrinks to Monday's pot.
- [ ] Two clients editing the same week: the second gets `ABORTED` with the new version.

## Notes

- PLN-08's `plans` table already has `kind` (`day` / `week`); the 7 day rows of a week use `kind = 'week'` and point to
  `week_plans`, so PLN-08's one-accepted-day-plan index does not block them. Check the index covers both kinds when a
  week replaces day plans.
