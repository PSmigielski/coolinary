# PLN-03-03 ImportDishes admin RPC

| | |
|---|---|
| Service | planner |
| Epic | [PLN-03 Library import, AI drafting and review](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-03-01, PLN-03-02, PLN-02-04, PLN-01-06 |
| Unblocks | ADM-04 (import upload) |

## Goal

Admin uploads an import file; planner creates or updates draft dishes, puts them in the review queue and reports what
happened per dish. A dry run shows the report without writing.

## Business rules

1. `ImportDishes(file_bytes, dry_run)`; dry run runs every step except the writes.
2. Per dish, by `slug`:
   - unknown slug → new dish, `source=import`, version 1, status `review`.
   - existing dish whose current (or open draft) content hash equals the file's → `unchanged`, no write.
   - existing published dish, different content → new draft version (PLN-02-04 rule 3), status of the version
     `review`; the published version stays live.
   - existing dish in `draft`/`review` → its open version is replaced by the file's content.
   - existing `retired` dish → `failed` with `DISH_RETIRED` (restore it first).
3. Content hash = SHA-256 of the canonical JSON of the parsed dish (keys sorted, units normalised), stored in
   `dish_versions.content_hash`; makes re-imports idempotent.
4. A dish is written only if it has no blocking problem; others in the same file still import.
5. Each run is recorded in `dish_imports` (who, when, file hash, counts) and every written dish in
   `admin_changes`.
6. Imported texts are `machine_translated=false`. Missing launch languages are filled later by PLN-03-05.
7. Nutrition (PLN-04-01) and derived flags (PLN-02-03) are computed as part of the write.

## Contract

```proto
rpc ImportDishes(ImportDishesRequest) returns (ImportDishesResponse);   // PlannerAdminService

message ImportDishesRequest { bytes file = 1; bool dry_run = 2; }
message ImportDishesResponse {
  string import_id = 1;                  // empty on dry run
  repeated ImportDishResult results = 2;
  repeated UnresolvedName unresolved = 3;
}
message ImportDishResult {
  string slug = 1;
  ImportOutcome outcome = 2;             // CREATED, NEW_VERSION, REPLACED_DRAFT, UNCHANGED, FAILED
  string dish_id = 3;
  repeated ValidationProblem problems = 4;
}
message UnresolvedName { string name = 1; string language = 2; repeated string slugs = 3; }
```

## Data

`planner/migrations/000NN_dish_imports.sql`

```sql
CREATE TABLE dish_imports (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_user_id  text NOT NULL,
    file_sha256    text NOT NULL,
    created        int NOT NULL,
    new_versions   int NOT NULL,
    unchanged      int NOT NULL,
    failed         int NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT now()
);
```

## Errors

| Case | gRPC status | Message |
|---|---|---|
| File > 2 MB or > 500 dishes | INVALID_ARGUMENT | `import file too large` |
| Not YAML / wrong `format` | INVALID_ARGUMENT | `not a coolinary-dishes/1 file` |
| Kitchen down | UNAVAILABLE | nothing written |

## Acceptance criteria

- [ ] Importing `example.yaml` twice gives 5 CREATED then 5 UNCHANGED.
- [ ] Changing one step of a published dish and re-importing gives NEW_VERSION; users still see the old version.
- [ ] A file with 1 broken and 4 good dishes imports 4 and reports 1 FAILED with its problems.
- [ ] Dry run writes no rows in any table.
