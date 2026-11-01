# PLN-12-07 Who sees which numbers

| | |
|---|---|
| Service | planner |
| Epic | [PLN-12 Per-person portions](README.md) |
| Release | R2 |
| Size | S |
| Depends on | PLN-12-06 |
| Unblocks | PLN-13-07, PLN-14-01, GW-07 |

## Goal

Every planner response with per-person numbers shows each viewer only what they may see: portions for everyone,
kcal and targets only where allowed (features.md "Privacy", "The child is never singled out").

## Business rules

1. The viewer is the member behind `x-user-id`, passed as `viewer_member_id` by the gateway (D3).
2. **Own numbers:** the viewer sees their own portion kcal, extras kcal, day total and target, unless their own
   `calories_hidden` flag is set; then they see portions only.
3. **Other adults:** portions and extras (without kcal) only. Never another adult's kcal, target, deviation or
   status flags such as `short`.
4. **Child profiles:** everyone sees portions only. An adult account member of the household also sees the child's
   kcal and target when the child's `calories_hidden` is false (it is true by default for children; kitchen sets
   the default).
5. **Teen with own account** (member kind `child` with an account): sees their own portions only, never kcal, a
   goal or deviation, whatever their flag says.
6. Dish-level kcal per standard portion (library card, "about 520 kcal") is not personal and stays visible, except
   for a viewer who hides calories, and for teen accounts, where it is removed too.
7. One function `redact(response, viewer)` is applied at the RPC boundary of every response type that carries
   per-person numbers (portions, family day, week plan, day balance). Internal callers get unredacted data.
8. Wording: no response field or text key uses "diet", "goal" or "lose" for child profiles.

## Contract

No new RPC. All per-person kcal fields in planner messages are `optional` so they can be omitted.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Viewer not found in snapshot | `PERMISSION_DENIED` |

## Acceptance criteria

- [ ] Worked example viewed by Ewa: Ewa sees her 780 kcal and target; for Adam only "1½"; for Ola and Kuba only
      portions; Kuba's extras listed without kcal (Verification 19, 20).
- [ ] Viewed by Adam with Ola's `calories_hidden` = false: Adam sees Ola's 520 kcal.
- [ ] Kuba's own account: only "2 portions, a glass of milk and a banana"; no number anywhere in the payload.
- [ ] A test walks every planner response type and fails if a per-person kcal field reaches a viewer it shouldn't.

## Notes

- Decided (Q203): identity has no parent–child link; every adult account member of the household sees a child's
  calories and target when they are not hidden (rule 4).
