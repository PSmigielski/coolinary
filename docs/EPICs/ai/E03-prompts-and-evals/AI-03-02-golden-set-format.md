# AI-03-02 Golden set format, storage and privacy

| | |
|---|---|
| Service | ai |
| Epic | [AI-03 Prompts, golden sets and evals](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-03-01 |
| Unblocks | AI-03-03, every golden-set task |

## Goal

One format for test cases of every job, stored in the repo, built only from data the team is allowed to keep.

## Business rules

1. Location: `ai/evals/<job>/cases/<case-id>.json`, binary inputs next to it (`<case-id>.1.jpg`, `.pdf`). Case ids
   are short slugs (`pl-fridge-07`).
2. Case file fields: `id`, `job`, `language`, `request` (the gRPC request as JSON, binary inputs referenced by
   file name), `expected` (job-specific: expected items, expected ids, forbidden ids, rules to check), `tags`
   (e.g. `pl`, `stt`, `injection`, `hard`).
3. Privacy rules for golden data:
   - Photos are taken by the team in their own kitchens; no people, faces, hands with rings, addresses, mail,
     screens or documents in frame. Each photo is checked before commit.
   - Receipts are real but redacted: card numbers, loyalty numbers, cashier names and store tax ids blacked out.
   - Flyer pages from chains are stored only under `ai/evals/` for testing, not republished anywhere, and removed
     if a retailer objects.
   - No production data ever becomes a golden case.
4. Binary inputs over 1 MiB are resized before commit (long edge 2048 px, JPEG quality 80); total size per job
   under 60 MiB. If the repo grows past 300 MiB of eval data, move binaries to Git LFS.
5. Each job's golden-set task states its size and the share per language (at least PL, DE, EN).

## Contract

No gRPC change.

## Data

No schema change. New folder `ai/evals/`.

## Errors

None at runtime; this task is eval or tooling work. A failed threshold blocks the release of the job (AI-03-04).

## Acceptance criteria

- [ ] A JSON schema `ai/evals/case.schema.json` validates every case file; a unit test runs it over all cases.
- [ ] `CONTRIBUTING` note in `ai/evals/README.md` lists the privacy rules above.
