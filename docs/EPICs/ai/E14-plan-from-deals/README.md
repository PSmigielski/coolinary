# AI-14 Plan from deals

| | |
|---|---|
| Service | ai |
| Release | R2 |
| Goal | Compose dishes from the fridge plus this week's offers, preferring the nearest shops; deals only fill gaps left after using what is at home (features.md §2 "Flyers and offers", §8 "Flyva" risk; Plus). |

## Flow

```
planner PLN-18: stock + profile + offers for candidate ingredients at the household's shops (OFF-09),
                shops ranked nearest first by the phone (D13), no distances on the server
  -> ai.ComposeFromDeals (Sonnet 5.5 medium) -> meals with stock uses, items to buy and the offers used
  validation (AI-14-02), planner rule checker, list badges from the used offers
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-14-01](AI-14-01-compose-from-deals.md) | ComposeFromDeals | R2 | L | AI-10-03 |
| [AI-14-02](AI-14-02-deals-validation-and-golden-set.md) | Deals validation and golden set | R2 | M | AI-14-01, AI-03-05 |
