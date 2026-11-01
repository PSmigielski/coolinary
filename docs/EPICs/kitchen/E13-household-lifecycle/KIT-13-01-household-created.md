# KIT-13-01 household.created: default profile and stock state

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-13 Household lifecycle](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-01-04 (consumer), KIT-10-01 |
| Unblocks | KIT-07 (stock state row), KIT-09 (timezone), KIT-11-01 |

## Goal

A household exists in kitchen as soon as identity creates it, with a default profile, so the first request after
sign-up never fails for lack of a row.

## Business rules

1. Handler for `household.created` v1 (`household_id`, `owner_user_id`, `owner_member_id`; no country, IDN-01-03):
   1. insert `kitchen_profiles` with defaults (KIT-10-01): empty country, language `en`, timezone
      `Europe/Brussels`, `is_default = true` — the app sets the real values in onboarding;
   2. insert `household_stock_state` with version 0;
   3. both with `ON CONFLICT DO NOTHING`: a request may have created them first (KIT-10-01 rule 1).
2. The owner's member row comes from the `household.member_added` that identity sends right after (KIT-11-01).
3. No plan is stored; kitchen asks identity when it needs it (KIT-01-06).
4. No `profile.changed` is written for a default profile (planner has nothing cached yet).

## Contract

No RPC. Consumes `household.created`.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Rows already exist | no-op |
| Payload without `household_id` | retried, then dead letter (KIT-01-05) |

## Acceptance criteria

- [ ] After a sign-up, `GetProfile` returns defaults without any earlier `UpdateProfile`.
- [ ] Replaying the event changes nothing.
- [ ] A `GetProfile` that arrives before the event creates the rows; the event then is a no-op.
