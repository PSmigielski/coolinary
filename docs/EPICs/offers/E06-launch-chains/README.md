# OFF-06 Launch chain adapters

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | Read the own flyer pages of the seven launch chains (features.md §3 "Sources per shop: launch chains", §7 release 2): Biedronka, Lidl and Kaufland in Poland; Lidl, Aldi Süd, Aldi Nord and Kaufland in Germany. |

## Flow per chain

1. `GATE:` the chain's terms of use are read and approved in admin (OFF-05-08) before the source is enabled
   (features.md §3: "Each retailer's terms must be checked before step 2 goes live for that chain").
2. The adapter is built against recorded pages (fixtures) and passes the shared contract test (OFF-05-02).
3. The source is enabled with the chain in `syncing`: flyers are detected, fetched and extracted, but offers are not
   shown to users.
4. A two-week shadow run checks detection before start dates (Verification 7).
5. The chain goes `live` only through the accuracy gate (OFF-08-02, ≥ 97% on 100 offers).

Each adapter starts with the extra check slots below. After the shadow run, each slot is moved to one hour after the
earliest publishing time actually observed for that chain.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-06-01](OFF-06-01-biedronka-pl.md) | Biedronka (PL) adapter | R2 | M | OFF-05-02, OFF-05-03 |
| [OFF-06-02](OFF-06-02-lidl-pl.md) | Lidl PL adapter | R2 | M | OFF-05-02, OFF-05-03 |
| [OFF-06-03](OFF-06-03-kaufland-pl.md) | Kaufland PL adapter (per store) | R2 | L | OFF-05-02, OFF-05-03, OFF-03-05, OFF-05-07 |
| [OFF-06-04](OFF-06-04-lidl-de.md) | Lidl DE adapter | R2 | M | OFF-05-02, OFF-05-03, OFF-06-02 |
| [OFF-06-05](OFF-06-05-aldi-sued-de.md) | Aldi Süd adapter | R2 | M | OFF-05-02, OFF-05-03 |
| [OFF-06-06](OFF-06-06-aldi-nord-de.md) | Aldi Nord adapter (PDF) | R2 | M | OFF-05-02, OFF-05-03 |
| [OFF-06-07](OFF-06-07-kaufland-de.md) | Kaufland DE adapter (per store) | R2 | S | OFF-06-03 |
