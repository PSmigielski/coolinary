# OFF-07-06 Hold-back queue and admin review

| | |
|---|---|
| Service | offers |
| Epic | [OFF-07 Offer extraction, validation and publishing](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-07-05, OFF-01-06 |
| Unblocks | ADM-06 (held offers screen) |

## Goal

The founder fixes or rejects held offers quickly, beside the flyer page they came from, and fixes feed back into
matching.

## Business rules

1. `ListHeldOffers`: filter by chain, edition, reason, and "starts within N days"; ordered by start date, then flyer,
   then page. Each item has the offer fields, reasons, the raw AI result, and a reference to the flyer page
   (`version_id`, `page_no`) for `GetFlyerPage`.
2. `UpdateHeldOffer`: admin may change product name, brand, pack, price, regular price, conditions text, dates and
   the ingredient id. Normalisation (OFF-07-03) and validation (OFF-07-05) run again. If no reason is left, the offer
   becomes `ready` and is published by the next publish run (OFF-07-07).
3. Setting an ingredient id by hand does not teach kitchen; the admin service also calls kitchen's admin RPC to store
   the name match for everyone (ADM-03 / KIT-03).
4. `ReleaseHeldOffer`: publish despite remaining reasons. Allowed only for `low_confidence` and `implausible_price`;
   other reasons must be fixed. The override is recorded.
5. `RejectHeldOffer`: state `rejected` with a reason text (kept for accuracy figures).
6. Bulk: `RejectHeldOffers` for all held offers of one version (a misread page).
7. Held offers whose end date passes without a decision become `expired`; they count as "not published" in the
   chain's accuracy record.
8. Every action writes `admin_changes`.

## Contract

`admin.proto`

```proto
rpc ListHeldOffers(ListHeldOffersRequest) returns (ListHeldOffersResponse);
rpc UpdateHeldOffer(UpdateHeldOfferRequest) returns (HeldOffer);
rpc ReleaseHeldOffer(ReleaseHeldOfferRequest) returns (HeldOffer);
rpc RejectHeldOffer(RejectHeldOfferRequest) returns (HeldOffer);
rpc RejectHeldOffers(RejectHeldOffersRequest) returns (RejectHeldOffersResponse);
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Release with a non-releasable reason | `FAILED_PRECONDITION` naming the reason |
| Offer not held | `FAILED_PRECONDITION` |

## Acceptance criteria

- [ ] Fixing the dates of a `dates_outside_flyer` offer makes it `ready` and it is published within 15 minutes.
- [ ] Releasing a `no_match` offer is refused.
- [ ] The queue item links to the right stored page.
