# PLN-08-01 Plan storage, accept and read

| | |
|---|---|
| Service | planner |
| Epic | [PLN-08 AI day plan](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-02-01 |
| Unblocks | PLN-08-07, PLN-09, PLN-10-02, PLN-12, PLN-14, PLN-15 (list from plan) |

## Goal

Store day plans so the household can review a proposal, accept it, see it later, cook from it and build the shopping
list from it.

## Business rules

1. A plan belongs to one household and one local date (`plan_date`, household time zone). Kind `day` in R1; `week`
   plans (PLN-14) group seven day plans.
2. Status: `proposed` → `accepted`, or `proposed`/`accepted` → `superseded`. Per household and date at most one
   `accepted` plan; accepting a new one supersedes the old in the same transaction.
3. Plan meals: meal, dish id and version (or, while a cuisine choice is open, up to 3 options and none chosen),
   servings (the pot), eater member ids, cuisine, `plates` and `leave_out` from the checker (PLN-07-03), allowed
   purchases, status `planned` / `cooked` / `skipped`, `source` (`ai`, `library`, `ai_repaired`, `library_fallback`,
   `manual` from SetPlanMeal, PLN-08-08), and the stock items past their date the meal uses (`check_first`, Q103).
4. A user can accept a plan with an open cuisine choice only after choosing (PLN-08-06).
5. Editing an accepted plan (swap a meal, PLN-09) changes that meal in place and keeps the plan accepted.
6. `GetPlans(from, to)` returns accepted plans plus the newest proposal per date, ≤ 31 days per call.
7. Proposed plans older than 7 days are deleted by a daily job; accepted plans are kept 1 year.
   Decided (Q118).
8. Plans store no names, weights or targets: eaters are member ids only.

## Contract

```proto
rpc GetPlans(GetPlansRequest) returns (GetPlansResponse);
rpc AcceptPlan(AcceptPlanRequest) returns (AcceptPlanResponse);
rpc DeletePlan(DeletePlanRequest) returns (DeletePlanResponse);

message Plan {
  string id = 1;
  string plan_date = 2;            // YYYY-MM-DD, household local
  PlanStatus status = 3;
  bool ai_used = 4;
  FallbackReason fallback_reason = 5;
  repeated PlanMeal meals = 6;
  string created_at = 7;
}
message PlanMeal {
  string id = 1;
  Meal meal = 2;
  string dish_id = 3;
  int32 dish_version = 4;
  double servings = 5;
  repeated string eater_member_ids = 6;
  coolinary.common.v1.Cuisine cuisine = 7;
  repeated MealOption options = 8;           // PLN-08-06; empty when a dish is set
  repeated Plate plates = 9;
  repeated LeaveOut leave_out = 10;
  repeated string allowed_purchase_ingredient_ids = 11;
  PlanMealStatus status = 12;
  PlanMealSource source = 13;
  Dish dish = 14;                            // list shape, for display
  repeated string check_first_stock_item_ids = 15;  // stock past its date this meal uses (Q103)
}
message Plate { string member_id = 1; map<string, StringList> toppings_by_group = 2; }
message LeaveOut { string member_id = 1; string ingredient_id = 2; bool required = 3; }
enum PlanMealSource { PLAN_MEAL_SOURCE_UNSPECIFIED = 0; PLAN_MEAL_SOURCE_AI = 1; PLAN_MEAL_SOURCE_LIBRARY = 2;
                      PLAN_MEAL_SOURCE_AI_REPAIRED = 3; PLAN_MEAL_SOURCE_LIBRARY_FALLBACK = 4;
                      PLAN_MEAL_SOURCE_MANUAL = 5; }   // MANUAL: SetPlanMeal, PLN-08-08
```

## Data

`planner/migrations/000NN_plans.sql`

```sql
CREATE TABLE plans (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    household_id    uuid NOT NULL,
    plan_date       date NOT NULL,
    kind            text NOT NULL DEFAULT 'day' CHECK (kind IN ('day','week')),
    status          text NOT NULL CHECK (status IN ('proposed','accepted','superseded')),
    ai_used         boolean NOT NULL,
    fallback_reason text,
    model           text,
    request_id      text,
    created_by      text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    accepted_at     timestamptz
);
CREATE UNIQUE INDEX plans_one_accepted ON plans (household_id, plan_date) WHERE status = 'accepted' AND kind = 'day';
CREATE INDEX plans_household_date ON plans (household_id, plan_date);
CREATE UNIQUE INDEX plans_request ON plans (household_id, request_id) WHERE request_id IS NOT NULL;

CREATE TABLE plan_meals (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id         uuid NOT NULL REFERENCES plans (id) ON DELETE CASCADE,
    meal            smallint NOT NULL,
    dish_id         uuid REFERENCES dishes (id),
    dish_version    int,
    servings        numeric(5,2) NOT NULL,
    eaters          uuid[] NOT NULL,
    cuisine         text,
    options         jsonb NOT NULL DEFAULT '[]',
    plates          jsonb NOT NULL DEFAULT '[]',
    leave_out       jsonb NOT NULL DEFAULT '[]',
    allowed_purchases text[] NOT NULL DEFAULT '{}',
    status          text NOT NULL DEFAULT 'planned' CHECK (status IN ('planned','cooked','skipped')),
    source          text NOT NULL CHECK (source IN ('ai','library','ai_repaired','library_fallback','manual')),
    check_first     uuid[] NOT NULL DEFAULT '{}',   -- stock items past their date (Q103)
    UNIQUE (plan_id, meal),
    CHECK ((dish_id IS NULL) = (options <> '[]'::jsonb))
);
```

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Accept with an open cuisine choice | FAILED_PRECONDITION | `choose an option first` |
| Plan of another household | NOT_FOUND | `plan not found` |
| Range > 31 days | INVALID_ARGUMENT | limits |

## Acceptance criteria

- [ ] Accepting a second plan for the same date supersedes the first; only one is `accepted`.
- [ ] `GetPlans` for a week returns accepted plans and the latest proposal for dates without an accepted plan.
- [ ] The database rejects a meal with both a dish and options.
