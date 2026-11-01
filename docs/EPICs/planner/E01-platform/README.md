# PLN-01 Platform: events, clients, snapshots, plan lookup

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | Give every later planner epic the plumbing it needs: consume events from identity, kitchen and offers; call kitchen, offers, ai and notifier safely; read the kitchen profile and stock cheaply; know each household's plan in background jobs; expose an admin gRPC service. |

## Flow

```
events.identity ─┐                         ┌─> handlers (PLN-04, PLN-21, PLN-17 ...)
events.kitchen  ─┼─> consumer group "planner" ┤
events.offers   ─┘                         └─> processed_events (idempotency)

planner ──gRPC + deadlines──> kitchen (GetPlanningProfile, ListStock, ...)   cached in redis-cache
        ──────────────────> offers, ai, notifier
admin   ──────────────────> PlannerAdminService (x-admin-user-id)
```

Planner publishes no stream events in R1–R3 (live list updates use pub/sub, PLN-15), so it gets a consumer but no
outbox. If a later epic needs to publish, it copies KIT-01's outbox.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-01-01](PLN-01-01-event-consumer.md) | Event consumer (Redis Streams, Go) | R1 | M | KIT-01 (consumer pattern) |
| [PLN-01-02](PLN-01-02-service-clients.md) | Clients for kitchen, offers, ai, notifier | R1 | M | — |
| [PLN-01-03](PLN-01-03-profile-snapshot.md) | Planning-profile snapshot cache | R1 | M | PLN-01-01, PLN-01-02, KIT-10 (GetPlanningProfile) |
| [PLN-01-04](PLN-01-04-stock-snapshot.md) | Stock snapshot cache | R1 | S | PLN-01-01, PLN-01-02, KIT-07 (ListStock) |
| [PLN-01-05](PLN-01-05-household-plan-copy.md) | Household plan lookup from identity | R1 | S | PLN-01-02, IDN-05 (GetEntitlements) |
| [PLN-01-06](PLN-01-06-admin-service.md) | PlannerAdminService scaffold and change log | R1 | S | — |
