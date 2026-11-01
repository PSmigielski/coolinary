# OFF-05 Flyer source framework and sync scheduler

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | Check each covered chain's own flyer page on a schedule, notice new and changed flyers, download their pages once and politely, keep them only as long as allowed, and tell the founder and the users when a source breaks (features.md §3 "Keeping flyers in sync", steps 1–3, 8, 9). |

## Flow

1. Each chain has a source: list-page URL, adapter, schedule, allowed hosts, and a retailer-terms approval. No source
   runs without the approval.
2. A daily check (plus extra checks just after the chain's usual publishing time) asks the adapter for the flyers on
   the list page, per wanted edition.
3. New flyer ids create flyers; changed dates or page counts create a new version to be read again.
4. Pages of each new version are fetched once with the polite fetcher and stored until the flyer's end date + 90 days (Q010).
5. The stored version goes to extraction (OFF-07). Failures change the source status, alert the founder and show
   "flyer not updated, scan it yourself" to users of that chain.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-05-01](OFF-05-01-sources-and-flyers-schema.md) | Flyer sources, flyers and versions | R2 | M | OFF-02-01, OFF-03-04 |
| [OFF-05-02](OFF-05-02-adapter-interface.md) | Adapter interface and fixture harness | R2 | M | OFF-05-01 |
| [OFF-05-03](OFF-05-03-polite-fetcher.md) | Polite fetcher | R2 | M | — |
| [OFF-05-04](OFF-05-04-sync-and-detection.md) | Sync checks and change detection | R2 | L | OFF-05-02, OFF-05-03, OFF-01-05 |
| [OFF-05-05](OFF-05-05-page-storage.md) | Flyer page storage and retention | R2 | M | OFF-05-04 |
| [OFF-05-06](OFF-05-06-source-health-status.md) | Source health, alerts and flyer status for users | R2 | M | OFF-05-04, OFF-01-04 |
| [OFF-05-07](OFF-05-07-on-demand-editions.md) | Read store editions only on demand | R2 | S | OFF-05-04, OFF-04-03 |
| [OFF-05-08](OFF-05-08-sources-admin.md) | Sources and flyers admin RPCs | R2 | M | OFF-05-01, OFF-05-05, OFF-01-06 |
