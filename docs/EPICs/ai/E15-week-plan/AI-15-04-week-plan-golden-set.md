# AI-15-04 Golden set: family weeks

| | |
|---|---|
| Service | ai |
| Epic | [AI-15 Week plan](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | AI-15-03, AI-03-05 |
| Unblocks | week plans going live |

## Goal

Meet features.md verification 18 on AI output before release 3: a test family of four with different goals and a
fixed fridge, a 7-day plan set to "every other day", at most 4 cooking sessions and one dish per meal unless an
allergy or diet forces a second.

## Business rules

1. 20 cases: the verification-18 family (from the worked example, 4 eaters) with 3 different fridges; 6 single-person
   weeks; 6 families with one diverging diet; 5 "twice a week" households.
2. The scorer replays the portion check of AI-12-04 rule 2 on library dishes and counts sessions as AI-15-03 rule 3.
3. Thresholds: violations 0; rhythm met 100%; second dishes only when forced 100%; leftovers within keeping time
   100%; portion check ≥ 0.90; expiring-stock used in days 1–3 ≥ 0.80; p95 latency ≤ 100 s; cost per call ≤ $0.15.

## Contract

No change.

## Data

`ai/evals/compose_week_plan/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Eval passes; baseline committed.
