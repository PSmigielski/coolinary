# OFF-12-04 Auchan PL adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-12 Next chain adapters](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-08-02 for `auchan-pl` |

## Goal

Auchan Poland's national flyer is read automatically.

## Business rules

1. Source `auchan-pl`: adapter `auchan_pl`, list URL `https://www.auchan.pl/pl/gazetki`, allowed hosts from the
   fixtures. Edition: national.
2. Auchan publishes national and per-store flyers (features.md §3). Decided (Q226): R3 reads only the
   national flyers; per-store flyers are skipped (the chain's `edition_model` stays `national`).
3. Daily check 06:00 `Europe/Warsaw`; extra checks set after the first shadow week.
4. `GATE: retailer terms checked` for auchan.pl.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled, terms not approved.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] Fixtures separate national from per-store flyers; only national ones are returned.
- [ ] Two-week shadow run detects every national flyer before its start date.
