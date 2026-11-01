# ADM-06-03 Accuracy sample of 100 offers

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R2 |
| Size | L |
| Status | todo |
| Depends on | ADM-06-02; OFF-08-01 (`CreateReviewSample`, `GetReviewSample`, `RecordReviewVerdict`, `AddMissedOffer`, `ListReviewSamples`), OFF-05-08 (`GetFlyerPage`) |
| Unblocks | ADM-06-04, ADM-06-02 (`go_live`) |

## Goal

Before a chain goes live, a human compares 100 extracted offers with the flyer pages; the chain may go live only at
≥ 97% correct prices and dates (features.md §8, verification 5).

## Business rules

1. Starting a go-live sample calls `CreateReviewSample(chain_id, kind = go_live)`: offers draws 100 offers at random
   from the chain's flyers of the last 21 days, at most `ceil(100 / number of flyers)` from one flyer and from at
   least 2 flyers, held-back and rejected offers included (OFF-08-01 rules 1–3). Too few offers → offers'
   `FAILED_PRECONDITION`, shown as 409.
2. `GetReviewSample` returns the items with their page references. For each item the reviewer sees the flyer page
   (proxied by admin, rule 7) and the extracted fields: product (name and size), price, regular price when shown,
   conditions ("2+1", "with app or card"), valid from and valid to.
3. Per item the reviewer judges five fields as `right`, `wrong` or `not_applicable` and may add a note:
   `RecordReviewVerdict(offer_id, product, price, regular_price, conditions, dates, note)`. Offers stores no
   corrected values; the note says what the right value is.
4. An item counts as correct when price and dates are right (OFF-08-01 rule 5, features.md "≥ 97% correct prices
   and dates"). Product and conditions errors are recorded and reported separately.
5. Missed offers: a product on a sampled page that was not extracted is added with `AddMissedOffer(version_id,
   page_no, note)` and counts as one incorrect offer (OFF-08-01 rule 6).
6. The sample closes by itself once every item has a verdict; score = correct ÷ (sampled + missed), one decimal.
   `go_live` (ADM-06-02) needs a closed go-live sample of the last 21 days with score ≥ 97.0. Open samples older
   than 14 days are abandoned by offers.
7. Flyer pages: `GET /admin/v1/offers/flyers/{flyerId}/pages/{n}` streams the bytes from offers (`GetFlyerPage`, up
   to 10 MB; admin's client raises its gRPC message limit to 16 MB for this call as offers does, OFF-05-08) with
   `Cache-Control: private, max-age=300`. A PDF page is returned as the PDF plus page number. Pages exist only as
   long as offers keeps them (until the offers expire plus 90 days, Q010); then 404 `page_gone`.
8. Several reviewers may work on one sample; each item keeps the last verdict; the audit log keeps all.
9. Audited: `offers.sample.start`, `offers.sample.verdict`, `offers.sample.missed_offer`.
10. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/offers/chains/{chainId}/samples` | `{ kind: "go_live" }` → sample |
| GET | `/admin/v1/offers/samples/{id}` | sample with items, verdicts, score when closed |
| POST | `/admin/v1/offers/samples/{id}/items/{offerId}/verdict` | `{ product, price, regular_price, conditions, dates, note? }` (each `right`\|`wrong`\|`not_applicable`) |
| POST | `/admin/v1/offers/samples/{id}/missed` | `{ version_id, page_no, note }` |
| GET | `/admin/v1/offers/flyers/{flyerId}/pages/{n}` | page bytes |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Page no longer stored | 404 `page_gone` |
| Too few offers, or offers from fewer than 2 flyers | 409 (offers' `FAILED_PRECONDITION`) |
| Verdict on a closed or abandoned sample | 409 |

## Acceptance criteria

- [ ] 97 correct of 100 with no missed offers → score 97.0, enough for `go_live`; 96 → not enough.
- [ ] An item with a wrong product name but right price and dates counts as correct and the product error is listed
      separately.
- [ ] 97 correct of 100 plus one missed offer → score 96.0 (97 ÷ 101).
- [ ] The page image loads through admin, not from the retailer's site.
