# IDN-07-04 Accept an invite

| | |
|---|---|
| Service | identity |
| Epic | [IDN-07 Invites](README.md) |
| Release | R2 |
| Size | L |
| Depends on | IDN-07-03, IDN-01-01, IDN-04-02 |
| Unblocks | GW-04, KIT-13 (moves stock, copies the person's profile), PLN-21, OFF-14 |

## Goal

The person joins the shared household in one step, with their fridge items if they chose so, and both phones then
show the same fridge and list (features.md Verification 15).

## Business rules

1. `AcceptInvite(code, actor_user_id, bring_fridge)`. In one transaction, locking the invite row and the target
   household row (`SELECT … FOR UPDATE`), re-check everything from the preview (IDN-07-03 rule 4). Concurrent
   accepts: the first wins; the second gets `invite_used` or `household_full`.
2. Refusals: anonymous actor (`sign_in_required`), already a member (`already_member`), owner of a household with
   other account members (`owns_shared_household`), household full (`household_full`), household no longer Plus
   (`plus_required`, IDN-05-04), expired or used invite.
3. `bring_fridge` is honoured only when the actor is the only account member of their current household
   (`can_bring_fridge`); otherwise it is ignored and the answer says `brought_fridge: false`.
4. The new member row: role `member`, display name copied from the actor's current member. If that name clashes in
   the new household, a number is added ("Ola 2") and the person can rename later.
5. The actor's old household O:
   - If the actor was its only account member: O is deleted with its profiles and pending invites.
   - If the actor was a non-owner member of a shared O: the actor's member row in O is removed; O stays.
6. Events, in this order, in the same transaction:
   1. If O stays: `household.member_removed` (O, old member id, `reason: left`, `new_household_id` = target).
   2. `household.member_added` (target, new member id, `has_account: true`, `role: member`,
      `brought_fridge_from` = O if the fridge comes along else null, `previous_member_id` = old member id).
   3. If O is deleted: `household.deleted` (O, `reason: joined_other`).
   Kitchen moves O's stock (when `brought_fridge_from` is set) and copies the person's diet and target from the old
   member id (`previous_member_id`) before it erases O.
7. The invite becomes `accepted` with `accepted_by_user_id` and `ended_at`.
8. A Plus plan on O (manual) is lost with O; the preview warned about it.
9. The answer carries the new household (as `GetHousehold`), `brought_fridge`, and `token_stale: true`; the app
   fetches a new JWT at once (IDN-04-03 rule 4).
10. Decided (Q184): the owner gets a push "Ewa joined your household", category `household` (NTF-03).

## Contract

```proto
rpc AcceptInvite(AcceptInviteRequest) returns (AcceptInviteResponse);

message AcceptInviteRequest { string code = 1; string actor_user_id = 2; bool bring_fridge = 3; }
message AcceptInviteResponse {
  Household household = 1;
  bool brought_fridge = 2;
  bool token_stale = 3;
}
```

## Data

No schema change.

## Events

Publishes `household.member_removed` (if O stays), `household.member_added`, `household.deleted` (if O is deleted),
in that order.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status / reason |
|---|---|
| Any refusal in rule 2 | `FAILED_PRECONDITION` with that reason |
| Unknown or cancelled code | `NOT_FOUND` `invite_not_found` |
| Too many unknown codes | `RESOURCE_EXHAUSTED` |

## Acceptance criteria

- [ ] Verification 15: a second test account joins by QR, by code and by link (three separate invites); a used or
      expired invite is refused; afterwards both accounts' JWTs carry the same `householdId`.
- [ ] With `bring_fridge: true` from a personal household, the events arrive as `member_added` (with
      `brought_fridge_from`) then `household.deleted`.
- [ ] Two accounts accepting the same invite at the same moment: exactly one joins.
- [ ] A sixth account can never join, even with two invites accepted at the same moment.
