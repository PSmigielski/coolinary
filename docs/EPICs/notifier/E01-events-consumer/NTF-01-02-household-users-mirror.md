# NTF-01-02 Household users mirror

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-01 Events consumer](README.md) |
| Release | R1 |
| Size | M |
| Depends on | NTF-01-01; IDN-04 (`household.created`, `household.member_added`); IDN-01-03 (payloads of `household.member_removed` and `household.deleted`, emitted later by IDN-08/IDN-09) |
| Unblocks | NTF-03-02 (household target) |

## Goal

notifier knows which accounts belong to each household, so `SendPush` can target "the household" without calling
identity (notifier calls no one).

## Business rules

1. The mirror holds pairs (`household_id`, `user_id`) for members with an account. Profiles without an account
   (children, people without a phone) have no `user_id` and are never stored.
2. `household.created` → add (`household_id`, `owner_user_id`).
3. `household.member_added` with a non-null `user_id` → add (`household_id`, `user_id`). With a null `user_id` →
   nothing.
4. `household.member_removed` with a non-null `user_id` → mark the pair removed (tombstone, rule 6).
5. `household.deleted` → mark every pair of that household removed.
6. Order safety: each pair stores `changed_at` = the `occurred_at` of the event that last changed it. An event older
   than the stored `changed_at` for the same pair is ignored (a retried old `member_added` cannot bring back a
   removed member). Removal keeps the row with `removed_at` set; tombstones older than 14 days are deleted by the
   daily job at 04:00 UTC.
7. Active members of a household = rows with `removed_at IS NULL`.
8. A user normally belongs to one household; the mirror does not enforce that (identity does).
9. `account.deleted` and `account.linked` change the mirror too; they are handled in NTF-01-03 and NTF-01-04.

## Contract

No gRPC change.

## Data

```sql
CREATE TABLE household_users (
  household_id uuid NOT NULL,
  user_id      text NOT NULL,
  changed_at   timestamptz NOT NULL,   -- occurred_at of the last applied event
  removed_at   timestamptz,
  PRIMARY KEY (household_id, user_id)
);
CREATE INDEX household_users_user ON household_users (user_id);
```

## Events

Consumes from `events.identity`: `household.created`, `household.member_added`, `household.member_removed`,
`household.deleted`.

## Errors

| Case | Behaviour |
|---|---|
| Payload missing `household_id` or `member_id` | Handler throws → retried, then `dead_events` (NTF-01-01) |

## Acceptance criteria

- [ ] created(owner A) → member_added(B) → member_added(profile, no user) → active members are A and B.
- [ ] member_removed(B) at t2, then a redelivered member_added(B) with t1 < t2 → B stays removed.
- [ ] household.deleted → no active members; tombstones gone after 14 days.
- [ ] Replaying all events twice gives the same table.

## Notes

docs/EPICs/README.md lists notifier as a consumer of `household.member_removed` only; this task also needs
`household.created`, `household.member_added` and `household.deleted`. To align in the README.
