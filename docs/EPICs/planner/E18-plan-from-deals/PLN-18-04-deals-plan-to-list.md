# PLN-18-04 Deals plan to the shopping list

| | |
|---|---|
| Service | planner |
| Epic | [PLN-18 Plan from deals](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | PLN-18-02, PLN-15-05, PLN-17-01 |
| Unblocks | — |

## Goal

When a deals plan's list is generated, the offers the plan relied on are already on the lines.

## Business rules

1. `GenerateList` (PLN-15-05) for dates with a deals plan adds, for each needed ingredient that has a plan offer item,
   a source of kind `deal` (ref offer id) and pins that offer's badge (PLN-17-01 rule 7).
2. Pack rounding uses the offer's pack size first when it covers the need.
3. If the offer ended before the list is generated, the line keeps no badge and gets the flag `deal_ended`.

## Contract

No new RPC.

## Data

No schema change.

## Errors

None of its own; runs inside list generation (PLN-15-05). If offers is down, lines are generated without badges and get them on the next badge refresh (PLN-17-03).

## Acceptance criteria

- [ ] Accepting the list for a deals plan shows the rice line with the pinned Lidl offer.
- [ ] An ended offer leaves the line without a badge and flagged `deal_ended`.
