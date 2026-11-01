# AI-02-01 Price table and cost calculator

| | |
|---|---|
| Service | ai |
| Epic | [AI-02 Cost tracking and spend limit](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-01-01 |
| Unblocks | AI-02-02, AI-02-03, AI-02-04 |

## Goal

Turn the usage of any completion into an exact cost in USD and EUR, including cache reads, cache writes, the
Haiku long-prompt rate and the batch discount.

## Business rules

1. Prices live in `src/coolinary_ai/prices.toml` (USD per million tokens), checked against the public pricing page
   at implementation time:

   | Model | Input ≤ 100K-token prompt | Output ≤ 100K | Input > 100K | Output > 100K |
   |---|---|---|---|---|
   | claude-haiku-5-5 | 0.10 | 0.50 | 0.50 | 2.50 |
   | claude-sonnet-5-5 | 2.00 | 10.00 | 2.00 | 10.00 |

   Other models that a Sonnet fallback can answer with (AI-01-03 rule 7) get rows too; an unknown `response.model`
   is priced at the most expensive row and logged at `warn`.
2. Multipliers on the input rate: cache read 0.1×, 5-minute cache write 1.25×, 1-hour cache write 2×.
3. Prompt size for the Haiku rate choice = `input_tokens + cache_read + cache_write` of that request.
4. Batch requests cost 50% of the same request's normal price (applied to all token kinds).
5. A repaired call (AI-01-06) is two completions; both are priced and summed.
6. Cost in EUR = USD × `AI_USD_TO_EUR` (config, default 0.92), stored as integer micro-euros (1 € = 1,000,000),
   rounded half up per completion.
7. Thinking tokens are billed as output and are already inside `output_tokens`; no separate field.

## Contract

`contracts/coolinary/ai/v1/common.proto` (new, package `coolinary.ai.v1`):

```proto
// Returned by every AI RPC.
message CallInfo {
  string model = 1;           // model that answered, after any fallback
  string prompt_version = 2;  // "<job>/v<N>"
  bool cached = 3;            // served from the result cache
  int64 cost_eur_micros = 4;  // 0 when cached
  bool repaired = 5;          // a repair call was needed
}
```

Every response message gets `CallInfo call = <next free field>;`. `ExtractItemsFromTextResponse` keeps `model`
and `cached` for compatibility and adds `call`.

## Data

No schema change.

## Errors

None new.

## Acceptance criteria

- [ ] Worked examples as unit tests: Haiku 1,000 input + 500 output = $0.00035; Haiku 90,000 cache read + 20,000 input + 1,000 output (prompt 110K) = 90,000×0.5×0.1 + 20,000×0.5 + 1,000×2.5 per million = $0.0170; Sonnet batch 10,000 input + 2,000 output = $0.02.
- [ ] Changing a price in the file changes the computed cost without a code change.
- [ ] `CallInfo.cost_eur_micros` is 0 for a cache hit and > 0 for a live call.
