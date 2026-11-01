# GW-04-04 Invite preview and accept

| | |
|---|---|
| Service | gateway |
| Epic | [GW-04 Household, profiles and invite routes](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | GW-03-04, GW-01-03, IDN-07 (`GetInvitePreview`, `AcceptInvite`) |
| Unblocks | MOB-01 (token refresh after joining), MOB-03 |

## Goal

The person joining sees whose household it is before confirming, chooses whether to bring their own fridge items
along, and joins (features.md §2 "Linking", verification 15).

## Business rules

1. Both routes accept the code from a QR scan, a typed code or a link. The gateway normalises: take the part after
   the last `/` if a URL was given, remove spaces and dashes, uppercase. Must then be 6 characters from identity's
   alphabet (IDN-07); otherwise 400 `invalid_code` without calling identity.
2. Registered accounts only (decided, Q035): an anonymous account gets 403 `registration_required`. Rate limit
   bucket `invite_lookup`: 20 per user per hour
   (GW-03-04), against code guessing.
3. **Preview** calls identity `GetInvitePreview(code, actor_user_id)` and returns the household name, owner's display
   name, number of accounts and profiles, and `expiresAt`. It does not consume the invite.
4. **Accept** body `{ "bringFridge": true|false }` calls identity `AcceptInvite(code, actor_user_id, bring_fridge)`.
   Identity moves the membership and publishes `household.member_added` with `brought_fridge_from` (kitchen moves the
   items, KIT-13) and `household.member_removed` for the old household as its rules say.
5. Decided (Q187): a household that dropped from Plus to Free keeps its current members, but nobody can join it until
   it is Plus again; identity refuses the accept with `owner_not_plus` (→ 422), and pending invites can't be used.
6. Accept response includes the new `householdId` and `tokenRefreshRequired: true`. The gateway also drops cached
   entitlements of both households (`gw:ent:*`).
7. Idempotency required on accept: a retry after a lost response returns the stored 200 rather than
   `invite_used`.
8. Identity's refusals pass through (GW-01-02): `invite_expired` and `invite_used` → 410, `household_full` → 422,
   `already_member` → 409, `owner_not_plus` → 422.

## Contract

| Method | Path | Body | 2xx |
|---|---|---|---|
| GET | `/v1/invites/{code}` | — | 200 preview |
| POST | `/v1/invites/{code}/accept` | `{ "bringFridge": true }` | 200 accepted |

```json
preview  { "householdName": "Kowalscy", "ownerDisplayName": "Adam", "accounts": 2, "profiles": 1,
           "expiresAt": "2026-11-03T10:00:00Z" }
accepted { "householdId": "0192...", "memberId": "0192...", "broughtFridge": true, "tokenRefreshRequired": true }
```

Tag `household`, `auth: "registered"`.

## Data

Deletes `gw:ent:<old householdId>`, `gw:ent:<new householdId>`.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Malformed code | 400 | `invalid_code` |
| Unknown code | 404 | `not_found` |
| Expired / already used | 410 | `invite_expired` / `invite_used` |
| Household full, owner not Plus (also after a Plus → Free drop) | 422 | `household_full` / `owner_not_plus` |
| Already a member | 409 | `already_member` |
| Too many lookups | 429 | `rate_limited` |
| Anonymous | 403 | `registration_required` |

## Acceptance criteria

- [ ] `https://coolinary.app/i/k7m-2qx`, `K7M 2QX` and `k7m2qx` all resolve to the same invite.
- [ ] Preview does not consume the invite; accept does; a second account accepting the same code gets 410
      `invite_used` (verification 15).
- [ ] An invite older than 48 h gets 410 `invite_expired`.
- [ ] After accept and a token refresh, both phones see the same fridge and list (end-to-end, verification 15).
- [ ] The 21st lookup in an hour gets 429.

## Notes

- Decided (Q035): only registered accounts can join (anonymous accounts are deleted after 30 days unused, which
  would silently drop a member).
