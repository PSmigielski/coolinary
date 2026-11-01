# OFF-02 Chains and shops from OpenStreetMap

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | Know every shop of the covered chains in Poland and Germany with its position, so the phone can show nearest shops first and the user can pick "my shops". Shop positions come from OpenStreetMap (open data, all EU countries); users can add or correct a shop. |

## Flow

1. A seed defines the chains (launch and next chains) and how to recognise them in OSM.
2. A monthly job imports each chain's shops per country from OSM and upserts them; shops missing twice are closed.
3. Nodes and building outlines of the same shop are merged.
4. A user can add a missing shop or propose a correction; it is visible to their household at once and to everyone
   after admin approval (ADM-06).
5. Every response with shop data carries the OSM attribution. The ODbL share-alike question is a gate before launch.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-02-01](OFF-02-01-chains-seed.md) | Chains table and seed | R2 | S | — |
| [OFF-02-02](OFF-02-02-shops-schema.md) | Shops schema rework | R2 | S | OFF-02-01 |
| [OFF-02-03](OFF-02-03-osm-import.md) | OSM shop import job | R2 | L | OFF-02-02, OFF-01-05 |
| [OFF-02-04](OFF-02-04-shop-dedupe.md) | Merge duplicate shops | R2 | S | OFF-02-03 |
| [OFF-02-05](OFF-02-05-user-adds-shop.md) | User adds a missing shop | R2 | M | OFF-02-02 |
| [OFF-02-06](OFF-02-06-user-corrects-shop.md) | User proposes a shop correction | R2 | S | OFF-02-02 |
| [OFF-02-07](OFF-02-07-shop-moderation-admin.md) | Shop moderation admin RPCs | R2 | S | OFF-02-05, OFF-02-06, OFF-01-06 |
| [OFF-02-08](OFF-02-08-osm-attribution-gate.md) | OSM attribution and ODbL gate | R2 | S | OFF-02-03 |
