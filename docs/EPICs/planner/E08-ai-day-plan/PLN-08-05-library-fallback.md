# PLN-08-05 Library day plan (non-AI fallback)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-08 AI day plan](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-06-02, PLN-08-03, PLN-07-03 |
| Unblocks | PLN-08-06, PLN-08-07 |

## Goal

A full day plan from the checked library when AI is not allowed (Free quota used), down, over budget or refuses — so the
app keeps working with AI switched off (features.md rule 1).

## Business rules

1. Meals are filled in allocation order (PLN-08-03 rule 1). For each meal: rank library dishes for that meal against
   the **remaining** stock (after earlier meals' allocation), with hard limits for the eaters and the meal's cuisine
   as a filter.
2. Pick the top-ranked dish. If no dish of that cuisine has share ≥ 0.5, the meal gets the three choices of PLN-08-06
   (library version) instead of a dish.
3. A dish already chosen for another meal of the same day is skipped.
4. Non-main meals prefer `no_cook` dishes (PLN-08-02 rule 2): among dishes within 0.1 of the top score, a `no_cook`
   dish wins.
5. Time budget: dishes over the meal's budget + 10 minutes are skipped.
6. `ai_used=false`; `fallback_reason` = `AI_NOT_ALLOWED`, `AI_UNAVAILABLE`, `AI_BUDGET`, `AI_REFUSED`, `AI_TIMEOUT` or
   `NONE` (library chosen on purpose, e.g. "plan without AI" button).
7. Runs in < 200 ms for a 2,000-dish library; no external call except cached profile and stock.

## Contract

Internal; result returned by `ComposeDayPlan` (PLN-08-07).

```proto
enum FallbackReason { FALLBACK_REASON_UNSPECIFIED = 0; FALLBACK_REASON_NONE = 1; FALLBACK_REASON_AI_NOT_ALLOWED = 2;
  FALLBACK_REASON_AI_UNAVAILABLE = 3; FALLBACK_REASON_AI_BUDGET = 4; FALLBACK_REASON_AI_REFUSED = 5;
  FALLBACK_REASON_AI_TIMEOUT = 6; FALLBACK_REASON_AI_INVALID = 7; }
```

## Data

No schema change.

## Errors

None of its own; it is the fallback path inside `ComposeDayPlan` (PLN-08-07). When the library has no dish for a meal, that meal is returned empty with `no_match = true`, not an error.

## Acceptance criteria

- [ ] With ai stopped and the Verification 1 test fridge, `ComposeDayPlan` returns three meals (or a choice where the
      cuisine doesn't fit) with `ai_used=false, fallback_reason=AI_UNAVAILABLE`.
- [ ] The same dish never appears twice in one day.
- [ ] Two meals needing the same 4 eggs out of 4 in stock: the second meal is ranked against 0 eggs.
