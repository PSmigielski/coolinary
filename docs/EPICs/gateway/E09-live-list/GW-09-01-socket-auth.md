# GW-09-01 `/v1/live` socket: auth, token refresh, limits, heartbeat

| | |
|---|---|
| Service | gateway |
| Epic | [GW-09 Live shopping list](README.md) |
| Release | R2 |
| Size | M |
| Depends on | GW-02-01, GW-03-04 |
| Unblocks | GW-09-02, MOB-05 |

## Goal

A WebSocket the app keeps open while it is in the foreground, authenticated with the same 15-minute JWT, without
putting tokens in URLs.

## Business rules

1. `GET /v1/live` upgrades with Elysia `.ws()` (Bun's WebSocket server). The upgrade itself needs no token; the
   first message must be `{"type":"auth","token":"<JWT>"}` within **5 seconds**, else close **4401**.
2. Token check as GW-02-01. Valid → the socket is bound to the token's `householdId`, `memberId`, `userId`; server
   sends `{"type":"ready","householdId","heartbeatSeconds":25}`.
3. **Refresh:** the client sends a new `auth` message before `exp`. Same household → keep going. Different
   household (the user joined or left) → close **4403**; the client reconnects with the new token. At `exp` without a
   refresh → close **4401**.
4. **Limits:** at most 5 open sockets per user and 30 per household; one more → close **4429** (oldest is kept).
   Upgrades are rate limited by bucket `ws_connect` (GW-03-04). Incoming messages: at most 10 per minute, each
   ≤ 4 KB; anything other than `auth` and `ping` is ignored; over the rate → close 4429.
5. **Heartbeat:** the server sends a WebSocket ping every 25 s and closes the socket if no pong arrives within 10 s.
   The client may also send `{"type":"ping"}`; the server answers `{"type":"pong"}`.
6. On gateway shutdown (SIGTERM), every socket is closed with **1001** before the process exits, so clients
   reconnect to the next instance.
7. Any account type, Free or Plus (a single user with two phones benefits too).
8. Metrics: gauge `gateway_ws_connections`, counter `gateway_ws_closes_total{code}`. No household ids in labels.
9. Logs: open and close with `user_hash` and close code; never the token.

## Contract

WebSocket `GET /v1/live`, documented in the gateway README (OpenAPI lists the route with a description only).

## Data

In memory per gateway process: `Map<householdId, Set<socket>>`, `Map<userId, Set<socket>>`.

## Errors

| Case | Close code |
|---|---|
| No auth within 5 s, bad or expired token | 4401 |
| Household changed on refresh | 4403 |
| Too many sockets or messages | 4429 |
| Gateway shutting down | 1001 |

## Acceptance criteria

- [ ] A socket that sends no auth message is closed with 4401 after 5 s.
- [ ] A socket whose token expires without a refresh is closed with 4401 at `exp`.
- [ ] Refreshing with a token for another household closes with 4403.
- [ ] The sixth socket of one user is refused with 4429.
- [ ] A client that stops answering pings is dropped within 35 s.
