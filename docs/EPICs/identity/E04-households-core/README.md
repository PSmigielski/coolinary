# IDN-04 Households and members core

A household is one shared fridge, shopping list and meal plan (features.md "Linking people into a household").
Today every new account gets a personal household (`createPersonalHousehold`) and `GetHousehold` returns it. This
epic adds the rules every later epic relies on: one household per account, exactly one owner, events on creation,
the JWT claims, membership checks and renaming.

Roles: `owner` (one per household, has an account), `member` (has an account), `profile` (no account: children,
people without a phone; IDN-06).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-04-01 | One household per account, one owner per household | R1 | S | — |
| IDN-04-02 | Household creation events and default names | R1 | S | IDN-01-01, IDN-04-01 |
| IDN-04-03 | JWT claims | R1 | S | IDN-04-01 |
| IDN-04-04 | Membership checks and richer `GetHousehold` | R1 | S | IDN-04-01 |
| IDN-04-05 | Rename household and display names | R1 | S | IDN-04-04 |
