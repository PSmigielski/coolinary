# GW-11 Devices and notification preference routes

Release R1 (expiry reminders are in the launch release).

## Goal

The app registers its Expo push token per installation, removes it on sign-out, and lets each user choose which
notifications they get and when they are quiet. Notifier owns devices and preferences (NTF-04).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-11-01](GW-11-01-devices.md) | Register and unregister a device | R1 | S | GW-02-01, NTF-04 |
| [GW-11-02](GW-11-02-notification-preferences.md) | Notification preferences and quiet hours | R1 | S | GW-02-01, NTF-04 |
