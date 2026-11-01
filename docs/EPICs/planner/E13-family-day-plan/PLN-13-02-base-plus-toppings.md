# PLN-13-02 Base-plus-toppings plates

| | |
|---|---|
| Service | planner |
| Epic | [PLN-13 Family day plan](README.md) |
| Release | R2 |
| Size | M |
| Depends on | PLN-13-01, PLN-02, PLN-12-08 |
| Unblocks | PLN-13-06, PLN-13-07, PLN-14-04 |

## Goal

When one person's diet narrows the choice, use dishes with one base for all and toppings chosen per plate (tacos,
bowls, wraps), so everyone still shares one dish (features.md "One dish per meal").

## Business rules

1. A dish is base-plus-toppings when PLN-02 gives it the effort tag `base_plus_toppings`: ingredient lines with role
   `base` and lines with role `topping` in named `topping_group`s (e.g. "protein", "sauce").
2. The base must pass the combined limits of all eaters (PLN-13-01).
3. Each eater gets every topping group that passes their own limits; a dish is usable only if every eater gets at
   least one topping group, or the dish marks toppings as optional.
4. Topping amounts scale with the sum of portions of the eaters who take that topping (PLN-12-08 rule 4). The
   base scales with the whole pot.
5. The kcal of a person's plate = base kcal per portion × portions + the kcal of their toppings per portion ×
   portions. PLN-12-03 uses this plate kcal instead of the dish's standard portion kcal for that eater.
6. A topping group that needs its own pan (a PLN-02 flag on the topping group, see Notes) counts as one extra cooking session (PLN-14-04) when at least
   one eater takes it.
7. When PLN-13-01 reports a narrowing eater, base-plus-toppings dishes get preference in the comparator
   (PLN-13-03, tie-break after the hard limits) for that meal.
8. Output per plate: member id, topping groups taken, portions; the dish card shows "Ewa: beans, cheese".

## Contract

`contracts/coolinary/planner/v1/family.proto`

```proto
message PlateToppings { string member_id = 1; repeated string topping_groups = 2; }
// PLN-08's plan meal message carries: repeated PlateToppings plates
```

## Data

No new table: plates are stored in PLN-08's `plan_meals.plates` (jsonb), as
`[{"member_key": "<member id>", "topping_groups": ["beans", "cheese"]}]`.

## Errors

| Case | Behaviour |
|---|---|
| An eater gets no topping and toppings are required | dish rejected for this meal |

## Acceptance criteria

- [ ] Tacos for Adam, Ewa (vegetarian), Ola, Kuba: Ewa's plate has beans and cheese, no chicken; chicken amount
      scales with 4½ portions (Adam 1½, Ola 1, Kuba 2), not with the pot of 6.
- [ ] A topping that needs its own pan adds one session to the day's count.
- [ ] Ewa's plate kcal is computed from her toppings, so her portion may differ from a meat eater's at the same target.

## Notes

- Needs a per-topping-group flag "needs its own pan" in PLN-02 (not in PLN-02-01's schema yet). Until it exists,
  no topping counts as extra cooking.
