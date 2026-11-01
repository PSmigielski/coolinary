# GW-01 Platform: request context, error model, idempotency keys, OpenAPI rules

Release R1. Everything else in the gateway builds on these five tasks.

## Goal

Every `/v1` route behaves the same way: it knows who is calling and in which language, forwards that to the
services, turns gRPC failures into one predictable JSON error, can be safely retried by the app's offline queue, and
is described in the OpenAPI document the app generates its types from.

## Flow of one request

1. `onRequest`: request id, client IP, language (GW-01-01).
2. Auth macro verifies the JWT (exists; extended in GW-02-01).
3. Rate limit (GW-03-04), then idempotency lookup for mutations (GW-01-03).
4. Quota reservation for AI routes (GW-03-01), Plus guard (GW-03-03).
5. gRPC call with forwarded metadata and the route's deadline class (GW-01-01).
6. Errors mapped to problem JSON (GW-01-02); success stored for idempotent replay (GW-01-03).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-01-01](GW-01-01-request-context.md) | Request context, metadata forwarding and deadline classes | R1 | S | — |
| [GW-01-02](GW-01-02-error-model.md) | Problem JSON error model and gRPC → HTTP mapping | R1 | M | GW-01-01 |
| [GW-01-03](GW-01-03-idempotency-keys.md) | Idempotency keys for mutations | R1 | M | GW-01-02 |
| [GW-01-04](GW-01-04-openapi-rules.md) | OpenAPI rules, tags and export check | R1 | S | GW-01-02 |
| [GW-01-05](GW-01-05-uploads.md) | Photo and PDF uploads | R1 | M | GW-01-02 |
