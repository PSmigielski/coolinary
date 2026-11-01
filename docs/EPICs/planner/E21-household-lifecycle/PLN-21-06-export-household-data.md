# PLN-21-06 ExportHouseholdData

| | |
|---|---|
| Service | planner |
| Epic | [PLN-21 Household lifecycle](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-15-01, PLN-08 |
| Unblocks | GW-02 (data export) |

## Goal

Return everything planner holds about a household, for the gateway's GDPR export (D9).

## Business rules

1. Content (JSON, one document): day and week plans with meals, dishes and portions; attendance marks and rules;
   the shopping list with lines and sources; favourites, ratings and cooking history; household-private dishes with
   ingredients and steps; top-ups; offer reminders sent.
2. Per-person numbers follow the visibility rules for the requesting member (PLN-12-07): the requester gets their
   own kcal numbers and their children's; never another adult's.
3. Dates in ISO format, amounts with units, ids as stored; ingredient names in the requester's language.
4. Limits: the last 24 months of plans and history; the response is streamed in pages of 1 MB at most.
5. Called only by the gateway; requires `household_id` and `member_id` of the requester.

## Contract

`contracts/coolinary/planner/v1/export.proto`

```proto
rpc ExportHouseholdData(ExportHouseholdDataRequest) returns (ExportHouseholdDataResponse);
message ExportHouseholdDataRequest { string household_id = 1; string member_id = 2; string page_token = 3; }
message ExportHouseholdDataResponse { bytes json_chunk = 1; string next_page_token = 2; string schema_version = 3; }
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Member not in household | `PERMISSION_DENIED` |
| Bad page token | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] An export for the worked-example household as Ewa contains her kcal numbers and no number of Adam's.
- [ ] Concatenated chunks form valid JSON with every section present (empty arrays where nothing exists).
