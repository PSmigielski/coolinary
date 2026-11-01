# OFF-02-07 Shop moderation admin RPCs

| | |
|---|---|
| Service | offers |
| Epic | [OFF-02 Chains and shops](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-02-05, OFF-02-06, OFF-01-06 |
| Unblocks | ADM-06 (shop review screen) |

## Goal

The founder reviews user-added shops and corrections in one queue, oldest first, and each decision takes effect
immediately.

## Business rules

1. The queue lists pending user-added shops and open corrections, oldest first, with the shop, the nearest 3 active
   shops of any chain (id, chain, distance in metres from the shop under review), and the count of households that
   use the shop in my shops. No household ids are shown.
2. Approve a user-added shop: status `active`, then OFF-02-04 merge check runs for it.
3. Reject a user-added shop: follows the rule chosen in OFF-02-05 rule 9.
4. Approve a correction: apply to the shop (and `shop_overrides` for OSM shops); `closed` sets status `closed`.
   Other open corrections of the same kind for that shop become `replaced`.
5. Reject a correction: status `rejected`; for "closed" reports, `reported_closed` is cleared.
6. Admin may edit a pending shop's name, chain or position before approving.
7. Every decision writes `admin_changes`.

## Contract

`admin.proto`

```proto
rpc ListShopReviews(ListShopReviewsRequest) returns (ListShopReviewsResponse);
rpc ReviewShop(ReviewShopRequest) returns (ReviewShopResponse);           // approve / reject / edit-and-approve
rpc ReviewShopCorrection(ReviewShopCorrectionRequest) returns (ReviewShopCorrectionResponse);
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Item already decided | `FAILED_PRECONDITION` |
| Unknown id | `NOT_FOUND` |

## Acceptance criteria

- [ ] Approving a user-added shop 30 m from an OSM shop of the same chain merges it into the OSM shop.
- [ ] Approving "closed" removes the shop from shop lists; households that had it keep a closed entry (OFF-04-05).
