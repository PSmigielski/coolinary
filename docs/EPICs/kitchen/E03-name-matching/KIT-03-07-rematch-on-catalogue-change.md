# KIT-03-07 Re-match after catalogue changes

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-03 Name matching](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | KIT-03-03, KIT-03-05, KIT-02-10 |
| Unblocks | — |

## Goal

When the catalogue grows or the normaliser changes, earlier "none" answers get another chance and keys stay valid.

## Business rules

1. When an ingredient or a name is added (admin, seed, translation batch), stored rows in that language with
   outcome `none`, method other than `admin`, and `similarity(key, new normalized name)` ≥ 0.5 are set to `queued`
   (picked up by the next batch, KIT-03-05).
2. Merges re-point matches through KIT-02-09; no extra work here.
3. A retired ingredient keeps its existing matches; it is no longer offered as a candidate.
4. On a normaliser version change a one-off job recomputes keys for all rows: rows whose new key is new are copied
   to the new version with their answer; when two old rows map to one new key, the `admin` answer wins, then the one
   with more hits. Old-version rows are deleted after the job.
5. The job runs in chunks of 5,000 rows and can resume after a restart.

## Contract

No RPC.

## Data

No schema change.

## Errors

None returned to users.

## Acceptance criteria

- [ ] Adding "pak choi" as a synonym re-queues the stored `none` row for "pak choy".
- [ ] A normaliser change that maps two old keys to one keeps the admin decision.
