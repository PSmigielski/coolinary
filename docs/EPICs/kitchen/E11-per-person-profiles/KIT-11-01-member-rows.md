# KIT-11-01 Member rows from identity events

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-11 Per-person profiles](README.md) |
| Release | R2 |
| Size | M |
| Depends on | KIT-01 (stream consumer), KIT-10-01, IDN-01-03 (event payloads) |
| Unblocks | KIT-11-02 … 06, KIT-13-02, KIT-13-03 |

## Goal

Kitchen keeps one profile row per household member identity knows about — owner, members and profiles without an
account — so per-person data has a home, without kitchen ever calling identity.

## Business rules

1. `household.member_added` (upsert by `member_id`): create the row with `household_id`, `user_id` (null for
   profiles), `has_account`, `role`. A repeated or late event changes nothing that the user set. The household's
   kitchen profile is created with defaults first if missing (KIT-13-01).
2. With `previous_member_id` set (a person who joined from another household or whose anonymous account was
   merged), the person's own data is copied from that old row — handled in KIT-13-02.
3. `household.snapshot` (upsert by household): rows for listed members are created or updated (`user_id`,
   `has_account`, `role`); rows of the household that are not listed are deleted with their per-person data. Used
   by identity to backfill and resync (IDN-01-04). Households created before R2 get their owner row this way.
4. `household.member_removed`: the row and its per-person data are deleted — KIT-13-03.
5. New rows: `target_kind = adult`, no target, `calories_hidden = false`, `goal_set_by_parent = false`, no
   personal restrictions. A profile without an account added by a parent as a child gets its kind through
   `UpdateMemberProfile` right after (KIT-11-02).
6. Every row change bumps the household's `planning_version` and writes `profile.changed` with section `members`
   (KIT-10-09).
7. Rows hold no names; display names stay in identity and reach the app through `/v1/household`.

## Contract

No RPC. Consumes `household.member_added`, `household.snapshot` (and `household.member_removed`, KIT-13-03).

## Data

`kitchen/migrations/000NN_member_profiles.sql`

```sql
-- +goose Up
CREATE TABLE member_profiles (
    member_id            uuid PRIMARY KEY,
    household_id         uuid NOT NULL REFERENCES kitchen_profiles (household_id) ON DELETE CASCADE,
    user_id              text,                        -- null for profiles without an account
    has_account          boolean NOT NULL,
    role                 text NOT NULL CHECK (role IN ('owner', 'member', 'profile')),
    target_kind          smallint NOT NULL DEFAULT 1 CHECK (target_kind IN (1, 2)),   -- 1 adult, 2 child
    goal_set_by_parent   boolean NOT NULL DEFAULT false,
    daily_kcal_target    integer,                     -- KIT-11-03
    consent_version      text,
    consent_given_at     timestamptz,
    consent_by_member_id uuid,
    calories_hidden      boolean NOT NULL DEFAULT false,
    version              integer NOT NULL DEFAULT 1,
    created_at           timestamptz NOT NULL DEFAULT now(),
    updated_at           timestamptz NOT NULL DEFAULT now(),
    CHECK (daily_kcal_target IS NULL OR (consent_version IS NOT NULL AND consent_given_at IS NOT NULL))
);
CREATE INDEX member_profiles_household ON member_profiles (household_id);
CREATE INDEX member_profiles_user ON member_profiles (user_id) WHERE user_id IS NOT NULL;

-- Per-person rows in profile_restrictions / profile_dislikes (KIT-10-04) reference the member:
ALTER TABLE profile_restrictions
    ADD CONSTRAINT profile_restrictions_member FOREIGN KEY (member_id) REFERENCES member_profiles (member_id) ON DELETE CASCADE;
ALTER TABLE profile_dislikes
    ADD CONSTRAINT profile_dislikes_member FOREIGN KEY (member_id) REFERENCES member_profiles (member_id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE profile_dislikes DROP CONSTRAINT profile_dislikes_member;
ALTER TABLE profile_restrictions DROP CONSTRAINT profile_restrictions_member;
DROP TABLE member_profiles;
```

## Events

Consumes `household.member_added`, `household.snapshot`; publishes `profile.changed` (`members`).

## Errors

| Case | Behaviour |
|---|---|
| Event payload missing `member_id` | consumer retries, then dead letter (KIT-01) |
| Member already exists | upsert, no change to user-set fields |

## Acceptance criteria

- [ ] `household.member_added` for an owner creates one adult row with no target; the same event again changes
      nothing.
- [ ] A snapshot listing members A and B for a household that has rows A, B, C deletes C and its restrictions.
- [ ] Each change bumps `planning_version` once and writes `profile.changed` with `members`.
- [ ] A household created in R1 gets its owner row from a `household.snapshot` backfill.
