# OFF-07-05 Validation without AI

| | |
|---|---|
| Service | offers |
| Epic | [OFF-07 Offer extraction, validation and publishing](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-07-04, OFF-07-08 |
| Unblocks | OFF-07-06, OFF-07-07, OFF-10-03 |

## Goal

No wrongly read offer reaches users unchecked: every offer passes the same rule checks, and anything that fails is held
back (features.md §3 step 5; README "AI output is a proposal").

## Business rules

1. Rules, each adding a held reason when it fails:

   | Rule | Reason |
   |---|---|
   | Price present and > 0 | `no_price` |
   | Printed price parses and equals the AI price (OFF-07-03 rule 1) | `price_mismatch`, `price_unparsed` |
   | Valid from ≥ flyer's start and valid to ≤ flyer's end (offers without own dates take the flyer's) | `dates_outside_flyer` |
   | Valid from ≤ valid to | `dates_reversed` |
   | Pack quantity or loose basis present | `no_unit` |
   | Product matched to an ingredient (OFF-07-04) | `no_match` |
   | AI confidence ≥ 0.6 | `low_confidence` |
   | Multi-buy: pay < buy, total > 0 | `bad_condition` |
   | Unit price within 0.2× to 5× the country reference price of the ingredient, when one exists (OFF-11-03) | `implausible_price` |

   Decided (Q157): the confidence threshold is 0.6 and the plausibility band 0.2×–5× of the reference price, on
   top of the four checks in features.md.
2. An offer with no reason becomes `ready`; with any reason `held_back`, keeping all reasons.
3. Validation runs again whenever an offer changes (admin edit, new match), clearing old reasons first.
4. Offers dated outside the flyer are never shortened silently to fit; they wait for review.
5. Metrics: `offers_validated_total{chain,result}`, `offers_held_total{chain,reason}`. A version with more than 20%
   held offers sends an alert `offers.held.<source_id>` (a sign the layout or the reading changed).

## Contract

No gRPC change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Any rule fails | Never an error: the failure becomes a held reason |

## Acceptance criteria

- [ ] An offer valid 10–20 October in a flyer valid 12–17 October is held with `dates_outside_flyer`.
- [ ] An offer without pack or loose basis is held with `no_unit`.
- [ ] A version of 100 offers with 25 held sends one alert.
- [ ] An offer with three failing rules shows all three reasons.
