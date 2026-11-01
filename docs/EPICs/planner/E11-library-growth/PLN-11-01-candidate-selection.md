# PLN-11-01 Promotion candidate job

| | |
|---|---|
| Service | planner |
| Epic | [PLN-11 Library growth from well-rated AI dishes](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | PLN-05-05, PLN-05-06, PLN-01-05 |
| Unblocks | PLN-11-02, PLN-11-03 |

## Goal

Find the household AI dishes worth adding to the library, without anyone reading households' data by hand.

## Business rules

1. Weekly job, Monday 04:00 UTC.
2. A household AI dish qualifies when: its household's rating (average of members, PLN-05-05 rule 4) ≥ 4.0; it was
   cooked ≥ 2 times by that household (PLN-05-06); it was created ≥ 14 days ago; it is not a rewrite or restyle of a
   library dish (those have a `parent_dish_id` in the library — the original is already there).
   Decided (Q214).
3. Qualifying dishes are clustered: same cuisine and ingredient-id set Jaccard ≥ 0.8 → one cluster. Cluster strength =
   number of distinct households, then average rating.
4. A cluster is skipped when a published or in-review library dish of the same cuisine has Jaccard ≥ 0.8 with the
   cluster's best dish (already in the library).
5. Each new cluster becomes a `promotion_candidates` row pointing at its best dish (highest rating, then most cooked).
   Existing open candidates are updated (strength), not duplicated.
6. Households of both plans count; the job reads only ids, ratings and counts.

## Contract

No RPC (job); results read through PLN-11-03.

## Data

`planner/migrations/000NN_promotion_candidates.sql`

```sql
CREATE TABLE promotion_candidates (
    id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    best_dish_id       uuid NOT NULL REFERENCES dishes (id) ON DELETE CASCADE,
    cuisine            text NOT NULL,
    households         int  NOT NULL,
    avg_rating         numeric(3,2) NOT NULL,
    status             text NOT NULL DEFAULT 'open' CHECK (status IN ('open','promoted','dismissed')),
    library_dish_id    uuid REFERENCES dishes (id),
    decided_by         text,
    decided_at         timestamptz,
    created_at         timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE promotion_cluster_members (
    candidate_id  uuid NOT NULL REFERENCES promotion_candidates (id) ON DELETE CASCADE,
    dish_id       uuid NOT NULL REFERENCES dishes (id) ON DELETE CASCADE,
    PRIMARY KEY (candidate_id, dish_id)
);
```

## Errors

| Case | Behaviour |
|---|---|
| Job fails halfway | the run is rolled back; the next scheduled run starts over (candidate selection is a pure read of ratings and cooks) |
| No candidates | the run ends with zero clusters; nothing is written |

## Acceptance criteria

- [ ] Two households' AI dishes with the same 9 of 10 ingredients and cuisine form one candidate with households = 2.
- [ ] A qualifying dish matching a published library dish (Jaccard ≥ 0.8) creates no candidate.
- [ ] Running the job twice in a week creates no duplicate candidates.
