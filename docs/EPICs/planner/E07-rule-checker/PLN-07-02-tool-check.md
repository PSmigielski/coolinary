# PLN-07-02 Tool check

| | |
|---|---|
| Service | planner |
| Epic | [PLN-07 Rule checker](README.md) |
| Release | R1 |
| Size | S |
| Depends on | PLN-07-01 |
| Unblocks | PLN-07-05, Verification 2 |

## Goal

No dish ever needs an unticked tool (features.md §1 "Kitchen tools: hard limit").

## Business rules

1. Required tools = the dish's listed tools ∪ the tool of every step.
2. Each required tool must be in the household's tool set (or the one-off tools of the request, intersected with the
   household's). Each missing tool → `TOOL_MISSING` with the tool code and the steps that use it.
3. Step texts are scanned, per language, for tool words of tools not ticked (a fixed word list per tool and language,
   e.g. oven: "oven", "piekarnik", "Backofen", "bake", "piec", "backen"); a hit → `TOOL_MISSING` with
   `found_in_text=true`. This catches AI steps that use an oven without listing it.
4. No substitutions are assumed (a hand blender does not stand in for a blender). Decided (Q105): no substitution
   table at launch.
5. The word list lives in `planner/internal/rules/toolwords_<lang>.txt`, reviewed with the language check.

## Contract

Internal.

## Data

No schema change.

## Errors

None of its own; internal check run by the checker core (PLN-07-01).

## Acceptance criteria

- [ ] Household with microwave + kettle; AI dish lists `[microwave]` but a step says "bake in the oven at 180 °C" →
      `TOOL_MISSING oven, found_in_text`.
- [ ] A step with tool `air_fryer` when only the oven is ticked fails.
- [ ] Polish step "upiecz w piekarniku" is caught.
