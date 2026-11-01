# KIT-07-02 Default storage place and use-by date

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-07 Fridge stock](README.md) |
| Release | R1 |
| Size | M |
| Depends on | KIT-07-01, KIT-02 (default storage place and shelf life per storage place), KIT-04 (storage hint and 2D-code date from `ResolveBarcode`), KIT-10-01 (household timezone) |
| Unblocks | KIT-07-03, KIT-07-04, KIT-07-06, GW-06 (add-item form) |

## Goal

Every new item gets a sensible storage place and use-by date without typing, so the fridge list can always say
what to use first (features.md §2 "Fridge and pantry", §4 "How one scan is resolved").

## Business rules

1. **Today** is the household's local date, from the timezone in the kitchen profile (KIT-10-01). Server time and
   UTC are never used as "today".
2. **Storage place**, first that applies:
   1. the place the user chose;
   2. the storage hint of the resolved product (KIT-04);
   3. the ingredient's default storage place (KIT-02);
   4. pantry.
3. **Use-by date**, first that applies:
   1. a date the user typed → `expiry_source = user`;
   2. a date read from a GS1 2D code on the pack (KIT-04) → `pack_code`;
   3. the ingredient's shelf life for the chosen storage place (KIT-02-03, which falls back to the default
      place's value and flags `storage_unusual`): today + shelf-life days → `default`;
   4. no shelf life known → no date, `none`.
4. A user date may lie at most 30 days in the past and at most 5 years ahead; otherwise `INVALID_ARGUMENT`.
5. **Changing the storage place** of an existing item (KIT-07-04):
   1. `user` dates never change;
   2. into the freezer, with a `default` or `pack_code` date: new date = move date + freezer shelf life of the
      ingredient (or no date if the catalogue has none), `expiry_source = default`;
   3. out of the freezer into the fridge, with a `default` date: new date = move date + min(fridge shelf life,
      2 days). Decided (Q078): thawed food gets at most 2 days, or the fridge shelf life if shorter;
   4. between fridge and pantry, with a `default` date: new date = the item's `added_at` local date + the new
      place's shelf life; when moving to the pantry, the new date is never later than the old one;
   5. `pack_code` dates stay when moving between fridge and pantry.
6. Shelf lives from the catalogue are already conservative ("conservative expiry guidance", business model §10);
   this task adds no extra margin.
7. Items without an ingredient (not matched) get storage rule 2.1 or pantry, and only a user or pack-code date.
8. The same rules are exposed as `SuggestStockDefaults` so the app can pre-fill the add form; nothing is stored by
   that call.

## Contract

`contracts/coolinary/kitchen/v1/stock.proto`

```proto
rpc SuggestStockDefaults(SuggestStockDefaultsRequest) returns (SuggestStockDefaultsResponse);

message SuggestStockDefaultsRequest {
  string household_id = 1;
  string ingredient_id = 2;     // one of ingredient_id / product_ean
  string product_ean = 3;
  StoragePlace storage = 4;     // optional: the user's choice
  string pack_code_expires_on = 5;  // optional: date read from a 2D code, YYYY-MM-DD
}

message SuggestStockDefaultsResponse {
  StoragePlace storage = 1;
  string expires_on = 2;        // empty = no date
  ExpirySource expiry_source = 3;
  // Shelf life per place, so the app can update the date when the user switches place.
  map<int32, int32> shelf_life_days_by_storage = 4;  // key = StoragePlace
}
```

The same logic is an internal Go function used by KIT-07-03, KIT-07-04 and KIT-07-06.

## Data

No schema change (reads KIT-02 and KIT-10 tables).

## Errors

| Case | gRPC status | Behaviour |
|---|---|---|
| Neither ingredient nor EAN given | `INVALID_ARGUMENT` | "ingredient_id or product_ean is required" |
| Unknown ingredient or EAN | `NOT_FOUND` | |
| Date outside −30 days … +5 years | `INVALID_ARGUMENT` | |
| Household has no profile yet | — | Profile is created with country defaults first (KIT-10-01) |

## Acceptance criteria

- [ ] Household in `Europe/Warsaw`, server clock at 23:30 UTC on 1 March (00:30 on 2 March in Warsaw): milk with a
      7-day fridge shelf life gets 9 March.
- [ ] A product with a storage hint "freezer" is suggested for the freezer even if the ingredient's default is
      the fridge; the user's choice "fridge" beats both.
- [ ] A 2D-code date beats the shelf life; a typed date beats the 2D-code date.
- [ ] Chicken (fridge 2 days, freezer 90 days) moved to the freezer on 5 March gets 3 June; moved back to the fridge
      on 1 April gets 3 April; a `user` date is unchanged in both moves.
- [ ] An ingredient with no shelf life gets no date and `expiry_source = none`.
- [ ] `SuggestStockDefaults` writes nothing to the database.
