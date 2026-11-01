# GW-04-01 `GET` and `PATCH /v1/household`

| | |
|---|---|
| Service | gateway |
| Epic | [GW-04 Household, profiles and invite routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-02-01, IDN-04 (`GetHousehold`, `UpdateHousehold`) |
| Unblocks | MOB-03 |

## Goal

The app shows the household with its members (who has an account, who is a profile, who is the owner, which one
is "me"), and the owner can rename it.

## Business rules

1. `GET` calls identity `GetHousehold(household_id)` (exists). Response adds `isMe` per member (member id equals
   the JWT `memberId`) and `myRole`. Errors move to the problem model (GW-01-02).
2. `PATCH` body `{ "name" }`: trimmed, 1–40 characters, no control characters. Calls identity
   `UpdateHousehold(household_id, actor_user_id, name)`. Only the owner may rename; identity answers
   `PERMISSION_DENIED` otherwise.
3. Renames are last-writer-wins: identity keeps no household version (IDN-04-05), so `PATCH` takes no `If-Match`.
4. Members' display names are personal data: returned only to members of the same household (identity checks the
   household id).

## Contract

`GET /v1/household` — `auth: true`, tag `household`.

```json
200 {
  "id": "0192...", "name": "Ewa's home", "plan": "free", "myRole": "owner",
  "members": [
    { "id": "0192...", "displayName": "Ewa", "role": "owner", "hasAccount": true, "isMe": true },
    { "id": "0192...", "displayName": "Ola", "role": "profile", "hasAccount": false, "isMe": false }
  ]
}
```

`PATCH /v1/household` — `auth: true`, idempotency optional.
Request `{ "name": "Kowalscy" }` → 200 same shape as `GET`.

RPC: `rpc UpdateHousehold(UpdateHouseholdRequest) returns (GetHouseholdResponse);` (identity, IDN-04).

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Name empty or > 40 chars | 400 | `validation` |
| Not the owner | 403 | `forbidden` |
| No household | 403 | `no_household` |

## Acceptance criteria

- [ ] `GET` marks exactly one member `isMe: true`.
- [ ] A member who is not the owner gets 403 on `PATCH`.
- [ ] Two renames in a row both succeed; `GET` returns the second name.
