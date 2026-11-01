# IDN-05-04 Downgrade of a shared household

| | |
|---|---|
| Service | identity |
| Epic | [IDN-05 Plan and entitlements](README.md) |
| Release | R2 |
| Size | S |
| Depends on | IDN-05-01, IDN-07-02, IDN-07-04 |
| Unblocks | — |

## Goal

When a shared household loses Plus, everyone stays and keeps the shared fridge, list and plans; only growing the
household needs Plus again (decision Q187).

## Business rules

1. Decided (Q187): when `changePlan` moves a household from Plus to Free, every member stays. There is no grace
   period, no move-out and no push.
2. While the household is Free and has more than one account member, creating invites (IDN-07-02) and accepting an
   invite into it (IDN-07-04) fail with `FAILED_PRECONDITION` `plus_required`. Leaving and removing still work
   (IDN-08).
3. Members keep full access to the shared fridge, list and plans. Every other Plus feature follows the Free plan for
   everyone, checked by each service through `GetEntitlements` (IDN-05-02, Q004).
4. When the household becomes Plus again, invites work again; nothing else changes.
5. Profiles without accounts stay. Their calorie goals become Free-limited in kitchen (KIT-11), and kitchen merges
   every person's allergens into the household list (Q195, KIT-11-04).
6. `GetHousehold` adds `sharing_limited`: true when the household is Free and has more than one account member, so
   the app can explain why invites are blocked.
7. This deviates from features.md "Sharing needs Plus, held by the owner" for households that were already shared;
   starting to share still needs Plus.

## Contract

No new RPC. `GetHousehold`'s `Household` message gains `bool sharing_limited`.

## Data

No schema change. (The `sharing_grace_until` column planned in IDN-05-01 is not added.)

## Events

None of its own; `household.plan_changed` comes from IDN-05-01. No member events.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18).

| Case | gRPC status |
|---|---|
| Invite created or accepted while Free | `FAILED_PRECONDITION` `plus_required` |

## Acceptance criteria

- [ ] A Plus household with 3 accounts set to Free keeps all 3 members; `sharing_limited` is true; nobody gets a
      push.
- [ ] Creating an invite, or accepting one into that household, returns `plus_required`; leaving still works.
- [ ] Setting the household back to Plus makes invites work again and `sharing_limited` false.

## Notes

Q188 (a `sharing_ends` push) is not needed: nobody loses access. Reason `plan_ended` in `household.member_removed`
is not emitted.
