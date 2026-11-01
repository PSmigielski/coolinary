# KIT-05-05 Admin product review

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-05 Unknown products](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-05-02, KIT-01-07 |
| Unblocks | ADM-03 |

## Goal

The founder settles disputed and reported products and corrects shared products, so the shared catalogue stays
trustworthy.

## Business rules

1. `ListProductReviews(status, page)` lists `disputed` products, newest first; each with every fingerprint's values,
   the number of households holding it, first and last date, and report reasons with counts. No household ids.
2. `GetSharedProduct(code)` returns the shared row, its fingerprints and its change history.
3. `ResolveProduct(code, action)`: `confirm_fingerprint(id)` (those values become the confirmed product),
   `set_values(fields)` (admin values, `source = 'admin'`, confirmed), or `retire` (the code is no longer answered
   from level 2; households keep their own rows). Report counts are reset.
4. `EditSharedProduct(code, fields)` corrects a confirmed product directly.
5. Every action writes `admin_changes` and deletes the code's cache key.
6. Changing a product's `ingredient_id` here also updates the stored match for its name (KIT-03-06 `set_ingredient`)
   so similar names follow.

## Contract

`contracts/coolinary/kitchen/v1/admin.proto`

```proto
rpc ListProductReviews(ListProductReviewsRequest) returns (ListProductReviewsResponse);
rpc GetSharedProduct(GetSharedProductRequest) returns (SharedProductDetail);
rpc ResolveProduct(ResolveProductRequest) returns (SharedProductDetail);
rpc EditSharedProduct(EditSharedProductRequest) returns (SharedProductDetail);
```

## Data

No schema change.

## Errors

| Case | Status | Message |
|---|---|---|
| Unknown code | `NOT_FOUND` | `product not found` |
| Fingerprint of another code | `INVALID_ARGUMENT` | `fingerprint does not match code` |

## Acceptance criteria

- [ ] Confirming one fingerprint of a disputed product makes level 2 answer it with those values.
- [ ] Retiring a product leaves households' own names working.
- [ ] Each action appears in `ListAdminChanges` for the product.
