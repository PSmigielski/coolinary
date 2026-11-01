# Coolinary mobile

Expo (SDK 58) + Expo Router app. Runs in **Expo Go**: no native build needed.
Standalone Bun project with its own `bun.lock` (not part of a workspace).

## Run

```bash
bun install
bun run api:types     # regenerate src/api/schema.d.ts from ../contracts/openapi/gateway.json
bunx expo start       # then scan the QR code with Expo Go (phone on the same Wi-Fi)
```

Checks: `bun run typecheck`, `bun run lint`, `bunx expo-doctor`.

## API URL

The app calls the gateway at `EXPO_PUBLIC_API_URL`, read from `mobile/.env` and
inlined into the bundle at build time (restart `expo start` after changing it).
A phone can't use `localhost`: use this machine's LAN IP plus `GATEWAY_HOST_PORT`
from the repo-root `.env`, e.g. `http://192.168.1.43:3080`. Put per-machine
overrides in `.env.local` (git-ignored). See `.env.example`.

## Layout

```
src/app/        routes (Expo Router): screens and layouts only
src/api/        generated OpenAPI types + typed client ($api)
src/features/   one folder per feature area: kitchen, planner, offers, household
src/stores/     Zustand stores (client state)
src/lib/        query client, secure storage, small helpers
```

## Libraries and why

- **TanStack Query** for server state: caching, retries, refetch on pull and on
  app foreground. The cache is persisted to AsyncStorage
  (`@tanstack/react-query-persist-client`), so the last data shows instantly after a reload.
- **openapi-typescript + openapi-fetch + openapi-react-query**: types generated from
  the gateway's OpenAPI document, so paths, params and responses are type-checked.
  Use `$api.useQuery("get", "/v1/...")` in components.
- **Zustand** for client state (UI preferences such as kcal/kJ), persisted with
  `persist` + AsyncStorage.
- **AsyncStorage** instead of MMKV because MMKV needs a development build and
  doesn't run in Expo Go. Swap in `react-native-mmkv` once we move to a dev build.
- **expo-secure-store** (`src/lib/secure.ts`) for sensitive data. Health data (age,
  height, weight, goals) stays on the device, only in secure storage. The only
  value sent to the server is the daily kcal target.
