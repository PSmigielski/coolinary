# GW-03-03 Plus-only route guard

| | |
|---|---|
| Service | gateway |
| Epic | [GW-03 Free-tier AI quota and rate limits](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-02-01 |
| Unblocks | Plus-only routes in GW-04, GW-07, GW-08, GW-10, GW-12 |

## Goal

Plus-only features answer Free users immediately with `plus_required`, before any service (and any AI) is called.
Services check again with `x-plan` (D5); the gateway check saves the call and gives one consistent error.

## Business rules

1. Route option `plus: "<feature>"`. If the JWT plan is not `plus` → 403 problem `plus_required`,
   `meta: { "feature": "<feature>", "paywall": true }`.
2. Plus-only routes and their feature names (from docs/EPICs/README.md "Free and Plus"):

   | Feature | Routes |
   |---|---|
   | `swap_meal` | `POST /v1/plans/{planId}/meals/{meal}/replace-ingredient`, `POST /v1/plans/wish` (a plain swap is open to Free, which gets a library dish without AI, PLN-09-01 rule 5) |
   | `tool_rewrite` | `POST /v1/dishes/{dishId}/rewrite` |
   | `restyle` | `POST /v1/dishes/{dishId}/restyle` |
   | `week_plan` | `POST /v1/plans/week` |
   | `plan_from_deals` | `POST /v1/plans/from-deals` |
   | `cheapest_top_up` | `POST /v1/list/cheapest-top-up` |
   | `own_flyer` | `POST /v1/flyers` |
   | `trip_plan` | `POST /v1/list/price-matrix` |
   | `shared_household` | `POST /v1/household/invites` |
   | `profiles` | `POST /v1/household/profiles` |
   | `savings_report` | `GET /v1/savings/reports`, `GET /v1/savings/card` |

3. Rules that depend on data, not on the route (Free may set a goal for itself but not for others; Free may pick 2
   shops; accepting an invite into a Plus household) are enforced by the owning service, which answers
   `FAILED_PRECONDITION` + `plus_required`; GW-01-02 maps it to the same 403.
4. The JWT plan can be up to 15 minutes old. The app, on `plus_required`, refreshes the JWT once and retries
   (MOB-02); the gateway does not call identity here.

## Contract

Elysia macro `plus` on route options; OpenAPI marks these routes with `x-plus: true` and lists the 403 response.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free or anonymous on a Plus route | 403 | `plus_required` |

## Acceptance criteria

- [ ] A Free token on `POST /v1/plans/week` gets 403 `plus_required`, `meta.feature: "week_plan"`, and planner
      receives no call.
- [ ] A Plus token reaches planner.
- [ ] A service answering `FAILED_PRECONDITION` + `plus_required` produces the same 403 shape.
- [ ] Every route in the table carries `x-plus: true` in the OpenAPI document.
