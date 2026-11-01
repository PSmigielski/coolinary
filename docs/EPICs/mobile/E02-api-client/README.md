# MOB-02 API client and error model

Goal: one typed client for every gateway call, with the same headers, timeouts, idempotency keys and a small set of
error types; the app tells offline from server-down and knows its Plus entitlements and AI quota.

## Flow

```
hook / outbox ──> openapi-fetch client ──> middleware: Bearer, x-request-id, Accept-Language,
                                           Idempotency-Key, x-client, timeout by route class
                                     ──> gateway ──> response
                                                      ├─ 2xx: quota headers → entitlements state
                                                      └─ error: toApiError() → offline | serverDown | quotaReached | plusRequired | conflict | …
expo-network ──> onlineManager (pause / resume queries, replay outbox)
network error while online ──> GET /healthz probe ──> serverDown
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [MOB-02-01](MOB-02-01-request-middleware.md) | Request middleware | R1 | M | MOB-01-04, GW-01 |
| [MOB-02-02](MOB-02-02-error-model.md) | Typed error model and retry policy | R1 | M | MOB-02-01, MOB-02-04, GW-01 |
| [MOB-02-03](MOB-02-03-openapi-workflow.md) | OpenAPI types workflow | R1 | S | GW-01 |
| [MOB-02-04](MOB-02-04-connectivity.md) | Connectivity and server reachability | R1 | S | MOB-01-01 |
| [MOB-02-05](MOB-02-05-entitlements-quota.md) | Entitlements and AI quota state | R1 | M | MOB-01-04, MOB-02-02, GW-02, GW-03, IDN-05 |
