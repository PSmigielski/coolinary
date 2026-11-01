# NTF-03-01 Push templates and categories

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-03 Push delivery](README.md) |
| Release | R1 (R2 templates marked) |
| Size | M |
| Depends on | NTF-02-01 |
| Unblocks | NTF-03-02, KIT-09 (expiry reminders), PLN-17 (offer-ends reminders), KIT-12 (savings), IDN-07/IDN-08 (household pushes) |

## Goal

One catalogue of push messages, each tied to a category the user can switch off, in three languages, short enough
to read on a lock screen.

## Business rules

1. Categories (proto enum `Category`): `expiry`, `offer_ends`, `household`, `savings`, `near_shop`. Every push
   template belongs to exactly one category.
2. Catalogue:

| Template id | Category | Expected sender | Release | Variables | en example |
|---|---|---|---|---|---|
| `expiry_tomorrow` | expiry | kitchen (KIT-09) | R1 | `count` int 1–999, `names` list (≤ 3 shown) | "3 items expire tomorrow" / "Milk, spinach, yoghurt" |
| `expiry_today` | expiry | kitchen (KIT-09) | R1 | `count`, `names` | "2 items expire today" / "Chicken thighs, cream" |
| `offer_ends_tomorrow` | offer_ends | planner (PLN-17) | R2 | `count` int, `shop` string ≤ 40 | "2 offers on your list end tomorrow" / "At Lidl" |
| `household_member_joined` | household | identity (IDN-07) | R2 | `member_name` string ≤ 40 | "Ola joined your household" |
| `household_member_left` | household | identity (IDN-08) | R2 | `member_name` | "Adam left your household" |
| `household_removed_you` | household | identity (IDN-08) | R2 | `household_name` string ≤ 40 | "You were removed from Ewa's home" / "You now have your own empty fridge." |
| `savings_monthly` | savings | kitchen (KIT-12) | R2 | `month` date, `amount` money, `kg` string ≤ 6 | "October: €38 and 4 kg saved" |

3. `near_shop` has no server template: near-shop reminders are created on the phone (D13). The category exists so
   the user's choice is stored and synced (NTF-04-02).
4. Length limits after rendering: title ≤ 50 characters, body ≤ 150 characters. When `names` makes the body too
   long, fewer names are shown and "+N" is appended ("Milk, spinach +3"). Other overflows are cut with "…".
5. Each message carries data: `path` (deep link in the app, starts with `/`, ≤ 200 characters, given by the sender),
   `category`, `template_id`. Total data ≤ 1 KB.
6. Android channel id = the category name (`expiry`, …); the app creates these channels (MOB-06). Sound: default
   for `expiry`, `offer_ends`, `household`; none for `savings`.
7. The sender in the table is the expected one; each sender's task confirms it. A template can be added only with
   all three languages (NTF-02-01 rule 8).

## Contract

`contracts/coolinary/notifier/v1/notifier.proto`

```proto
enum Category {
  CATEGORY_UNSPECIFIED = 0;
  CATEGORY_EXPIRY = 1;
  CATEGORY_OFFER_ENDS = 2;
  CATEGORY_HOUSEHOLD = 3;
  CATEGORY_SAVINGS = 4;
  CATEGORY_NEAR_SHOP = 5;
}
```

## Data

No schema change.

## Errors

As NTF-02-01.

## Acceptance criteria

- [ ] Snapshot per template per language; Polish plurals right for `count` 1, 2, 5, 22.
- [ ] `expiry_tomorrow` with 12 long names stays ≤ 150 characters and ends with "+N".
- [ ] A test fails if a template's category is missing or it lacks a language.
