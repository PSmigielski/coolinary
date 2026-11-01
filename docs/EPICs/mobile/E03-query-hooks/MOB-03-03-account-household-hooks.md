# MOB-03-03 Account and household hooks

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R1 (me, household, rename), R2 (profiles, invites, leave, remove) |
| Size | M |
| Status | todo |
| Depends on | MOB-03-01, MOB-03-02, GW-02, GW-04 |
| Unblocks | later household screens |

## Goal

Typed hooks for everything about the account and the household, including the R2 sharing flow (features.md
"Linking people into a household").

## Business rules

1. Hooks and routes (assumed GW-02/GW-04 paths):

   | Hook | HTTP | Rel. | Notes |
   |---|---|---|---|
   | `useMe()` | `GET /v1/me` | R1 | user-scoped, not persisted (holds the e-mail) |
   | `useHousehold()` | `GET /v1/household` | R1 | persisted; members' display names and roles |
   | `useRenameHousehold()` | `PATCH /v1/household` `{ name }` | R1 | owner only; name 1–40 chars (GW-04-01) |
   | `useAddProfile()` | `POST /v1/household/profiles` `{ displayName }` | R2 | profile without account (child); does not count toward 5 |
   | `useEditProfile()` / `useDeleteProfile()` | `PATCH` / `DELETE /v1/household/profiles/{memberId}` | R2 | |
   | `useInvites()` | `GET /v1/household/invites` | R2 | not persisted |
   | `useCreateInvite()` | `POST /v1/household/invites` | R2 | returns `{ id, code, link, expiresAt }`; the code is 6 characters, valid 48 h, one person |
   | `useCancelInvite()` | `DELETE /v1/household/invites/{id}` | R2 | |
   | `useInvitePreview(code)` | `GET /v1/invites/{code}` | R2 | shows whose household before joining |
   | `useAcceptInvite()` | `POST /v1/invites/{code}/accept` `{ bringFridge: boolean }` | R2 | then `forceRefresh()` → `householdChanged` |
   | `useLeaveHousehold()` | `POST /v1/household/leave` `{ newOwnerMemberId? }` | R2 | the owner of a shared household must name a new owner; then `forceRefresh()`; the user starts with an empty fridge |
   | `useTransferOwnership()` | `POST /v1/household/owner` `{ memberId }` | R2 | owner only |
   | `useRemoveMember()` | `DELETE /v1/household/members/{memberId}` | R2 | owner only |
2. Invite codes are entered case-insensitively and normalised to upper case; links `coolinary://invite/{code}` and
   the https link `https://coolinary.app/i/{code}` (decided, Q181) resolve to the same preview. QR content is the https link (rendering is a frontend task).
3. On Free, `useCreateInvite()` is gated by `usePlusGate("shared_household")` (MOB-02-05); the server still answers
   `403 plus_required`. Decided (Q187): when a shared household drops from Plus to Free, its members stay and keep
   sharing; only new invites and joining are blocked (`plus_required`, `owner_not_plus`) until it is Plus again.
4. Decided (Q035): only registered accounts can join. `useAcceptInvite()` for an anonymous user returns
   `registration_required` without a request; the flow asks the user to sign in with an e-mail code first
   (MOB-01-05), then accepts.
5. Accept, leave and remove are not queued offline (they change which household the phone acts for).
6. Before `useAcceptInvite()` runs, the outbox for the current household is flushed (same rule as MOB-01-05 rule 6).
7. Errors surfaced by name: `invite_expired`, `invite_used`, `household_full`, `owner_not_plus`, `plus_required`,
   `registration_required`, `not_owner`, `owner_must_transfer`.

## Contract

As in the table; types from `schema.d.ts` aliases in `src/api/types/household.ts`.

## Data

Query cache per MOB-03-01 (household persisted, invites and me not).

## Errors

| Case | Behaviour |
|---|---|
| `410`/`409` invite expired or used | `conflict` with `code` from the table in rule 6 |
| Accept while offline | `offline`, nothing queued |

## Acceptance criteria

- [ ] Verification 15 from the app side: a second test account joins by code and by link; both phones then read the
      same `GET /v1/stock`; after removal the second phone moves to its own household on its next request.
- [ ] An expired invite gives `invite_expired`; a used one gives `invite_used`.
