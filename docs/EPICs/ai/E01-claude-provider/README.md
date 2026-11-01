# AI-01 Claude provider foundation

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | Every AI job in Coolinary runs through one Claude provider that returns schema-checked results, records usage, and fails in a way callers can turn into their non-AI fallback. |

## Flow

```
caller (kitchen / planner / offers)
  -> gRPC AiService.<Job>  (deadline, x-plan, x-household-id, x-language)
     -> request context + limits (AI-01-08)
     -> result cache lookup (AI-01-09)
     -> spend check (AI-02-04)
     -> job runner: prompt from registry (AI-03-01), untrusted input wrapped (AI-01-07),
        cached system prefix (AI-01-05)
     -> provider: Claude via the anthropic SDK (AI-01-03) or the stub (tests, offline dev)
     -> output validation, at most one repair call (AI-01-06)
     -> cost + metrics + Langfuse (AI-02)
  <- typed response + CallInfo, or a gRPC error with a reason token (AI-01-04)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-01-01](AI-01-01-sdk-config-and-job-registry.md) | SDK dependency, settings and job registry | R1 | M | — |
| [AI-01-02](AI-01-02-provider-interface.md) | Provider interface for structured jobs | R1 | M | AI-01-01 |
| [AI-01-03](AI-01-03-claude-provider.md) | Claude provider (Messages API, structured output, thinking, refusals) | R1 | L | AI-01-02 |
| [AI-01-04](AI-01-04-error-mapping-and-circuit-breaker.md) | Error mapping, retries, deadlines and circuit breaker | R1 | M | AI-01-03 |
| [AI-01-05](AI-01-05-prompt-caching-layout.md) | Prompt caching and request layout | R1 | M | AI-01-03 |
| [AI-01-06](AI-01-06-output-validation-and-repair.md) | Output validation, calorie guard and one repair call | R1 | M | AI-01-02 |
| [AI-01-07](AI-01-07-untrusted-input.md) | Untrusted input handling (prompt injection) | R1 | S | AI-01-02 |
| [AI-01-08](AI-01-08-grpc-limits-and-context.md) | gRPC limits, request context and concurrency | R1 | M | AI-01-01 |
| [AI-01-09](AI-01-09-result-cache-rules.md) | Result cache rules per job | R1 | S | AI-01-02 |

## Notes

- The existing `ExtractItemsFromText` path (cache, tracing, metrics interceptor, stub) is the pattern every job follows; these tasks generalise it.
- The stub provider stays. `LLM_PROVIDER=stub` keeps `make dev`, tests and CI free of API calls and cost.
