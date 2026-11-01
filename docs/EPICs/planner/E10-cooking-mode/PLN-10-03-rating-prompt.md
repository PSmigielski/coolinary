# PLN-10-03 Rating prompt after cooking

| | |
|---|---|
| Service | planner |
| Epic | [PLN-10 Cooking mode and "cooked it"](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | PLN-10-02, PLN-05-05 |
| Unblocks | PLN-11-01 (ratings of AI dishes) |

## Goal

Ask for a rating at the moment people remember the dish, without nagging, so the library sorts well and good AI dishes
can join it.

## Business rules

1. `MarkCookedResponse.rating_prompt = true` when the account has not rated this dish yet, or rated it more than 90 days
   ago, and has not dismissed the prompt for this dish.
2. At most one prompt per account per day: a second "cooked it" the same day returns `false`.
3. `DismissRatingPrompt(dish_id)` stores the dismissal; it is never asked again for that dish.
4. AI dishes (household AI) are always prompted on their first cook, ignoring rule 2 (their ratings feed PLN-11).

## Contract

```proto
rpc DismissRatingPrompt(DismissRatingPromptRequest) returns (DismissRatingPromptResponse);
message DismissRatingPromptRequest { string household_id = 1; string user_id = 2; string dish_id = 3; }
```

## Data

`planner/migrations/000NN_rating_prompts.sql`

```sql
CREATE TABLE rating_prompts (
    user_id       text NOT NULL,
    dish_id       uuid NOT NULL REFERENCES dishes (id) ON DELETE CASCADE,
    last_prompted date,
    dismissed     boolean NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, dish_id)
);
```

## Errors

| Case | gRPC status / behaviour |
|---|---|
| Missing `user_id` or `dish_id` | `INVALID_ARGUMENT` |
| Unknown dish | `NOT_FOUND` |
| Dismissing twice | success, no change |

## Acceptance criteria

- [ ] First cook of an unrated dish → prompt; second dish cooked the same day → no prompt (unless it is an AI dish's first cook).
- [ ] After dismissing, cooking the dish again never prompts.
