# AI-06 Photo to items

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | A photo of the fridge, a shelf, the freezer or an unpacked shopping bag becomes an item list the user confirms (features.md §2, with AI, per request; Free: 2 photo scans a week, counted by the gateway). |

## Flow

```
phone resizes photos (long edge 2048 px, JPEG) -> gateway (quota GW-03) -> kitchen (KIT-08)
  -> ai.ExtractItemsFromPhoto (Haiku 5.5, effort medium), photos processed in memory, never stored (D6)
  <- items {name, brand, estimated quantity, packaging, hints, confidence, image_index}
kitchen matches to the catalogue (KIT-03), suggests storage and use-by, returns proposed items to confirm
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-06-01](AI-06-01-photo-to-items.md) | ExtractItemsFromPhoto | R1 | L | AI-05-01, AI-01-08 |
| [AI-06-02](AI-06-02-photo-golden-set.md) | Golden set: 60 fridge and shelf photos | R1 | M | AI-06-01, AI-03-03 |
