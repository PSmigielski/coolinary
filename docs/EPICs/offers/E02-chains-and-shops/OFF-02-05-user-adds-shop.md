# OFF-02-05 User adds a missing shop

| | |
|---|---|
| Service | offers |
| Epic | [OFF-02 Chains and shops](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-02-02 |
| Unblocks | OFF-02-07, OFF-10-05 |

## Goal

When OSM lacks a shop, the user adds it once and can use it straight away; after review everyone nearby sees it
(features.md §3 "Users can add or correct a shop").

## Business rules

1. Any account (Free or Plus, anonymous included) may add a shop: chain (one of the seeded chains, or "other chain"
   with a free-text name up to 60 characters), latitude and longitude chosen on the phone's map, and optionally street
   and postcode.
2. The phone sends the shop's position, not the user's. offers does not ask where the user is.
3. The new shop has `source = 'user'`, `status = 'pending'`, `added_by_household` = the caller's household. It is
   returned to that household in shop lists and can be put in its my shops at once; other households do not see it.
4. A household may have at most 5 pending shops. The 6th is refused until one is reviewed.
5. If an active shop of the same chain is within 50 m, the call does not create a shop and returns the existing one
   (`existing_shop`), so the user can pick it instead.
6. Postcode: when not given, it is filled from the nearest postcode centroid (OFF-03-01).
7. "Other chain" shops never get synced flyers; they can only receive user-supplied flyers (OFF-10).
8. Text fields are trimmed, at most 120 characters, and stored as given; they are not logged.
9. On approval (OFF-02-07) the shop becomes `active` for everyone. On rejection (Decided, Q145): the shop becomes
   `rejected`, is hidden for everyone, and the adding household's my-shops entry is removed.

## Contract

`contracts/coolinary/offers/v1/shops.proto`

```proto
rpc SubmitShop(SubmitShopRequest) returns (SubmitShopResponse);

message SubmitShopRequest {
  string household_id = 1;       // from the JWT (README D3)
  string chain_id = 2;           // empty when other_chain_name is set
  string other_chain_name = 3;
  double latitude = 4;
  double longitude = 5;
  string street = 6;
  string postcode = 7;
  string country_code = 8;
}

message SubmitShopResponse {
  Shop shop = 1;                 // the new pending shop, or the existing one
  bool existing_shop = 2;
}
```

"Other chain" needs a chain row: a per-country placeholder chain `other-<cc>` (coverage `user_flyers_only`) is
seeded, and `other_chain_name` is stored on the shop as its `name`.

## Data

No new table; uses `shops` columns from OFF-02-02. Migration seeds chains `other-pl` and `other-de`.

## Errors

| Case | gRPC status |
|---|---|
| Position outside the country | `INVALID_ARGUMENT` |
| Neither or both of `chain_id` and `other_chain_name` | `INVALID_ARGUMENT` |
| 5 pending shops already | `FAILED_PRECONDITION` "too many shops waiting for review" |

## Acceptance criteria

- [ ] A user-added Biedronka appears in that household's shop list and my shops immediately, and in another
      household's list only after approval.
- [ ] Adding a Lidl 30 m from an existing Lidl returns the existing shop with `existing_shop = true`.
- [ ] The 6th pending shop is refused.
