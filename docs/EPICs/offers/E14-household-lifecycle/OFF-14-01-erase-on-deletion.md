# OFF-14-01 Erase data on household or account deletion

| | |
|---|---|
| Service | offers |
| Epic | [OFF-14 Household lifecycle](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-01-02, IDN-08 / IDN-09 (household.deleted, account.deleted) |
| Unblocks | GDPR deletion end to end |

## Goal

When a household or an account is deleted, offers keeps nothing that belongs to it or points at the person.

## Business rules

1. `household.deleted` (one transaction): delete `my_shops`, `shopping_preferences`, `user_flyers`
   and their offers (published ones expire first, with `offers.expired`), and the household's pending and rejected
   user-added shops that no other household uses. Active user-added shops approved for everyone stay, with
   `added_by_household` set to null. Open shop corrections of the household are deleted.
2. `account.deleted`: `user_flyers.created_by` of that user is set to `deleted`; nothing else in offers is per user.
   Household data stays (the household may have other members).
3. Price observations, review samples and admin changes hold no household or user ids and are not touched.
4. Both handlers are idempotent (OFF-01-02) and finish within 1 minute of the event for a household with 100 user
   flyers.
5. A deletion audit row records event id, type and counts deleted (no ids of the person).

## Contract

Consumes `household.deleted`, `account.deleted` (docs/EPICs/README.md "Events").

## Data

`offers/migrations/000NN_deletion_log.sql`

```sql
CREATE TABLE deletion_log (
    event_id   uuid PRIMARY KEY,
    event_type text NOT NULL,
    counts     jsonb NOT NULL,
    done_at    timestamptz NOT NULL DEFAULT now()
);
```

## Errors

| Case | Behaviour |
|---|---|
| Handler fails | Retried, then dead-lettered with an alert (OFF-01-02); deletion must not be silently skipped |

## Acceptance criteria

- [ ] After `household.deleted`, a search of every offers table for the household id finds nothing except
      `deletion_log` (which holds only the event id).
- [ ] Its private offers no longer appear anywhere and planner got `offers.expired` for them.
- [ ] A user-added shop the household created and admin approved is still listed for others.
