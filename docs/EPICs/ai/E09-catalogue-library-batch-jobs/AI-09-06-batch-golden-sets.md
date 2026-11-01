# AI-09-06 Golden sets for catalogue and library jobs

| | |
|---|---|
| Service | ai |
| Epic | [AI-09 Catalogue and library batch jobs](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-09-01 … AI-09-05, AI-03-03 |
| Unblocks | the batch jobs going live |

## Goal

Measure each batch job before its output reaches the catalogue or the library.

## Business rules

1. All runs use `--batch` (50% cost).
2. Sets and thresholds:

   | Job | Set | Thresholds |
   |---|---|---|
   | translate_ingredients | 150 ingredients × pl, de; expected names from a native speaker | name accepted ≥ 0.95; no wrong-food translations |
   | match_nutrition_entries | 200 ingredients with candidates from Ciqual/BLS/USDA, labelled by hand | accuracy ≥ 0.90; precision of "none" ≥ 0.90; no raw/cooked mix-ups on the 30 state traps |
   | draft_dishes | 30 briefs over the 5 launch cuisines, reviewed by the founder | 0 hard violations after caller checks (scorer AI-03-05); ≥ 60% approved without edits; 0 claim words |
   | translate_dishes | 30 dishes × pl, de, reviewed by native speakers | numbers preserved 100%; accepted without edits ≥ 0.85 |
   | suggest_effort_tags | 80 dishes labelled by hand | precision ≥ 0.90 per tag; `KEEPS_3_DAYS` false positives = 0 on the 15 food-safety traps |

3. Reviews by native speakers or the founder are recorded in the case files (`expected` + reviewer initials), not in
   chat.

## Contract

No change.

## Data

`ai/evals/<job>/` for each job.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Each job's eval passes and its baseline is committed.
