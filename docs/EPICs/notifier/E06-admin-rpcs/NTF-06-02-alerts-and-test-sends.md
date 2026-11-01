# NTF-06-02 Alerts list, acknowledge and test sends (admin)

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-06 Admin RPCs](README.md) |
| Release | R1 |
| Size | S |
| Depends on | NTF-06-01, NTF-05-01, NTF-02-03, NTF-03-02 |
| Unblocks | ADM-07-01, ADM-07-04 |

## Goal

The founder sees open alerts in admin, marks them as seen, and can send a test e-mail or push to themselves to check
a template or a delivery path.

## Business rules

1. Same access rule as NTF-06-01 (`x-admin-user-id` required).
2. `ListAlerts(state)`: `open` (default) or `all` (last 30 days), newest `last_at` first, at most 200.
3. `AcknowledgeAlert(id)` sets `acknowledged_at` and `acknowledged_by` = the admin user id. It does not close the
   alert and does not stop throttled repeats (only `resolved` or 7 days of silence do, NTF-05-01).
4. `SendTestEmail(to, template_id, language)`: renders the template with its built-in sample variables and sends it.
   `to` must equal the admin's own e-mail address passed by admin (taken from the admin's JWT); notifier checks
   that the request's `to` equals metadata `x-admin-email`, otherwise `PERMISSION_DENIED`. The delivery row has
   `test = true`.
5. `SendTestPush(template_id, language)` sends to the devices of the user in `x-admin-user-id`, with sample
   variables, ignoring preferences and quiet hours; delivery rows have `test = true`. No device → `NOT_FOUND`.
6. At most 10 test sends per admin per hour → `RESOURCE_EXHAUSTED`.

## Contract

`contracts/coolinary/notifier/v1/admin.proto` (service in NTF-06-01)

```proto
message ListAlertsRequest { string state = 1; }        // "open" | "all"
message Alert {
  string id = 1; string key = 2; string severity = 3; string source = 4; string summary = 5; string details = 6;
  string link = 7; int32 count = 8; string first_at = 9; string last_at = 10; string resolved_at = 11;
  string acknowledged_at = 12; string acknowledged_by = 13;
}
message ListAlertsResponse { repeated Alert alerts = 1; }
message AcknowledgeAlertRequest { string id = 1; }
message SendTestEmailRequest { string to = 1; string template_id = 2; string language = 3; }
message SendTestEmailResponse { string delivery_id = 1; }
message SendTestPushRequest { string template_id = 1; string language = 2; }
message SendTestPushResponse { int32 queued = 1; }
```

New metadata key `x-admin-email`, set only by admin on `SendTestEmail`.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Unknown alert id | `NOT_FOUND` |
| `to` is not the admin's own address | `PERMISSION_DENIED` |
| Admin has no device | `NOT_FOUND` |
| Over 10 tests per hour | `RESOURCE_EXHAUSTED` |

## Acceptance criteria

- [ ] Acknowledging an open alert keeps it open and records who acknowledged it.
- [ ] Test e-mail of `sign_in_code` in `pl` arrives in Mailpit with sample code `123 456`.
- [ ] Test e-mail to another address is refused.
