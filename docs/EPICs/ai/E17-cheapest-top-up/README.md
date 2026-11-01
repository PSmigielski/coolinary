# AI-17 Cheapest top-up

| | |
|---|---|
| Service | ai |
| Release | R3 |
| Goal | "Cheapest top-up": the smallest basket that completes the week plan using current offers, including swapping an ingredient for a discounted equivalent (features.md §2 "Shopping list", with AI; Plus). Totals and savings are computed by planner, never by AI. |

## Flow

```
planner PLN-20: week-plan lines still missing after stock, offers for them and for same-category alternatives (OFF-09),
                ranked shops, dish facts
  -> ai.CheapestTopUp (Sonnet 5.5 medium) -> per line: an offer, a swap to a discounted equivalent, or no offer
  validation (AI-17-02); planner prices both baskets and shows the AI one only if it is cheaper
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-17-01](AI-17-01-cheapest-top-up.md) | CheapestTopUp | R3 | L | AI-14-01 |
| [AI-17-02](AI-17-02-top-up-validation-and-golden-set.md) | Top-up validation and golden set | R3 | M | AI-17-01, AI-03-05 |
