# OFF-07-04 Match offer products to catalogue ingredients

| | |
|---|---|
| Service | offers |
| Epic | [OFF-07 Offer extraction, validation and publishing](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-07-03, KIT-03 (MatchNames), KIT-02 (GetIngredients) |
| Unblocks | OFF-07-05, OFF-09-05, OFF-11-01 |

## Goal

Each offer points at a catalogue ingredient, so it can be matched to shopping-list lines and recipes ("Hähnchenschenkel
1 kg" → chicken thigh). Flyers carry no barcodes, so matching is by name (features.md §4 table, "Flyers: No EANs").

## Business rules

1. Names sent for matching: the product name without pack text and without the brand when the brand is a separate
   field, in the flyer's language. Context `offer`.
2. Calls go in chunks of at most 200 names per `MatchNames` call (README D1 edge `offers → kitchen`). kitchen stores
   each normalised name's match once for everyone (README decision 13), so repeat products cost nothing.
3. Result per name: an ingredient id with a confidence, or "none". The ingredient id and confidence are stored on the
   offer. "none" → held reason `no_match` (OFF-07-05).
4. The offer's `category` is the matched ingredient's category from `GetIngredients` (one call per chunk of ids).
   Decided (Q156): offer categories are the catalogue's ingredient categories, one taxonomy across the app.
5. kitchen down or slow (`UNAVAILABLE`, `DEADLINE_EXCEEDED`): offers stay `extracted` and the step retries every
   30 minutes. An outage never produces `no_match`.
6. Offers held back as `no_match` are re-sent once a night (02:30 local), because kitchen's review queue may have
   added the match meanwhile (KIT-03); a new match removes the reason and re-runs validation.
7. When `catalogue.ingredient_merged` arrives, ingredient ids are rewritten (OFF-14-02).

## Contract

Uses `KitchenService.MatchNames` and `GetIngredients`.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| kitchen down | Offers stay `extracted`; retry every 30 minutes |
| kitchen returns an unknown ingredient id | Treated as "none" and logged with the id |

## Acceptance criteria

- [ ] A fixture flyer with 150 offers needs one `MatchNames` call and one `GetIngredients` call.
- [ ] With kitchen stopped, no offer gets `no_match`; after kitchen returns, matching completes within 30 minutes.
- [ ] A held `no_match` offer whose name kitchen learns during the day is matched and validated the next night.
