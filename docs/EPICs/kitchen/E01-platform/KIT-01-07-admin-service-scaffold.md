# KIT-01-07 KitchenAdminService scaffold and change records

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-01 Platform](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-01-01 |
| Unblocks | KIT-01-05, KIT-02-09, KIT-02-10, KIT-03-06, KIT-05-05, KIT-06-06, ADM-03 |

## Goal

Kitchen exposes a separate gRPC service for the back office, reachable only from admin, and records who changed
what, so catalogue and nutrition edits can be traced and explained.

## Business rules

1. `coolinary.kitchen.v1.KitchenAdminService` is defined in `contracts/coolinary/kitchen/v1/admin.proto` and mounted
   in `main.go` on the same gRPC port as `KitchenService`. Feature epics add their RPCs to it.
2. An interceptor on this service rejects any call without a non-empty `x-admin-user-id` with `PERMISSION_DENIED`.
   Kitchen does not check the admin role itself; admin verifies it with identity (ADM-02). The gateway never routes
   to admin RPCs (GW-01).
3. Every admin RPC that changes data writes one `admin_changes` row in the same transaction: admin user id, action
   name (the RPC name), entity type and id, and the entity's JSON before and after (only catalogue fields; never
   household data).
4. Read-only admin RPCs write nothing.
5. `ListAdminChanges(entity_type, entity_id, page)` returns the history of one entity, newest first, 50 per page.
6. Admin list RPCs follow the paging rule in `.claude/rules/go.md`: default page size 50, maximum 200.
7. Rows in `admin_changes` are kept for 3 years, then deleted by a monthly job.
8. `GetQueueCounts` returns one number per work queue for the admin overview (ADM-07-01), as a map with stable
   keys: `match_reviews` (KIT-03-06), `nutrition_suggestions` (KIT-06-06), `unreviewed_names` (KIT-02-10),
   `product_reviews` (KIT-05-05), `parked_items` (items parked after 3 failed batch attempts, KIT-02-08/KIT-06-05),
   `open_batches` (KIT-01-08), `dead_letters` (KIT-01-05). A queue whose task isn't built yet is left out of the
   map. Each count is one indexed `count(*)`; the whole call answers within 200 ms on the starter catalogue.

## Contract

`contracts/coolinary/kitchen/v1/admin.proto`

```proto
service KitchenAdminService {
  rpc ListAdminChanges(ListAdminChangesRequest) returns (ListAdminChangesResponse);
  rpc GetQueueCounts(GetQueueCountsRequest) returns (GetQueueCountsResponse);   // rule 8
  // Feature RPCs are added by KIT-01-05, KIT-02-09, KIT-02-10, KIT-03-06, KIT-05-05, KIT-06-06.
}
message AdminChange { int64 id = 1; string admin_user_id = 2; string action = 3; string entity_type = 4;
  string entity_id = 5; string before_json = 6; string after_json = 7; string created_at = 8; }
message GetQueueCountsRequest {}
message GetQueueCountsResponse { map<string, int64> counts = 1; string computed_at = 2; }
```

## Data

`kitchen/migrations/000NN_admin_changes.sql`

```sql
CREATE TABLE admin_changes (
    id             bigserial PRIMARY KEY,
    admin_user_id  text NOT NULL,
    action         text NOT NULL,
    entity_type    text NOT NULL,
    entity_id      text NOT NULL,
    before         jsonb,
    after          jsonb,
    created_at     timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX admin_changes_entity ON admin_changes (entity_type, entity_id, created_at DESC);
```

## Errors

| Case | Status | Message |
|---|---|---|
| Missing `x-admin-user-id` | `PERMISSION_DENIED` | `admin only` |

## Acceptance criteria

- [ ] A call without `x-admin-user-id` to any `KitchenAdminService` RPC is refused.
- [ ] A sample write RPC (added with KIT-02-10) produces one `admin_changes` row with before/after.
- [ ] `make grpc s=kitchen args=list` shows both services in dev.
- [ ] With 3 open match reviews and 1 open batch, `GetQueueCounts` returns `match_reviews: 3`, `open_batches: 1`.
