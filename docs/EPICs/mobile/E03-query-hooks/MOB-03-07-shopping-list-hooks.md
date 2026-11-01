# MOB-03-07 Shopping list hooks

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R1 (list), R2 (offer badges shown from the response) |
| Size | M |
| Status | todo |
| Depends on | MOB-03-01, MOB-03-02, MOB-04-01, MOB-04-04, GW-08 |
| Unblocks | MOB-05-02, later list screens |

## Goal

Hooks for the shopping list (features.md "Shopping list"): proposals from the plan, manual lines, ticking in the shop
with the item moving into the fridge, all instant and usable without signal in the shop.

## Business rules

1. Hooks and routes (GW-08 route catalogue in docs/EPICs/gateway/README.md):

   | Hook | HTTP | Offline |
   |---|---|---|
   | `useList()` | `GET /v1/list` → `{ id, version, lines[] }` | cache + overlay |
   | `useGenerateList()` | `POST /v1/list/proposals` `{ from, to, meals? }` → proposal with keyed lines (kept 24 h by planner) | never |
   | `useAcceptProposal()` | `POST /v1/list/proposals/{id}/accept` `{ wantedKeys }` | never (the proposal expires) |
   | `useAddLine()` | `POST /v1/list/lines` `{ ingredientId \| name, quantity? }` (one line per call) | queued |
   | `useEditLine()` | `PATCH /v1/list/lines/{id}` `{ quantity, baseVersion? }` | queued |
   | `useTickLine()` / `useUntickLine()` | `POST /v1/list/lines/{id}/tick` `{ tickedAt, shopId? }` / `POST /v1/list/lines/{id}/untick` | queued |
   | `useDeleteLine()` | `DELETE /v1/list/lines/{id}` | queued |
   | `useClearTicked()` | `DELETE /v1/list/lines?state=ticked` | queued |
   | `useAddBasePantry()` | `POST /v1/list/base-pantry` `{ cuisine }` | queued |
   | `useAddOfferToList()` (R2) | `POST /v1/list/offer-lines` `{ offerId }` | queued |
2. Lines get client-generated UUIDs on creation; the server accepts them (decided, Q116; PLN-15, GW-08).
3. The list has a `version` that increases with every change on the server; each line has its own `version`.
   Quantity edits send `baseVersion` (the line version the user saw); planner applies them last-writer-wins and
   returns the line as it now is (GW-08-01 rule 3). Ticks and deletes send no version: a tick is idempotent and a
   delete of a deleted line is a success (MOB-04-03).
4. Ticking is optimistic and immediate; the server moves the item into the fridge with a default use-by date
   (features.md list rule 7), so a tick also invalidates stock (MOB-03-02).
5. Proposed lines from `useGenerateList()` are kept in memory until the user accepts; only the wanted keys are sent
   (`useAcceptProposal`). An expired proposal (410) is generated again.
6. Each line shows its source (meal, cuisine base, who added it) and, from R2, its offer badge, exactly as returned
   by the server; the phone does no merging or rounding.
7. Shared lists (R2) also receive live changes (MOB-05-02); a local optimistic change and the server echo of the same
   change are matched by `Idempotency-Key` so the line doesn't flicker.

## Contract

```ts
interface LineInput { id: string; ingredientId?: string; text: string; quantity?: { value: number; unit: "gram" | "millilitre" | "piece" }; source: { kind: "meal" | "cuisine_base" | "manual" | "offer"; refId?: string } }
```

## Data

Query cache: list persisted; proposals memory only.

## Errors

| Case | Behaviour |
|---|---|
| `409` on edit (line changed elsewhere) | MOB-04-03: server copy wins, user told |
| Generate offline | `offline`; the existing list still works |

## Acceptance criteria

- [ ] Verification 12 from the app side: generating from a test plan returns exactly the missing ingredients; ticking
      one moves it into `GET /v1/stock`.
- [ ] Ticking 10 lines in airplane mode shows all ticks at once and sends them on reconnect, in order.
