# GW-04 Household, profiles and invite routes

Releases R1 (household read and rename) and R2 (profiles without accounts, invites, leaving, removing).

## Goal

The app can show and rename the household, add children and others without a phone as profiles, invite people by
QR code, 6-character code or link, and let people leave or be removed (features.md §2 "Linking people into a
household"). All rules live in identity; the gateway maps routes to identity RPCs and handles the JWT refresh that
follows a household change.

## Household change and the JWT

Joining, leaving and being removed change the user's `householdId`. Those responses carry
`tokenRefreshRequired: true`; the app fetches a new JWT before its next call (MOB-01). A removed member's old JWT
still names the old household for up to 15 minutes; decided (Q191): the gateway closes that window at once by
checking identity's revocation key on every request (GW-04-05 rule 6).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-04-01](GW-04-01-household-read-rename.md) | `GET` and `PATCH /v1/household` | R1 | S | GW-02-01, IDN-04 |
| [GW-04-02](GW-04-02-profiles-without-accounts.md) | Profiles without accounts | R2 | S | GW-03-03, IDN-06 |
| [GW-04-03](GW-04-03-invites-create-list-cancel.md) | Invites: create, list, cancel | R2 | M | GW-03-03, IDN-07 |
| [GW-04-04](GW-04-04-invite-preview-accept.md) | Invite preview and accept | R2 | M | GW-03-04, GW-01-03, IDN-07 |
| [GW-04-05](GW-04-05-leave-remove.md) | Leave the household, remove a member | R2 | S | GW-01-03, IDN-08 |
