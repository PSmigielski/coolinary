# ADM-03-01 Name-match review queue

| | |
|---|---|
| Service | admin |
| Epic | [ADM-03 Catalogue, matching and nutrition review](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | ADM-02-03; KIT-03-06 (`KitchenAdminService.ListMatchReviews`, `DecideMatch`) |
| Unblocks | — |

## Goal

Names that Haiku could not match to a catalogue ingredient with confidence (product names, receipt lines, flyer
offers, recipe import lines) are resolved by a human once, and the answer then serves everyone (decision 13).

## Business rules

1. The queue lists what kitchen puts in it (KIT-03 decides when: "none" answers and answers under its confidence
   threshold). Each item shows what `ListMatchReviews` returns (KIT-03-06 rule 1): example name, normalised key,
   language, contexts (`product`, `receipt`, `offer`, `recipe`), current outcome and ingredient, confidence, method,
   hit count, and the top 5 trigram candidates.
2. Order as kitchen returns it: most seen first, then most recently seen. Filters: language, context. Only open
   items are listed; decided ones are seen in the item's history (`ListAdminChanges`, KIT-01-07).
3. Actions on one item, each one `DecideMatch(match_id, action)` call:
   - `approve`: keep the current answer (the AI pick);
   - `set_ingredient`: another ingredient (`ingredient_id`, searched with `KitchenService.SearchIngredients`);
     kitchen also adds the name as a synonym (KIT-03-06 rule 3);
   - `set_none`: the name is not food or not useful (e.g. "deposit bottle"); kitchen stores it so the name is not
     asked again;
   - `create_ingredient`: the admin fills the ingredient form of ADM-03-03 and kitchen creates and matches it in one
     call.
4. Bulk: up to 50 items decided with the same action and ingredient in one request; each item gets its own result.
5. Two admins deciding the same item: the later decision wins; both are in kitchen's `admin_changes` and in the
   admin audit log.
6. Each resolution is audited (`kitchen.match.resolve`, item id, action, ingredient id).
7. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/kitchen/match-reviews?status=&language=&context=&page_token=` | `{ items: MatchReview[], next_page_token }` |
| POST | `/admin/v1/kitchen/match-reviews/resolve` | `{ items: [{ id, action: approve\|set_ingredient\|set_none\|create_ingredient, ingredient_id?, ingredient? }] }` → `{ results: [{ id, ok, error? }] }` |

Calls `KitchenAdminService.ListMatchReviews`, `KitchenAdminService.DecideMatch` (one call per item, at most 5 in
parallel; KIT-03-06).

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| More than 50 items | 400 |
| `set_ingredient` without `ingredient_id`, or `create_ingredient` without the form | 400 |
| Unknown or retired ingredient | per-item 400 from kitchen's `INVALID_ARGUMENT` / `NOT_FOUND` |

## Acceptance criteria

- [ ] Deciding "Hähnchenschenkel 1kg" (de, offer) with `set_ingredient chicken_thigh` removes it from the list and
      kitchen matches the next identical name without asking.
- [ ] A bulk request of 3 items where one names an unknown ingredient returns 2 ok and 1 error.
- [ ] Every resolution appears in the audit log.
