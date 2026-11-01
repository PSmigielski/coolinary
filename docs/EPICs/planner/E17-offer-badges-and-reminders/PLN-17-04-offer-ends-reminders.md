# PLN-17-04 "Offer ends tomorrow" reminders

| | |
|---|---|
| Service | planner |
| Epic | [PLN-17 Offer badges and offer-ends reminders](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | PLN-17-01, NTF-03 (SendPush) |
| Unblocks | — |

## Goal

Remind the household when an offer on their list ends tomorrow (features.md "Flyers and offers": "Offer ends tomorrow
reminders"). Planner runs it because it knows the list (D1, D12).

## Business rules

1. An hourly job selects households whose local time is between 18:00 and 18:59.
2. For each, it finds open lines whose primary or pinned badge has `valid_to` = tomorrow (local), excluding upcoming
   badges and offers already reminded.
3. One push per household per day, sent with NTF-03 `SendPush` to the household (all account members), category
   `offer_ends`, template `offer_ends_tomorrow`, data: number of offers and the shop name when they are all at one
   shop ("2 offers on your list end tomorrow at Lidl"). No item names.
4. notifier applies each person's preferences and quiet hours; planner doesn't check them.
5. Each reminded (household, offer) pair is stored; it is never reminded again.
6. Free and Plus both get the reminder.
7. notifier down: the job retries that household in the next run until 21:00 local, then gives up for the day.
8. The push opens the shopping list (deep link `coolinary://list`).

## Contract

Uses `SendPush` (NTF-03). No new planner RPC.

## Data

`planner/migrations/000NN_offer_reminders.sql`

```sql
CREATE TABLE offer_reminders_sent (
    household_id uuid NOT NULL,
    offer_id     text NOT NULL,
    sent_on      date NOT NULL,
    PRIMARY KEY (household_id, offer_id)
);
```

## Errors

| Case | Behaviour |
|---|---|
| notifier down | the job retries that household in the next run until 21:00 local, then gives up for the day (rule 7) |
| offers down | the run skips households whose offer data it can't confirm; no reminder is better than a wrong one |

## Acceptance criteria

- [ ] A list with two Lidl offers ending tomorrow gets one push at 18:xx local time with count 2.
- [ ] The next day nothing is sent for the same offers.
- [ ] A ticked line is not reminded.
- [ ] Logs show household id hashes and counts only.

## Notes

- Decided (Q179, Q178): reminders are sent at 18:00 local time, to Free and Plus households.
