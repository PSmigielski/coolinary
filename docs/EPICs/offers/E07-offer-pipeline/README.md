# OFF-07 Offer extraction, validation and publishing

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | Turn each fetched flyer into checked offers: AI reads the pages once for everyone in a batch (50% cheaper), code without AI checks every offer, failures wait in a hold-back queue, and offers go live on their start date and disappear after their end date (features.md §3 steps 4–7, §2 "Flyers and offers … With AI"). |

## Flow

1. A fetched flyer version is sent to `ai` as one Message Batches job (`extract_flyer_offers`); urgent flyers use the
   per-request call instead.
2. Results are normalised by code: prices to minor units, pack sizes, unit prices, multi-buy and loyalty conditions.
3. Product names are matched to catalogue ingredients through kitchen (`MatchNames`).
4. Each offer is validated without AI; any failed rule holds it back for review.
5. Ready offers of live chains are published (shown as "next week" before they start), and expire after their end
   date. `offers.published` and `offers.expired` tell planner.

## Offer states

`extracted` → `ready` or `held_back` → `published` → `expired`; also `rejected` (admin or duplicate) and
`superseded` (a newer version of the flyer replaced it).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-07-01](OFF-07-01-offers-schema.md) | Offers schema rework | R2 | M | OFF-05-01, OFF-03-04 |
| [OFF-07-02](OFF-07-02-batch-extraction.md) | Submit flyers for extraction and collect results | R2 | L | OFF-07-01, OFF-05-05, OFF-01-04, AI-04, AI-13 |
| [OFF-07-03](OFF-07-03-normalise-offers.md) | Normalise prices, packs and conditions | R2 | L | OFF-07-02 |
| [OFF-07-04](OFF-07-04-match-products.md) | Match offer products to catalogue ingredients | R2 | M | OFF-07-03, KIT-03 (MatchNames), KIT-02 (GetIngredients) |
| [OFF-07-05](OFF-07-05-validation-rules.md) | Validation without AI | R2 | M | OFF-07-04, OFF-07-08 |
| [OFF-07-06](OFF-07-06-hold-back-queue.md) | Hold-back queue and admin review | R2 | M | OFF-07-05, OFF-01-06 |
| [OFF-07-07](OFF-07-07-publish-expire.md) | Publish, "next week" and expiry | R2 | L | OFF-07-05, OFF-01-01, OFF-01-05 |
| [OFF-07-08](OFF-07-08-duplicate-offers.md) | Duplicate and conflicting offers | R2 | S | OFF-07-03 |
