# PLN-14 Week plan with leftovers and fewest cooking sessions

Release R3, Plus only. A 7-day plan that cooks as few times as the household's rhythm allows: on cooking days the
main dish is cooked in a double amount when it keeps or freezes, and the second half becomes a later meal, split
again for whoever eats then (features.md §2 "Calories and goals" with AI, "One dish per meal, few cooking sessions",
Verification 18).

## Flow

1. Cooking days from the rhythm (PLN-14-02).
2. Compose 7 days: with AI through `AI-15 (ComposeWeekPlan)` (PLN-14-07) or from the library (PLN-14-06), using the
   family-day building blocks (PLN-13) and cook-once-eat-twice (PLN-14-03).
3. Balance each person across days and the week (PLN-14-05), count sessions and stove time (PLN-14-04).
4. Store, read and edit the week (PLN-14-01).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-14-01](PLN-14-01-week-plan-model.md) | Week plan model, read and edit | R3 | M | PLN-08, PLN-12-06 |
| [PLN-14-02](PLN-14-02-cooking-days.md) | Cooking days from the rhythm | R3 | S | PLN-14-01, KIT-10 (GetPlanningProfile) |
| [PLN-14-03](PLN-14-03-cook-once-eat-twice.md) | Cook once, eat twice | R3 | M | PLN-14-02, PLN-12-08, PLN-02 |
| [PLN-14-04](PLN-14-04-effort-summary.md) | Effort summary (sessions, stove time) | R3 | S | PLN-14-03, PLN-13-02 |
| [PLN-14-05](PLN-14-05-week-balance.md) | Week balance | R3 | M | PLN-12-05, PLN-14-03 |
| [PLN-14-06](PLN-14-06-non-ai-week.md) | Week plan without AI | R3 | L | PLN-14-02 … PLN-14-05, PLN-13-06 |
| [PLN-14-07](PLN-14-07-ai-week.md) | Week plan with AI | R3 | L | PLN-14-06, AI-15 (ComposeWeekPlan), PLN-07 |
