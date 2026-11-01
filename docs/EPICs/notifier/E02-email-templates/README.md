# NTF-02 E-mail templates and languages

**Goal.** Callers send e-mail by template id, language and variables; notifier owns every text in Polish, German
and English. The first and most important e-mail is the sign-in code (decision 12: e-mail code is the only way to
sign in). The delivery log stops storing addresses and message text.

**Flow.** identity → `SendEmail(template_id, language, variables, idempotency_key)` → template renders subject, text
and HTML → SMTP (Mailpit in dev) → delivery row with a hashed recipient.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [NTF-02-01](NTF-02-01-template-registry.md) | Template registry and renderer | R1 | M | — |
| [NTF-02-02](NTF-02-02-delivery-log-privacy.md) | Delivery log without personal data | R1 | M | — |
| [NTF-02-03](NTF-02-03-send-email-templated.md) | SendEmail with templates, languages and idempotency | R1 | M | NTF-02-01, NTF-02-02 |
| [NTF-02-04](NTF-02-04-sign-in-code-email.md) | Sign-in code e-mail | R1 | S | NTF-02-01, NTF-02-03 |
| [NTF-02-05](NTF-02-05-account-deleted-email.md) | Account deletion confirmation e-mail | R1 | S | NTF-02-01, NTF-02-03 |
