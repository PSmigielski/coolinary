# OFF-03 Postcode areas and flyer editions

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | Turn a postcode or town into an area, list the shops around it with their positions (the phone computes distances, README D13), and know which flyer edition each shop gets: national, regional, per store (Kaufland) or per store format (Carrefour). |

## Flow

1. A postcode dataset gives each postcode and town a centre point.
2. The phone sends a postcode (from its own reverse geocoding, or typed) or a town name; offers resolves it to an area.
3. offers returns every shop in a box of ±15 km around the area's centre, with coordinates; the phone sorts by
   distance and filters by the travel range.
4. Each shop maps to one or more flyer editions; offers of a shop are the offers of its editions.
5. Kaufland shops are linked to Kaufland store numbers, because their flyer depends on the store.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-03-01](OFF-03-01-postcode-dataset.md) | Postcode and town dataset import | R2 | M | OFF-01-05 |
| [OFF-03-02](OFF-03-02-resolve-area.md) | Resolve a postcode or town to an area | R2 | S | OFF-03-01 |
| [OFF-03-03](OFF-03-03-shops-near-area.md) | Shops near an area | R2 | M | OFF-03-02, OFF-02-03, OFF-02-04, OFF-03-04 |
| [OFF-03-04](OFF-03-04-editions-model.md) | Flyer editions and shop-to-edition rules | R2 | M | OFF-02-01, OFF-02-02 |
| [OFF-03-05](OFF-03-05-chain-store-numbers.md) | Link shops to chain store numbers | R2 | M | OFF-03-04, OFF-02-03 |
