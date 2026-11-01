# PLN-15-05 Generate, review and accept proposed lines

| | |
|---|---|
| Service | planner |
| Epic | [PLN-15 Shopping list](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | PLN-15-04 |
| Unblocks | GW-08, PLN-15-10, PLN-18-04 |

## Goal

The user reviews the proposed lines before they join the list and can untick any of them; each line shows where it
came from (features.md "Shopping list": "Also without AI").

## Business rules

1. `GenerateList` runs PLN-15-02 → PLN-15-04 for the requested dates and returns a proposal; nothing is written to the
   list yet.
2. The proposal is stored for 24 hours with an id, the plan range, and the list version and stock version it was
   built on.
3. Proposed lines are ticked (wanted) by default, except optional-only lines, which start unticked.
4. Needs already on the list from the same plan meal (a source `plan_meal` with that ref) are not proposed again;
   only the difference is proposed (so generating twice does not double anything).
5. `AcceptListProposal` takes the proposal id and the keys of lines the user wants. For each wanted line:
   - an ingredient with an open line on the list: contributions are added to it and its packs re-rounded;
   - otherwise a new line is created.
   Unwanted lines are dropped and remembered on the proposal only (not on the list).
6. Accepting an expired proposal returns `FAILED_PRECONDITION` (`proposal_expired`); the app regenerates.
7. If the stock or the list changed since the proposal was built, accept still applies it but recomputes the
   stock subtraction for the wanted lines first; lines that became fully covered are skipped and reported.
8. One accept is one list change (one version bump) and one live notice (PLN-15-09).
9. Each mutation carries a `client_mutation_id` (PLN-15-08).
10. Any account member may generate and accept.

## Contract

`contracts/coolinary/planner/v1/list.proto`

```proto
rpc GenerateList(GenerateListRequest) returns (ListProposal);
rpc AcceptListProposal(AcceptListProposalRequest) returns (ShoppingList);

message GenerateListRequest { string household_id = 1; string from = 2; string to = 3; repeated Meal meals = 4; }
message ListProposal {
  string id = 1; string expires_at = 2;
  repeated ProposedLine lines = 3;
  repeated string covered_ingredient_ids = 4;   // fully at home, shown as "you have these"
}
message ProposedLine {
  string key = 1; string ingredient_id = 2; string name = 3; string category = 4;
  coolinary.common.v1.Quantity amount = 5; coolinary.common.v1.Quantity pack_size = 6; int32 pack_count = 7;
  repeated LineSource sources = 8; bool wanted = 9; repeated string flags = 10;
}
message AcceptListProposalRequest {
  string household_id = 1; string member_id = 2; string proposal_id = 3;
  repeated string wanted_keys = 4; string client_mutation_id = 5;
}
```

## Data

`planner/migrations/000NN_list_proposals.sql`

```sql
CREATE TABLE list_proposals (
    id            uuid PRIMARY KEY,
    household_id  uuid NOT NULL,
    range_from    date NOT NULL,
    range_to      date NOT NULL,
    list_version  bigint NOT NULL,
    stock_version text NOT NULL,
    lines         jsonb NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    expires_at    timestamptz NOT NULL
);
CREATE INDEX list_proposals_expiry ON list_proposals (expires_at);
```

Expired proposals are deleted by an hourly job.

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Unknown or other household's proposal | `NOT_FOUND` |
| Expired | `FAILED_PRECONDITION` `proposal_expired` |
| kitchen down at generate | `UNAVAILABLE` |

## Acceptance criteria

- [ ] Verification 12: for the fixed test plan and test fridge, the proposal holds exactly the missing ingredients,
      each once, with amounts merged across dishes; nothing at home appears.
- [ ] Unticking garlic and accepting: garlic is not on the list; everything else is, with sources "dinner Tue: chicken
      curry" etc.
- [ ] Generating and accepting the same day twice leaves the list unchanged the second time.
- [ ] Accepting with the same `client_mutation_id` twice returns the same list version.
