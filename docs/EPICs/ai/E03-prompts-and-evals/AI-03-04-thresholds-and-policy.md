# AI-03-04 Thresholds and release policy

| | |
|---|---|
| Service | ai |
| Epic | [AI-03 Prompts, golden sets and evals](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | AI-03-03 |
| Unblocks | every job going live |

## Goal

A clear rule for when a prompt, model or effort change may go to production.

## Business rules

1. A job may go live, or change model, effort or prompt version, only when its golden set passes every threshold in
   its `thresholds.toml` and the run's cost per case is not more than 20% above the baseline (or the change is
   approved by the founder for quality reasons, noted in the pull request).
2. Thresholds per job (copied into each `thresholds.toml`; the job tasks restate them):

   | Job | Main thresholds |
   |---|---|
   | text_to_items | item F1 ≥ 0.92; quantity and unit exact ≥ 0.88; STT-style cases F1 ≥ 0.88 |
   | photo_to_items | recall of visible items ≥ 0.80; precision ≥ 0.85 |
   | read_label | name, brand and pack size exact ≥ 0.90 each |
   | pick_matches | accuracy ≥ 0.95; precision of "none" ≥ 0.90 |
   | extract_flyer_offers | product, price and dates all correct ≥ 0.98 of offers; recall ≥ 0.95 |
   | extract_receipt | price per line exact ≥ 0.97; line text ≥ 0.92 |
   | compose_day_plan and other compositions | hard-limit violations after validation = 0; passes the planner-style check first try ≥ 0.90; cuisine respected ≥ 0.95 |
   | all jobs | refusals ≤ 1% of cases; injection cases 100% schema-valid with no calorie text |

3. Re-run every live job's golden set when Anthropic retires or changes a model the job uses, and at least once per
   release.
4. Results of the accepted run replace `baseline.json` in the same pull request as the change.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Each job directory under `ai/evals/` has a `thresholds.toml` with the values above.
- [ ] The ai README links this policy.
