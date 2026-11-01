# AI-03-01 Prompt registry and versions

| | |
|---|---|
| Service | ai |
| Epic | [AI-03 Prompts, golden sets and evals](README.md) |
| Release | R1 |
| Size | S |
| Depends on | AI-01-02 |
| Unblocks | every job epic, AI-03-02 |

## Goal

Prompts are code: reviewed in pull requests, versioned, and traceable from every trace and cache entry back to the
exact text that produced it.

## Business rules

1. Layout: `ai/src/coolinary_ai/prompts/<job>/v<N>/system.md` (instructions, in English) and
   `prompts/<job>/schema.py` (pydantic output model, shared across versions until the schema changes; a schema
   change is a new version too). `prompts/<job>/__init__.py` exports `CURRENT = "v<N>"`.
2. The version string is `<job>/v<N>`. It goes into `CallInfo.prompt_version`, the cache key (AI-01-09) and the
   Langfuse tags (AI-02-03).
3. Any edit to `system.md` or the schema of a released version is forbidden; copy to `v<N+1>`, edit, eval, then move
   `CURRENT`. A unit test stores a hash of each released version and fails if it changes.
4. Prompts are written in English; the requested output language is given per request
   (`Write names and steps in: pl`). Output language codes are BCP 47 primary tags (`pl`, `de`, `en`).
5. Every system prompt includes three shared fragments from `prompts/_shared/`: the untrusted-input rule
   (AI-01-07), the "never state calorie or kilojoule numbers" rule, and the output-only-schema rule.
6. Prompts contain no personal data, no real household examples, no brand names of real users' products beyond
   generic examples.

## Contract

No gRPC change.

## Data

No schema change.

## Errors

None new.

## Acceptance criteria

- [ ] Changing a released `system.md` fails the hash test.
- [ ] `CallInfo.prompt_version` of a `text_to_items` call reads `text_to_items/v1`.
- [ ] Every job's rendered system prompt contains the three shared fragments (test).
