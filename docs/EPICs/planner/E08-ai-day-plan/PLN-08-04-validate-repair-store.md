# PLN-08-04 Validate, repair once, store AI dishes

| | |
|---|---|
| Service | planner |
| Epic | [PLN-08 AI day plan](README.md) |
| Release | R1 |
| Size | L |
| Status | todo |
| Depends on | PLN-07-05, PLN-08-03, PLN-04-01, PLN-02-04, AI-10 (ComposeDayPlan) |
| Unblocks | PLN-08-06, PLN-08-07, PLN-09 |

## Goal

Turn an AI answer into a day plan users can trust: check every meal, give AI one chance to fix what broke, replace any
meal that still fails with a library dish, and store new AI dishes with calories from the tables.

## Business rules

1. For each meal in the AI response: library dish id → load it and run the checker for these eaters (library dishes
   can still break a dislike or a stock rule); new dish → map to `CandidateDish` and run the full checker.
2. Plan-level allocation (PLN-08-03) runs over all meals before the stock check, so cross-meal shortages are found.
3. If any meal fails, planner makes **one** repair call: `ComposeDayPlan` again with the passing meals fixed and, for
   each failing meal, its violation codes and the lines/tools/members concerned (English, no user text). The repair
   request's other fields are identical. Only the failing meals are replaced from the repair answer.
4. After the repair, any meal still failing is replaced by the best library dish for that meal against the remaining
   stock (PLN-08-05 logic) with `source=library_fallback`. If the repair call itself fails (unavailable, refusal,
   budget), the same replacement applies.
5. Meals that passed first time have `source=ai`; repaired ones `ai_repaired`.
6. New AI dishes are stored as `household_ai` dishes (status `household`, the household's id) with version 1, the
   names and steps AI wrote, derived flags (PLN-02-03) and nutrition computed by kitchen (PLN-04-01). Nothing
   numeric about energy is taken from AI.
7. AI names for new ingredients are not accepted: AI may only use ids it was given (PLN-07-04 rule 3). A dish that
   references a name instead of an id is a `STRUCTURE` violation.
8. `ai_used=true` whenever at least one meal came from AI (`ai` or `ai_repaired`). If every meal ended as
   `library_fallback`, `ai_used=true` still (the call was made and paid) but `fallback_reason=AI_INVALID`.
   Decided (Q120): such a plan does not use Free quota; the gateway counts only when at least one AI meal survived
   (field `ai_meals_kept`).
9. Metrics from PLN-07-05; additionally `planner_day_plans_total{result}` (`ai`, `ai_repaired`, `partial_fallback`,
   `full_fallback`, `library`).

## Contract

Internal; uses `AiService.ComposeDayPlan` with a repair section (AI-10 must support: fixed meals + failed meals with
violations). The response field `ai_meals_kept` is added to `ComposeDayPlanResponse` (PLN-08-07).

## Data

Uses `dishes`, `dish_versions` etc. (PLN-02-01) and `plans` (PLN-08-01).

## Errors

| Case | Behaviour |
|---|---|
| AI response not parseable | Treated like a failed repair: library fallback for all meals |
| Kitchen down while storing nutrition | AI dish stored with `pending=true` (PLN-04-01 rule 8) |

## Acceptance criteria

- [ ] Recorded AI answer where supper uses an oven in a microwave-only kitchen → one repair call; if the repair passes,
      supper is `ai_repaired`; if not, `library_fallback`. Breakfast and dinner unchanged.
- [ ] At most 2 ai calls per day plan, ever.
- [ ] A stored AI dish has calories from `ComputeNutrition` and none of its texts contain a calorie number.
- [ ] Verification 2 and 3 at plan level: 50 recorded day plans pass through, and no shown dish breaks tools or
      "no gluten, no nuts".
