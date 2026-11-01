# KIT-11-05 Calories hidden and "eat lighter as a family"

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-11 Per-person profiles](README.md) |
| Release | R2 |
| Size | S |
| Depends on | KIT-11-02, KIT-10-06 (`eat_lighter_as_family` column) |
| Unblocks | KIT-11-06, PLN-13-04 (eat lighter), GW-05 |

## Goal

Each person can hide calories; children never see calories or goals; and when a parent chooses "eat lighter as a
family", the whole family's extras change, not only the child's (features.md §2 "Calories and goals",
"Children and teenagers").

## Business rules

1. **`calories_hidden`** per person. Default false for adults. Setting `target_kind = child` sets it to true in the
   same update unless the request sets it explicitly. Who may change it: KIT-11-02 rule 5 (a locked child cannot).
2. The flag is a display rule for the app and the gateway: the app shows that person only portions. Planner still
   gets the target (it needs it for portions).
3. **`eat_lighter_as_family`** is household-level (in the `household` section, KIT-10-06). Any adult account
   member may switch it; children may not (`not_allowed`). Plus only (family goals are Plus): on Free →
   `plus_required`.
4. The app offers it when a child has the goal "grow into a healthier weight" or "set by our doctor or dietitian";
   the goal is on the phone, so kitchen does not check that condition and accepts the flag from any adult.
5. Each change bumps `planning_version` and writes `profile.changed` (`members` or `settings`).

## Contract

Field `calories_hidden` of `UpdateMemberProfile` (KIT-11-02); field `household.eat_lighter_as_family` of
`UpdateProfile` (KIT-10-01).

## Data

No schema change.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Reason |
|---|---|---|
| Child switching eat-lighter or their own hidden flag when locked | `PERMISSION_DENIED` | `not_allowed` / `goal_locked` |
| Free household switching eat-lighter on | `FAILED_PRECONDITION` | `plus_required` |

## Acceptance criteria

- [ ] Making Ola a child sets `calories_hidden = true`; an adult can switch it off for her.
- [ ] Kuba (locked child) cannot switch his hidden flag off.
- [ ] Turning eat-lighter on changes `GetPlanningProfile.household.eat_lighter_as_family` for the whole household
      (features.md Verification 19, last clause, is planner's check on top of this).
