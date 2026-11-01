# ADM-04-05 Dietitian sign-off

| | |
|---|---|
| Service | admin |
| Epic | [ADM-04 Recipe import and dish review](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | ADM-02-02 (role `dietitian`), ADM-04-03; PLN-03-08 (`RecordDietitianCheck`, `GetDietitianReport`), PLN-03-07 (`ListDishes`) |
| Unblocks | ADM-04-04 |

## Goal

A dietitian checks the calorie values of the 200 starter dishes before launch (features.md §5 open points), and the
result is recorded per dish, so verification 16 ("within 10% of a dietitian's calculation") can be shown.

## Business rules

1. The queue is `ListDishes(flags = [dietitian_pending])` (PLN-03-07 rule 10): dishes tagged `starter` (the first
   200 published dishes, or tagged by an admin) whose current version has no check. A new published version clears
   the old check (PLN-03-08 rule 5), so the dish comes back.
2. The dietitian sees the dish (read-only view of ADM-04-03) and enters their own kcal per portion, optionally
   protein in grams, and a note (≤ 500 characters). Admin calls `RecordDietitianCheck(dish_id, version,
   kcal_per_portion, protein_g, note)`; one check per version, which can be updated.
3. Verdict is computed by planner, not chosen: `ok` when planner's kcal is within 10% of the dietitian's, otherwise
   `deviation` (PLN-03-08 rule 2, features.md Verification 16). The dietitian cannot override it.
4. A `deviation` doesn't change the dish's status. It is listed in the report (rule 5) and the admin fixes the
   recipe through ADM-04-03; the new version needs a new check.
5. The report page is `GetDietitianReport`: starter dishes total, checked, ok, deviation (with both numbers) and
   whether the launch gate passes (PLN-03-08 rule 4, Q094: every starter dish checked and at least 95% within 10%).
6. Roles `dietitian` and `admin` may record checks and read the report.
7. Audited (`planner.dietitian_check.record`, dish id, version, dietitian kcal, verdict).

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/planner/dietitian-queue?page_token=` | `{ items: [{ dish_id, version, name, kcal_per_portion }], next_page_token }` |
| POST | `/admin/v1/planner/dishes/{id}/dietitian-check` | `{ version, kcal_per_portion, protein_g?, note? }` → `{ verdict, difference_pct }` |
| GET | `/admin/v1/planner/dietitian-report` | `{ total, checked, ok, deviation: [...], gate_passed }` |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| kcal outside 1–5,000 | 400 |
| Unknown dish or version | 404 |

## Acceptance criteria

- [ ] Computed 520 kcal, dietitian 560 → `ok` (7.1%); dietitian 600 → `deviation` (13.3%), deviation measured
      against the dietitian's value as PLN-03-08 rule 2 defines it.
- [ ] Publishing a new version of a checked dish puts it back in the queue.
