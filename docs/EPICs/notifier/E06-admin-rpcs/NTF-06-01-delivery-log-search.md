# NTF-06-01 Delivery log search and statistics (admin)

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-06 Admin RPCs](README.md) |
| Release | R1 |
| Size | M |
| Depends on | NTF-02-02 |
| Unblocks | ADM-07-04 |

## Goal

When a user says "I never got my code", support can see whether notifier sent it and what the mail server or Expo
answered, without notifier storing the address.

## Business rules

1. `NotifierAdminService` is called only by admin. Every call must carry `x-admin-user-id` metadata; without it →
   `PERMISSION_DENIED`.
2. `ListDeliveries` filters (all optional, combined with AND): `channel`, `status`, `template_id`, `category`,
   `user_id`, `email` (hashed by notifier with `DELIVERY_HASH_KEY`, never stored or logged), `from` / `to`
   (timestamps; the range is at most 31 days; default the last 24 h). Newest first, page size 1–100 (default 50),
   opaque page token.
3. Each row returns the fields of NTF-02-02 except `recipient_hash`.
4. `GetDeliveryStats(from, to)` returns counts per day (UTC), channel, template and status; range at most 90 days.
5. Searching by `email` is a personal-data read; admin audits it (ADM-02-03).

## Contract

`contracts/coolinary/notifier/v1/admin.proto`

```proto
service NotifierAdminService {
  rpc ListDeliveries(ListDeliveriesRequest) returns (ListDeliveriesResponse);
  rpc GetDeliveryStats(GetDeliveryStatsRequest) returns (GetDeliveryStatsResponse);
  // NTF-06-02
  rpc ListAlerts(ListAlertsRequest) returns (ListAlertsResponse);
  rpc AcknowledgeAlert(AcknowledgeAlertRequest) returns (Alert);
  rpc SendTestEmail(SendTestEmailRequest) returns (SendTestEmailResponse);
  rpc SendTestPush(SendTestPushRequest) returns (SendTestPushResponse);
}

message ListDeliveriesRequest {
  string channel = 1; string status = 2; string template_id = 3; string category = 4;
  string user_id = 5; string email = 6;
  string from = 7; string to = 8;                       // RFC 3339
  coolinary.common.v1.PageRequest page = 9;
}
message Delivery {
  string id = 1; string channel = 2; string template_id = 3; string category = 4; string language = 5;
  string user_id = 6; string status = 7; string provider_id = 8; string error_code = 9;
  int32 attempts = 10; bool test = 11; string created_at = 12; string finished_at = 13;
}
message ListDeliveriesResponse { repeated Delivery deliveries = 1; coolinary.common.v1.PageResponse page = 2; }

message GetDeliveryStatsRequest { string from = 1; string to = 2; }
message DeliveryStat { string day = 1; string channel = 2; string template_id = 3; string status = 4; int32 count = 5; }
message GetDeliveryStatsResponse { repeated DeliveryStat stats = 1; }
```

## Data

No schema change (indexes from NTF-02-02).

## Errors

| Case | gRPC status |
|---|---|
| Missing `x-admin-user-id` | `PERMISSION_DENIED` |
| Range over 31 / 90 days, bad timestamp | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] Searching `Ewa@Example.com` finds deliveries sent to `ewa@example.com`.
- [ ] Without the metadata key every admin RPC is refused.
- [ ] Stats for 7 days return one row per day, channel, template and status that has messages.
