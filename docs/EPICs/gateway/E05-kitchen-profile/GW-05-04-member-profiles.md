# GW-05-04 Per-person profiles and daily kcal targets

| | |
|---|---|
| Service | gateway |
| Epic | [GW-05 Kitchen profile and catalogue routes](README.md) |
| Release | R2 |
| Size | M |
| Depends on | GW-05-01, GW-04-02, KIT-11 (`GetMemberProfiles`, `UpdateMemberProfile`) |
| Unblocks | MOB-03 (people and goals screens) |

## Goal

Each person (accounts and profiles) gets their own diet, allergens, dislikes, calorie display setting and daily
kcal target. Only the target reaches the server; the phone calculates it from age, height and weight it keeps
(features.md §2 "Planning for each person", "Privacy").

## Business rules

1. `GET /v1/profile/members` calls kitchen `GetMemberProfiles(household_id, member_id)`. Kitchen filters by
   visibility (features.md §2 "Privacy"): the requester's own target; children's (profiles') targets for adult
   members; **never another adult's target or goal** — for them `dailyKcalTarget` is omitted and `hasTarget` says
   whether one exists.
2. `PUT /v1/profile/members/{memberId}` calls kitchen `UpdateMemberProfile`. Body fields: `diet`
   (`flags`, `allergens`, `dislikeIngredientIds`), `dailyKcalTarget` (integer kcal or `null` to delete),
   `targetKind` (`adult` | `child`), `goalSetByParent` (boolean), `caloriesHidden` (boolean), `healthConsent`
   (`{ "version": "2026-11", "givenAt": RFC 3339 }`, required whenever `dailyKcalTarget` is set; recorded by kitchen).
3. **Privacy guard:** the body schema has `additionalProperties: false`; it has no field for age, birth date, sex,
   height, weight, activity or goal type. Any such field → 400. Acceptance test included (verification 20).
4. Limits and Free/Plus are kitchen's rules (KIT-11), passed through:
   - adult target below 1,200 kcal → kitchen `FAILED_PRECONDITION` `target_below_minimum` → 422 (verification 19);
   - Free household setting a target for anyone but the requester → `plus_required` → 403;
   - per-person diets on Free → `plus_required` → 403;
   - a teenager's own account changing a parent-set goal → `goal_locked` → 403.
5. Deleting a target (`dailyKcalTarget: null`) removes it on the server (verification 20, "deleting a goal removes
   the target from the server").
6. Concurrency: `If-Match` with the member profile `version`, as in GW-05-01.
7. Values are never logged (the log line names the route only).

## Contract

`GET /v1/profile/members`, `PUT /v1/profile/members/{memberId}` — `auth: true`, tag `kitchen-profile`,
idempotency optional.

```json
GET 200 { "members": [
  { "memberId": "0192...", "isMe": true, "version": 2,
    "diet": { "flags": [], "allergens": ["peanuts"], "dislikeIngredientIds": [] },
    "hasTarget": true, "dailyKcalTarget": 2000, "targetKind": "adult", "caloriesHidden": false },
  { "memberId": "0192...", "isMe": false, "version": 1,
    "diet": { "flags": ["vegetarian"], "allergens": [], "dislikeIngredientIds": [] },
    "hasTarget": true, "targetKind": "adult", "caloriesHidden": false }
] }
PUT body { "diet": {...}, "dailyKcalTarget": 1400, "targetKind": "child", "goalSetByParent": true,
           "caloriesHidden": true, "healthConsent": { "version": "2026-11", "givenAt": "2026-11-01T09:00:00Z" } }
PUT 200 one member object
```

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Health field in body, wrong type | 400 | `validation` |
| Target set without `healthConsent` | 422 | `consent_required` |
| Adult target < 1,200 | 422 | `target_below_minimum` |
| Free, other member's target or any per-person diet | 403 | `plus_required` |
| Teen changing a parent-set goal | 403 | `goal_locked` |
| Version mismatch | 409 | `version_conflict` |

## Acceptance criteria

- [ ] `PUT` with `{"weightKg": 80}` returns 400 and kitchen is not called.
- [ ] Adult B's `GET` shows adult A with `hasTarget: true` and no `dailyKcalTarget`.
- [ ] An adult target of 1,150 returns 422 `target_below_minimum` (verification 19).
- [ ] Setting `dailyKcalTarget: null` and reading again returns `hasTarget: false`.
- [ ] Network capture of the app shows only `dailyKcalTarget` leaving the phone (verification 20, with MOB-03).
