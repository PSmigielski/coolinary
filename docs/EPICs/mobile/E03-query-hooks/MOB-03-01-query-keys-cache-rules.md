# MOB-03-01 Query keys and cache rules

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | MOB-01-02, MOB-02-02 |
| Unblocks | MOB-03-02 … MOB-03-09, MOB-01-07 |

## Goal

Every cached response has a predictable key scoped to its household, a stated freshness, and an explicit rule on
whether it may be written to the phone's unencrypted storage.

## Business rules

1. One key factory `src/api/keys.ts`. Every household-scoped key starts with `["h", householdId]`; user-scoped keys
   start with `["u", userId]`; shared reference data (catalogue, presets, dishes) starts with `["ref", language]`.

   ```ts
   export const keys = {
     me: (u: string) => ["u", u, "me"] as const,
     stock: (h: string) => ["h", h, "stock"] as const,
     list: (h: string) => ["h", h, "list"] as const,
     dish: (lang: string, id: string) => ["ref", lang, "dish", id] as const,
     // … one entry per hook in MOB-03-03 … MOB-03-09
   };
   ```
2. Hooks use `$api.useQuery(...)` with `queryKey` overridden by the factory, so invalidation never depends on
   openapi-react-query's internal key shape.
3. Freshness defaults (overridable per hook, listed in each hook task):

   | Data | `staleTime` | Persist |
   |---|---|---|
   | stock, list, plans, profile | 30 s | yes |
   | household, members (display names) | 5 min | yes |
   | entitlements | 5 min | yes |
   | dishes, library pages, presets, cuisines | 24 h | yes (first library page, last 50 dishes opened) |
   | barcode lookups | 7 days | yes (last 500) |
   | ingredient autocomplete | 1 h | no |
   | offers, shops | 1 h | yes (current week only) |
   | price matrix, intake proposals, invites, export links | 0 | **no** |
4. **Persistence allow-list.** `persistOptions.dehydrateOptions.shouldDehydrateQuery` persists a query only when its
   hook sets `meta: { persist: true }` and the query succeeded. Never persisted: invites (their codes work as
   passwords for 48 h), intake proposals, price matrix, export links, anything with an e-mail address.
5. **Calorie targets.** Member profiles that carry a daily kcal target are never persisted in the query cache or
   AsyncStorage (decided, Q045: a target is derived from health data and AsyncStorage is unencrypted). Offline
   portion display reads them from secure storage or refetches (a later frontend task).
6. `gcTime` stays 24 h (≥ `maxAge`, as in `src/lib/query.ts`). `maxAge` of the persisted cache stays 24 h.
7. The persister `buster` becomes `v<cacheVersion>-<app version>` (`expo-constants`), so a new app release starts
   with an empty persisted cache; the outbox (MOB-04) is separate and survives releases.
8. On `householdChanged(old)` remove every key under `["h", old]`; on `userChanged(old)` remove `["u", old]`; on wipe
   (MOB-01-07) `queryClient.clear()` and remove the persisted client.
9. Queries that need auth set `enabled: status !== "booting" && householdId != null`.

## Contract

```ts
export const keys: { /* factory */ };
export const queryMeta = { persist: { persist: true } } as const;
```

## Data

AsyncStorage `coolinary.query-cache` (existing), now filtered by the allow-list.

## Errors

None of its own.

## Acceptance criteria

- [ ] After using invites and a price matrix, the persisted cache contains neither (inspect AsyncStorage).
- [ ] Switching household removes all old-household entries from memory and from the persisted copy.
- [ ] Installing a new app version starts with an empty persisted cache but keeps the outbox.
