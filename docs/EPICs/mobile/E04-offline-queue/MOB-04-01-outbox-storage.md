# MOB-04-01 Outbox: storage and allowed changes

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-04 Offline mutation queue](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | MOB-01-02, MOB-02-01 |
| Unblocks | MOB-04-02, MOB-04-03, MOB-04-04, MOB-03-04 … MOB-03-09 |

## Goal

Changes the user makes without signal (in the shop, in the cellar) are saved on the phone and survive an app
restart, and only changes that are safe to send later can be queued.

## Business rules

1. The outbox is an ordered list persisted in AsyncStorage key `coolinary.outbox.v1`, written after every change
   (append, status update, removal). It is separate from the query cache and survives app updates.
2. Entry shape:

   ```ts
   interface OutboxEntry {
     id: string;                 // UUIDv4
     idempotencyKey: string;     // sent as Idempotency-Key on every attempt
     kind: OutboxKind;
     userId: string | null;      // null only before the anonymous account exists (MOB-01-03)
     householdId: string | null;
     method: "POST" | "PUT" | "PATCH" | "DELETE";
     path: string;               // concrete path, e.g. /v1/list/lines/7f…
     body?: unknown;
     createdAt: string;          // ISO 8601, phone time
     attempts: number;
     nextAttemptAt?: string;
     status: "pending" | "sending" | "failed";
     lastError?: { kind: string; code?: string };
   }
   ```
3. **Allowed kinds** (everything else is sent online only and fails with `offline`):

   | Kind | Route |
   |---|---|
   | `stock.add`, `stock.edit`, `stock.remove` | `POST /v1/stock`, `PATCH`/`DELETE /v1/stock/{id}` (with `If-Match`) |
   | `product.name`, `product.confirm` | `POST /v1/products`, `POST /v1/products/{ean}/confirm` |
   | `list.lines.add`, `list.line.edit`, `list.line.tick`, `list.line.untick`, `list.line.delete`, `list.clearTicked`, `list.basePantry`, `list.offerLine` | `POST /v1/list/lines`, `PATCH`/`DELETE /v1/list/lines/{id}`, `POST /v1/list/lines/{id}/tick\|untick`, `DELETE /v1/list/lines?state=ticked`, `POST /v1/list/base-pantry`, `POST /v1/list/offer-lines` |
   | `cooked` | `POST /v1/cooked` |
   | `profile.edit`, `member.edit` (R2) | `PATCH /v1/profile`, `PUT /v1/profile/members/{id}` |
   | `favourite`, `rating`, `ratingPrompt.dismiss` | `PUT`/`DELETE /v1/dishes/{id}/favourite`, `PUT`/`DELETE /v1/dishes/{id}/rating`, `POST /v1/dishes/{id}/rating-prompt/dismiss` |
   | `plan.away` (R2) | `PUT /v1/attendance` |
   | `myShops.set`, `shoppingPrefs.set`, `notificationPrefs.set` | GW-10, GW-11 |

   Never queued: any AI call (plans, swaps, rewrites, intake, flyers), sign-in, invites, leave/remove, account
   deletion, export, uploads.
4. Limits: at most 500 entries and 1 MB serialised; an entry lives at most 24 hours (decided, Q021; MOB-04-02
   rule 10). Adding beyond either limit fails with `outbox_full` and the UI
   asks the user to go online. A single body is at most 64 KB.
5. Coalescing when an entry is added, only among `pending` entries of the same household:
   - a newer `list.line.tick` for the same line replaces the older one;
   - a newer `profile.edit` replaces the older one (whole-profile `PUT`);
   - a `*.delete`/`stock.remove` of something created by a still-pending add removes both entries;
   - nothing else is merged.
6. Entries are never encrypted; they must not contain health data. `member.edit` bodies pass the privacy guard of
   MOB-03-04 rule 3 before being stored.
7. `enqueue()` returns at once; the caller updates the cache optimistically (MOB-04-04).

## Contract

```ts
export function enqueue(entry: Omit<OutboxEntry, "id" | "attempts" | "status" | "createdAt">): OutboxEntry;
export function pending(householdId: string): OutboxEntry[];
export function useOutboxCount(): number;   // for a later "3 changes not sent yet" badge
```

## Data

AsyncStorage `coolinary.outbox.v1`, versioned; a future shape change migrates or drops it with a user notice.

## Errors

| Case | Behaviour |
|---|---|
| 500 entries or 1 MB reached | `outbox_full` |
| AsyncStorage write fails | The change is not accepted; the user sees "couldn't save" |

## Acceptance criteria

- [ ] 20 changes made in airplane mode survive a force-quit and are all pending after restart.
- [ ] Ticking and unticking a line 5 times offline leaves one entry.
- [ ] Adding then deleting a fridge item offline leaves no entry.
- [ ] Trying to compose a day plan offline gives `offline` and adds no entry.
