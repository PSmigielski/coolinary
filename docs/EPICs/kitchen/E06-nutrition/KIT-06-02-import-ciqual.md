# KIT-06-02 Import Ciqual

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-06 Nutrition](README.md) |
| Release | R1 (GATE: licence and credit wording for a paid app) |
| Size | M |
| Status | todo |
| Depends on | KIT-06-01 |
| Unblocks | KIT-06-05 |

## Goal

Load the French Ciqual table (ANSES; 3,484 foods, 74 nutrients, update of November 2025), the first source for plain
ingredients (features.md §5 table).

## Business rules

1. Command `server import-nutrition --source ciqual --dir <path> --version <label> [--dry-run]`, run by hand when a
   new release appears. The files (Ciqual's XML release: foods, composition, constituents) are downloaded by the
   operator, not by the service.
2. Mapping by Ciqual constituent code (confirm the codes against the constituent file of the release imported):
   energy EU kcal (328) and kJ (327), protein N × Jones factor (25000), carbohydrate (31000), sugars (32000),
   fat (40000), saturated fatty acids (40302), fibre (34100), salt (10004). Name: French and English food names.
3. Value notation: `-` → `NULL`; `traces` and `< x` per KIT-06-01 rule 3; decimal commas accepted.
4. `state` from the English name: contains "raw" → `raw`; "dried" or "dry" → `dry`; "cooked", "boiled", "roasted",
   "fried", "baked" → `cooked`; "prepared" or "ready-to-eat" → `prepared`; else `other`.
5. Upsert by (`ciqual`, food code). Foods missing from the new release are marked `retired`, not deleted (approved
   links keep working and show in admin as "source retired").
6. If values of an entry with an approved link change, that ingredient's `nutrition_version` is raised (KIT-06-10).
7. The whole import is one transaction; a parse error aborts it with the file and record.
8. `nutrition_sources` row for `ciqual` is updated with version, licence and credit text.

## Contract

No RPC (CLI).

## Data

No schema change.

## Errors

CLI exits non-zero with `<file>:<record>: <reason>`.

## Acceptance criteria

- [ ] Importing the release creates about 3,484 entries; importing it again changes nothing.
- [ ] A sample of 10 foods matches the published values exactly.
- [ ] A value change on a linked entry raises the ingredient's `nutrition_version` and produces one
      `catalogue.nutrition_changed`.
