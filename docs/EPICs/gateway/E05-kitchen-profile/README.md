# GW-05 Kitchen profile and catalogue routes

Releases R1 (household profile, options, ingredient search) and R2 (per-person profiles and calorie targets).

## Goal

The app sets up and edits the kitchen profile (features.md §1): cuisines per meal, tools, diet and allergens,
dislikes, staples, household settings, and in R2 each person's diet and daily kcal target. Kitchen owns all of it
(decision 3); the gateway maps routes to kitchen RPCs and enforces the privacy shape of the bodies.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-05-01](GW-05-01-household-profile.md) | `GET` and `PATCH /v1/profile` | R1 | M | GW-01-02, GW-01-03, KIT-10 |
| [GW-05-02](GW-05-02-profile-options.md) | `GET /v1/profile/options` (presets and reference lists) | R1 | S | GW-01-01, KIT-10 |
| [GW-05-03](GW-05-03-ingredient-search.md) | `GET /v1/catalogue/ingredients` (autocomplete) | R1 | S | GW-01-01, KIT-02 |
| [GW-05-04](GW-05-04-member-profiles.md) | Per-person profiles and daily kcal targets | R2 | M | GW-05-01, GW-04-02, KIT-11 |
