# KIT-03-04 MatchNames: Haiku pick and fallbacks

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-03 Name matching](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-03-03, KIT-01-01, AI-08 (`PickCatalogueMatches`) |
| Unblocks | KIT-03-05, KIT-04-10, KIT-08 |

## Goal

Names the rules can't settle are decided by Claude Haiku choosing from the trigram candidates, and the answer is
checked without AI and stored for everyone.

## Business rules

1. In `INTERACTIVE` mode, names left by KIT-03-03 step 7 are sent in one `PickCatalogueMatches` call (up to 20 names
   per call; more are split into parallel calls, at most 3 at a time).
2. Each name is sent with its key, the raw name, `context`, and its candidates: id, primary name in the language and
   in `en`, category, up to 3 synonyms. No household id or user data is sent.
3. Deadline 4 s for the whole AI step.
4. Checks without AI: the chosen id must be one of that name's candidates, else the answer counts as `none` with
   `review = needed`. Confidence is clamped to 0–1.
5. Decided (Q064): confidence bands: ≥ 0.8 → `matched`, no review; 0.5–0.8 → `matched`, `review = needed`;
   below 0.5 or "none" → `none`, `review = needed`. Stored with method `ai`.
6. Fallback when ai fails, times out, refuses, or returns `RESOURCE_EXHAUSTED` (spend limit, README D11):
   - best trigram score ≥ 0.6 → answer `provisional` with that candidate (stored as `provisional`, `attempts + 1`);
   - else → `queued` (picked up by the batch, KIT-03-05).
   After 3 failed attempts within 24 hours a `provisional` row becomes `queued`.
7. AI matching is shared work: it never counts against a household's AI quota.
8. Spend guard: at most 1,000 names per hour go through the per-request AI path per kitchen instance; beyond that
   they are `queued`.
9. Metrics: `kitchen_match_total{method,outcome}`, `kitchen_match_ai_seconds`.

## Contract

Uses ai `PickCatalogueMatches` (AI-08). No change to `MatchNames`.

## Data

No schema change.

## Errors

AI failures never fail `MatchNames`; they produce `provisional` or `queued` answers.

## Acceptance criteria

- [ ] Stub ai picking a candidate with 0.9 → `matched`, stored, method `ai`.
- [ ] Stub ai picking an id that wasn't a candidate → `none`, review needed.
- [ ] Stub ai down → best trigram 0.7 gives `provisional`; best 0.4 gives `queued`.
- [ ] The 1,001st name in an hour is `queued`.
