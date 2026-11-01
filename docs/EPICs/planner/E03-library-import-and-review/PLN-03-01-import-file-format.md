# PLN-03-01 Import file format and parser

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-02-02 |
| Unblocks | PLN-03-03 |

## Goal

A YAML format the founder (or a hired recipe writer) can type dishes in, with a parser that turns it into dish
versions and reports every mistake with its line number.

## Business rules

1. File = UTF-8 YAML, ≤ 2 MB, ≤ 500 dishes, top-level `format: coolinary-dishes/1` and `dishes: [...]`.
2. Each dish:

   ```yaml
   - slug: chicken-curry
     cuisine: asian
     sub_cuisine: indian
     meals: [dinner]
     tools: [hob]
     total_minutes: 40
     active_minutes: 20
     servings: 4
     effort_tags: [keeps_3_days, freezes_well]
     ingredients:
       - { name: chicken thigh, amount: 600, unit: g }
       - { name: basmati rice, amount: 300, unit: g }
       - { name: coriander, amount: 1, unit: tbsp, optional: true, note: { en: chopped } }
       - { id: ing_0193, amount: 1, unit: tsp }            # catalogue id instead of a name
     steps:
       - { tool: hob }
       - { tool: hob, timer: 20m }
     texts:
       en: { name: Chicken curry with rice, intro: "", steps: ["Fry the onion ...", "Simmer ..."] }
       pl: { name: Kurczak curry z ryżem, steps: ["...", "..."] }
     source_language: en
   ```
3. Units accepted: `g`, `kg`, `ml`, `l`, `pc`, `tbsp`, `tsp`. `kg` and `l` are stored as `g` and `ml` (× 1,000).
   Timers: seconds as an integer, or `90s`, `20m`, `1h30m`.
4. Ingredient by `name` is resolved in `source_language` (default `en`) through PLN-03-02; by `id` it is used as is.
   Each line has exactly one of `name` or `id`.
5. `role` defaults to `main`; `topping` requires `group`.
6. Every field name is checked; an unknown key is an error (catches typos such as `ingrediants`).
7. The parser output is passed to the validator (PLN-02-02). Problems are reported with `slug`, YAML line, path and
   code; parse errors stop that dish but not the file.
8. The five skeleton seed dishes (owsianka, shakshuka, chicken-curry, pierogi, salad-bowl) are kept as
   `planner/testdata/import/example.yaml`, the format's worked example.

## Contract

Internal parser package; used by `ImportDishes` (PLN-03-03).

## Data

No schema change.

## Errors

None at the RPC level; parse and validation problems are returned in the `ImportDishes` report (PLN-03-03) with slug, YAML line, path and code (rule 7). A file that is not valid YAML at all fails the whole import with `INVALID_ARGUMENT` in PLN-03-03.

## Acceptance criteria

- [ ] `example.yaml` parses into 5 dishes with no problems once the catalogue holds their ingredients.
- [ ] A misspelled key reports `UNKNOWN_FIELD` with the right line number.
- [ ] `timer: 1h30m` becomes 5,400 seconds; `amount: 1.5, unit: kg` becomes 1,500 g.
- [ ] A line with both `name` and `id` is rejected.
