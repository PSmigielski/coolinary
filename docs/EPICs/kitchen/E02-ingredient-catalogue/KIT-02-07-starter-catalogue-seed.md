# KIT-02-07 Starter catalogue seed and seed command

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-02-01, KIT-02-02, KIT-02-03, KIT-02-04 |
| Unblocks | PLN-03 (recipe import needs ingredients to exist), KIT-06-05, launch |

## Goal

Load a reviewed starter catalogue from files kept in the repository, repeatably, so every environment has the same
ingredients and the 200 starter dishes (features.md §7 release 1) can reference them.

## Business rules

1. Seed files live in `kitchen/seed/catalogue/`: `ingredients.csv` (slug, category, base unit, allergens,
   vegetarian, vegan, contains lactose, staple candidate, density, default storage, shelf lives), `names.csv`
   (slug, language, name, kind), `measures.csv` (slug, measure, grams), `packs.csv` (slug, country, value,
   is default).
2. The command `server seed catalogue --dir seed/catalogue [--dry-run]` (run with `make seed s=kitchen`) upserts by
   `slug`. It never deletes; a slug missing from the files is left alone. `--dry-run` prints the changes and writes
   nothing.
3. Seeded names get `origin = 'seed'`, `reviewed = true`. A seed run never changes rows with `origin = 'admin'`
   (admin edits win over the files).
4. The whole run is one transaction; any invalid row aborts it with the file name and line number.
5. Every change is written to `admin_changes` with `admin_user_id = 'seed'`, and the catalogue version is raised
   once per run.
6. Allergens and diet flags in the seed are checked by a person before merge (safety: they are hard limits).
7. Coverage rule, checked by the command: every ingredient referenced by the planner starter dish files (PLN-03)
   exists in the seed, and each seeded ingredient has a primary `en` name.

## Contract

No RPC. CLI subcommand only.

## Data

No schema change.

## Errors

CLI exits non-zero with `<file>:<line>: <reason>` on invalid rows.

## Acceptance criteria

- [ ] Running the seed twice in a row makes no changes the second time.
- [ ] A row with an unknown allergen code aborts the run and names the line.
- [ ] An admin-edited name survives a seed run.
- [ ] Dev and prod produce the same catalogue from the same files.

## Notes

Decided (Q053): the starter content is about 700 ingredients (everything the 200 starter dishes use plus common
fridge items in Poland and Germany), drafted with AI help outside the product, then checked by the founder (names,
allergens, diet flags) and by a native speaker per language. Writing the content is
an estimated 4–6 days of non-coding work on top of this task.
- Size M covers the code only; writing the starter content itself is separate work (see the notes above).
