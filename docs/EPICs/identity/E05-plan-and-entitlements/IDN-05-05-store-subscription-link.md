# IDN-05-05 Store subscription link

| | |
|---|---|
| Service | identity |
| Epic | [IDN-05 Plan and entitlements](README.md) |
| Release | R1 |
| Size | L |
| Depends on | IDN-05-01 |
| Unblocks | Paywall purchase flow (frontend, later) |

`GATE:` decision 8 — payments are decided later (RevenueCat, direct App Store / Google Play notifications, or
another route). Nothing here is built before that decision. The rules below hold whatever the route.

## Goal

A purchase or renewal in the App Store or Google Play turns the buyer's household into Plus, and a cancellation,
refund or lapse turns it back, without manual work.

## Business rules

1. Only the owner can buy Plus for a household; the purchase is tied to the owner's account
   (features.md "Sharing needs Plus, held by the owner").
2. An anonymous account must sign in with an e-mail code before buying, so the purchase can be restored on a new
   phone.
   Decided (Q026): sign-in comes first; anonymous accounts can't buy.
3. Plans: monthly and annual, annual first, with a 7-day trial on the annual plan (business-model §3). Trial time
   counts as Plus.
4. Every store event is applied through `changePlan` with `source = store` and `expiresAt` = end of the paid or
   trial period; renewals move the expiry.
5. Store notifications are idempotent by the provider's event id; replaying one changes nothing.
6. Refund, chargeback or lapse beyond the store's grace period → Free (with IDN-05-04 for shared households).
7. A manual Plus (IDN-10-03) and a store Plus can overlap: the later expiry wins; a store lapse doesn't end a
   manual Plus that runs longer.
8. If the owner leaves or transfers ownership (IDN-08), what happens to a store subscription is decided together
   with this task.
9. Prices per country tier (Tier A €4.99 / €34.99, Tier B ≈ €2.99 / ≈ €19.99) are set in the stores, not in
   identity.

## Contract

To be designed with the payment route: a webhook endpoint (through the gateway or Caddy) and a "restore purchases"
call.

## Data

To be designed; at least a table of applied store events (provider, event id, household id, time).

## Errors

To be completed with the payment design. Known cases now:

| Case | Status / behaviour |
|---|---|
| Purchase attempt by a non-owner | `PERMISSION_DENIED`, `x-error-code: owner_only` |
| Purchase attempt by an anonymous account | `FAILED_PRECONDITION`, `x-error-code: sign_in_required` (rule 2, pending its decision) |
| Store notification with a bad signature | rejected (HTTP 401 at the webhook); nothing applied; alert after 5 in an hour |
| Replayed store notification | success, no change (rule 5) |
| Notification for an unknown household or account | stored as `unmatched` for admin review; plan unchanged |

## Acceptance criteria

To be written with the design. At minimum: purchase → Plus within 1 minute; refund → Free; replayed notification →
no change.

## Notes

- Size L is a placeholder: re-estimate once the payment route is chosen.
