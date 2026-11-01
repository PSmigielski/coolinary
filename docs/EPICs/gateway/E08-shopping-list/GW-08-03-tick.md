# GW-08-03 Tick in the shop → into the fridge

| | |
|---|---|
| Service | gateway |
| Epic | [GW-08 Shopping list routes](README.md) |
| Release | R1 |
| Size | S |
| Depends on | GW-08-01, PLN-15 (`TickListLine`, `UntickListLine`) |
| Unblocks | MOB-04 (ticks queued offline in the shop) |

## Goal

Ticking an item in the shop moves it into the fridge with a default use-by date (features.md §2 step 7, the fastest
way to stock the fridge, §4).

## Business rules

1. `POST /v1/list/lines/{lineId}/tick` body `{ "tickedAt", "shopId"? }` → planner `TickListLine(household_id,
   member_id, client_mutation_id, line_id, shop_id, ticked_at)`, which adds the bought item to the fridge through
   kitchen `AddStockItems` with a default use-by date (features.md list rule 7). One line per call; the offline queue
   (MOB-04) sends a shop visit as a series of calls in order.
2. `tickedAt` RFC 3339 device time; planner stores it if it is within the last 7 days and not in the future,
   otherwise server time (PLN-15-08 rule 9).
3. Ticking is idempotent: the idempotency key is planner's `client_mutation_id`, so a replayed tick adds nothing
   twice. A tick on a line removed in the last 7 days restores it as ticked ("bought beats removed").
4. Response: the line, `addedToFridge`, and `fridgePending` (kitchen was down; planner adds the item when it is back).
5. `POST /v1/list/lines/{lineId}/untick` → `UntickListLine`. Within 24 h of the tick it also removes the stock item
   the tick created, unless that item was already used or changed (`stockKept`, PLN-15-07 rule 6).

## Contract

| Method | Path | 2xx |
|---|---|---|
| POST | `/v1/list/lines/{lineId}/tick` | 200 `{ "version", "line": line, "addedToFridge": true, "fridgePending": false }` |
| POST | `/v1/list/lines/{lineId}/untick` | 200 `{ "version", "line": line }` |

Tag `list`, `auth: true`, idempotency required.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| `tickedAt` not RFC 3339 | 400 | `validation` |
| Line not in this household (and not removed in the last 7 days) | 404 | `not_found` |
| Planner down | 503 | `service_unavailable` (kitchen down is not an error: `fridgePending`) |

## Acceptance criteria

- [ ] Ticking "jajka" adds eggs to the fridge with the default date (verification 12: "ticking an item in the shop
      adds it to the fridge").
- [ ] Sending the same tick twice (same key) adds the item once.
- [ ] Untick within 24 hours removes the stock item added by the tick, unless it was already used (`stockKept`).
- [ ] With kitchen stopped, a tick returns `fridgePending: true` and the item appears in the fridge once kitchen is back.
