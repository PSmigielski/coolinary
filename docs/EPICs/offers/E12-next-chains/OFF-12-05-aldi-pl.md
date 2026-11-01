# OFF-12-05 Aldi PL adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-12 Next chain adapters](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-08-02 for `aldi-pl` |

## Goal

Aldi Poland's flyers are read automatically.

## Business rules

1. Source `aldi-pl`: list URL `https://www.aldi.pl/gazetki.html`. features.md §3: "Page exists; details not yet
   checked". This task starts with a check of the page: flyer format (flipbook, PDF, images), URL patterns, rhythm,
   edition. The findings are written into this file and into docs/features.md §3 before the adapter is built.
2. If the page uses the same viewer as Aldi Nord (Aldi Poland is run by Aldi Nord), the `aldi_nord` adapter is reused
   with Polish settings.
3. Edition: national unless the check shows otherwise.
4. `GATE: retailer terms checked` for aldi.pl.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled, terms not approved.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] The check's findings are recorded.
- [ ] Fixtures pass the contract test; two-week shadow run detects every flyer before its start date.
