# AI-03 Prompts, golden sets and evals

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | Every prompt is versioned in the repo, and every job has a golden set with pass thresholds, so a model, effort or prompt change is measured before it ships instead of guessed. |

## Flow

```
prompts/<job>/v<N>/system.md + schema.py      (AI-03-01)
evals/<job>/cases/*.json (+ images)            (AI-03-02)
uv run python -m coolinary_ai.evals run <job>  (AI-03-03)  -> report.json, report.md, pass/fail per threshold (AI-03-04)
composition scorer (AI-03-05) replays the hard-limit checks the planner does
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-03-01](AI-03-01-prompt-registry.md) | Prompt registry and versions | R1 | S | AI-01-02 |
| [AI-03-02](AI-03-02-golden-set-format.md) | Golden set format, storage and privacy | R1 | S | AI-03-01 |
| [AI-03-03](AI-03-03-eval-runner.md) | Eval runner CLI | R1 | M | AI-03-02, AI-02-01 |
| [AI-03-04](AI-03-04-thresholds-and-policy.md) | Thresholds and release policy | R1 | S | AI-03-03 |
| [AI-03-05](AI-03-05-composition-scorer.md) | Composition scorer (tools, allergens, diet, stock, cuisine) | R1 | M | AI-03-03 |

Golden sets per job are their own tasks in each job epic (AI-05-02, AI-06-02, …).
