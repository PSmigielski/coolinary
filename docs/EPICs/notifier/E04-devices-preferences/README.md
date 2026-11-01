# NTF-04 Devices and notification preferences

**Goal.** Each phone registers its push token with language and time zone; each account chooses which categories it
gets and its quiet hours. The gateway exposes both (GW-11); the app calls them (MOB-06). The GDPR export includes
this data.

**Flow.** App start / sign-in → gateway → `RegisterDevice` (upsert by token and install). Settings screen → gateway →
`Get/SetPreferences`. Sign-out → `UnregisterDevice`.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [NTF-04-01](NTF-04-01-devices.md) | Register and unregister devices | R1 | M | — |
| [NTF-04-02](NTF-04-02-preferences.md) | Notification preferences and quiet hours | R1 | M | NTF-03-01 |
| [NTF-04-03](NTF-04-03-export-user-data.md) | Export a user's notifier data | R1 | S | NTF-02-02, NTF-04-01, NTF-04-02 |
