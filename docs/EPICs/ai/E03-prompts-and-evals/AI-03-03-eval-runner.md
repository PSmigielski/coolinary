# AI-03-03 Eval runner CLI

| | |
|---|---|
| Service | ai |
| Epic | [AI-03 Prompts, golden sets and evals](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-03-02, AI-02-01 |
| Unblocks | AI-03-04, AI-03-05, every golden-set task |

## Goal

One command runs a job's golden set against a chosen model, effort and prompt version, shows the cost first, and
writes a report that says pass or fail against the job's thresholds.

## Business rules

1. Command: `uv run python -m coolinary_ai.evals run <job> [--model M] [--effort E] [--prompt vN] [--limit N]
   [--tags t1,t2] [--batch] [--yes]`. Also `make ai-eval job=<job>` in the root Makefile, run inside the ai
   container.
2. Before calling the API it prints the estimated cost (AI-02-04 rule 3 estimate × cases) and stops unless `--yes`
   is given. Evals never run in CI and never with `LLM_PROVIDER=stub` (stub runs are a separate unit-test target).
3. `--batch` submits the cases through Message Batches (AI-04) at 50% cost and waits for results (poll every 60 s,
   give up after 24 h).
4. Each case runs through the real `run_job` path (cache bypassed), so validation and repair are part of what is
   measured.
5. Scoring functions are per job (`evals/<job>/score.py`), returning per-case metrics; the runner aggregates them.
6. Output in `ai/evals/<job>/runs/<date>-<model>-<effort>-<prompt>/`: `report.json` (per case and totals: metrics,
   tokens, cost, repairs, refusals, latency p50/p95) and `report.md` (summary table, failures with case ids). Runs
   are git-ignored except the file `ai/evals/<job>/baseline.json`, which holds the last accepted totals.
7. Exit code 0 if every threshold in `evals/<job>/thresholds.toml` passes, 1 otherwise; the report shows the delta to
   the baseline.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] Running without `--yes` prints the estimate and makes no API call.
- [ ] A run on a 3-case toy set writes both report files with tokens and cost per case.
- [ ] A threshold failure exits 1 and names the metric.
