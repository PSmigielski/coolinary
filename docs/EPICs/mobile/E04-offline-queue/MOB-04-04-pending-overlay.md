# MOB-04-04 Pending changes overlay

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-04 Offline mutation queue](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | MOB-04-01, MOB-03-01 |
| Unblocks | MOB-03-05, MOB-03-07 |

## Goal

Changes waiting in the outbox stay visible on every screen, even after a refetch or an app restart replaces the
cached server data.

## Business rules

1. The overlay is a pure function per data type: `applyPending(serverData, pendingEntries) → displayedData`. It is
   applied in the `select` of the stock, list, profile, favourites and plans queries, never written into the cache.
2. Per kind:
   - `stock.add` / `list.lines.add`: append the item if its client id is not in the server data;
   - `stock.edit` / `list.line.edit` / `list.line.tick`: patch the matching item;
   - `stock.remove` / `list.line.delete`: hide the item;
   - `cooked`: no overlay on stock (the server computes consumption); the plan meal shows as cooked;
   - `profile.edit`, `member.edit`: the queued body replaces the matching section;
   - `favourite`, `rating`: patch the dish.
3. Items changed by a pending entry carry `pending: true` so a later UI can mark them.
4. Entries with status `failed` are still applied until the user discards them.
5. The overlay is computed from the outbox in memory (loaded once at start, kept in sync on every outbox change).
6. When an entry succeeds, the cache is updated from the response before the overlay entry disappears, so nothing
   flickers.

## Contract

```ts
export function applyPending<T>(kindGroup: "stock" | "list" | "profile" | "dishes" | "plans", data: T, entries: OutboxEntry[]): T;
```

## Data

Memory only.

## Errors

None of its own: a pure function over server data and the outbox. An entry whose kind it doesn't know is not overlaid and is logged once in development builds.

## Acceptance criteria

- [ ] Offline: add 2 fridge items, restart the app, pull to refresh while still offline and then online: the items
      are visible throughout and appear once after the replay.
- [ ] A pending tick stays ticked after a background refetch of the list.
