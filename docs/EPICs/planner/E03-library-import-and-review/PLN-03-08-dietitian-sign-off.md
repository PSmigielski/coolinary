# PLN-03-08 Dietitian sign-off

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-03-07, PLN-04-01 |
| Unblocks | Launch gate (features.md §5 "Open points", Verification 16) |

## Goal

Before launch a dietitian checks the calories of the 200 starter dishes (features.md §5 "Open points"). Planner records
their own calculation per dish version and reports where Coolinary's number is more than 10% off.

## Business rules

1. `RecordDietitianCheck(dish_id, version, kcal_per_portion, protein_g?, note)` by an admin acting for the dietitian;
   one row per version, updatable.
2. Deviation = |planner kcal − dietitian kcal| / dietitian kcal. ≤ 10% → `ok`; above → `deviation` (Verification 16:
   within 10%).
3. A dish flagged `starter` (set on the first 200 published dishes, or by admin) shows `dietitian_pending` in the
   queue until checked. Decided (Q093): checking is not required to publish, but launch waits until every starter
   dish is checked (rule 4).
4. `GetDietitianReport` returns: starter dishes total, checked, ok, deviation (with both numbers), and whether the
   launch gate passes: all starter dishes checked and ≥ 95% within 10%. Decided (Q094): 95%, with every starter
   dish checked (features.md's 20-dish check, Verification 16, is part of it).
5. A new published version of a checked dish clears the check (the numbers may have changed).

## Contract

```proto
// PlannerAdminService
rpc RecordDietitianCheck(RecordDietitianCheckRequest) returns (RecordDietitianCheckResponse);
rpc GetDietitianReport(GetDietitianReportRequest) returns (GetDietitianReportResponse);
```

## Data

`planner/migrations/000NN_dietitian_checks.sql`

```sql
ALTER TABLE dishes ADD COLUMN starter boolean NOT NULL DEFAULT false;
CREATE TABLE dietitian_checks (
    dish_id          uuid NOT NULL,
    version          int  NOT NULL,
    kcal_per_portion int  NOT NULL CHECK (kcal_per_portion > 0),
    protein_g        numeric(6,1),
    note             text NOT NULL DEFAULT '',
    recorded_by      text NOT NULL,
    recorded_at      timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (dish_id, version),
    FOREIGN KEY (dish_id, version) REFERENCES dish_versions (dish_id, version) ON DELETE CASCADE
);
```

## Errors

| Case | gRPC status / behaviour |
|---|---|
| `kcal_per_portion` ≤ 0 or above 5,000, or `protein_g` < 0 | `INVALID_ARGUMENT` |
| Unknown dish or version | `NOT_FOUND` |
| Version is not published | `FAILED_PRECONDITION` (`version_not_published`) |
| Planner's own calories are incomplete for that version (PLN-04) | `FAILED_PRECONDITION` (`calories_incomplete`); fix the dish first |
| Called without `x-admin-user-id` | `PERMISSION_DENIED` |

## Acceptance criteria

- [ ] Planner 520 kcal vs dietitian 560 kcal → 7.1% → `ok`; vs 600 → 13.3% → `deviation`.
- [ ] Publishing version 2 of a checked dish puts it back to `dietitian_pending`.
- [ ] The report's gate is false while any starter dish is unchecked.
