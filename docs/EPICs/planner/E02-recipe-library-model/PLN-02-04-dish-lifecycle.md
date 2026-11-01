# PLN-02-04 Dish lifecycle and versions

| | |
|---|---|
| Service | planner |
| Epic | [PLN-02 Recipe library model](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-02-01, PLN-02-02 |
| Unblocks | PLN-03-03, PLN-03-07, PLN-08-04, PLN-09, PLN-11-02 |

## Goal

Clear states for every dish so users only ever see checked dishes, editing a published dish never changes what a
household already planned or cooked, and household AI dishes stay private.

## Business rules

1. Library states: `draft` → `review` → `published` → `retired`. Allowed moves: draft→review, review→draft
   (rejected with a reason), review→published, published→retired, retired→published (restore). Any other move is
   FAILED_PRECONDITION.
2. Household AI dishes have status `household` from creation and never enter the library states; PLN-11 copies them
   into a new library dish instead.
3. A published dish is never edited in place. Editing creates `draft_version = current_version + 1`, copied from the
   current version; users keep seeing `current_version` until the draft is published, which sets
   `current_version = draft_version`, `draft_version = NULL`.
4. Only one open draft version per dish. A second edit request edits the same draft.
5. Retiring hides a dish from browse, suggestions and AI candidate lists. Plans, history, favourites and ratings that
   point at it still read it; favourites show it as "no longer available".
6. Versions are never deleted while any plan, history or rating row references them. A daily job deletes rejected
   draft versions older than 90 days that nothing references.
7. Household AI dishes are deleted with their household (`household.deleted`, PLN-21) unless they were copied into
   the library (the library copy holds no household id).
8. All state moves of library dishes are admin RPCs (PLN-03-07) and write `admin_changes` (PLN-01-06).

## Contract

Internal Go package `dishes` (state machine + version copy). The RPCs that call it are in PLN-03-07.

## Data

No schema change beyond PLN-02-01 (`status`, `current_version`, `draft_version`, timestamps).

## Errors

| Case | gRPC status | Message |
|---|---|---|
| Move not allowed | FAILED_PRECONDITION | `cannot move dish from <a> to <b>` |
| Publish with blocking validation problems | FAILED_PRECONDITION | `dish not publishable`, problems attached |

## Acceptance criteria

- [ ] Editing a published dish leaves `GetDish` (no version given) returning the old content until the draft is published.
- [ ] A plan made before a publish still shows the old version afterwards.
- [ ] retired→published restores the dish to browse results.
- [ ] The cleanup job keeps a rejected draft version that a rating references.
