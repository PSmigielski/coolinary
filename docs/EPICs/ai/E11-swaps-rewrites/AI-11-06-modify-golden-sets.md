# AI-11-06 Golden sets for swaps, wishes and rewrites

| | |
|---|---|
| Service | ai |
| Epic | [AI-11 Swaps, wishes, tool rewrite, restyle](README.md) |
| Release | R1–R3 |
| Size | M |
| Depends on | AI-11-01 … AI-11-05, AI-03-05 |
| Unblocks | each mode going live |

## Goal

Each mode passes its own small golden set before it is offered.

## Business rules

1. Sets (all on the fixture catalogue, scored by AI-03-05 plus mode checks):

   | Mode | Cases | Thresholds |
   |---|---|---|
   | SWAP (R1) | 30 days from the day-plan set with one meal swapped | violations 0; different dish 100%; stock fits with other meals 100% |
   | EXPAND_ALTERNATIVE (R1) | 20 alternatives | ingredients/stock/to-buy unchanged 100% |
   | REPLACE_INGREDIENT (R1) | 30 dishes × 1 missing ingredient | first suggestion accepted by reviewer ≥ 0.80; hard-limit violations 0 |
   | WISH (R1) | 40 wishes (pl/de/en, 10 conflicting, 6 injection) | violations 0; correct `wish_status` ≥ 0.90; time wishes met 100% when possible |
   | TOOLS (R2) | 20 dishes × 2 tool sets | tools respected 100%; ingredient set unchanged 100%; reviewer accepts ≥ 0.80 |
   | RESTYLE (R3) | 15 dishes × 2 cuisines | MAIN/BASE kept 100%; reviewer judges the cuisine right ≥ 0.80 |

2. Reviewer judgements are recorded in the case files.

## Contract

No change.

## Data

`ai/evals/modify_meal/`, `ai/evals/rewrite_dish/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] R1 modes pass before the R1 release; TOOLS before R2; RESTYLE before R3.
