# OFF-10-04 Confirm and save private offers

| | |
|---|---|
| Service | offers |
| Epic | [OFF-10 User-supplied flyers](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-10-03, OFF-07-07 |
| Unblocks | PLN-17, PLN-18 (offers.published with household_id) |

## Goal

The user confirms what was read, and those offers then behave like any offer for that household: in browsing, list
badges and plan from deals, until they end.

## Business rules

1. `ConfirmUserFlyer`: the user sends the flyer's start and end date and, per draft offer, keep or drop, with optional
   corrections (product name, price, pack, conditions, dates).
2. Dates: required, end ≥ today, end ≥ start, span at most 31 days. Offer dates must lie inside them.
3. Kept offers are normalised and validated again. Offers that still fail a blocking rule (`no_price`, `no_unit`,
   `dates_reversed`, `dates_outside_flyer`, `bad_condition`) are returned as errors and nothing is saved; the user
   fixes or drops them.
4. Saved offers: source kind `user_flyer`, the household id, the shop id, state `published`. One `offers.published`
   event with the household id in the envelope.
5. Visibility (Decided, Q173): only the uploading household sees these offers; they are never shown to other
   households and never used as public price observations (OFF-11).
6. They expire with the normal expiry run (OFF-07-07). `DeleteUserFlyer` expires them at once with `offers.expired`.
7. After a Plus → Free downgrade, saved offers stay until they expire (OFF-04-04 rule 5).

## Contract

```proto
rpc ConfirmUserFlyer(ConfirmUserFlyerRequest) returns (UserFlyer);
message ConfirmUserFlyerRequest {
  string household_id = 1;
  string user_flyer_id = 2;
  string valid_from = 3;
  string valid_to = 4;
  repeated DraftDecision decisions = 5;
}
message DraftDecision {
  string draft_id = 1;
  bool keep = 2;
  Offer corrected = 3;            // only the editable fields are read
}
```

## Data

No schema change (offers rows with source kind `user_flyer`).

## Errors

| Case | gRPC status |
|---|---|
| Draft expired or already confirmed | `FAILED_PRECONDITION` |
| Blocking rule fails on a kept offer | `INVALID_ARGUMENT` with per-offer reasons |
| Another household's flyer | `NOT_FOUND` |

## Acceptance criteria

- [ ] Confirming 12 of 15 draft offers saves 12; they appear in `BrowseOffers` and `FindOffersForIngredients` for the
      household and for no other household.
- [ ] planner receives one `offers.published` with the household id and 12 offer ids.
- [ ] Deleting the flyer makes them disappear and sends `offers.expired`.
