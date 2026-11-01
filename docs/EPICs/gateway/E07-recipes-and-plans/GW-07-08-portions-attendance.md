# GW-07-08 Per-person portions and "away" marking

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R2 |
| Size | S |
| Depends on | GW-07-05, PLN-12 (`SetMealAttendance`, `ListMealAttendance`, `SetAttendanceRule`, `DeleteAttendanceRule`, `GetMealPortions`, `GetScaledIngredients`), PLN-13 (family day plan) |
| Unblocks | MOB-03 |

## Goal

Plans show each person's portion of the one shared dish, and a person can be marked away for a meal (school lunch,
eating out); the pot shrinks and the shopping list follows (features.md §2 "Planning for each person",
verification 17).

## Business rules

1. **Portions of a plan meal:** `GET /v1/plan-meals/{planMealId}/portions` → planner `GetMealPortions(household_id,
   plan_meal_id, viewer_member_id)`: the pot, each eater's portion and extra, warnings, and "calories incomplete"
   (PLN-12-06). Family day plans (GW-07-04 rule 9) already carry the same `portions` per meal.
2. **Visibility is planner's job** and the gateway passes the requester's `memberId` as `viewer_member_id`: `kcal`
   is present only for the requester, for profiles when the requester is an adult member, and never for another
   adult; members with hidden calories get no `kcal` at all, including their own (features.md §2 "Privacy",
   PLN-12-07). The gateway does not filter.
3. **Scaled ingredients:** `GET /v1/plan-meals/{planMealId}/ingredients` → planner `GetScaledIngredients`: the dish's
   ingredients scaled to the pot, with extras per person listed separately (PLN-12-08).
4. **One-off attendance:** `PUT /v1/attendance` body `{ "memberId", "date", "meal", "status": "away" | "present" |
   "default" }` → planner `SetMealAttendance(household_id, member_id, date, meal, status)`. `default` clears the
   one-off mark. Response `{ "attendanceVersion", "recomputedPlanMealIds" }`; the app refetches those meals.
5. **Attendance overview:** `GET /v1/attendance?from=&to=` → `ListMealAttendance`: resolved away marks and the
   weekly rules.
6. **Weekly rules** ("Ola has school lunch on weekdays"): `POST /v1/attendance/rules` (create) and
   `PUT /v1/attendance/rules/{ruleId}` body `{ "memberId", "meal", "weekdays": [1–7], "validFrom", "validTo"? }` →
   `SetAttendanceRule`; `DELETE /v1/attendance/rules/{ruleId}` → `DeleteAttendanceRule`.
7. Who may mark whom: planner's rule (PLN-12-01 rule 4: any account member for any member; Free households only
   their own account).
8. Changing attendance never calls AI and never uses quota. Idempotency optional on `PUT`/`DELETE`, required on
   `POST`.

## Contract

| Method | Path | Body | 2xx |
|---|---|---|---|
| GET | `/v1/plan-meals/{planMealId}/portions` | — | 200 `{ "potPortions": 6, "eaters": [ { "memberId", "portions": 1.5, "kcal"?: 780, "extra"? } ], "warnings": [], "caloriesIncomplete": false }` |
| GET | `/v1/plan-meals/{planMealId}/ingredients` | — | 200 `{ "factor": 1.5, "potPortions": 6, "ingredients": [...], "extras": [...] }` |
| PUT | `/v1/attendance` | `{ "memberId": "0192...ola", "date": "2026-11-03", "meal": "dinner", "status": "away" }` | 200 `{ "attendanceVersion": 7, "recomputedPlanMealIds": ["pm_..."] }` |
| GET | `/v1/attendance?from=&to=` | — | 200 `{ "entries": [...], "rules": [...] }` |
| POST, PUT | `/v1/attendance/rules`, `/v1/attendance/rules/{ruleId}` | rule | 200 rule |
| DELETE | `/v1/attendance/rules/{ruleId}` | — | 204 |

Tag `plans`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Member not in household | 422 | `rule_violated` |
| Free household marking another member | 403 | `plus_required` |
| Weekdays outside 1–7, `validTo` before `validFrom` | 400 | `validation` |

## Acceptance criteria

- [ ] Worked example (features.md §2): pot 6, Adam 1½, Ewa 1½, Ola 1, Kuba 2; marking Ola away makes the pot 5 and
      the next list proposal smaller (verification 17); `recomputedPlanMealIds` names the dinner.
- [ ] Ewa's response shows her own `kcal`, Ola's (child) `kcal`, and no `kcal` for Adam.
- [ ] A member with `caloriesHidden` sees portions without any `kcal`.
