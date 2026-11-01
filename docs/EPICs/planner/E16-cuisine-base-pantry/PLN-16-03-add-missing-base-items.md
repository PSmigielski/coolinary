# PLN-16-03 Add missing base items to the list

| | |
|---|---|
| Service | planner |
| Epic | [PLN-16 Cuisine base pantry](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | PLN-16-02, PLN-15-04, PLN-15-08 |
| Unblocks | GW-08 |

## Goal

One tap adds the missing base items of a cuisine to the shopping list (features.md §2 "Cuisine choice").

## Business rules

1. Adds every item that is missing and not already on the list; items already on the list are left as they are.
2. Amount per item: one pack, the smallest pack size known (PLN-15-04 rules); without a pack size, 1 piece or the
   rounding rule's minimum.
3. Source kind `base_pantry`, ref = cuisine code ("Italian base" in the app).
4. Goes straight onto the list (the user chose it in one tap); the response lists what was added.
5. One list change, one notice, idempotent by `client_mutation_id` (PLN-15-08).

## Contract

`contracts/coolinary/planner/v1/pantry.proto`

```proto
rpc AddBasePantryToList(AddBasePantryToListRequest) returns (AddBasePantryToListResponse);
message AddBasePantryToListRequest { string household_id = 1; string member_id = 2; coolinary.common.v1.Cuisine cuisine = 3; string client_mutation_id = 4; }
message AddBasePantryToListResponse { repeated string added_line_ids = 1; int64 list_version = 2; }
```

## Data

No schema change.

## Errors

Reason codes travel as trailing gRPC metadata `x-error-code` (README decision 18).

| Case | gRPC status / behaviour |
|---|---|
| `cuisine` unspecified or `client_mutation_id` missing | `INVALID_ARGUMENT` |
| No base pantry defined for that cuisine | `NOT_FOUND` |
| kitchen down | `UNAVAILABLE`; nothing added |
| Nothing missing | success with an empty `added_line_ids`; no version bump |
| Same `client_mutation_id` again | the first response is returned (PLN-15-08) |
| List version conflict while writing | retried once inside planner; then `ABORTED` (`version_mismatch`) |

## Acceptance criteria

- [ ] Asian with soy sauce and sesame oil missing: two lines added with source "asian", one pack each.
- [ ] Tapping twice adds nothing the second time.
