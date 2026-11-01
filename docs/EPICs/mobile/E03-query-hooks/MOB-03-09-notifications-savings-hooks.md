# MOB-03-09 Notification preference and savings hooks

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R1 (notification preferences), R2 (savings), R3 (share card) |
| Size | S |
| Depends on | MOB-03-01, MOB-03-02, GW-11, GW-12 |
| Unblocks | later settings and savings screens |

## Goal

Hooks for what the user wants to be notified about and for the savings counters ("€38 and 4 kg saved this month").

## Business rules

1. Hooks and routes (assumed GW-11/GW-12 paths):

   | Hook | HTTP | Rel. | Offline |
   |---|---|---|---|
   | `useNotificationPreferences()` | `GET /v1/notification-preferences` | R1 | cached |
   | `useSetNotificationPreferences()` | `PUT /v1/notification-preferences` `{ categories: { expiry, offer_ends, household, savings_report }: boolean, quietHours?: { from: "22:00", to: "07:00" } }` | R1 | queued |
   | `useSavings(month)` | `GET /v1/savings?month=YYYY-MM` | R2 | cached |
   | `useSavingsCard(month)` | `GET /v1/savings/card?month=` → card data; the image is drawn on the phone (GW-12-02) | R3 | never |
2. Preferences are per user (user-scoped key). Categories that need Plus (expiry reminders) are shown as locked on
   Free by `usePlusGate`; the server ignores them for Free households anyway.
3. Quiet hours are local times in the user's time zone; the time zone is sent with the device registration
   (MOB-06-03).
4. Money is shown from `Money` (`amountMinor`, `currency`) as given; the phone does no currency conversion.

## Contract

Types aliased in `src/api/types/notifications.ts` and `src/api/types/savings.ts`.

## Data

Query cache: preferences and savings persisted.

## Errors

| Case | HTTP | Behaviour |
|---|---|---|
| Free household saves a Plus-only category (expiry) | 200 | accepted but stays locked in the UI; the server ignores it for Free (rule 2) |
| Invalid quiet hours (same start and end, bad format) | 400 `validation` | roll back the optimistic change; show the field error |
| Savings for a month before the household existed or in the future | 404 / 400 | show an empty month |
| Savings report or card on Free | 403 `plus_required` | open the Plus offer (MOB-02 error model) |
| Savings card while offline | — | not queued (offline: never); show offline state |

## Acceptance criteria

- [ ] Turning off expiry reminders stops them from the next daily run (checked against NTF-04).
- [ ] `useSavings` for a month without data returns zeros, not an error.
