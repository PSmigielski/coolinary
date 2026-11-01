# OFF-05-08 Sources and flyers admin RPCs

| | |
|---|---|
| Service | offers |
| Epic | [OFF-05 Flyer source framework and sync scheduler](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-01, OFF-05-05, OFF-01-06 |
| Unblocks | ADM-06 |

## Goal

The founder can see every source's health, approve a chain's terms, change schedules, re-run a check and look at
any stored flyer page, all from the admin service.

## Business rules

1. `ListSources`: every source with chain, coverage, status, last success, last error, consecutive failures, and the
   last 10 checks.
2. `UpdateSource`: list URL, allowed hosts, check time, extra checks, days ahead, enabled. Enabling requires terms
   approval (OFF-05-01 rule 2).
3. `ApproveSourceTerms`: sets `terms_approved = true` with a required note (at least 20 characters: which page, which
   date, what it allows). `RevokeSourceTerms` sets it false and disables the source in the same transaction.
4. `TriggerSync`: runs a check now for a source, or for one edition of it. Refused while a check of that source is
   running.
5. `ListFlyers`: by chain, edition, state and date range, with version count, page count and offer counts by state.
6. `GetFlyerPage`: returns one stored page (bytes, MIME type) or the PDF plus page number. Returns `NOT_FOUND` once
   the file is deleted (OFF-05-05 rule 6).
7. `ListEditions`: by chain, with demand, wanted state and last check.
8. Every change writes `admin_changes`.

## Contract

`admin.proto`

```proto
rpc ListSources(ListSourcesRequest) returns (ListSourcesResponse);
rpc UpdateSource(UpdateSourceRequest) returns (FlyerSource);
rpc ApproveSourceTerms(ApproveSourceTermsRequest) returns (FlyerSource);
rpc RevokeSourceTerms(RevokeSourceTermsRequest) returns (FlyerSource);
rpc TriggerSync(TriggerSyncRequest) returns (TriggerSyncResponse);
rpc ListFlyers(ListFlyersRequest) returns (ListFlyersResponse);
rpc GetFlyerPage(GetFlyerPageRequest) returns (GetFlyerPageResponse);
rpc ListEditions(ListEditionsRequest) returns (ListEditionsResponse);

message GetFlyerPageResponse {
  bytes content = 1;
  string mime_type = 2;
  int32 pdf_page = 3;      // set when content is a PDF
}
```

`GetFlyerPage` responses can reach 10 MB; the offers server and the admin client raise the gRPC message limit to
16 MB for this RPC only.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Enable without approval | `FAILED_PRECONDITION` |
| Approval note too short | `INVALID_ARGUMENT` |
| Check already running | `FAILED_PRECONDITION` |
| Page deleted | `NOT_FOUND` |

## Acceptance criteria

- [ ] Revoking terms disables the source and the next scheduled check does not run.
- [ ] `GetFlyerPage` returns a stored JPEG page byte-identical to the file.
- [ ] Each change appears in `ListAdminChanges` with the admin's id.
