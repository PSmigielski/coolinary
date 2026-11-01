# GW-04-03 Invites: create, list, cancel

| | |
|---|---|
| Service | gateway |
| Epic | [GW-04 Household, profiles and invite routes](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | GW-03-03, IDN-07 (`CreateInvite`, `ListInvites`, `CancelInvite`) |
| Unblocks | MOB-03 |

## Goal

The owner creates an invite that can be shown as a QR code, read out as a 6-character code, or sent as a link
through any messaging app; waiting invites are listed and can be cancelled (features.md §2 "Linking").

## Business rules

1. Plus-only (`plus: "shared_household"`) and registered only. Who may invite (owner only, or any adult member) is
   identity's rule (IDN-07).
2. Identity creates the invite: one person, valid 48 hours, a 6-character code. The gateway builds:
   - `link` = `INVITE_LINK_BASE` + `/` + code (config, e.g. `https://coolinary.app/i`);
   - `qrPayload` = the same link (the app renders the QR code).
3. `GET` lists waiting (not used, not expired, not cancelled) invites, newest first, at most 10.
4. `DELETE` cancels a waiting invite; cancelling a used or expired one → 404.
5. Household full (5 accounts, counting waiting invites as identity defines) → identity `FAILED_PRECONDITION`
   `household_full` → 422.

## Contract

| Method | Path | Body | 2xx |
|---|---|---|---|
| POST | `/v1/household/invites` | `{}` | 201 invite |
| GET | `/v1/household/invites` | — | 200 `{ "invites": [invite] }` |
| DELETE | `/v1/household/invites/{inviteId}` | — | 204 |

```json
invite { "id": "0192...", "code": "K7M2QX", "link": "https://coolinary.app/i/K7M2QX",
         "qrPayload": "https://coolinary.app/i/K7M2QX", "expiresAt": "2026-11-03T10:00:00Z",
         "createdByMemberId": "0192..." }
```

Tag `household`, `auth: "registered"`, idempotency required on POST.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free | 403 | `plus_required` |
| Not allowed to invite | 403 | `forbidden` |
| Household full | 422 | `household_full` |
| Unknown or not waiting | 404 | `not_found` |

## Acceptance criteria

- [ ] A Plus owner creates an invite; the code is 6 characters and the link ends with it; `expiresAt` is 48 h later.
- [ ] `GET` lists it; after `DELETE` it is gone and accepting it fails (GW-04-04).
- [ ] A Free owner gets 403 `plus_required` (features.md §2: "On Free, 'Invite someone' opens the Plus offer").

## Notes

- Decided (Q181): invite links are `https://coolinary.app/i/<code>` (config `INVITE_LINK_BASE`); the domain must be
  bought and serve the universal-link and app-link files before R2. Decided (Q187): a household that dropped to
  Free keeps its members but can't create invites (403 `plus_required`) until it is Plus again.
