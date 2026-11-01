# KIT-07-09 Link items typed without a catalogue match

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-07-03, KIT-03 (`MatchNames`) |
| Unblocks | — |

## Goal

Items the user typed as free text ("grandma's plum jam") get linked to a catalogue ingredient in the background, so
they start counting for recipes and the shopping list without the user doing anything.

## Business rules

1. A job runs every 10 minutes and takes up to 200 distinct (normalised name, household language) pairs from
   unmatched items, oldest first. Normalised = lower case, whitespace collapsed, trimmed.
2. It calls `MatchNames` (KIT-03) in one batch (up to 100 names per call, KIT-03 limit) with context
   `MATCH_CONTEXT_TYPED`.
3. A name matched with the confidence KIT-03 counts as "matched" links every unmatched item of that name in that
   language, in every household that has it: `ingredient_id` set, `matched = true`, the user's `name` kept,
   `version` + 1. Each touched household gets one `stock_version` bump and one `stock.changed`.
4. A name KIT-03 answers "none" for is tried again after 1 day, 7 days and 30 days, then never again (attempt count
   kept per name). The item stays unmatched; the app shows it as "not used for recipes" and the user can pick an
   ingredient with `UpdateStockItem`.
5. Names are user text: they are never logged or put in metrics or traces. KIT-03 sends them to AI only as the
   matching task requires.
6. If `ai` is down or over budget, KIT-03 still returns trigram-only results or "none"; the job does not fail, it
   retries on the next schedule step.

## Contract

No new RPC (internal job).

## Data

`kitchen/migrations/000NN_unmatched_names.sql`

```sql
-- +goose Up
CREATE TABLE unmatched_name_attempts (
    name_norm     text NOT NULL,
    language      text NOT NULL,
    attempts      smallint NOT NULL DEFAULT 0,
    next_try_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (name_norm, language)
);

-- +goose Down
DROP TABLE unmatched_name_attempts;
```

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | Behaviour |
|---|---|
| `MatchNames` fails | batch left for the next run; metric `kitchen_unmatched_link_failures_total` + 1 |

## Acceptance criteria

- [ ] An item typed as "Hähnchenschenkel" in a German household is linked to chicken thigh within one run and its
      name stays "Hähnchenschenkel".
- [ ] A name with no match is retried after 1, 7 and 30 days, then no more.
- [ ] Logs of a run contain counts only, never names.
- [ ] Two households with the same unmatched name are both linked by one `MatchNames` call.
