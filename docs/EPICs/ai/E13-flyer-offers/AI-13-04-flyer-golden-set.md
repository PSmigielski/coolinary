# AI-13-04 Golden set per launch chain

| | |
|---|---|
| Service | ai |
| Epic | [AI-13 Flyer offer extraction](README.md) |
| Release | R2 |
| Size | M |
| Depends on | AI-13-02, AI-13-03, AI-03-03 |
| Unblocks | each chain going live (OFF-08 gate) |

## Goal

Meet the chain go-live gate with margin: offers requires ≥ 97% correct prices and dates on a 100-offer sample per
chain (features.md §8, verification 5); ai aims for ≥ 98% before that sample is even taken.

## Business rules

1. Per launch chain (Biedronka, Lidl PL, Kaufland PL, Lidl DE, Aldi Süd, Aldi Nord, Kaufland DE): 2 recent flyers,
   hand-labelled until at least 100 offers per chain; stored only under `ai/evals/` (AI-03-02 rule 3).
2. Plus 20 user-style cases: photos of printed flyers at an angle, in poor light, PDFs from other chains and
   countries.
3. Scoring per offer: product name (token overlap ≥ 0.8), price exact, pack size exact, conditions type exact, dates
   correct after applying flyer dates; recall (labelled offers found) and precision (found offers that exist).
4. Thresholds per chain: product + price + dates all correct ≥ 0.98; recall ≥ 0.95; precision ≥ 0.97; loyalty
   conditions found ≥ 0.95. User-style cases: all-correct ≥ 0.92.
5. Rerun when a chain changes its flyer layout (OFF-05 detects layout failures).

## Contract

No change.

## Data

`ai/evals/extract_flyer_offers/<chain>/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Every launch chain passes before OFF-08's gate sample is taken; results per chain in the report.
