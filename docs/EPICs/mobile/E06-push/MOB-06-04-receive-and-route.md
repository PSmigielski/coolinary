# MOB-06-04 Receiving notifications and opening the right place

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-06 Push notifications](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | MOB-06-03, MOB-03-02, NTF-03 (push payload) |
| Unblocks | later screens linked from notifications |

## Goal

Tapping a reminder opens the screen it is about (the expiring items, the list, the invite), also when the app was
closed, and a notification for another household never opens the wrong data.

## Business rules

1. Payload in the notification's `data` (assumed NTF-03; align when written):

   ```ts
   interface PushData { v: 1; category: "expiry" | "offer_ends" | "household" | "savings_report" | "general"; route: string; householdId?: string; refreshKeys?: string[] }
   ```
2. Routing allow-list (Expo Router paths; the screens come later, unknown routes open the home screen):
   `/fridge?filter=expiring`, `/list`, `/offers/[offerId]`, `/household`, `/household/invites`,
   `/savings/[month]`. A `route` outside the list is ignored.
3. Taps: `Notifications.addNotificationResponseReceivedListener` while running; on cold start
   `Notifications.getLastNotificationResponseAsync()` once after auth boot (MOB-01-02) has finished, so the route opens
   with a valid session.
4. If `householdId` is set and differs from the current household, first `forceRefresh()` the JWT; if it still
   differs, open the home screen instead of the route.
5. A notification received in the foreground invalidates the matching data (`expiry` → stock; `offer_ends` → list,
   offers; `household` → household; `savings_report` → savings) so the open screen is current.
6. Each response is handled once (track the last handled `request.identifier` in memory).
7. Payloads never carry personal data or item names that the app needs; notifier writes the visible text, the app only
   uses `data` for routing.

## Contract

As in rule 1.

## Data

Memory only.

## Errors

| Case | Behaviour |
|---|---|
| Malformed `data` | Open the home screen |

## Acceptance criteria

- [ ] With the app killed, tapping an expiry reminder opens the expiring-items route after sign-in boot.
- [ ] A notification carrying another household's id opens the home screen.
- [ ] A foreground expiry notification refreshes the fridge list.
