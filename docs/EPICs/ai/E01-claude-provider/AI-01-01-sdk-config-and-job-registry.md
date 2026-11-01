# AI-01-01 SDK dependency, settings and job registry

| | |
|---|---|
| Service | ai |
| Epic | [AI-01 Claude provider foundation](README.md) |
| Release | R1 |
| Size | M |
| Depends on | — |
| Unblocks | AI-01-02, AI-01-08, AI-02-01 |

## Goal

Add the official `anthropic` Python SDK and one place where every AI job's model, effort, token limit, timeout and
cache lifetime is configured, so changing a model is a config change, not a code change.

## Business rules

1. Add `anthropic` (latest 1.x at implementation time, pinned in `uv.lock`) to `ai/pyproject.toml`. No other LLM SDK.
2. New settings in `config.py` (each also readable as `X_FILE`, per conventions):
   - `ANTHROPIC_API_KEY` — required when `LLM_PROVIDER=anthropic`; prod reads `ANTHROPIC_API_KEY_FILE=/run/secrets/anthropic_api_key`.
   - `LLM_PROVIDER` — `anthropic` or `stub`. Default stays `stub`. Unknown value: the service refuses to start.
   - `AI_JOBS_FILE` — path to the job registry, default `src/coolinary_ai/jobs.toml`.
   - `AI_USD_TO_EUR` — conversion rate for cost in euros, default `0.92`, see AI-02-01.
3. The job registry `jobs.toml` holds one table per job (names in the job table of `docs/EPICs/ai/README.md`):
   `model`, `effort` (`low` … `max`), `max_tokens`, `timeout_seconds`, `max_retries`, `cache_ttl_seconds`,
   `enabled` (bool), `fallbacks` (bool, Sonnet jobs only), and optional `[<job>.by_plan.free]` /
   `[<job>.by_plan.plus]` overrides for `model` and `effort`. The plan is read from `x-plan` (AI-01-08); an empty or
   unknown plan uses the `free` override. Decision 15: `compose_day_plan` ships with
   `[compose_day_plan.by_plan.free] model = "claude-haiku-5-5"`, and that override may not be removed by env.
4. Starting values (tuned later by evals, AI-03):

   | Job | model | effort | max_tokens | timeout_s |
   |---|---|---|---|---|
   | text_to_items | claude-haiku-5-5 | low | 4000 | 15 |
   | photo_to_items | claude-haiku-5-5 | medium | 8000 | 25 |
   | read_label | claude-haiku-5-5 | low | 3000 | 15 |
   | pick_matches | claude-haiku-5-5 | low | 6000 | 20 |
   | compose_day_plan | claude-haiku-5-5 (Free, fixed); Plus: claude-haiku-5-5 until AI-10-05 shows Sonnet 5.5 clearly better | medium | 16000 | 55 |
   | modify_meal | claude-sonnet-5-5 | low | 8000 | 40 |
   | rewrite_dish | claude-sonnet-5-5 | low | 8000 | 40 |
   | compose_family_day | claude-sonnet-5-5 | medium | 16000 | 55 |
   | compose_from_deals | claude-sonnet-5-5 | medium | 16000 | 55 |
   | compose_week_plan | claude-sonnet-5-5 | high | 32000 | 110 |
   | cheapest_top_up | claude-sonnet-5-5 | medium | 12000 | 55 |
   | extract_flyer_offers | claude-haiku-5-5 | medium | 12000 | 40 |
   | extract_receipt | claude-haiku-5-5 | medium | 6000 | 25 |
   | batch kinds | see the kinds table in [AI-04](../E04-batch-jobs/README.md) | | | |

5. Model ids are written exactly `claude-haiku-5-5` and `claude-sonnet-5-5`, never with a date suffix. The
   registry rejects any model id not in the price table (AI-02-01) at start-up.
6. Per-request overrides through env: `AI_JOB_<JOB>_MODEL`, `AI_JOB_<JOB>_EFFORT`, `AI_JOB_<JOB>_ENABLED`
   (upper-case job name). Env wins over the file.
7. A job with `enabled=false` returns `UNAVAILABLE` with reason `ai_disabled` (AI-01-04) without calling Claude, so
   any AI feature can be switched off in prod by env and the caller's non-AI fallback runs.
8. The registry is loaded once at start-up and validated: effort in the allowed set, `max_tokens` between 256 and
   64000, `timeout_seconds` between 5 and 120, `max_retries` 0–3. Any error stops the service with a log line naming
   the job and field.
9. `.env.example`, `.env.dev` (`LLM_PROVIDER=stub`), `.env.prod` (`LLM_PROVIDER=anthropic`,
   `ANTHROPIC_API_KEY_FILE`) and `scripts/make-secrets.sh` (placeholder `secrets/anthropic_api_key` that must be
   filled by hand, like `smtp_url`) are updated.

## Contract

No gRPC change.

## Data

No schema change. New file `ai/src/coolinary_ai/jobs.toml`.

## Errors

| Case | Behaviour |
|---|---|
| `LLM_PROVIDER=anthropic` and no key | Service exits at start-up: `ANTHROPIC_API_KEY is required` |
| Invalid registry value | Service exits at start-up naming job and field |
| Job disabled | `UNAVAILABLE`, `ai_disabled` |

## Acceptance criteria

- [ ] `uv run pytest` passes with `LLM_PROVIDER=stub` and no API key set.
- [ ] Setting `AI_JOB_COMPOSE_DAY_PLAN_MODEL=claude-sonnet-5-5` changes the Plus model used, visible in the response's `CallInfo.model` (AI-01-03), while a call with `x-plan: free` still runs on `claude-haiku-5-5`.
- [ ] A registry with `effort = "extreme"` fails start-up with a clear message.
- [ ] `AI_JOB_PHOTO_TO_ITEMS_ENABLED=false` makes `ExtractItemsFromPhoto` return `UNAVAILABLE ai_disabled` with no outbound HTTP call (checked with a mocked transport).
- [ ] The API key never appears in logs, traces or `/healthz`.

## Notes

- `by_plan` exists because the cost caps differ by plan. Decision 15 (README): Free day plans compose on Haiku 5.5;
  Plus moves to Sonnet 5.5 only on the measured comparison in AI-10-05.
