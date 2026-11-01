# OFF-12-01 Dino (PL) adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-12 Next chain adapters](README.md) |
| Release | R3 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-08-02 for `dino-pl` |

## Goal

Dino's weekly flyer is read automatically; Dino is strong in small towns where the launch chains are thin.

## Business rules

1. Source `dino-pl`: adapter `dino`, list URL `https://marketdino.pl/gazetki/lista/standard`, allowed hosts
   `marketdino.pl` plus image or PDF hosts found in the fixtures. Edition: national.
2. Rhythm: weekly, Wednesday to Tuesday. Daily check 06:00 `Europe/Warsaw`; extra check Tuesday 12:00.
3. `external_id`, dates and pages from the list page and flyer page; dates from the title when not in the data.
4. Only the "standard" list is read; other lists (themes) are skipped unless added to the source config.
5. `GATE: retailer terms checked` for marketdino.pl.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled, terms not approved.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] Fixtures pass the contract test with Wednesday–Tuesday dates.
- [ ] Two-week shadow run detects each flyer before its Wednesday start.
