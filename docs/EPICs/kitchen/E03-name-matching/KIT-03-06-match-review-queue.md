# KIT-03-06 Match review queue (admin)

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-03 Name matching](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | KIT-03-02, KIT-01-07 |
| Unblocks | ADM-03 |

## Goal

The founder checks uncertain matches, most frequent first, and fixes them so each fix helps every household.

## Business rules

1. `ListMatchReviews(language?, context?, page)` returns rows with `review = needed`, ordered by `hit_count`
   descending, then `last_seen_at` descending; 50 per page, maximum 200. Each item shows example, key, language,
   contexts, current outcome and ingredient, confidence, method, hit count, and the top 5 trigram candidates.
2. `DecideMatch(match_id, action)` with action:
   - `approve`: keep the current answer;
   - `set_ingredient(ingredient_id)`: matched to that active ingredient;
   - `set_none`: nothing fits;
   - `create_ingredient(fields)`: creates the ingredient (KIT-02-10 rules) and matches to it.
   The row gets `method = 'admin'`, `confidence = 1`, `review = 'done'`, `decided_by`, `decided_at`, and an
   `admin_changes` record.
3. `set_ingredient` also adds the example's key as a `synonym` of the ingredient in that language when no other
   ingredient has it, so similar names match exactly next time.
4. Decided (Q058): when an admin changes an existing `matched` answer, kitchen writes the event
   `catalogue.match_changed` `{language, key, from_ingredient_id, to_ingredient_id}` (README event table) so offers
   and planner fix references they stored.
5. Gauge `kitchen_match_reviews_open{language}`.

## Contract

`contracts/coolinary/kitchen/v1/admin.proto`

```proto
rpc ListMatchReviews(ListMatchReviewsRequest) returns (ListMatchReviewsResponse);
rpc DecideMatch(DecideMatchRequest) returns (NameMatchRecord);
message DecideMatchRequest { int64 match_id = 1;
  oneof action { bool approve = 2; string set_ingredient_id = 3; bool set_none = 4; CreateIngredientRequest create_ingredient = 5; } }
```

## Data

No schema change.

## Events

Possibly `catalogue.match_changed` v1 (see rule 4).

## Errors

| Case | Status | Message |
|---|---|---|
| Unknown match id | `NOT_FOUND` | `match not found` |
| Target ingredient not active | `FAILED_PRECONDITION` | `ingredient is not active` |

## Acceptance criteria

- [ ] A review row with 50 hits is listed before one with 3.
- [ ] After `set_ingredient`, `MatchNames` for the same name returns that ingredient with method `admin`, and a
      batch result arriving later does not change it.
- [ ] `set_ingredient` adds the synonym; a second, similar name now matches on the exact path.
