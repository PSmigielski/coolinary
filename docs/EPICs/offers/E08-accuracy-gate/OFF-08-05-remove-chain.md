# OFF-08-05 Remove a chain within a day

| | |
|---|---|
| Service | offers |
| Epic | [OFF-08 Accuracy gate and weekly spot-checks](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-07-07, OFF-05-05 |
| Unblocks | — |

## Goal

When a retailer objects, its flyers stop being read and its offers vanish at once, and its stored pages are deleted
within 24 hours; user-supplied flyers keep working (features.md Risks "A retailer objects: Remove the chain within a
day; user-supplied flyers keep working").

## Business rules

1. `RemoveChain(chain_id, reason)`, in one transaction:
   - coverage → `removed`; the source → disabled with terms approval revoked;
   - every published chain-flyer or feed offer of the chain → `expired` (reason `chain_removed`), with
     `offers.expired` events;
   - pending extraction jobs of the chain are ignored when they finish;
   - price observations from its flyers stop being used by `GetPrices` (receipt prices stay; they are users' own
     data and not the retailer's).
2. Stored flyer pages of the chain are deleted by a job that starts immediately and must finish within 24 h; when
   done, an info alert confirms it. Extracted raw rows are deleted with them.
3. Users of the chain's shops see flyer status `not_covered` (OFF-05-06) and can scan flyers themselves (OFF-10).
4. Shops of the chain stay (they are OSM data, not the retailer's).
5. `RestoreChain` sets coverage back to `user_flyers_only`; going live again needs new terms approval and a go-live
   sample.
6. The reason text and the admin are recorded.

## Contract

`admin.proto`

```proto
rpc RemoveChain(RemoveChainRequest) returns (Chain);    // chain_id, reason (required)
rpc RestoreChain(RestoreChainRequest) returns (Chain);
```

## Data

No schema change.

## Errors

| Case | Behaviour |
|---|---|
| Chain already removed | `FAILED_PRECONDITION` |
| Empty reason | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] After `RemoveChain`, `BrowseOffers` returns no offer of the chain, and planner gets `offers.expired` for all of
      them within one minute.
- [ ] Within 24 h the chain's page files are gone from the volume and an alert confirms it.
- [ ] A household's private user-flyer offers for a shop of that chain are still visible.
