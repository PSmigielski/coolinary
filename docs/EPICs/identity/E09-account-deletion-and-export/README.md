# IDN-09 Account deletion and data export

GDPR: a person can delete their account at any time and get a copy of their data. Identity deletes what it holds
and tells every other service through events (`account.deleted`, `household.deleted`; derived decision D9). The
export is assembled by the gateway from each service's export RPC; identity provides `ExportUserData`.

## Flow (deletion)

1. App → gateway → `RequestAccountDeletion`: a 6-digit code is e-mailed (accounts with an e-mail only).
2. App → gateway → `DeleteAccount(code)`: identity handles the household (refused with `owner_must_transfer` while the owner hasn't named a new owner; otherwise removal of the member or deletion of a one-person household), deletes the user,
   emits events; the app signs out.
3. kitchen, planner, offers and notifier erase their data on the events.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-09-01 | Deletion confirmation code | R1 | S | IDN-02-01, NTF-02 (template `account_deletion_code`) |
| IDN-09-02 | Delete an account | R1 | L | IDN-09-01, IDN-01-01, IDN-04-01 |
| IDN-09-03 | `ExportUserData` | R1 | S | IDN-04-04 |
