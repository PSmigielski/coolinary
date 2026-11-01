# PLN-14-02 Cooking days from the rhythm

| | |
|---|---|
| Service | planner |
| Epic | [PLN-14 Week plan](README.md) |
| Release | R3 |
| Size | S |
| Status | todo |
| Depends on | PLN-14-01, KIT-10 (GetPlanningProfile) |
| Unblocks | PLN-14-03, PLN-14-06, PLN-14-07, PLN-13-03 |

## Goal

Turn the household's chosen cooking rhythm into the days on which the main hot meal is cooked, and into the
session allowance the comparator uses (features.md §1 "Household": every day, every other day, or twice a week).

## Business rules

1. Rhythm comes from the profile snapshot: `every_day`, `every_other_day`, `twice_a_week`.
2. Days are numbered 1–7 from the week's start date.
   - `every_day`: cooking days 1–7, allowance 7 sessions; no double cooking unless it saves a session the household
     asked to skip.
   - `every_other_day`: days 1, 3, 5, 7, allowance 4.
   - `twice_a_week`: days 1 and 4, allowance 2.
3. Cooking days whose main meal has no eaters (everyone away) move to the next day with eaters; the allowance is
   unchanged.
4. Breakfast and the other meal are "no cooking" unless the household cooks them (profile); if they are cooked, each
   cooked one adds its own session to the allowance.
5. Main meals not covered by a cooking day or its leftover get a "no cooking" dish; if none passes the limits, an
   extra session is planned and counted over the allowance (comparator key 2 penalises it).
6. The user may pin cooking days for a week (e.g. Saturday instead of Thursday); pinned days replace the defaults,
   same count.

## Contract

Internal (`compose.CookingDays(rhythm, start, attendance, pins)`); pins are a field on `CreateWeekPlanRequest`
(PLN-14-07): `repeated int32 pinned_cooking_days = 6;`.

## Data

No schema change.

## Errors

None of its own; internal logic called by the week composers (PLN-14-06, PLN-14-07). A rhythm with no cooking day in the week (all meals pinned or everyone away) returns an empty day list, not an error.

## Acceptance criteria

- [ ] `every_other_day` from Monday: Mon, Wed, Fri, Sun; allowance 4 (Verification 18).
- [ ] `twice_a_week` from Monday: Mon and Thu; Tue and Fri dinners are leftovers; Wed, Sat, Sun get no-cook dishes.
- [ ] Everyone away for Wednesday dinner: that session moves to Thursday.

## Notes

- Decided (Q217): `twice_a_week` cooks on days 1 and 4, double amounts only; the other main meals are no-cook.
