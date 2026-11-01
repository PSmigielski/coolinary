# MOB-03-02 Invalidation map

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | MOB-03-01 |
| Unblocks | MOB-03-03 … MOB-03-09, MOB-04-02, MOB-05-02 |

## Goal

One table says which cached data goes stale after each change, so a cooked dish, a ticked list line or an accepted
invite updates every screen that shows it.

## Business rules

1. `src/api/invalidations.ts` exports `invalidate(kind, ctx)`; mutations, the outbox replay and live notices all call
   it, never `invalidateQueries` directly.
2. The map (household `h`, user `u`):

   | Change | Invalidates |
   |---|---|
   | `stock.add`, `stock.edit`, `stock.remove` | stock, suggestions, list (missing items), savings (R2) |
   | `cooked` | stock, suggestions, history, plans (the meal is marked cooked), savings (R2) |
   | `list.generate.accept`, `list.line.*` | list |
   | `list.tick` | list, stock (ticked items move into the fridge) |
   | `plan.compose`, `plan.swap`, `plan.replace`, `plan.away` | plans, list proposals |
   | `profile.edit` | profile, suggestions, library pages (tool/diet filters), plans |
   | `member.edit` (R2) | profile, plans |
   | `household.*` (rename, profiles, invites, leave, remove; R2) | household, entitlements; leave/remove also triggers `forceRefresh()` |
   | `myShops.edit` (R2) | my shops, offers, price matrix |
   | `favourite`, `rating` | dish, library pages, favourites |
   | `notificationPrefs.edit` | notification preferences |
   | `plan.changed` (JWT plan claim) | entitlements, every Plus-gated query |
   | live `stock.changed` (R2) | stock, suggestions |
   | live `list.changed` with a gap (R2) | list |
3. Invalidation marks queries stale and refetches only active ones (`refetchType: "active"`).
4. While offline, invalidation is deferred: queries refetch on reconnect through `onlineManager`.

## Contract

```ts
export type ChangeKind = "stock.add" | "stock.edit" | "stock.remove" | "cooked" | "list.tick" | /* … table rows */ string;
export function invalidate(kind: ChangeKind, ctx: { householdId: string; userId?: string }): Promise<void>;
```

## Data

No storage.

## Errors

| Case | Behaviour |
|---|---|
| Unknown change kind passed to `invalidate` | throws in development builds; in production logs a warning (no user data) and invalidates nothing |
| Offline when invalidating | deferred until reconnect (rule 4) |
| A refetch after invalidation fails | the query keeps its last data and shows the standard error state (MOB-02 error model); no retry loop beyond TanStack's `retry: 1` |

## Acceptance criteria

- [ ] Unit test: each `ChangeKind` invalidates exactly the keys in the table.
- [ ] Marking a dish cooked updates the fridge list and the "what can I make now" list without a manual refresh.
