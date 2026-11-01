# KIT-10-06 Household rhythm and meal shares

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-10-01 |
| Unblocks | KIT-10-08, PLN-08 (day plan), PLN-12 (portions), PLN-14 (week plan) |

## Goal

How the household cooks: how many portions, which meal is the main hot one, how long each meal may take, how often
to cook, and how each person's day target is shared over the meals (features.md §1 "Household", §2 "How one meal
is split", step 1).

## Business rules

1. **Portions:** 1–12, default 2. In R1 it is the household's portion count. From R2 planner sizes the pot from
   people and targets and uses this value only for eaters without a profile (PLN-12-02 rule 7).
2. **Main hot meal:** `dinner` or `supper`, default `dinner` in every country (Polish obiad at midday; the German cold
   supper). Breakfast is never the main meal.
3. **Minutes per meal:** breakfast, dinner, supper, each 5–180, defaults 15, 45 and 20.
4. **Cooking rhythm:** `daily` (default), `every_other_day`, `twice_a_week`.
5. **Cook the other meals:** `cook_non_main_meals`, default false: breakfast and the meal that is not the main one
   come from "no cooking" dishes unless this is true (features.md §2 "One dish per meal").
6. **Meal shares** of each person's day target, whole percentages: breakfast, main hot meal, the other meal, snacks
   and drinks. Default 25 / 40 / 25 / 10.
   1. The four must sum to exactly 100.
   2. Decided (Q062): bounds breakfast 10–40, main 25–60, other 10–40, snacks 0–20.
   3. Violations → `INVALID_ARGUMENT` with reason `shares_invalid` (gateway → 422).
   4. Decided (Q063): one default, 25 / 40 / 25 / 10, for every country for now; a per-country default (features.md
      §5 open point) can be added later as config.
7. These values are part of the `household` section (KIT-10-01), replaced as a whole on save.

## Contract

`contracts/coolinary/kitchen/v1/profile.proto`

```proto
enum MainMeal { MAIN_MEAL_UNSPECIFIED = 0; MAIN_MEAL_DINNER = 1; MAIN_MEAL_SUPPER = 2; }

enum CookingRhythm {
  COOKING_RHYTHM_UNSPECIFIED = 0;
  COOKING_RHYTHM_DAILY = 1;
  COOKING_RHYTHM_EVERY_OTHER_DAY = 2;
  COOKING_RHYTHM_TWICE_A_WEEK = 3;
}

message MinutesPerMeal { int32 breakfast = 1; int32 dinner = 2; int32 supper = 3; }

message MealShares { int32 breakfast = 1; int32 main = 2; int32 other = 3; int32 snacks = 4; }
```

## Data

`kitchen/migrations/000NN_profile_household_settings.sql`

```sql
-- +goose Up
ALTER TABLE kitchen_profiles
    ADD COLUMN portions            smallint NOT NULL DEFAULT 2  CHECK (portions BETWEEN 1 AND 12),
    ADD COLUMN main_meal           smallint NOT NULL DEFAULT 1  CHECK (main_meal IN (1, 2)),
    ADD COLUMN minutes_breakfast   smallint NOT NULL DEFAULT 15 CHECK (minutes_breakfast BETWEEN 5 AND 180),
    ADD COLUMN minutes_dinner      smallint NOT NULL DEFAULT 45 CHECK (minutes_dinner BETWEEN 5 AND 180),
    ADD COLUMN minutes_supper      smallint NOT NULL DEFAULT 20 CHECK (minutes_supper BETWEEN 5 AND 180),
    ADD COLUMN cooking_rhythm      smallint NOT NULL DEFAULT 1  CHECK (cooking_rhythm IN (1, 2, 3)),
    ADD COLUMN cook_non_main_meals boolean  NOT NULL DEFAULT false,
    ADD COLUMN share_breakfast     smallint NOT NULL DEFAULT 25 CHECK (share_breakfast BETWEEN 10 AND 40),
    ADD COLUMN share_main          smallint NOT NULL DEFAULT 40 CHECK (share_main BETWEEN 25 AND 60),
    ADD COLUMN share_other         smallint NOT NULL DEFAULT 25 CHECK (share_other BETWEEN 10 AND 40),
    ADD COLUMN share_snacks        smallint NOT NULL DEFAULT 10 CHECK (share_snacks BETWEEN 0 AND 20),
    ADD COLUMN eat_lighter_as_family boolean NOT NULL DEFAULT false,
    ADD CONSTRAINT shares_sum_100 CHECK (share_breakfast + share_main + share_other + share_snacks = 100);

-- +goose Down
ALTER TABLE kitchen_profiles
    DROP CONSTRAINT shares_sum_100,
    DROP COLUMN eat_lighter_as_family, DROP COLUMN share_snacks, DROP COLUMN share_other, DROP COLUMN share_main,
    DROP COLUMN share_breakfast, DROP COLUMN cook_non_main_meals, DROP COLUMN cooking_rhythm,
    DROP COLUMN minutes_supper, DROP COLUMN minutes_dinner, DROP COLUMN minutes_breakfast,
    DROP COLUMN main_meal, DROP COLUMN portions;
```

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Shares not summing to 100 or out of bounds | `INVALID_ARGUMENT` | `shares_invalid` |
| Portions, minutes or enums out of range | `INVALID_ARGUMENT` | `invalid_value` |

## Acceptance criteria

- [ ] Defaults: 2 portions, dinner main, 15/45/20 minutes, daily, 25/40/25/10.
- [ ] Shares 30/40/20/10 save; 30/40/25/10 (sum 105) fail with `shares_invalid`; 5/60/25/10 fail (breakfast < 10).
- [ ] features.md worked example ("dinner is the main meal", 40% share) reads 40 from here.
