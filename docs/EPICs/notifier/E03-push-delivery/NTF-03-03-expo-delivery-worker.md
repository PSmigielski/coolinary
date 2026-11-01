# NTF-03-03 Expo delivery worker

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-03 Push delivery](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | NTF-03-02 |
| Unblocks | NTF-03-04 |

## Goal

Queued pushes are sent to the Expo Push API in batches, retried when Expo or the network fails, and dead push tokens
are removed.

## Business rules

1. The worker runs inside notifier. It wakes on `NOTIFY push_queue` and at least every 1 s.
2. Each round takes up to 100 rows with `status = 'pending'`, `next_attempt_at <= now()`, `deliver_after <= now()`,
   oldest first, with `FOR UPDATE SKIP LOCKED`, and sets them to `sending`.
3. Rows past `expires_at` become `expired` (delivery row too) and are not sent.
4. Rows whose device no longer exists become `cancelled`.
5. One request per round: `POST https://exp.host/--/api/v2/push/send` with a JSON array of at most 100 messages
   (`to` = token, `title`, `body`, `data`, `sound` (`"default"` or omitted), `channelId` = category on Android,
   `ttl` = seconds left until `expires_at`, `priority` `"default"`). Body gzip-compressed when over 1 KB. Header
   `Authorization: Bearer <EXPO_ACCESS_TOKEN>` when the token is set (secret `expo_access_token` in prod). Timeout
   10 s.
6. Throughput cap: at most 600 messages per second (Expo's per-project limit), so at most 6 requests per second.
7. Tickets come back in the same order as the messages:
   - `ok` → row `sent`, `ticket_id` stored, `sent_at` set; delivery row `sent`.
   - error `DeviceNotRegistered` → the device is deleted (NTF-04-01), the row and the device's other pending rows
     become `cancelled`, delivery row `failed` / `expo_DeviceNotRegistered`.
   - `MessageTooBig` → `failed` and an `error` log (a template bug).
   - `MessageRateExceeded` → back to `pending`, retried with backoff (rule 8).
   - `InvalidCredentials` → `failed`, founder alert `notifier.expo_credentials` (critical).
   - any other error → `failed` with `expo_<name>`.
8. Whole-request failure (network error, timeout, HTTP 429 or 5xx) → every row back to `pending`, `attempts + 1`,
   `next_attempt_at = now + min(2^(attempts-1) s, 300 s)`, or later if `Retry-After` says so. After 8 attempts →
   `failed` / `expo_unavailable`.
9. Crash recovery: at start-up, rows in `sending` for more than 60 s go back to `pending`. A message may then arrive
   twice (accepted: at-least-once).
10. Logs carry `delivery_id`, `status`, `error_code`; never tokens, titles or bodies.

## Contract

No gRPC change. New env: `EXPO_ACCESS_TOKEN` (optional in dev), `EXPO_PUSH_URL` (default as above, overridable for
tests).

## Data

No schema change.

## Errors

See rules 7–8.

## Acceptance criteria

- [ ] Against a mock Expo server: 250 queued rows go out in 3 requests of 100, 100, 50.
- [ ] A `DeviceNotRegistered` ticket deletes the device and cancels its other pending rows.
- [ ] Mock returns 503 three times then 200 → sent with `attempts = 3`.
- [ ] Killing notifier mid-request and restarting resends those rows.
- [ ] Metrics: `push_sent_total{category,result}`, `push_queue_depth` (pending rows), `expo_request_duration_seconds`.

## Notes

Check the Expo Push API docs at implementation for current limits and error names (mobile/AGENTS.md: Expo changes
between SDK releases).
