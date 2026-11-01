# PLN-01-06 PlannerAdminService scaffold and change log

| | |
|---|---|
| Service | planner |
| Epic | [PLN-01 Platform](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | — |
| Unblocks | PLN-03-03, PLN-03-07, PLN-03-08, PLN-03-09, PLN-11-03, ADM-04 |

## Goal

A separate gRPC service for back-office calls, so review and import RPCs are never reachable through the gateway, and
every admin change is recorded with who made it.

## Business rules

1. New `PlannerAdminService` in `contracts/coolinary/planner/v1/admin.proto`, served on the same port as
   `PlannerService`. The gateway has no client for it; only admin does (decision 9).
2. Every admin RPC requires metadata `x-admin-user-id`; missing → PERMISSION_DENIED.
3. Every admin RPC that changes data writes one `admin_changes` row in the same transaction: admin user id, RPC
   name, target (`dish:<id>`, `import:<id>`, ...), a short JSON summary (ids and counts, no free text from users).
4. Admin change rows are kept 1 year, then deleted by a monthly job. Decided (Q016).
5. `ListAdminChanges(target, page)` lets admin show the history of one dish.
6. `GetQueueCounts` returns one number per work queue for the admin overview (ADM-07-01), as a map with stable
   keys: `review_queue` (PLN-03-07), `dietitian_pending` (PLN-03-08), `effort_tags_pending` (PLN-03-06),
   `promotion_candidates` (PLN-11-03, R2), `open_batches` (`planner_batches` not yet processed, PLN-03-04). A queue
   whose task isn't built yet is left out of the map.

## Contract

```proto
service PlannerAdminService {
  rpc ListAdminChanges(ListAdminChangesRequest) returns (ListAdminChangesResponse);
  rpc GetQueueCounts(GetQueueCountsRequest) returns (GetQueueCountsResponse);   // rule 6
  // Review, import and promotion RPCs are added by PLN-03 and PLN-11.
}

message ListAdminChangesRequest {
  string target = 1;                         // e.g. "dish:3f2c..."
  coolinary.common.v1.PageRequest page = 2;
}
message AdminChange {
  string admin_user_id = 1;
  string rpc = 2;
  string target = 3;
  string summary_json = 4;
  string created_at = 5;                     // RFC 3339
}
message ListAdminChangesResponse {
  repeated AdminChange changes = 1;
  coolinary.common.v1.PageResponse page = 2;
}
message GetQueueCountsRequest {}
message GetQueueCountsResponse { map<string, int64> counts = 1; string computed_at = 2; }
```

## Data

`planner/migrations/000NN_admin_changes.sql`

```sql
CREATE TABLE admin_changes (
    id            bigserial PRIMARY KEY,
    admin_user_id text NOT NULL,
    rpc           text NOT NULL,
    target        text NOT NULL,
    summary       jsonb NOT NULL DEFAULT '{}',
    created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX admin_changes_target ON admin_changes (target, created_at DESC);
```

## Errors

| Case | gRPC status | Message |
|---|---|---|
| No `x-admin-user-id` | PERMISSION_DENIED | `admin only` |

## Acceptance criteria

- [ ] Calling any `PlannerAdminService` RPC without `x-admin-user-id` returns PERMISSION_DENIED.
- [ ] The gateway's generated client list does not include `PlannerAdminService`.
- [ ] Health and reflection (dev) list both services.
- [ ] With 4 versions in review and one unprocessed batch, `GetQueueCounts` returns `review_queue: 4`,
      `open_batches: 1`.
