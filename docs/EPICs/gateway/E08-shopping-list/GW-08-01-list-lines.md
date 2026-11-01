# GW-08-01 Read the list; add, edit and delete lines

| | |
|---|---|
| Service | gateway |
| Epic | [GW-08 Shopping list routes](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-01-03, PLN-15 (`GetList`, `AddListLine`, `UpdateListLine`, `RemoveListLine`, `ClearTickedLines`, `GetListChanges`) |
| Unblocks | GW-08-02, GW-08-03, GW-09-02, MOB-03, MOB-04 |

## Goal

Read the list grouped as planner returns it, and change it by hand, safely from several phones at once and from the
offline queue.

## Business rules

1. `GET /v1/list?includeTicked=` → planner `GetList(household_id, include_ticked)`: `version` (integer, +1 on every
   change), lines, `updatedAt`. Lines: id, ingredient (or free-text name), amount, packs (whole packs when the
   catalogue knows the pack size), `sources` (`meal` | `cuisine_base` | `manual` with the adding member id |
   `offer`), state, ticked time and member, flags (e.g. `plan_changed`), and from R2 the offer badge (PLN-17).
2. `POST /v1/list/lines` body `{ "ingredientId" | "name", "quantity"? }` (one line per call) → planner
   `AddListLine(household_id, member_id, client_mutation_id, ingredient_id | free_text, amount)`. Without an amount
   planner uses one pack; free text is matched to an ingredient once; adding an ingredient that has an open line
   merges into it and returns that line (PLN-15-06 rules 1–3). Free text ≤ 100 characters, ≤ 300 active lines.
   Adding a line from an offer uses `POST /v1/list/offer-lines` (GW-10-03).
3. `PATCH /v1/list/lines/{lineId}` body `{ "quantity", "baseVersion"? }` → `UpdateListLine`. **Last writer wins**
   (PLN-15-08 rule 4): an edit with an older `baseVersion` is still applied and the response shows the line as it
   now is; the gateway asks for no `If-Match` (two people in a shop must never block each other). An edit on a
   removed line → 404 (the phone drops it).
4. `DELETE /v1/list/lines/{lineId}` → `RemoveListLine`; 204 with header `X-List-Version`. Removing a missing line →
   404 (the offline queue treats it as done).
5. `DELETE /v1/list/lines?state=ticked` → `ClearTickedLines`; 200 list. Planner also clears ticked lines 48 h after
   ticking.
6. `GET /v1/list/changes?since=<version>` → `GetListChanges(household_id, since_version)`: lines changed after that
   version (removed lines as tombstones) and the current version, or `fullResync: true` with the whole list when the
   gap is too old or too big (PLN-15-09). Used by the live client after a version gap (GW-09-02).
7. **Idempotency:** every write requires an idempotency key, passed to planner as `client_mutation_id`; planner
   keeps results for 7 days, so offline replays (MOB-04) change nothing twice. Mutations older than 7 days are
   refused by planner (`mutation_too_old` → 422) and the app resyncs.
8. Every write response returns the list `version` and the changed line, so the client can tell whether a live
   notice (GW-09) is its own change.
9. Free and Plus: same routes. Live sharing between members needs a shared household (Plus) by nature.

## Contract

| Method | Path | 2xx |
|---|---|---|
| GET | `/v1/list?includeTicked=true` | 200 `{ "version": 42, "updatedAt": "...", "lines": [line] }` |
| POST | `/v1/list/lines` | 201 `{ "version": 43, "line": line }` (200 when merged into an open line) |
| PATCH | `/v1/list/lines/{lineId}` | 200 `{ "version": 44, "line": line }` |
| DELETE | `/v1/list/lines/{lineId}` | 204 (header `X-List-Version: 45`) |
| DELETE | `/v1/list/lines?state=ticked` | 200 `{ "version", "lines": [line] }` |
| GET | `/v1/list/changes?since=42` | 200 `{ "version": 45, "changed": [line], "fullResync": false, "list": null }` |

```json
line { "id": "lin_...", "ingredientId": "ing_onion", "name": "cebula", "quantity": { "value": 1000, "unit": "g" },
       "packs": { "count": 1, "packSize": { "value": 1000, "unit": "g" } }, "state": "open",
       "sources": [ { "kind": "meal", "label": "wt. obiad: Kurczak curry", "planId": "pln_...", "meal": "dinner" } ],
       "tickedAt": null, "tickedByMemberId": null, "offer": null, "version": 43 }
```

Tag `list`, `auth: true`. RPCs: planner PLN-15-01, 15-06, 15-09.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| No ingredient and no name, or name > 100 characters | 400 | `validation` |
| Line not in this household | 404 | `not_found` |
| 300 active lines reached | 422 | planner's code |
| Mutation older than 7 days | 422 | `mutation_too_old` |
| Edit of a removed line | 404 | `not_found` |

## Acceptance criteria

- [ ] Adding "cebula" twice by hand gives one line with the summed quantity (the second call returns the same line id).
- [ ] Two phones edit the same line's quantity; the later write is kept and both see the same `version`
      afterwards.
- [ ] A replayed `POST` with the same idempotency key does not add a second line.
- [ ] Every line shows where it came from (features.md §2 "Each line shows where it came from").
- [ ] `GET /v1/list/changes?since=42` after three edits returns exactly the three changed lines and version 45.
