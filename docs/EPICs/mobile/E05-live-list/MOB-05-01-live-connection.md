# MOB-05-01 Live connection manager

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-05 Live shopping list client](README.md) |
| Release | R2 |
| Size | M |
| Depends on | MOB-01-04, MOB-02-04, GW-09 (`GET /v1/live` WebSocket) |
| Unblocks | MOB-05-02, MOB-05-03 |

## Goal

While a shared household uses the app, the phone keeps one WebSocket to the gateway, so a line ticked in the shop
appears on the other person's phone within seconds (features.md: "The list is shared live").

## Business rules

1. Connect only when all are true: the app is in the foreground, the user is online, and the household has at least
   2 members with accounts (from `useHousehold()`). Single-person households never open a socket.
2. URL: `API_URL` with `http` → `ws` / `https` → `wss`, path `/v1/live`. Release builds use `wss` only.
3. Authentication (assumed GW-09; align when written): the JWT in an `Authorization: Bearer` header via React
   Native's `WebSocket(url, protocols, { headers })` option. If GW-09 chooses first-message auth instead, send
   `{ "type": "auth", "token": "<jwt>" }` within 5 s of opening.
4. The JWT lives 15 min but the socket longer: 120 s before the token's expiry, after refreshing it (MOB-01-04), the
   client sends `{ "type": "auth", "token": "<new jwt>" }`. The server closes with code `4401` if the token expires
   anyway.
5. Heartbeat: the client sends `{ "type": "ping" }` every 25 s; if no `pong` or other message arrives for 60 s, the
   socket is treated as dead, closed and reconnected.
6. Reconnect on any close except `1000` sent by the client: delays 1 s, 2 s, 4 s, 8 s, 16 s, then 30 s steady,
   ± 20 % jitter; the counter resets after the connection has been up for 60 s.
7. Close codes from the server (assumed GW-09): `4401` → `forceRefresh()` then reconnect at once; `4403` (not a
   member) → `forceRefresh()`, `householdChanged` follows, connect to the new household if rule 1 holds; `1001` or
   `1012` (server restart) → reconnect with the backoff.
8. App goes to the background → close with `1000` after 10 s in the background (so a quick app switch doesn't
   reconnect); foreground → connect again and do a full list refetch (MOB-05-02 rule 4).
9. `householdChanged` → close and reopen for the new household.
10. Nothing received over the socket is persisted except through the normal query cache rules.

## Contract

Client → server messages: `auth`, `ping`. Server → client: `hello { connectionId }`, `pong`, and the change notices
of MOB-05-02 and MOB-05-03.

```ts
export function useLiveConnection(): { state: "off" | "connecting" | "open" | "reconnecting"; connectionId?: string };
```

## Data

Memory only.

## Errors

| Case | Behaviour |
|---|---|
| Gateway unreachable | Reconnect backoff; queries still work by polling (`staleTime` 30 s) |
| Too many sockets for the user (server close 4429) | Wait 60 s before reconnecting |

## Acceptance criteria

- [ ] Two phones in a shared household: a tick on phone A appears on phone B in under 2 s on Wi-Fi.
- [ ] Toggling airplane mode on B reconnects within 2 s of the network returning, and B shows ticks made meanwhile.
- [ ] A socket open for 40 minutes keeps receiving changes (token re-auth works).
- [ ] A single-person household opens no socket (checked in the gateway logs).
