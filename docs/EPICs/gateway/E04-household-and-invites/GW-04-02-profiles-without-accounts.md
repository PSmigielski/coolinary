# GW-04-02 Profiles without accounts

| | |
|---|---|
| Service | gateway |
| Epic | [GW-04 Household, profiles and invite routes](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | GW-03-03, IDN-06 (`CreateProfile`, `DeleteProfile`), IDN-04 (`UpdateMember`) |
| Unblocks | GW-05-04, MOB-03 |

## Goal

A parent adds a child (or anyone without a phone) as a profile, for portions and goals only (features.md §2:
"These profiles do not count towards the 5").

## Business rules

1. `POST` body `{ "displayName" }` (1–30 characters, trimmed). Plus-only (`plus: "profiles"`, GW-03-03) and
   registered accounts only.
2. Who may add, rename and remove: identity's rule (IDN-06). The gateway passes `actor_user_id` from the JWT:
   `POST` → identity `CreateProfile(household_id, actor_user_id, display_name)`; `PATCH` → identity
   `UpdateMember(household_id, actor_user_id, member_id, display_name)`; `DELETE` → identity
   `DeleteProfile(household_id, actor_user_id, member_id)`.
3. Removing a profile publishes `household.member_removed` with `reason: profile_deleted` (identity); kitchen and
   planner erase the person's profile and portions from that event.
4. Profile count limit: identity's rule; over it → identity `FAILED_PRECONDITION` `profile_limit` → 422.
5. No age, birth date, height, weight or goal goes through these routes (body schema allows only `displayName`).

## Contract

| Method | Path | Body | 2xx |
|---|---|---|---|
| POST | `/v1/household/profiles` | `{ "displayName": "Ola" }` | 201 member object (as in GW-04-01) |
| PATCH | `/v1/household/profiles/{memberId}` | `{ "displayName": "Ola K." }` | 200 member object |
| DELETE | `/v1/household/profiles/{memberId}` | — | 204 |

Tag `household`, `auth: "registered"`, idempotency required on POST, optional on PATCH/DELETE.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free household | 403 | `plus_required` |
| Anonymous | 403 | `registration_required` |
| `memberId` is not a profile of this household | 404 | `not_found` |
| Limit reached | 422 | `profile_limit` |

## Acceptance criteria

- [ ] A Plus owner adds "Ola"; `GET /v1/household` lists her with `role: "profile"`, `hasAccount: false`.
- [ ] A Free owner gets 403 `plus_required` and identity is not called.
- [ ] A body with `{ "displayName": "Ola", "age": 8 }` is refused with 400.
- [ ] Deleting a profile returns 204 and `household.member_removed` is published (integration test).
