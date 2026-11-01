# IDN-02-01 Notifier gRPC client

| | |
|---|---|
| Service | identity |
| Epic | [IDN-02 Sign-in with e-mail code](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | — |
| Unblocks | IDN-02-02, IDN-09-01, IDN-05-04, IDN-01-02 (alerts) |

## Goal

Identity can call notifier (`SendEmail`, `SendPush`, `SendAlert`) with the same client rules as every other
service: keepalive, deadlines, forwarded metadata and client metrics.

## Business rules

1. Copy the channel options from `gateway/src/grpc/clients.ts` and the call wrapper from `gateway/src/grpc/unary.ts`
   into `identity/src/grpc/clients.ts` (services stay self-contained; conventions in `service-conventions.md`).
2. Address from `NOTIFIER_GRPC_ADDR` (already in `config.ts`).
3. Deadlines: `SendEmail` 5 s, `SendPush` 5 s, `SendAlert` 3 s.
4. Forward `x-request-id` and `traceparent` from the incoming call when present; set `x-user-id` and
   `x-household-id` when the caller knows them.
5. No retries in the client. Callers decide (sign-in code: the user taps "resend"; alerts: next occurrence).
6. Count every call in `rpc_client_requests_total{rpc_service,rpc_method,rpc_code}`.
7. `/healthz` does not check notifier (a service is healthy when its own stores answer).

## Contract

Uses `coolinary.notifier.v1.NotifierService` as extended by NTF-02, NTF-03 and NTF-05. No identity contract change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Notifier down or deadline hit | The wrapper throws `RpcError(UNAVAILABLE)`; callers map it (IDN-02-02 returns 503) |

## Acceptance criteria

- [ ] A test call to `SendEmail` against the dev stack lands in Mailpit.
- [ ] With notifier stopped, the call fails within 5 s with `UNAVAILABLE` and the metric shows `rpc_code="unavailable"`.
- [ ] `x-request-id` from an incoming identity RPC appears in notifier's log line for the forwarded call.
