# IDN-04-02 Household creation events and default names

| | |
|---|---|
| Service | identity |
| Epic | [IDN-04 Households and members core](README.md) |
| Release | R1 |
| Size | S |
| Depends on | IDN-01-01, IDN-04-01 |
| Unblocks | IDN-03-01, IDN-08-01, KIT-13, PLN-21, OFF-14 |

## Goal

Every new household is announced to kitchen, planner and offers, so they can prepare its default profile and
member rows without calling identity.

## Business rules

1. `createPersonalHousehold(tx, user)` (exists in `src/households.ts`) is the only way a household is created. It is
   used by: sign-up with a new e-mail, anonymous sign-in, a member who leaves or is removed (IDN-08-01), and
   self-healing (IDN-04-01).
2. In the same transaction it emits `household.created` (`household_id`, `owner_user_id`, `owner_member_id`) and
   `household.member_added` for the owner (`role: owner`, `has_account: true`, `brought_fridge_from: null`,
   `previous_member_id`: given by the caller, else null).
3. It returns `{ householdId, memberId }`.
4. New households start on plan Free with `plan_source = default` (IDN-05-01).
5. Default names: the household name is stored as `""` and the app shows a translated default ("Mój dom",
   "Mein Zuhause", "My home") until someone renames it. The owner's display name is the user's name, trimmed, or
   `""` if none was given.
   Decided (Q023): stored empty; the app shows the translated default. Today's English "<name>'s home" goes.

## Contract

No RPC change. `GetHousehold` returns `name: ""` for unnamed households.

## Data

No schema change.

## Events

Publishes `household.created`, `household.member_added`.

## Errors

None beyond the transaction failing (the sign-up fails and can be retried).

## Acceptance criteria

- [ ] A new e-mail sign-up emits `household.created` then `household.member_added` with the same `household_id`.
- [ ] Both events are absent when the creating transaction is rolled back.
- [ ] A household created without a user name has `name = ""` and owner `display_name = ""`.
