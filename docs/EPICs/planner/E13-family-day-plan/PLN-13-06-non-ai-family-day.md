# PLN-13-06 Family day without AI

| | |
|---|---|
| Service | planner |
| Epic | [PLN-13 Family day plan](README.md) |
| Release | R2 |
| Size | L |
| Status | todo |
| Depends on | PLN-13-01, PLN-13-02, PLN-13-03, PLN-13-04, PLN-13-05, PLN-06, PLN-12-06 |
| Unblocks | PLN-13-07, PLN-14-06 |

## Goal

Compose a family day from the recipe library alone. It is the fallback when `ai` is down, over budget or refuses,
and the result whenever AI output fails the rule check twice (features.md rule 1: "The app stays useful with AI
switched off").

## Business rules

1. Inputs: date, the profile snapshot, attendance, stock (KIT-07 `ListStock`), and the household's cuisines per
   meal.
2. Candidates per meal: published library dishes for that meal type that pass the shared limits (PLN-13-01), the
   household's tools, and the meal's time limit. Breakfast and the meal that is not the main hot one use only dishes
   tagged "no cooking", unless the profile says the household cooks those meals too.
3. Cuisine: dishes in the meal's chosen cuisines first; other cuisines are used only when fewer than 3 candidates
   remain, and the result then says `cuisine_relaxed` for that meal.
4. Each candidate is checked for portions (PLN-12-03): candidates that are `UNSUITABLE` for an eater are dropped.
5. Take the top 8 candidates per meal by stock score (share of ingredients on hand, weighted by urgency, PLN-06),
   and evaluate every combination of breakfast × dinner × supper (at most 512) as a day.
6. Stock allocation across the day: an ingredient used by an earlier meal is reserved and is not on hand for a later
   one, so meals do not compete for the same items (features.md "Dishes and the day plan" with AI, same rule here).
7. Each day is allocated (PLN-12-06 logic, unsaved) and ranked with the comparator (PLN-13-03, PLN-13-04). The best
   day wins; second dishes only via PLN-13-05.
8. If a meal has no candidate at all, the day is returned with that meal empty and `no_candidate` set; it is never
   filled with a dish that breaks a hard limit.
9. The result is saved as a day plan (PLN-08 model) with `source = library`, `ai_used = false`, and allocations
   stored.
10. Runtime budget: at most 1 second for 512 combinations; otherwise fall back to 4 candidates per meal.

## Contract

Internal composer, called by `CreateFamilyDayPlan` (PLN-13-07) and by the week composer (PLN-14-06).

## Data

No schema change (uses PLN-08 and PLN-12-06 tables).

## Errors

| Case | Behaviour |
|---|---|
| kitchen down | `UNAVAILABLE`: no plan without stock and limits |
| Library empty for a meal | meal returned empty with `no_candidate` |

## Acceptance criteria

- [ ] With AI disabled and 10 items in the fridge, the family day has a dish for each meal, filtered by cuisine and
      tools, and portions for each person (Verification 1).
- [ ] A test fridge with 300 g chicken: if both dinner and supper candidates need 300 g chicken, only one meal uses it.
- [ ] Microwave-and-kettle preset: no chosen dish needs another tool (Verification 2).
- [ ] "No gluten, no nuts": zero violations across 50 composed days (Verification 3).
- [ ] Every person's day is within 10% where the library allows; otherwise flagged.
