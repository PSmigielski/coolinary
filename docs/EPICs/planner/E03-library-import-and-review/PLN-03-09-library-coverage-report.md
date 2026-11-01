# PLN-03-09 Library coverage report

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-02-03, PLN-05-02 |
| Unblocks | ADM-04 (what to draft next) |

## Goal

Show where the library is thin, so the founder knows what to import or draft before launch: the app must give at least
3 matches per meal with AI off (Verification 1) for every launch cuisine, tool preset and common diet.

## Business rules

1. Counts of published library dishes per cell of: cuisine (the 5 launch cuisines: polish, german, italian, asian,
   mediterranean) × meal (breakfast, dinner, supper) × tool preset × diet.
2. Tool presets (from KIT-10): `student_room` = {microwave, kettle}; `basic_kitchen` = {hob, oven, kettle, toaster, hand_blender, baking_tins} (Q060);
   `fully_equipped` = all tools. A dish counts for a preset when its tools ⊆ the preset. Decided (Q060): kitchen
   (KIT-10) owns the preset lists and the report reads them from kitchen; "basic kitchen" = hob, oven, kettle, toaster,
   hand blender, baking tins.
3. Diets reported: none, vegetarian, vegan, gluten_free, lactose_free, nut_free (using `diets_capable`).
4. Targets: overall 200 dishes; every cuisine × meal × preset `basic_kitchen` cell ≥ 8; every cuisine × meal for
   `student_room` ≥ 4; every cuisine × meal × diet (vegetarian, gluten_free) ≥ 3. Decided (Q100): these targets
   (they leave room for ≥ 3 matches after stock filtering, Verification 1).
5. Cells below target are listed with the gap, ready to turn into briefs (PLN-03-04).
6. A second section simulates Verification 1 on fixed test fridges stored in `planner/testdata/fridges/*.yaml`
   (10 items each, one per cuisine): for each, the number of library matches per meal with share on hand ≥ 0.5.

## Contract

```proto
// PlannerAdminService
rpc GetLibraryCoverage(GetLibraryCoverageRequest) returns (GetLibraryCoverageResponse);
```

## Data

No schema change.

## Errors

| Case | gRPC status / behaviour |
|---|---|
| kitchen down while reading the tool presets (rule 2) | report still returned with `presets_unavailable = true` and preset cells empty |
| A test fridge file missing or invalid | the simulation section lists that fridge as `error` with the reason; the rest of the report is returned |
| Empty library | success; every cell shows its full gap |

## Acceptance criteria

- [ ] With the 5 example dishes only, total = 5 and every cell is listed as below target.
- [ ] The fridge simulation reports per test fridge and meal the number of matches; the launch check passes when every
      value is ≥ 3.
