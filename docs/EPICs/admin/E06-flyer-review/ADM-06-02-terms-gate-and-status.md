# ADM-06-02 Terms gate and chain status

| | |
|---|---|
| Service | admin |
| Epic | [ADM-06 Flyer sources, offers and shops review](README.md) |
| Release | R2 |
| Size | M |
| Depends on | ADM-06-01; OFF-05-08 (`ApproveSourceTerms`, `RevokeSourceTerms`, `UpdateSource`), OFF-08-02 (`SetChainLive`, `SetChainSyncing`), OFF-08-05 (`RemoveChain`, `RestoreChain`) |
| Unblocks | ADM-06-03 |

## Goal

No chain goes live before its terms have been checked and its extraction passed the accuracy test, and a chain that
objects can be removed within a day (features.md §3 and §8).

## Business rules

1. Terms per source: `approve` → `ApproveSourceTerms` with a note of at least 20 characters (OFF-05-08 rule 3);
   admin builds it from the terms URL, the date checked and the founder's note (≤ 1,000 characters). `revoke` →
   `RevokeSourceTerms`, which also disables the source. Approval is required before a source can be enabled
   (features.md §3: "Each retailer's terms must be checked before step 2 goes live for that chain"). `GATE:` the
   approval records a legal check done outside the system.
2. Chain and source actions (offers' states, OFF-08-02):
   - `enable`: `UpdateSource(enabled = true)`; needs approved terms. The chain moves from `user_flyers_only` to
     `syncing`: offers reads its flyers without publishing.
   - `go_live`: `SetChainLive(chain_id, sample_id)` with a closed go-live sample of the last 21 days scoring
     ≥ 97.0 (ADM-06-03) and the source status `ok` (OFF-08-02 rule 2).
   - `back_to_syncing`: `SetChainSyncing(chain_id, reason)`: offers hides the chain's live offers and stops
     publishing; going live again needs a new go-live sample (OFF-08-02 rule 5). Used e.g. after a bad weekly check
     (ADM-06-04).
   - `disable`: `UpdateSource(enabled = false)`: stops checking the source.
   - `remove`: `RemoveChain(chain_id, reason)`: unpublishes every offer of the chain at once, disables the source
     and revokes its terms; stored flyer pages are deleted within 24 h (features.md §8 "Remove the chain within a
     day", OFF-08-05).
   - `restore`: `RestoreChain(chain_id)`: back to `user_flyers_only`; going live needs new terms and a new sample.
3. A reason (≤ 300 characters) is required for `back_to_syncing` and `remove`.
4. `remove` asks for the chain name typed again (`confirm_chain`).
5. offers enforces the conditions and returns `FAILED_PRECONDITION` with the unmet ones; admin shows them.
6. Other source settings (list URL, allowed hosts, check time, extra checks, days ahead) are edited with
   `UpdateSource` from the source page.
7. Audited: `offers.source.terms`, `offers.source.update`, `offers.chain.status` (action, reason).
8. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| POST | `/admin/v1/offers/sources/{id}/terms` | `{ action: approve\|revoke, terms_url?, checked_on?, note? }` → source |
| PATCH | `/admin/v1/offers/sources/{id}` | `{ list_url?, allowed_hosts?, check_time?, extra_checks?, days_ahead?, enabled? }` → source |
| POST | `/admin/v1/offers/chains/{chainId}/status` | `{ action: go_live\|back_to_syncing\|remove\|restore, sample_id?, reason?, confirm_chain? }` → chain |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Unmet condition | 409 with `unmet[]` |
| Missing reason, wrong `confirm_chain` | 400 |

## Acceptance criteria

- [ ] `go_live` without a passed go-live sample from the last 21 days → 409 with offers' unmet conditions.
- [ ] `remove` takes the chain's offers off `GET /v1/offers` immediately.
