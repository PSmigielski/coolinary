# KIT-09-03 Send the reminder and record it

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-09 Expiry reminders](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | KIT-09-02, NTF-03 (`SendPush`, templates `expiry_today` / `expiry_tomorrow`) |
| Unblocks | — |

## Goal

Hand the reminder to notifier in one call, so preferences, quiet hours, languages and devices are handled in one
place.

## Business rules

1. One `SendPush` per household and day: target `household_id`; category `CATEGORY_EXPIRY`; template
   `expiry_today` or `expiry_tomorrow` (KIT-09-01 kind); variables `count` and `names` (encoded as SendPush expects);
   `path = "/fridge?filter=expiring"`; `dedup_key = "expiry:<local date>"`; `ttl_seconds = 43200` (12 h: a reminder
   delivered the next morning is useless).
2. Deadline 2 s. `UNAVAILABLE` or deadline → the run stays open (KIT-09-02 rule 5). `INVALID_ARGUMENT` → logged as a
   bug, run row `skipped`.
3. `unknown_household = true` or 0 queued messages is still a success (`sent`): users may have turned expiry pushes
   off or have no device.
4. Item display names go to notifier and Expo; they are food names, not health data. No other item data is sent.
5. The reminder is not repeated for the same household and date even if notifier later fails to deliver.

## Contract

Calls `notifier.v1.NotifierService/SendPush` (NTF-03-02). No new kitchen RPC.

## Data

No schema change (run row from KIT-09-02).

## Errors

| Case | Behaviour |
|---|---|
| notifier unavailable | retry next tick |
| bad request (template mismatch) | `skipped`, error log |

## Acceptance criteria

- [ ] A household with 3 items expiring today produces one `SendPush` with template `expiry_today`, count 3 and
      3 names, path `/fridge?filter=expiring`.
- [ ] A user who turned expiry pushes off (NTF-04) gets nothing, and the run is recorded as `sent`.
- [ ] Re-running the same day (simulated) sends nothing new (run row + `dedup_key`).
