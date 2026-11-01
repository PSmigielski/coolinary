# identity events

Identity publishes household and account changes to the Redis stream `events.identity` on redis-state. kitchen,
planner, offers and notifier follow them. Each event is written to the outbox in the same transaction as the change
(`emit()` in `src/events/outbox.ts`) and published in order, at least once.

The definitions are in `src/events/types.ts`. Shared rules for all events: `docs/service-conventions.md` "Events".

## Envelope

Each stream entry has two fields: `type` (the event type) and `envelope` (JSON):

```json
{
  "id": "0193a0d2-5c7e-7b41-9f2a-3e8c1d6b4a90",
  "type": "household.deleted",
  "version": 1,
  "occurred_at": "2026-11-02T08:15:30.123Z",
  "producer": "identity",
  "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9",
  "payload": { "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9", "reason": "owner_deleted" }
}
```

`id` is a UUIDv7 and the idempotency key. `household_id` is the household the event is about; for `account.*`
events it is the account's household at the time, or null if it has none.

## Rules for consumers

- Drop duplicates by envelope `id`: delivery is at least once.
- `household.member_added` and `household.snapshot` are upserts by `member_id` / `household_id`: applying the same
  one twice, or a member that already exists, is not an error.
- Within `events.identity`, events about one household arrive in the order identity committed them. Example:
  `member_added` with `brought_fridge_from = X` always comes before `household.deleted` for X. Don't rely on order
  across households.
- Acknowledge and ignore unknown event types. Ignore unknown payload fields.

## Versioning

Adding an optional field keeps the version. Removing, renaming or changing the meaning of a field creates version 2;
identity emits v2 only after every consumer accepts both versions (one release).

## Payloads

Payloads hold ids, enums and booleans only: no e-mail address, name or other user text. Field names are snake_case.
Ids are UUID strings, except user ids, which are Better Auth text ids (letters and digits). A field that may be null
is always present.

JSON Schemas (2020-12): `contracts/events/identity/<type>.v1.json`. They are generated from `src/events/types.ts`
with `bun run events:schemas` (in `identity/`); commit them with the change. identity's tests fail while a file is
out of date. Use them to validate test fixtures, not live events: they reject any field they don't list, and a live
event may carry an optional field added after your copy.

### `household.created`

A household was created: on sign-up, on an anonymous first start, and for a person who leaves or is removed from a
household. The owner's `household.member_added` follows in the same transaction.

| Field             | Type    | Meaning                     |
| ----------------- | ------- | --------------------------- |
| `household_id`    | UUID    | The new household           |
| `owner_user_id`   | user id | The account that owns it    |
| `owner_member_id` | UUID    | The owner's member id in it |

```json
{
  "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9",
  "owner_user_id": "Xq7LmP2vT9sKcR4wYb8NdF3hJ6gA1eZu",
  "owner_member_id": "07cdc834-8deb-4937-bafa-b4394c952d0b"
}
```

### `household.member_added`

A person became a member of a household: the owner of a new household, someone who accepted an invite, a new
profile, or an anonymous account merged into an existing one. An upsert by `member_id`.

| Field                 | Type                           | Meaning                                                                                          |
| --------------------- | ------------------------------ | ------------------------------------------------------------------------------------------------ |
| `household_id`        | UUID                           | The household                                                                                    |
| `member_id`           | UUID                           | The new member                                                                                   |
| `user_id`             | user id or null                | The member's account; null for a profile                                                         |
| `has_account`         | boolean                        | True for `owner` and `member`, false for `profile`                                               |
| `role`                | `owner` / `member` / `profile` | A profile is a person without an account (a child, someone without a phone)                      |
| `brought_fridge_from` | UUID or null                   | A household whose stock moves into this one; its `household.deleted` follows                     |
| `previous_member_id`  | UUID or null                   | The same person's member id in the household they came from, to carry over their diet and target |

```json
{
  "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9",
  "member_id": "a4bba725-063d-41ec-9220-96582fd724e0",
  "user_id": "Hn4sW8kPz2QmV6tR1cY9bL3xJ7fD5aGe",
  "has_account": true,
  "role": "member",
  "brought_fridge_from": "944f0589-901c-46f2-99cd-d27b41985853",
  "previous_member_id": "2c7e3f56-5129-428a-b5ff-a58c0707c551"
}
```

### `household.member_removed`

A person stopped being a member of a household.

| Field              | Type                                                  | Meaning                                                                           |
| ------------------ | ----------------------------------------------------- | --------------------------------------------------------------------------------- |
| `household_id`     | UUID                                                  | The household they left                                                           |
| `member_id`        | UUID                                                  | Their member id in it                                                             |
| `user_id`          | user id or null                                       | Their account; null for a profile                                                 |
| `reason`           | `left` / `removed` / `profile_deleted` / `plan_ended` | `plan_ended` is reserved and never emitted: members stay after a downgrade (Q187) |
| `new_household_id` | UUID or null                                          | The household they move to; null if they move nowhere (a deleted profile)         |

