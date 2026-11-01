# PLN-08 AI day plan and cuisine choice

| | |
|---|---|
| Service | planner |
| Release | R1 |
| Goal | The core of the product: "one look in your fridge, three meals sorted". Compose breakfast, dinner and supper from current stock, oldest items first, so the three meals don't compete for the same ingredients, in the cuisine chosen per meal — with AI when allowed and available, checked by rules, and from the library when not. |

## Flow

```
ComposeDayPlan(date, meals, ai_allowed)                                   (gateway sets ai_allowed from the Free quota)
  ├─ ai_allowed and ai up ─> assemble request (PLN-08-02) ─> AI-10 ComposeDayPlan (Haiku 5.5 on Free; Plus per AI-10-05)
  │     ─> CandidateDish per meal ─> rule checker + plan allocation (PLN-07, PLN-08-03)
  │     ─> failures? one repair call for the failed meals only (PLN-08-04)
  │     ─> still failing meal ─> library dish for that meal
  │     ─> stock doesn't fit a meal's cuisine ─> three choices (PLN-08-06)
  └─ otherwise ─> library day plan (PLN-08-05)
  ─> plan stored as `proposed` (PLN-08-01), AI dishes stored as household AI dishes with calories from tables
  ─> response says ai_used, so the gateway counts quota only for AI plans
```

By hand (Q117): `SetPlanMeal` puts a dish the user picked into a plan meal: rule-checked, no AI, no quota (PLN-08-08).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-08-01](PLN-08-01-plan-storage.md) | Plan storage, accept and read | R1 | M | PLN-02-01 |
| [PLN-08-02](PLN-08-02-compose-request.md) | Compose request assembly | R1 | M | PLN-01-03, PLN-01-04, PLN-06-02, AI-10 (ComposeDayPlan) |
| [PLN-08-03](PLN-08-03-plan-allocation.md) | Ingredient allocation across meals | R1 | M | PLN-06-01 |
| [PLN-08-04](PLN-08-04-validate-repair-store.md) | Validate, repair once, store AI dishes | R1 | L | PLN-07-05, PLN-08-03, PLN-04-01, PLN-03-02 |
| [PLN-08-05](PLN-08-05-library-fallback.md) | Library day plan (non-AI fallback) | R1 | M | PLN-06-02, PLN-08-03 |
| [PLN-08-06](PLN-08-06-cuisine-choices.md) | Cuisine mismatch: three choices | R1 | M | PLN-08-04, PLN-08-05 |
| [PLN-08-07](PLN-08-07-compose-day-plan-rpc.md) | ComposeDayPlan RPC, quota signal, idempotency | R1 | M | PLN-08-01 … PLN-08-06 |
| [PLN-08-08](PLN-08-08-set-plan-meal.md) | SetPlanMeal: put a chosen dish into a plan meal | R1 | M | PLN-08-01, PLN-08-03, PLN-07-01 … PLN-07-04 |
