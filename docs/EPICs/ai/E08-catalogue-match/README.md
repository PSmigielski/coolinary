# AI-08 Catalogue match pick

| | |
|---|---|
| Service | ai |
| Release | R1 |
| Goal | For a product name, receipt line or flyer offer, pick the one catalogue ingredient it is (or say "none") from the trigram candidates kitchen found (decision 13). Each name is matched once for everyone and stored by kitchen. |

## Flow

```
kitchen MatchNames (KIT-03): normalise name -> stored match? -> pg_trgm top 10 candidates per language
   -> ai.PickCatalogueMatches (sync, up to 50 names)         for live scans, intake and imports
   -> or SubmitBatch(pick_catalogue_matches)                  for bulk work (offers OFF-07, catalogue imports)
   <- per name: ingredient_id or none, confidence, reason
kitchen applies its thresholds (auto-accept / review queue / none) and stores the match
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-08-01](AI-08-01-pick-catalogue-matches.md) | PickCatalogueMatches (sync) | R1 | M | AI-01-03, AI-01-06, AI-01-09 |
| [AI-08-02](AI-08-02-batch-pick-matches.md) | Batch kind pick_catalogue_matches | R1 | S | AI-08-01, AI-04-03 |
| [AI-08-03](AI-08-03-match-golden-set.md) | Golden set: 500 names | R1 | M | AI-08-01, AI-03-03 |
