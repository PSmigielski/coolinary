# AI-01-04 Error mapping, retries, deadlines and circuit breaker

| | |
|---|---|
| Service | ai |
| Epic | [AI-01 Claude provider foundation](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-01-03 |
| Unblocks | every job epic; callers' fallbacks (KIT-08, PLN-08, OFF-07) |

## Goal

Every failure reaches the caller as one predictable gRPC status with a short reason token, fast enough that the
caller can still show its non-AI result within its own deadline.

## Business rules

1. Status message format: `<reason>: <short detail>`; the reason token is the first word. Callers branch on the
   status code and the reason, never on the detail text.
2. Mapping:

   | Source | gRPC status | Reason |
   |---|---|---|
   | Request field invalid (empty text, too many images, bad custom_id) | `INVALID_ARGUMENT` | `ai_bad_request` |
   | `ProviderRefusal` | `FAILED_PRECONDITION` | `ai_refused` |
   | `ProviderInvalidOutput` / `ProviderTruncated` after the one repair (AI-01-06) | `UNAVAILABLE` | `ai_invalid_output` |
   | `RateLimitError` (429) after SDK retries | `UNAVAILABLE` | `ai_rate_limited` |
   | `APIStatusError` 500/529 overloaded after retries | `UNAVAILABLE` | `ai_overloaded` |
   | `APIConnectionError`, `APITimeoutError` | `UNAVAILABLE` | `ai_unreachable` |
   | `AuthenticationError`, `PermissionDeniedError` | `UNAVAILABLE` | `ai_misconfigured` (logged at `error`) |
   | `BadRequestError` (400 from the API) | `INTERNAL` | `ai_request_bug` |
   | Daily spend limit hit (AI-02-04) | `RESOURCE_EXHAUSTED` | `ai_budget` |
   | Job disabled (AI-01-01) | `UNAVAILABLE` | `ai_disabled` |
   | Circuit open (rule 6) | `UNAVAILABLE` | `ai_circuit_open` |
   | Caller deadline passed | `DEADLINE_EXCEEDED` | — |

3. Rule for callers (documented in the ai README): **fall back to the non-AI result on every status except
   `INVALID_ARGUMENT`**, which is a caller bug.
4. Timeouts: the SDK request timeout is `min(job.timeout_seconds, grpc_time_remaining - 1.0 s)`. If less than
   3 seconds remain when the handler starts, return `DEADLINE_EXCEEDED` without calling Claude.
5. Retries: the SDK's own retries (429, 5xx, connection errors) with `max_retries` from the job spec, and only while
   the remaining deadline allows another full attempt; never retry `ai_refused` or `ai_bad_request`.
6. Circuit breaker per model: after 5 consecutive `ai_overloaded` / `ai_unreachable` / `ai_rate_limited` results
   within 60 seconds, the circuit opens for 30 seconds; calls to that model fail fast with `ai_circuit_open`. After
   30 seconds one trial call is let through; success closes the circuit, failure reopens it for 60 seconds
   (doubling up to 5 minutes).
7. Every error increments `ai_errors_total{job, reason}` and logs one line with `job`, `reason`, `request_id`,
   `duration_ms`; no request content.

## Contract

No gRPC change. Reason tokens are listed in `contracts/coolinary/ai/v1/ai.proto` as a comment block above the
service, so callers can see them.

## Data

No schema change.

## Errors

This task defines them (table above).

## Acceptance criteria

- [ ] Table-driven test: each source row produces the listed status and reason.
- [ ] With a mocked API returning 529 six times in a row, the sixth call fails in under 50 ms with `ai_circuit_open`; after 30 s a trial call is made.
- [ ] A call arriving with 2 s of deadline left returns `DEADLINE_EXCEEDED` and makes no HTTP request.
- [ ] Refusals are never retried (request count = 1).
