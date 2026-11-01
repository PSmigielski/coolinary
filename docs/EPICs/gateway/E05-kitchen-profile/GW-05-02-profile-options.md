# GW-05-02 `GET /v1/profile/options` (presets and reference lists)

| | |
|---|---|
| Service | gateway |
| Epic | [GW-05 Kitchen profile and catalogue routes](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-01-01, KIT-10 (`GetProfileOptions`), KIT-06 (`ListNutritionSources`) |
| Unblocks | MOB-03 (onboarding screens) |

## Goal

The app gets every list it needs for the profile screens, in the user's language, in one cacheable call: cuisines
(with Asian sub-cuisines), tools and tool presets, diet flags, allergens, meals, cooking rhythms.

## Business rules

1. Public route (no JWT), so onboarding can render before sign-in finishes. Language from `x-language`; country
   optional query `country` (ISO 3166-1 alpha-2) for country defaults.
2. Calls kitchen `GetProfileOptions(language, country)`.
3. Cached in `redis-cache` per `language:country` for 1 hour; response header `Cache-Control: public,
   max-age=3600` and an `ETag` (hash of the body); `If-None-Match` → 304.
4. Tool presets from features.md §1: `student_room` (microwave + kettle), `basic_kitchen`, `fully_equipped`; their
   tool lists come from kitchen.
5. `GET /v1/catalogue/nutrition-sources` (public) → kitchen `ListNutritionSources`: each nutrition table with its
   name, version, licence, credit text and URL, for the app's methodology page (KIT-06-01 rule 8; app store reviews
   of health accuracy claims). Cached like rule 3.

## Contract

`GET /v1/profile/options?country=PL` — public, tag `kitchen-profile`. `GET /v1/catalogue/nutrition-sources` —
public, tag `catalogue`, 200 `{ "sources": [ { "code", "name", "version", "licence", "credit", "url" } ] }`.

```json
200 {
  "cuisines": [{ "code": "asian", "name": "Azjatycka", "sub": [{ "code": "asian:thai", "name": "Tajska" }] }],
  "tools": [{ "code": "air_fryer", "name": "Frytkownica beztłuszczowa" }],
  "toolPresets": [{ "code": "student_room", "name": "Pokój w akademiku", "tools": ["microwave", "kettle"] }],
  "dietFlags": [{ "code": "vegan", "name": "Wegańska" }],
  "allergens": [{ "code": "gluten", "name": "Gluten" }],
  "meals": ["breakfast", "dinner", "supper"],
  "cookingRhythms": ["daily", "every_other_day", "twice_week"],
  "defaults": { "homeCuisine": "polish", "mainMeal": "dinner" }
}
```

## Data (redis-cache)

| Key | Value | TTL |
|---|---|---|
| `gw:opts:<language>:<country>` | response JSON | 1 h |

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bad `country` | 400 | `validation` |
| Kitchen down, nothing cached | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] `Accept-Language: de` returns German names; `pl` Polish.
- [ ] A second identical request within an hour does not call kitchen.
- [ ] `If-None-Match` with the current ETag returns 304.
- [ ] `student_room` preset lists exactly `microwave` and `kettle`.
