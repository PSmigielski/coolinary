# NTF-05-01 SendAlert with deduplication and throttling

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-05 Founder alerts](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | NTF-05-02 |
| Unblocks | OFF-05 (flyer source failures), KIT-04 (lookup source down), NTF-01-01, NTF-03-03, NTF-06-02 |

## Goal

Any service can tell the founder that something needs a human (a flyer source changed or blocked access, Expo
credentials failed, events were lost) without flooding the inbox when the same problem repeats.

## Business rules

1. Request fields: `key` (required, 1–120 characters, `[a-z0-9._:-]`, e.g. `offers.source_failed.lidl_de`),
   `severity` (`info`, `warning`, `critical`), `source` (service name), `summary` (1–200 characters), `details`
   (≤ 2,000 characters), `link` (optional `https://…` URL, or `http://10.77.0.x…` in dev), `resolved` (bool).
2. No personal data: any e-mail address pattern in `summary` or `details` is replaced with `[redacted]` before
   storing or sending.
3. One open alert per key. A call for a key without an open alert opens one; a call for a key with an open alert
   adds to its `count` and updates `last_at`.
4. Sending: `warning` and `critical` alerts are e-mailed when opened. While open, a repeat is e-mailed again only
   after the throttle window since `last_sent_at`: `critical` 1 h, `warning` 6 h. That e-mail says
   "repeated N times since <first_at>". Decided (Q040): critical 1 h, warning 6 h.
5. `info` alerts are stored only (visible in admin, NTF-06-02), never e-mailed.
6. `resolved = true` closes the open alert of that key and sends one "resolved" e-mail if the alert had been
   e-mailed; with no open alert it is ignored.
7. An open alert with no new call for 7 days is closed automatically (`resolved_at`, reason `stale`), no e-mail.
8. Recipients: `FOUNDER_ALERT_EMAILS`, comma-separated, at most 5 addresses. When empty, alerts are stored and a
   `warn` line is logged at start-up.
9. E-mails go through the same SMTP sender with template `founder_alert` (NTF-05-02), are exempt from the
   per-address limit (NTF-02-03 rule 5) and are logged in the delivery log like other e-mails.
10. notifier also calls this internally (dead events, Expo credentials).
11. `ai` calls no service, so its problems surface through metrics and Grafana, not this RPC.

## Contract

`contracts/coolinary/notifier/v1/notifier.proto`

```proto
rpc SendAlert(SendAlertRequest) returns (SendAlertResponse);

enum Severity { SEVERITY_UNSPECIFIED = 0; SEVERITY_INFO = 1; SEVERITY_WARNING = 2; SEVERITY_CRITICAL = 3; }

message SendAlertRequest {
  string key = 1;
  Severity severity = 2;
  string source = 3;
  string summary = 4;
  string details = 5;
  string link = 6;
  bool resolved = 7;
}
message SendAlertResponse {
  string alert_id = 1;
  bool emailed = 2;      // an e-mail went out for this call
  bool suppressed = 3;   // counted but throttled
}
```

## Data

```sql
CREATE TABLE alerts (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  key              text NOT NULL,
  severity         text NOT NULL,
  source           text NOT NULL,
  summary          text NOT NULL,
  details          text NOT NULL DEFAULT '',
  link             text,
  count            integer NOT NULL DEFAULT 1,
  first_at         timestamptz NOT NULL DEFAULT now(),
  last_at          timestamptz NOT NULL DEFAULT now(),
  last_sent_at     timestamptz,
  resolved_at      timestamptz,
  resolved_reason  text,            -- resolved | stale
  acknowledged_at  timestamptz,
  acknowledged_by  text             -- admin user id (NTF-06-02)
);
CREATE UNIQUE INDEX alerts_open_key ON alerts (key) WHERE resolved_at IS NULL;
```

## Errors

| Case | gRPC status |
|---|---|
| Bad key, severity unspecified, empty summary, too long fields | `INVALID_ARGUMENT` |
| SMTP down | Alert stored, `emailed = false`; metric `alerts_email_failed_total` + 1; status OK |

## Acceptance criteria

- [ ] 50 `warning` calls with the same key within 1 h → one e-mail, `count = 50`.
- [ ] A call 6 h after the last e-mail → a second e-mail with "repeated 49 times".
- [ ] `resolved` → one "resolved" e-mail; the next call with the key opens a new alert.
- [ ] `details` containing `ewa@example.com` is stored as `[redacted]`.
- [ ] Metric `alerts_total{source,severity}`.
