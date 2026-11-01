# NTF-04-03 Export a user's notifier data

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-04 Devices and notification preferences](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | NTF-02-02, NTF-04-01, NTF-04-02 |
| Unblocks | GW-02 (data export, D9) |

## Goal

The GDPR data export assembled by the gateway (D9) includes what notifier holds about the person.

## Business rules

1. `ExportUserData(user_id)` returns:
   - devices: platform, language, time zone, app version, created, last seen, and the token shortened to its last
     6 characters (`…a1B2c3`);
   - preferences: the full `Preferences` (defaults when none stored);
   - deliveries of the last 90 days, counted per channel, template and status (no message content is stored).
2. Unknown user → empty lists and default preferences (not `NOT_FOUND`).
3. The response is at most 1 MB; with 10 devices and counts it is far smaller.

## Contract

`contracts/coolinary/notifier/v1/notifier.proto`

```proto
rpc ExportUserData(ExportUserDataRequest) returns (ExportUserDataResponse);

message ExportUserDataRequest { string user_id = 1; }
message ExportUserDataResponse {
  repeated ExportedDevice devices = 1;
  Preferences preferences = 2;
  repeated DeliveryCount deliveries = 3;
}
message ExportedDevice {
  string platform = 1; string language = 2; string timezone = 3; string app_version = 4;
  string created_at = 5; string last_seen_at = 6; string token_suffix = 7;
}
message DeliveryCount { string channel = 1; string template_id = 2; string status = 3; int32 count = 4; }
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Empty `user_id` | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] Export for a user with 2 devices shows both with token suffixes only.
- [ ] After `account.deleted` the export is empty.

## Notes

docs/EPICs/README.md lists notifier among the services with export RPCs but gives it no lifecycle epic; the export
lives here.
