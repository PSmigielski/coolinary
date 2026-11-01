# PLN-15-01 List model and GetList

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | M |
| Depends on | PLN-01, KIT-02 (GetIngredients) |
| Unblocks | PLN-15-02 … PLN-15-10, PLN-16-03, PLN-17-01, PLN-19-03, GW-08 |

## Goal

One shopping list per household with lines that know what they are, how much, where they came from, and whether they
were bought.

## Business rules

1. Each household has exactly one active list, created on first use. The list has a `version` (starts at 0, +1 per
   committed change, PLN-15-08).
2. A line is either an ingredient line (catalogue `ingredient_id`) or a free-text line (≤ 100 characters,
   trimmed, never logged).
3. Amount: `amount` in the ingredient's base unit (g, ml or piece, from `GetIngredients`), plus the pack
   rounding result `pack_size` × `pack_count` when known (PLN-15-04). `manual_amount = true` when a person typed the
   amount.
4. Sources: one or more per line, each with kind `plan_meal`, `extra`, `base_pantry`, `manual`, `offer`, `deal` or
   `topup`, its reference (plan meal id, cuisine code, member id, offer id, …) and its amount contribution. The
   line's amount is the sum of contributions unless `manual_amount`.
5. States: active (not ticked), ticked (with `ticked_at` and `ticked_by` member id), removed (tombstone kept 7 days
   for sync).
6. At most one active, unticked ingredient line per ingredient; adding the same ingredient merges into it.
7. At most 300 active lines; more returns `RESOURCE_EXHAUSTED`.
8. `GetList` returns active and ticked lines (ticked first hidden by the app), each with the ingredient's name in
   `x-language`, its catalogue category (for grouping by shop section), amount, packs, sources, state, `added_by`
   member id, line version, and badges when present (PLN-17). Names and categories come from `GetIngredients`,
   cached in redis-cache for 24 hours per language.
9. Any account member of the household can read and edit the list (features.md "shared live between household
   members"; sharing itself needs Plus in identity).

## Contract

`contracts/coolinary/planner/v1/list.proto`

```proto
rpc GetList(GetListRequest) returns (ShoppingList);

message GetListRequest { string household_id = 1; bool include_ticked = 2; }
message ShoppingList { string id = 1; int64 version = 2; repeated ListLine lines = 3; }
message ListLine {
  string id = 1;
  string ingredient_id = 2;               // empty for free text
  string name = 3;                         // ingredient name in x-language, or the free text
  string category = 4;                     // catalogue category, "" for free text
  coolinary.common.v1.Quantity amount = 5;
  coolinary.common.v1.Quantity pack_size = 6;
  int32 pack_count = 7;
  bool manual_amount = 8;
  repeated LineSource sources = 9;
  LineState state = 10;
  string ticked_at = 11; string ticked_by_member_id = 12;
  string added_by_member_id = 13;
  int64 version = 14;
  repeated OfferBadge badges = 15;         // PLN-17-01
  repeated string flags = 16;              // "check_at_home", "plan_changed", "optional"
}
message LineSource { SourceKind kind = 1; string ref = 2; coolinary.common.v1.Quantity amount = 3; string member_id = 4; string date = 5; Meal meal = 6; }
enum SourceKind { SOURCE_KIND_UNSPECIFIED = 0; SOURCE_KIND_PLAN_MEAL = 1; SOURCE_KIND_EXTRA = 2; SOURCE_KIND_BASE_PANTRY = 3; SOURCE_KIND_MANUAL = 4; SOURCE_KIND_OFFER = 5; SOURCE_KIND_DEAL = 6; SOURCE_KIND_TOPUP = 7; }
enum LineState { LINE_STATE_UNSPECIFIED = 0; LINE_STATE_ACTIVE = 1; LINE_STATE_TICKED = 2; LINE_STATE_REMOVED = 3; }
```

## Data

`planner/migrations/000NN_shopping_list.sql`

```sql
CREATE TABLE shopping_lists (
    id           uuid PRIMARY KEY,
    household_id uuid NOT NULL UNIQUE,
    version      bigint NOT NULL DEFAULT 0,
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE list_lines (
    id                 uuid PRIMARY KEY,
    list_id            uuid NOT NULL REFERENCES shopping_lists ON DELETE CASCADE,
    ingredient_id      text,
    free_text          text CHECK (free_text IS NULL OR length(free_text) <= 100),
    normalized_text    text,
    amount_value       double precision,
    amount_unit        smallint NOT NULL DEFAULT 0,
    pack_value         double precision,
    pack_unit          smallint NOT NULL DEFAULT 0,
    pack_count         integer,
    manual_amount      boolean NOT NULL DEFAULT false,
    state              smallint NOT NULL DEFAULT 1,
    ticked_at          timestamptz,
    ticked_by          uuid,
    stock_item_ids     text[] NOT NULL DEFAULT '{}',
    added_by           uuid,
    flags              text[] NOT NULL DEFAULT '{}',
    version            bigint NOT NULL,
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now(),
    removed_at         timestamptz,
    CHECK ((ingredient_id IS NULL) <> (free_text IS NULL))
);
CREATE UNIQUE INDEX list_lines_one_open_per_ingredient
    ON list_lines (list_id, ingredient_id) WHERE state = 1 AND ingredient_id IS NOT NULL;
CREATE INDEX list_lines_version ON list_lines (list_id, version);
CREATE TABLE list_line_sources (
    line_id      uuid NOT NULL REFERENCES list_lines ON DELETE CASCADE,
    kind         smallint NOT NULL,
    ref          text NOT NULL,
    amount_value double precision,
    amount_unit  smallint NOT NULL DEFAULT 0,
    member_id    uuid,
    meal_date    date,
    meal         smallint,
    PRIMARY KEY (line_id, kind, ref)
);
```

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| kitchen down (names) | lines returned with the cached or empty name and flag `name_unavailable` |

## Acceptance criteria

- [ ] A new household gets an empty list with version 0 on first `GetList`.
- [ ] Two ingredient lines for the same ingredient cannot both be active (unique index test).
- [ ] Names follow `x-language` (pl, de, en).
