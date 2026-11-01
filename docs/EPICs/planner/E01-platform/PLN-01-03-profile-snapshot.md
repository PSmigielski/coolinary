# PLN-01-03 Planning-profile snapshot cache

| | |
|---|---|
| Service | planner |
| Epic | [PLN-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-01-01, PLN-01-02, KIT-10 (GetPlanningProfile, profile.changed) |
| Unblocks | PLN-05-02, PLN-06, PLN-07, PLN-08, PLN-12 |

## Goal

Planner needs the kitchen profile (cuisines, tools, diets, allergens, dislikes, staples, household settings, people
and targets) on almost every request. One cached snapshot per household keeps that to one kitchen call per change.

## Business rules

1. `profiles.Get(household_id)` returns the snapshot from `KIT-10 GetPlanningProfile`. Planner never stores the
   profile in Postgres; it owns none of it (decision 3).
2. Cache key `planner:profile:<household_id>` in redis-cache, value = the protobuf bytes plus kitchen's
   `profile_version`; TTL 10 minutes.
3. On `profile.changed` for a household, the key is deleted. Also deleted on `household.member_added`,
   `household.member_removed`, `household.plan_changed` and `household.deleted` (the snapshot lists people and Plus-only
   parts).
4. In-process single-flight: concurrent misses for the same household make one kitchen call.
5. If kitchen is `Unavailable` and redis-cache holds no entry, the request fails with UNAVAILABLE ("profile
   unavailable"); planner never guesses hard limits (tools, allergens). A stale entry is never served past its TTL.
6. The snapshot holds only what kitchen sends: daily kcal targets, never age, weight or height (features.md
   "Privacy").
7. The snapshot exposes helpers used by later epics: `Tools()`, `Staples()`, `EatersFor(meal)` (all people unless
   marked away, PLN-12), `HardLimits(memberIDs)` (union of allergens, diets and dislikes of the given people),
   `CuisineFor(meal)`, `TimeBudget(meal)`, `MealShares()`, `Timezone()`, `Language()`.
8. In R1 a household has one profile (one person, or the household as a whole); `HardLimits` then returns that profile's
   limits.

## Contract

No planner contract change; consumes `KitchenService.GetPlanningProfile`.

## Data

No schema change.

## Events

Consumes `profile.changed`, `household.member_added`, `household.member_removed`, `household.plan_changed`,
`household.deleted` (cache delete only).

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Kitchen down, no cache | UNAVAILABLE | `profile unavailable` |
| Kitchen returns NOT_FOUND (household has no profile yet) | — | Use kitchen's documented defaults for the country (KIT-10); never fail |

## Acceptance criteria

- [ ] Two dish searches in a row for one household make one `GetPlanningProfile` call.
- [ ] After a `profile.changed` event, the next request fetches a fresh snapshot (new `profile_version`).
- [ ] 50 concurrent requests on a cold cache make one kitchen call.
- [ ] With kitchen down and no cache, `SearchDishes` returns UNAVAILABLE instead of unfiltered dishes.
- [ ] Test asserts the snapshot type has no field for age, height or weight.
