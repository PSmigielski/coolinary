# AI-01-06 Output validation, calorie guard and one repair call

| | |
|---|---|
| Service | ai |
| Epic | [AI-01 Claude provider foundation](README.md) |
| Release | R1 |
| Size | M |
| Depends on | AI-01-02 |
| Unblocks | every job epic |

## Goal

Nothing leaves the ai service unless it matches its schema and the job's own checks. A failed check gets one
repair attempt; then the caller gets `ai_invalid_output` and uses its fallback.

## Business rules

1. Two layers run on every completion:
   - **Schema**: the pydantic model (types, enums, required fields, string lengths, list sizes).
   - **Job validator**: a function `validate(request, parsed) -> list[Violation]` per job, defined in each job task
     (ids must come from the request, quantities in range, dates ISO, etc.).
2. **Calorie guard** (features.md rule 2, "AI never states a calorie number"), applied to every free-text field of
   every job output (dish names, steps, notes, reasons): reject any text matching a number followed by an energy
   unit — `kcal`, `cal`, `calories`, `kalorii`, `kalorie`, `kj`, `kilojoule`, `kilokalorie`, `Kalorien`, `kJ`,
   case-insensitive, with or without a space or `~`. Numbers given to the model in the request (e.g. library
   `kcal_per_portion` in R2) are never echoed in text; the guard applies to them too.
3. Violations are graded:
   - **Fixable by dropping**: an output list item that alone breaks a rule (an extracted item with an empty name,
     an offer with price 0) is dropped and counted in `ai_output_dropped_total{job, rule}`. If more than 30% of
     items would be dropped, treat it as not fixable.
   - **Not fixable**: anything else (unknown id in a composed dish, calorie text, schema error, truncation).
4. Not fixable → one repair call: the same request plus one extra user turn listing the violations by rule code
   and path (e.g. `meals[1].ingredients[3].ingredient_id: unknown id`), never quoting user text back. The repair
   call uses the same job spec; a truncated output repairs with `max_tokens` doubled (cap 64000).
5. If the repair also fails, return `UNAVAILABLE ai_invalid_output`. Both calls are billed and traced; the cost of
   both counts towards the spend limit.
6. The repair is skipped when less than `job.timeout_seconds / 2` of the gRPC deadline remains; the error is
   returned instead.
7. `ai_repairs_total{job, outcome="fixed"|"failed"}` is exported. A repair rate above 10% for a job over a day is a
   prompt bug (alert rule in Prometheus).

## Contract

No gRPC change.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Still invalid after repair | `UNAVAILABLE`, `ai_invalid_output` |

## Acceptance criteria

- [ ] Calorie guard test cases: `"about 520 kcal"`, `"520kcal"`, `"~500 Kalorien"`, `"2100 kJ"`, `"350 cal"` are rejected; `"bake 20 min at 180°C"`, `"2 cups"` pass.
- [ ] A fake provider returning one invalid then one valid output yields success with `ai_repairs_total{outcome="fixed"}` = 1 and two cost records.
- [ ] The repair turn contains rule codes and paths only (test asserts the user's original text is not in the repair message).
- [ ] Dropping 4 of 10 items fails as not fixable; dropping 2 of 10 succeeds with 8 items.
