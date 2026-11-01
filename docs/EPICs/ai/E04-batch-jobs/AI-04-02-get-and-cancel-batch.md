# AI-04-02 GetBatch and CancelBatch

| | |
|---|---|
| Service | ai |
| Epic | [AI-04 Batch jobs](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | AI-04-01 |
| Unblocks | AI-04-03 |

## Goal

Callers can see a batch's progress and stop one they no longer need (e.g. a flyer that was withdrawn).

## Business rules

1. `GetBatch(batch_id)` calls `client.messages.batches.retrieve` and maps `processing_status`: `in_progress` →
   `BATCH_STATUS_IN_PROGRESS`, `canceling` → `CANCELING`, `ended` → `ENDED`.
2. It returns the request counts (`processing`, `succeeded`, `errored`, `canceled`, `expired`) and the times
   `created_at`, `ended_at`, `expires_at` (processing deadline, 24 h after creation) and
   `results_available_until` (29 days after creation).
3. `CancelBatch(batch_id)` calls the cancel endpoint; already-ended batches return their status unchanged.
   Items finished before the cancel still have results.
4. Recommended caller polling: every 10 minutes; most batches end within 1 hour, the hard limit is 24 hours.
   A caller must fetch results before `results_available_until`.
5. Batch ids are checked against `^(msgbatch_|stub_batch_)[A-Za-z0-9_-]+$`.

## Contract

```proto
rpc GetBatch(GetBatchRequest) returns (GetBatchResponse);
rpc CancelBatch(CancelBatchRequest) returns (GetBatchResponse);

enum BatchStatus {
  BATCH_STATUS_UNSPECIFIED = 0;
  BATCH_STATUS_IN_PROGRESS = 1;
  BATCH_STATUS_CANCELING = 2;
  BATCH_STATUS_ENDED = 3;
}

message GetBatchRequest { string batch_id = 1; }
message CancelBatchRequest { string batch_id = 1; }

message GetBatchResponse {
  string batch_id = 1;
  BatchStatus status = 2;
  int32 processing = 3;
  int32 succeeded = 4;
  int32 errored = 5;
  int32 canceled = 6;
  int32 expired = 7;
  string created_at = 8;              // RFC 3339
  string ended_at = 9;
  string expires_at = 10;
  string results_available_until = 11;
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Malformed id | `INVALID_ARGUMENT`, `ai_bad_request` |
| Unknown id (404 from the API) | `NOT_FOUND` |

## Acceptance criteria

- [ ] Recorded responses for each status map correctly.
- [ ] Cancel on an ended batch returns `ENDED` without error.
