# MOB-01-05 Sign-in with an e-mail code and anonymous upgrade

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-01 Auth session](README.md) |
| Release | R1 |
| Size | M |
| Depends on | MOB-01-02, MOB-01-04, MOB-04-02, IDN-02 (e-mail OTP), IDN-03 (link anonymous → account), GW-02 |
| Unblocks | later sign-in screens |

## Goal

A user signs in with a 6-digit code sent by e-mail, the only sign-in method (README decision 12). An anonymous user
who signs in keeps their household, fridge and list.

## Business rules

1. Two calls, exposed as mutations `useRequestCode()` and `useVerifyCode()`:
   1. `authClient.emailOtp.sendVerificationOtp({ email, type: "sign-in" })`.
   2. `authClient.signIn.emailOtp({ email, otp })`.
2. The e-mail is trimmed and lowercased before sending; it must match a simple check (`x@y.z`, ≤ 254 chars) before
   any call. It is kept in memory for the flow only, never logged or persisted.
3. The code is exactly 6 digits; anything else is rejected on the phone without a call.
4. Server limits (IDN-02): a code is valid **10 minutes** (decided, Q031) and allows 3 attempts. The app shows a
   countdown from 600 s; after expiry or `TOO_MANY_ATTEMPTS` the user must request a new code.
5. Resend cooldown on the phone: 60 s between `sendVerificationOtp` calls for the same e-mail (decided, Q032), the
   same as identity's limit of one code per minute (IDN-02-03).
6. Sign-in needs a connection. Before calling `signIn.emailOtp`, the outbox (MOB-04) is flushed for the current
   household; if it cannot be flushed within 10 s, the user is told that unsent changes will be sent first and the
   flow waits.
7. **Upgrade of an anonymous account.** When `status === "anonymous"`, the same calls link the account: Better Auth
   creates the real user and runs identity's `onLinkAccount`, which moves the household to the new user (IDN-03).
   The user id changes. After success:
   1. `forceRefresh()` the JWT (MOB-01-04); `userChanged` fires; `householdChanged` fires only if identity put the
      user in a different household (the e-mail already had an account with its own household). In that case
      (decided, Q036) identity moved the anonymous fridge items into that household and dropped the anonymous plans,
      shopping list and profile; the app refetches stock and shows once: "Your fridge items were added to your
      account. Plans and the shopping list from this phone were not kept."
   2. `GET /v1/me` is refetched; `status = "signedIn"`.
   3. Push registration is redone for the new user (MOB-06-03).
8. Signing in when `status === "signedIn"` with a different e-mail first runs sign-out (MOB-01-07).
9. Sign-in on a fresh install with an e-mail that has an account simply signs in; no anonymous account is created
   first (rule 1 of MOB-01-03 does not fire because sign-in itself creates the session).

## Contract

| Call | HTTP | Success | Failure |
|---|---|---|---|
| Request code | `POST /api/auth/email-otp/send-verification-otp` `{ email, type: "sign-in" }` | `200` | `400` invalid e-mail, `429` too many requests |
| Verify code | `POST /api/auth/sign-in/email-otp` `{ email, otp }` | `200` + session cookie | `400`/`401` invalid or expired code, `TOO_MANY_ATTEMPTS` |

```ts
type SignInError = "invalid_email" | "invalid_code" | "code_expired" | "too_many_attempts" | "rate_limited" | "offline" | "server_down";
```

## Data

No new storage. The session cookie is replaced by `expoClient` in SecureStore.

## Errors

| Case | Behaviour |
|---|---|
| Wrong code (attempt 1–2) | `invalid_code`, attempts left = 3 − attempts made |
| `TOO_MANY_ATTEMPTS` | `too_many_attempts`; code field disabled until a new code is requested |
| Code older than 10 minutes (local countdown or server error) | `code_expired` |
| 429 on request | `rate_limited`, wait `Retry-After` or 60 s |
| Offline | `offline`; the request is not queued |

## Acceptance criteria

- [ ] Anonymous user with 5 fridge items signs in with a new e-mail: after sign-in `GET /v1/stock` still returns
      the 5 items and `GET /v1/me` shows the e-mail.
- [ ] Anonymous user signs in with an e-mail that already owns a household: the app switches to that household,
      the old household's cache is removed, and no queued change is lost (it was flushed first).
- [ ] Three wrong codes give `too_many_attempts`; a new code then works.
- [ ] No e-mail address or code appears in logs, AsyncStorage or the query cache.
