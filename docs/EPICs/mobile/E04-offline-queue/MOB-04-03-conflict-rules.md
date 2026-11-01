# MOB-04-03 Conflict rules

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-04 Offline mutation queue](README.md) |
| Release | R1 |
| Size | M |
| Depends on | MOB-04-02, GW-01 (error codes), GW-08 (list versions), GW-06 |
| Unblocks | MOB-05-02 |

## Goal

When a queued change no longer fits what is on the server (someone else changed or deleted it), the outcome is
predictable and the user is told only when it matters.

## Business rules

1. Rules per answer and kind; the server's state always wins:

   | Answer | Kind | Outcome | Tell the user |
   |---|---|---|---|
   | `404` | `*.delete`, `stock.remove`, `list.line.delete` | success (already gone) | no |
   | `404` | `stock.edit`, `list.line.edit`, `list.line.tick` | drop; refetch | yes: "X was removed on another phone" |
   | `404` | `cooked` (dish retired) | drop | yes |
   | `409` `version_conflict` | `list.line.edit` | drop; put `current` from the error into the cache | yes: "X was changed on another phone" |
   | `409` | `list.line.tick` | never expected (ticks carry no version); treat as success and refetch | no |
   | `409` `already_exists` | `stock.add`, `list.lines.add` with the same client id | success (an earlier attempt landed) | no |
   | `409` `insufficient_stock` | `cooked` | success; server consumed what was there | no |
   | `400` / `422` | any | drop; roll back the optimistic change | yes, with the field |
   | `403` `plus_required` | any | drop; roll back; invalidate entitlements | yes |
   | `403` `not_member` | any | drop; `forceRefresh()` | yes, once for all dropped entries |
2. Dropped entries are reported through one event `outboxDropped({ entry, reason })` so a later UI can show one
   summary instead of many toasts. Messages are translation keys.
3. Ticks are idempotent and order-independent: the last tick state the user set wins.
4. A dropped entry also drops the pending overlay for that change (MOB-04-04).

## Contract

```ts
export type DropReason = "removed_elsewhere" | "changed_elsewhere" | "invalid" | "plus_required" | "not_member";
export function onOutboxDropped(fn: (e: { entry: OutboxEntry; reason: DropReason }) => void): () => void;
```

## Data

No new storage.

## Errors

This task is the error handling for replayed entries; its outcomes are the table in rule 1. Anything not in that table:

| Case | Behaviour |
|---|---|
| `5xx`, timeout or no network | keep the entry; retry with the MOB-04 backoff |
| `401` | pause the queue, refresh the session once (MOB-01); on a second `401`, keep entries and sign-in is required |
| `429` | keep the entry; wait for `Retry-After` before the next attempt |
| Unknown `4xx` code | drop the entry; roll back; report through `outboxDropped` with reason `unknown` |

## Acceptance criteria

- [ ] Phone A deletes a line while phone B, offline, edits its quantity; when B reconnects the line disappears on B
      and B is told once.
- [ ] Phone B offline ticks a line that A already ticked: no message, the line stays ticked.
- [ ] A `stock.add` whose first attempt landed but timed out is not added twice.
