# AI-05-02 Golden set: typed and spoken item lists

| | |
|---|---|
| Service | ai |
| Epic | [AI-05 Text to items](README.md) |
| Release | R3 |
| Size | S |
| Depends on | AI-05-01, AI-03-03 |
| Unblocks | typed and voice intake going live in R3 (KIT-08-03) |

## Goal

Prove that typed and spoken item lists are read well enough in Polish, German and English before users see them.

## Business rules

1. 150 cases: 50 `pl`, 50 `de`, 50 `en`. In each language: 30 typed (with punctuation, typos, abbreviations like
   "2kg ziemn."), 15 speech-recognition style (no punctuation, number words, filler: "no i jeszcze"), 5 with
   non-food items.
2. Plus 6 injection cases (AI-01-07) per language.
3. Expected per case: the item list (name normalised to a lemma for comparison, quantity, unit, hints).
4. Scoring: item F1 by name match (lemma, case-insensitive); quantity+unit exact on matched items;
   `ignored_count` correct.
5. Thresholds: F1 ≥ 0.92 overall and ≥ 0.88 on speech-style cases; quantity+unit exact ≥ 0.88; injection cases 100%
   schema-valid with no calorie text.

## Contract

No change.

## Data

`ai/evals/text_to_items/` (cases, `score.py`, `thresholds.toml`, `baseline.json`).

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] `make ai-eval job=text_to_items` passes on Haiku 5.5 at effort `low`; the report and baseline are committed.
- [ ] The report also lists cost per case (expected ≈ $0.0003).

## Notes

- Decided (Q219): typed and voice intake ship in release 3, so this golden set is built in R3, before KIT-08-03
  goes live.
