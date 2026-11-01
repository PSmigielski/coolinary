# MOB-05-03 Live stock notices

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-05 Live shopping list client](README.md) |
| Release | R2 |
| Size | S |
| Depends on | MOB-05-01, MOB-03-02, GW-09 |
| Unblocks | — |

## Goal

When someone else ticks an item into the fridge, cooks or edits the stock, the other phones update without a
manual refresh.

## Business rules

1. Notice (GW-09-02, KIT-07-08): `{ "kind": "stock.changed", "household_id", "version" }`. Only two kinds exist,
   `stock.changed` (here) and `list.changed` (MOB-05-02).
2. `stock.changed` with a `version` above the cached stock version → `invalidate("live stock.changed")` (stock,
   suggestions), debounced 2 s so a batch of ticks causes one refetch. A version at or below the cached one is
   ignored.
3. Household membership changes have no live notice. The phone learns about them when it refetches the household
   on foreground, and a removed member's next request gets 401 `token_stale` (GW-04-05, IDN-08-04), which forces a
   token refresh and a household refetch (MOB-01).
4. Unknown `kind` values are ignored (forward compatible).

## Contract

As in rule 1.

## Data

None.

## Errors

| Case | Behaviour |
|---|---|
| Notice for another household id | ignored (the socket may still be closing after a household change) |
| Malformed JSON or missing fields | ignored and counted in a dev-only log |
| Unknown `kind` | ignored (rule 4) |
| Refetch after a notice fails | standard query error state; the next notice or foreground triggers it again |

## Acceptance criteria

- [ ] Ticking 5 lines on phone A causes one stock refetch on phone B.
- [ ] A `stock.changed` notice with an old version causes no refetch.
- [ ] After phone A's owner removes account B, B's next request gets 401 `token_stale` and B shows its new household.
