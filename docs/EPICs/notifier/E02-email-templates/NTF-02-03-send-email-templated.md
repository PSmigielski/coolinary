# NTF-02-03 SendEmail with templates, languages and idempotency

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-02 E-mail templates and languages](README.md) |
| Release | R1 |
| Size | M |
| Depends on | NTF-02-01, NTF-02-02 |
| Unblocks | NTF-02-04, NTF-02-05, IDN-02 (sign-in code e-mail) |

## Goal

Callers send e-mail by template id, language and variables. Retries by the caller never send the same e-mail twice,
and a permanently rejected address is reported as such.

## Business rules

1. `template_id` is required and must be an e-mail template; requests without it → `INVALID_ARGUMENT`
   "template_id is required". The old raw fields (`subject`, `text_body`, `html_body`) are kept in the proto as
   deprecated so `make proto-breaking` passes, and are ignored.
2. `to` must be an e-mail address (existing pattern) of at most 254 characters.
3. `language` is resolved per NTF-02-01 rule 2.
4. Idempotency: when `idempotency_key` (at most 128 characters) is set, notifier runs
   `SET ntf:email:idem:<key> <delivery_id> NX EX 600` in `redis-state`. If the key exists, nothing is sent and the
   stored result is returned with `deduplicated = true`.
5. Per-address safety limit, on top of identity's own limits: at most 5 e-mails per address per 15 minutes and 20
   per day, counted by `recipient_hash` in `redis-state`. Above it → `RESOURCE_EXHAUSTED` and no send. Founder alert
   e-mails (NTF-05) are exempt. Decided (Q030): this net stays, on top of identity's limits.
6. SMTP send with a 10 s timeout. On a transient failure (connection error, timeout, SMTP 4xx) notifier retries once
   after 1 s, but only if the caller's deadline leaves at least 2 s.
7. SMTP 5xx (e.g. mailbox does not exist) → `FAILED_PRECONDITION` "recipient rejected", no retry, so identity can
   tell the user to check the address.
8. Final transient failure → `UNAVAILABLE`.
9. Every attempt writes or updates one delivery row (NTF-02-02) with `user_id` when given.
10. Headers: `From` = `MAIL_FROM`; `Content-Language` = resolved language; `Auto-Submitted: auto-generated`. No
    `Reply-To`.

## Contract

`contracts/coolinary/notifier/v1/notifier.proto`

```proto
message SendEmailRequest {
  string to = 1;
  string subject = 2 [deprecated = true];
  string text_body = 3 [deprecated = true];
  string html_body = 4 [deprecated = true];
  string template_id = 5;                 // e.g. "sign_in_code"
  string language = 6;                    // BCP 47
  map<string, string> variables = 7;      // typed by the template, sent as strings
  string idempotency_key = 8;             // optional, <= 128 chars
  string user_id = 9;                     // optional, for the delivery log and erasure
}

message SendEmailResponse {
  string message_id = 1;    // SMTP message id
  string delivery_id = 2;
  bool deduplicated = 3;
}
```

Variables travel as strings; `money` is `"<minor>:<CUR>"` (e.g. `"3800:EUR"`), `list` is a JSON array string.

## Data

No schema change beyond NTF-02-02.

## Errors

| Case | gRPC status |
|---|---|
| Missing/unknown template, bad variable, bad address | `INVALID_ARGUMENT` |
| Over the per-address limit | `RESOURCE_EXHAUSTED` |
| Address rejected by the mail server | `FAILED_PRECONDITION` |
| SMTP down after retry | `UNAVAILABLE` |

## Acceptance criteria

- [ ] Two calls with the same `idempotency_key` within 10 minutes deliver one e-mail to Mailpit; the second returns
      `deduplicated = true` with the same `delivery_id`.
- [ ] The 6th e-mail to one address within 15 minutes is refused.
- [ ] With Mailpit stopped, the call returns `UNAVAILABLE` after one retry and the delivery row is `failed` with
      `smtp_unavailable`.
- [ ] A request with only the deprecated raw fields is refused.
