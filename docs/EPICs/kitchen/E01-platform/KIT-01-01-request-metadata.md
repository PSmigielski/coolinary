# KIT-01-01 Request metadata and household guard

| | |
|---|---|
| Service | kitchen (shared platform code; copy to planner and offers) |
| Epic | [KIT-01 Platform](README.md) |
| Release | R1 |
| Size | S |
| Depends on | — |
| Unblocks | KIT-01-06, KIT-01-07, KIT-01-08, every user-scoped kitchen RPC |

## Goal

Every RPC can read who is calling, for which household, on which plan and in which language, and outgoing calls
pass the same values on. A request can never act on a household other than the caller's.

## Business rules

1. The server and client interceptors in `internal/platform/rpc` forward these metadata keys: `x-request-id`,
   `traceparent`, `x-user-id`, `x-household-id`, `x-plan`, plus the new `x-language` (README D4) and
   `x-admin-user-id` (admin calls only).
2. Helpers return typed values from `ctx`: `rpc.UserID`, `rpc.HouseholdID`, `rpc.Plan` (`free` when missing or
   unknown), `rpc.Language` (default `en`), `rpc.AdminUserID` (empty when missing).
3. `rpc.Language` accepts a BCP 47 tag and keeps only the primary subtag, lowercased (`pl-PL` → `pl`). Values
   outside `pl`, `de`, `en` are passed through unchanged so language-aware code can fall back to `en` itself.
4. Household guard `rpc.CheckHousehold(ctx, requestHouseholdID)`: if both the request field and `x-household-id`
   are set and differ, the call fails with `PERMISSION_DENIED`. If only one is set, that one is used. If neither is
   set on a user-scoped RPC, the call fails with `INVALID_ARGUMENT` ("household_id is required").
5. Service-to-service calls (planner → kitchen) may omit `x-household-id`; then the request field alone is used
   (README D3).
6. Household ids must be UUIDs; anything else is `INVALID_ARGUMENT`.
7. None of these values is written to logs except `request_id`; `household_id` may be logged only as the first 8
   characters of its SHA-256 hash.

## Contract

No proto change. Metadata keys only.

## Data

No schema change.

## Errors

| Case | Status | Message |
|---|---|---|
| Request field and metadata household differ | `PERMISSION_DENIED` | `household mismatch` |
| No household on a user-scoped RPC | `INVALID_ARGUMENT` | `household_id is required` |
| Household id not a UUID | `INVALID_ARGUMENT` | `household_id must be a UUID` |

## Acceptance criteria

- [ ] Table tests for every helper, including missing, malformed and mixed-case values.
- [ ] A kitchen → ai call made inside a request carries all six keys that were present on the incoming request
      (checked with a test server that echoes headers).
- [ ] `ResolveBarcode` with `household_id=A` and `x-household-id=B` returns `PERMISSION_DENIED`.
- [ ] The identical files are copied into planner and offers in the same change set or noted as a follow-up in
      PLN-01 / OFF-01.

## Notes

`docs/service-conventions.md` lists the forwarded keys; add `x-language` and `x-admin-user-id` there
(README "Existing docs to update").
