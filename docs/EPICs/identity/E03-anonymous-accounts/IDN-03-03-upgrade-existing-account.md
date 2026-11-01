# IDN-03-03 Upgrade into an existing account

| | |
|---|---|
| Service | identity |
| Epic | [IDN-03 Anonymous accounts and upgrade](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | IDN-03-02 |
| Unblocks | KIT-13 (moves stock on `household.member_added.brought_fridge_from`), MOB-01 |

## Goal

A returning user who reinstalled the app, used it anonymously for a while and then signs in to their existing
account ends up in their real household, with the items they added meanwhile moved into its fridge.

## Business rules

1. Trigger: a request to `/api/auth/sign-in/email-otp` with a valid anonymous session, for an e-mail that already
   has an account. Better Auth calls `onLinkAccount` with that existing user as `newUser`.
2. The existing account keeps its household (which may be shared). The anonymous household A is merged into it:
   fridge items move over; A's plans, shopping list and kitchen profile are dropped.
   Decided (Q036): the fridge moves over; plans, list and profile are dropped. The app does not ask.
3. One transaction:
   1. Emit `household.member_added` for the existing household E with the existing member id,
      `has_account: true`, the member's current role, `brought_fridge_from: A`, `previous_member_id` = A's owner
      member id. Consumers treat it as an upsert (nothing changes for the member) and kitchen moves A's stock into E.
   2. Delete household A (its member row, invites; profiles in A are deleted with it).
   3. Emit `household.deleted` for A with reason `linked_merge` (after the event above, so kitchen moves stock
      before erasing A).
   4. Delete the anonymous user row.
   5. Emit `account.linked` with `household_id: E`, `anonymous_household_id: A`, `merged: true`.
4. A's plan does not carry over (A is Free unless an admin changed it; if A was Plus and E is Free, the Plus is lost
   and an `info` log line records the household ids).
5. Quota usage counted on A is not added to E.
6. Failure handling as IDN-03-02 rule 6.

## Contract

No new endpoint.

## Data

No schema change.

## Events

Publishes, in this order: `household.member_added` (E, `brought_fridge_from: A`), `household.deleted` (A),
`account.linked` (`merged: true`).

## Errors

| Case | HTTP |
|---|---|
| Link transaction fails | 500; anonymous session still valid |

## Acceptance criteria

- [ ] Existing account U in household E; a new anonymous account adds items, then signs in as U → JWT
      `householdId = E`; A no longer exists in identity.
- [ ] The three events appear in the order above in `events.identity`.
- [ ] Signing in to an account whose household is shared with two other people leaves those members untouched.
