# AI-01-08 gRPC limits, request context and concurrency

| | |
|---|---|
| Service | ai |
| Epic | [AI-01 Claude provider foundation](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-01-01 |
| Unblocks | AI-06-01, AI-13-02, AI-04-01 |

## Goal

The gRPC server accepts the payload sizes the jobs need, reads the caller's context in one place, and limits how
many Claude calls run at once so one burst can't exhaust rate limits for everyone.

## Business rules

1. Server option `grpc.max_receive_message_length` = 32 MiB (photos, PDFs and batch submissions; the gRPC default is
   4 MiB). Send limit stays 4 MiB except `GetBatchResults`, which streams.
2. A `RequestContext` is built per call from metadata: `x-request-id`, `traceparent`, `x-household-id`,
   `x-user-id`, `x-plan` (`free`/`plus`, anything else treated as `free`), `x-language`. A request field
   `household_id` wins over metadata when both are present (D3). Missing `x-plan` is logged at `warn` once per
   caller and treated as `free`, so the cheaper rules apply.
3. Concurrency: one semaphore per model for sync calls — 24 in flight for Haiku, 12 for Sonnet (config
   `AI_MAX_INFLIGHT_HAIKU`, `AI_MAX_INFLIGHT_SONNET`). A call waits at most 2 seconds for a slot, then fails
   `UNAVAILABLE ai_overloaded`. Batch submission does not use these slots.
4. `/healthz` keeps checking only Redis (conventions: a service is healthy even when what it calls is down);
   Anthropic reachability is a metric (`ai_provider_up{model}` from the circuit breaker, AI-01-04), not a health
   check.
5. The metrics interceptor in `grpc_server.py` also covers server-streaming RPCs (needed for `GetBatchResults`).

## Contract

No gRPC change.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Message over 32 MiB | `RESOURCE_EXHAUSTED` (gRPC runtime) |
| No concurrency slot within 2 s | `UNAVAILABLE`, `ai_overloaded` |

## Acceptance criteria

- [ ] A 10 MiB photo request is accepted; a 33 MiB one is rejected by the runtime.
- [ ] With the Sonnet semaphore at 12 and 13 slow fake calls, the 13th fails after about 2 s with `ai_overloaded`.
- [ ] `x-plan: plus` and `x-household-id` reach the tracing call (asserted with a fake tracer).
