# KIT-08-01 Shared intake pipeline: match, propose, metrics

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-08 AI intake](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-07-06 (proposal builder), KIT-03 (`MatchNames`), KIT-01 (ai client with deadlines) |
| Unblocks | KIT-08-02, KIT-08-03, KIT-08-04 |

## Goal

One pipeline turns AI-extracted items into ticked or unticked proposals with storage places and dates, the same way
for photos, text and receipts, and measures how often users accept them.

## Business rules

1. Input: extracted items (name, optional quantity and unit, confidence 0–1, optional `is_food`), the household,
   the language (`x-language`, else the profile language), an optional storage place hint, and the input kind.
2. Clean-up: names trimmed to 120 characters; empty names dropped; items with `is_food = false` dropped; at most
   60 items kept in the order AI returned them (`truncated = true` when more).
3. Items with the same normalised name are merged; quantities are added when the units match, otherwise the first
   quantity is kept.
4. All names go to `MatchNames` (KIT-03) in one call (≤ 100 names), context `MATCH_CONTEXT_PHOTO`,
   `MATCH_CONTEXT_TYPED` (text and voice) or `MATCH_CONTEXT_RECEIPT`.
5. Proposals come from the KIT-07-06 builder with the matched ingredient id (or the bare name when "none").
6. **Quantity**: the AI quantity converted to g, ml or pieces; when AI gave none, the ingredient's default pack size
   for the household's country (KIT-02-03) with note `estimated_quantity`, else 1 piece with the same note.
7. **Storage**: the storage hint from the request (where the photo was taken), when given, is the user's choice for
   every proposal; otherwise KIT-07-02 picks per item.
8. **Pre-tick**: `confidence ≥ 0.5` and matched → ticked; otherwise unticked with note `low_confidence` or
   `not_matched`. Decided (Q109): the pre-tick threshold is 0.5.
9. **Nothing is stored.** The response carries `intake_ref` (random 16-byte id, base32) and each proposal a
   `proposal_id` = `<intake_ref>:<index>`. Proposals have no server-side expiry: they are just pre-filled
   `NewStockItem`s. The app confirms with `AddStockItems(source, intake_ref)`.
10. `ai_used` is true whenever the AI call returned a result (even with zero items); false when a fallback was used.
    The gateway counts the Free quota only when it is true (GW-03-01).
11. **Metrics** (no names, no household ids): `kitchen_intake_proposals_total{kind, ticked}`,
    `kitchen_intake_confirmed_items_total{kind}` (counted in `AddStockItems` when `intake_ref` is set),
    `kitchen_intake_duration_seconds{kind}`, `kitchen_intake_ai_unavailable_total{kind}`.
12. User text, item names and images are never logged, traced or cached by kitchen.

## Contract

`contracts/coolinary/kitchen/v1/intake.proto` — shared response of the intake RPCs (KIT-08-02 … 04).

```proto
enum IntakeKind {
  INTAKE_KIND_UNSPECIFIED = 0;
  INTAKE_KIND_PHOTO = 1;
  INTAKE_KIND_TEXT = 2;
  INTAKE_KIND_VOICE = 3;
  INTAKE_KIND_RECEIPT = 4;
}

message IntakeResponse {
  string intake_ref = 1;
  repeated IntakeProposal proposals = 2;
  bool ai_used = 3;
  bool rules_fallback = 4;        // text only: the non-AI parser produced the proposals
  bool truncated = 5;
}

message IntakeProposal {
  string proposal_id = 1;         // "<intake_ref>:<index>"
  StockProposal proposal = 2;     // KIT-07-06
}
```

## Data

No schema change.

## Errors

| Case | gRPC status | Reason (error details) |
|---|---|---|
| ai unreachable, deadline passed, spend limit, refusal — and no fallback for the kind | `UNAVAILABLE` | `ai_unavailable` (gateway → 503, quota released) |
| `MatchNames` fails | — | proposals built unmatched (unticked); not an error |

## Acceptance criteria

- [ ] AI returns "Milch", "milch ", "Eier" (6), "Spülmittel" (`is_food = false`) → 2 proposals (milk once, eggs 6).
- [ ] An item with confidence 0.3 is unticked with `low_confidence`; a matched one with 0.8 is ticked.
- [ ] 75 extracted items → 60 proposals and `truncated = true`.
- [ ] Confirming 5 of 8 proposals via `AddStockItems(intake_ref)` adds 5 to `kitchen_intake_confirmed_items_total`.
- [ ] Log lines of an intake contain kind, counts and duration only.
