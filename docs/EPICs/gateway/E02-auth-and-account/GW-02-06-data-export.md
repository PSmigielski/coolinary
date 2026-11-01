# GW-02-06 `POST /v1/me/export` (GDPR export)

| | |
|---|---|
| Service | gateway |
| Epic | [GW-02 Auth and account routes](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | GW-02-01, IDN-09 (`ExportUserData`), NTF-04 (`ExportUserData`), KIT-13, PLN-21, OFF-14 (`ExportHouseholdData`) |
| Unblocks | MOB-03 |

## Goal

A user downloads everything Coolinary stores about them and their household in one JSON file (GDPR right of
access and portability), assembled by the gateway from each service (D9).

## Business rules

1. Registered accounts only (decided, Q142): anonymous accounts get 403 `registration_required` and can delete
   their account instead.
2. The gateway calls in parallel, each with the `heavy` deadline (10 s), using each owner's own messages:
   identity `ExportUserData(actor_user_id)` → `json` (string, IDN-09-03); notifier `ExportUserData(user_id)` → typed
   fields (devices, preferences, delivery counts, NTF-04-03), which the gateway turns into JSON with the standard
   proto JSON mapping; kitchen `ExportHouseholdData(household_id, member_id)` → `json` + `schema_version`
   (KIT-13-07); planner `ExportHouseholdData(household_id, member_id, page_token)` → pages of `json_chunk` of at most
   1 MB, fetched until `next_page_token` is empty and joined in order (PLN-21-06); offers
   `ExportHouseholdData(household_id)` → `json` (OFF-14-03; offers only from R2, before that the call is skipped by
   config `EXPORT_SERVICES`).
3. Each service filters by visibility itself: shared household data, the requester's own profile and target,
   children's profiles when the requester is a parent, never another adult's target (features.md §2 "Privacy").
   The gateway does not filter or inspect contents.
4. All parts must succeed; any failure → 503 `export_unavailable` (no partial export).
5. Response is one JSON document, `Content-Disposition: attachment; filename="coolinary-export-<YYYY-MM-DD>.json"`:
   `{ "generatedAt", "format": "coolinary-export/1", "notes": [...], "identity": {...}, "notifier": {...},
   "kitchen": {...}, "planner": {...}, "offers": {...} }`. `notes` always includes: "Age, height, weight and goals are
   kept only on your phone and are not part of this export."
6. Limit: 3 exports per user per 24 h (rate-limit bucket `export`, GW-03-04). Maximum assembled size 50 MB; larger →
   500 logged as a bug (not expected).
7. Nothing is stored or cached; the body is streamed out.

## Contract

`POST /v1/me/export` — `auth: "registered"`, tag `account`.

200 `application/json` (attachment, shape above).

Owner RPCs as listed in rule 2 (shapes differ per service; the owners' tasks define them).

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Anonymous | 403 | `registration_required` |
| Any service fails or times out | 503 | `export_unavailable` |
| More than 3 in 24 h | 429 | `rate_limited` |

## Acceptance criteria

- [ ] Export of a test household contains sections from every configured service and the phone-only note.
- [ ] With planner stopped, the route returns 503 and nothing partial.
- [ ] The fourth export within 24 h gets 429 with `Retry-After`.
- [ ] A second adult's export does not contain the first adult's daily target (kitchen's filtering, checked end to end).

## Notes

- Decided (Q142): no export for anonymous accounts; nothing links them to a person, and deletion is always available.
- Decided (Q143): the export is a **synchronous JSON download**; one household's data is small.
- Decided (Q138): each service returns its own JSON, which the gateway nests under the service's name.
