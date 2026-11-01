# ADM-01-05 gRPC clients, error model and OpenAPI

| | |
|---|---|
| Service | admin |
| Epic | [ADM-01 Service skeleton](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | ADM-01-01; GW-01 (error model and gRPC → HTTP mapping) |
| Unblocks | every ADM-02 … ADM-07 route |

## Goal

Every admin route calls an owning service's admin RPC the same way: with a deadline, the right metadata, and errors
turned into the same HTTP error shape the app API uses.

## Business rules

1. Clients for `IdentityAdminService`, `NotifierAdminService`, `KitchenAdminService`, `PlannerAdminService`,
   `OffersAdminService`, `AiAdminService` (AI-02-04), the public `IdentityService` (`GetHousehold`, ADM-05-01) and
   `KitchenService` (`SearchIngredients`, `GetIngredients`, ADM-03), plus `grpc.health.v1.Health` for each service,
   with the gateway's channel
   options (keepalive 30 s, timeout 10 s) and its `unary` wrapper. Channels connect lazily: admin starts when the
   others are down.
2. Default deadline 5 s; 30 s for import and batch-start routes (ADM-04-01, ADM-04-02, ADM-03-06).
3. Metadata on every call: `x-request-id` (from the request or new), `traceparent` when present,
   `x-admin-user-id` (ADM-02-02), `x-language` from `Accept-Language`.
4. Error mapping and body: the gateway's GW-01 table and error JSON. If GW-01 is not done yet: `INVALID_ARGUMENT`
   400, `NOT_FOUND` 404, `ALREADY_EXISTS` 409, `FAILED_PRECONDITION` 409, `PERMISSION_DENIED` 403,
   `RESOURCE_EXHAUSTED` 429, `UNAVAILABLE` 503, `DEADLINE_EXCEEDED` 504, other 500. The gRPC detail text is passed
   through (owning services must not put personal data in it).
5. All routes live under `/admin/v1/…`; JSON only, except uploads (multipart) and flyer images (binary).
6. OpenAPI document at `/openapi` (Elysia OpenAPI plugin), one tag per area (`auth`, `kitchen`, `planner`,
   `identity`, `offers`, `notifier`, `system`). `make openapi` also exports `contracts/openapi/admin.json` for the
   later admin UI.
7. A route whose downstream service is not deployed yet answers 503 with code `service_unavailable`; it never
   crashes the process.

## Contract

`GET /openapi`, `GET /openapi/json`. Error JSON as GW-01.

## Data

None.

## Errors

See rule 4.

## Acceptance criteria

- [ ] With kitchen stopped, a kitchen route answers 503 within 5 s; other routes keep working.
- [ ] Each downstream call carries `x-admin-user-id` and `x-request-id` (checked in the callee's logs).
- [ ] `contracts/openapi/admin.json` is produced by `make openapi`.
