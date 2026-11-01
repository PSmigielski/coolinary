# KIT-09-01 Which items a reminder lists

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-09 Expiry reminders](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-07-05 (expiry status), KIT-01-06 (plan from identity) |
| Unblocks | KIT-09-02 |

## Goal

Decide, for one household and one local day, whether a reminder goes out and what it says, so users get at most one
useful nudge a day and never a reminder they cannot act on.

## Business rules

1. Only households whose plan, asked from identity for this run (KIT-01-06, Q004), is `plus` get reminders. Free and anonymous households never do
   (features.md §6). A plan change takes effect at the next selection.
2. "Today" and "tomorrow" are local dates in the household's timezone (KIT-10-01).
3. **One reminder a day at most**, of one of two kinds:
   1. `today`: at least one item has `expires_on = today` → the reminder lists the items expiring today;
   2. otherwise `tomorrow`: at least one item has `expires_on = tomorrow` → lists those;
   3. otherwise no reminder.
   An item due tomorrow is therefore usually named twice: the day before and on the day.
4. Items already expired are never listed (the fridge list shows them in red).
5. Freezer items are listed only when their own date is today or tomorrow (rare by nature).
6. Unmatched items are listed like the others.
7. **Content:** `count` = number of items of the chosen kind; `names` = up to 3 display names, ordered by storage
   place (fridge, pantry, freezer), then name; notifier adds "+N" for the rest.
8. A household with more than 999 due items sends `count = 999` (template limit, NTF-03-01).
9. Selection reads stock only; it writes nothing.

## Contract

Internal function `selectReminder(household, localDate) → (kind, count, names) | none`. No RPC.

## Data

No schema change.

## Errors

None (pure selection).

## Acceptance criteria

- [ ] Plus household, items: milk today, spinach tomorrow, yoghurt tomorrow → kind `today`, count 1, names [milk].
- [ ] Same household next day (milk removed) → kind `today`, count 2, names [spinach, yoghurt].
- [ ] Items only expired yesterday → no reminder.
- [ ] Free household with items expiring today → no reminder.
- [ ] Household in `Europe/Berlin` at 23:30 UTC on 31 March (01:30 local on 1 April, after the DST change) uses
      1 April as today.
