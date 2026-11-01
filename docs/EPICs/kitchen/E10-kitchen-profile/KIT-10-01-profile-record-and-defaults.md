# KIT-10-01 Profile record, country defaults, GetProfile and UpdateProfile

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-01 (request metadata, outbox) |
| Unblocks | KIT-07-02 (timezone), KIT-09, KIT-10-02 … 09, GW-05 (`GET`/`PATCH /v1/profile`) |

## Goal

Every household has one kitchen profile with its country, timezone and language, filled with country defaults until
the user saves, and edited section by section without two phones overwriting each other.

## Business rules

1. One profile per household. It is created by KIT-13-01 on `household.created`, or on first read if that event has
   not arrived yet. Until the first successful `UpdateProfile`, `is_default = true`.
2. **Country defaults.** `household.created` carries no country (IDN-01-03), so a new profile has an empty country
   until the app saves one in onboarding. While `is_default = true` and the country is empty, `GetProfile` uses the
   request's `country_hint` (the device region, sent by the app through the gateway) to pick the defaults below:

   | Country | Language | Timezone | Currency | Home cuisine |
   |---|---|---|---|---|
   | PL | pl | Europe/Warsaw | PLN | Polish |
   | DE | de | Europe/Berlin | EUR | German |
   | AT | de | Europe/Vienna | EUR | German |
   | other / unknown | en | Europe/Brussels | EUR | none ("no preference") |

   Decided (Q059): Austria defaults like Germany; every other country gets English, EUR, Europe/Brussels and no
   home cuisine. Onboarding is never refused for a country.
3. The household section holds `country`, `timezone` (IANA name, validated against the Go tz database),
   `language` (`pl`, `de`, `en`) and the cooking settings of KIT-10-06. `currency` is derived from the country and
   read-only. The app sends the device timezone on its first save; a later change of the device timezone does not
   change the profile. The household section also holds `contribute_to_off` (Decided, Q076): `true` by default;
   `false` keeps this household's confirmed product namings out of what is sent to Open Food Facts (KIT-05-04).
   It is not part of the planning snapshot and does not bump `planning_version`.
4. **Sections:** `cuisines` (KIT-10-02), `tools` (KIT-10-03), `diet` (KIT-10-04), `staples` (KIT-10-05),
   `household` (KIT-10-06 + rule 3). `UpdateProfile` takes any subset; each section sent replaces that whole
   section; sections not sent stay.
5. **Concurrency:** the profile has a `version` (starts at 1). `UpdateProfile` must send it; a mismatch fails with
   `ABORTED` and the current version in the details (the gateway turns it into 409, GW-05-01). Each successful update
   adds 1.
6. **Who may edit:** any account member of the household, anonymous accounts included. Profiles without an account
   never call. Restrictions for children's accounts come in R2 (KIT-11-02).
7. A successful update writes one `profile.changed` (KIT-10-09) in the same transaction, except an update that
   changes only `contribute_to_off` (rule 3), which planning doesn't see.
8. `GetProfile` returns the whole profile with `version` and `is_default`; `member_id` is the caller (used from R2
   to say which sections they may edit).

## Contract

`contracts/coolinary/kitchen/v1/profile.proto`

```proto
rpc GetProfile(GetProfileRequest) returns (GetProfileResponse);
rpc UpdateProfile(UpdateProfileRequest) returns (GetProfileResponse);

message HouseholdSettings {
  string country = 1;                 // ISO 3166-1 alpha-2
  string timezone = 2;                // IANA
  string language = 3;                // pl | de | en
  string currency = 4;                // read-only, ISO 4217
  // Cooking settings, KIT-10-06:
  int32 portions = 5;
  MainMeal main_meal = 6;
  MinutesPerMeal minutes_per_meal = 7;
  CookingRhythm cooking_rhythm = 8;
  MealShares meal_shares = 9;
  bool cook_non_main_meals = 10;
  bool eat_lighter_as_family = 11;    // KIT-11-05 (R2)
  bool contribute_to_off = 12;        // Q076, default true (KIT-05-04)
}

message KitchenProfile {
  int64 version = 1;
  bool is_default = 2;
  CuisineSection cuisines = 3;        // KIT-10-02
  repeated coolinary.common.v1.KitchenTool tools = 4;   // KIT-10-03
  DietSection diet = 5;               // KIT-10-04
  repeated string staple_ingredient_ids = 6;            // KIT-10-05
  HouseholdSettings household = 7;
}

message GetProfileRequest {
  string household_id = 1;
  string member_id = 2;
  string country_hint = 3;            // ISO 3166-1 alpha-2; used only while is_default and no country stored
}
message GetProfileResponse { KitchenProfile profile = 1; }

message ToolSection { repeated coolinary.common.v1.KitchenTool tools = 1; }
message StapleSection { repeated string ingredient_ids = 1; }

message UpdateProfileRequest {
  string household_id = 1;
  string member_id = 2;
  int64 version = 3;
  optional CuisineSection cuisines = 4;
  optional ToolSection tools = 5;
  optional DietSection diet = 6;
  optional StapleSection staples = 7;
  optional HouseholdSettings household = 8;
}
```

## Data

`kitchen/migrations/000NN_kitchen_profiles.sql`

```sql
-- +goose Up
CREATE TABLE kitchen_profiles (
    household_id          uuid PRIMARY KEY,
    country               text NOT NULL DEFAULT '' CHECK (country = '' OR country ~ '^[A-Z]{2}$'),
    timezone              text NOT NULL,
    language              text NOT NULL CHECK (language IN ('pl', 'de', 'en')),
    is_default            boolean NOT NULL DEFAULT true,
    contribute_to_off     boolean NOT NULL DEFAULT true,   -- Q076, KIT-05-04
    version               bigint NOT NULL DEFAULT 1,       -- ETag of GetProfile / UpdateProfile
    planning_version      bigint NOT NULL DEFAULT 1,       -- KIT-10-09: any change that planning sees
    created_at            timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now()
);
-- Cooking settings columns are added by KIT-10-06; section tables by KIT-10-02 … 05.

-- +goose Down
DROP TABLE kitchen_profiles;
```

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Version mismatch | `ABORTED` | `version_conflict` (details: current version) |
| Unknown timezone, language, country format | `INVALID_ARGUMENT` | `invalid_value` |
| Caller not an account member of the household | `PERMISSION_DENIED` | |

## Acceptance criteria

- [ ] A new household read with `country_hint = PL` gets `is_default = true`, language `pl`, timezone
      `Europe/Warsaw`, currency `PLN`, Polish pre-selected; saving `household.country = DE` later switches the
      currency to EUR.
- [ ] Saving only `tools` keeps every other section unchanged and sets `is_default = false`.
- [ ] Two phones save with version 4; the second gets `ABORTED` with current version 5.
- [ ] `timezone = "Mars/Olympus"` is rejected.
- [ ] Every successful update writes exactly one `profile.changed` outbox row, except one changing only
      `contribute_to_off`.
- [ ] A new profile has `contribute_to_off = true`; saving the household section with `false` keeps it `false`
      and does not change `planning_version`.
