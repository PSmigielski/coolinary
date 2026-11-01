# KIT-03-03 MatchNames: stored, exact and trigram paths

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-03 Name matching](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-03-02, KIT-02-04 |
| Unblocks | KIT-03-04, KIT-04-10, KIT-08, OFF-07, PLN-03 |

## Goal

`MatchNames` answers most names without AI: from the store, by exact name, or by a clear trigram winner.

## Business rules

1. Input: up to 100 names (`name`, optional `brand`), `language`, `context` (`product`, `receipt`, `offer`, `recipe`,
   `typed`, `photo`), `mode` (`INTERACTIVE` or `BULK`). More than 100 names is `INVALID_ARGUMENT`.
2. Each name is normalised (KIT-03-01); a normaliser error is returned for that name only (`outcome = error` with a
   reason). Names with the same key in one call are resolved once.
3. Stored path: a row with the current normaliser version and outcome `matched` or `none` is the answer. Rows with
   outcome `queued` answer `queued` in `BULK` mode; in `INTERACTIVE` mode they go on to step 7. Rows with outcome
   `provisional` go on to step 7 in both modes.
4. Exact path: the key equals the normalised name of exactly one active ingredient in the language or in `en` →
   `matched`, confidence 1.0, method `exact` (primary name) or `synonym`. Stored.
5. If the key equals names of two or more ingredients, those ingredients become the candidates for step 7.
6. Trigram path: the 10 best distinct active ingredients by `similarity(key, normalized)` ≥ 0.3 over names in the
   language and in `en`.
   - No candidate ≥ 0.3 → `none`, confidence 0, `review = needed`. Stored. No AI call.
   - Decided (Q057): best ≥ 0.85 and best − second ≥ 0.15 → `matched`, method `trigram`, confidence = best,
     no review. Stored. The thresholds are config values, tuned from review data after launch.
7. Otherwise the name goes to the AI path (KIT-03-04) in `INTERACTIVE` mode or the queue (KIT-03-05) in `BULK` mode.
8. Each answer returns: input index, key, outcome, ingredient id, confidence, method, review flag, up to 3
   candidates with scores (for the app's "did you mean"), and the extracted quantity.
9. `context` is added to the row's `contexts`. No household id is stored.
10. p95 ≤ 30 ms for 20 names that are all stored.

## Contract

`contracts/coolinary/kitchen/v1/matching.proto`

```proto
rpc MatchNames(MatchNamesRequest) returns (MatchNamesResponse);
message MatchNamesRequest { repeated NameInput names = 1; string language = 2; MatchContext context = 3; MatchMode mode = 4; }
message NameInput { string name = 1; string brand = 2; }
enum MatchContext { MATCH_CONTEXT_UNSPECIFIED = 0; MATCH_CONTEXT_PRODUCT = 1; MATCH_CONTEXT_RECEIPT = 2;
  MATCH_CONTEXT_OFFER = 3; MATCH_CONTEXT_RECIPE = 4; MATCH_CONTEXT_TYPED = 5; MATCH_CONTEXT_PHOTO = 6; }
enum MatchMode { MATCH_MODE_UNSPECIFIED = 0; MATCH_MODE_INTERACTIVE = 1; MATCH_MODE_BULK = 2; }
enum MatchOutcome { MATCH_OUTCOME_UNSPECIFIED = 0; MATCH_OUTCOME_MATCHED = 1; MATCH_OUTCOME_NONE = 2;
  MATCH_OUTCOME_QUEUED = 3; MATCH_OUTCOME_PROVISIONAL = 4; MATCH_OUTCOME_ERROR = 5; }
message NameMatch { int32 index = 1; string key = 2; MatchOutcome outcome = 3; string ingredient_id = 4;
  float confidence = 5; string method = 6; bool needs_review = 7; repeated Candidate candidates = 8;
  coolinary.common.v1.Quantity quantity = 9; string error = 10; }
message Candidate { string ingredient_id = 1; string name = 2; float score = 3; }
message MatchNamesResponse { repeated NameMatch matches = 1; }
```

## Data

No schema change (KIT-03-02).

## Errors

| Case | Status | Message |
|---|---|---|
| More than 100 names, or none | `INVALID_ARGUMENT` | `1-100 names` |
| One bad name | per-name `outcome = ERROR` | `empty_name`, `too_long`, `not_a_name` |

## Acceptance criteria

- [ ] "Zucchini" in `en` matches by exact name; "courgette" by synonym.
- [ ] "pepper" (two ingredients) never matches on the exact path.
- [ ] A name with no candidate ≥ 0.3 returns `none` and creates a review row without calling ai (checked with a
      stub that fails on any call).
- [ ] The second call for the same name reads the stored row.
