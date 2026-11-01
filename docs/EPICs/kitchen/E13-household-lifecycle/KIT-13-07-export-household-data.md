# KIT-13-07 ExportHouseholdData

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-07, KIT-08-05, KIT-10, KIT-11-06 (visibility), KIT-12 |
| Unblocks | GW-02 (`POST /v1/me/export`, GW-02-06) |

## Goal

A user can download everything kitchen holds about their household, filtered so it never shows another adult's
target (GDPR right of access; features.md §2 "Privacy").

## Business rules

1. `ExportHouseholdData(household_id, member_id)`; `member_id` is the requester. Called only by the gateway (D9).
2. The JSON contains:
   1. `profile`: every section of the kitchen profile;
   2. `members`: rows as `GetMemberProfiles` shows them to the requester (KIT-11-06), plus, for the requester's own
      row, the consent version and time;
   3. `stock`: every current item;
   4. `movements`: the last 24 months (date, ingredient, amount, reason, value), without internal ids;
   5. `savings`: every month;
   6. `receipt_names`: the receipt-name memory;
   7. household products and namings from KIT-04/KIT-05 through their export hook.
3. Keys are snake_case; dates ISO 8601; amounts with their units; ingredient ids with the ingredient's name in the
   requester's language.
4. `schema_version = "kitchen-export/1"`.
5. Size limit 10 MB; a larger export is an `INTERNAL` error (not expected: 1,000 items, 24 months).
6. Read-only; deadline 10 s (gateway `heavy`).
7. No age, weight, height or goal exists in kitchen, so none can appear (the gateway adds the phone-only note).

## Contract

`contracts/coolinary/kitchen/v1/export.proto`

```proto
rpc ExportHouseholdData(ExportHouseholdDataRequest) returns (ExportDataResponse);

message ExportHouseholdDataRequest { string household_id = 1; string member_id = 2; }
message ExportDataResponse { bytes json = 1; string schema_version = 2; }   // same shape in every service (GW-02-06)
```

## Data

No schema change.

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| Requester not a member of the household | `PERMISSION_DENIED` | |
| Export larger than 10 MB | `INTERNAL` | logged as a bug |

## Acceptance criteria

- [ ] A test household's export holds every section listed, with the requester's own target and consent time.
- [ ] A second adult's export does not contain the first adult's target (GW-02-06 end-to-end check).
- [ ] A parent's export contains the children's targets; a locked child's export (if allowed by the gateway) does not
      contain their own target.
- [ ] The export validates against the JSON Schema committed with the task.
