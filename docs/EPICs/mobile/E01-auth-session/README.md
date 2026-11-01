# MOB-01 Auth session

Goal: the app always has a valid identity for its calls: an anonymous account from the first real action, sign-in
with an e-mail code that keeps the anonymous user's data, a 15-minute JWT that refreshes itself, and a clean
sign-out. Decisions 11 and 12 in [docs/EPICs/README.md](../../README.md).

## Flow

```
first action ──> ensureSession() ──> POST /api/auth/sign-in/anonymous ──> GET /api/auth/token
                                                     │                          │
                                         session cookie in SecureStore    JWT in memory (15 min,
                                                                          refreshed 120 s early)
e-mail ──> send code ──> verify code ──> identity links anonymous → account ──> forceRefresh() JWT
                                                                               └─> householdChanged / userChanged
401 ──> forceRefresh + retry once ──> token 401 ──> session lost ──> wipe (+ new anonymous account)
sign-out ──> flush outbox ──> unregister device ──> POST /api/auth/sign-out ──> wipe
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [MOB-01-01](MOB-01-01-auth-client-setup.md) | Better Auth client setup | R1 | S | GW-02, IDN-02, IDN-03 |
| [MOB-01-02](MOB-01-02-auth-state.md) | Auth state and app start | R1 | M | MOB-01-01 |
| [MOB-01-03](MOB-01-03-anonymous-start.md) | Anonymous start | R1 | S | MOB-01-02, IDN-03 |
| [MOB-01-04](MOB-01-04-jwt-provider.md) | JWT provider and refresh | R1 | M | MOB-01-01, MOB-01-02 |
| [MOB-01-05](MOB-01-05-email-code-sign-in.md) | Sign-in with an e-mail code and anonymous upgrade | R1 | M | MOB-01-02, MOB-01-04, MOB-04-02, IDN-02, IDN-03 |
| [MOB-01-06](MOB-01-06-session-loss.md) | 401 recovery and session loss | R1 | S | MOB-01-04, MOB-02-01 |
| [MOB-01-07](MOB-01-07-sign-out.md) | Sign-out and local wipe | R1 | M | MOB-01-02, MOB-03-01, MOB-04-01, MOB-06-03 |
| [MOB-01-08](MOB-01-08-account-deletion-export.md) | Account deletion and data export calls | R1 | S | MOB-01-07, MOB-02-01, IDN-09 |
