# PLN-20-04 Applying a top-up

| | |
|---|---|
| Service | planner |
| Epic | [PLN-20 Cheapest top-up](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | PLN-20-02, PLN-14-01, PLN-15-10, PLN-12-06 |
| Unblocks | GW-10 |

## Goal

The user picks which swaps to take; the week plan's dishes and the shopping list change accordingly.

## Business rules

1. Input: top-up id and the swap keys the user accepts (any subset). Expired top-ups → `FAILED_PRECONDITION`
   (`topup_expired`).
2. Each accepted swap creates a household-private variant of the dish with the swapped ingredient (PLN-02 private
   dish, `source = topup`, linked to the original), checked again with PLN-07; if the check now fails (the profile
   changed), that swap is skipped and reported.
3. The planned meals using the dish point to the variant; portions are recomputed (PLN-12-06), which updates list
   lines through PLN-15-10.
4. Lines for the basket's offers get a `topup` source and the offer badge pinned (PLN-17-01).
5. Cooked meals are never changed.
6. Everything is one plan version bump and one list change.

## Contract

`contracts/coolinary/planner/v1/topup.proto`

```proto
rpc ApplyTopUp(ApplyTopUpRequest) returns (ApplyTopUpResponse);
message ApplyTopUpRequest { string household_id = 1; string member_id = 2; string topup_id = 3; repeated string swap_keys = 4; string client_mutation_id = 5; }
message ApplyTopUpResponse { repeated string applied_swap_keys = 1; repeated string skipped_swap_keys = 2; int64 list_version = 3; int64 week_plan_version = 4; }
```

## Data

No schema change (PLN-02 private dishes; PLN-14-01 plan meals).

## Errors

| Case | gRPC status |
|---|---|
| Top-up expired | `FAILED_PRECONDITION` `topup_expired` |
| Week plan changed since the top-up | the swap is skipped if its meals no longer use the dish |

## Acceptance criteria

- [ ] Accepting "chicken breast → thighs" changes Monday's and Tuesday's curry (cooked once, eaten twice) and the
      list shows thighs with the pinned offer, breast removed.
- [ ] Accepting none of the swaps only pins the offers on the existing lines.
