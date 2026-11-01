# MOB-06 Push notifications

Goal: reminders (use-by in R1, offers ending and household events in R2) reach the phone, permission is asked at a
moment that makes sense, and taps open the right place. Remote push needs a **development build** on SDK 58; Expo Go
cannot receive it.

## Flow

```
start ──> channels (Android) + foreground handler
trigger (Plus + dated items, category on, shared household) ──> explain ──> system prompt
granted ──> getExpoPushTokenAsync({ projectId }) ──> PUT /v1/devices/{installationId} (token, language, time zone)
token rotates / user changes / 7 days ──> register again;  sign-out ──> DELETE /v1/devices/{installationId}
tap ──> data.route (allow-list) ──> open after auth boot
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [MOB-06-01](MOB-06-01-dev-build-notification-config.md) | Development build and notification setup | R1 | M | — |
| [MOB-06-02](MOB-06-02-permission-timing.md) | When to ask for notification permission | R1 | S | MOB-06-01, MOB-02-05 |
| [MOB-06-03](MOB-06-03-device-registration.md) | Device and push token registration | R1 | M | MOB-06-01, MOB-06-02, MOB-01-02, NTF-04 via GW-11 |
| [MOB-06-04](MOB-06-04-receive-and-route.md) | Receiving notifications and opening the right place | R1 | M | MOB-06-03, MOB-03-02, NTF-03 |
