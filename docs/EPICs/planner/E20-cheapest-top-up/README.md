# PLN-20 Cheapest top-up

Release R3, Plus only. "Cheapest top-up": the smallest basket that completes the week plan using current offers,
including swapping an ingredient for a discounted equivalent (features.md §2 "Shopping list" with AI). AI proposes
swaps; rules check every swapped dish; prices come from offers, never from AI.

## Flow

1. Work out what the week still needs and which swaps are allowed (PLN-20-01).
2. Ask `AI-17 (CheapestTopUp)` for a basket with swaps (PLN-20-02), or build a basket without swaps (PLN-20-03).
3. The user accepts some or all swaps; plan and list change (PLN-20-04).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-20-01](PLN-20-01-needs-and-allowed-swaps.md) | Top-up needs and allowed swaps | R3 | M | PLN-15-02, PLN-15-03, PLN-18-01, PLN-07 |
| [PLN-20-02](PLN-20-02-ai-cheapest-top-up.md) | Cheapest top-up with AI | R3 | L | PLN-20-01, AI-17 (CheapestTopUp), PLN-19-02 |
| [PLN-20-03](PLN-20-03-non-ai-top-up.md) | Top-up without AI | R3 | S | PLN-20-01, PLN-19-02 |
| [PLN-20-04](PLN-20-04-apply-top-up.md) | Applying a top-up | R3 | M | PLN-20-02, PLN-14-01, PLN-15-10, PLN-12-06 |
