# AI-01-05 Prompt caching and request layout

| | |
|---|---|
| Service | ai |
| Epic | [AI-01 Claude provider foundation](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-01-03 |
| Unblocks | every job epic |

## Goal

Lay out every request so the stable part (instructions, schema, catalogue or library context) is read from
Claude's prompt cache at about 0.1× the input price, and prove it with a metric.

## Business rules

1. Render order is `tools → system → messages`. Jobs use no tools, so the cached prefix is the system prompt plus,
   where a job has it, a stable context block.
2. `system_blocks` (job instructions, output rules, allowed enums) never contain per-request data: no dates, no
   request ids, no household data, no language name if it varies — the language goes in the per-request part.
3. `context_blocks` (e.g. the cuisine/tool vocabulary, a catalogue slice for drafting) are serialised
   deterministically: JSON with `sort_keys=True`, lists sorted by id, no sets. A context block carries its own
   version string (e.g. `vocab:2026-11-01`) so a change is a deliberate new prefix.
4. One `cache_control: {"type": "ephemeral"}` breakpoint goes on the last stable block. Sync jobs use the default
   5-minute lifetime; batch jobs use `"ttl": "1h"` (AI-04-01). At most 4 breakpoints per request.
5. The stable prefix of every job is at least 512 tokens (the minimum for Haiku 5.5 and Sonnet 5.5); a unit test
   counts tokens of each job's prefix with `messages.count_tokens` (recorded once, stored as a fixture) and fails
   if a prefix falls under 512 — under that it silently doesn't cache.
6. Per-request content (user text, images, stock lists) goes in the user message after the breakpoint.
7. Metrics: `ai_prompt_cache_tokens_total{job, kind="read"|"write"}`. A Grafana panel shows the read share per job;
   a sync job with more than 20 calls in an hour and a read share under 50% is a bug to investigate.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

None new.

## Acceptance criteria

- [ ] A test builds the same job request twice with different user content and asserts the bytes before the breakpoint are identical.
- [ ] Live smoke (manual): two `text_to_items` calls 10 s apart show `cache_read_input_tokens > 0` on the second.
- [ ] No job prefix contains a timestamp, uuid or household id (test scans the rendered prefix for patterns).
