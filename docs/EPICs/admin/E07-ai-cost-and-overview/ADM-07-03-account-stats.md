# ADM-07-03 Accounts and plans statistics

| | |
|---|---|
| Service | admin |
| Epic | [ADM-07 AI cost and system overview](README.md) |
| Release | R1 |
| Size | S |
| Depends on | ADM-01-05; IDN-10-04 (`IdentityAdminService.GetAccountStats`) |
| Unblocks | ADM-07-02 |

## Goal

Counts the founder tracks against the business-model gates (100 payers for Phase 1, 1,000 for Phase 2), and the
active-household numbers the AI guardrails need.

## Business rules

1. From identity's `AccountStats` (IDN-10-04): e-mail accounts and anonymous accounts; anonymous accounts created
   and upgraded in the last 7 days; e-mail accounts created in the last 30 days; households by plan (Free, Plus) and
   Plus set by hand; shared households and households with profiles; pending invites; households active in the last
   30 days by plan (a member's session used in that time); anonymous accounts removed by cleanup in the last 7 days.
   Sign-ups by country are not shown: identity doesn't know the country (it is in kitchen's profile).
2. Aggregates only; no ids or names.
3. Cached in admin for 5 minutes.
4. Role `admin` only; not audited.

## Contract

| Method | Path | Response |
|---|---|---|
| GET | `/admin/v1/stats/accounts` | `{ accounts: { email, anonymous, anonymous_created_7d, anonymous_upgraded_7d, email_created_30d }, households: { free, plus, plus_manual, shared, with_profiles }, active_30d: { free, plus }, pending_invites, anonymous_cleaned_7d, computed_at }` |

## Data

None.

## Errors

| Case | HTTP |
|---|---|
| identity unreachable | 503 |

## Acceptance criteria

- [ ] Creating an anonymous account raises `accounts.anonymous` by one after the cache expires.
