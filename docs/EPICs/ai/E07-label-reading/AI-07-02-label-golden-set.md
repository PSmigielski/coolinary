# AI-07-02 Golden set: 50 labels

| | |
|---|---|
| Service | ai |
| Epic | [AI-07 Product label reading](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | AI-07-01, AI-03-03 |
| Unblocks | label reading going live |

## Goal

Check label reading on real Polish and German packs, including own-brand and multipack products.

## Business rules

1. 50 cases: 20 Polish packs, 20 German, 10 mixed-language or imported; at least 8 multipacks, 8 own-brand
   products (Biedronka, Lidl, Aldi, Kaufland own brands), 5 unreadable photos.
2. Scoring: exact match (normalised whitespace and case) for brand and pack size; product name match by token
   overlap ≥ 0.8; `unreadable` correct.
3. Thresholds: product name ≥ 0.90, brand ≥ 0.90, pack size ≥ 0.90, multipack `unit_count` ≥ 0.90, no false
   best-before dates.

## Contract

No change.

## Data

`ai/evals/read_label/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Eval passes on Haiku 5.5 at `low`; report and baseline committed.
