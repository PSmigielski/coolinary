# NTF-03-02 SendPush and the push queue

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-03 Push delivery](README.md) |
| Release | R1 |
| Size | L |
| Status | todo |
| Depends on | NTF-01-02, NTF-02-02, NTF-03-01, NTF-04-01, NTF-04-02 |
| Unblocks | NTF-03-03, KIT-09, PLN-17, KIT-12 |

## Goal

Other services ask for a push to some users or to a whole household with one call. notifier applies the user's
choices (category switched off, quiet hours, no repeats) and queues one message per device. The call never waits
for Expo.

## Business rules

1. Target is exactly one of: `user_ids` (1–50 distinct ids) or `household_id`. A household resolves to its active
   members in the mirror (NTF-01-02). `exclude_user_ids` are removed afterwards (e.g. the person who made the
   change).
2. A household unknown to the mirror is not an error: the response has all counts 0 and `unknown_household = true`.
3. `category` must not be `UNSPECIFIED`; `template_id` must be a push template of that category; variables are
   validated by the template. Otherwise `INVALID_ARGUMENT`.
4. Plan rules are not checked here; senders enforce Free/Plus (D5).
5. Per user, in this order:
   1. category switched off in preferences (NTF-04-02) → counted in `skipped_preferences`;
   2. `dedup_key` (optional, ≤ 128 characters): `SET ntf:push:dedup:<user_id>:<dedup_key> 1 NX EX 86400` in
      `redis-state`; when the key exists → counted in `deduplicated`;
   3. no devices → counted in `skipped_no_device`;
   4. otherwise one queue row per device.
6. Per device: the text is rendered in the device's language (NTF-04-01), else the user's preference language, else
   `en`.
7. Quiet hours (NTF-04-02) are checked in the device's time zone. Inside quiet hours, `deliver_after` = the next end
   of quiet hours; otherwise now.
8. `ttl_seconds`: default 86,400 (24 h), minimum 60, maximum 604,800 (7 days). `expires_at = now + ttl`. If
   `deliver_after > expires_at`, the row is stored directly as `expired` and never sent.
9. More than 500 device rows for one call → `INVALID_ARGUMENT` (a household has at most 5 accounts × 10 devices).
10. Queue rows and their delivery rows (status `queued`) are written in one transaction, then `NOTIFY push_queue`
    wakes the worker (NTF-03-03).
11. The rendered title, body and data are kept in the queue row only until the message is final (NTF-03-04 rule 4).
12. Typical caller deadline is 2 s; the call does only database and Redis work.

## Contract

`contracts/coolinary/notifier/v1/notifier.proto`

```proto
rpc SendPush(SendPushRequest) returns (SendPushResponse);

message SendPushRequest {
  oneof target {
    UserTargets users = 1;
    string household_id = 2;
  }
  repeated string exclude_user_ids = 3;
  Category category = 4;
  string template_id = 5;
  map<string, string> variables = 6;   // encoding as in SendEmail
  string path = 7;                     // deep link, starts with "/"
  string dedup_key = 8;                // optional; same key + user within 24 h is sent once
  int32 ttl_seconds = 9;               // 0 = default 86400
}

message UserTargets { repeated string user_ids = 1; }

message SendPushResponse {
  int32 queued = 1;                // device messages queued (including deferred by quiet hours)
  int32 deferred = 2;              // of those, held until quiet hours end
  int32 skipped_preferences = 3;   // users
  int32 skipped_no_device = 4;     // users
  int32 deduplicated = 5;          // users
  int32 expired = 6;               // device messages that could not be sent within ttl
  bool unknown_household = 7;
}
```

## Data

```sql
CREATE TABLE push_queue (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  delivery_id     uuid NOT NULL REFERENCES deliveries (id) ON DELETE CASCADE,
  device_id       uuid NOT NULL,             -- devices.id; no FK so a deleted device leaves the row to be cancelled
  user_id         text NOT NULL,
  category        text NOT NULL,
  template_id     text NOT NULL,
  title           text,                      -- cleared when final
  body            text,
  data            jsonb,
  sound           boolean NOT NULL,
  deliver_after   timestamptz NOT NULL,
  expires_at      timestamptz NOT NULL,
  status          text NOT NULL DEFAULT 'pending'
                  CHECK (status IN ('pending','sending','sent','delivered','failed','expired','cancelled','unknown')),
  attempts        integer NOT NULL DEFAULT 0,
  next_attempt_at timestamptz NOT NULL DEFAULT now(),
  ticket_id       text,
  error_code      text,
  sent_at         timestamptz,
  final_at        timestamptz,
  created_at      timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX push_queue_due ON push_queue (next_attempt_at) WHERE status = 'pending';
CREATE INDEX push_queue_receipts ON push_queue (sent_at) WHERE status = 'sent';
CREATE INDEX push_queue_user ON push_queue (user_id);
```

## Errors

| Case | gRPC status |
|---|---|
| No target, both targets, > 50 user ids, bad category/template/variables, bad `path`, > 500 devices | `INVALID_ARGUMENT` |
| Postgres down | `UNAVAILABLE` |

## Acceptance criteria

- [ ] Household of 3 accounts (one with expiry off, one with no device, one with 2 devices) → `queued = 2`,
      `skipped_preferences = 1`, `skipped_no_device = 1`.
- [ ] Same `dedup_key` twice within 24 h → second call `deduplicated = 1`, nothing queued.
- [ ] Device in Europe/Warsaw, quiet hours 21:00–08:00, call at 22:30 local → row `deliver_after` = 08:00 local
      next day, `deferred = 1`; with `ttl_seconds = 3600` → `expired = 1` instead.
- [ ] One device set to `de`, one to `pl` → each row is rendered in its own language.
- [ ] p99 latency under 100 ms for a household call (metric `rpc_server_duration_seconds`).
