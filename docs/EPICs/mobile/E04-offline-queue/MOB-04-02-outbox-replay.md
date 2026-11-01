# MOB-04-02 Outbox replay

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-04 Offline mutation queue](README.md) |
| Release | R1 |
| Size | M |
| Depends on | MOB-04-01, MOB-02-02, MOB-02-04, MOB-01-03, MOB-01-04, MOB-03-02, GW-01-03 (24 h replay window) |
| Unblocks | MOB-04-03, MOB-01-05, MOB-01-07 |

## Goal

Queued changes reach the server in the order the user made them, exactly once, as soon as there is a connection.

## Business rules

1. Replay starts when: the app goes online (MOB-02-04), the app comes to the foreground, a new entry is added while
   online, the server becomes reachable again, and every 30 s while entries are pending and the app is in the
   foreground. It never runs in the background.
2. One replay loop at a time (single-flight). Entries are sent **one by one, oldest first, per household**; entry
   n+1 is sent only after entry n has a final outcome. This keeps "add line" before "tick line".
3. Before sending, the loop ensures a session (MOB-01-03 for entries with `userId: null`, then fills `userId` and
   `householdId` of those entries from the new JWT) and a token (MOB-01-04).
4. Entries whose `householdId` differs from the current household are sent only if the user is still a member there
   (the gateway answers `not_member` otherwise → the entry is dropped and the user told). Entries whose `userId`
   differs from the current user are not sent; they are dropped after 24 h (MOB-01-06 rule 3).
5. Each attempt sends `Idempotency-Key = entry.idempotencyKey` and the default 15 s timeout.
6. Outcome per response:

   | Response | Outcome |
   |---|---|
   | 2xx | remove entry; `invalidate(kind)` (MOB-03-02); apply the response to the cache |
   | `offline`, `timeout`, `serverDown`, `rateLimited` | keep; stop the loop; retry later with backoff |
   | `401` | MOB-01-06 handles the token; then retry once |
   | `409`, `404`, `400`/`422`, `403` | conflict rules (MOB-04-03) |
   | `500` | keep; backoff; after 5 attempts mark `failed` |
7. Backoff for kept entries: 2 s, 4 s, 8 s, 16 s, … doubling, at most 5 min, ± 20 % jitter; reset to 2 s after any
   success.
8. `failed` entries stay visible (count badge) and can be retried or discarded by the user; they no longer block the
   entries behind them.
9. `flush(householdId, timeoutMs)` sends all pending entries now and resolves when the outbox is empty or the
   timeout passes; used by sign-in, invite accept and sign-out.
10. **Maximum age (decided, Q021): 24 hours.** Before each replay loop, entries whose `createdAt` is more than
    24 hours old are removed without sending (the gateway keeps idempotency keys for 24 h, GW-01-03, so a later
    replay could apply twice). Their optimistic cache changes are rolled back and the user sees one summary
    ("3 offline changes were older than a day and were not saved"). `failed` entries follow the same rule.

## Contract

```ts
export function startReplay(): void;
export function flush(householdId: string, timeoutMs: number): Promise<{ remaining: number }>;
```

## Data

Updates entries in `coolinary.outbox.v1`.

## Errors

See rule 6.

## Acceptance criteria

- [ ] Add line, tick it, edit fridge item offline; on reconnect the gateway receives them in that order.
- [ ] Killing the app while an entry is `sending` and restarting sends it again with the same `Idempotency-Key`; the
      server applies it once.
- [ ] Gateway stopped: entries stay pending with growing delays; after restart they go out within 30 s.
- [ ] An entry created 25 hours ago (test clock) is removed without a request, its optimistic change is rolled
      back and the summary message appears once; an entry 23 hours old is sent.
