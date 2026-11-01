# OFF-12-06 Penny DE adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-12 Next chain adapters](README.md) |
| Release | R3 |
| Size | M |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-08-02 for `penny-de` |

## Goal

Penny's weekly offers are read, including next week's from Friday.

## Business rules

1. Source `penny-de`: adapter `penny`, list URL `https://www.penny.de/angebote`, allowed hosts from the fixtures.
2. Rhythm: Monday to Saturday; next week shown from Friday (features.md §3). Daily check 06:00 `Europe/Berlin`;
   extra check Friday 12:00.
3. Penny has extra local offers per store. Decided (Q227): only the national offers are read (edition national);
   local offers are skipped.
4. If the page lists offers as structured web items instead of flyer pages, the adapter returns them as structured
   offers (no AI reading), and they go through normalisation, matching and validation like feed offers (OFF-13-02).
   This needs an optional adapter method `ListOffers` added to the interface of OFF-05-02.
5. `GATE: retailer terms checked` for penny.de.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled, terms not approved.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] A Friday fixture yields this week and next week.
- [ ] Two-week shadow run detects next week's offers by Saturday.
