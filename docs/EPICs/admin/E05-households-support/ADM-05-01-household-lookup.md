# ADM-05-01 Household lookup

| | |
|---|---|
| Service | admin |
| Epic | [ADM-05 Households and plans support](README.md) |
| Release | R1 |
| Size | M |
| Depends on | ADM-02-03; IDN-10-02 (`IdentityAdminService.FindAccount`), IDN-04 (`IdentityService.GetHousehold`) |
| Unblocks | ADM-05-02, ADM-05-03 |

## Goal

When a user writes to support, the founder can find their household and see its plan and members, and nothing more
than support needs (no fridge, plans, goals or health data).

## Business rules

1. Exactly one search key: `email` (exact match after lower-casing and trimming), `user_id`, or `household_id`.
   Partial or wildcard e-mail search is not offered.
2. Result, from two identity calls:
   - `IdentityAdminService.FindAccount(key)` (IDN-10-02): the found account (user id, e-mail, created at, anonymous
     flag, last seen) and its household summary (id, name, the account's role, plan, plan source `default` /
     `manual` / `store`, plan expiry, account count, profile count, pending invites, sharing grace end). For a
     `household_id` key the account is the owner's.
   - `IdentityService.GetHousehold(household_id)`: the members (member id, display name, role owner / member /
     profile, has account).
   The country is not shown: identity doesn't hold it (it is in kitchen's profile).
3. Never returned: fridge contents, plans, lists, calorie targets, diets, location, e-mail addresses of other
   members.
4. No match → 404.
5. Every lookup is audited as a `read` (`identity.household.lookup`) with the key type and the found household id;
   an e-mail key is stored as "by e-mail", not the address (ADM-02-03 rule 5).
6. Role `admin` only.

## Contract

| Method | Path | Response |
|---|---|---|
| GET | `/admin/v1/households?email=` \| `?user_id=` \| `?household_id=` | `{ household: {...}, members: [...] }` |

Calls `IdentityAdminService.FindAccount` (IDN-10-02), then `IdentityService.GetHousehold` (IDN-04).

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| No key or more than one key | 400 |
| Not found | 404 |

## Acceptance criteria

- [ ] Searching `Ewa@Example.com ` finds the household of `ewa@example.com`.
- [ ] The response contains no fridge, plan or goal fields.
- [ ] Each lookup creates one `read` audit row without the e-mail address.
