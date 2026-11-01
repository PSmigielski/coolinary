# KIT-11-03 Daily kcal target: limits, consent, deletion

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-11 Per-person profiles](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | KIT-11-02 |
| Unblocks | KIT-11-06, PLN-12 (portions) |

## Goal

Store each person's daily kcal target — the only number derived from health data that leaves the phone — with the
floor that protects adults, a record of consent, and a delete that really deletes (features.md §2 "Each person's
daily target", "Privacy"; Verification 19 and 20).

## Business rules

1. The target is a whole number of kcal. The phone computes it (Mifflin–St Jeor × activity, −500 / +300, EFSA values
   for children, or "my own number"); the server receives only the result.
2. **Adult floor:** with `target_kind = adult`, a target below 1,200 kcal is refused on every route, "my own number"
   included: `FAILED_PRECONDITION` reason `target_below_minimum` (gateway → 422). Decided (Q193): upper bound
   5,000 kcal.
3. **Children:** with `target_kind = child`, the server accepts 500–4,000 kcal. The age-based limits (EFSA value at
   the lowest activity level; the per-age lower limit for a doctor's number) are enforced on the phone, which knows
   the age. Decided (Q194): the server's absolute child bounds are 500–4,000 kcal; a paediatric dietitian reviews
   them before R2 (features.md §5 open points) and may narrow them (config values).
4. Changing `target_kind` re-checks an existing target against the new bounds; if it no longer fits, the call fails
   with `target_below_minimum` (or `invalid_value` above the maximum) and nothing changes.
5. **Consent:** setting a target requires `health_consent` (version and time) given on the phone. Kitchen stores
   the version, the time and which member gave it (the caller: the person themself or a parent). No consent → 
   `FAILED_PRECONDITION` reason `consent_required`.
6. **No history:** only the current value is kept, in the member row; no audit table, no old values in logs, events
   or metrics (features.md: "no weight history is kept"; the same applies to targets).
7. **Deletion:** `clear_target = true` sets the target and the consent fields to null in the same transaction,
   deletes the household's cached planning snapshots in `redis-cache` (`kitchen:planprof:<household_id>:*`), bumps
   `planning_version` and writes `profile.changed`. After the call no server copy of the value exists
   (Verification 20).
8. **Free:** only the caller's own target (the one personal goal, features.md §6); KIT-11-02 rule 8.
9. Unit: always kcal on the server; kJ is a display option on the phone.

## Contract

Fields `daily_kcal_target`, `clear_target`, `health_consent` of `UpdateMemberProfile` (KIT-11-02).

## Data

`kitchen/migrations/000NN_member_target_bounds.sql`

```sql
-- +goose Up
ALTER TABLE member_profiles ADD CONSTRAINT member_target_bounds CHECK (
    daily_kcal_target IS NULL
    OR (target_kind = 1 AND daily_kcal_target BETWEEN 1200 AND 5000)
    OR (target_kind = 2 AND daily_kcal_target BETWEEN 500 AND 4000));

-- +goose Down
ALTER TABLE member_profiles DROP CONSTRAINT member_target_bounds;
```

The service checks the same bounds first so the error reason is precise; the constraint is the safety net.

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Adult target < 1,200 | `FAILED_PRECONDITION` | `target_below_minimum` |
| Child target < 500 | `FAILED_PRECONDITION` | `target_below_minimum` |
| Above the maximum | `INVALID_ARGUMENT` | `invalid_value` |
| No consent with a new target | `FAILED_PRECONDITION` | `consent_required` |
| Both a target and `clear_target` | `INVALID_ARGUMENT` | |

## Acceptance criteria

- [ ] Adult target 1,199 → `target_below_minimum`; 1,200 → saved (Verification 19).
- [ ] Child target 1,400 with consent given by a parent → saved with `consent_by_member_id` = the parent.
- [ ] A child with a target of 900 switched to `adult` → `target_below_minimum`, nothing changed.
- [ ] `clear_target` → the row has no target and no consent; `GetPlanningProfile` returns no target; no Redis key of
      the household holds the old value (Verification 20).
- [ ] grep of logs and the outbox after a test run finds no target values.
