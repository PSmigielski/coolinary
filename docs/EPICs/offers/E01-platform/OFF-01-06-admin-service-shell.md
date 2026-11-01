# OFF-01-06 OffersAdminService shell and change log

| | |
|---|---|
| Service | offers |
| Epic | [OFF-01 Platform](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | — |
| Unblocks | every admin task in OFF-02 … OFF-13, ADM-06 |

## Goal

The admin service has one place to call for everything about chains, flyers, shops and review queues, and every
admin change in offers can be traced to the person who made it.

## Business rules

1. A second connect-go handler `coolinary.offers.v1.OffersAdminService` is served on the same port 50051.
2. Every admin RPC requires metadata `x-admin-user-id` (non-empty). Without it: `PERMISSION_DENIED`. offers does not
   check the role itself; admin (ADM-02) does that before calling.
3. Every admin RPC that changes data inserts one row in `admin_changes` in the same transaction: admin user id, RPC
   name, entity type, entity id, a JSON summary of the change (ids and changed field names with old and new values;
   no user personal data).
4. `admin_changes` is kept for 2 years.
5. Read-only admin RPCs are not logged.
6. The admin handler is excluded from gateway access by network design: the gateway has no client for it.
7. `GetQueueCounts` returns one number per work queue for the admin overview (ADM-07-01), as a map with stable
   keys: `failing_sources` (OFF-05-06), `open_samples` (OFF-08-01), `held_offers` (OFF-07-06), `shop_reviews`
   (OFF-02-07), `store_link_queue` (OFF-03-05), `feed_store_queue` (OFF-13-01, R3). A queue whose task isn't built
   yet is left out of the map. Read-only, not logged (rule 5).

## Contract

`contracts/coolinary/offers/v1/admin.proto`

```proto
syntax = "proto3";
package coolinary.offers.v1;

// Called only by the admin service. Requires x-admin-user-id metadata.
service OffersAdminService {
  rpc ListJobs(ListJobsRequest) returns (ListJobsResponse);
  rpc TriggerJob(TriggerJobRequest) returns (TriggerJobResponse);
  rpc ListAdminChanges(ListAdminChangesRequest) returns (ListAdminChangesResponse);
  rpc GetQueueCounts(GetQueueCountsRequest) returns (GetQueueCountsResponse);   // rule 7
  // Further RPCs are added by the tasks that need them.
}

message ListAdminChangesRequest {
  string entity_type = 1;   // optional filter: "chain", "source", "offer", "shop", ...
  string entity_id = 2;     // optional
  coolinary.common.v1.PageRequest page = 3;
}

message GetQueueCountsRequest {}
message GetQueueCountsResponse { map<string, int64> counts = 1; string computed_at = 2; }
```

## Data

`offers/migrations/000NN_admin_changes.sql`

```sql
CREATE TABLE admin_changes (
    id            bigserial PRIMARY KEY,
    admin_user_id text NOT NULL,
    rpc           text NOT NULL,
    entity_type   text NOT NULL,
    entity_id     text NOT NULL,
    summary       jsonb NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX admin_changes_entity ON admin_changes (entity_type, entity_id, created_at DESC);
```

## Errors

| Case | gRPC status |
|---|---|
| Missing `x-admin-user-id` | `PERMISSION_DENIED` |

## Acceptance criteria

- [ ] `grpcurl` without the metadata gets `PERMISSION_DENIED` on every admin RPC.
- [ ] `TriggerJob` with the metadata writes one `admin_changes` row naming the admin and the job.
- [ ] `rpc_service` label shows `coolinary.offers.v1.OffersAdminService` for admin calls.
- [ ] With 2 held offers and 1 open sample, `GetQueueCounts` returns `held_offers: 2`, `open_samples: 1`.
