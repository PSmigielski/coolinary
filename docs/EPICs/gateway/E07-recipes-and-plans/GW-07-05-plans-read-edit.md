# GW-07-05 Read, accept, edit and delete plans

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-07-04, PLN-08 (`GetPlans`, `AcceptPlan`, `DeletePlan`, `SetPlanMeal`) |
| Unblocks | MOB-03, GW-08-02 |

## Goal

The app shows the planned days, lets the user accept a proposed plan (which makes it the plan for that date) and
delete a plan, without AI and without quota.

## Business rules

1. `GET /v1/plans?from=&to=`: dates, range at most 31 days (planner's limit, PLN-08-01 rule 6), default today to
   today + 6. Calls planner `GetPlans(household_id, from, to)`, which returns accepted plans plus the newest
   proposal per date.
2. `POST /v1/plans/{planId}/accept` → planner `AcceptPlan`. Accepting supersedes the date's earlier accepted plan in
   the same transaction (PLN-08-01 rule 2). A plan with an open cuisine choice can be accepted only after an option
   is chosen (GW-07-04 rule 7); planner refuses with `FAILED_PRECONDITION` and its `x-error-code` → 422.
3. `DELETE /v1/plans/{planId}` → planner `DeletePlan`. Response 204.
4. AI edits to a meal go through swap, replace-ingredient and wish (GW-07-06) and attendance (GW-07-08).
5. **Manual placement (decided, Q117):** `PUT /v1/plans/{planId}/meals/{meal}` body `{ "dishId" }` puts a library
   dish the user picked into that meal of a day plan. Calls planner `SetPlanMeal(household_id, plan_id, meal,
   dish_id)` (PLN-08-01). Planner runs the rule checker (tools, allergens, diets of everyone eating, time); a dish
   that breaks a hard limit is refused with `FAILED_PRECONDITION` and its `x-error-code` (→ 422, `meta.problems`
   lists the rules broken). Stock shortfalls are warnings in the returned plan, not refusals. No AI, no quota.
   `meal` is `breakfast`, `dinner` or `supper`. Response 200 with the updated plan.
6. Idempotency: required on accept, optional on delete and on the `PUT` (it sets a value, so a repeat is harmless).
   None of these routes touches the quota.

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/plans?from=2026-11-02&to=2026-11-08` | 200 `{ "plans": [plan] }` |
| POST | `/v1/plans/{planId}/accept` | 200 plan |
| PUT | `/v1/plans/{planId}/meals/{meal}` `{ "dishId": "0192..." }` | 200 plan |
| DELETE | `/v1/plans/{planId}` | 204 |

Plan shape as in GW-07-04. Tag `plans`, `auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Range > 31 days | 400 | `validation` |
| Accepting with an open cuisine choice | 422 | planner's code |
| Plan or dish not found | 404 | `not_found` |
| Unknown meal | 400 | `validation` |
| Dish breaks a hard limit (tool, allergen, diet) | 422 | planner's code, `meta.problems` |

## Acceptance criteria

- [ ] Accepting a proposal makes it the plan returned for that date; the date's earlier accepted plan is no longer
      returned.
- [ ] A plan with an open cuisine choice can't be accepted until an option is chosen.
- [ ] None of these routes changes the quota.
- [ ] A 32-day range returns 400.
- [ ] `PUT` of a library dish into dinner returns the plan with that dish; the quota is unchanged.
- [ ] `PUT` of a dish that needs an oven in a household without one returns 422 with the tool rule in
      `meta.problems`.
