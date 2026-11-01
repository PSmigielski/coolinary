# MOB-03-06 Library, plans and cooking hooks

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R1 (library, suggestions, day plan, swaps, cooked it), R2 (away, tool rewrite, family plan), R3 (week plan, restyle) |
| Size | L |
| Status | todo |
| Depends on | MOB-03-01, MOB-03-02, MOB-02-05, MOB-04-01, GW-07 |
| Unblocks | later library, plan and cooking-mode screens |

## Goal

Hooks for the recipe library, "what can I make now", AI day and week plans, swaps and the "cooked it" button, with
the AI calls clearly separated from the offline-safe calls.

## Business rules

1. Hooks and routes (GW-07 route catalogue in docs/EPICs/gateway/README.md):

   | Hook | HTTP | Rel. | Offline |
   |---|---|---|---|
   | `useDishes(filters)` | `GET /v1/dishes?meal=&cuisine=&tools=&maxMinutes=&kcal=lt400\|400to600\|gt600&protein25=&effort=&q=&pageToken=` | R1 | first page cached |
   | `useDish(id, portions?)` | `GET /v1/dishes/{id}?portions=` (steps, timers, nutrition) | R1 | last 50 cached |
   | `useFavourite()` | `PUT` / `DELETE /v1/dishes/{id}/favourite` | R1 | queued |
   | `useRateDish()` | `PUT` / `DELETE /v1/dishes/{id}/rating` `{ stars: 1..5 }` | R1 | queued |
   | `useDismissRatingPrompt()` | `POST /v1/dishes/{id}/rating-prompt/dismiss` | R1 | queued |
   | `useHistory()` | `GET /v1/history?pageToken=` | R1 | cached |
   | `useSuggestions(meal)` | `GET /v1/suggestions?meal=breakfast\|dinner\|supper&limit=20` | R1 | cached; recomputed online |
   | `useComposeDay()` | `POST /v1/plans/day` `{ date, meals?, cuisines?, useAi?, replace? }` | R1 | never (AI) |
   | `useChooseMealOption()` | `POST /v1/plans/{id}/meals/{meal}/choice` `{ option }` (cuisine choice) | R1 | never |
   | `usePlans(from, to)` | `GET /v1/plans?from=&to=` (≤ 31 days) | R1 | cached |
   | `useAcceptPlan()` / `useDeletePlan()` | `POST /v1/plans/{id}/accept` / `DELETE /v1/plans/{id}` | R1 | never / queued |
   | `useSetPlanMeal()` | `PUT /v1/plans/{id}/meals/{meal}` `{ dishId }` (a library dish picked by hand; rule-checked, no AI, no quota; decided Q117) | R1 | never; a `422` lists the broken rules |
   | `useSwapMeal()` | `POST /v1/plans/{id}/meals/{meal}/swap` `{}` (Free gets a library dish) | R1 | never (AI) |
   | `useReplaceIngredient()` | `POST /v1/plans/{id}/meals/{meal}/replace-ingredient` `{ ingredientId }` (Plus) | R1 | never (AI) |
   | `useWish()` | `POST /v1/plans/wish` `{ planId? \| date?, meal, text }` (Plus) | R1 | never (AI) |
   | `useCookingSteps(dishId)` | `GET /v1/cooking/{dishId}?version=&planMealId=&servings=` | R1 | last 20 cached |
   | `useCooked()` | `POST /v1/cooked` `{ dishId, dishVersion, planMealId?, servings, cookedAt, chosenToppings?, leftOutIngredientIds? }` | R1 | queued |
   | `useMealPortions(planMealId)` | `GET /v1/plan-meals/{id}/portions` (per-person portions, visibility applied by the server) | R2 | cached |
   | `useMarkAway()` | `PUT /v1/attendance` `{ memberId, date, meal, status: "away" \| "present" \| "default" }` | R2 | queued |
   | `useAttendanceRules()` | `GET /v1/attendance`, `POST`/`PUT`/`DELETE /v1/attendance/rules[/{id}]` | R2 | cached / queued |
   | `useRewriteForTools()` | `POST /v1/dishes/{id}/rewrite` `{ tools }` | R2 | never (AI) |
   | `useWeekPlan()` / `useComposeWeek()` | `GET` / `POST /v1/plans/week` `{ startDate, replace?, useAi?, pinnedCookingDays? }` | R3 | cached / never (AI) |
   | `useReplaceWeekMeal()` | `PUT /v1/plans/week/{id}/meals` `{ date, meal, dishId, baseVersion }` | R3 | never |
   | `useRestyle()` | `POST /v1/dishes/{id}/restyle` `{ cuisine }` | R3 | never (AI) |
2. Library pages use `useInfiniteQuery` with the server's `pageToken`; page size 20; filters are part of the key, and
   the profile's hard filters (tools, allergens, diet) are applied by the server, not by the phone.
3. AI hooks: check `usePlusGate` first (day plan: `day_plan` quota; swaps, rewrites, week plan: Plus on Free per the
   README table), set the AI timeout class (MOB-02-01), and pass a fresh `Idempotency-Key` per tap. A response flag
   `aiUsed: false` means planner returned the non-AI fallback; the hook exposes it so the UI can say so.
4. Free-text wishes (`useWish`) are 1–200 characters; they go only to the gateway, never into logs or the cache key.
5. `useCooked()` sends `cookedAt` (the phone's time, ISO 8601 with offset) so a cooked-it tapped offline reduces stock
   with the right timestamp when replayed; the server is idempotent on the `Idempotency-Key`.
6. Plan responses are cached per date range; a compose or swap writes the new plan into the cache from its response
   instead of refetching.
7. Portions per person (R2) are shown from the dish/plan response as given by planner; the phone does no portion
   arithmetic in this layer.

## Contract

Types aliased in `src/api/types/planner.ts`. All AI responses include `{ aiUsed: boolean }` (assumed).

## Data

Query cache per MOB-03-01; AI responses cached only as the resulting plan.

## Errors

| Case | Behaviour |
|---|---|
| Day plan over Free quota | `quotaReached` |
| Swap on Free | `plusRequired` |
| AI unavailable | planner falls back (`aiUsed: false`); only routes without fallback give `aiUnavailable` |

## Acceptance criteria

- [ ] Verification 1 from the app side: with ai stopped, `useComposeDay` returns library-based meals with
      `aiUsed: false` and the free quota is unchanged.
- [ ] Cooked-it tapped offline reduces stock after reconnect, once, even if replayed twice.
- [ ] Scrolling the library loads pages of 20 and keeps the first page available offline.
