# PLN-14-07 Week plan with AI

| | |
|---|---|
| Service | planner |
| Epic | [PLN-14 Week plan](README.md) |
| Release | R3 |
| Size | L |
| Status | todo |
| Depends on | PLN-14-06, AI-15 (ComposeWeekPlan), PLN-07, PLN-14-01 |
| Unblocks | GW-07, PLN-20-01 |

## Goal

Compose a week plan with AI: the fewest cooking sessions, dinner leftovers reused, every person's target met, stock
first (features.md §2 with AI: "Week plan, including dinner leftovers reused in the next supper", "Week plan with
the fewest cooking sessions").

## Business rules

1. Plus only (`plus_required` otherwise).
2. Request to `AI-15 (ComposeWeekPlan)`: cooking days and allowance (PLN-14-02), per day and meal the eaters as
   anonymous keys with kcal budgets, shared limits per meal, stock with days to expiry (items past their date flagged
   `past_date`, Q103), tools, cuisines per meal,
   effort tags required for doubling, and up to 60 library candidates. No names, ages, weights or goal types.
3. AI returns per cooking day a dish (library id or new dish), which later meal is its leftover, and dishes for the
   other meals. AI may propose a leftover in a later supper ("dinner leftovers reused in the next supper") only for
   dishes tagged "good cold the next day" or "keeps 3 days".
4. Checks without AI, per dish: PLN-07 rule check; doubling only with the tags in PLN-14-03; session count ≤
   allowance + 1; stock not used twice; portions via PLN-12; then week balance (PLN-14-05).
5. Failing days are retried once with their violations; still failing days are filled by PLN-14-06 for those days
   only (`partial_fallback`).
6. `ai` unavailable, over budget or refusing: the whole week comes from PLN-14-06 (`ai_used = false`).
7. The AI week is compared with the non-AI week by the comparator (PLN-13-03); the AI week is dropped if worse on
   keys 1–3 (`ai_plan_rejected`).
8. Creating a week over existing day plans requires `replace = true`; cooked meals are never replaced.
9. Response visibility-filtered (PLN-12-07).

## Contract

`contracts/coolinary/planner/v1/week.proto`

```proto
rpc CreateWeekPlan(CreateWeekPlanRequest) returns (WeekPlan);
message CreateWeekPlanRequest {
  string household_id = 1;
  string viewer_member_id = 2;
  string start_date = 3;              // today … today + 13
  bool replace = 4;
  bool use_ai = 5;
  repeated int32 pinned_cooking_days = 6; // PLN-14-02
}
```

## Data

No new tables (PLN-14-01).

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Free household | `FAILED_PRECONDITION` `plus_required` |
| Overlap without `replace` | `ALREADY_EXISTS` |
| kitchen down | `UNAVAILABLE` |

## Acceptance criteria

- [ ] Verification 18 passes with AI on and with AI off.
- [ ] An AI week doubling a dish without "keeps 3 days" or "freezes well" is corrected before it is shown.
- [ ] The serialized AI request contains no names, ages or weights.
- [ ] Langfuse tag `week_plan`; the effort summary matches PLN-14-04 for the saved week.
