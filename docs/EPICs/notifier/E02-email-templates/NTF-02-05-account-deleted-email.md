# NTF-02-05 Account deletion confirmation e-mail

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-02 E-mail templates and languages](README.md) |
| Release | R1 |
| Size | S |
| Depends on | NTF-02-01, NTF-02-03 |
| Unblocks | IDN-09 (if it chooses to confirm deletions) |

## Goal

After an account is deleted, its former owner gets a short confirmation, so a deletion they did not make is noticed.

## Business rules

1. Template id `account_deleted`, channel `email`, languages pl / de / en, no variables.
2. identity sends it right before it erases the address (IDN-09), with `idempotency_key = "account_deleted:<user_id>"`.
3. Text (en): subject "Your Coolinary account has been deleted"; body "Your account and its data have been deleted.
   If you shared a household, its fridge, list and plans stay with the other members. If you didn't ask for this,
   contact us at <support address>." The support address comes from config `SUPPORT_EMAIL`.
4. No link to the app or the store.
5. `user_id` is not passed (the account is gone); the delivery row has only the hash.

## Contract

Uses `SendEmail` (NTF-02-03).

## Data

New config `SUPPORT_EMAIL` in `notifier/.env.*`.

## Errors

As NTF-02-03.

## Acceptance criteria

- [ ] Snapshot per language.
- [ ] Sending twice with the same idempotency key delivers one e-mail.

## Notes

Decided (Q039): identity sends this e-mail after every deletion of an account with an e-mail address.
