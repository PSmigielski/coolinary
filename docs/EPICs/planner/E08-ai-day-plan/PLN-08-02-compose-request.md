# PLN-08-02 Compose request assembly

| | |
|---|---|
| Service | planner |
| Epic | [PLN-08 AI day plan](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-01-03, PLN-01-04, PLN-06-02, PLN-05-06, AI-10 (ComposeDayPlan) |
| Unblocks | PLN-08-04, PLN-09 |

## Goal

Send AI exactly what it needs to compose a good, checkable day — stock oldest first, limits, cuisines, tools, a short
list of library dishes it may reuse — and nothing personal.

## Business rules

1. Meals: requested meals (default breakfast, dinner, supper). For each: cuisine (request override, else the profile's
   cuisine for that meal, else "no preference"), time budget from the profile, whether it is the household's main hot
   meal.
2. Non-main meals carry the preference "prefer no-cooking dishes" unless the household turned it off (features.md
   "One dish per meal"). Decided (Q119): this applies as a soft preference from R1, not only in family plans.
3. Servings: the household's portions (R1). Family targets are added by PLN-13.
4. Tools: the household's tool codes. Hard limits: union of allergens and diets of the eaters; dislikes as ingredient
   ids with names.
5. Stock: all matched items, including those past their date, as `{ingredient_id, name, quantity, unit, days_left,
   past_date}` (`days_left` < 0 and `past_date = true` when past the date, Q103), sorted oldest first
   (PLN-01-04); at most 80 lines (most urgent first, then most quantity), summed per ingredient. Unmatched items are
   left out.
6. Staples: `{ingredient_id, name}`. Allowed purchases: none, unless the request has them.
7. Library candidates: per meal, the top 10 ranked library dishes (PLN-06-02) after hard limits, as `{dish_id, name,
   share_on_hand, missing names}`. AI may answer with a library dish id instead of a new dish (cheaper to check and
   already has calories).
8. Recently cooked dish names (last 7 days), to avoid repeats.
9. Language of the output = the request's `x-language`; ingredient names are sent in that language.
10. Never sent: member names, ages, weights, targets (R1), user ids, free text other than a wish (PLN-09-03).
    Household id only as the hashed trace user (ai side, existing).
11. The request is built deterministically (stable ordering), so ai's result cache (existing `cache.py`) hits when the
    same household asks twice with unchanged stock and profile.
12. Deadline 45 s (PLN-01-02).

## Contract

Calls `AiService.ComposeDayPlan`. The fields above are what planner needs from AI-10's request message; AI-10 owns the
exact message. The response must give per meal either `library_dish_id` or a new dish using only the ingredient ids
sent, plus a per-meal `cuisine_fit` flag and, when false, the three choices of PLN-08-06.

## Data

No schema change.

## Errors

None of its own; builds the request for `AiService.ComposeDayPlan` inside PLN-08-07. An ai error or a deadline over 45 s (rule 12) is handled there by the library fallback (PLN-08-05).

## Acceptance criteria

- [ ] A request for a household with 120 stock items sends 80, the 80 most urgent.
- [ ] A golden test asserts the serialised request has no field named or containing age, weight, height or e-mail.
- [ ] Two calls with unchanged stock and profile produce byte-identical requests.
