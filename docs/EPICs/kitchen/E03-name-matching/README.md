# KIT-03 Name matching (trigram + Haiku)

| | |
|---|---|
| Service | kitchen |
| Release | R1 |
| Unblocks | KIT-04-10 (product → ingredient), KIT-08 (photo/text/receipt items), OFF-07 (offers), PLN-03 (recipe import) |

## Goal

Turn any product, receipt, offer or typed name ("Hähnchenschenkel 1kg", "MLEKO LACIATE 3,2%") into one catalogue
ingredient, or an honest "none", once for everyone (README decision 13). Each normalised name is matched once and
stored; later requests are a database lookup.

## Flow

```
name ──normalise──> key ──stored match?──yes──> answer
                          │ no
                          ├─ exact name/synonym (one ingredient) ──> matched, stored
                          ├─ pg_trgm top 10 ──clear winner──> matched, stored
                          ├─ no candidate ≥ 0.3 ──> none, review
                          └─ otherwise: interactive → Haiku pick (AI-08) ; bulk → queued → batch (AI-04)
low confidence / none ──> admin review queue (ADM-03)
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-03-01](KIT-03-01-name-normaliser.md) | Name normaliser per language | R1 | M | — |
| [KIT-03-02](KIT-03-02-match-store.md) | Stored matches | R1 | S | KIT-03-01, KIT-02-01 |
| [KIT-03-03](KIT-03-03-match-names-deterministic.md) | MatchNames: stored, exact and trigram paths | R1 | M | KIT-03-02, KIT-02-04 |
| [KIT-03-04](KIT-03-04-match-names-ai-pick.md) | MatchNames: Haiku pick and fallbacks | R1 | M | KIT-03-03, KIT-01-01, AI-08 |
| [KIT-03-05](KIT-03-05-bulk-matching-batch.md) | Bulk matching through Message Batches | R1 | M | KIT-03-04, KIT-01-08 |
| [KIT-03-06](KIT-03-06-match-review-queue.md) | Match review queue (admin) | R1 | M | KIT-03-02, KIT-01-07 |
| [KIT-03-07](KIT-03-07-rematch-on-catalogue-change.md) | Re-match after catalogue changes | R1 | S | KIT-03-03, KIT-03-05, KIT-02-10 |

## Decisions (taken 1 November 2026)

- Q052 (KIT-03-01): No splitting; Haiku resolves compounds when text matching is unsure.
- Q057 (KIT-03-03): Best score ≥ 0.85 and ahead of the runner-up by ≥ 0.15.
- Q058 (KIT-03-06): Yes: new event catalogue.match_changed.
- Q064 (KIT-03-04): ≥ 0.8 accept; 0.5–0.8 accept and review; below 0.5 none, review.
- Q065 (KIT-03-05): Ask again: offers re-calls MatchNames in bulk mode every 30 minutes.
