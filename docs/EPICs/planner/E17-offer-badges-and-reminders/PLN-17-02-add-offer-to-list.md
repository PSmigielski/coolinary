# PLN-17-02 Add an offer to the list

| | |
|---|---|
| Service | planner |
| Epic | [PLN-17 Offer badges and offer-ends reminders](README.md) |
| Release | R2 |
| Size | S |
| Depends on | PLN-17-01, PLN-15-08, OFF-09 (GetOffers) |
| Unblocks | GW-10 |

## Goal

Tap an offer to add it to the shopping list (features.md "Flyers and offers").

## Business rules

1. Input: offer id. Planner fetches the offer from offers (`GetOffers` by id) — the client's copy is not trusted.
2. The offer must be current or upcoming (valid_to ≥ today) and at one of the household's shops; otherwise
   `FAILED_PRECONDITION` (`offer_not_available`).
3. Offer matched to a catalogue ingredient → ingredient line; if an open line for that ingredient exists, the offer
   is pinned on it and its amount is raised to at least one offer pack. Unmatched → free-text line with the offer's
   product name.
4. Amount: one offer pack. Multi-buy conditions ("2+1") are shown, not applied to the amount.
5. Source kind `offer`, ref offer id; the badge is pinned (PLN-17-01 rule 7).
6. Idempotent by `client_mutation_id`; adding the same offer twice keeps one contribution.

## Contract

`contracts/coolinary/planner/v1/list.proto`

```proto
rpc AddOfferToList(AddOfferToListRequest) returns (ListLine);
message AddOfferToListRequest { string household_id = 1; string member_id = 2; string offer_id = 3; string client_mutation_id = 4; }
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Offer unknown | `NOT_FOUND` |
| Expired or not at the household's shops | `FAILED_PRECONDITION` `offer_not_available` |
| offers down | `UNAVAILABLE` |

## Acceptance criteria

- [ ] Adding "Basmati rice 1 kg, 2+1" adds a rice line of 1 kg with a pinned badge showing "2+1".
- [ ] Adding an unmatched "Party pack mixed sweets" offer adds a free-text line with a badge.
- [ ] An expired offer is refused.

## Notes

- Needs an offers RPC to read offers by id (`GetOffers`), not in the README's RPC table. Flagged to offers (OFF-09).
