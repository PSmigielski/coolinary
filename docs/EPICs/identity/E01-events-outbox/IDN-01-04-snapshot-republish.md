# IDN-01-04 Household snapshot republish

| | |
|---|---|
| Service | identity |
| Epic | [IDN-01 Events: outbox and publisher](README.md) |
| Release | R1 |
| Size | M |
| Depends on | IDN-01-02, IDN-01-03, IDN-10-01 |
| Unblocks | ADM-05 (button "republish household snapshots"), consumer bootstraps in KIT-13, PLN-21, OFF-14 |

## Goal

A consumer that starts late (a new service, a lost consumer group, more than 7 days of downtime) can rebuild its
copy of households, members and plans without calling identity directly.

## Business rules

1. Admin RPC `RepublishHouseholdSnapshots` starts a background job. With an empty `household_ids` list it covers
   every household; otherwise only the listed ones (at most 1,000 ids).
2. Only one job runs at a time. A second call while one runs returns `ALREADY_EXISTS` with the running job id.
3. The job walks households in id order, 500 per transaction, and emits one `household.snapshot` per household
   through the outbox, so snapshots are ordered with the regular events of the same household.
4. Throughput is capped at 1,000 snapshots per second so the stream and consumers are not flooded.
5. Snapshot content: current plan and every member (account members and profiles) with role.
6. Progress (`total`, `done`, `started_at`, `finished_at`, `status`) is kept in redis-state key
   `identity:snapshot_job` for 7 days. `GetSnapshotJob` reads it.
7. A restart of identity stops the job; it is not resumed automatically. The admin starts it again. Consumers
   treat snapshots as upserts, so repeating is harmless.
8. Every call is logged with the admin user id from `x-admin-user-id`.

## Contract

`contracts/coolinary/identity/v1/admin.proto`:

```proto
rpc RepublishHouseholdSnapshots(RepublishHouseholdSnapshotsRequest) returns (SnapshotJob);
rpc GetSnapshotJob(GetSnapshotJobRequest) returns (SnapshotJob);

message RepublishHouseholdSnapshotsRequest { repeated string household_ids = 1; }
message GetSnapshotJobRequest { string job_id = 1; }
message SnapshotJob {
  string job_id = 1;
  string status = 2;          // "running", "finished", "failed", "stopped"
  int64 total = 3;
  int64 done = 4;
  string started_at = 5;      // RFC 3339
  string finished_at = 6;
}
```

## Data

No schema change.

## Events

Publishes `household.snapshot` (IDN-01-03).

## Errors

| Case | gRPC status |
|---|---|
| A job is running | `ALREADY_EXISTS` |
| More than 1,000 ids, or an id that isn't a UUID | `INVALID_ARGUMENT` |
| Caller is not an admin | `PERMISSION_DENIED` (IDN-10-01) |

## Acceptance criteria

- [ ] With 2,000 test households, a full run emits exactly 2,000 `household.snapshot` events and reports
      `done = total = 2000`.
- [ ] A second call during the run returns `ALREADY_EXISTS`.
- [ ] Snapshots for a household are published after any earlier event of that household (stream order check).
