# MOB-06-02 When to ask for notification permission

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-06 Push notifications](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | MOB-06-01, MOB-02-05 |
| Unblocks | MOB-06-03 |

## Goal

The app asks for permission at a moment the user sees why (a reminder before food goes off), so more people say yes
and nobody is asked at first start.

## Business rules

1. Never ask during onboarding or at first start.
2. Ask (system prompt) only after an in-app explanation that the user accepts, at the first of:
   - the household has Plus (expiry reminders are Plus) and the user has just saved fridge items with use-by dates;
   - the user turns on any notification category in settings;
   - (R2) the user joins or creates a shared household.
   Decided (Q135): only these triggers. Free households have no expiry reminders, so asking them early would get a
   "no" for little value.
3. If the user declines the in-app explanation, do not show it again for 30 days. If the system permission is denied,
   never call `requestPermissionsAsync()` again automatically; settings show a link to the system settings instead.
4. Check permission with `Notifications.getPermissionsAsync()`; on iOS read `ios.status` (provisional counts as
   granted for registration).
5. The result is stored in AsyncStorage `coolinary.push-ask`: `{ lastAskedAt, outcome: "granted" | "denied" | "dismissed" }`.

## Contract

```ts
export function shouldAskForPush(trigger: "fridge_items_with_dates" | "category_enabled" | "shared_household"): boolean;
export function askForPush(): Promise<"granted" | "denied">;
```

## Data

AsyncStorage `coolinary.push-ask`.

## Errors

| Case | Behaviour |
|---|---|
| System permission denied | never prompt again automatically; settings show a link to the system settings (rule 3) |
| User dismisses the in-app explanation | don't show it again for 30 days (rule 3) |
| `getPermissionsAsync` or `requestPermissionsAsync` throws (e.g. Expo Go without push support) | treat as `denied` for this session, store nothing, and let MOB-06-01 report "push needs a development build" in dev |
| AsyncStorage write of `coolinary.push-ask` fails | the ask still happens; next start may show the explanation once more |

## Acceptance criteria

- [ ] A fresh install never shows the system prompt before the first trigger.
- [ ] Dismissing the explanation suppresses it for 30 days (clock-shifted test).
