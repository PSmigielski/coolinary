# MOB-01-01 Better Auth client setup

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-01 Auth session](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | GW-02 (`/api/auth/*` proxy), IDN-02 (e-mail OTP plugin), IDN-03 (anonymous plugin) |
| Unblocks | MOB-01-02, MOB-01-03, MOB-01-04, MOB-01-05 |

## Goal

One configured Better Auth client in the app, talking to identity through the gateway, so every other auth task
calls typed methods instead of hand-written requests.

## Business rules

1. Install with `bunx expo install better-auth @better-auth/expo expo-network` (expo-secure-store is already there).
   `expo-network` is required by the Expo client and is also used by MOB-02-04.
2. One module `src/lib/auth/client.ts` exports `authClient`, built with `createAuthClient` from `better-auth/react`:
   - `baseURL` = `API_URL` from `src/api/client.ts` (the gateway; it proxies `/api/auth/*` to identity), `basePath`
     `/api/auth`.
   - Plugins: `expoClient({ scheme: "coolinary", storagePrefix: "coolinary", storage: SecureStore })`,
     `anonymousClient()`, `emailOTPClient()`, `jwtClient()`.
3. `scheme` is the `scheme` in `app.json` (`coolinary`). identity must list `coolinary://` in `trustedOrigins`
   (IDN-02); dev builds and Expo Go also need the `exp://` origin in dev only, never in prod.
4. The session cookie lives only in SecureStore (written by `expoClient`). No auth value is ever written to
   AsyncStorage, the query cache, Zustand persist or logs.
5. Do not change Metro's `unstable_enablePackageExports`; Better Auth needs package exports, which SDK 58 enables by
   default.
6. Every auth request uses `credentials: "omit"` (the Expo client sets the cookie header itself; `include` interferes).
7. No other module imports `better-auth` directly; they use `authClient` or the hooks from MOB-01-02.

## Contract

Endpoints used through the client (all under the gateway's `/api/auth/*` proxy, GW-02):

| Client method | HTTP | Used by |
|---|---|---|
| `authClient.signIn.anonymous()` | `POST /api/auth/sign-in/anonymous` | MOB-01-03 |
| `authClient.emailOtp.sendVerificationOtp({ email, type: "sign-in" })` | `POST /api/auth/email-otp/send-verification-otp` | MOB-01-05 |
| `authClient.signIn.emailOtp({ email, otp })` | `POST /api/auth/sign-in/email-otp` | MOB-01-05 |
| `authClient.token()` | `GET /api/auth/token` → `{ token }` | MOB-01-04 |
| `authClient.getSession()` | `GET /api/auth/get-session` | MOB-01-02 |
| `authClient.signOut()` | `POST /api/auth/sign-out` | MOB-01-07 |

```ts
// src/lib/auth/client.ts
export const authClient: ReturnType<typeof createAuthClient<{ plugins: [...] }>>;
```

## Data

SecureStore keys written by `expoClient` under the `coolinary` prefix (exact key names are Better Auth's). No new
AsyncStorage keys.

## Errors

| Case | Behaviour |
|---|---|
| `EXPO_PUBLIC_API_URL` empty | `authClient` creation throws in dev with a clear message; in a release build the app shows the "can't reach server" state |
| identity rejects the origin (403 from Better Auth) | Logged in dev as a configuration error naming `trustedOrigins` |

## Acceptance criteria

- [ ] `bun run typecheck` and `bun run lint` pass; `bunx expo-doctor` reports no incompatible versions.
- [ ] On a device, `authClient.signIn.anonymous()` against `make dev` creates a session and a SecureStore entry; a
      restart of the app keeps the session.
- [ ] Searching the AsyncStorage contents after sign-in finds no cookie, token or e-mail address.
- [ ] Release builds contain no `exp://` origin assumptions.

## Notes

- Better Auth Expo docs: https://www.better-auth.com/docs/integrations/expo; plugins: anonymous, email-otp, jwt.
