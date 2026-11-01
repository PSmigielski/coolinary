# NTF-04-02 Notification preferences and quiet hours

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-04 Devices and notification preferences](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | NTF-03-01 (categories) |
| Unblocks | NTF-03-02, GW-11, MOB-03 |

## Goal

Each person chooses which kinds of push they get and when they must not be disturbed. The choice is stored once per
account and applies to all their phones (features.md §2 "Savings and household": notifications in settings).

## Business rules

1. Preferences belong to a `user_id` (from the JWT). Without a stored row the defaults apply.
2. Defaults: `expiry` on, `offer_ends` on, `household` on, `savings` on, `near_shop` **off** (opt-in, features.md §8
   "near-shop reminders are opt-in"); quiet hours on, 21:00–08:00; language not set (device language is used).
   Decided (Q133): the default quiet hours are 21:00–08:00.
3. `SetPreferences` is partial: only fields present in the request change (`optional` fields).
4. Quiet hours: `start` and `end` are local times with minute precision (`HH:MM`), may cross midnight, and must
   differ; otherwise `INVALID_ARGUMENT`. They apply in each device's own time zone.
5. `language`, when set, is one of `pl`, `de`, `en`; empty string clears it.
6. Turning a category off cancels that user's pending queue rows of that category (deferred messages included).
7. Preferences apply to Free and Plus alike; whether a message is sent at all is the sender's rule (D5).
8. `GetPreferences` always returns every field, filled with defaults where nothing is stored.

## Contract

`contracts/coolinary/notifier/v1/notifier.proto`

```proto
rpc GetPreferences(GetPreferencesRequest) returns (Preferences);
rpc SetPreferences(SetPreferencesRequest) returns (Preferences);

message GetPreferencesRequest { string user_id = 1; }

message Preferences {
  bool expiry = 1;
  bool offer_ends = 2;
  bool household = 3;
  bool savings = 4;
  bool near_shop = 5;
  bool quiet_hours_enabled = 6;
  string quiet_hours_start = 7;   // "21:00"
  string quiet_hours_end = 8;     // "08:00"
  string language = 9;            // "" = device language
}

message SetPreferencesRequest {
  string user_id = 1;
  optional bool expiry = 2;
  optional bool offer_ends = 3;
  optional bool household = 4;
  optional bool savings = 5;
  optional bool near_shop = 6;
  optional bool quiet_hours_enabled = 7;
  optional string quiet_hours_start = 8;
  optional string quiet_hours_end = 9;
  optional string language = 10;
}
```

## Data

```sql
CREATE TABLE notification_preferences (
  user_id             text PRIMARY KEY,
  expiry              boolean NOT NULL DEFAULT true,
  offer_ends          boolean NOT NULL DEFAULT true,
  household           boolean NOT NULL DEFAULT true,
  savings             boolean NOT NULL DEFAULT true,
  near_shop           boolean NOT NULL DEFAULT false,
  quiet_hours_enabled boolean NOT NULL DEFAULT true,
  quiet_hours_start   time NOT NULL DEFAULT '21:00',
  quiet_hours_end     time NOT NULL DEFAULT '08:00',
  language            text,
  updated_at          timestamptz NOT NULL DEFAULT now()
);
```

## Errors

| Case | gRPC status |
|---|---|
| Empty `user_id`, bad time format, start = end, unsupported language | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] A new user gets the defaults; `near_shop` is off.
- [ ] Setting only `savings = false` leaves the rest unchanged.
- [ ] Quiet hours 22:00–06:30 accepted; 08:00–08:00 refused.
- [ ] Turning `expiry` off cancels that user's pending expiry pushes.
