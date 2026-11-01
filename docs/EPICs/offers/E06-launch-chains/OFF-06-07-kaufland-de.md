# OFF-06-07 Kaufland DE adapter (per store)

| | |
|---|---|
| Service | offers |
| Epic | [OFF-06 Launch chain adapters](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-06-03 (shared `kaufland` adapter) |
| Unblocks | OFF-08-02 for `kaufland-de` |

## Goal

Kaufland Germany's per-store flyers are read for the stores households use, with the code built for Kaufland Poland.

## Business rules

1. Source `kaufland-de`: adapter `kaufland` with German settings, list URL `https://filiale.kaufland.de/prospekte.html`,
   flyer URL pattern `https://leaflets.kaufland.com/de-DE/[id]/ar/[store]`, allowed hosts `filiale.kaufland.de`,
   `leaflets.kaufland.com` plus image hosts from the fixtures.
2. Editions per store, read on demand (OFF-05-07). `ListStores` from Kaufland DE's public store finder for
   OFF-03-05 (same gate rule as OFF-06-03 rule 4).
3. Flyers run Thursday to Wednesday. Daily check 06:00 `Europe/Berlin`; extra checks Tuesday 12:00 and Wednesday
   12:00.
4. `GATE: retailer terms checked` for kaufland.de and leaflets.kaufland.com.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] German fixtures pass the contract test with the shared adapter and no Poland-specific code paths.
- [ ] Two Kaufland DE stores with different editions show different offers once live (Verification 8).
