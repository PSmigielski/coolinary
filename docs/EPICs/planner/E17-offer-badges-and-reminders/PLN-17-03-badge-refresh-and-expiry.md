# PLN-17-03 Badge refresh and expiry

| | |
|---|---|
| Service | planner |
| Epic | [PLN-17 Offer badges and offer-ends reminders](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | PLN-17-01, PLN-21-05 |
| Unblocks | Verification 6 |

## Goal

Offers disappear on their end date, and so do their badges; new flyers update badges (features.md §3 "Keeping
flyers in sync" step 7, Verification 6).

## Business rules

1. `offers.expired` (handled in PLN-21-05) deletes badges with those offer ids at once; pinned or not.
2. A nightly sweep at 00:05 in each household's time zone deletes badges whose `valid_to` is before today, in case
   an event was missed.
3. `offers.published` (PLN-21-05) marks every list with open ingredient lines `badges_stale`.
4. Stale lists are recomputed: lazily on the next `GetList` (in the background; the stale badges are returned
   meanwhile), and by a sweep that recomputes up to 20 lists per second, oldest first.
5. Upcoming badges become current on their start date at the next recompute.
6. Removing a badge never changes the line itself.
7. Each badge change that is visible bumps the list version and sends a live notice (PLN-15-09); sweeps batch
   changes per list into one version bump.

## Contract

No new RPC.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| An offers event handler fails | rolled back and retried by the consumer framework; after 5 failures it goes to `dead_events` (PLN-01-01) |
| offers down during a badge refresh | badges keep their last state; the refresh retries at the next run; expired badges are still removed by end date |

## Acceptance criteria

- [ ] An offer ending Saturday: its badge is gone on Sunday morning, and the line is still on the list
      (Verification 6).
- [ ] After `offers.published` for Lidl, a list with chicken gets the new Lidl badge within one sweep cycle.
