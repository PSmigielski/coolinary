# Task file template

Copy this shape for every task. File name: `<TASK-ID>-<short-slug>.md`, inside its epic folder. Write plainly:
short sentences, numbers stated exactly, no "should probably". Leave out a section only where marked optional.

````markdown
# KIT-04-03 Open Food Facts lookup

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-04 Barcode lookup chain](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-04-01, KIT-02-01 |
| Unblocks | GW-06-02 |

## Goal

One or two sentences: what the user or the system can do once this is done, and why it matters.

## Business rules

1. Numbered, each one testable on its own.
2. Every limit, default, rounding rule, ordering and tie-break is written out with its number.
3. Who may do it (owner, member, any account, admin, Free/Plus) and what happens otherwise.
4. What happens at the edges: empty input, duplicates, concurrent edits, the dependency being down.
5. Rules taken from docs/features.md cite the section, e.g. (features.md §4 "How one scan is resolved").

## Contract

`contracts/coolinary/kitchen/v1/barcode.proto` — new or changed RPCs and messages.

```proto
rpc ResolveBarcode(ResolveBarcodeRequest) returns (ResolveBarcodeResponse);
```

Gateway and mobile tasks use HTTP instead: method, path, request and response JSON, status codes.

## Data

`kitchen/migrations/000NN_off_cache.sql` — tables, columns, constraints, indexes. Write "No schema change" if none.

```sql
CREATE TABLE ...
```

## Events (optional)

Published or consumed events, with the envelope fields used. Name events exactly as in docs/EPICs/README.md.

## Errors

| Case | gRPC status / HTTP | Message or behaviour |
|---|---|---|

## Acceptance criteria

- [ ] Concrete checks a reviewer can run: given / when / then, with test data.
- [ ] Reference the verification items in docs/features.md where they apply ("Verification 14").
- [ ] Metrics, logs and privacy checks where relevant (no personal data in logs).

## Notes (optional)

Background, links to sources, and any `DECISION NEEDED:` with the recommended choice and the alternative.
````

## Sizes

`S` ≤ 1 day · `M` 2–3 days · `L` 4–5 days. Split anything larger.

## Status markers

- `DECISION NEEDED:` a choice the docs don't settle; the task states the recommended rule. Listed again in the
  service README.
- `GATE:` work that must not go live until an outside condition is met (a licence checked, a retailer's terms
  approved, an access contract signed).
