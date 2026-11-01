# GW-07-04 AI day plan and cuisine choices

| | |
|---|---|
| Service | gateway |
| Epic | [GW-07 Recipe library, plans and cooking routes](README.md) |
| Release | R1 |
| Size | M |
| Depends on | GW-03-01, GW-03-02, GW-01-03, PLN-08 (`ComposeDayPlan`, `ChooseMealOption`), PLN-13 (`CreateFamilyDayPlan`, R2) |
| Unblocks | MOB-03 |

## Goal

The core feature: one call composes breakfast, dinner and supper from what is at home, oldest items first, in the
chosen cuisines (features.md §2 "Dishes and the day plan", "Cuisine choice"). Free gets 2 a week; when AI is
unavailable planner answers from the library and the quota is not used.

## Business rules

1. Body `{ "date": "YYYY-MM-DD", "meals"?: [...] (default all three), "cuisines"?: { "<meal>": "<cuisine code>" }
   (one cuisine per meal, overrides for this plan only), "useAi"?: true }`. `date` from today to today + 13 days
   (household time zone is applied by planner).
2. Free and anonymous households (and Plus households until R2): order rate limit → idempotency (required) → quota
   reserve `dayPlan` (only when `useAi` is true) → planner `ComposeDayPlan(household_id, date, meals,
   cuisine_overrides, ai_allowed, request_id = idempotency key, language)` with the `ai` deadline (planner's own
   end-to-end limit is 55 s) → **commit the reserved unit only when planner answers `quota_consumed = true`**
   (AI used, at least one AI meal kept, not a replay; PLN-08-07 rule 5), else release.
3. `useAi: false` → `ai_allowed = false`, no quota is reserved, and planner answers from the library
   (`fallbackReason: "ai_not_allowed"`). A Free household over its weekly limit gets 429 `quota_reached`
   (GW-03-02) with `meta.libraryPlanAvailable: true`, so the app can offer a plan without AI.
4. A new proposal for a date replaces that date's earlier proposal; an accepted plan is never replaced by composing
   (PLN-08-07 rule 6). The user accepts the new proposal with `POST /v1/plans/{planId}/accept` (GW-07-05), which
   supersedes the old accepted plan.
5. A second compose for the same household while one is running waits in planner; after 50 s planner answers
   `ABORTED` "plan in progress" → 409 `plan_in_progress`, and the reserved unit is released.
6. Response: the plan (status `proposed`) with its meals (dish summary, `source`: `ai` | `library` | `ai_repaired` |
   `library_fallback`), `aiUsed`, `aiMealsKept`, `fallbackReason`. A meal whose stock does not fit the chosen
   cuisine has no dish and up to 3 `options` (features.md §2 "Cuisine choice"): `closest_now`, `buy_few` (with the
   1–3 ingredients to buy) and `other_cuisine` (naming the cuisine), each a full checked dish (PLN-08-06).
7. Choosing an option: `POST /v1/plans/{planId}/meals/{meal}/choice` body
   `{ "option": "closest_now" | "buy_few" | "other_cuisine" }` → planner `ChooseMealOption(household_id, plan_id,
   meal, kind)`. It costs no AI call and no quota (PLN-08-06 rule 7). Default deadline class; idempotency required.
8. AI dishes in the plan are already checked by planner's rule checker (PLN-07); the gateway does not inspect them.
9. R2, Plus: the same route calls planner `CreateFamilyDayPlan(household_id, viewer_member_id, date, replace,
   use_ai)` (PLN-13-07), which plans one dish per meal for everyone eating with per-person portions (GW-07-08).
   No quota applies. Body field `replace: true` confirms replacing an existing plan for that date; without it planner
   answers `ALREADY_EXISTS` → 409 `plan_exists` with `meta.planId`. Decided (Q125): from R2 every Plus household uses
   the family route, single-person households included (one code path; a single person is just one plate).

## Contract

`POST /v1/plans/day` — `auth: true`, tag `plans`, idempotency required, quota `dayPlan` (Free, `useAi` true).

```json
200 {
  "plan": { "id": "pln_...", "date": "2026-11-02", "status": "proposed", "meals": [
    { "meal": "breakfast", "dish": { "id": "dish_ai_...", "name": "Szakszuka z papryką", "kcalPerPortion": 420,
      "...": "" }, "source": "ai", "options": [] },
    { "meal": "supper", "dish": null, "source": null, "cuisine": "asian:japanese",
      "options": [
        { "kind": "closest_now", "dish": { "...": "" } },
        { "kind": "buy_few", "dish": { "...": "" }, "toBuyIngredientIds": ["ing_nori"] },
        { "kind": "other_cuisine", "cuisine": "polish", "dish": { "...": "" } } ] } ] },
  "aiUsed": true, "aiMealsKept": 2, "fallbackReason": null,
  "quota": { "kind": "dayPlan", "limit": 2, "used": 1, "remaining": 1, "nextFreeAt": "2026-11-05T18:42:10Z" }
}
```

`POST /v1/plans/{planId}/meals/{meal}/choice` → 200 `{ "plan": {...} }`.

RPCs: planner `ComposeDayPlan`, `ChooseMealOption` (PLN-08), `CreateFamilyDayPlan` (PLN-13-07, R2).

## Data

No storage in the gateway.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Date out of range | 400 | `validation` |
| Compose already running for the household | 409 | `plan_in_progress` |
| Plan exists for the date (family route, R2) | 409 | `plan_exists` |
| Free quota used up | 429 | `quota_reached` |
| Planner down | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] With `ai` stopped, the route returns 200 with library dishes, `aiUsed: false`,
      `fallbackReason: "ai_unavailable"`, and the quota is unchanged.
- [ ] Free: the third successful AI day plan of the week gets 429.
- [ ] For each launch cuisine, a plan from the fixed test fridge is in that cuisine or returns
      `options` naming the missing items (verification 4).
- [ ] With "no gluten, no nuts", 50 composed plans contain no violation (verification 3, planner rule, end to end).
- [ ] Choosing an option does not change `quota.used`.
- [ ] Replaying a compose with the same idempotency key returns the stored plan and does not change `quota.used`
      (planner `replayed`, `quota_consumed = false`).
- [ ] `useAi: false` from a Free household at its limit returns a library plan with `fallbackReason:
      "ai_not_allowed"` and no quota change.

## Notes

- Choosing a cuisine option never counts as a day plan: settled by planner (PLN-08-06 rule 7).
- planner's two compose RPCs treat an existing plan differently (ComposeDayPlan supersedes the earlier proposal;
  CreateFamilyDayPlan needs `replace`); the gateway follows each as written.
