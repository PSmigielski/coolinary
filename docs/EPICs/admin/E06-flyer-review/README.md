# ADM-06 Flyer sources, offers and shops review

**Goal.** The human checks that features.md §3 and §8 require before and after a chain goes live: the retailer's
terms (a gate), a 100-offer accuracy sample at ≥ 97% correct prices and dates, a weekly 20-offer spot-check, the
queue of offers that failed the non-AI checks, user corrections to shops, linking shops to chain store numbers,
reference prices, and a status board that shows failing sources the same day. Removing a chain takes effect at
once. From R3 also the licensed feed's retailer and store mapping.

**Flow.** offers syncs flyers (OFF-05) → extraction (OFF-07) → held offers / samples (OFF-08) → admin reviews with the
flyer page image → offers updates source status and publishes.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [ADM-06-01](ADM-06-01-source-status-board.md) | Flyer source status board | R2 | M | ADM-01-05, OFF-05-08, OFF-08-01, OFF-08-04 |
| [ADM-06-02](ADM-06-02-terms-gate-and-status.md) | Terms gate and chain status | R2 | M | ADM-06-01, OFF-05-08, OFF-08-02, OFF-08-05 |
| [ADM-06-03](ADM-06-03-accuracy-sample.md) | Accuracy sample of 100 offers | R2 | L | ADM-06-02, OFF-05-08, OFF-08-01 |
| [ADM-06-04](ADM-06-04-weekly-spot-check.md) | Weekly spot-check of 20 offers per chain | R2 | S | ADM-06-03, OFF-08-01, OFF-08-03, OFF-08-04 |
| [ADM-06-05](ADM-06-05-held-offers-queue.md) | Held-back offers queue | R2 | M | ADM-06-03, OFF-07-06, KIT-03-06 |
| [ADM-06-06](ADM-06-06-shop-corrections.md) | Shop additions and corrections from users | R2 | S | ADM-02-03, OFF-02-04, OFF-02-07 |
| [ADM-06-07](ADM-06-07-store-linking.md) | Link shops to chain store numbers | R2 | S | ADM-02-03, OFF-03-05 |
| [ADM-06-08](ADM-06-08-reference-prices.md) | Reference prices | R2 | S | ADM-02-03, OFF-11-03 |
| [ADM-06-09](ADM-06-09-feed-mapping.md) | Licensed feed: retailer and store mapping | R3 | S | ADM-06-07, OFF-13-01 |
