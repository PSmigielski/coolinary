# PLN-07-05 AI output checks, time budget, metrics, 50-dish harness

| | |
|---|---|
| Service | planner |
| Epic | [PLN-07 Rule checker](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-07-01 … PLN-07-04, PLN-02-02 |
| Unblocks | PLN-08-04, PLN-09 |

## Goal

Close the remaining gaps in AI output (unknown ingredients, calorie numbers, wrong meal, time over budget), measure how
often AI breaks rules, and prove Verification 2 and 3 with a repeatable test.

## Business rules

1. All structural rules of the validator (PLN-02-02 rules 1–10) run on AI dishes; blocking problems become
   `STRUCTURE` (or the specific code).
2. Unknown ingredient id → `UNKNOWN_INGREDIENT`.
3. Any calorie or kilojoule number in any text → `CALORIES_IN_TEXT` (AI never states a calorie number, features.md
   rule 2). The text is not edited; the dish fails.
4. The dish's meals must include the meal it was composed for → else `MEAL_MISMATCH`.
5. Time: profile time budget for the meal (minutes). `total_minutes` > budget + 10 → `TIME_OVER`; budget < total ≤
   budget + 10 → warning `TIME_TIGHT`. A wish with its own time ("in 20 minutes") replaces the budget for that meal.
   Decided (Q108).
6. Cuisine differs from the requested one → warning `CUISINE_DIFFERS` (cuisine choice is handled in PLN-08-06).
7. Metric `planner_rule_violations_total{code, source}` (`source` = feature: `day_plan`, `swap`, `wish`, `rewrite`,
   ...), and `planner_ai_dishes_checked_total{source, result}` (`pass`, `repaired`, `fallback`).
8. Harness `go test ./internal/rules -run Verification -tags eval`: replays recorded AI responses from
   `planner/testdata/ai/*.json` (at least 50 per scenario) through the checker. Scenarios: "microwave + kettle" preset
   (Verification 2) and "no gluten, no nuts" (Verification 3). The test asserts: every dish that passes the checker
   needs no other tool / contains zero violations. Recording new responses is a manual, cost-approved step via AI-03.

## Contract

Internal; metrics on `/metrics`.

## Data

No schema change.

## Errors

None of its own; internal checks run by the checker core (PLN-07-01) on AI output. Metrics only.

## Acceptance criteria

- [ ] A recorded response with "about 450 kcal" in the intro fails with `CALORIES_IN_TEXT`.
- [ ] The Verification 2 harness passes with 50 recorded dishes: 0 passing dishes need a tool outside microwave/kettle.
- [ ] The Verification 3 harness passes: 0 passing dishes contain gluten or nuts in a non-optional line.
- [ ] Metrics show one increment per violation code per checked dish.
