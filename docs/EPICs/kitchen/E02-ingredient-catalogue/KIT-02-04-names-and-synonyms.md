# KIT-02-04 Names, synonyms and plurals per language

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-02 Ingredient catalogue](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-02-01, KIT-03-01 (normaliser) |
| Unblocks | KIT-02-05, KIT-02-06, KIT-02-08, KIT-03-03 |

## Goal

Every ingredient has a display name in each launch language plus the other words people and shops use for it, so
autocomplete, matching and recipes work in Polish, German and English.

## Business rules

1. A name row has `language` (`pl`, `de`, `en`, later others), `name` (1–60 characters, as displayed), `kind`
   (`primary`, `synonym`, `plural`) and `origin` (`seed`, `admin`, `machine`).
2. Exactly one `primary` name per ingredient and language. Every active ingredient has a primary `en` name; the
   seed and the translation batch add `pl` and `de`.
3. `normalized` is the matching key from the normaliser (KIT-03-01), recomputed when the normaliser version changes.
4. The same normalised name may belong to two ingredients ("pepper": bell pepper, black pepper). That is allowed;
   matching then never takes the exact path for it (KIT-03-03).
5. One ingredient can't have the same normalised name twice in one language.
6. Display name for language L: L's primary, else `en` primary (with `translated = false` in responses).
7. `reviewed` is true for `seed` and `admin` names and false for `machine` names until an admin approves them
   (ADM-03). A `machine` name never replaces a `seed` or `admin` name.
8. Changes raise the catalogue version (KIT-02-06) so caches refresh.

## Contract

No RPC here; names are returned by KIT-02-05 and KIT-02-06.

## Data

`kitchen/migrations/000NN_ingredient_names.sql`

```sql
CREATE TABLE ingredient_names (
    id             bigserial PRIMARY KEY,
    ingredient_id  uuid NOT NULL REFERENCES ingredients (id) ON DELETE CASCADE,
    language       text NOT NULL CHECK (language ~ '^[a-z]{2}$'),
    name           text NOT NULL CHECK (char_length(name) BETWEEN 1 AND 60),
    kind           text NOT NULL CHECK (kind IN ('primary','synonym','plural')),
    normalized     text NOT NULL,
    origin         text NOT NULL CHECK (origin IN ('seed','admin','machine')),
    reviewed       boolean NOT NULL DEFAULT false,
    created_at     timestamptz NOT NULL DEFAULT now(),
    UNIQUE (ingredient_id, language, normalized)
);
CREATE UNIQUE INDEX ingredient_names_one_primary ON ingredient_names (ingredient_id, language) WHERE kind = 'primary';
CREATE INDEX ingredient_names_trgm ON ingredient_names USING gin (normalized gin_trgm_ops);
CREATE INDEX ingredient_names_lang ON ingredient_names (language, normalized);
```

## Errors

Validation errors come with the admin RPCs (KIT-02-10).

## Acceptance criteria

- [ ] Two primary names for one ingredient and language are refused.
- [ ] "pepper" may be a synonym of two ingredients.
- [ ] Display name for `de` falls back to `en` when no German primary exists, flagged `translated = false`.
- [ ] A normaliser version change recomputes `normalized` for all rows in one migration job.
