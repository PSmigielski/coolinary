# KIT-10 Kitchen profile and household settings

| | |
|---|---|
| Service | kitchen |
| Release | R1 |
| Depends on | KIT-01 (outbox, plan from identity), KIT-02 (ingredients for staples and dislikes) |

## Goal

The household's kitchen profile — cuisines, tools, diet and allergens, staples, how and when it cooks — set once in at
most four screens with presets, and read by planner as one snapshot (features.md §1 "Kitchen profile").
Per-person profiles and calorie targets are KIT-11 (R2).

## Flow

```
app ─GET /v1/profile─► GetProfile (country defaults until first save, is_default = true)
app ─GET /v1/profile/options─► GetProfileOptions (codes, names in the user's language, presets, defaults)
app ─PATCH /v1/profile─► UpdateProfile(section…, version) ─► profile.changed ─► planner drops its snapshot
planner ─► GetPlanningProfile (one snapshot: settings, cuisines, tools, diet, staples, members, plan)
```

Onboarding screens (at most 4, each skippable with the defaults): 1 cuisines, 2 kitchen tools (preset), 3 diet and
allergens, 4 household (portions, main meal, cooking rhythm). People and goals are not asked at first start.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-10-01](KIT-10-01-profile-record-and-defaults.md) | Profile record, country defaults, GetProfile and UpdateProfile | R1 | M | KIT-01 |
| [KIT-10-02](KIT-10-02-cuisines-per-meal.md) | Cuisines per meal | R1 | S | KIT-10-01 |
| [KIT-10-03](KIT-10-03-kitchen-tools-and-presets.md) | Kitchen tools and presets | R1 | S | KIT-10-01 |
| [KIT-10-04](KIT-10-04-household-diet-allergens-dislikes.md) | Household diet, allergens and dislikes | R1 | M | KIT-10-01, KIT-02 |
| [KIT-10-05](KIT-10-05-staples.md) | Staples ("always at home") | R1 | S | KIT-10-01, KIT-02 |
| [KIT-10-06](KIT-10-06-household-rhythm-and-meal-shares.md) | Household rhythm and meal shares | R1 | M | KIT-10-01 |
| [KIT-10-07](KIT-10-07-profile-options.md) | Profile options and onboarding defaults | R1 | S | KIT-10-02 … 06 |
| [KIT-10-08](KIT-10-08-planning-profile-snapshot.md) | GetPlanningProfile snapshot for planner | R1 | M | KIT-10-02 … 06, KIT-01-06 (plan from identity) |
| [KIT-10-09](KIT-10-09-profile-changed.md) | profile.changed event and planning version | R1 | S | KIT-10-01, KIT-01 |

## Decisions (taken 1 November 2026)

- Q059 (KIT-10-01): Austria like Germany; others English, EUR, Brussels time, no home cuisine.
- Q060 (KIT-10-03): Hob, oven, kettle, toaster, hand blender, baking tins.
- Q061 (KIT-10-03): Yes, the kettle is a tool (the "student room" preset uses it).
- Q062 (KIT-10-06): Breakfast 10–40%, main 25–60%, other 10–40%, snacks 0–20%.
- Q063 (KIT-10-06): One default, 25/40/25/10, for every country for now.
