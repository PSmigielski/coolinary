# OFF-05-01 Flyer sources, flyers and versions

| | |
|---|---|
| Service | offers |
| Epic | [OFF-05 Flyer source framework and sync scheduler](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-02-01, OFF-03-04 |
| Unblocks | OFF-05-02, OFF-05-04, OFF-06-*, OFF-07-01 |

## Goal

One model describes where each chain's flyers come from, which flyers exist, and which version of each flyer was
read, so a flyer that changes is read again and nothing is read twice.

## Business rules

1. A source belongs to one chain and names its adapter (`biedronka`, `lidl_pl`, ...), list-page URL, allowed hosts
   (redirects and page images may only come from these), local check time (default 06:00 in the chain's time zone),
   extra check slots (weekday + local time), and how far ahead flyers are accepted (default 21 days, Lidl DE
   publishes three weeks ahead).
2. A source can only be `enabled` when `terms_approved = true`; approval records the admin, the date and a note
   (which terms page was read). Disabling never needs approval. (`GATE:` per chain in OFF-06 and OFF-12.)
3. A flyer is one flyer of one edition, identified by the chain's own id (`external_id`) and carries title, start and
   end date (inclusive, local dates), and the public URL where users can see it.
4. A flyer version is one reading of it: page count, a fingerprint of the page list, fetch state and extraction
   state. Version 1 is created when the flyer is first seen; a new version when its dates or page count change
   (OFF-05-04).
5. A flyer's `current_version_id` is the version whose offers are (or will be) published.
6. Flyer states: `detected`, `fetched`, `extracting`, `extracted`, `published`, `expired`, `withdrawn`, `failed`.
   Only the pipeline tasks change them: OFF-05-04 (detected, withdrawn), OFF-05-05 (fetched), OFF-07-02 (extracting,
   extracted, failed), OFF-07-07 (published, expired).

## Contract

`contracts/coolinary/offers/v1/flyers.proto` — messages used by admin and status RPCs:

```proto
message FlyerSource {
  string id = 1;
  string chain_id = 2;
  string adapter = 3;
  string list_url = 4;
  bool enabled = 5;
  bool terms_approved = 6;
  string terms_note = 7;
  SourceStatus status = 8;
  string last_success_at = 9;
  string last_error = 10;
}
enum SourceStatus { SOURCE_STATUS_UNSPECIFIED = 0; SOURCE_STATUS_OK = 1; SOURCE_STATUS_STALE = 2; SOURCE_STATUS_FAILING = 3; SOURCE_STATUS_DISABLED = 4; }

message Flyer {
  string id = 1;
  string edition_id = 2;
  string external_id = 3;
  string title = 4;
  string valid_from = 5;        // YYYY-MM-DD, local
  string valid_to = 6;
  string public_url = 7;
  FlyerState state = 8;
  int32 current_version = 9;
  int32 page_count = 10;
}
```

## Data

`offers/migrations/000NN_flyers.sql`

```sql
CREATE TABLE flyer_sources (
    id              text PRIMARY KEY,                 -- "lidl-de"
    chain_id        text NOT NULL REFERENCES chains (id),
    adapter         text NOT NULL,
    list_url        text NOT NULL,
    allowed_hosts   text[] NOT NULL,
    check_time      time NOT NULL DEFAULT '06:00',
    extra_checks    jsonb NOT NULL DEFAULT '[]',      -- [{"weekday":2,"time":"12:00"}]
    accept_days_ahead smallint NOT NULL DEFAULT 21,
    enabled         boolean NOT NULL DEFAULT false,
    terms_approved  boolean NOT NULL DEFAULT false,
    terms_note      text NOT NULL DEFAULT '',
    terms_approved_by text,
    terms_approved_at timestamptz,
    status          text NOT NULL DEFAULT 'disabled' CHECK (status IN ('ok', 'stale', 'failing', 'disabled')),
    last_success_at timestamptz,
    last_error      text,
    consecutive_failures smallint NOT NULL DEFAULT 0,
    CHECK (NOT enabled OR terms_approved)
);
CREATE TABLE flyers (
    id           uuid PRIMARY KEY,
    source_id    text NOT NULL REFERENCES flyer_sources (id),
    edition_id   text NOT NULL REFERENCES editions (id),
    external_id  text NOT NULL,
    title        text NOT NULL DEFAULT '',
    valid_from   date NOT NULL,
    valid_to     date NOT NULL CHECK (valid_to >= valid_from),
    public_url   text NOT NULL,
    state        text NOT NULL DEFAULT 'detected',
    current_version_id uuid,
    missing_checks smallint NOT NULL DEFAULT 0,
    first_seen_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (edition_id, external_id)
);
CREATE TABLE flyer_versions (
    id          uuid PRIMARY KEY,
    flyer_id    uuid NOT NULL REFERENCES flyers (id) ON DELETE CASCADE,
    version     int NOT NULL,
    page_count  int NOT NULL DEFAULT 0,
    fingerprint text NOT NULL DEFAULT '',
    valid_from  date NOT NULL,
    valid_to    date NOT NULL,
    fetch_state text NOT NULL DEFAULT 'pending' CHECK (fetch_state IN ('pending', 'done', 'failed')),
    extract_state text NOT NULL DEFAULT 'pending' CHECK (extract_state IN ('pending', 'submitted', 'done', 'failed')),
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (flyer_id, version)
);
ALTER TABLE flyers ADD FOREIGN KEY (current_version_id) REFERENCES flyer_versions (id);
```

## Errors

| Case | Behaviour |
|---|---|
| Enable without terms approval | `FAILED_PRECONDITION` (RPC); CHECK violation (database) |
| Same `(edition, external_id)` again | Upsert, no error |

## Acceptance criteria

- [ ] Enabling a source without terms approval fails at the database and at the RPC.
- [ ] One flyer id seen in two editions creates two flyers.
