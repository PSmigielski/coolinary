# AI-16-03 Golden set: 40 receipts

| | |
|---|---|
| Service | ai |
| Epic | [AI-16 Receipt reading](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | AI-16-02, AI-03-03 |
| Unblocks | receipt intake going live |

## Goal

Receipt reading is right often enough that prices from receipts can be trusted for the trip plan.

## Business rules

1. 40 redacted receipts (AI-03-02 rule 3): 20 Polish (Biedronka, Lidl, Kaufland, Dino, Żabka), 20 German (Lidl,
   Aldi Nord, Aldi Süd, Rewe, Edeka); 8 long receipts in 2–3 photos; 5 digital receipts from loyalty apps; 4 crumpled
   or faded.
2. Scoring: line recall; unit price and line total exact; line text token overlap ≥ 0.8; total exact; mismatch flag
   correct; privacy leaks (must be 0).
3. Thresholds: line recall ≥ 0.95; prices exact ≥ 0.97; line text ≥ 0.92; total exact ≥ 0.95; privacy leaks 0.

## Contract

No change.

## Data

`ai/evals/extract_receipt/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Eval passes; baseline committed.
