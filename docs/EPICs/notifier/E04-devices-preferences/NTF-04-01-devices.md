# NTF-04-01 Register and unregister devices

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-04 Devices and notification preferences](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | — |
| Unblocks | NTF-03-02, GW-11 (device routes), MOB-06 |

## Goal

Each phone that allows notifications registers its Expo push token with its language and time zone, so pushes reach
the right phones in the right language at the right local time.

## Business rules

1. `RegisterDevice` stores: `user_id` (filled by the gateway from the JWT), `push_token`, `platform` (`ios` /
   `android`), `installation_id` (a random id the app keeps per install, MOB-06), `language`, `timezone`,
   `app_version`.
2. `push_token` must match `^Expo(nent)?PushToken\[[A-Za-z0-9_-]{1,200}\]$`; otherwise `INVALID_ARGUMENT`.
3. `language` is resolved per NTF-02-01 rule 2 and stored resolved (`pl`, `de`, `en`). `timezone` must be an IANA
   name known to `Intl.supportedValuesOf("timeZone")`; otherwise `INVALID_ARGUMENT`.
4. Upsert rules, in order:
   1. the token already exists for another user → it moves to this user (the phone changed account);
   2. the same `user_id` + `installation_id` exists with another token → that row's token is replaced (token
      rotated);
   3. otherwise a new row.
   Every register sets `last_seen_at = now()`; the app registers on each start and sign-in.
5. Decided (Q131): at most 10 devices per user. Registering an 11th deletes the one with the oldest `last_seen_at`.
6. Decided (Q130): `UnregisterDevice(user_id, installation_id)` deletes the row of that user and install, so the app
   doesn't have to keep the push token until sign-out. `push_token` is accepted instead of `installation_id` for old
   app builds; one of the two is required. Unknown install or token, or one of another user → OK, nothing changes
   (idempotent; the app calls it on sign-out).
7. Decided (Q132): a daily job at 04:15 UTC deletes devices not seen for 90 days.
8. Deleting a device cancels its pending queue rows.
9. Anonymous accounts register devices like any account; NTF-01-04 moves them on sign-in.
10. Tokens are never logged; logs use the device `id`.

## Contract

`contracts/coolinary/notifier/v1/notifier.proto`

```proto
rpc RegisterDevice(RegisterDeviceRequest) returns (RegisterDeviceResponse);
rpc UnregisterDevice(UnregisterDeviceRequest) returns (UnregisterDeviceResponse);

enum Platform { PLATFORM_UNSPECIFIED = 0; PLATFORM_IOS = 1; PLATFORM_ANDROID = 2; }

message RegisterDeviceRequest {
  string user_id = 1;
  string push_token = 2;
  Platform platform = 3;
  string installation_id = 4;   // <= 64 chars
  string language = 5;          // BCP 47
  string timezone = 6;          // IANA, e.g. "Europe/Warsaw"
  string app_version = 7;       // <= 32 chars
}
message RegisterDeviceResponse { string device_id = 1; }

message UnregisterDeviceRequest { string user_id = 1; string push_token = 2; string installation_id = 3; }  // one of push_token / installation_id
message UnregisterDeviceResponse {}
```

## Data

```sql
CREATE TABLE devices (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         text NOT NULL,
  push_token      text NOT NULL UNIQUE,
  platform        text NOT NULL CHECK (platform IN ('ios', 'android')),
  installation_id text NOT NULL,
  language        text NOT NULL,
  timezone        text NOT NULL,
  app_version     text NOT NULL DEFAULT '',
  created_at      timestamptz NOT NULL DEFAULT now(),
  last_seen_at    timestamptz NOT NULL DEFAULT now(),
  UNIQUE (user_id, installation_id)
);
CREATE INDEX devices_user ON devices (user_id);
```

## Errors

| Case | gRPC status |
|---|---|
| Empty `user_id`, bad token, platform unspecified, bad time zone | `INVALID_ARGUMENT` |
| Unregister with neither `installation_id` nor `push_token` | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] Register A on phone 1, then B on phone 1 with the same token → one row, owned by B.
- [ ] Same install registers a new token → one row with the new token.
- [ ] 11th device removes the least recently seen one.
- [ ] Unregister twice → both OK.
- [ ] Unregister by `installation_id` alone removes that install's row; another user's install id changes nothing.
