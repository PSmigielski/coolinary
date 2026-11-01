# OFF-06-04 Lidl DE adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-06 Launch chain adapters](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-05-02, OFF-05-03, OFF-06-02 (shared `lidl` adapter) |
| Unblocks | OFF-08-02 for `lidl-de` |

## Goal

Lidl Germany's flyers, published up to three weeks ahead, are found and read, so German users can plan next week's
shopping before the offers start (features.md §3 step 6 "next week").

## Business rules

1. Source `lidl-de`: adapter `lidl` with German settings, list URL `https://www.lidl.de/c/online-prospekte/s10005610`,
   flyer URL pattern `https://www.lidl.de/l/prospekte/[name]/ar/0`, allowed hosts `www.lidl.de` plus flipbook and
   image hosts from the fixtures. Edition: national.
2. `external_id` = `[name]`; dates from the flyer data or title ("gültig ab 12.10.", "12.10. – 17.10.2026").
3. Flyers run Monday to Saturday; three weeks are published ahead (checked 4 October 2026), so
   `accept_days_ahead = 21`.
4. Schedule: daily 06:00 `Europe/Berlin`; extra check Sunday 10:00 and Wednesday 10:00.
5. Non-food theme flyers skipped by title skip list (German keywords, source config).
6. `GATE: retailer terms checked` for lidl.de.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] Fixtures give three weeks of flyers; the two future weeks are detected and later shown as "next week"
      (OFF-09-01) once the chain is live.
- [ ] Contract test passes.
- [ ] Two-week shadow run detects every flyer before its start date.
