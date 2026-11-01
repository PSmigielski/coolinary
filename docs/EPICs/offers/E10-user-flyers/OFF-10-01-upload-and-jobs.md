# OFF-10-01 Upload photos or a PDF

| | |
|---|---|
| Service | offers |
| Epic | [OFF-10 User-supplied flyers](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-07-01 |
| Unblocks | OFF-10-02, OFF-10-03, GW-10, MOB-07 |

## Goal

A Plus household can hand over a flyer as photos or a PDF and get a job it can follow, without the server keeping the
files.

## Business rules

1. Plus only (plan from identity, `PlanFor`, OFF-01-03); Free → `FAILED_PRECONDITION` (`plus_required`). Not part of the Free photo-scan quota.
2. Input: household id, shop id (OFF-10-05), and either 1 to 10 images (JPEG, PNG, WebP; the phone resizes to at most
   2,000 px on the long side, MOB-07) or one PDF of at most 40 pages (Decided, Q022). Total upload at most 25 MB; the server accepts
   messages up to 32 MB for this RPC only.
3. The call creates a `user_flyers` row (state `processing`) and returns its id at once; extraction runs in the
   background in this process (OFF-10-03). The bytes live only in memory until the AI call returns. A restart during
   processing marks the job `failed` (`interrupted`); the user uploads again.
4. Limits per household: at most 2 jobs processing at once and at most 10 uploads per day (Decided, Q171: an abuse
   guard, stated in the terms). Across the service at most 4 extractions run at once; others wait in memory in arrival order.
5. `GetUserFlyer(id)` returns state (`processing`, `draft`, `confirmed`, `failed`, `expired`), failure reason, and the
   draft when ready. Only the owning household can read it.
6. `DeleteUserFlyer(id)`: a draft is discarded; a confirmed flyer's offers expire at once (OFF-10-04 rule 6).
7. Kept data: shop, dates, counts, state, creating user id. Rows are deleted 30 days after the flyer's end date
   (or 30 days after failure or draft expiry).

## Contract

`contracts/coolinary/offers/v1/user_flyers.proto`

```proto
rpc UploadUserFlyer(UploadUserFlyerRequest) returns (UserFlyer);
rpc GetUserFlyer(GetUserFlyerRequest) returns (UserFlyer);
rpc DeleteUserFlyer(DeleteUserFlyerRequest) returns (DeleteUserFlyerResponse);

message UploadUserFlyerRequest {
  string household_id = 1;
  string user_id = 2;
  string shop_id = 3;
  repeated UploadedPage images = 4;
  bytes pdf = 5;
}
message UploadedPage { bytes content = 1; string mime_type = 2; }

message UserFlyer {
  string id = 1;
  string shop_id = 2;
  UserFlyerState state = 3;
  string failure_reason = 4;          // "interrupted", "ai_unavailable", "no_offers_found", "too_large", ...
  UserFlyerDraft draft = 5;           // OFF-10-03
  string valid_from = 6;
  string valid_to = 7;
  int32 offer_count = 8;
}
enum UserFlyerState {
  USER_FLYER_STATE_UNSPECIFIED = 0; USER_FLYER_STATE_PROCESSING = 1; USER_FLYER_STATE_DRAFT = 2;
  USER_FLYER_STATE_CONFIRMED = 3; USER_FLYER_STATE_FAILED = 4; USER_FLYER_STATE_EXPIRED = 5;
}
```

## Data

`offers/migrations/000NN_user_flyers.sql`

```sql
CREATE TABLE user_flyers (
    id           uuid PRIMARY KEY,
    household_id text NOT NULL,
    created_by   text NOT NULL,              -- user id
    shop_id      text NOT NULL REFERENCES shops (id),
    source       text NOT NULL CHECK (source IN ('images', 'pdf', 'link')),
    pages        int NOT NULL,
    state        text NOT NULL DEFAULT 'processing'
                 CHECK (state IN ('processing', 'draft', 'confirmed', 'failed', 'expired')),
    failure_reason text,
    valid_from   date,
    valid_to     date,
    draft        jsonb,                      -- extracted offers waiting for confirmation, no images
    created_at   timestamptz NOT NULL DEFAULT now(),
    delete_after date
);
CREATE INDEX user_flyers_household ON user_flyers (household_id, created_at DESC);
ALTER TABLE offers ADD FOREIGN KEY (user_flyer_id) REFERENCES user_flyers (id) ON DELETE CASCADE;
```

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Free household | `FAILED_PRECONDITION` (`plus_required`) |
| No files, or images and PDF together | `INVALID_ARGUMENT` |
| More than 10 images, PDF over 40 pages, upload over 25 MB | `INVALID_ARGUMENT` (`too_large`) |
| Unsupported file type | `INVALID_ARGUMENT` |
| 2 jobs processing or 10 uploads today | `RESOURCE_EXHAUSTED` |

## Acceptance criteria

- [ ] A Plus household uploads 4 photos and gets an id with state `processing`; a Free household is refused.
- [ ] After processing, no image bytes exist in Postgres, Redis or on disk (checked by search for the file hash).
- [ ] Restarting offers during processing leaves the job `failed` with `interrupted`.
