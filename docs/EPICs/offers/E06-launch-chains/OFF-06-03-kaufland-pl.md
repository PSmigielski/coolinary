# OFF-06-03 Kaufland PL adapter (per store)

| | |
|---|---|
| Service | offers |
| Epic | [OFF-06 Launch chain adapters](README.md) |
| Release | R2 |
| Size | L |
| Depends on | OFF-05-02, OFF-05-03, OFF-03-05, OFF-05-07 |
| Unblocks | OFF-06-07, OFF-08-02 for `kaufland-pl` |

## Goal

Each Kaufland Poland store that a household uses gets its own flyer edition, read from Kaufland's own leaflet pages
(features.md §3: "Per store version"; Verification 8: two stores with different editions show different offers).

## Business rules

1. Source `kaufland-pl`: adapter `kaufland` (shared with Kaufland DE), list URL
   `https://sklep.kaufland.pl/gazeta-reklamowa.html`, flyer URL pattern
   `https://leaflets.kaufland.com/pl-PL/[id]/ar/[store]`, allowed hosts `sklep.kaufland.pl`, `leaflets.kaufland.com`
   and the image hosts found in the fixtures.
2. Editions: one per store (`kaufland-pl:store:<store>`), read only when wanted (OFF-05-07).
3. `ListFlyers(edition)` returns the flyers available for that store number; `external_id` = `[id]`; the same `[id]`
   in two stores is two flyers (OFF-05-01 rule 3). Identical page sets are extracted once (OFF-05-07 rule 6).
4. `ListStores` returns Kaufland PL's store list (store number, name, address, postcode, position) from Kaufland's
   public store finder, for linking (OFF-03-05). `GATE:` covered by the same terms approval; if the store finder may
   not be used, stores are linked by hand.
5. Rhythm: flyers run Thursday to Wednesday and are online from Tuesday (checked 4 October 2026). Daily check 06:00
   `Europe/Warsaw`; extra checks Tuesday 12:00 and Wednesday 12:00 (features.md §3 step 1 example).
6. `GATE: retailer terms checked` for kaufland.pl and leaflets.kaufland.com.

## Contract

No gRPC change.

## Data

Migration inserts the source, disabled.

## Errors

| Case | Behaviour |
|---|---|
| Store number not known to Kaufland any more | Edition check returns no flyers; after 2 such checks the store link is flagged for admin (OFF-03-05 rule 6) |

## Acceptance criteria

- [ ] Fixtures for two stores with different flyers give two editions with different pages.
- [ ] Fixtures for two stores with the same pages lead to one extraction (checked by counting `SubmitBatch` calls).
- [ ] `ListStores` fixture links at least 90% of fixture shops automatically (OFF-03-05).
- [ ] Two-week shadow run for the wanted stores detects each flyer by Wednesday evening before its Thursday start.
