# OFF-14-03 Export household data

| | |
|---|---|
| Service | offers |
| Epic | [OFF-14 Household lifecycle](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-04-01, OFF-04-02, OFF-10-04 |
| Unblocks | GW-02 (data export, README D9) |

## Goal

A household's data export includes everything offers holds about it.

## Business rules

1. `ExportHouseholdData(household_id)` returns JSON with: my shops (shop id, name, chain, position, state, order),
   shopping preferences, user-added shops and corrections the household submitted (with their status), and user
   flyers (shop, dates, state, created at) with their saved offers.
2. No photos (never stored), no other household's data, no internal fields (match confidence, held reasons).
3. At most 10 MB; the gateway assembles it with the other services' parts (GW-02).

## Contract

```proto
rpc ExportHouseholdData(ExportHouseholdDataRequest) returns (ExportHouseholdDataResponse);
message ExportHouseholdDataRequest { string household_id = 1; }
message ExportHouseholdDataResponse { bytes json = 1; }
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Unknown household | Empty export (not an error) |

## Acceptance criteria

- [ ] A test household with 2 shops, preferences and one confirmed user flyer gets all of it in the export and nothing
      of a second test household.
