# NTF-02-01 Template registry and renderer

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-02 E-mail templates and languages](README.md) |
| Release | R1 |
| Size | M |
| Depends on | — |
| Unblocks | NTF-02-03, NTF-02-04, NTF-02-05, NTF-03-01, NTF-05-02 |

## Goal

Every e-mail and push text is written once per language in notifier, and callers only send a template id and
variables. Callers never build message text.

## Business rules

1. A template has an id (`snake_case`, e.g. `sign_in_code`), a channel (`email` or `push`), a list of variables and
   one render function per language. Templates are TypeScript modules in `notifier/src/templates/<id>.ts`, reviewed
   like code. No templates in the database.
2. Supported languages: `pl`, `de`, `en`. Language resolution: take the primary subtag of the requested BCP 47 tag
   (`de-AT` → `de`, `pl-PL` → `pl`); anything else, or empty, → `en`.
3. Variable types: `string` (with a max length), `int` (with min and max), `money` (amount in minor units + ISO 4217
   currency, formatted with `Intl.NumberFormat` for the language, e.g. `38,00 €`, `38,00 zł`), `list` of strings
   (with max items shown), `date` (ISO `YYYY-MM-DD`, formatted per language).
4. Validation before rendering: a missing required variable, an unknown variable, a wrong type or an out-of-range
   value → `INVALID_ARGUMENT` naming the variable, never its value.
5. Plurals use `Intl.PluralRules` for the language (Polish needs one / few / many: "1 produkt", "2 produkty",
   "5 produktów"). A helper `plural(lang, n, forms)` is provided.
6. E-mail output: `subject`, `text`, `html`. Push output: `title`, `body`. Variables are HTML-escaped in `html`.
7. E-mail HTML uses one shared layout: max width 560 px, inline CSS only, system font stack, the word "Coolinary"
   as text (no remote images, no tracking pixels, no click tracking, no links unless the template needs one),
   footer with one line saying why the e-mail was sent.
8. Every template exists in all three languages; a unit test fails if one is missing or if the languages declare
   different variables.
9. Snapshot tests render each template in each language with fixed sample variables.
10. Rendered texts are never logged.

## Contract

No gRPC change here; NTF-02-03 and NTF-03-02 expose templates.

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Unknown template id | `INVALID_ARGUMENT` "unknown template" |
| Variable problem | `INVALID_ARGUMENT` "variable <name>: <reason>" |

## Acceptance criteria

- [ ] `de-AT` renders German, `fr` renders English.
- [ ] Polish plural forms are correct for 1, 2, 5, 12, 22, 25.
- [ ] `money` 3800 EUR renders `38,00 €` in de and pl, `€38.00` in en.
- [ ] A test fails when a template lacks a language.
- [ ] HTML output contains no `<img` with a remote `src` and no `http` link unless declared by the template.

## Notes

`GATE:` Polish and German texts get a native-speaker check before release 1 (business model §6, one-off per
language).
