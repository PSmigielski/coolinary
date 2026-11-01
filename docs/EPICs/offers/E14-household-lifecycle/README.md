# OFF-14 Household lifecycle: event consumers and data export

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | offers forgets households and accounts that are deleted, keeps its ingredient references valid when kitchen merges ingredients, and hands over a household's data on request (GDPR; README D9). |

Plan changes are read from identity (OFF-01-03, Decided Q004) and applied by OFF-04-04 (my-shops rules). `household.member_added` and
`household.member_removed` need no action in offers: my shops, preferences and private offers belong to the household,
not to a member.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-14-01](OFF-14-01-erase-on-deletion.md) | Erase data on household or account deletion | R2 | M | OFF-01-02, IDN-09 (account.deleted), IDN-08 (household.deleted) |
| [OFF-14-02](OFF-14-02-ingredient-merged.md) | Follow merged catalogue ingredients | R2 | S | OFF-01-02, KIT-02 (catalogue.ingredient_merged) |
| [OFF-14-03](OFF-14-03-export-household-data.md) | Export household data | R2 | S | OFF-04-01, OFF-04-02, OFF-10-04 |
