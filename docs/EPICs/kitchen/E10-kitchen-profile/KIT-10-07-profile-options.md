# KIT-10-07 Profile options and onboarding defaults

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-10-02 … KIT-10-06 |
| Unblocks | GW-05 (`GET /v1/profile/options`) |

## Goal

The app gets every choice of the profile screens — codes, names in the user's language, presets and the country's
defaults — in one cacheable call, so onboarding fits in four screens and works before the user has saved anything.

## Business rules

1. `GetProfileOptions(language, country)` needs no household; it is public data (the gateway serves it without a
   JWT, GW-05-02).
2. It returns, with names in `language` (`pl`, `de`, `en`; anything else → `en`):
   1. cuisines with Asian's six sub-cuisines nested under it;
   2. the 17 kitchen tools;
   3. the three tool presets with their tool lists (KIT-10-03 rule 3);
   4. the five diet flags and the 14 allergens;
   5. meals (breakfast, dinner, supper), main-meal choices, cooking rhythms, minutes bounds, meal-share bounds and
      defaults;
   6. the country's defaults: home cuisine, timezone, language, currency, main meal, default preset, the default
      staple ingredient ids.
3. Names are fixed strings kept in kitchen (one translation table in code, reviewed by native speakers per language
   wave, business model §6). No AI.
4. The response is the same for every caller of one (language, country) pair; the gateway caches it 1 hour.
5. **Onboarding:** the app shows four screens (cuisines, tools, diet and allergens, household) pre-filled from
   these defaults; each screen can be skipped. The app saves once at the end with `UpdateProfile` (all sections),
   which sets `is_default = false`. People and goals are not part of onboarding.

## Contract

`contracts/coolinary/kitchen/v1/profile.proto`

```proto
rpc GetProfileOptions(GetProfileOptionsRequest) returns (GetProfileOptionsResponse);

message GetProfileOptionsRequest { string language = 1; string country = 2; }

message NamedCode { int32 code = 1; string key = 2; string name = 3; repeated NamedCode sub = 4; }

message ToolPreset { string key = 1; string name = 2; repeated coolinary.common.v1.KitchenTool tools = 3; }

message GetProfileOptionsResponse {
  repeated NamedCode cuisines = 1;
  repeated NamedCode tools = 2;
  repeated ToolPreset tool_presets = 3;
  repeated NamedCode diets = 4;
  repeated NamedCode allergens = 5;
  repeated NamedCode cooking_rhythms = 6;
  MealShares default_meal_shares = 7;
  MinutesPerMeal default_minutes = 8;
  KitchenProfile country_defaults = 9;   // what GetProfile returns for a new household of this country
}
```

## Data

No schema change.

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| Unknown language | — | English names |
| Unknown or empty country | — | "other" defaults (KIT-10-01 rule 2) |

## Acceptance criteria

- [ ] `(pl, PL)` returns Polish names, home cuisine Polish, the three presets, and Asian with six sub-cuisines.
- [ ] `(de, DE)` returns German names and home cuisine German.
- [ ] `country_defaults` equals what `GetProfile` returns for a brand-new household of that country.
- [ ] Every enum value has a name in pl, de and en (test fails on a missing translation).
