# PLN-13-04 Eat lighter as a family

| | |
|---|---|
| Service | planner |
| Epic | [PLN-13 Family day plan](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | PLN-13-03, PLN-12-04, KIT-02 (GetIngredients) |
| Unblocks | PLN-13-06, PLN-13-07 |

## Goal

When the household turns on "eat lighter as a family", the whole family eats lighter together: dishes with more
vegetables and fewer kcal per portion, and lighter no-cook extras for everyone. Each adult keeps their own target;
the child is never singled out (features.md "Children and teenagers under 18").

## Business rules

1. The flag `eat_lighter_as_family` comes from the kitchen profile snapshot (KIT-11). The phone decides when to offer it (a child's
   goal is health data and stays on the phone); planner only reads the flag.
2. Vegetable share of a dish = grams of ingredients whose catalogue category is vegetables, fruit or legumes ÷ total
   grams, per standard portion (categories from KIT-02 `GetIngredients`; pieces and spoons converted with typical
   weights).
3. Lighter score for a candidate day = Σ over its dishes of `vegetable share band` (0–9, one per 10%) × 10 −
   `kcal per standard portion` ÷ 100. Used as key 4 of the comparator (PLN-13-03).
4. Extras use the lighter ranks (PLN-12-04 rule 2): fruit, raw vegetables, milk, yoghurt; no bread, cheese or nuts,
   for everyone in the household, not only the child.
5. Targets are unchanged; portions still follow each person's target.
6. No response or text key mentions why the flag is on, and nothing links it to a member.
7. Library browse (PLN-05) is not filtered by the flag; only composers use it.

## Contract

No RPC. Reads `eat_lighter_as_family` from `GetPlanningProfile`.

## Data

No schema change.

## Errors

None of its own; a preference read from `GetPlanningProfile` inside the family composers (PLN-13-06, PLN-13-07).

## Acceptance criteria

- [ ] With the flag on, a test day picks a vegetable-heavy curry over a cream pasta of equal stock score.
- [ ] Extras for Adam (an adult) also switch to fruit and vegetables (Verification 19: "changes the extras for
      everyone, not only the child").
- [ ] Turning the flag off restores the default extras order at the next recompute.
