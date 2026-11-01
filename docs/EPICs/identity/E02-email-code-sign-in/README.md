# IDN-02 Sign-in with e-mail code; remove passwords

Decision 12: the only sign-in method is a 6-digit code sent by e-mail (Better Auth `emailOTP` plugin). Signing in
with an unknown e-mail creates the account, so sign-up and sign-in are one flow. Password sign-up is removed.

## Flow

1. App → gateway `/api/auth/email-otp/send-verification-otp` `{email, type: "sign-in"}` → identity checks limits,
   stores a hashed code and asks notifier to send the e-mail in the user's language.
2. App → `/api/auth/sign-in/email-otp` `{email, otp}` → Better Auth verifies, creates the account if new (with its
   personal household, IDN-04-02) or links an anonymous session (IDN-03), and returns a session.
3. App → `/api/auth/token` → 15-minute JWT the gateway accepts.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-02-01 | Notifier gRPC client | R1 | S | — |
| IDN-02-02 | E-mail code plugin and code e-mail | R1 | M | IDN-02-01, NTF-02 (SendEmail with template) |
| IDN-02-03 | Sign-in rate limits and throttles | R1 | M | IDN-02-02 |
| IDN-02-04 | Remove password sign-up | R1 | S | IDN-02-02 |
| IDN-02-05 | Session lifetime and JWT issuance | R1 | S | IDN-02-02 |
