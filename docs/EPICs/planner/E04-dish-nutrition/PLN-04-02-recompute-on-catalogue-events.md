# PLN-04-02 Recompute on catalogue events

| | |
|---|---|
| Service | planner |
| Epic | [PLN-04 Dish nutrition and calories](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-04-01, PLN-01-01, KIT-06 (catalogue.nutrition_changed), KIT-02 (catalogue.ingredient_merged) |
| Unblocks | — |

## Goal

When kitchen fixes a nutrition match or merges two ingredients, every dish using them shows the corrected values within
minutes, without anyone re-saving dishes.

## Business rules

1. `catalogue.nutrition_changed {ingredient_ids}`: enqueue every dish version that uses any of the ids and is
   referenced by something live — the current version of published dishes, open draft versions, and household AI
   dishes of the last 90 days — into `nutrition_recompute_queue` (dish_id, version), de-duplicated.
2. `catalogue.ingredient_merged {from, into}`: in one transaction rewrite `dish_ingredients.ingredient_id` from →
   into for all versions (an id rename, so no new version), rewrite `recipe_name_matches` (PLN-03-02), then enqueue
   the affected versions as in rule 1 and recompute their derived flags (PLN-02-03). If a version then has the same
   ingredient twice with the same role and group, the amounts are summed into the first line (same unit) or both lines
   are kept (different units).
3. A worker drains the queue every 30 s, 200 versions per run, 1 `ComputeNutrition` call per version; on kitchen
   `Unavailable` it stops the run and retries at the next one.
4. Old versions referenced only by history or old plans are not recomputed: history keeps the values it showed.
   Decided (Q095).
5. Household AI dishes stored with `pending=true` (PLN-04-01 rule 8) are also queued at insert.
6. A metric `planner_nutrition_queue_depth` and an alert threshold of 5,000 rows for 1 hour (alert via Prometheus
   rule, not `SendAlert`).

## Contract

No RPC.

## Data

`planner/migrations/000NN_nutrition_queue.sql`

```sql
CREATE TABLE nutrition_recompute_queue (
    dish_id     uuid NOT NULL,
    version     int  NOT NULL,
    reason      text NOT NULL,          -- 'nutrition_changed' | 'ingredient_merged' | 'pending'
    enqueued_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (dish_id, version)
);
```

## Events

Consumes `catalogue.nutrition_changed` and `catalogue.ingredient_merged` (handlers enqueue only; the worker calls
kitchen outside the event transaction, PLN-01-01 rule 9).

## Errors

| Case | Behaviour |
|---|---|
| kitchen `ComputeNutrition` returns `UNAVAILABLE` | the run stops and retries at the next run; already updated dishes keep their new values |
| A handler error on the event | rolled back and retried by the consumer framework; after 5 failures the event goes to `dead_events` (PLN-01-01) |

## Acceptance criteria

- [ ] Changing the nutrition match of "basmati rice" in kitchen updates the kcal of every published dish with rice
      within 2 minutes in the dev stack.
- [ ] Merging "tomato passata" into "passata" leaves no `dish_ingredients` row with the old id and creates no new
      dish version.
- [ ] Replaying the same event twice enqueues each version once.
