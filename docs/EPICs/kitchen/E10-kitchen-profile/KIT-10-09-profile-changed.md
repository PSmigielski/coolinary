# KIT-10-09 profile.changed event and planning version

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-10-01, KIT-01 (outbox) |
| Unblocks | PLN-01-03 (drop the cached snapshot), KIT-10-08 |

## Goal

Whenever something planning depends on changes, planner learns it at once and fetches a new snapshot; nothing is
planned on an old allergen list.

## Business rules

1. `planning_version` (column of `kitchen_profiles`) increases by 1 in every transaction that changes anything in
   the snapshot: any profile section, a member row or its target or restrictions (KIT-11), a plan change seen as
   `household.plan_changed` (KIT-13-06), an ingredient merge that touched the household's staples or dislikes (KIT-13-09).
2. The same transaction writes one `profile.changed` to the outbox with `household_id`, `planning_version` and
   `sections` (one or more of `settings`, `cuisines`, `tools`, `diet`, `staples`, `members`, `plan`).
3. `version` (the ETag of `GetProfile`) changes only on `UpdateProfile`; member and plan changes do not make a
   pending household edit fail.
4. Events carry no values (no targets, no allergen lists), only what changed.
5. One event per transaction, however many sections changed.

## Contract

Event `profile.changed` (docs/EPICs/README.md), payload extended with `planning_version`.

```json
{ "type": "profile.changed", "producer": "kitchen", "household_id": "…",
  "payload": { "household_id": "…", "planning_version": 12, "sections": ["tools", "diet"] } }
```

## Data

No schema change (`planning_version` column from KIT-10-01).

## Errors

None of its own; a failed outbox insert fails the whole update.

## Acceptance criteria

- [ ] Saving tools and diet in one call writes one event with sections `["tools","diet"]` and `planning_version` + 1.
- [ ] A plan change (KIT-13-06) bumps `planning_version` but not `version`; a pending `UpdateProfile` with the old
      `version` still succeeds.
- [ ] The event payload never contains values (test on a household with a target and allergens).
