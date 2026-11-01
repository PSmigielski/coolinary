# PLN-12-02 Meal shares per person

| | |
|---|---|
| Service | planner |
| Epic | [PLN-12 Per-person portions](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | PLN-01, KIT-11 (GetPlanningProfile) |
| Unblocks | PLN-12-03, PLN-12-05 |

## Goal

Turn each eater's daily kcal target into a kcal share for each meal, using the household's meal shares
(features.md "How one meal is split", step 1).

## Business rules

1. Inputs from the profile snapshot (KIT-11 `GetPlanningProfile`): members with `daily_kcal_target` (or none),
   member kind (`adult` / `child`), the household's meal shares and which meal is the main hot one (`dinner` or
   `supper`), and the household portions setting.
2. Default shares: breakfast 25%, main hot meal 40%, the other of dinner/supper 25%, snacks and drinks 10%. The
   household may change them; kitchen validates that they sum to 100. Planner rejects a snapshot whose shares don't
   sum to 100 with `FAILED_PRECONDITION` (`profile_invalid`) rather than guessing.
3. `meal_target(person, meal) = daily_kcal_target × share(meal) / 100`, kept as a decimal (no rounding here).
4. The snacks share is never planned. The "meal part" of the day target is `daily_kcal_target × (100 − snacks share) / 100`.
5. A person without a target ("just show calories", or no goal set) has no meal target; they get 1 standard portion
   in PLN-12-03.
6. Free households: at most one member has a target (the one personal goal, features.md §6). If the snapshot shows
   more than one target on a Free household (e.g. right after a downgrade), only the target of the account member
   is used; the others are treated as having no target.
7. Unnamed eaters: if the household portions setting is larger than the number of members in the snapshot, the
   difference is added as unnamed eaters without a target (1 standard portion each), so a Free couple where one
   person has a goal still cooks for two.
8. A member marked away for a meal (PLN-12-01) gets no share for that meal.

## Contract

Internal package `planner/internal/portions` (`MealShares(snapshot, date, meal, attendance) []EaterShare`). No RPC.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Shares don't sum to 100 | `FAILED_PRECONDITION` `profile_invalid` returned by the calling RPC |
| Target ≤ 0 | treated as no target, warning logged with member id |

## Acceptance criteria

- [ ] Worked example: dinner is the main meal; Adam 2,000 → 800, Ewa 1,800 → 720, Ola 1,400 → 560, Kuba 3,200 →
      1,280 kcal dinner shares.
- [ ] Changing the shares to 30/35/25/10 changes the dinner share of Adam to 700.
- [ ] Main meal set to supper: supper gets 40% and dinner 25%.
- [ ] Household portions 2, one member with a 1,800 target on Free: two eaters, one with a 720 share, one unnamed.
- [ ] Table-driven unit tests; no I/O in the package.

## Notes

- Decided (Q197): extra eaters beyond the household's people count as unnamed eaters with 1 standard portion each
  (rule 7).
- Meal shares by country (features.md §5 open point) are kitchen's defaults; planner only reads the household's values.
