# PLN-14-04 Effort summary (sessions, stove time)

| | |
|---|---|
| Service | planner |
| Epic | [PLN-14 Week plan](README.md) |
| Release | R3 |
| Size | S |
| Status | todo |
| Depends on | PLN-14-03, PLN-13-02 |
| Unblocks | PLN-14-06, PLN-14-07 |

## Goal

The week plan shows the effort: "3 cooking sessions this week, about 2 hours at the stove" (features.md "One dish
per meal").

## Business rules

1. A cooking session is one cooked dish for one meal on one day: every `cooking_sessions` row, plus every cooked
   (not "no cooking") breakfast or other meal, plus every topping group that needs its own pan (PLN-13-02), plus a
   cooked second dish (PLN-13-05).
2. Leftover meals and "no cooking" dishes count 0.
3. Stove time per session = the dish's active minutes (PLN-02) × 1.25 when the amount factor is 2, else × 1; a
   separate topping pan adds 10 minutes.
4. Week stove time = Σ sessions, rounded to the nearest 15 minutes for display ("about 2 hours"); the exact minutes
   are also returned.
5. The summary is recomputed on every week edit.

## Contract

```proto
message EffortSummary { int32 cooking_sessions = 1; int32 stove_minutes = 2; int32 stove_minutes_rounded = 3; }
```

## Data

No schema change.

## Errors

None of its own: the summary is computed inside the week-plan RPCs (PLN-14-01, PLN-14-07), whose errors apply. A dish without active minutes counts 0 minutes and sets `stove_minutes_incomplete = true` rather than failing.

## Acceptance criteria

- [ ] A week with 3 sessions of 30, 40 and 25 active minutes, the first doubled: 37.5 + 40 + 25 = 102.5 → shown as
      105 minutes, "about 1¾ hours".
- [ ] A taco night with a separate meat pan counts 2 sessions.

## Notes

- Decided (Q218): a double amount takes ×1.25 the time; +10 minutes per topping that needs its own pan. Replace with
  per-dish values if PLN-02 stores them later.
