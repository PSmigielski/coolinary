# OFF-08-01 Review samples and verdicts

| | |
|---|---|
| Service | offers |
| Epic | [OFF-08 Accuracy gate and weekly spot-checks](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-07-05, OFF-05-05, OFF-01-06 |
| Unblocks | OFF-08-02, OFF-08-03, OFF-08-04, ADM-06 |

## Goal

The founder checks a random sample of a chain's offers against the flyer pages, field by field, and the result is a
number the go-live gate and the weekly check can use.

## Business rules

1. Sample kinds: `go_live` (100 offers) and `weekly` (20 offers).
2. Sampling: offers of the chain's extracted versions in the window (go-live: all versions from the last 21 days;
   weekly: versions published in the last 7 days), states `ready`, `published`, `held_back` and `rejected` with
   reason other than `duplicate`. Drawn at random, spread over flyers: at most `ceil(size / number of flyers)` offers
   from any one flyer; a go-live sample needs offers from at least 2 flyers. If fewer offers exist, the sample cannot
   be created (`FAILED_PRECONDITION`).
3. Held and rejected offers are in the sample on purpose: an offer the checks held back still counts as a reading
   error.
4. Each sampled offer is judged on five fields: product (name and size), price, regular price (when shown),
   conditions, dates. Each field is `right`, `wrong` or `not_applicable`. The reviewer may add a note.
5. An offer is correct when price and dates are right (features.md: "≥ 97% correct prices and dates"). Product and
   conditions errors are recorded and reported separately.
6. Missed offers: the reviewer may add products that are on the sampled pages but were not extracted; each counts as
   one incorrect offer in the score (the sample size grows by one).
7. Score = correct offers / (sampled offers + missed offers), with one decimal.
8. A sample is `open` until every item has a verdict, then `closed` with its score. Open samples older than 14 days
   are `abandoned`.
9. Pages of sampled offers must still be stored; the sample window keeps inside the page retention (end + 90 days, Q010).

## Contract

`admin.proto`

```proto
rpc CreateReviewSample(CreateReviewSampleRequest) returns (ReviewSample);   // chain_id, kind
rpc GetReviewSample(GetReviewSampleRequest) returns (ReviewSample);         // items with page refs
rpc RecordReviewVerdict(RecordReviewVerdictRequest) returns (ReviewSample); // offer_id, five field verdicts, note
rpc AddMissedOffer(AddMissedOfferRequest) returns (ReviewSample);           // version_id, page_no, note
rpc ListReviewSamples(ListReviewSamplesRequest) returns (ListReviewSamplesResponse);

message ReviewSample {
  string id = 1;
  string chain_id = 2;
  string kind = 3;          // "go_live", "weekly"
  string state = 4;         // "open", "closed", "abandoned"
  int32 size = 5;
  int32 judged = 6;
  double score_percent = 7; // set when closed
  repeated ReviewItem items = 8;
}
```

## Data

`offers/migrations/000NN_review_samples.sql`

```sql
CREATE TABLE review_samples (
    id          uuid PRIMARY KEY,
    chain_id    text NOT NULL REFERENCES chains (id),
    kind        text NOT NULL CHECK (kind IN ('go_live', 'weekly')),
    state       text NOT NULL DEFAULT 'open' CHECK (state IN ('open', 'closed', 'abandoned')),
    created_by  text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    closed_at   timestamptz,
    score_percent numeric(4,1)
);
CREATE TABLE review_items (
    sample_id  uuid NOT NULL REFERENCES review_samples (id) ON DELETE CASCADE,
    offer_id   text REFERENCES offers (id),
    missed_version_id uuid, missed_page_no int,   -- for missed offers
    product    text CHECK (product IN ('right', 'wrong', 'not_applicable')),
    price      text CHECK (price IN ('right', 'wrong', 'not_applicable')),
    regular_price text CHECK (regular_price IN ('right', 'wrong', 'not_applicable')),
    conditions text CHECK (conditions IN ('right', 'wrong', 'not_applicable')),
    dates      text CHECK (dates IN ('right', 'wrong', 'not_applicable')),
    note       text NOT NULL DEFAULT '',
    judged_by  text,
    judged_at  timestamptz
);
```

## Errors

| Case | gRPC status |
|---|---|
| Too few offers for the sample | `FAILED_PRECONDITION` with the count found |
| Verdict on a closed sample | `FAILED_PRECONDITION` |

## Acceptance criteria

- [ ] A go-live sample from 3 fixture flyers holds 100 offers with at most 34 from one flyer.
- [ ] 97 offers right on price and dates, 3 wrong → score 97.0.
- [ ] 97 right, 3 wrong, plus 1 missed offer → score 96.0 (97 / 101, rounded to one decimal).
