# OFF-10 User-supplied flyers

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | A Plus household photographs, shares (PDF) or links any flyer that is not covered yet and gets its offers extracted, in any country, from the first flyer release (features.md §2 "Flyers and offers … With AI", §3 source 1 "Flyers the user supplies"). |

## Rules for the whole epic

- **Plus only** (features.md §6, "own-flyer scan" in the Plus list). Free households get `FAILED_PRECONDITION`
  (`plus_required`). Own-flyer scans never count toward the Free photo-scan quota (docs/EPICs/README.md
  "Free and Plus"; the gateway does not meter this route).
- Photos and PDFs are processed in memory and never stored (README D6).
- Extracted offers are checked without AI and confirmed by the user before they are saved.

## Flow

1. The user picks the shop (or adds one), then uploads photos or a PDF, or pastes a link.
2. offers sends the pages to `ai` (`ExtractFlyerOffers`, per request) and checks the result with the same code as
   chain flyers.
3. The user sees a draft: offers with warnings, detected dates. They untick or fix lines and confirm the dates.
4. Confirmed offers are saved for that household and shop until their end date, and feed badges and deals for that
   household.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-10-01](OFF-10-01-upload-and-jobs.md) | Upload photos or a PDF | R2 | M | OFF-07-01 |
| [OFF-10-02](OFF-10-02-flyer-link.md) | Flyer from a link | R2 | M | OFF-10-01 |
| [OFF-10-03](OFF-10-03-extract-draft.md) | Extract and check a draft | R2 | M | OFF-10-01, OFF-07-03, OFF-07-04, OFF-07-05, AI-13 (ExtractFlyerOffers) |
| [OFF-10-04](OFF-10-04-confirm-private-offers.md) | Confirm and save private offers | R2 | M | OFF-10-03, OFF-07-07 |
| [OFF-10-05](OFF-10-05-shop-for-user-flyer.md) | Shop for a user flyer, in any country | R2 | S | OFF-02-05 |
