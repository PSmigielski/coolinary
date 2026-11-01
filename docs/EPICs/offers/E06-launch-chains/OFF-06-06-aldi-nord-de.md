# OFF-06-06 Aldi Nord adapter (PDF)

| | |
|---|---|
| Service | offers |
| Epic | [OFF-06 Launch chain adapters](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-08-02 for `aldi-nord-de` |

## Goal

Aldi Nord's flyers (northern and eastern Germany) are read, using the PDF Aldi Nord offers, so a user near Hamburg
gets Aldi Nord offers (Verification 8).

## Business rules

1. Source `aldi-nord-de`: adapter `aldi_nord`, list URL `https://www.aldi-nord.de/prospekte.html`, flyer URL pattern
   `https://magazine.aldi-nord.de/aldi-nord/aldi-aktuell/[year]/…`, allowed hosts `www.aldi-nord.de`,
   `magazine.aldi-nord.de` plus the PDF host from the fixtures.
2. A PDF is offered (checked 4 October 2026). `FetchFlyer` returns the PDF when one is linked, else the page images.
   One PDF is cheaper to read than many images and is accepted by the extraction step (AI-13, OFF-07-02).
3. Edition: national only. Regional extras are not read in R2 (Decided, Q148; OFF-03-04 rule 2).
4. Flyers run Monday to Saturday, up to 14 days ahead: `accept_days_ahead = 14`.
5. Schedule: daily 06:00 `Europe/Berlin`; extra check Sunday 10:00.
6. `GATE: retailer terms checked` for aldi-nord.de and magazine.aldi-nord.de.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] Fixtures return a PDF for the current and the next week with correct dates.
- [ ] The stored version has the PDF's page count.
- [ ] Two-week shadow run detects every flyer before its start date.
