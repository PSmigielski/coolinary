# OFF-06-05 Aldi Süd adapter

| | |
|---|---|
| Service | offers |
| Epic | [OFF-06 Launch chain adapters](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-05-02, OFF-05-03 |
| Unblocks | OFF-08-02 for `aldi-sued-de` |

## Goal

Aldi Süd's weekly flyers (southern and western Germany) are read, so a user near Munich gets Aldi Süd offers
(Verification 8).

## Business rules

1. Source `aldi-sued-de`: adapter `aldi_sued`, list URL `https://www.aldi-sued.de/de/angebote/prospekte.html`,
   flyer URL pattern `https://prospekt.aldi-sued.de/kw[week]-[year]-…`, allowed hosts `www.aldi-sued.de`,
   `prospekt.aldi-sued.de` plus image hosts from the fixtures. Edition: national (the chain itself only exists in its
   region).
2. `external_id` = the `kw[week]-[year]-…` path segment. When no dates are given, they come from the ISO week: Monday
   to Saturday of week `[week]` of `[year]`.
3. This week plus two weeks ahead are published (checked 4 October 2026): `accept_days_ahead = 21`.
4. Schedule: daily 06:00 `Europe/Berlin`; extra check Sunday 10:00.
5. `GATE: retailer terms checked` for aldi-sued.de.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled.

## Errors

| Case | Behaviour |
|---|---|
| Expected list markers missing | `ErrLayoutChanged` (OFF-05-02) |

## Acceptance criteria

- [ ] `kw42-2026-…` without other dates gives 12 to 17 October 2026.
- [ ] A week-1 flyer published in late December gets the next year's dates when `[year]` says so.
- [ ] Two-week shadow run detects every flyer before its start date.
