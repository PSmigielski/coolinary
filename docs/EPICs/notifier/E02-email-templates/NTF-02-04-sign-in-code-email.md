# NTF-02-04 Sign-in code e-mail

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-02 E-mail templates and languages](README.md) |
| Release | R1 |
| Size | S |
| Depends on | NTF-02-01, NTF-02-03 |
| Unblocks | IDN-02 (e-mail code sign-in) |

## Goal

The only way to sign in is a code sent by e-mail (decision 12), so this e-mail must arrive fast, read clearly on a
phone and never look like phishing.

## Business rules

1. Template id `sign_in_code`, channel `email`, languages pl / de / en.
2. Variables: `code` (string, 4–8 digits, shown in groups of 3 when it has 6 digits: `482 913`), `minutes` (int,
   1–60: how long the code is valid, sent by identity).
3. The e-mail contains no links and no buttons, only the code.
4. Decided (Q037): the subject contains the code, so it can be read from the phone's notification.
5. Decided (Q038): German uses "du".
6. Texts (code 482913, 10 minutes, Q031):

| | Subject | Body (text version) |
|---|---|---|
| en | `482 913 is your Coolinary code` | Enter this code in the app to sign in: 482 913. It is valid for 10 minutes. If you didn't ask for it, ignore this e-mail; nobody can sign in without the code. |
| pl | `482 913 to Twój kod do Coolinary` | Wpisz ten kod w aplikacji, aby się zalogować: 482 913. Kod jest ważny przez 10 minut. Jeśli to nie Ty, zignoruj tę wiadomość; bez kodu nikt się nie zaloguje. |
| de | `482 913 ist dein Coolinary-Code` | Gib diesen Code in der App ein, um dich anzumelden: 482 913. Er ist 10 Minuten gültig. Wenn du ihn nicht angefordert hast, ignoriere diese E-Mail; ohne den Code kann sich niemand anmelden. |

7. HTML shows the code at 32 px, monospace, with letter spacing, and the same sentences.
8. The footer line: "You get this e-mail because someone entered this address in the Coolinary app." (per
   language).

## Contract

Uses `SendEmail` (NTF-02-03) with `template_id = "sign_in_code"`.

## Data

No schema change.

## Errors

As NTF-02-01 (bad `code` or `minutes` → `INVALID_ARGUMENT`).

## Acceptance criteria

- [ ] Snapshot per language; the Polish minutes form is right for 1, 2, 5 ("minutę", "minuty", "minut").
- [ ] The HTML contains no `<a` element.
- [ ] Through Mailpit the e-mail shows the grouped code in subject and body.
