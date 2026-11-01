# ADM-07-01 System overview

| | |
|---|---|
| Service | admin |
| Epic | [ADM-07 AI cost and system overview](README.md) |
| Release | R1 |
| Size | M |
| Depends on | ADM-01-05; KIT-01-07, PLN-01-06, OFF-01-06 (`GetQueueCounts`); NTF-06-02 (`ListAlerts`) |
| Unblocks | — |

## Goal

The admin start page answers "is everything running, and what is waiting for me?" in one look.

## Business rules

1. Health: gRPC `grpc.health.v1.Health/Check` on identity, notifier, kitchen, planner, offers and ai, and HTTP
   `GET http://gateway:8081/healthz`, each with a 1 s timeout; shown as `SERVING`, `NOT_SERVING` or `unreachable`.
2. Work waiting, from each owning service's `GetQueueCounts` (2 s deadline each). Each answers a map with stable
   keys and leaves out queues it hasn't built yet; admin shows the keys it knows and ignores others:
   - kitchen (KIT-01-07 rule 8): `match_reviews`, `nutrition_suggestions`, `unreviewed_names`, `product_reviews`,
     `parked_items`, `open_batches`, `dead_letters`;
   - planner (PLN-01-06): `review_queue`, `dietitian_pending`, `effort_tags_pending`, `promotion_candidates` (R2),
     `open_batches`;
   - offers (OFF-01-06 rule 7, R2): `failing_sources`, `open_samples`, `held_offers`, `shop_reviews`,
     `store_link_queue`, `feed_store_queue` (R3);
   - notifier: open alerts from `ListAlerts(state = open)`, unacknowledged first.
3. Each part is fetched in parallel; a part that fails shows "unavailable" and the rest still renders (the route
   never fails as a whole because of one service).
4. Links: Grafana dashboards (*Coolinary Services*, *Containers*, *Data Stores*) at `GRAFANA_PUBLIC_URL`, Langfuse at
   `LANGFUSE_PUBLIC_URL`, Mailpit in dev.
5. Readable by roles `admin` and `dietitian` (a dietitian sees only planner's `dietitian_pending` and health).
6. Not audited (no personal data).

## Contract

| Method | Path | Response |
|---|---|---|
| GET | `/admin/v1/overview` | `{ health: { <svc>: status }, queues: { kitchen: {...}, planner: {...}, offers: {...}, notifier: {...} }, unavailable: [svc], links: {...} }` |

Calls `KitchenAdminService.GetQueueCounts`, `PlannerAdminService.GetQueueCounts`, `OffersAdminService.GetQueueCounts`,
`NotifierAdminService.ListAlerts(state=open)`.

## Data

None.

## Errors

None at route level (rule 3).

## Acceptance criteria

- [ ] With offers stopped, the page shows offers `unreachable` and offers queues "unavailable"; everything else is
      filled.
- [ ] Response time under 2.5 s with all services up.
