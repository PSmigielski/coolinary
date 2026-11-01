# IDN-10-04 Account numbers

| | |
|---|---|
| Service | identity |
| Epic | [IDN-10 Admin RPCs and admin role](README.md) |
| Release | R1 |
| Size | S |
| Depends on | IDN-10-01 |
| Unblocks | ADM-07 |

## Goal

The founder sees at a glance how many people use the app and how many move from anonymous to a real account, the
first step of the business model's funnel.

## Business rules

1. `GetAccountStats` returns counts:
   - accounts with an e-mail; anonymous accounts;
   - anonymous accounts created in the last 7 days; anonymous accounts upgraded in the last 7 days;
   - households by plan (Free, Plus) and by plan source;
   - households with 2 or more accounts; households with at least one profile;
   - pending invites;
   - e-mail accounts created in the last 30 days;
   - households active in the last 30 days, by plan (Free, Plus): at least one member's session `updated_at` in the
     last 30 days. Admin divides AI cost by these numbers for the cost guardrails (ADM-07-02);
   - anonymous accounts removed by the cleanup job (IDN-03-04) in the last 7 days, counted from
     `household.deleted` rows with reason `anonymous_cleanup` in `outbox_events` (same 7-day window as rule 2).
2. "Upgraded in the last 7 days" is counted from `account.linked` rows in `outbox_events` (kept 7 days after
   publishing, IDN-01-02 rule 9).
3. The result is cached in redis-cache for 5 minutes.
4. Only counts; no ids, no e-mails.

## Contract

```proto
message GetAccountStatsRequest {}
message AccountStats {
  int64 email_accounts = 1;
  int64 anonymous_accounts = 2;
  int64 anonymous_created_7d = 3;
  int64 anonymous_upgraded_7d = 4;
  int64 households_free = 5;
  int64 households_plus = 6;
  int64 households_plus_manual = 7;
  int64 households_shared = 8;
  int64 households_with_profiles = 9;
  int64 pending_invites = 10;
  string computed_at = 11;
  int64 email_accounts_created_30d = 12;
  int64 households_active_30d_free = 13;
  int64 households_active_30d_plus = 14;
  int64 anonymous_cleaned_7d = 15;
}
```

## Data

No schema change.

## Errors

None beyond `PERMISSION_DENIED` (IDN-10-01).

## Acceptance criteria

- [ ] With fixture data (3 e-mail, 2 anonymous, 1 upgraded yesterday), the counts match.
- [ ] A second call within 5 minutes returns the same `computed_at`.
- [ ] A Free household whose only session was last used 31 days ago is not in `households_active_30d_free`; one
      used yesterday is.
