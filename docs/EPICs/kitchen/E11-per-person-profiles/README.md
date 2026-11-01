# KIT-11 Per-person profiles and calorie targets

| | |
|---|---|
| Service | kitchen |
| Release | R2 |
| Depends on | KIT-10 (profile), KIT-01 (consumer, plan from identity), IDN-04/06/07 (`household.member_added`, `household.member_removed`, `household.snapshot`) |

## Goal

Each person who eats at home — children and people without a phone included — can have their own diet, allergens,
dislikes and a daily kcal target, so one shared dish is portioned for everyone (features.md §1 "People and goals",
§2 "Planning for each person"). Age, height, weight and goals stay on the phone; the server only ever stores the
daily kcal target (features.md "Privacy").

## Flow

```
identity ─household.member_added / member_removed / snapshot─► member rows (one per identity member)
phone: computes the target from age, height, weight, goal (on device) ─► UpdateMemberProfile(daily_kcal_target,
       target_kind, goal_set_by_parent, calories_hidden, diet, health_consent)
gateway ─► GetMemberProfiles (filtered: never another adult's target)
planner ─► GetPlanningProfile.members (targets, kinds, restrictions) ─► portions (PLN-12)
```

## What the server can and cannot check

The server never knows a person's age. It enforces: the adult floor of 1,200 kcal on every route (Verification 19);
absolute bounds for children; who may change what; Free/Plus. The phone enforces the age-based child rules (EFSA
value at the lowest activity level, the per-age lower limit for a doctor's number, no "lose" goal under 18), because
only the phone has the age.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-11-01](KIT-11-01-member-rows.md) | Member rows from identity events | R2 | M | KIT-01, KIT-10-01 |
| [KIT-11-02](KIT-11-02-edit-rights.md) | UpdateMemberProfile and who may change what | R2 | M | KIT-11-01 |
| [KIT-11-03](KIT-11-03-daily-kcal-target.md) | Daily kcal target: limits, consent, deletion | R2 | M | KIT-11-02 |
| [KIT-11-04](KIT-11-04-per-person-restrictions.md) | Per-person diet, allergens and dislikes | R2 | M | KIT-11-02, KIT-10-04 |
| [KIT-11-05](KIT-11-05-display-and-family-flags.md) | Calories hidden and "eat lighter as a family" | R2 | S | KIT-11-02 |
| [KIT-11-06](KIT-11-06-member-visibility-and-snapshot.md) | GetMemberProfiles visibility and members in the planning snapshot | R2 | M | KIT-11-03 … 05, KIT-10-08 |

## Decisions (taken 1 November 2026)

- Q192 (KIT-11-02): Only the household owner; any adult for profiles without an account.
- Q193 (KIT-11-03): 5,000 kcal.
- Q194 (KIT-11-03): 500–4,000 kcal on the server; age-based limits checked on the phone.
- Q195 (KIT-11-04): On a downgrade every person's allergens are added to the household allergen filter; per-person diets stop applying; editing per-person data waits for Plus. (changed from the first draft)
