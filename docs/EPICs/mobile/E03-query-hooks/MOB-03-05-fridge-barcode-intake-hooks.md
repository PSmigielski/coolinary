# MOB-03-05 Fridge, barcode and intake hooks

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-03 Query keys, hooks and cache rules](README.md) |
| Release | R1 (stock, barcode, unknown products, photo, label), R3 (typed text, voice, receipt) |
| Size | L |
| Status | todo |
| Depends on | MOB-03-01, MOB-03-02, MOB-04-01, MOB-04-04, MOB-07-02, GW-06 |
| Unblocks | later fridge, scan and unpack screens |

## Goal

Hooks for stocking the fridge in every way features.md §4 lists (tick in shop, unpack mode, receipt, shelf photo,
typing), with instant optimistic updates that also work offline.

## Business rules

1. Hooks and routes (GW-06 route catalogue in docs/EPICs/gateway/README.md):

   | Hook | HTTP | Rel. | Offline |
   |---|---|---|---|
   | `useStock()` | `GET /v1/stock` | R1 | cache + pending overlay (MOB-04-04) |
   | `useAddStock()` | `POST /v1/stock` `{ items: StockItemInput[], source, intakeRef? }` (1–100 items, used by typing, unpack and intake confirm) | R1 | queued |
   | `useEditStock()` | `PATCH /v1/stock/{id}` (quantity, storage, expiry), `If-Match` item version | R1 | queued |
   | `useRemoveStock()` | `DELETE /v1/stock/{id}?reason=used\|thrown\|mistake`, `If-Match` item version | R1 | queued |
   | `useStockPreview()` | `POST /v1/stock/preview` `{ source, candidates }` (unpack mode review list) | R1 | never |
   | `useStockDefaults()` | `GET /v1/stock/defaults?ingredientId=\|ean=&storage=` (pre-fill the add form) | R1 | cached 1 h |
   | `useResolveBarcode(code)` | `GET /v1/products/resolve?code=` (EAN or GS1 2D content, URL-encoded) | R1 | cache only |
   | `useHouseholdProducts()` | `GET /v1/products/mine?updatedSince=` (offline copy of this household's products) | R1 | persisted |
   | `useNameProduct()` | `POST /v1/products` `{ ean, name, packSize? }` | R1 | queued |
   | `useConfirmProduct()` | `POST /v1/products/{ean}/confirm` `{ fingerprintId, accept, edits? }` | R1 | queued |
   | `useReportProduct()` | `POST /v1/products/{ean}/report` `{ reason }` | R1 | queued |
   | `usePhotoIntake()` | `POST /v1/intake/photo` (multipart, MOB-07-02) `kind=fridge\|shelf` | R1 | never; AI quota |
   | `useLabelIntake()` | `POST /v1/intake/label` (multipart) `ean` | R1 | never; AI quota |
   | `useTextIntake()` | `POST /v1/intake/text` `{ text }` | R3 | never; AI |
   | `useReceiptIntake()` | `POST /v1/intake/receipt` (multipart, 1–3 images) | R3 | never; AI quota |
2. Intake routes return **proposed items** that are not saved (kitchen stores nothing, KIT-08-01). They are kept in
   memory only (no persistence) until the user confirms; confirmation calls `useAddStock()` with the ticked items
   and the response's `intakeRef`. Unpack mode collects scanned barcodes
   on the phone and confirms them with one `useAddStock()` call.
3. Stock items get a client-generated id (UUIDv4) when created, so offline edits can refer to them before the server
   has answered. Decided (Q116): the server accepts phone-made UUIDs for new stock items on `POST /v1/stock`
   (KIT-07, GW-06); no id remapping after replay.
4. Text intake (R3, decided Q219: typed sentences ship together with voice): 1–2,000 characters (the ai service
   limit); longer text is refused on the phone. Voice input is transcribed by the phone's speech recognition and sent
   through `useTextIntake()` (README D7, decided Q005). Each text intake counts as one Free scan (Q112).
5. Barcode lookup: digits only, 8–14 characters, or a GS1 2D payload; the hook rejects other input without a call.
   Results are cached 7 days (MOB-03-01); a `404` is cached 1 h so a repeated unknown scan doesn't call again.
6. Optimistic updates for add, edit and remove: `onMutate` cancels the stock query, applies the change, and keeps a
   snapshot; rollback only on a definitive server rejection (`validation`, `notFound` for edit). Queued changes stay
   visible through the overlay (MOB-04-04).
7. AI intake calls check `usePlusGate("photo_scan")` first (Free: 2 in any 7 days) and surface `quotaReached`.

## Contract

```ts
interface StockItemInput {
  id: string;                      // client-generated UUID (rule 3)
  ingredientId?: string; ean?: string; name: string;
  quantity: { value: number; unit: "gram" | "millilitre" | "piece" };
  storage: "fridge" | "freezer" | "pantry";
  expiresOn?: string;              // YYYY-MM-DD; server fills a default shelf life when absent
  source: "typed" | "barcode" | "unpack" | "photo" | "receipt" | "text" | "list_tick";
}
```

## Data

Query cache: stock persisted; barcode results persisted (last 500); proposals memory only.

## Errors

| Case | Behaviour |
|---|---|
| Edit of an item deleted on another phone | `notFound` → drop the queued edit, refetch stock (MOB-04-03) |
| Photo scan over quota | `quotaReached` with `nextFreeAt` |
| AI down | Intake returns `serverDown`/`aiUnavailable`; typing and barcode still work |

## Acceptance criteria

- [ ] Adding 3 items in airplane mode shows them at once; after reconnect `GET /v1/stock` returns them with the same
      ids.
- [ ] A shelf photo returns proposals; confirming 2 of 5 adds exactly 2 items.
- [ ] An unknown barcode named once is found by `useResolveBarcode` on a second account (verification 14).
