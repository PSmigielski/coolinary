# GW-07 Recipe library, plans and cooking routes

Releases R1 (library, what-can-I-make, AI day plan, swaps, cooking mode), R2 (per-person portions, tool rewrite,
plan from deals) and R3 (week plan, restyle).

## Goal

The app browses the recipe library, asks what it can make now, gets an AI day plan (2 a week on Free), swaps a meal
(Plus), cooks with step timers and marks a dish cooked, which reduces stock. Planner owns all of it (decision 2);
the gateway maps routes, applies the quota and the Plus guard, and passes `member_id` so planner can apply
per-person visibility.

## AI routes in this epic

| Route | Quota (Free) | Plus-only | Rel. |
|---|---|---|---|
| `POST /v1/plans/day` | `dayPlan` 2/week | no | R1 |
| `POST /v1/plans/{planId}/meals/{meal}/swap`, `/replace-ingredient` | — | yes (`swap_meal`) | R1 |
| `POST /v1/dishes/{dishId}/rewrite` | — | yes (`tool_rewrite`) | R2 |
| `POST /v1/plans/from-deals` | — | yes (`plan_from_deals`) | R2 |
| `POST /v1/plans/week` | — | yes (`week_plan`) | R3 |
| `POST /v1/dishes/{dishId}/restyle` | — | yes (`restyle`) | R3 |

All AI routes use the `ai` deadline class, require an idempotency key, and return `aiUsed`; when planner falls back
to the library (`ai_used: false`) the response says why in `fallbackReason` (`ai_unavailable`, `ai_invalid`,
`ai_spend_limit`).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-07-01](GW-07-01-library-browse.md) | Library browse and dish detail | R1 | M | GW-01-04, PLN-05, PLN-04 |
| [GW-07-02](GW-07-02-favourites-ratings-history.md) | Favourites, ratings, history | R1 | S | GW-01-03, PLN-05 |
| [GW-07-03](GW-07-03-suggestions.md) | What can I make now | R1 | S | GW-01-02, PLN-06 |
| [GW-07-04](GW-07-04-ai-day-plan.md) | AI day plan and cuisine choices | R1 | M | GW-03-01, GW-03-02, PLN-08 |
| [GW-07-05](GW-07-05-plans-read-edit.md) | Read, accept, edit and delete plans | R1 | S | GW-07-04, PLN-08 |
| [GW-07-06](GW-07-06-swaps-wishes.md) | Swap a meal, replace an ingredient, free-text wish | R1 | M | GW-03-03, PLN-09 |
| [GW-07-07](GW-07-07-cooking-and-cooked.md) | Cooking mode and "cooked it" | R1 | M | GW-01-03, PLN-10 |
| [GW-07-08](GW-07-08-portions-attendance.md) | Per-person portions and "away" marking | R2 | S | GW-07-05, PLN-12, PLN-13 |
| [GW-07-09](GW-07-09-tool-rewrite.md) | Rewrite a dish for the tools at hand | R2 | S | GW-03-03, PLN-09 |
| [GW-07-10](GW-07-10-plan-from-deals.md) | Plan from deals | R2 | M | GW-03-03, GW-10-02, PLN-18 |
| [GW-07-11](GW-07-11-week-plan.md) | Week plan | R3 | M | GW-03-03, PLN-14 |
| [GW-07-12](GW-07-12-restyle.md) | Restyle a dish | R3 | S | GW-03-03, PLN-09 |
