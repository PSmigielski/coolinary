# OFF-12-02 Netto PL adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-12 Next chain adapters](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-08-02 for `netto-pl` |

## Goal

Netto Poland's Monday and Thursday flyers are read automatically.

## Business rules

1. Source `netto-pl`: adapter `netto_pl`, list URL `https://netto.pl/gazetka-netto/`, flyer pages
   `?leafletid=[id]&page=[n]`, allowed hosts `netto.pl` plus image hosts from the fixtures. Edition: national.
2. `external_id` = `[id]`; pages are fetched with `page=1…N`, N from the flyer data or until a page returns no image
   (at most 120).
3. Rhythm: Monday and Thursday flyers. Daily check 06:00 `Europe/Warsaw`; extra checks Saturday 10:00 and Tuesday
   10:00.
4. `GATE: retailer terms checked` for netto.pl.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled, terms not approved.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] Fixtures give both weekly flyers with all pages in order.
- [ ] Two-week shadow run detects every flyer before its start date.
