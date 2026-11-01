# GW-11-01 Register and unregister a device

| | |
|---|---|
| Service | gateway |
| Epic | [GW-11 Devices and notification preference routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-02-01, NTF-04 (`RegisterDevice`, `UnregisterDevice`) |
| Unblocks | MOB-06 |

## Goal

Push notifications (expiry reminders, offer-ends reminders, household events) reach the right phones, and stop
reaching a phone after sign-out.

## Business rules

1. `PUT /v1/devices/{installationId}` body `{ "expoPushToken", "platform": "ios" | "android", "appVersion",
   "timezone", "language"? }`. `timezone` is the phone's IANA zone (notifier uses it for quiet hours). `installationId`: UUID the app generates once per install (MOB-06). `expoPushToken` must match
   `^(Exponent|Expo)PushToken\[[A-Za-z0-9_-]+\]$`. `language` defaults to `x-language`.
2. Calls notifier `RegisterDevice(user_id, push_token, platform, installation_id, language, timezone, app_version)`.
   Upsert: the same installation re-registering replaces its token; the same token under another user moves to the
   new user (a shared phone after sign-in) — notifier's rule.
3. `DELETE /v1/devices/{installationId}` (no body) → `UnregisterDevice(user_id, installation_id)` → 204; an unknown
   installation → 204 (sign-out must never fail on this, NTF-04-01 rule 6). Decided (Q130): notifier unregisters by
   installation id, so the app needs no push token at sign-out.
4. Anonymous accounts may register (the first expiry reminder can reach them).
5. Tokens are never logged.

## Contract

| Method | Path | 2xx |
|---|---|---|
| PUT | `/v1/devices/{installationId}` | 200 `{ "installationId", "registeredAt" }` |
| DELETE | `/v1/devices/{installationId}` | 204 |

Tag `notifications`, `auth: true`, idempotency optional.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Bad token or installation id | 400 | `validation` |
| Notifier down | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] Registering twice from the same installation leaves one device in notifier.
- [ ] After `DELETE`, a push to the user does not reach that installation (integration with NTF-03).
- [ ] `DELETE` of an unknown installation returns 204.
- [ ] `DELETE` with no body unregisters the device; no push token is needed.
