# OFF-01 Platform: events, clients, plan lookup, jobs, admin service

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | Give offers the shared machinery every later epic uses: events in and out, clients for kitchen, ai and notifier, a local copy of each household's plan, a scheduled-job runner, and the `OffersAdminService` shell. |

## Flow

1. Business changes write rows and an outbox entry in one transaction; a relay publishes the outbox to `events.offers`.
2. A consumer group `offers` reads `events.identity` and `events.kitchen`; handlers are idempotent by event `id`.
3. Plan checks ask identity: `PlanFor` calls `IdentityService.GetEntitlements` per request (memoised within it) and per household per job run; offers keeps no plan copy (Decided, Q004).
4. Scheduled jobs (OSM refresh, flyer sync, publishing, spot-check sampling) run from one job table, safe if the container restarts.
5. Admin calls arrive on `OffersAdminService`, always with `x-admin-user-id`, and every change is recorded.

The outbox, publisher and consumer are copies of kitchen's Go code (KIT-01), kept identical per `docs/service-conventions.md`.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-01-01](OFF-01-01-event-outbox-publisher.md) | Event outbox and Redis Streams publisher | R2 | M | KIT-01 (outbox pattern) |
| [OFF-01-02](OFF-01-02-event-consumer.md) | Event consumer group and idempotent handlers | R2 | M | OFF-01-01, KIT-01 |
| [OFF-01-03](OFF-01-03-household-plan-copy.md) | Household plan lookup through identity | R2 | S | OFF-01-04, IDN-05 (`GetEntitlements`) |
| [OFF-01-04](OFF-01-04-service-clients.md) | Clients for identity, kitchen, ai and notifier | R2 | S | — |
| [OFF-01-05](OFF-01-05-scheduled-jobs.md) | Scheduled job runner | R2 | M | — |
| [OFF-01-06](OFF-01-06-admin-service-shell.md) | OffersAdminService shell and change log | R2 | S | — |
