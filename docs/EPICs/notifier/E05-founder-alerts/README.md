# NTF-05 Founder alerts

**Goal.** One RPC for "a human must look at this" (flyer source changed or blocked, Expo credentials failed, events
lost), deduplicated by key and throttled, e-mailed to the founder and listed in admin.

**Flow.** service → `SendAlert(key, severity, summary, …)` → open alert per key → e-mail on open and after the
throttle window → `resolved` closes it.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [NTF-05-02](NTF-05-02-alert-email.md) | Founder alert e-mail and recipients | R1 | S | NTF-02-01 |
| [NTF-05-01](NTF-05-01-send-alert.md) | SendAlert with deduplication and throttling | R1 | M | NTF-05-02 |
