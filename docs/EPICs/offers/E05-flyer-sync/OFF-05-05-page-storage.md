# OFF-05-05 Flyer page storage and retention

| | |
|---|---|
| Service | offers |
| Epic | [OFF-05 Flyer source framework and sync scheduler](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-04 |
| Unblocks | OFF-07-02, OFF-08-01, OFF-05-08 |

## Goal

Each flyer version's pages are downloaded once and kept just long enough for extraction and review, then deleted
(README D6). Pages are never republished to users (features.md §3 "no flyer pages are republished").

## Business rules

1. Storage (Decided, Q153): a named Docker volume `offers-flyers` mounted at `/data/flyers` (the prod root
   filesystem is read-only; the volume is writable), path `<chain_id>/<flyer_id>/<version>/<page>.<ext>`.
2. A version's pages are fetched with `FetchFlyer` in page order. Each page is stored with its SHA-256, size,
   MIME type and pixel size. A PDF is stored as one file; its page count comes from the PDF.
3. Version fingerprint = SHA-256 over the ordered page hashes. If it equals the previous version's fingerprint, the
   new version is marked `duplicate`, its files are deleted, and the flyer keeps the previous version (no re-read,
   no AI cost).
4. Limits: at most 120 pages and 300 MB per version. Over the limit: version `failed` with reason `too_large`,
   alert.
5. A partial fetch (any page failed after retries) marks the version `failed`; nothing is extracted from a
   half-fetched flyer. The next check retries it.
6. Retention (Decided, Q010): every file gets `delete_after = valid_to + 90 days` (valid_to of its version).
   Withdrawn flyers: `withdrawn_at + 90 days`. A daily job at 02:00 UTC deletes due files and their rows' paths.
   With 90 days the launch chains are expected to need about 5–10 GB (to be measured in the shadow weeks); an alert
   `offers.storage` fires when the volume passes 80% full.
7. Files of a removed chain are deleted within 24 h (OFF-08-05).
8. Pages are readable only through the admin RPC `GetFlyerPage` (OFF-05-08) and by the extraction step; no public
   route serves them.
9. The volume is not included in database backups.

## Contract

No public RPC.

## Data

`offers/migrations/000NN_flyer_pages.sql`

```sql
CREATE TABLE flyer_pages (
    version_id   uuid NOT NULL REFERENCES flyer_versions (id) ON DELETE CASCADE,
    page_no      int NOT NULL,
    path         text,                  -- null once deleted
    sha256       text NOT NULL,
    bytes        int NOT NULL,
    mime_type    text NOT NULL,
    width        int,
    height       int,
    delete_after date NOT NULL,
    PRIMARY KEY (version_id, page_no)
);
CREATE INDEX flyer_pages_due ON flyer_pages (delete_after) WHERE path IS NOT NULL;
ALTER TABLE flyer_versions ADD COLUMN duplicate_of uuid REFERENCES flyer_versions (id);
```

## Errors

| Case | Behaviour |
|---|---|
| A page fails after retries | Version `failed`; retried at the next check |
| Over 120 pages or 300 MB | Version `failed` (`too_large`), alert |
| Disk write fails or volume full | Version `failed`, alert `offers.storage` |

## Acceptance criteria

- [ ] A 24-page fixture flyer is stored as 24 files with hashes; a re-fetch with identical pages becomes a duplicate
      version and costs no extraction.
- [ ] On `valid_to + 90 days` the files are still there; on `valid_to + 91 days` they are gone and `path` is null.
- [ ] Killing the fetch after 10 of 24 pages leaves the version `failed` and no extraction is submitted.
