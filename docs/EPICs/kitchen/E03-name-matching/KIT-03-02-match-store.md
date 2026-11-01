# KIT-03-02 Stored matches

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-03 Name matching](README.md) |
| Release | R1 |
| Size | S |
| Depends on | KIT-03-01, KIT-02-01 |
| Unblocks | KIT-03-03 … KIT-03-07, KIT-02-09 |

## Goal

Keep one decision per (language, key), so each name is matched once for everyone and the decisions can be reviewed.

## Business rules

1. One row per `(language, key, normalizer_version)`.
2. `outcome`: `matched` (has an ingredient), `none` (nothing fits), `queued` (waiting for a batch), `provisional`
   (trigram guess while AI was unavailable; retried).
3. `method`: `exact`, `synonym`, `trigram`, `ai`, `admin`. `confidence` 0–1.
4. `review`: `none`, `needed`, `done`. Rows decided by an admin have `method = 'admin'` and are never changed by any
   automatic path.
5. `contexts` collects where the name was seen: `product`, `receipt`, `offer`, `recipe`, `typed`, `photo`.
6. `example` keeps one raw name as seen, cut to 120 characters, for the reviewer.
7. `hit_count` and `last_seen_at` are updated at most once per minute per row (batched), to avoid write load.
8. On household deletion nothing here changes: the table holds no household id.

## Contract

No RPC.

## Data

`kitchen/migrations/000NN_name_matches.sql`

```sql
CREATE TABLE name_matches (
    id                  bigserial PRIMARY KEY,
    language            text NOT NULL,
    key                 text NOT NULL,
    normalizer_version  integer NOT NULL,
    ingredient_id       uuid REFERENCES ingredients (id),
    outcome             text NOT NULL CHECK (outcome IN ('matched','none','queued','provisional')),
    method              text NOT NULL CHECK (method IN ('exact','synonym','trigram','ai','admin')),
    confidence          real NOT NULL CHECK (confidence BETWEEN 0 AND 1),
    review              text NOT NULL DEFAULT 'none' CHECK (review IN ('none','needed','done')),
    contexts            text[] NOT NULL DEFAULT '{}',
    example             text NOT NULL,
    hit_count           integer NOT NULL DEFAULT 1,
    attempts            integer NOT NULL DEFAULT 0,
    first_seen_at       timestamptz NOT NULL DEFAULT now(),
    last_seen_at        timestamptz NOT NULL DEFAULT now(),
    decided_at          timestamptz,
    decided_by          text,
    UNIQUE (language, key, normalizer_version),
    CHECK ((outcome = 'matched') = (ingredient_id IS NOT NULL) OR outcome = 'provisional')
);
CREATE INDEX name_matches_review ON name_matches (hit_count DESC) WHERE review = 'needed';
CREATE INDEX name_matches_queued ON name_matches (hit_count DESC) WHERE outcome = 'queued';
CREATE INDEX name_matches_ingredient ON name_matches (ingredient_id);
```

## Errors

None (storage only).

## Acceptance criteria

- [ ] Inserting the same (language, key, version) twice fails.
- [ ] `ingredient_id` is registered in KIT-02-09's reference registry, so merges re-point matches.
- [ ] Hit counting for 1,000 lookups of one key in a minute writes the row at most twice.
