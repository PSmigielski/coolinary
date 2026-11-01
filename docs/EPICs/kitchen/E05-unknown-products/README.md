# KIT-05 Unknown products and shared-catalogue confirmation

| | |
|---|---|
| Service | kitchen |
| Release | R1 |
| Unblocks | GW-06 (naming and label routes), Verification 14 ("an unknown product named once is recognised on the next scan from a second account") |

## Goal

When no source knows a code, the user names the product once (typing, or a label photo read by AI), and it joins
the shared catalogue for everyone once a second household accepts it unchanged (features.md §4 step 5). Products
confirmed this way can be sent back to Open Food Facts.

## Flow

```
no hit ──> user types name (or label photo ──AI-07──> proposal) ──NameProduct──> household product + submission
next household scans ──ResolveBarcode──> proposal "Is this X?" ──ConfirmProposal──> accepted unchanged → confirmed
conflicts / reports ──> admin review (ADM-03) ;  confirmed user products ──(GATE)──> Open Food Facts
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-05-01](KIT-05-01-name-product.md) | Name an unknown product once | R1 | M | KIT-04-04, KIT-04-05, KIT-04-09 |
| [KIT-05-02](KIT-05-02-submissions-and-confirmation.md) | Submissions, proposals and confirmation by a second household | R1 | L | KIT-05-01 |
| [KIT-05-03](KIT-05-03-label-photo.md) | Read a product label photo | R1 | M | KIT-05-01, AI-07 |
| [KIT-05-04](KIT-05-04-contribute-to-off.md) | Send confirmed products to Open Food Facts | R1 (GATE) | M | KIT-05-02, KIT-04-07 |
| [KIT-05-05](KIT-05-05-admin-product-review.md) | Admin product review | R1 | M | KIT-05-02, KIT-01-07 |

## Decisions (taken 1 November 2026)

- Q073 (KIT-05-01): 50 per day, then refused until the next day.
- Q074 (KIT-05-02): 3 conflicting names within 30 days, or 3 households reporting it wrong.
- Q075 (KIT-05-02): Yes: a word list per language (pl, de, en) before names are shared.
- Q076 (KIT-05-04): A household setting, on by default, keeps its confirmed product names out of what is sent to Open Food Facts. (changed from the first draft)
- Q110 (KIT-05-03): No: read name, brand and size only in R1.
