# PLN-21 Household lifecycle: event consumers and data export

Releases R1–R2. Planner reacts to household, account, catalogue and offer events (Redis Streams, consumer group
`planner`, idempotent by event id, PLN-01) and gives the gateway its part of a GDPR export (D9).

Recipe-side handling of `catalogue.ingredient_merged` and `catalogue.nutrition_changed` (dish ingredients and
calories) belongs to PLN-02 / PLN-04; this epic covers lists, plans, portions, base pantries, extras and offers.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-21-01](PLN-21-01-household-and-account-deletion.md) | Household and account deletion | R1 | M | PLN-01, PLN-15-01 |
| [PLN-21-02](PLN-21-02-member-added-and-removed.md) | Member added and removed | R2 | S | PLN-01, PLN-12-06, PLN-12-01 |
| [PLN-21-03](PLN-21-03-plan-downgrade.md) | Plus → Free downgrade | R2 | S | PLN-01, PLN-12-06 |
| [PLN-21-04](PLN-21-04-ingredient-merge-for-lists.md) | Ingredient merge for lists, pantries and extras | R1 | S | PLN-01, PLN-15-01 |
| [PLN-21-05](PLN-21-05-offers-events.md) | Offers events | R2 | S | PLN-01, PLN-17-01 |
| [PLN-21-06](PLN-21-06-export-household-data.md) | `ExportHouseholdData` | R1 | M | PLN-15-01, PLN-08 |
