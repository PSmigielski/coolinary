# AI-05 Text to items

| | |
|---|---|
| Service | ai |
| Release | R1–R3 |
| Goal | A typed or spoken sentence ("bought eggs, two peppers, some cheese") becomes a list of items with quantities that the user confirms (features.md §2 "Fridge and pantry", with AI, per request). |

## Flow

```
phone: typed text, or speech recognition on the phone (D7)
  -> gateway -> kitchen (KIT-08) -> ai.ExtractItemsFromText (Haiku 5.5, effort low)
  <- items {name, quantity, unit, hints, confidence}
kitchen matches names to the catalogue (KIT-03) and returns proposed items; nothing is stored until the user confirms
```

Decided (Q219): typed and voice intake reach users in release 3. AI-05-01 is still built in R1 because photo intake
(AI-06-01) reuses its item rules.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-05-01](AI-05-01-text-to-items-claude.md) | ExtractItemsFromText on Claude | R1 | M | AI-01-03, AI-01-06, AI-01-07, AI-01-09, AI-03-01 |
| [AI-05-02](AI-05-02-text-golden-set.md) | Golden set: typed and spoken item lists | R3 | S | AI-05-01, AI-03-03 |
