# PLN-14-05 Week balance

| | |
|---|---|
| Service | planner |
| Epic | [PLN-14 Week plan](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | PLN-12-05, PLN-14-03 |
| Unblocks | PLN-14-06, PLN-14-07 |

## Goal

Each person's days land within 10% of their target and their week within 5% (features.md "How one meal is split",
step 6; Verification 18).

## Business rules

1. Days are allocated in date order. A person's remainder at the end of a day (the gap after supper, PLN-12-05)
   becomes the carry-in of their first meal the next day, capped at ±10% of that day's meal part.
2. Week check: `planned week = Σ planned day` over the 7 days; `week target = Σ day target part` (PLN-12-05 rule 7).
   In range when `|planned week − week target| ≤ 5% × week target`; otherwise `week_out_of_range` with the deviation
   in percent.
3. Day checks stay as in PLN-12-05 (10%).
4. Leftover meals have fixed pots; their portions follow PLN-14-03 rule 5, and their gaps carry over like any meal.
5. Output per person: 7 day totals and in-range flags, the week total, deviation, in range. All kcal fields pass the
   visibility filter (PLN-12-07).
6. When a person is out of range, the composers (PLN-14-06, PLN-14-07) try alternatives through comparator key 3;
   the flags are shown only to viewers who may see that person's numbers.

## Contract

```proto
message PersonWeek {
  string member_id = 1;
  repeated PersonDay days = 2;          // PLN-13-07 message
  optional int32 planned_kcal = 3;
  optional int32 target_kcal = 4;
  optional int32 deviation_percent = 5;
  optional bool in_range = 6;
}
```

## Data

`planner/migrations/000NN_week_balance.sql`

```sql
CREATE TABLE week_person_balance (
    week_plan_id      uuid NOT NULL REFERENCES week_plans ON DELETE CASCADE,
    member_key        text NOT NULL,
    planned_kcal      integer NOT NULL,
    target_kcal       integer NOT NULL,
    deviation_percent smallint NOT NULL,
    PRIMARY KEY (week_plan_id, member_key)
);
```

## Errors

None of its own: computed inside the week-plan RPCs (PLN-14-06, PLN-14-07), whose errors apply. A person without a target (no goal set) is returned with the kcal fields empty, not an error.

## Acceptance criteria

- [ ] Test family of four with different goals: every day within 10% and every week within 5% (Verification 18).
- [ ] A person 150 kcal over on Monday starts Tuesday's breakfast with −150 carry-in (if within the 10% cap).
- [ ] A teen account viewing the week sees no kcal, no target and no deviation.
