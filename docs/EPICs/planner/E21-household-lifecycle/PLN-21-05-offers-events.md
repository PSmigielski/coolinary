# PLN-21-05 Offers events

| | |
|---|---|
| Service | planner |
| Epic | [PLN-21 Household lifecycle](README.md) |
| Release | R2 |
| Size | S |
| Depends on | PLN-01, PLN-17-01 |
| Unblocks | PLN-17-03 |

## Goal

Keep badges and deal references in step with the offers service (`offers.published`, `offers.expired`).

## Business rules

1. `offers.published`: set `badges_stale` on every list that has open ingredient lines (PLN-17-03 recomputes), and
   drop cached deal candidates (PLN-18-01) by bumping an offers epoch in redis-state.
2. `offers.expired`: delete badges with those offer ids (pinned too); plan meal offer items keep the reference but are
   shown as ended (PLN-18-04 `deal_ended`).
3. Both handlers are idempotent and batch by list (one version bump per list).

## Contract

Event consumers only.

## Data

No schema change.

## Errors

Handler errors roll back and are retried by the consumer framework; after 5 failures the event goes to `dead_events` (PLN-01-01). Offer ids planner doesn't know are ignored.

## Acceptance criteria

- [ ] `offers.expired` for an offer on three lists removes three badges with one version bump per list.
- [ ] `offers.published` followed by `GetList` returns the old badges and, within one sweep cycle, the new ones.

## Notes

- If the `offers.published` payload carried `ingredient_ids`, only lists with those ingredients would need a
  recompute. Suggested to the offers writer; not required.
