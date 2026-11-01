# GW-05-01 `GET` and `PATCH /v1/profile`

| | |
|---|---|
| Service | gateway |
| Epic | [GW-05 Kitchen profile and catalogue routes](README.md) |
| Release | R1 |
| Size | M |
| Depends on | GW-01-02, GW-01-03, KIT-10 (`GetProfile`, `UpdateProfile`) |
| Unblocks | MOB-03 (onboarding and settings) |

## Goal

Read the household's kitchen profile in one call and change any section of it, with safe concurrent edits from
several household members.

## Business rules

1. `GET` calls kitchen `GetProfile(household_id, member_id)`. A household that never saved a profile gets kitchen's
   defaults for its country (kitchen decides; e.g. home cuisine pre-selected, features.md §1) with `isDefault: true`.
2. `PATCH` takes any subset of the sections below; sections not sent are unchanged. Inside a section the whole
   section is replaced (no deep merge), so arrays are always sent complete.
3. Concurrency: `GET` returns `ETag: "<version>"`; `PATCH` requires `If-Match`. Kitchen compares and answers
   `ABORTED` on mismatch → 409 `version_conflict` with `meta.currentVersion`. Missing `If-Match` → 428
   `precondition_required`.
4. Validation of values (allowed codes, 4-screen presets, meal share sums) is kitchen's (KIT-10). The gateway checks
   only types, enums from the OpenAPI schema and sizes: at most 30 tools, 50 staples, 50 dislikes, 12 cuisines per
   meal.
5. In R1 diet, allergens and dislikes are household-level (`household.diet`); per-person values come in R2
   (GW-05-04).
6. Successful `PATCH` makes kitchen publish `profile.changed` (KIT-10); the gateway does nothing extra.
7. Who may edit: any member with an account (anonymous included, it is their own household); kitchen may restrict
   sections (KIT-10).
8. Decided (Q076): `household.contributeToOpenFoodFacts` (kitchen `contribute_to_off`, KIT-10-01, KIT-05-04),
   default `true`. Set to `false`, the household's confirmed product names are never sent to Open Food Facts.

## Contract

`GET /v1/profile`, `PATCH /v1/profile` — `auth: true`, tag `kitchen-profile`, idempotency optional.

```json
{
  "version": 7, "isDefault": false,
  "cuisines": {
    "breakfast": ["polish"], "dinner": ["polish", "italian"], "supper": ["asian:japanese", "mediterranean"],
    "noPreference": false
  },
  "tools": ["hob", "oven", "microwave", "air_fryer"],
  "diet": { "flags": ["vegetarian"], "allergens": ["peanuts", "gluten"], "dislikeIngredientIds": ["ing_..."] },
  "staples": ["ing_salt", "ing_oil", "ing_flour"],
  "household": {
    "portions": 2, "mainMeal": "dinner",
    "minutesPerMeal": { "breakfast": 15, "dinner": 45, "supper": 20 },
    "cookingRhythm": "every_other_day",
    "mealShares": { "breakfast": 25, "main": 40, "other": 25, "snacks": 10 },
    "country": "PL", "timezone": "Europe/Warsaw", "language": "pl",
    "eatLighterAsFamily": false, "contributeToOpenFoodFacts": true
  }
}
```

`PATCH` body: any of `cuisines`, `tools`, `diet`, `staples`, `household` (each whole). 200 returns the full profile.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Unknown enum value, wrong type, too many items | 400 | `validation` |
| `If-Match` missing | 428 | `precondition_required` |
| Version mismatch | 409 | `version_conflict` |
| Kitchen rule broken (e.g. shares not 100) | 422 | kitchen's code |

## Acceptance criteria

- [ ] A new household gets a profile with `isDefault: true` and Polish cuisine pre-selected for country `PL`.
- [ ] Two members patch `tools` with the same `If-Match`: one succeeds, the other gets 409 with `currentVersion`.
- [ ] Patching only `staples` leaves `tools` unchanged.
- [ ] A `tools` array of 31 entries is refused with 400.
- [ ] A new profile has `household.contributeToOpenFoodFacts: true`; patching `household` with `false` keeps it
      `false` on the next `GET`.
