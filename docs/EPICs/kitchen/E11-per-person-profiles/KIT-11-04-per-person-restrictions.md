# KIT-11-04 Per-person diet, allergens and dislikes

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-11 Per-person profiles](README.md) |
| Release | R2 |
| Size | M |
| Depends on | KIT-11-02, KIT-10-04, KIT-13-06 (plan-change handler), KIT-01-06 (plan from identity) |
| Unblocks | KIT-11-06, PLN-13-01 (hard limits for a shared meal) |

## Goal

Each person can add their own diet, allergens and dislikes, so a shared dish meets the limits of everyone eating it
and nobody else's meals are restricted needlessly (features.md §2 "One dish per meal": "The shared dish meets the
allergens and diet of everyone eating it").

## Business rules

1. A person's section has the same shape and limits as the household's (KIT-10-04): diet flags, the 14 allergens,
   0–50 dislikes. Vegan implies vegetarian.
2. **Meaning:** the household section applies to every meal; a person's section applies to the meals that person
   eats. Planner builds the union for the eaters of a meal (PLN-01-03 `HardLimits`); kitchen does not merge them.
3. **Plus only to edit** (features.md §6: per-person diet profiles are Plus): on Free → `plus_required`.
4. Decided (Q195) **after a downgrade to Free** (`household.plan_changed`, run by KIT-13-06):
   1. every person's allergens are merged into the household allergen list in one transaction (a union; nothing is
      removed), so nobody's allergy is dropped. Decided (Q230): the health-protecting diet flags gluten-free,
      lactose-free and nut-free are merged the same way, as household diet flags;
   2. per-person lifestyle diets (vegetarian, vegan) and dislikes stop applying: `GetPlanningProfile` leaves them
      out while the household is Free (KIT-10-08);
   3. nothing is deleted: the per-person rows stay and apply again after an upgrade; the merged household list is
      not undone on upgrade (the owner can edit it);
   4. per-person data can't be edited until the household is Plus again; removing the person (identity) removes
      their rows.
5. Who may edit: KIT-11-02 rule 5 (a locked child may change only their own dislikes).
6. The section is replaced as a whole on save; rows live in `profile_restrictions` / `profile_dislikes` with the
   member id (KIT-10-04 schema).
7. These values are visible to every member of the household (they decide what the shared dish may contain); targets
   are not (KIT-11-06).

## Contract

Field `diet` of `UpdateMemberProfile` (KIT-11-02), message `DietSection` (KIT-10-04).

## Data

No schema change (tables from KIT-10-04 with the member foreign key from KIT-11-01).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status | Reason |
|---|---|---|
| Free household | `FAILED_PRECONDITION` | `plus_required` |
| Unknown code or ingredient, more than 50 dislikes | `INVALID_ARGUMENT` | `invalid_value` / `unknown_ingredient` |

## Acceptance criteria

- [ ] Household diet: none; Ola: vegetarian; Kuba: nut-free. The snapshot lists both per person and none at household
      level.
- [ ] After a downgrade, the household allergens include nuts (from Kuba) and Kuba's row is unchanged; Ola's
      vegetarian is not in the snapshot; editing either returns `plus_required`.
- [ ] After a new upgrade, Ola's vegetarian applies again and the household allergen list still includes nuts.
- [ ] A household with no per-person allergens is unchanged by the downgrade.
- [ ] Saving a person's section with `[vegan]` returns `[vegetarian, vegan]`.
