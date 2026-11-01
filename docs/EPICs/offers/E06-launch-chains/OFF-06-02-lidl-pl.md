# OFF-06-02 Lidl PL adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-06 Launch chain adapters](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-06-04, OFF-08-02 for `lidl-pl` |

## Goal

Lidl Poland's Monday and Thursday flyers are found and read from Lidl's own flyer page.

## Business rules

1. Source `lidl-pl`: adapter `lidl` (shared with Lidl DE, country-specific settings), list URL
   `https://www.lidl.pl/c/nasze-gazetki/s10008614`, allowed hosts: `www.lidl.pl` plus the flipbook and image hosts
   found while building the fixtures (written into the source config). Edition: national.
2. Flyers are shown in an online flipbook (checked 4 October 2026). The adapter reads flyer ids, titles, dates and page
   image URLs from the data the flipbook page delivers in its HTML or its JSON requests, without running JavaScript
   (OFF-05-03 rule 9). If that proves impossible, the chain stays `user_flyers_only` and this task records why.
3. `external_id` = Lidl's flyer id or slug from the flipbook URL; `public_url` = the flipbook URL.
4. Dates from the flyer data, else from the title ("od 12.10 do 17.10", "ważna od poniedziałku 12.10").
5. Schedule: daily 06:00 `Europe/Warsaw`; extra checks Saturday 10:00 and Tuesday 10:00 (before the Monday and
   Thursday starts).
6. Non-food theme flyers are skipped by the source's title skip list (see OFF-06-01 rule 5).
7. `GATE: retailer terms checked` for lidl.pl.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled.

## Errors

| Case | Behaviour |
|---|---|
| Flipbook data not found | `ErrLayoutChanged` |

## Acceptance criteria

- [ ] Fixtures give the Monday and the Thursday flyer of one week with correct dates and all pages.
- [ ] Contract test passes; no request in the recorded session runs scripts or uses a log-in.
- [ ] Two-week shadow run detects every flyer before its start date.
