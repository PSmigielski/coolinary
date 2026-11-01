# KIT-01 Platform: events, plan from identity, admin service, batch jobs

| | |
|---|---|
| Service | kitchen |
| Release | R1 |
| Unblocks | every other kitchen epic; planner PLN-01 and offers OFF-01 copy this pattern |

## Goal

Give kitchen the shared machinery the feature epics need: forwarded request metadata, a transactional outbox that
publishes events to Redis Streams, a consumer that handles other services' events exactly once in effect, a local
copy of each household's plan for background jobs, the `KitchenAdminService` with change records, and a tracker for
AI batch jobs. The code lives in `internal/platform/*` (identical in kitchen, planner and offers, per
`.claude/rules/go.md`) or in `internal/kitchen` where it is kitchen-specific.

## Flow

```
RPC handler ── tx ──> business tables + outbox row ──commit──> relay ──XADD──> events.kitchen
events.identity ──XREADGROUP (group "kitchen")──> consumer ── tx ──> handler + processed_events ──XACK
admin ──(x-admin-user-id)──> KitchenAdminService ──> admin_changes
job ──SubmitBatch──> ai ... poll GetBatch ... GetBatchResults ──> kind handler
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-01-01](KIT-01-01-request-metadata.md) | Request metadata and household guard | R1 | S | — |
| [KIT-01-02](KIT-01-02-event-envelope-outbox.md) | Event envelope and outbox writer | R1 | M | — |
| [KIT-01-03](KIT-01-03-outbox-relay.md) | Outbox relay to Redis Streams | R1 | M | KIT-01-02 |
| [KIT-01-04](KIT-01-04-stream-consumer.md) | Stream consumer with idempotent handlers | R1 | L | KIT-01-02 |
| [KIT-01-05](KIT-01-05-dead-letters-lag.md) | Dead letters, lag metrics and alerts | R1 | S | KIT-01-03, KIT-01-04, NTF-05 |
| [KIT-01-06](KIT-01-06-household-plan-copy.md) | Household plan from identity | R1 | S | KIT-01-01, IDN-05 |
| [KIT-01-07](KIT-01-07-admin-service-scaffold.md) | KitchenAdminService scaffold and change records | R1 | M | KIT-01-01 |
| [KIT-01-08](KIT-01-08-ai-batch-tracker.md) | AI client and batch job tracker | R1 | M | KIT-01-01, AI-04 |

## Decisions (taken 1 November 2026)

- Q004 (KIT-01-06): kitchen keeps no local plan copy; it asks identity (`GetEntitlements`) on each request that needs
  the plan and per household per job run. (changed from the first draft)
- Q015 (KIT-01-02): each stream entry has two fields, `type` and `envelope`.
- Q042 (KIT-01-05): no automatic resync after a stream gap in R1; the alert links to a runbook.
