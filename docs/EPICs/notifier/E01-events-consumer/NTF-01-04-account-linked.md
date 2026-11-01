# NTF-01-04 Move an anonymous account's devices to the signed-in account

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-01 Events consumer](README.md) |
| Release | R1 |
| Size | S |
| Depends on | NTF-01-01, NTF-01-02, NTF-03-02, NTF-04-01, NTF-04-02; IDN-03 (`account.linked`) |
| Unblocks | — |

## Goal

A person who starts anonymously and later signs in with an e-mail code keeps their push notifications and settings.

## Business rules

1. On `account.linked` (`anonymous_user_id`, `user_id`, `household_id`), in one transaction:
   1. devices: every device of `anonymous_user_id` moves to `user_id`;
   2. preferences: if `user_id` has no preferences row, the anonymous row is moved to `user_id`; if `user_id`
      already has one (they signed in on another phone first), it is kept and the anonymous row is deleted;
   3. household mirror: each active pair (`h`, `anonymous_user_id`) becomes (`h`, `user_id`); if (`h`, `user_id`)
      already exists and is active, the anonymous pair is marked removed;
   4. queued pushes and delivery log rows of `anonymous_user_id` get `user_id`.
2. Same `anonymous_user_id` and `user_id` → no-op. Unknown anonymous user → no-op.
3. The device limit (NTF-04-01 rule 5) is applied after the move: if `user_id` now has more than 10 devices, the
   ones with the oldest `last_seen_at` are deleted.

## Contract

No gRPC change.

## Data

No schema change.

## Events

Consumes `account.linked` from `events.identity`.

## Errors

| Case | Behaviour |
|---|---|
| Postgres error | Rolled back, retried |

## Acceptance criteria

- [ ] Anonymous user with 1 device and expiry reminders off → after linking, the signed-in user has the device and
      expiry still off.
- [ ] Signed-in user already had preferences (expiry on) → stays on; the anonymous preferences row is gone.
- [ ] A household push sent after linking reaches the device.
