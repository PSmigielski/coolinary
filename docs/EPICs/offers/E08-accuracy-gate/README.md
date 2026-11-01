# OFF-08 Accuracy gate and weekly spot-checks

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | A chain's offers reach users only after a 100-offer sample shows at least 97% correct prices and dates; every live chain is spot-checked on 20 offers each week; a chain can be taken down within a day (features.md §3 step 5, Risks "A flyer price is read wrongly" and "A retailer objects", Verification 5). |

## Flow

1. While a chain is `syncing`, admin draws a 100-offer sample from its extracted offers and marks each one right or
   wrong against the flyer page.
2. At ≥ 97% correct (and terms approved) the chain goes `live`; its ready offers are published.
3. Every Monday a 20-offer sample per live chain is drawn and reviewed within the week; results build the chain's
   accuracy record.
4. A kill switch removes a chain's offers at once and deletes its stored pages within 24 hours.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-08-01](OFF-08-01-review-samples.md) | Review samples and verdicts | R2 | M | OFF-07-05, OFF-05-05, OFF-01-06 |
| [OFF-08-02](OFF-08-02-go-live-gate.md) | Go-live gate per chain | R2 | M | OFF-08-01, OFF-07-07 |
| [OFF-08-03](OFF-08-03-weekly-spot-checks.md) | Weekly spot-checks | R2 | S | OFF-08-01, OFF-01-05 |
| [OFF-08-04](OFF-08-04-accuracy-record.md) | Accuracy record and metrics | R2 | S | OFF-08-01 |
| [OFF-08-05](OFF-08-05-remove-chain.md) | Remove a chain within a day | R2 | M | OFF-07-07, OFF-05-05 |
