# AI-15-02 Week composition prompt and call

| | |
|---|---|
| Service | ai |
| Epic | [AI-15 Week plan](README.md) |
| Release | R3 |
| Size | L |
| Depends on | AI-15-01 |
| Unblocks | AI-15-03; PLN-14 |

## Goal

Plan the week so the household cooks as little as its rhythm allows while every day stays varied and uses the
oldest stock first.

## Business rules

1. Job `compose_week_plan`, Sonnet 5.5 at effort `high`, `max_tokens` 32000, so the call streams
   (AI-01-03 rule 5). Prompt `compose_week_plan/v1`.
2. Rules in the order features.md §2 sets when aims conflict: hard limits (allergens, diets, tools) → the household's
   cooking rhythm → each eater's day within reach (planner sizes portions; AI prefers dishes that portion well, as
   AI-12-02 rule 2.4) → oldest stock first → this week's offers (not used here; deals stay in AI-14).
3. Cooking days: `EVERY_DAY` → the main meal cooked daily; `EVERY_OTHER_DAY` → at most 4 cooking days in 7;
   `TWICE_A_WEEK` → at most 2. Preferred dates are used when given.
4. On a cooking day the main dish is cooked as a double batch when it keeps 3 days (`KEEPS_3_DAYS` tag or
   `keeps_days ≥ 3`) or freezes well; the second half becomes one later main meal within its keeping window, or any
   later day if frozen. Dinner leftovers may also become a later supper.
5. Breakfast and the non-main meal are no-cook dishes unless the meal request says otherwise.
6. Library first (cheaper, checked); at most `max_new_dishes` new dishes in the week; no dish more than twice in the
   week counting leftovers; stock expiring early is used in the first days.
7. Expected size ≈ 10,000 input and ≈ 10,000 output tokens (≈ $0.12 on Sonnet; AI-02).

## Contract

As AI-15-01.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Input limits broken, unknown ids or codes in the request | `INVALID_ARGUMENT`, `ai_bad_request` |
| Output still invalid after repair | `UNAVAILABLE`, `ai_invalid_output` |
| Others | AI-01-04; the caller falls back to its non-AI result |

## Acceptance criteria

- [ ] Every-other-day rhythm for 7 days gives ≤ 4 cooking days, with leftovers on the days between.
- [ ] A dish without keep or freeze tags is never used as a leftover more than 1 day later.
- [ ] The call streams and completes within the 110 s timeout on the golden set (p95).
