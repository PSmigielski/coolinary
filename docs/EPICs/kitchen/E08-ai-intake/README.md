# KIT-08 AI intake: photo, text, receipt

| | |
|---|---|
| Service | kitchen |
| Release | R1 (photo); R3 (text/voice, receipt) |
| Depends on | KIT-07-06 (proposals), KIT-03 (`MatchNames`), AI-06 (`ExtractItemsFromPhoto`), AI-05 (`ExtractItemsFromText`), AI-16 (`ExtractReceipt`) |

## Goal

A photo of a shelf, a sentence ("bought eggs, two peppers, some cheese") or a receipt becomes a list of items the
user confirms. AI only proposes; nothing reaches the fridge until the user confirms, and when AI is down the user can
still type or scan (features.md §2 "Fridge and pantry, with AI", rule 1 "the app stays useful with AI switched off").

## Flow

```
gateway (quota reserved) ─► IntakeFromPhoto / IntakeFromText / IntakeFromReceipt
                              │
                              ├─► ai: extract items (Haiku 5.5)          ── down/over budget ─► UNAVAILABLE (photo, receipt)
                              │                                            └─► rule-based parser (text)
                              ├─► MatchNames (KIT-03) / receipt-name memory (KIT-08-05)
                              └─► proposal builder (KIT-07-06): storage, use-by date, ticks
                              ▼
                       proposals + intake_ref + ai_used  ──► app ──► user confirms ──► AddStockItems (KIT-07-03)
                                                                                           │ receipt only
                                                                                           ▼
                                                                   receipt-name memory + receipt.scanned (offers)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-08-01](KIT-08-01-intake-pipeline.md) | Shared intake pipeline: match, propose, metrics | R1 | M | KIT-07-06, KIT-03 |
| [KIT-08-02](KIT-08-02-photo-intake.md) | Photo intake (fridge or shelf) | R1 | M | KIT-08-01, AI-06 |
| [KIT-08-03](KIT-08-03-text-intake.md) | Text and voice intake with a non-AI fallback | R3 | M | KIT-08-01, AI-05, KIT-02 |
| [KIT-08-04](KIT-08-04-receipt-intake.md) | Receipt intake and line matching | R3 | L | KIT-08-01, AI-16, KIT-04 |
| [KIT-08-05](KIT-08-05-receipt-name-memory.md) | Receipt-name memory per household and chain | R3 | M | KIT-08-04 |
| [KIT-08-06](KIT-08-06-receipt-scanned-event.md) | receipt.scanned event for shop prices | R3 | M | KIT-08-04, KIT-01 |

## Decisions (taken 1 November 2026)

- Q109 (KIT-08-01): Matched items with confidence ≥ 0.5 are pre-ticked.
- Q112 (KIT-08-03): Yes, every AI reading request counts as one scan.
- Q219 (KIT-08-03): Release 3, together with voice, as features.md groups them.
- Q221 (KIT-08-04): Yes: a 2-hour copy in the cache; no image kept.
- Q222 (KIT-08-06): Yes: no household, user or receipt id in the event.
