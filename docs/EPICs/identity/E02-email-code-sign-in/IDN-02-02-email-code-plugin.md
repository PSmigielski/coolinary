# IDN-02-02 E-mail code plugin and code e-mail

| | |
|---|---|
| Service | identity |
| Epic | [IDN-02 Sign-in with e-mail code](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | IDN-02-01, NTF-02 (SendEmail with template `sign_in_code`) |
| Unblocks | IDN-02-03, IDN-02-04, IDN-02-05, IDN-03-02, IDN-03-03, GW-02, MOB-01 |

## Goal

A person signs in, or creates an account, by typing a 6-digit code that arrives by e-mail in their language.

## Business rules

1. Add Better Auth's `emailOTP` plugin (`better-auth/plugins`) in `src/auth.ts`, then `bun run auth:schema` and
   `bun run db:generate` if the plugin changes the schema.
2. Code: 6 digits (`otpLength: 6`), valid for 10 minutes (`expiresIn: 600`), 3 wrong tries allowed
   (`allowedAttempts: 3`); after the third wrong try the code is dead and Better Auth returns `TOO_MANY_ATTEMPTS`.
   Decided (Q031): codes are valid 10 minutes; the app's countdown uses 600 s (MOB-01-05).
3. Codes are stored hashed (`storeOTP: "hashed"`), never in plain text.
4. A new code replaces the previous one for the same e-mail; only the latest code works.
5. `disableSignUp: false`: a correct code for an e-mail with no account creates the account. The personal
   household is created by the user-create hook (IDN-04-02), except when an anonymous session is being upgraded
   (IDN-03-02, IDN-03-03).
6. `sendVerificationOTP` sends only for `type = "sign-in"`. For `"email-verification"` and `"forget-password"` it
   sends nothing and returns an error (there are no passwords, and signing in with a code verifies the e-mail).
7. E-mail addresses are trimmed and lower-cased before anything else. More than 254 characters, or not a valid
   address, is rejected with 400 before any code is created.
8. The e-mail is sent with notifier `SendEmail`: `template = "sign_in_code"`, `variables = {code, valid_minutes: 10}`,
   `language` = the `x-language` request header if the gateway set it, else the first of `pl`, `de`, `en` found in
   `Accept-Language`, else `en`.
9. The answer to "send code" is the same whether the e-mail has an account or not (no account enumeration).
10. If notifier is unavailable, "send code" returns 503 `code_not_sent`; the stored code is deleted so the next
    request starts clean.
11. A successful sign-in sets `email_verified = true` on the user.
12. Neither the code nor the e-mail address is ever logged or put in a metric label.
13. Metric `auth_codes_sent_total{language}`; `auth_sign_ins_total{result}` with `result` in `ok`, `wrong_code`,
    `expired`, `too_many_attempts`.

## Contract

HTTP (Better Auth, reached through the gateway's `/api/auth/*` proxy):

| Method and path | Body | Answer |
|---|---|---|
| `POST /api/auth/email-otp/send-verification-otp` | `{ "email": "ewa@coolinary.test", "type": "sign-in" }` | `200 {"success": true}` |
| `POST /api/auth/sign-in/email-otp` | `{ "email": "...", "otp": "123456", "name": "Ewa" }` (`name` optional, used only for a new account) | `200` session (Expo plugin handles the token) |

## Data

Uses Better Auth's `verification` table (exists). No other change.

## Errors

| Case | HTTP | Body `code` |
|---|---|---|
| Invalid e-mail | 400 | `invalid_email` |
| Wrong code | 400 | `INVALID_OTP` (Better Auth) |
| Code older than 10 min | 400 | `OTP_EXPIRED` |
| Third wrong try | 403 | `TOO_MANY_ATTEMPTS` |
| Notifier down | 503 | `code_not_sent` |
| Type other than `sign-in` | 400 | `unsupported_type` |

## Acceptance criteria

- [ ] With `Accept-Language: pl-PL`, the code e-mail in Mailpit uses the Polish template.
- [ ] A first sign-in with a new e-mail creates one user, one household and one owner member.
- [ ] A second sign-in with the same e-mail creates nothing new and returns the same household in the JWT.
- [ ] Three wrong codes kill the code; the right code afterwards fails until a new one is sent.
- [ ] A code used after 10 minutes fails with `OTP_EXPIRED` (test clock).
- [ ] The `verification` table holds no plain 6-digit value.
- [ ] Log lines of a full sign-in contain no e-mail address and no code.
