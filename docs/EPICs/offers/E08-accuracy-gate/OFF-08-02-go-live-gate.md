# OFF-08-02 Go-live gate per chain

| | |
|---|---|
| Service | offers |
| Epic | [OFF-08 Accuracy gate and weekly spot-checks](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-08-01, OFF-07-07 |
| Unblocks | users seeing synced offers of each chain |

## Goal

A chain's synced offers become visible only after its terms are approved and a 100-offer sample passes 97%
(features.md Risks: "a chain goes live only at ≥ 97% correct prices and dates on a 100-offer sample", Verification 5).

## Business rules

1. Coverage moves: `user_flyers_only` → `syncing` when an admin enables its source (terms approved, OFF-05-08);
   `syncing` → `live` only through `SetChainLive`.
2. `SetChainLive` succeeds only when: the source's terms are approved, a `go_live` sample for the chain was closed in
   the last 21 days with score ≥ 97.0, and the source status is `ok`.
3. On going live, every `ready` offer of the chain that is still valid is published at once (OFF-07-07 rule 2), and
   `offers.published` events are sent.
4. A failed sample keeps the chain `syncing`; the founder fixes the adapter or prompts (AI-13) and draws a new sample.
   There is no limit on attempts.
5. `SetChainSyncing` takes a live chain back to `syncing` (offers hidden, `offers.expired` sent for the visible ones),
   e.g. after a bad weekly check (OFF-08-03). Going live again needs a new go-live sample.
6. Every transition writes `admin_changes` and sends an info alert (`offers.chain.<chain_id>`).

## Contract

`admin.proto`

```proto
rpc SetChainLive(SetChainLiveRequest) returns (Chain);       // chain_id, sample_id
rpc SetChainSyncing(SetChainSyncingRequest) returns (Chain); // chain_id, reason
```

## Data

No schema change (`chains.coverage`).

## Errors

| Case | gRPC status |
|---|---|
| No passing sample in 21 days | `FAILED_PRECONDITION` "needs a go-live sample ≥ 97%" |
| Terms not approved or source not ok | `FAILED_PRECONDITION` |

## Acceptance criteria

- [ ] With a closed sample at 96.9, `SetChainLive` is refused; at 97.0 it succeeds and the ready offers become
      visible in `BrowseOffers`.
- [ ] `SetChainSyncing` hides the chain's offers immediately and planner receives `offers.expired`.
