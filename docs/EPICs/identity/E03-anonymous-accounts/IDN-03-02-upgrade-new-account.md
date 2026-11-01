# IDN-03-02 Upgrade to a new e-mail account

| | |
|---|---|
| Service | identity |
| Epic | [IDN-03 Anonymous accounts and upgrade](README.md) |
| Release | R1 |
| Size | M |
| Depends on | IDN-03-01, IDN-02-02, IDN-01-01 |
| Unblocks | MOB-01, NTF-01 (`account.linked`) |

## Goal

Someone who started anonymously and then signs in with an e-mail that has no account keeps their fridge, plans,
list and quota counters, because they keep the same household and member.

## Business rules

1. Trigger: a request to `/api/auth/sign-in/email-otp` that carries a valid anonymous session, for an e-mail with no
   account. Better Auth creates the new user and calls `onLinkAccount({ anonymousUser, newUser })`.
2. The user-create hook must not create a personal household for this new user. It skips when the request carries
   an anonymous session. If Better Auth's hook context can't tell, the hook creates one and step 4 deletes it again,
   emitting `household.deleted` with reason `linked_merge`; the acceptance test proves only one household survives.
3. In `onLinkAccount`, one transaction:
   1. `households.owner_user_id` and `household_members.user_id`: anonymous user id → new user id.
   2. The member's `display_name` stays; if it is empty and the new user gave a `name`, it becomes that name
      (trimmed, at most 30 characters).
   3. Delete the anonymous user row (its sessions cascade). The household is safe because it no longer references
      that user.
   4. Emit `account.linked` with `anonymous_user_id`, `user_id`, `household_id`, `anonymous_household_id` =
      `household_id`, `merged: false`.
4. No `account.deleted` and no `household.member_*` events are emitted: the member id and household id don't change,
   so kitchen, planner and offers have nothing to do. Notifier moves devices and preferences to the new user id
   (NTF-01).
5. The household's plan and the gateway's quota counters (keyed by household) carry over unchanged.
6. If the transaction fails, the sign-in answers 500 and the anonymous session stays valid; the person can retry.
   A new user left without a household by such a failure gets one on their next token (IDN-04-01 rule 4).
7. Afterwards the JWT has `anon: false`, the `email` claim, and the same `householdId` and `memberId` as before.

## Contract

No new endpoint; behaviour of `POST /api/auth/sign-in/email-otp` with an anonymous session.

## Data

No schema change.

## Events

Publishes `account.linked` (`merged: false`).

## Errors

| Case | HTTP |
|---|---|
| Link transaction fails | 500; anonymous session still valid |

## Acceptance criteria

- [ ] Anonymous user with household H and member M signs in with a new e-mail → JWT `householdId = H`,
      `memberId = M`, `anon: false`.
- [ ] Exactly one household exists for the new user afterwards; the anonymous user row is gone.
- [ ] `events.identity` has one `account.linked` and no `account.deleted` for this flow.
- [ ] A plan set to Plus on H by admin before the upgrade is still Plus after it.
