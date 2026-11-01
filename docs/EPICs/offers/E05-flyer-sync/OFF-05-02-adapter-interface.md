# OFF-05-02 Adapter interface and fixture harness

| | |
|---|---|
| Service | offers |
| Epic | [OFF-05 Flyer source framework and sync scheduler](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-01 |
| Unblocks | OFF-06-*, OFF-12-* |

## Goal

Each chain needs only a small adapter that knows its pages; everything else (scheduling, fetching rules, storage,
extraction) is shared. Adapters are tested offline against recorded pages, so a layout change shows up as a failing
fixture, not as wrong offers.

## Business rules

1. Interface (Go, `internal/flyers/adapter`):

   ```go
   type Adapter interface {
       Name() string
       // Flyers on the chain's list page for one edition.
       ListFlyers(ctx context.Context, f Fetcher, e Edition) ([]FlyerRef, error)
       // Page images or one PDF of a flyer.
       FetchFlyer(ctx context.Context, f Fetcher, ref FlyerRef) (FlyerContent, error)
       // Store list for per-store chains (OFF-03-05); others return ErrNotSupported.
       ListStores(ctx context.Context, f Fetcher) ([]ChainStore, error)
   }
   type FlyerRef struct {
       ExternalID, Title, PublicURL string
       ValidFrom, ValidTo           civil.Date // local, inclusive
       PageCount                    int        // 0 = unknown until fetched
   }
   type FlyerContent struct {
       Pages []Page   // image pages in order, or
       PDF   []byte   // one PDF
   }
   ```

2. Adapters fetch only through the given `Fetcher` (OFF-05-03); direct HTTP is not allowed (lint rule or code review
   checklist).
3. An adapter must return dates for every flyer. When it cannot read them (from the page data or the title), it
   returns the flyer with `ErrNoDates`; the flyer is not created and the source records a parse warning.
4. An adapter returns `ErrLayoutChanged` when the list page has none of the markers it expects (e.g. no flyer
   elements at all, or required fields missing in every item). That is a failure (OFF-05-06), never "no flyers".
5. Dates in titles are parsed with the chain's language and time zone; parsers live in a shared package
   (`od 14.10`, `ab Montag, 12.10.`, `12.10.–17.10.2026`), with the year inferred as the nearest future-or-current date
   within 60 days.
6. Fixtures: recorded responses under `internal/flyers/adapter/<name>/testdata/`, with a `record` command that
   refreshes them by hand (`make offers-record a=<name>`), using the polite fetcher. Tests never touch the network.
7. A shared contract test runs every adapter against its fixtures and checks: at least one flyer, all with
   non-empty id, dates, URL on an allowed host, and page fetch returning at least 1 page.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| `ErrNoDates` | That flyer skipped; warning in the check record |
| `ErrLayoutChanged` | Check failed (OFF-05-06) |
| `ErrNotSupported` from `ListStores` | Expected for chains that are not per-store |

## Acceptance criteria

- [ ] A dummy adapter used in tests passes the contract test; one that returns a flyer without dates fails it.
- [ ] The title-date parser handles the three example formats above and the year change (a flyer "od 29.12" seen on
      20 December gives the current year; "od 02.01" seen on 20 December gives the next year).
