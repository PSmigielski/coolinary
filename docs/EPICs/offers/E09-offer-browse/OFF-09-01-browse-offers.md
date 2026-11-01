# OFF-09-01 Browse offers of my shops

| | |
|---|---|
| Service | offers |
| Epic | [OFF-09 Offer browse, search and lookup by ingredient](README.md) |
| Release | R2 |
| Size | L |
| Status | todo |
| Depends on | OFF-07-07, OFF-04-01, OFF-03-04 |
| Unblocks | OFF-09-02, OFF-09-03, OFF-09-04, GW-10 |

## Goal

The offers screen shows the current (or next week's) offers of the user's shops, nearest shop first, filterable by
shop, category and end date (features.md §2 "Browse and search this week's offers by shop, category, and end date";
"Nearest shops first").

## Business rules

1. Input: shop ids in nearest-first order (empty = the household's active my shops in their stored order), view
   `current` (valid today), `upcoming` (start after today: "next week"), or `all`; optional category; optional
   `ending_by` date (offers that end on or before it, for "ends soon"); paging.
2. Free (plan from identity, `PlanFor`, OFF-01-03): shop ids must be among the household's active my shops (at most 2). Any other id →
   `FAILED_PRECONDITION` (`plus_required`). Plus: any active shops, at most 40 per call (the phone sends the shops
   inside its travel range).
3. Each shop maps to its editions (`EditionsForShop`, OFF-03-04). The result is the visible offers (OFF-07-07 rule 3)
   of those editions, plus the household's own user-flyer offers for those shops (OFF-10-04).
4. One offer appears once, with `shop_ids` = the requested shops that carry it, in request order.
5. Order: by the position of the offer's first shop in the request (nearest first), then category, then
   `discount_percent` descending, then product name. The server never computes a distance.
6. Shops of chains that are not live are not an error; they just add no offers. The response lists, per requested
   shop, its flyer status code (OFF-05-06) so the app can show "scan it yourself" next to it.
7. Page size default 50, maximum 100; the page token is opaque (encodes the sort key of the last item).
8. Results are cached in `redis-cache` per (edition set, view, category, ending_by, page) for 10 minutes; the key
   includes each edition's cache version (bumped by OFF-07-07). Household-private offers are merged after the cache
   read.
9. Product names are shown as printed (flyer language). Decided (Q170): the response also carries the matched
   ingredient's id so the app shows its name in the user's language underneath.
10. `ListOffers` is marked deprecated. The gateway switches to `BrowseOffers` (GW-10); `ListOffers` is then removed
    in a later release, with the breaking change accepted explicitly in `make proto-breaking`.

## Contract

`contracts/coolinary/offers/v1/offers.proto`

```proto
rpc BrowseOffers(BrowseOffersRequest) returns (BrowseOffersResponse);
rpc ListOffers(ListOffersRequest) returns (ListOffersResponse) { option deprecated = true; }

message BrowseOffersRequest {
  string household_id = 1;
  repeated string shop_ids = 2;        // nearest first; empty = my shops
  OfferView view = 3;
  string category = 4;
  string ending_by = 5;                // YYYY-MM-DD
  string query = 6;                    // OFF-09-02
  coolinary.common.v1.PageRequest page = 7;
}
enum OfferView { OFFER_VIEW_UNSPECIFIED = 0; OFFER_VIEW_CURRENT = 1; OFFER_VIEW_UPCOMING = 2; OFFER_VIEW_ALL = 3; }

message BrowseOffersResponse {
  repeated Offer offers = 1;
  repeated ShopFlyerStatus shop_statuses = 2;
  coolinary.common.v1.PageResponse page = 3;
  Attribution attribution = 4;
}
```

## Data

No schema change (index `offers_visible`).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Free and a shop outside active my shops | `FAILED_PRECONDITION` (`plus_required`) |
| More than 40 shops | `INVALID_ARGUMENT` |
| Unknown shop id | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] From a fixed test position, with shops sent nearest first, offers appear in order of shop distance
      (Verification 9).
- [ ] A Free household with Lidl and Biedronka sees only their offers; adding a third shop id is refused.
- [ ] A national Lidl offer requested for 3 Lidl shops appears once with 3 shop ids.
- [ ] `view = upcoming` returns next week's Lidl DE flyer offers before they start.
- [ ] p95 latency under 150 ms for 10 shops on the dev database with 20,000 published offers.
