# OFF-06-01 Biedronka (PL) adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-06 Launch chain adapters](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-08-02 for `biedronka-pl` |

## Goal

Biedronka's flyers are found and read automatically from Biedronka's own flyer page.

## Business rules

1. Source `biedronka-pl`: adapter `biedronka`, list URL `https://www.biedronka.pl/pl/gazetki`, allowed hosts
   `www.biedronka.pl`, `biedronka.pl`, `cdn.biedronka.pl`. Edition: national (`biedronka-pl:national`).
2. List page → one `FlyerRef` per flyer link of the form `biedronka.pl/pl/press,id,[id],title,[name]`. `external_id`
   = `[id]`; `public_url` = the link; title = the link text.
3. Dates: the start date is in the title (checked 4 October 2026). The end date is read from the flyer page when
   shown. Decided (Q160): when only a start date is found, the flyer's end date is start + 6 days, and the dates
   the AI reads on each offer may narrow it but never extend it (OFF-07-05).
4. Pages: page images on `cdn.biedronka.pl`, in the order the flyer page lists them.
5. Biedronka publishes several flyers a week. Flyers whose title matches the source's skip list (non-food themes,
   e.g. "dom", "ogród", "ubrania", "zabawki") are not fetched. The skip list is source config, edited in admin;
   the rule for non-food offers inside food flyers is in OFF-07-03.
6. Schedule: daily check 06:00 `Europe/Warsaw`; extra check every day at 18:00 (several flyers a week).
7. `GATE: retailer terms checked` — Biedronka's site terms (reuse of content, automated access) approved in admin
   with a note before the source is enabled.

## Contract

No gRPC change. Config row in `flyer_sources` (migration or admin).

## Data

`offers/migrations/000NN_source_biedronka.sql` inserts the source, disabled, terms not approved.

## Errors

| Case | Behaviour |
|---|---|
| No `press,id,` links on the list page | `ErrLayoutChanged` |
| A link without a readable start date | That flyer skipped with `ErrNoDates`; others continue |

## Acceptance criteria

- [ ] Fixture of the list page yields every flyer with id, title, start date and URL; the contract test passes.
- [ ] Fixture of one flyer yields all pages in order from `cdn.biedronka.pl`.
- [ ] A title like "Gazetka od poniedziałku 12.10" gives start 12 October of the right year.
- [ ] Two-week shadow run: every Biedronka flyer on the list page is detected before its start date.
