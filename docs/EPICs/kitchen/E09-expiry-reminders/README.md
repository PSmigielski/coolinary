# KIT-09 Expiry reminders

| | |
|---|---|
| Service | kitchen |
| Release | R1 |
| Depends on | KIT-07 (stock), KIT-10-01 (timezone), KIT-01 (plan from identity KIT-01-06, job lock), NTF-03 (`SendPush`, templates `expiry_today`, `expiry_tomorrow`) |

## Goal

Plus households get one short push a day when food is about to go off, so it gets eaten instead of thrown away
(features.md §2 "Expiry reminders"; Plus only, §6).

## Flow

```
every 15 min ─► households whose local time passed 09:00 and not yet run today (Plus only)
              ─► items expiring today (else tomorrow) ─► SendPush(household, expiry, template, names) ─► run row + log
```

Notifier applies each user's preferences, quiet hours, language and deduplication (NTF-03-02).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-09-01](KIT-09-01-reminder-selection.md) | Which items a reminder lists | R1 | S | KIT-07-05, KIT-01 |
| [KIT-09-02](KIT-09-02-reminder-scheduler.md) | Daily reminder scheduler per household timezone | R1 | M | KIT-09-01, KIT-10-01 |
| [KIT-09-03](KIT-09-03-reminder-push.md) | Send the reminder and record it | R1 | S | KIT-09-02, NTF-03 (`SendPush`) |

## Decisions (taken 1 November 2026)

- Q134 (KIT-09-02): 09:00 local time, in time to plan the day's meals.
