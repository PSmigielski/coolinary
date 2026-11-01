# AI-08-03 Golden set: 500 names

| | |
|---|---|
| Service | ai |
| Epic | [AI-08 Catalogue match pick](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-08-01, AI-03-03 |
| Unblocks | matching going live |

## Goal

Matching quality decides whether recipes, stock and offers line up; measure it on the kinds of names users and
flyers actually produce.

## Business rules

1. 500 cases on the fixture catalogue (AI-03-05 rule 2): 150 product names from packs (pl/de), 100 receipt lines
   with abbreviations, 150 flyer offer names, 50 typed names in en, 50 traps (composite products, near misses such as
   "masło orzechowe" vs "masło", where the right answer is a specific candidate or none).
2. Candidates per case are produced the way KIT-03 will: top 10 by pg_trgm similarity in the case's language (a
   script in `ai/evals/pick_matches/make_candidates.py` runs against a local Postgres with the fixture catalogue).
3. Scoring: accuracy (right id, or none when none is right); precision and recall of "none"; calibration table
   (accuracy per confidence band 0.6–0.9, ≥ 0.9).
4. Thresholds: accuracy ≥ 0.95; precision of "none" ≥ 0.90; accuracy in the ≥ 0.90 confidence band ≥ 0.98 (so kitchen
   can auto-accept there).

## Contract

No change.

## Data

`ai/evals/pick_matches/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Eval passes on Haiku 5.5 at `low`; report includes the calibration table.
