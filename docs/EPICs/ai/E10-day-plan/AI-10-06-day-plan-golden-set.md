# AI-10-06 Golden set and eval for day plans

| | |
|---|---|
| Service | ai |
| Epic | [AI-10 Day plan composition and cuisine choices](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-10-04, AI-03-05 |
| Unblocks | AI-10-05; day plans going live |

## Goal

Prove day plans keep the hard limits, use what expires, respect cuisines and stay within cost, on the cases
features.md's verification list names.

## Business rules

1. 60 cases on the fixture catalogue (AI-03-05):
   - 17 "student room" cases (tools: microwave, kettle) — 51 meals, covering verification 2 (50 dishes need no other
     tool).
   - 17 "no gluten, no nuts" cases — 51 meals, covering verification 3 (zero violations).
   - 10 cuisine-fit cases: 2 per launch cuisine, one fridge that fits and one that doesn't (verification 4).
   - 10 expiry cases: 3–5 lines expiring within 2 days among 20–40 lines.
   - 6 library-first cases where a candidate fits fully.
   - Injection is not relevant (no user text), but 2 cases have candidate names containing instruction-like text,
     to check names are treated as labels.
2. Fridges are Polish and German shopping baskets; languages split pl / de / en.
3. Metrics: violations after validation (must be 0), first-try pass rate (no repair), expiry-use share (of lines
   expiring within 2 days that fit some allowed dish, share used), cuisine respected or correct alternatives,
   library-first rate, cost per call, p95 latency.
4. Thresholds: violations 0; first-try pass ≥ 0.90; expiry-use ≥ 0.80; cuisine ≥ 0.95; library-first ≥ 0.80;
   p95 latency ≤ 40 s.

## Contract

No change.

## Data

`ai/evals/compose_day_plan/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Eval passes on the routing chosen in AI-10-05; baseline committed.
