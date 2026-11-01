# OFF-10-03 Extract and check a draft

| | |
|---|---|
| Service | offers |
| Epic | [OFF-10 User-supplied flyers](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-10-01, OFF-07-03, OFF-07-04, OFF-07-05, AI-13 (ExtractFlyerOffers) |
| Unblocks | OFF-10-04 |

## Goal

The uploaded flyer becomes a list of offers the user can check, read by AI and checked by the same code as chain
flyers (README "AI output is a proposal").

## Business rules

1. offers calls `ExtractFlyerOffers` (per request, AI-13) with the pages or PDF, the language of the shop's country
   (`pl`, `de`; for other countries the household's language from `x-language`), the shop's chain name if known, and
   the currency of the shop's country. Deadline 120 s.
2. The AI result also proposes the flyer's chain name and its start and end date when printed.
3. Each offer is normalised (OFF-07-03), de-duplicated within the upload (OFF-07-08 rule 2), matched (OFF-07-04) and
   validated (OFF-07-05) with these differences:
   - dates are checked against the dates the user will confirm, so `dates_outside_flyer` is reported as a warning;
   - `no_match` is a warning, not a block: unmatched offers can be saved but get no list badges;
   - `implausible_price` uses the reference price when one exists for the country.
4. The draft holds every offer with its warnings, and the proposed dates. It is stored in `user_flyers.draft`
   (text and numbers only) and the state becomes `draft`.
5. If the shop's edition already has a published flyer valid on those dates, the draft sets `already_covered = true`
   (the app can say so); the user may still save.
6. AI unavailable, over its spend limit, or refusing → state `failed` with `ai_unavailable`; zero offers found →
   `failed` with `no_offers_found`. The user can upload again; failed jobs do not count towards the daily limit.
7. A draft not confirmed within 24 h becomes `expired` and its draft data is deleted.

## Contract

```proto
message UserFlyerDraft {
  repeated DraftOffer offers = 1;
  string proposed_valid_from = 2;
  string proposed_valid_to = 3;
  string detected_chain_name = 4;
  bool already_covered = 5;
}
message DraftOffer {
  string draft_id = 1;            // stable within the draft
  Offer offer = 2;
  repeated string warnings = 3;   // validation reasons from OFF-07-05
}
```

## Data

No schema change (the draft is stored in `user_flyers.draft`).

## Errors

| Case | Result |
|---|---|
| ai `UNAVAILABLE`, `RESOURCE_EXHAUSTED`, refusal | job `failed` (`ai_unavailable`) |
| Nothing extracted | job `failed` (`no_offers_found`) |

## Acceptance criteria

- [ ] With the stub provider and a fixture, a draft with offers, warnings and proposed dates is produced.
- [ ] With ai stopped, the job fails with `ai_unavailable` within the deadline and the daily count is unchanged.
- [ ] An unmatched offer appears in the draft with warning `no_match`.
