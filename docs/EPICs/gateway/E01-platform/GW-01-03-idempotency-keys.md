# GW-01-03 Idempotency keys for mutations

| | |
|---|---|
| Service | gateway |
| Epic | [GW-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Depends on | GW-01-02 |
| Unblocks | MOB-04 (offline queue), every mutating route marked "idempotency required" |

## Goal

The app can resend a mutation (after a timeout, a crash, or from its offline queue) without adding a fridge item
twice, ticking a line twice, reducing stock twice, or spending an AI quota twice.

## Business rules

1. Applies to `POST`, `PUT`, `PATCH`, `DELETE` under `/v1`. Header `Idempotency-Key`: 16–64 characters of
   `[A-Za-z0-9-]` (the app sends a UUIDv7). Anything else → 400 `validation`.
2. Each route declares `idempotency: "required" | "optional" | "none"`. Required routes without the header → 400
   `idempotency_key_required`. Required on: stock add/edit/remove, list lines add/edit/tick/delete, proposal accept,
   "cooked it", every AI route (day plan, swaps, intake, flyer, deals, week plan, top-up), invite accept, leave.
   Optional elsewhere; `none` for `/api/auth/*`, WebSocket and GETs.
3. Keys are scoped per user: storage key `gw:idem:<userId>:<key>`. Two users may use the same key.
4. **Fingerprint** = SHA-256 of `method + " " + route template + " " + canonical JSON body` (keys sorted; for
   multipart, the SHA-256 of each part's bytes in order). Path params are part of the route values and included.
5. **First arrival.** `SET gw:idem:<user>:<key> {"state":"pending","fp":...} NX PX <deadline + 10 s>`. If set, run the
   route. When it finishes:
   - 2xx or 4xx (except 409 `idempotency_*`, 429): store `{"state":"done","fp","status","headers","body"}` with TTL
     **24 hours**, replacing the pending marker. Stored headers: `Content-Type`, `Location`, `ETag`.
   - 5xx, 429, 499 (client left): delete the key, so a retry runs again.
6. **Repeat arrival.**
   - Stored `done` and same fingerprint → return the stored status and body, plus header `Idempotent-Replayed: true`.
     No gRPC call, no quota reservation, no rate-limit charge beyond the request itself.
   - Stored `pending` and same fingerprint → 409 `idempotency_in_progress` with `Retry-After: 2`.
   - Any state and different fingerprint → 422 `idempotency_key_reused`.
7. Responses larger than 256 KB are not stored; the key is kept as `{"state":"done","fp","status","tooLarge":true}`
   and a repeat gets 409 `idempotency_replay_unavailable` (no route is expected to hit this; a metric alerts if one does).
8. **Order with the quota (GW-03-01):** idempotency lookup runs before the quota reservation, so a replayed AI result
   never counts twice.
9. If `redis-state` is unreachable: routes with `idempotency: "required"` return 503 `service_unavailable` (safer
   than risking a double write); `optional` routes run without storing.
10. Metrics: `gateway_idempotency_total{outcome}` with `stored`, `replayed`, `in_progress`, `reused`, `skipped`.

## Contract

Request header `Idempotency-Key`; response header `Idempotent-Replayed: true` on replays. Documented once as a
reusable OpenAPI header parameter and attached to every route with `required`/`optional`.

## Data (redis-state)

| Key | Value | TTL |
|---|---|---|
| `gw:idem:<userId>:<key>` | JSON `{state, fp, status?, headers?, body?, tooLarge?}` | pending: route deadline + 10 s; done: 24 h |

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Header missing on a required route | 400 | `idempotency_key_required` |
| Malformed header | 400 | `validation` |
| Same key, still running | 409 | `idempotency_in_progress` |
| Same key, different request | 422 | `idempotency_key_reused` |
| Stored response too large to replay | 409 | `idempotency_replay_unavailable` |
| redis-state down on a required route | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] Two identical `POST /v1/stock` with the same key create one item; the second response is byte-identical with
      `Idempotent-Replayed: true` (fake kitchen counts one call).
- [ ] Two concurrent identical requests: one runs, the other gets 409 `idempotency_in_progress`.
- [ ] Same key with a changed body gets 422 `idempotency_key_reused`.
- [ ] A 503 from kitchen is not stored: the retry with the same key calls kitchen again.
- [ ] A replayed `POST /v1/plans/day` does not change the quota count (GW-03-01 test).
- [ ] Keys of user A are invisible to user B.
- [ ] Stored entries expire after 24 h (TTL asserted).

## Notes

- Decided (Q021): the replay window is **24 hours**. MOB-04 drops queued mutations older than 24 hours, so a replay
  never arrives after its key has expired.
