# PLN-01-02 Clients for kitchen, offers, ai, notifier, identity

| | |
|---|---|
| Service | planner |
| Epic | [PLN-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Depends on | — |
| Unblocks | PLN-01-03, PLN-01-04, PLN-03, PLN-04, PLN-08, PLN-10, PLN-17 |

## Goal

One place that builds connect-go clients to every service planner calls, with deadlines, metadata forwarding and
predictable failure types, so feature code can fall back to non-AI results without guessing.

## Business rules

1. Clients: kitchen (`KITCHEN_GRPC_ADDR`), offers (`OFFERS_GRPC_ADDR`), ai (`AI_GRPC_ADDR`), notifier
   (`NOTIFIER_GRPC_ADDR`, new, default `notifier:50051`; D1 adds the `planner → notifier` edge) and identity
   (`IDENTITY_GRPC_ADDR`, new, default `identity:50051`; Q004 adds the `planner → identity` edge for plan lookups). All use
   `rpc.NewHTTPClient()` and `rpc.ClientOptions()` (gRPC protocol, h2c, keepalive).
2. Every call sets a deadline: the smaller of the caller's remaining deadline minus 200 ms and the per-call default:
   kitchen 2 s, offers 2 s, identity 1 s, notifier 3 s, ai text/match calls 15 s, ai compose calls 45 s, ai batch submit/poll 10 s.
3. Metadata forwarded from the incoming request when present: `x-request-id`, `traceparent`, `x-user-id`,
   `x-household-id`, `x-plan`, `x-language`. Background jobs set a fresh `x-request-id` (`job-<uuid>`) and the
   household's plan looked up from identity (PLN-01-05).
4. Errors are mapped once into three kinds the feature code switches on: `Unavailable` (UNAVAILABLE,
   DEADLINE_EXCEEDED, connection refused), `Rejected` (INVALID_ARGUMENT, FAILED_PRECONDITION, NOT_FOUND,
   PERMISSION_DENIED, RESOURCE_EXHAUSTED — the original code is kept), and `Bug` (INTERNAL, UNKNOWN, anything else).
5. No automatic retries for ai calls (each costs money). Kitchen and offers read calls retry once after 100 ms on
   `Unavailable`; writes (`ConsumeStock`, `AddStockItems`) never retry automatically — they carry an idempotency key
   instead (PLN-10-02).
6. A client never logs request or response bodies.
7. Config is read in `internal/config`; `.env.example`, `.env.dev`, `.env.prod` gain `NOTIFIER_GRPC_ADDR` and
   `IDENTITY_GRPC_ADDR`.

## Contract

No contract change; uses generated clients of kitchen, offers, ai, notifier, identity.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Deadline exceeded | Kind `Unavailable`; `rpc_client_requests_total{rpc_code="deadline_exceeded"}` |
| ai returns RESOURCE_EXHAUSTED (spend limit, D11) | Kind `Rejected` with code kept; callers fall back to non-AI |

## Acceptance criteria

- [ ] With kitchen stopped, a kitchen read returns kind `Unavailable` within 2.3 s (one retry included).
- [ ] An incoming request with `x-plan: plus` and `x-language: de` reaches ai with both keys (checked with a test server).
- [ ] An ai call is attempted exactly once even when it fails with UNAVAILABLE.
- [ ] `rpc_client_requests_total` is labelled with the target `rpc_service`, `rpc_method` and `rpc_code`.
