# AI-13 Flyer offer extraction

| | |
|---|---|
| Service | ai |
| Release | R2 |
| Goal | Read flyer pages into offers — product, size, price, conditions ("2+1", "with app or card"), valid dates — once a week for every covered chain (batch, once for everyone), and on request for any flyer a user photographs or shares (features.md §2 "Flyers and offers", §3 "Keeping flyers in sync" step 4). |

## Flow

```
chain flyers:  offers OFF-05/06 fetch pages -> SubmitBatch(extract_flyer_offers, one item per page)   (AI-13-03)
               -> GetBatchResults -> offers OFF-07 checks without AI (price, dates, unit, catalogue match) -> hold back or publish
user flyers:   phone photo / PDF / link -> gateway -> offers OFF-10 -> ExtractFlyerOffers (sync, chunked)  (AI-13-02)
both use the same offer schema and checks (AI-13-01); ai never fetches URLs
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-13-01](AI-13-01-offer-schema.md) | Offer schema and checks | R2 | M | AI-01-06, AI-06-01 |
| [AI-13-02](AI-13-02-extract-flyer-offers-sync.md) | ExtractFlyerOffers for user flyers (sync, chunked) | R2 | L | AI-13-01, AI-01-08 |
| [AI-13-03](AI-13-03-flyer-batch.md) | Batch kind extract_flyer_offers for chain flyers | R2 | S | AI-13-01, AI-04-03 |
| [AI-13-04](AI-13-04-flyer-golden-set.md) | Golden set per launch chain | R2 | M | AI-13-02, AI-13-03, AI-03-03 |
