# MOB-02-04 Connectivity and server reachability

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-02 API client and error model](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | MOB-01-01 (installs `expo-network`) |
| Unblocks | MOB-02-02, MOB-04-02, MOB-05-01 |

## Goal

The app tells "you're offline" apart from "our server is down", pauses queries while offline, and replays queued
changes as soon as the connection is back.

## Business rules

1. `src/lib/network.ts` subscribes with `Network.addNetworkStateListener` (`expo-network`, in Expo Go) and reads the
   first state with `Network.getNetworkStateAsync()` at start.
2. Online means `isConnected && isInternetReachable !== false`. (On iOS `isInternetReachable` equals
   `isConnected`; on Android it is the validated-network flag.)
3. The state feeds TanStack Query's `onlineManager.setEventListener`, so queries pause offline and refetch on
   reconnect, and the outbox replay (MOB-04-02) starts on reconnect.
4. Changes are debounced by 1 s so a flapping connection doesn't trigger many refetches.
5. **Reachability probe.** When a request fails with a network error while the phone is online, the module probes
   `GET /healthz` on the gateway with a 3 s timeout. Probe fails → `serverReachable = false` and errors become
   `serverDown`; probe succeeds → the original error stands. While unreachable, the probe repeats every 30 s in the
   foreground; first success sets `serverReachable = true` and resumes queries and the outbox.
6. Exposes `useConnectivity(): { online: boolean; serverReachable: boolean }` for a later offline banner.

## Contract

`GET /healthz` (gateway, public, already exists). No new route.

## Data

Memory only.

## Errors

None of its own.

## Acceptance criteria

- [ ] Airplane mode on: queries stop; off: they refetch within 2 s, and queued changes start replaying.
- [ ] `make down` with Wi-Fi on: the next failed request gives `serverDown`; after `make dev`, queries resume
      within 30 s without user action.
