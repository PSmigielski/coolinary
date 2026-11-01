# PLN-12-01 Meal attendance (marking people away)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-12 Per-person portions](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | PLN-01, PLN-08 |
| Unblocks | PLN-12-06, PLN-15-10, GW-07 |

## Goal

A person can be marked away for a meal (school lunch, eating out). They are left out of that meal, the pot shrinks,
and the shopping list follows (features.md §2 "Calories and goals", Verification 17).

## Business rules

1. Attendance is per household, member, date and meal (`breakfast`, `dinner`, `supper`). Default: everyone in the
   household eats every meal.
2. A one-off mark sets one member away for one date and meal. A whole day is three marks.
3. A weekly rule sets a member away on chosen weekdays for one meal, from a start date to an optional end date
   (e.g. "Ola, dinner, Mon–Fri" for school lunch). A one-off "present" mark overrides a rule for that date.
4. Any account member of the household may set or clear attendance for any member, including profiles without
   accounts (children). Free households can mark only their own account (they have no other members).
5. Marks for dates before today (household time zone from the profile snapshot) cannot be changed: `FAILED_PRECONDITION`.
6. Marks for meals already marked "cooked" (PLN-10) cannot be changed.
7. Setting or clearing a mark bumps the household's attendance version and triggers portion recompute for every
   affected planned meal (PLN-12-06), which in turn updates list lines that came from those meals (PLN-15-10).
8. If every eater of a meal is away, the meal keeps its dish but its pot becomes 0 and it adds nothing to the list.
9. A member removed from the household loses all marks and rules (PLN-21-02).
10. At most 20 weekly rules per household; more returns `RESOURCE_EXHAUSTED`.

## Contract

`contracts/coolinary/planner/v1/portions.proto`

```proto
rpc SetMealAttendance(SetMealAttendanceRequest) returns (SetMealAttendanceResponse);
rpc ListMealAttendance(ListMealAttendanceRequest) returns (ListMealAttendanceResponse);
rpc SetAttendanceRule(SetAttendanceRuleRequest) returns (AttendanceRule);
rpc DeleteAttendanceRule(DeleteAttendanceRuleRequest) returns (DeleteAttendanceRuleResponse);

message SetMealAttendanceRequest {
  string household_id = 1;
  string member_id = 2;
  string date = 3;            // YYYY-MM-DD
  Meal meal = 4;
  AttendanceStatus status = 5; // AWAY, PRESENT, or DEFAULT (clears the one-off mark)
}
enum AttendanceStatus { ATTENDANCE_STATUS_UNSPECIFIED = 0; ATTENDANCE_STATUS_DEFAULT = 1; ATTENDANCE_STATUS_AWAY = 2; ATTENDANCE_STATUS_PRESENT = 3; }
message SetMealAttendanceResponse { int64 attendance_version = 1; repeated string recomputed_plan_meal_ids = 2; }

message ListMealAttendanceRequest { string household_id = 1; string from = 2; string to = 3; }
message ListMealAttendanceResponse {
  repeated MealAttendanceEntry entries = 1; // resolved: rules + one-off marks, AWAY only
  repeated AttendanceRule rules = 2;
}
message MealAttendanceEntry { string member_id = 1; string date = 2; Meal meal = 3; bool from_rule = 4; }
message AttendanceRule {
  string id = 1; string household_id = 2; string member_id = 3; Meal meal = 4;
  repeated int32 weekdays = 5; // ISO 1 = Monday … 7 = Sunday
  string valid_from = 6; string valid_to = 7; // valid_to empty = open-ended
}
message SetAttendanceRuleRequest { AttendanceRule rule = 1; } // id empty = create
message DeleteAttendanceRuleRequest { string household_id = 1; string rule_id = 2; }
message DeleteAttendanceRuleResponse {}
```

## Data

`planner/migrations/000NN_meal_attendance.sql`

```sql
CREATE TABLE meal_attendance (
    household_id uuid     NOT NULL,
    member_id    uuid     NOT NULL,
    date         date     NOT NULL,
    meal         smallint NOT NULL,          -- planner.v1.Meal
    status       smallint NOT NULL,          -- 2 away, 3 present
    updated_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, member_id, date, meal)
);
CREATE TABLE attendance_rules (
    id           uuid PRIMARY KEY,
    household_id uuid     NOT NULL,
    member_id    uuid     NOT NULL,
    meal         smallint NOT NULL,
    weekdays     smallint NOT NULL CHECK (weekdays BETWEEN 1 AND 127), -- bit 0 = Monday
    valid_from   date     NOT NULL,
    valid_to     date     CHECK (valid_to IS NULL OR valid_to >= valid_from)
);
CREATE INDEX attendance_rules_household ON attendance_rules (household_id);
CREATE TABLE attendance_versions (
    household_id uuid PRIMARY KEY,
    version      bigint NOT NULL DEFAULT 0
);
```

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| Unknown meal, bad date, empty weekdays | `INVALID_ARGUMENT` | |
| Member not in the household (profile snapshot) | `NOT_FOUND` | |
| Date in the past or meal already cooked | `FAILED_PRECONDITION` | `attendance_locked` |
| More than 20 rules | `RESOURCE_EXHAUSTED` | |

## Acceptance criteria

- [ ] Worked example household (Adam, Ewa, Ola, Kuba) with chicken curry for dinner: pot 6; marking Ola away for
      dinner makes the pot 5 and the dinner's list contribution shrinks by one portion (Verification 17).
- [ ] A rule "Ola, dinner, Mon–Fri" removes Ola from every weekday dinner in a week plan; a one-off PRESENT on
      Wednesday puts her back for that day only.
- [ ] Clearing a one-off mark (DEFAULT) restores the rule's or the default state.
- [ ] Changing a mark for yesterday fails with `attendance_locked`.
- [ ] Marks never contain names; logs show member ids only.

## Notes

- Decided (Q196): weekly rules exist ("school lunch" repeats every weekday), at most 20 per household.
