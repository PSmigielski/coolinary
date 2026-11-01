# OFF-01-04 Clients for identity, kitchen, ai and notifier

| | |
|---|---|
| Service | offers |
| Epic | [OFF-01 Platform](README.md) |
| Release | R2 |
| Size | S |
| Depends on | — (uses `internal/platform/rpc/client.go`) |
| Unblocks | OFF-01-03, OFF-05-06, OFF-07-02, OFF-07-04, OFF-09-02, OFF-10-03 |

## Goal

offers can call identity (plan checks, Decided Q004), kitchen (new edge, README D1), ai and notifier with the same
deadlines, metadata forwarding and metrics as every other service.

## Business rules

1. New env vars `IDENTITY_GRPC_ADDR` (default `identity:50051`) and `KITCHEN_GRPC_ADDR` (default `kitchen:50051`),
   added to `.env.example`, `.env.dev`, `.env.prod`; `AI_GRPC_ADDR` and `NOTIFIER_GRPC_ADDR` already exist.
2. Deadlines per call (caller may pass a shorter one):

   | Call | Deadline |
   |---|---|
   | identity `GetEntitlements` | 2 s |
   | kitchen `MatchNames`, `SearchIngredients`, `GetIngredients` | 5 s |
   | ai `SubmitBatch`, `GetBatch` | 30 s |
   | ai `GetBatchResults` (one page) | 30 s |
   | ai `ExtractFlyerOffers` (per request) | 120 s |
   | notifier `SendAlert` | 5 s |

3. Forward `x-request-id`, `traceparent`, `x-household-id`, `x-plan`, `x-language` when the incoming context has them
   (docs/service-conventions.md, README D4).
4. Retry only `UNAVAILABLE`, at most 2 times with 200 ms and 800 ms waits, and never for `ExtractFlyerOffers` or
   `SubmitBatch` (not idempotent).
5. identity, kitchen or ai being down never makes offers unhealthy (`/healthz` checks only Postgres and Redis).
6. `rpc_client_requests_total` is recorded for every call.

## Contract

No contract change in offers. Uses `IdentityService.GetEntitlements` (IDN-05), `KitchenService` (KIT-02, KIT-03), `AiService` (AI-04, AI-13),
`NotifierService.SendAlert` (NTF-05).

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Callee down | `UNAVAILABLE` returned to the calling task after the retries |
| Deadline passed | `DEADLINE_EXCEEDED`, never retried |

## Acceptance criteria

- [ ] With kitchen stopped, a `MatchNames` call fails with `UNAVAILABLE` after 3 attempts in under 7 s and offers'
      `/healthz` stays 200.
- [ ] A request with `x-request-id: abc` reaches kitchen with the same id (checked in kitchen's logs).
- [ ] `make grpc s=offers` still works; the new env var is documented.
