# IDN-06 Profiles without accounts

Children, and anyone without a phone, can be added to a household as profiles: a name only, used for portions and
goals (features.md "Linking people into a household"). They have a member id like everyone else, so kitchen keeps
their diet and calorie target and planner sizes their portion. They don't count toward the 5 accounts.

Identity stores no age, sex, minor status or goal for a profile: those are health data and stay on the phone; the
daily kcal target and its kind (adult/child) live in kitchen (KIT-11).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-06-01 | Add a profile | R2 | S | IDN-04-04, IDN-04-05, IDN-05-01, IDN-01-01 |
| IDN-06-02 | Delete a profile | R2 | S | IDN-06-01 |

Renaming a profile is part of `UpdateMember` (IDN-04-05).
