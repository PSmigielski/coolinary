# MOB-05-02 Applying live list changes

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-05 Live shopping list client](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | MOB-05-01, MOB-03-07, MOB-04-03, PLN-15 (list versions, `GetListChanges`), GW-08-01 (`GET /v1/list/changes`), GW-09 |
| Unblocks | later list screens |

## Goal

Changes from other phones are applied to the cached list in the right order; if anything was missed, the phone
fetches the whole list instead of showing a wrong one.

## Business rules

1. Notice shape (GW-09-02, docs/EPICs/README.md "Live notices"):

   ```ts
   interface ListChanged {
     kind: "list.changed";
     household_id: string;
     version: number;                       // list version after the change
   }
   ```
   Notices carry no line data. Messages with a `type` key are control messages (MOB-05-01), not notices.
2. Let `v` be the cached list version (from `GET /v1/list`, from `GET /v1/list/changes`, or from the version a list
   mutation returned):
   - `version <= v` → ignore (already applied, including this phone's own changes, whose response already set `v`);
   - `version > v` → fetch `GET /v1/list/changes?since=v` once (planner `GetListChanges`, GW-08-01) and apply the
     changed lines and tombstones; if the answer says `fullResync`, replace the cached list with the list it
     returns. Notices arriving during that fetch are only compared with the new `v` afterwards (no second fetch
     when the answer already covers them).
3. Own changes: a list mutation's response carries the new list version and sets `v`, so the echo of the phone's
   own change is ignored by rule 2 and the optimistic state does not flicker.
4. On every (re)connect and on foreground, fetch `GET /v1/list/changes?since=v` once, because notices sent while
   disconnected are lost (pub/sub has no history).
5. Notices for a household other than the current one are ignored.
6. An upsert that conflicts with a pending local edit of the same line: the pending edit stays in the overlay
   (MOB-04-04) and is resolved by the server when it replays (MOB-04-03).

## Contract

As in rule 1.

## Data

Writes to the list query cache only.

## Errors

| Case | Behaviour |
|---|---|
| Malformed notice | Ignore it and fetch `GET /v1/list/changes?since=v` |

## Acceptance criteria

- [ ] Simulated gap (versions 4 → 6) triggers exactly one `GET /v1/list/changes?since=4`.
- [ ] A `fullResync` answer replaces the cached list.
- [ ] Phone A's own tick doesn't flicker when its echo arrives.
- [ ] Two phones ticking different lines at the same moment both end with both lines ticked.
