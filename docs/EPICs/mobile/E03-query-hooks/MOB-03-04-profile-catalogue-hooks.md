# MOB-03-04 Kitchen profile and catalogue hooks

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R1 (household profile, catalogue), R2 (per-person profiles and targets) |
| Size | M |
| Depends on | MOB-03-01, MOB-03-02, MOB-04-01, GW-05 |
| Unblocks | later onboarding and settings screens, the later on-device target calculator |

## Goal

Hooks for the kitchen profile (features.md §1) and ingredient search, with a hard guarantee that only a daily kcal
target, never age, height, weight or a goal, leaves the phone.

## Business rules

1. Hooks and routes (GW-05 route catalogue in docs/EPICs/gateway/README.md):

   | Hook | HTTP | Rel. | Offline |
   |---|---|---|---|
   | `useProfile()` | `GET /v1/profile` | R1 | read from cache |
   | `useUpdateProfile()` | `PATCH /v1/profile` with only the changed sections (cuisines per meal, tools, diet, staples, household settings), `If-Match` profile version | R1 | queued (MOB-04) |
   | `useUpdateMember()` | `PUT /v1/profile/members/{memberId}` | R2 | queued |
   | `useProfileOptions()` | `GET /v1/profile/options?country=` (cuisines incl. Asian sub-cuisines, tools, tool presets, diet flags, allergens, meals) | R1 | cached 24 h |
   | `useNutritionSources()` | `GET /v1/catalogue/nutrition-sources` (methodology page credits) | R1 | cached 24 h |
   | `useIngredientSearch(q)` | `GET /v1/catalogue/ingredients?q=&limit=10` | R1 | not offline |
2. Ingredient search runs only for `q.trim().length >= 2`, debounced 250 ms, keeps the previous results while
   loading (`placeholderData: keepPreviousData`), `staleTime` 1 h, not persisted.
3. **Privacy guard for member updates.** The request body is built from an explicit allow-list:

   ```ts
   interface MemberProfileUpdate {
     diet?: string[];             // e.g. "vegetarian", "gluten_free"
     allergens?: string[];
     dislikes?: string[];         // ingredient ids
     dailyKcalTarget?: number | null;   // the only number derived from health data
     caloriesHidden?: boolean;
   }
   ```
   A runtime guard throws in dev and drops the field in release if any key outside this list is present (in
   particular `age`, `birthYear`, `sex`, `height`, `weight`, `activity`, `goal`). A unit test covers each forbidden
   key.
4. `dailyKcalTarget` comes from the later on-device calculator (frontend task). The hook only sends it; the server
   refuses adult targets under 1,200 kcal (KIT-11) and the hook surfaces that as `validation` with code
   `target_too_low`.
5. Free households may set one personal target only (README "Free and Plus"); `useUpdateMember()` for a second member
   is gated by `usePlusGate("member_goals")`.
6. Profile updates are optimistic: the cache is updated at once and rolled back if the server rejects the change
   (`validation`), not when the call is merely queued offline.
7. Only the changed sections are sent with `PATCH`, with `If-Match` set to the profile version the user saw
   (GW-05-01, kitchen `UpdateProfile` with optional sections and a version). A `409 version_conflict` refetches the
   profile and re-applies the user's change to the fresh copy once; the response replaces the cache.
8. Decided (Q076): the household settings section carries `contributeToOpenFoodFacts` (default `true`, GW-05-01);
   `useUpdateProfile()` sends it like any other household setting.

## Contract

As in the table; types in `src/api/types/profile.ts`.

## Data

Query cache: profile persisted without member targets (MOB-03-01 rule 5).

## Errors

| Case | Behaviour |
|---|---|
| Target under 1,200 kcal | `validation` / `target_too_low` |
| Second member's goal on Free | `plusRequired` |

## Acceptance criteria

- [ ] Verification 20 from the app side: a network capture while saving a member shows only the allow-listed keys;
      age, height and weight never appear.
- [ ] Typing "pap" shows ingredient suggestions after 250 ms with one request.
- [ ] A profile change made offline appears at once and reaches the server after reconnect.
- [ ] Turning `contributeToOpenFoodFacts` off sends only the `household` section and survives a refetch.
