# PLN-16-01 Base pantry lists and admin edits

| | |
|---|---|
| Service | planner |
| Epic | [PLN-16 Cuisine base pantry](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-01, KIT-02 (GetIngredients) |
| Unblocks | PLN-16-02, ADM-04 |

## Goal

Store the base pantry of each cuisine as catalogue ingredients, editable in the back office.

## Business rules

1. Cuisine codes as in PLN-02 and the kitchen profile (KIT-10): `polish`, `italian`, `asian`, `mediterranean`,
   `french`, `spanish`, `german`, `mexican`, `middle_eastern`, `american`; sub-cuisines only under `asian`:
   `chinese`, `japanese`, `thai`, `indian`, `vietnamese`, `korean`. A pantry key is the cuisine, or
   `asian/<sub-cuisine>`.
2. A base pantry is an ordered list of 4–12 ingredient ids per key. An Asian sub-cuisine without its own list
   uses `asian`.
3. Seed for the 5 launch cuisines (features.md §7): Polish, German, Italian, Asian, Mediterranean; e.g. Asian: soy
   sauce, rice, ginger, sesame oil; Italian: olive oil, passata, pasta, hard cheese (features.md §2). The rest of the
   seed is written with the starter library and checked in review.
4. Every ingredient id must exist in the kitchen catalogue (`GetIngredients`); unknown ids are rejected.
5. Edits go through `PlannerAdminService` and record `x-admin-user-id` and time.
6. `catalogue.ingredient_merged` rewrites ids here (PLN-21-04).

## Contract

`contracts/coolinary/planner/v1/admin.proto` (`PlannerAdminService`)

```proto
rpc ListBasePantries(ListBasePantriesRequest) returns (ListBasePantriesResponse);
rpc SetBasePantry(SetBasePantryRequest) returns (BasePantry);
message BasePantry { coolinary.common.v1.Cuisine cuisine = 1; repeated string ingredient_ids = 2; string updated_by = 3; string updated_at = 4; }
message SetBasePantryRequest { coolinary.common.v1.Cuisine cuisine = 1; repeated string ingredient_ids = 2; }
message ListBasePantriesRequest {}
message ListBasePantriesResponse { repeated BasePantry pantries = 1; }
```

## Data

`planner/migrations/000NN_base_pantries.sql`

```sql
CREATE TABLE base_pantries (
    cuisine       text NOT NULL,
    position      smallint NOT NULL,
    ingredient_id text NOT NULL,
    PRIMARY KEY (cuisine, ingredient_id)
);
CREATE TABLE base_pantry_changes (
    cuisine     text NOT NULL,
    changed_by  text NOT NULL,          -- admin user id
    changed_at  timestamptz NOT NULL DEFAULT now(),
    ingredient_ids text[] NOT NULL
);
```

## Errors

| Case | gRPC status |
|---|---|
| Fewer than 4 or more than 12 items | `INVALID_ARGUMENT` |
| Unknown ingredient id | `INVALID_ARGUMENT` with the ids |

## Acceptance criteria

- [ ] The 5 launch cuisines have seeded lists after migration.
- [ ] `asian/thai` without its own list returns the `asian` list.
- [ ] An edit by an admin is stored with their id.
