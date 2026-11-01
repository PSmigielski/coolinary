# OFF-04-05 Closed and merged shops in my shops

| | |
|---|---|
| Service | offers |
| Epic | [OFF-04 My shops and shopping preferences](README.md) |
| Release | R2 |
| Size | S |
| Depends on | OFF-04-01, OFF-02-04 |
| Unblocks | — |

## Goal

When a shop closes or turns out to be a duplicate, households that use it keep a correct list without doing anything
except, for a closure, picking a replacement.

## Business rules

1. Merged shop: the entry moves to the kept shop at the same position (OFF-02-04 rule 3). If the household already has
   the kept shop, the better (lower) position stays and the other entry is removed.
2. Closed shop (OSM import, admin approval of a "closed" correction): the entry's state becomes `closed`. Closed
   entries do not count towards the Free limit of 2 and get no offers.
3. `ListMyShops` returns closed entries after active and inactive ones, so the app can say "Lidl, ul. Prosta 2 has
   closed — choose another shop".
4. A closed entry is removed when the household sets a new list without it, or after 30 days.
5. A shop that reopens (reappears in OSM) does not bring its closed entries back to active; the user re-adds it.

## Contract

No contract change.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Any step fails | The merge or close transaction rolls back and is retried; nothing is visible to users |

## Acceptance criteria

- [ ] A Free household whose shop 1 closes can add a new shop without hitting the limit.
- [ ] A merged node entry becomes the way entry at the same position.
