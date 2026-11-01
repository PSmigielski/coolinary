# AI-15 Week plan

| | |
|---|---|
| Service | ai |
| Release | R3 |
| Goal | A week plan with the fewest cooking sessions: cook a double amount on cooking days when the dish keeps 3 days or freezes, and plan the second half for a later meal; leftovers of dinner reused at a later supper (features.md §2 "Dishes and the day plan", "Calories and goals", "One dish per meal, few cooking sessions"; Plus). |

## Flow

```
planner PLN-14: rhythm (every day / every other day / twice a week), 7 dates, meals per day, stock, library candidates
                with effort tags, eaters (family) when goals are set
  -> ai.ComposeWeekPlan (Sonnet 5.5 high, streamed)
     <- per day and meal: COOK (dish, double batch?) | LEFTOVER of an earlier COOK | NO_COOK dish
  validation (AI-15-03); planner splits portions per eater, counts sessions and time, builds the list
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [AI-15-01](AI-15-01-week-plan-contract.md) | ComposeWeekPlan contract | R3 | S | AI-12-01 |
| [AI-15-02](AI-15-02-week-plan-compose.md) | Week composition prompt and call | R3 | L | AI-15-01 |
| [AI-15-03](AI-15-03-week-plan-validation.md) | Validation of sessions and leftovers | R3 | M | AI-15-02, AI-12-03 |
| [AI-15-04](AI-15-04-week-plan-golden-set.md) | Golden set: family weeks | R3 | M | AI-15-03, AI-03-05 |
