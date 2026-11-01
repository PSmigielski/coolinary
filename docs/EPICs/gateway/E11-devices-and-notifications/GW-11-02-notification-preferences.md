# GW-11-02 Notification preferences and quiet hours

| | |
|---|---|
| Service | gateway |
| Epic | [GW-11 Devices and notification preference routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-02-01, NTF-04 (`GetPreferences`, `SetPreferences`) |
| Unblocks | MOB-03 |

## Goal

Each user chooses which kinds of notification they receive and a quiet period (features.md §2 "Savings and
household": "Notifications, language, units").

## Business rules

1. Preferences are per user (not per household). Categories (booleans, notifier's defaults when never set):
   `expiry` (R1), `offerEnds` (R2), `household` (R2: someone joined, left, removed you), `savingsReport` (R2).
   Near-shop reminders (R3) are on the phone and not stored here.
2. `quietHours` `{ "start": "HH:MM", "end": "HH:MM" }` in the user's `timezone` (IANA name, sent by the app);
   `null` = none. Start and end may cross midnight (22:00–07:00).
3. A category that the plan doesn't include (e.g. `expiry` on Free, features.md §6) can still be stored; notifier
   and the sending service decide whether anything is sent. The response marks it `availableOnPlan: false`.
4. `PUT` replaces the whole object; `If-Match` with `version`.
5. Calls notifier `GetPreferences(user_id)` / `SetPreferences(user_id, ...)`.

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/notification-preferences` | 200 prefs |
| PUT | `/v1/notification-preferences` | 200 prefs |

```json
prefs { "version": 2, "timezone": "Europe/Warsaw",
        "categories": { "expiry": { "enabled": true, "availableOnPlan": false },
                        "offerEnds": { "enabled": true, "availableOnPlan": false },
                        "household": { "enabled": true, "availableOnPlan": true },
                        "savingsReport": { "enabled": false, "availableOnPlan": false } },
        "quietHours": { "start": "22:00", "end": "07:00" } }
```

`PUT` body sends `categories` as `{ "<name>": true|false }`, `timezone`, `quietHours`. Tag `notifications`,
`auth: true`.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bad time or time zone | 400 | `validation` |
| Version mismatch | 409 | `version_conflict` |

## Acceptance criteria

- [ ] Turning `expiry` off stops expiry pushes for that user only (integration with KIT-09/NTF-03).
- [ ] `quietHours` 22:00–07:00 is accepted and returned unchanged.
- [ ] `timezone: "Mars/Base"` is refused with 400.
