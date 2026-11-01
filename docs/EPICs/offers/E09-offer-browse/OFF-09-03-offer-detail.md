# OFF-09-03 Offer detail

| | |
|---|---|
| Service | offers |
| Epic | [OFF-09 Offer browse, search and lookup by ingredient](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-09-01 |
| Unblocks | GW-10 |

## Goal

Tapping an offer shows its full facts and a link to the original flyer (features.md §2 "Every offer links to the
original flyer"; Risks "Link to the source flyer on every offer").

## Business rules

1. `GetOffer(offer_id)` returns the offer with: chain, edition, flyer dates, page number, printed conditions and
   structured conditions, regular price and discount, unit price, matched ingredient id, and `shop_ids` among the
   household's active my shops that carry it.
2. `source_url` is the flyer's public URL on the chain's site (with the page anchor when the chain's viewer supports
   one). User-flyer offers have no URL; they show "from your photo of <date>".
3. A visible offer only (OFF-07-07 rule 3). An expired offer returns `NOT_FOUND`, so the app can drop it.
4. A user-flyer offer is visible only to its household; others get `NOT_FOUND`.

## Contract

```proto
rpc GetOffer(GetOfferRequest) returns (GetOfferResponse);
message GetOfferRequest { string household_id = 1; string offer_id = 2; }
message GetOfferResponse { Offer offer = 1; Chain chain = 2; Flyer flyer = 3; }
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Unknown, expired, or another household's private offer | `NOT_FOUND` |

## Acceptance criteria

- [ ] A Kaufland store-edition offer links to `leaflets.kaufland.com/…/ar/<store>`.
- [ ] Another household asking for a private user-flyer offer gets `NOT_FOUND`.
