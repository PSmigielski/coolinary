# MOB-06-01 Development build and notification setup

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-06 Push notifications](README.md) |
| Release | R1 |
| Size | M |
| Depends on | — |
| Unblocks | MOB-06-02, MOB-06-03, MOB-06-04 |

## Goal

The app can receive remote push notifications at all. On SDK 58 this needs a development build: the Expo docs state
that remote push is not available in Expo Go on Android since SDK 53 and that "a development build is required to
use push notifications". Local notifications still work in Expo Go.

## Business rules

1. Install `expo-notifications` with `bunx expo install expo-notifications`, add its config plugin in `app.json`:
   `["expo-notifications", { "enableRemoteNotifications": true, "enableBackgroundRemoteNotifications": false }]`,
   with the notification icon and colour from the design language.
2. Link the project to EAS (`bunx eas-cli init`) so `extra.eas.projectId` exists in the app config; the push token
   call passes `projectId` explicitly (MOB-06-03).
3. Add an EAS `development` profile (`eas.json`) with `developmentClient: true` for iOS and Android. `make mobile`
   keeps starting the Expo server; a new `make mobile-dev-build` documents how to build and install the development
   client. `mobile/README.md` explains that Expo Go remains fine for everything except remote push.
4. Push credentials: FCM v1 service account for Android and an APNs key for iOS, uploaded to EAS (no secrets in the
   repo). The gateway/notifier side uses the Expo Push API (NTF-03).
5. Android notification channels, created at start before any token request (Android 13+ shows the permission
   prompt only after a channel exists). Ids without `-` or `:`:

   | Channel id | Name | Importance |
   |---|---|---|
   | `expiry` | Use-by reminders | default |
   | `offers` | Offers ending | low |
   | `household` | Household | default |
   | `general` | General | low |
6. `Notifications.setNotificationHandler` is set once at start: while the app is in the foreground, show
   `shouldShowBanner: true`, `shouldShowList: true`, `shouldPlaySound: false`, `shouldSetBadge: false` (the
   deprecated `shouldShowAlert` is not used). The handler answers immediately (it must answer within 3 s).

## Contract

App config and `eas.json` changes only.

## Data

None.

## Errors

| Case | Behaviour |
|---|---|
| Running in Expo Go | Push registration (MOB-06-03) is skipped and logged in dev: "remote push needs a development build" |

## Acceptance criteria

- [ ] A development build on an Android 14 and an iOS device receives a test push sent with Expo's push tool.
- [ ] The four channels appear in Android's app notification settings.
- [ ] Expo Go still starts the app without errors.

## Notes

- https://docs.expo.dev/versions/v58.0.0/sdk/notifications/ and https://docs.expo.dev/push-notifications/push-notifications-setup/
