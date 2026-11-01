# PLN-07-01 Checker core, candidate dish and verdict

| | |
|---|---|
| Service | planner |
| Epic | [PLN-07 Rule checker](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-02-02, PLN-01-03 |
| Unblocks | PLN-07-02 … PLN-07-05, PLN-08-04, PLN-09, PLN-03-07 |

## Goal

One checker with one result type, used for AI dishes (always) and library dishes (in review and before reuse in a plan
for new eaters), so every rule lives in exactly one place.

## Business rules

1. Input: a `CandidateDish` (same fields as a dish version: meals, cuisine, tools, times, servings, ingredient lines with
   ids, amounts, units, optional, role, topping group, steps with tools and texts) + a `CheckContext`: eaters (member
   ids), the household's tools, hard limits per eater (PLN-01-03), time budget for the meal, stock mode (`FROM_STOCK`
   or `ANY`), allowed purchases (ingredient ids the plan may buy), the candidate ingredient ids sent to AI, language.
2. AI responses are converted to `CandidateDish` by one mapper; an AI field that can't be mapped (unknown enum,
   missing field) is a `STRUCTURE` violation, never a crash.
3. Output `Verdict`: `pass` (no violation), `violations` (blocking), `warnings` (shown or logged, never blocking),
   `plates` (per eater for base-plus-toppings dishes: allowed topping options, PLN-07-03), `leave_out` per eater.
4. Each finding: `code`, severity, the line / step / tool / member it concerns, and a short English text for logs and
   admin (no user free text).
5. Checks run in a fixed order and all run (the verdict lists every problem, so one repair request can fix them all).
6. The checker is pure (no network) apart from the inputs; ingredient flags and conversions are passed in, pre-fetched.
7. Codes: `STRUCTURE`, `UNKNOWN_INGREDIENT`, `NOT_ALLOWED_INGREDIENT`, `CALORIES_IN_TEXT`, `TOOL_MISSING`,
   `ALLERGEN`, `DIET`, `DISLIKE`, `STOCK_SHORT`, `TIME_OVER`, `MEAL_MISMATCH`; warnings `TIME_TIGHT`,
   `CUISINE_DIFFERS`, `LOW_STOCK_MARGIN`, `CHECK_FIRST` (stock past its date, Q103).

## Contract

Internal Go package `rules`. Codes are also exposed in `ValidationProblem.code` (PLN-02-02) for admin views.

## Data

No schema change.

## Errors

None of its own; internal logic called by PLN-08-04, PLN-09 and PLN-03-07. Rule violations are results (`ValidationProblem` codes), not errors.

## Acceptance criteria

- [ ] An AI response with an unknown tool string yields one `STRUCTURE` violation and the other checks still run.
- [ ] A dish with a missing tool and an allergen returns both violations in one verdict.
- [ ] Unit tests run with fixtures only (no kitchen, no ai).
