# GW-07-11 Week plan

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R3 |
| Size | M |
| Depends on | GW-03-03, GW-01-03, PLN-14 (`CreateWeekPlan`, `GetWeekPlan`, `ReplaceWeekMeal`, `DeleteWeekPlan`) |
| Unblocks | MOB-03 |

## Goal

A 7-day plan with the fewest cooking sessions: cook a double amount on cooking days and eat it later, dinner
leftovers reused in a later meal (features.md §2, verification 18). Plus only.

## Business rules

1. Plus-only (`plus: "week_plan"`). `POST /v1/plans/week` body `{ "startDate", "replace"?: false, "useAi"?: true,
   "pinnedCookingDays"?: [1–7] }`. `startDate` from today to today + 13 (planner's range); default next Monday.
   A week covers 7 consecutive dates; weeks may not overlap (PLN-14-01 rule 1).
2. Calls planner `CreateWeekPlan(household_id, viewer_member_id, start_date, replace, use_ai,
   pinned_cooking_days)` with a **longer deadline: 120 s** (a dedicated `ai_long` class, config); idempotency
   required. An existing week for that start date needs `replace: true`; without it planner refuses and the gateway
   answers 409 `plan_exists`.
3. Response: planner's `WeekPlan`: `id`, `version`, the 7 day plans (family day shape, GW-07-04 rule 9), cooking
   `sessions` (cook day and meal, amount factor, the later meal it feeds and fridge/freezer), `effort`
   (features.md: "3 cooking sessions this week, about 2 hours at the stove"), per-person week totals filtered by
   visibility, `aiUsed`, `fallbackReason`.
4. AI unavailable → planner builds a non-AI week from the library (PLN-14-06), `aiUsed: false`.
5. `GET /v1/plans/week?startDate=` → `GetWeekPlan`. A household that dropped to Free can still read its weeks.
6. `PUT /v1/plans/week/{weekPlanId}/meals` body `{ "date", "meal", "dishId", "baseVersion" }` → `ReplaceWeekMeal`
   (no AI, Plus). Replacing a cooking meal's dish also replaces its leftover meal (PLN-14-01). A stale
   `baseVersion` → planner `ABORTED` → 409 `version_conflict` with the current version.
7. `DELETE /v1/plans/week/{weekPlanId}` → `DeleteWeekPlan`; day plans of dates not yet cooked go with it. 204.

## Contract

| Method | Path | 2xx |
|---|---|---|
| POST | `/v1/plans/week` | 200 week plan |
| GET | `/v1/plans/week?startDate=2026-11-09` | 200 week plan |
| PUT | `/v1/plans/week/{weekPlanId}/meals` | 200 week plan |
| DELETE | `/v1/plans/week/{weekPlanId}` | 204 |

```json
200 { "id": "wk_...", "startDate": "2026-11-09", "version": 1, "days": [ { "...": "" } ],
      "sessions": [ { "date": "2026-11-09", "meal": "dinner", "dishId": "...", "amountFactor": 2,
                      "leftoverDate": "2026-11-10", "leftoverMeal": "dinner", "storage": "fridge" } ],
      "effort": { "cookingSessions": 3, "minutesAtStove": 120 }, "aiUsed": true, "fallbackReason": "" }
```

Tag `plans`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| Bad start date | 400 | `validation` |
| A week already starts on that date (no `replace`) or weeks overlap | 409 | `plan_exists` |
| Stale `baseVersion` on a meal edit | 409 | `version_conflict` |

## Acceptance criteria

- [ ] For the test family with "every other day", the 7-day plan has at most 4 cooking sessions (verification 18).
- [ ] Creating a week for a start date that already has one returns 409 unless `replace: true`.
- [ ] Replacing a cooking meal's dish with an old `baseVersion` returns 409; with the current one, the leftover meal
      changes too.
- [ ] A request waits up to 120 s before 504.
