# ADM-04-01 Recipe import upload

| | |
|---|---|
| Service | admin |
| Epic | [ADM-04 Recipe import and dish review](README.md) |
| Release | R1 |
| Size | M |
| Depends on | ADM-01-05, ADM-02-03; PLN-03-01 (import file format), PLN-03-03 (`PlannerAdminService.ImportDishes`) |
| Unblocks | ADM-04-03 |

## Goal

Hand-written dishes (decision 10) go into the library through an import file: checked first, then committed into
the review queue, never published directly.

## Business rules

1. Upload is a multipart file of at most 2 MB holding at most 200 dishes, in the format PLN-03 defines.
2. Two steps:
   1. `mode=validate`: admin calls `ImportDishes(file, dry_run = true)`; planner parses and checks the file and
      returns, per dish, the outcome it would have and its problems, plus every ingredient name it could not resolve
      (`unresolved`, with the slugs that use it). Nothing is stored. admin returns the result and the file's SHA-256.
   2. `mode=commit`: the same file again, with `validated_sha256`. admin refuses unless the hash equals the one it
      validated for this admin within the last hour (kept in memory) and no dish had a blocking problem, then calls
      `ImportDishes(file, dry_run = false)`.
3. On commit planner writes each dish (PLN-03-03 rule 2) and returns per slug the outcome (`CREATED`,
   `NEW_VERSION`, `REPLACED_DRAFT`, `UNCHANGED`, `FAILED`) and the dish id; written dishes appear in the review queue
   (ADM-04-03).
4. A dish whose slug already exists follows planner's rule: a published dish gets a new version in review (the live
   one stays), an unpublished one is replaced, a retired one fails with `DISH_RETIRED`; admin only shows it.
5. Audited (`planner.import.commit`, file hash, dish count, outcome counts and dish ids). Validation is not audited.
6. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/planner/imports?mode=validate` | multipart `file` → `{ sha256, results: [{ slug, outcome, problems[] }], unresolved: [{ name, language, slugs[] }] }` |
| POST | `/admin/v1/planner/imports?mode=commit` | multipart `file`, field `validated_sha256` → `{ import_id, results: [{ slug, outcome, dish_id }] }` |

Calls `PlannerAdminService.ImportDishes(file, dry_run)` (PLN-03-03) with a 30 s deadline.

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| File over 2 MB or over 200 dishes | 413 / 400 |
| Commit without a matching, error-free validation in the last hour | 409 `validate_first` |

## Acceptance criteria

- [ ] A file with one unknown ingredient name validates with that name in `unresolved`; commit is refused.
- [ ] After fixing the file, validate then commit creates the dishes in `review`.
- [ ] Commit with a different file than the validated one → 409.
