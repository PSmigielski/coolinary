# KIT-05-02 Submissions, proposals and confirmation by a second household

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-05 Unknown products](README.md) |
| Release | R1 |
| Size | L |
| Status | todo |
| Depends on | KIT-05-01 |
| Unblocks | KIT-05-04, KIT-05-05, Verification 14 |

## Goal

A product named by one household is offered to the next as "Is this X?", and becomes part of the shared catalogue
when a second household accepts it unchanged (features.md §4 step 5).

## Business rules

1. Each household has at most one submission per code: its latest values (name, brand, pack size, units per pack,
   ingredient) and a `fingerprint` = hash of (normalised name key, normalised brand key, pack value and unit, units
   per pack, ingredient id).
2. A submission is `named` (the household typed it), `accepted` (it accepted a proposal unchanged) or `edited` (it
   changed a proposal; the edit is a new fingerprint).
3. Confirmation: when two distinct households hold submissions with the same fingerprint, the shared product
   (KIT-04-05) gets those values, `status = 'confirmed'`, `confirmed_at = now()`. Earlier `pending` values are
   replaced.
4. Proposal: a household that scans a code with no confirmed product but pending submissions from other households
   gets the fingerprint held by the most households (tie: the earliest) as `proposal` in `ResolveBarcode`. It never
   says who named it.
5. A proposal is shown only if its name and brand pass the normaliser's `not_a_name` rule and the profanity list.
   Decided (Q075): profanity word lists per language (pl, de, en) in
   `internal/kitchen/names/profanity/<lang>.txt`; a name or brand containing a listed word is never proposed to
   other households.
6. `ConfirmProposal(household_id, code, fingerprint_id, accept, edits?)`: `accept = true` with no edits creates an
   `accepted` submission; with edits an `edited` one. Either way the household's own product row is written
   (KIT-04-04).
7. Decided (Q074): disputes. Three distinct fingerprints within 30 days without one reaching two
   households → `status = 'disputed'`, sent to admin review (KIT-05-05). A confirmed product reported wrong
   (`ReportProduct`) by three distinct households → `disputed`. A disputed product is not answered from level 2
   until resolved; households keep their own rows.
8. A confirmed product never changes automatically after confirmation; only admin changes it.
9. On household deletion (KIT-13) its submissions keep their values but lose the household link
   (`household_id = NULL`); they still count for confirmation that already happened, not for new ones.

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto`

```proto
rpc ConfirmProposal(ConfirmProposalRequest) returns (ConfirmProposalResponse);
rpc ReportProduct(ReportProductRequest) returns (ReportProductResponse);
message ProductProposal { int64 fingerprint_id = 1; string name = 2; string brand = 3;
  coolinary.common.v1.Quantity pack_size = 4; int32 units_per_pack = 5; string ingredient_id = 6; }
message ConfirmProposalRequest { string household_id = 1; string code = 2; int64 fingerprint_id = 3; bool accept = 4;
  NameProductRequest edits = 5; }
message ConfirmProposalResponse { HouseholdProduct product = 1; bool product_confirmed = 2; }
message ReportProductRequest { string household_id = 1; string code = 2; string reason = 3; } // reason: wrong_name|wrong_size|not_food|other
```

## Data

`kitchen/migrations/000NN_product_submissions.sql`

```sql
CREATE TABLE product_fingerprints (
    id              bigserial PRIMARY KEY,
    gtin            char(14) NOT NULL,
    fingerprint     text NOT NULL,
    name            text NOT NULL,
    brand           text NOT NULL DEFAULT '',
    pack_value      double precision,
    pack_unit       smallint NOT NULL DEFAULT 0,
    units_per_pack  integer NOT NULL DEFAULT 1,
    ingredient_id   uuid REFERENCES ingredients (id),
    created_at      timestamptz NOT NULL DEFAULT now(),
    UNIQUE (gtin, fingerprint)
);
CREATE TABLE product_submissions (
    gtin            char(14) NOT NULL,
    household_id    uuid,
    fingerprint_id  bigint NOT NULL REFERENCES product_fingerprints (id),
    kind            text NOT NULL CHECK (kind IN ('named','accepted','edited')),
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    UNIQUE (gtin, household_id)
);
CREATE TABLE product_reports (
    gtin          char(14) NOT NULL,
    household_id  uuid,
    reason        text NOT NULL CHECK (reason IN ('wrong_name','wrong_size','not_food','other')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (gtin, household_id)
);
```

## Errors

| Case | Status | Message |
|---|---|---|
| Fingerprint not for this code | `INVALID_ARGUMENT` | `proposal does not match code` |
| Household confirms its own proposal | `FAILED_PRECONDITION` | `cannot confirm your own product` |
| Report on a code that isn't confirmed | `FAILED_PRECONDITION` | `product is not shared` |

## Acceptance criteria

- [ ] Verification 14: household A names a code; household B scans it, gets the proposal, accepts unchanged; household
      C's scan now resolves from the shared catalogue.
- [ ] B editing the size instead creates a second fingerprint; the product stays pending.
- [ ] Three conflicting fingerprints put the product in review.
- [ ] The proposal response contains no household id or account data.