```json
{
  "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9",
  "member_id": "a4bba725-063d-41ec-9220-96582fd724e0",
  "user_id": "Hn4sW8kPz2QmV6tR1cY9bL3xJ7fD5aGe",
  "reason": "left",
  "new_household_id": "944f0589-901c-46f2-99cd-d27b41985853"
}
```

### `household.plan_changed`

A household's plan changed. React to the up- or downgrade; read the plan itself from identity (`GetEntitlements`,
Q004).

| Field           | Type                           | Meaning                                                                                                           |
| --------------- | ------------------------------ | ----------------------------------------------------------------------------------------------------------------- |
| `household_id`  | UUID                           | The household                                                                                                     |
| `plan`          | `free` / `plus`                | The plan from now on                                                                                              |
| `previous_plan` | `free` / `plus`                | The plan before                                                                                                   |
| `source`        | `default` / `manual` / `store` | `default`: Free with nothing bought; `manual`: set by an admin; `store`: reserved for store purchases (IDN-05-05) |

```json
{
  "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9",
  "plan": "plus",
  "previous_plan": "free",
  "source": "manual"
}
```

### `household.deleted`

A household and everything in it is gone. Delete its data.

| Field          | Type                                                                    | Meaning                                                                                                                                                                                                                                      |
| -------------- | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `household_id` | UUID                                                                    | The deleted household                                                                                                                                                                                                                        |
| `reason`       | `owner_deleted` / `joined_other` / `linked_merge` / `anonymous_cleanup` | `owner_deleted`: its only account was deleted; `joined_other`: its only account joined another household; `linked_merge`: an anonymous household merged into an existing account; `anonymous_cleanup`: a stale anonymous account was removed |

```json
{
  "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9",
  "reason": "owner_deleted"
}
```

### `account.linked`

An anonymous account signed in with an e-mail code. With `merged: false` it was upgraded in place and nothing else
changes (`household_id` = `anonymous_household_id`). With `merged: true` it signed in to an existing account: the
`household.member_added` (with `brought_fridge_from`) and `household.deleted` (`linked_merge`) for the anonymous
household come first.

| Field                    | Type    | Meaning                                          |
| ------------------------ | ------- | ------------------------------------------------ |
| `anonymous_user_id`      | user id | The anonymous account                            |
| `user_id`                | user id | The account it became or was merged into         |
| `household_id`           | UUID    | The account's household after linking            |
| `anonymous_household_id` | UUID    | The anonymous account's household before linking |
| `merged`                 | boolean | True when it was merged into an existing account |

```json
{
  "anonymous_user_id": "Hn4sW8kPz2QmV6tR1cY9bL3xJ7fD5aGe",
  "user_id": "Xq7LmP2vT9sKcR4wYb8NdF3hJ6gA1eZu",
  "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9",
  "anonymous_household_id": "944f0589-901c-46f2-99cd-d27b41985853",
  "merged": true
}
```

### `account.deleted`

An account was deleted. Delete what belongs to the account and to these members. If it was the only account of its
household, `household.deleted` comes first.

| Field        | Type          | Meaning                                                |
| ------------ | ------------- | ------------------------------------------------------ |
| `user_id`    | user id       | The deleted account                                    |
| `member_ids` | list of UUIDs | The account's member ids; empty if it had no household |

```json
{
  "user_id": "Xq7LmP2vT9sKcR4wYb8NdF3hJ6gA1eZu",
  "member_ids": ["07cdc834-8deb-4937-bafa-b4394c952d0b"]
}
```

### `household.snapshot`

The household's current state, republished on an admin's request (IDN-01-04) for consumers that fell behind the
stream's 7 days. An upsert by `household_id`.

| Field          | Type            | Meaning                                                                                                                                                |
| -------------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `household_id` | UUID            | The household                                                                                                                                          |
| `plan`         | `free` / `plus` | Its plan                                                                                                                                               |
| `members`      | list            | Every member, accounts and profiles: `member_id` (UUID), `user_id` (user id or null), `has_account` (boolean), `role` (`owner` / `member` / `profile`) |

```json
{
  "household_id": "69220c46-f782-4cf9-b1e9-007cc895e0f9",
  "plan": "free",
  "members": [
    {
      "member_id": "07cdc834-8deb-4937-bafa-b4394c952d0b",
      "user_id": "Xq7LmP2vT9sKcR4wYb8NdF3hJ6gA1eZu",
      "has_account": true,
      "role": "owner"
    },
    { "member_id": "cc53d33c-2e9c-416d-a9de-666ebbce5764", "user_id": null, "has_account": false, "role": "profile" }
  ]
}
```
