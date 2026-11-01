# KIT-10-05 Staples ("always at home")

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | KIT-10-01, KIT-02 (`is_staple_candidate`) |
| Unblocks | KIT-07-07, KIT-10-08, PLN-06, PLN-15 (never counted as missing) |

## Goal

Salt, oil and flour are "always at home": they never show as missing, never land on the shopping list, and need not
be tracked in the fridge (features.md §1 "Staples").

## Business rules

1. A staple is a catalogue ingredient id. 0–50 per household (gateway limit, GW-05-01). No duplicates.
2. **Defaults** before the first save: every ingredient with `is_staple_candidate = true` (KIT-02-01: salt, black
   pepper, sugar, wheat flour, cooking oil — rapeseed in Poland, sunflower in Germany if the catalogue has both —
   water). The app shows them ticked; the user unticks what they don't keep.
3. **Water** is always a staple: it is added to every snapshot (KIT-10-08) and cannot be removed; it is not stored.
4. Staples are never counted as missing by planner (PLN-06, PLN-15). ConsumeStock skips a staple only when there is
   no stock of it (KIT-07-07 rule 7).
5. The section is replaced as a whole on save.
6. Merged ingredients are re-pointed by the KIT-02-09 registry; a duplicate after re-pointing is dropped.

## Contract

Section `staples` of `UpdateProfile` (KIT-10-01): `StapleSection { repeated string ingredient_ids }`.

## Data

`kitchen/migrations/000NN_profile_staples.sql`

```sql
-- +goose Up
CREATE TABLE profile_staples (
    household_id  uuid NOT NULL REFERENCES kitchen_profiles (household_id) ON DELETE CASCADE,
    ingredient_id uuid NOT NULL REFERENCES ingredients (id),
    PRIMARY KEY (household_id, ingredient_id)
);

-- +goose Down
DROP TABLE profile_staples;
```

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Unknown ingredient, more than 50 | `INVALID_ARGUMENT` | `unknown_ingredient` / `invalid_value` |

## Acceptance criteria

- [ ] A new household reads the staple candidates as its staples, with `is_default = true`.
- [ ] Saving staples without water still returns water in `GetPlanningProfile`.
- [ ] features.md Verification 12 ("nothing already at home appears") holds for staples: a plan needing salt and
      flour puts neither on the list.
