# PLN-15-09 Live change notices and GetListChanges

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-15-08 |
| Unblocks | GW-09, MOB-05 |

## Goal

The list is shared live between household members: every change is announced so open phones update within a second,
and a phone that missed notices can catch up from a version (features.md "Shopping list": "shared live").

## Business rules

1. After each committed list change, planner publishes on Redis pub/sub (redis-state) channel
   `live.household.<household_id>` the JSON
   `{"kind": "list.changed", "household_id": "...", "version": N, "list_id": "...", "occurred_at": "..."}`
   (the README "Live notices" shape plus `list_id` and `occurred_at`) (same `kind` key as kitchen's
   `stock.changed` notice on this channel). No item names, no member names.
2. Publishing failures are logged and counted (`live_publish_failures_total`) but never fail the mutation; phones
   catch up through rule 3.
3. `GetListChanges(since_version)` returns every line whose `version > since_version` (including tombstones) and
   the current list version.
4. If `since_version` is older than the oldest kept tombstone (7 days) or more than 500 lines changed, the response
   sets `full_resync = true` and returns the whole list.
5. `since_version` greater than the current version returns `full_resync = true` (the phone's state is wrong).
6. The gateway (GW-09) relays notices to WebSocket clients; planner knows nothing about sockets.
7. R1 publishes the same notices even though sharing starts in R2 (one phone may have the list open on two
   screens); the channel name is shared with other household notices.

## Contract

`contracts/coolinary/planner/v1/list.proto`

```proto
rpc GetListChanges(GetListChangesRequest) returns (ListChanges);
message GetListChangesRequest { string household_id = 1; int64 since_version = 2; }
message ListChanges {
  int64 version = 1;
  repeated ListLine changed = 2;      // includes LINE_STATE_REMOVED tombstones
  bool full_resync = 3;
  ShoppingList list = 4;              // set only when full_resync
}
```

## Data

No schema change (uses `list_lines.version` index).

## Errors

| Case | gRPC status / behaviour |
|---|---|
| `since_version` < 0 | `INVALID_ARGUMENT` |
| `since_version` older than kept tombstones, more than 500 changes, or greater than the current version | success with `full_resync = true` (rules 4–5) |
| No list yet for the household | `NOT_FOUND` |
| Redis publish fails | mutation still succeeds; counted in `live_publish_failures_total` (rule 2) |

## Acceptance criteria

- [ ] Ticking a line publishes exactly one notice with the new version.
- [ ] `GetListChanges(since = v)` after two changes returns the two changed lines and version v + 2.
- [ ] `since` older than 7 days returns `full_resync`.
- [ ] With Redis pub/sub down, mutations still succeed.
