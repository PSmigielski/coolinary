# GW-03-01 Quota counters with reserve / commit / release

| | |
|---|---|
| Service | gateway |
| Epic | [GW-03 Free-tier AI quota and rate limits](README.md) |
| Release | R1 |
| Size | M |
| Depends on | GW-01-03, IDN-05 (`GetEntitlements`) |
| Unblocks | GW-03-02, GW-06-04/05/06/07, GW-07-04 |

## Goal

Count Free AI use per household exactly: two parallel requests can't both take the last unit, a failed or
non-AI answer costs nothing, and a retry never counts twice.

## Business rules

1. Two quota kinds (features.md §6, business model §3):
   - `dayPlan`: `POST /v1/plans/day` (GW-07-04).
   - `scan`: `POST /v1/intake/photo` (fridge or shelf, GW-06-04), `POST /v1/intake/label` (GW-06-05),
     `POST /v1/intake/text` (GW-06-06, R3), `POST /v1/intake/receipt` (GW-06-07, R3).
   - Own-flyer scans (`POST /v1/flyers`, GW-10-04) are Plus-only and **never** count towards `scan`.
2. Limit for Free and anonymous households: **2 per kind in any 7 days**. Decided (Q025): the limits come from
   identity's `GetEntitlements` (`aiDayPlansPerWeek`, `photoScansPerWeek`), read through the entitlements cache key
   `gw:ent:<householdId>` in redis-cache (60 s, the key GW-02-04 also uses; never older than 5 minutes); config `FREE_DAY_PLANS_PER_WEEK` / `FREE_SCANS_PER_WEEK` (2) is used only
   when identity is down and nothing is cached. Plus: no limit; units are still counted for metrics.
   Decided (Q112): every AI reading request counts as one `scan`: fridge and shelf photos, label photos of unknown
   products, typed or spoken sentences (R3) and receipts (R3).
3. **Window (decided, Q111): rolling 7 days.** A committed use counts for exactly 7 × 24 h from its commit time and
   then frees up; there is no fixed reset day and no time zone. `used` = committed uses in `(now − 7 days, now]`.
   `nextFreeAt` = the oldest counted use + 7 days (when the next slot frees), or `null` when `used = 0`.
4. The quota belongs to the household, not the user: after joining a Plus household, the new member is unlimited;
   after leaving, the person's new personal household starts its own count.
5. **Reserve** (atomic Lua script on `redis-state`):
   1. Remove expired members from the pending set (score < now) and uses older than 7 days from the used set
      (score ≤ now − 7 days).
   2. If `|used| + |pending| >= limit` → reject (`quota_reached`).
   3. Otherwise add a reservation id (UUIDv7) to the pending set with score `now + deadline + 10 s`.
6. **Commit** when the service's response has `ai_used = true` — for day plans when planner's `quota_consumed = true`
   (AI used, at least one AI meal kept, not a replay, PLN-08-07 rule 5): `ZADD used <now> <id>`, `ZREM pending <id>`,
   in one script.
   **Release** in every other case (`ai_used = false`, any error, timeout, client gone): `ZREM pending <id>`.
7. A reservation not committed or released (gateway crash) expires by its score and stops counting.
8. Plus households skip reserve/commit for enforcement but run `ZADD used <now> <id>` after a successful AI
   response (for GW-03-05).
9. `redis-state` unreachable: Free → 503 `service_unavailable` for AI routes (the cost cap is never bypassed);
   Plus → the route runs, counting is skipped and a metric records it.
10. Services must report `ai_used` on every AI-capable response (planner `ComposeDayPlan`, kitchen intake RPCs). A
    missing field is treated as `true` on success (counts), so a service bug can't give free AI away.

## Contract

Internal module `src/quota.ts`:

```ts
type QuotaKind = "dayPlan" | "scan";
reserve(householdId, kind, plan, deadlineMs): Promise<{ ok: true; reservationId: string } | { ok: false; status: QuotaStatus }>
commit(householdId, kind, reservationId): Promise<void>
release(householdId, kind, reservationId): Promise<void>
status(householdId, kind, plan): Promise<QuotaStatus>   // { limit, used, remaining, nextFreeAt: string | null }
```

## Data (redis-state)

| Key | Type | Value | TTL |
|---|---|---|---|
| `gw:quota:<householdId>:<kind>:used` | sorted set | committed reservation ids, score = commit time ms | 8 days, refreshed on every commit |
| `gw:quota:<householdId>:<kind>:pending` | sorted set | reservation ids, score = expiry ms | 1 day, refreshed on every reserve |

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Free limit reached | 429 | `quota_reached` (GW-03-02) |
| redis-state down (Free) | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] Free household: two successful day plans succeed, the third gets 429 `quota_reached`.
- [ ] A day plan where planner fell back to the library (`ai_used: false`) leaves `used` unchanged.
- [ ] Ten parallel requests with `used = 1` produce exactly one success and nine `quota_reached` (load test).
- [ ] A planner timeout releases the reservation: `remaining` is back to its old value immediately.
- [ ] A reservation left by a killed gateway stops counting after deadline + 10 s.
- [ ] Rolling window (time-travel unit test): uses at day 0 10:00 and day 2 09:00 → a third is refused until day 7
      10:00, when exactly one slot frees; `nextFreeAt` reads day 7 10:00 before that and day 9 09:00 after it.
- [ ] A label photo and a typed sentence each count as one `scan`.
- [ ] Plus household: no 429 after 10 day plans; `used` reads 10.
- [ ] With redis-state stopped: Free AI route → 503; Plus AI route → 200.
- [ ] Verification 11 (features.md): the counters give the AI calls per Free household per 7 days used in the cost check.

## Notes

- Decided (Q111): rolling 7-day window, not a Monday reset. The app shows "next free plan on <nextFreeAt>".
- Decided (Q112): every AI reading request counts as one scan (rule 2).
