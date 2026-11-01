# KIT-05-03 Read a product label photo

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-05 Unknown products](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-05-01, AI-07 (`ReadProductLabel`) |
| Unblocks | GW-06 (label route, counts as a photo scan), MOB-07 (uploads) |

## Goal

Instead of typing, the user photographs the label of an unknown product and gets a proposed name, brand and size to
confirm (features.md §4 step 5 "or photographs the label and AI reads it (per request)").

## Business rules

1. `ReadLabel(household_id, code, image, mime_type, language)`. The image is JPEG, PNG or WebP (the phone converts
   HEIC), at most 3 MB, long edge at most 1,568 px (larger images cost more and are resized by Claude anyway).
2. The image is held in memory for the call and never stored or logged (README D6).
3. kitchen calls ai `ReadProductLabel` with a 15 s deadline. The result is a proposal only; nothing is saved until
   the user calls `NameProduct` with `source = label` (KIT-05-01).
4. Checks without AI on the result: name 2–120 characters and passes `not_a_name`; brand ≤ 60; pack size parsed by
   the normaliser into a valid range; anything failing is dropped from the proposal. An empty proposal is
   `found = false`.
5. The proposal's name is matched with `MatchNames` (`INTERACTIVE`, context `photo`) and returned with up to 3
   candidate ingredients.
6. Decided (Q110): no nutrition values are read from labels in R1; only name, brand and size. A misread number
   would break "AI never states a calorie number" in spirit.
7. Quota: the gateway counts this call as one AI photo scan (README Free and Plus). Kitchen applies no quota.
8. ai `UNAVAILABLE`, `RESOURCE_EXHAUSTED` or a refusal → kitchen returns `UNAVAILABLE` with message
   `label reading unavailable`; the app falls back to typing.
9. `ai_used` is true on every successful answer (label reading has no non-AI path). It is in the response so the
   gateway counts the Free scan quota the same way as for the intake RPCs (KIT-08-01 rule 10, GW-03-01).

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto`

```proto
rpc ReadLabel(ReadLabelRequest) returns (ReadLabelResponse);
message ReadLabelRequest { string household_id = 1; string code = 2; bytes image = 3; string mime_type = 4; string language = 5; }
message ReadLabelResponse { bool found = 1; string name = 2; string brand = 3; coolinary.common.v1.Quantity pack_size = 4;
  int32 units_per_pack = 5; repeated Candidate ingredient_candidates = 6; bool ai_used = 7; }
```

## Data

No schema change.

## Errors

| Case | Status | Message |
|---|---|---|
| Image too large or wrong type | `INVALID_ARGUMENT` | `image must be JPEG, PNG or WebP up to 3 MB` |
| ai failing | `UNAVAILABLE` | `label reading unavailable` |

## Acceptance criteria

- [ ] With a stub ai, a label photo returns a proposal; nothing is written to the database.
- [ ] A 4 MB image is refused before calling ai.
- [ ] Logs and traces for the call contain no image bytes and no proposed text.
- [ ] A successful read returns `ai_used = true`; an `ai` outage returns `UNAVAILABLE` and no quota is counted.
