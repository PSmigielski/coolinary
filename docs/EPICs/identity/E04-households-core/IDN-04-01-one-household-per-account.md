# IDN-04-01 One household per account, one owner per household

| | |
|---|---|
| Service | identity |
| Epic | [IDN-04 Households and members core](README.md) |
| Release | R1 |
| Size | S |
| Depends on | — |
| Unblocks | IDN-04-02 … IDN-04-05, IDN-07-04, IDN-08-01 |

## Goal

Every account belongs to exactly one household at any moment and every household has exactly one owner, so "your
household" in the JWT is never ambiguous.

## Business rules

1. An account (user id) appears in at most one `household_members` row. Enforced by a unique partial index.
2. Every household has exactly one member with role `owner`, that member has an account, and
   `households.owner_user_id` equals that member's `user_id`. Enforced by a unique partial index on
   `(household_id) WHERE role = 'owner'` plus code paths that change the owner only through one function
   (`setOwner(tx, householdId, memberId)`, used by IDN-08-03 and IDN-09-02).
3. `household_members.created_at` is the join time ("joined_at" in RPCs). No new column.
4. Self-healing: when a JWT is issued (`definePayload`) for an account that has no household — for example after a
   failed upgrade (IDN-03-02 rule 6) — identity first creates a personal household (IDN-04-02), logs a warning and
   increments `identity_household_healed_total`. A JWT is never issued without `householdId`.
5. The migration fails if existing data breaks rules 1–2 (only possible in dev; reset the dev database).

## Contract

No RPC change.

## Data

`identity/src/db/schema.ts`:

```ts
uniqueIndex("household_members_user_unique").on(t.userId).where(sql`user_id IS NOT NULL`),
uniqueIndex("household_members_one_owner").on(t.householdId).where(sql`role = 'owner'`),
```

## Errors

| Case | Result |
|---|---|
| Code tries to add a second household for an account | unique violation → `INTERNAL` (a bug) |

## Acceptance criteria

- [ ] Inserting a second member row for the same user id fails at the database.
- [ ] Inserting a second owner for one household fails at the database.
- [ ] Deleting a test user's household and then requesting a token creates a new personal household and logs the
      warning; the JWT carries the new `householdId`.
