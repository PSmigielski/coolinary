# OFF-05-04 Sync checks and change detection

| | |
|---|---|
| Service | offers |
| Epic | [OFF-05 Flyer source framework and sync scheduler](README.md) |
| Release | R2 |
| Size | L |
| Depends on | OFF-05-02, OFF-05-03, OFF-01-05 |
| Unblocks | OFF-05-05, OFF-05-06, OFF-07-02 |

## Goal

Every new flyer on a chain's list page is found before its start date, and a flyer that the chain changes is read
again (features.md §3 steps 1–2, Verification 7).

## Business rules

1. Each enabled source gets a job `flyer_check:<source_id>` at its daily `check_time` and at each extra check slot,
   in the chain's time zone. Admin can trigger it (`TriggerSync`, OFF-05-08); a store edition with new demand gets
   one too (OFF-05-07).
2. A check runs `ListFlyers` for each edition in scope: the national edition, or every wanted store/format edition
   (OFF-05-07).
3. Flyers ending before today, or starting more than `accept_days_ahead` days from today, are ignored.
4. For each returned flyer:
   - new `(edition, external_id)` → create the flyer (state `detected`) and version 1 with the reported dates and page
     count;
   - known flyer with different start date, end date or a different page count (when known) → new version; the flyer
     goes back to `detected` while its current published version stays visible until the new one is published
     (OFF-07-07);
   - known flyer, nothing changed → only `missing_checks = 0`.
5. A known flyer that is missing from the list page: `missing_checks + 1`. Decided (Q152): if it is missing in 2
   checks in a row and has not started yet, it becomes `withdrawn` and its unpublished offers are dropped; once
   started it is kept until its end date (chains often drop running flyers from the list early).
6. A check whose adapter returns `ErrLayoutChanged`, `ErrBlocked`, `ErrRobots` or a fetch failure changes no flyers
   (nothing is counted missing) and is a failed check for OFF-05-06.
7. Detected versions are fetched right after the check (OFF-05-05), one flyer at a time per source.
8. A check records: start, end, flyers seen, new, changed, missing, errors. Kept 90 days.
9. Metrics: `flyer_checks_total{source,result}`, `flyers_detected_total{source,kind}` (`new`, `changed`).

## Contract

No public RPC change. Admin `TriggerSync` in OFF-05-08.

## Data

`offers/migrations/000NN_flyer_checks.sql`

```sql
CREATE TABLE flyer_checks (
    id          bigserial PRIMARY KEY,
    source_id   text NOT NULL REFERENCES flyer_sources (id),
    edition_id  text REFERENCES editions (id),
    started_at  timestamptz NOT NULL DEFAULT now(),
    finished_at timestamptz,
    result      text CHECK (result IN ('ok', 'not_modified', 'failed')),
    seen int NOT NULL DEFAULT 0, new int NOT NULL DEFAULT 0,
    changed int NOT NULL DEFAULT 0, missing int NOT NULL DEFAULT 0,
    error       text
);
CREATE INDEX flyer_checks_recent ON flyer_checks (source_id, started_at DESC);
```

## Errors

| Case | Behaviour |
|---|---|
| Adapter error | Check `failed`, no flyer changes, counts towards source health |
| Same flyer reported twice in one list | Second one ignored |

## Acceptance criteria

- [ ] Fixture day 1 with flyers A and B, day 2 with A, B, C: C is created once; A and B unchanged.
- [ ] Fixture where B's end date moves by a day creates version 2 of B; the published version 1 stays visible.
- [ ] A fixture returning `ErrLayoutChanged` leaves all flyers untouched and marks the check failed.
- [ ] Over a two-week shadow run per launch chain, every flyer on the list page is detected before its start date
      (Verification 7, first part).
