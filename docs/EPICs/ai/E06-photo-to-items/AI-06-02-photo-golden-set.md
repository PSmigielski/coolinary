# AI-06-02 Golden set: 60 fridge and shelf photos

| | |
|---|---|
| Service | ai |
| Epic | [AI-06 Photo to items](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-06-01, AI-03-03 |
| Unblocks | photo intake going live |

## Goal

Measure recall and precision of photo reading on realistic Polish and German kitchens before launch.

## Business rules

1. 60 cases following the privacy rules in AI-03-02: 25 fridges, 10 freezers, 10 pantry shelves, 10 unpacked bags,
   5 hard cases (dim light, cluttered, reflections). Products from Polish and German shops, at least 20 cases each.
2. Expected per case: the list of visible food items (name lemma, count where countable, brand when legible), the
   non-food count, unusable flags.
3. Scoring: recall and precision on items (name match by lemma or synonym list in the case file), count accuracy on
   countable items, brand accuracy where expected.
4. Thresholds: recall ≥ 0.80, precision ≥ 0.85, count exact ≥ 0.70; no invented use-by dates (any
   `use_by_visible` not in the expected set fails the case).
5. Also run once at effort `low` and record the result in the report, to decide whether `medium` is worth its cost.

## Contract

No change.

## Data

`ai/evals/photo_to_items/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Eval passes on Haiku 5.5 at `medium`; report and baseline committed with the `low`-effort comparison.
