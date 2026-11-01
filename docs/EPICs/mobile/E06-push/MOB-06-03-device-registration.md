# MOB-06-03 Device and push token registration

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-06 Push notifications](README.md) |
| Release | R1 |
| Size | M |
| Depends on | MOB-06-01, MOB-06-02, MOB-01-02, NTF-04 (`RegisterDevice` via GW-11 `PUT /v1/devices/{installationId}`) |
| Unblocks | MOB-01-07, MOB-06-04 |

## Goal

notifier always knows the current push token, language and time zone of each phone of a user, and forgets the phone
when the user signs out.

## Business rules

1. **Installation id:** a UUIDv4 created on first launch and kept in SecureStore key `device.installation` (add the key
   to `SecureKey` in `src/lib/secure.ts`). It identifies this app install; it is not a hardware id.
2. Token: `Notifications.getExpoPushTokenAsync({ projectId })`, only in a development or release build, only with
   permission granted, only when `status` is `anonymous` or `signedIn`.
3. Registration call (assumed GW-11): `PUT /v1/devices/{installationId}` with

   ```ts
   interface DeviceRegistration { expoPushToken: string; platform: "ios" | "android"; appVersion: string; language: "pl" | "de" | "en"; timeZone: string }
   ```
   `timeZone` is the IANA name from `expo-localization` (`getCalendars()[0].timeZone`). The gateway binds the device to
   the JWT's user and household.
4. Register when any of these change compared with the last successful registration, stored as a fingerprint in
   AsyncStorage `coolinary.device-registration` (`{ userId, token, language, timeZone, appVersion, at }`, no e-mail):
   user id, token, language, time zone, app version; and at least every 7 days.
5. Listen with `Notifications.addPushTokenListener`; a rotated token re-registers at once.
6. On `userChanged` (anonymous upgrade, sign-in) register again for the new user. On sign-out
   `DELETE /v1/devices/{installationId}` (best effort, MOB-01-07) and delete the fingerprint.
7. Failures retry with 5 s, 30 s, 5 min, then on the next app start; registration is never queued in the outbox.
8. Permission revoked in system settings: on the next foreground, `getPermissionsAsync()` reports denied →
   `DELETE /v1/devices/{installationId}`.

## Contract

| Call | Request | Responses |
|---|---|---|
| `PUT /v1/devices/{installationId}` | `DeviceRegistration` | `204`, `400 validation` |
| `DELETE /v1/devices/{installationId}` | — | `204` (also when unknown) |

## Data

SecureStore `device.installation`; AsyncStorage `coolinary.device-registration`.

## Errors

| Case | Behaviour |
|---|---|
| Expo Go | Skip with a dev log (MOB-06-01) |
| `getExpoPushTokenAsync` throws (no network, FCM issue) | Retry per rule 7 |

## Acceptance criteria

- [ ] After sign-in, notifier's device table shows this installation for the new user and not the anonymous one.
- [ ] Changing the phone's language re-registers once with the new language.
- [ ] After sign-out, a push to the old user does not reach the phone.
