# OFF-04-02 Shopping preferences ("How you shop")

| | |
|---|---|
| Service | offers |
| Epic | [OFF-04 My shops and shopping preferences](README.md) |
| Release | R2 |
| Size | S |
| Depends on | — |
| Unblocks | GW-10, MOB-03 (the phone's trip plan reads them) |

## Goal

The household's travel range and trip settings are stored once and shared by every device in the household; the
phone uses them to filter shops and to plan the trip (features.md §1 "My shops and location", §2 "Planning the
shopping trip").

## Business rules

1. Settings and allowed values:

   | Setting | Allowed | Default |
   |---|---|---|
   | `range_km` | 1 (walking), 3 (short trip), 10 (by car) | 3 (Decided, Q164) |
   | `walk_between_m` | 0 (one shop only) to 1,500, in steps of 100 | 500 (Decided, Q165) |
   | `shops_per_trip` | 1, 2, 3 | 2 (Decided, Q163) |
   | `min_extra_shop_saving` | 0 to 50.00 in the country's currency, in minor units | EUR 2.00 in euro countries; PLN 8.00 in Poland (Decided, Q166) |

2. `walk_between_m = 0` forces `shops_per_trip = 1` (only single shops, Verification 13); the stored value is then 1.
3. The currency of `min_extra_shop_saving` is the currency of the household's my shops; with no shops yet it is the
   currency of the country passed by the phone.
4. A household without stored settings gets the defaults; reading never fails.
5. `SetShoppingPreferences` replaces all four values; `updated_at` is returned. Last write wins.
6. Any member with an account may change them. Free households may store them too; the trip plan itself is Plus and
   runs on the phone (README decision 7).
7. These values are not used by offers for any filtering (offers never knows distances).

## Contract

```proto
rpc GetShoppingPreferences(GetShoppingPreferencesRequest) returns (ShoppingPreferences);
rpc SetShoppingPreferences(SetShoppingPreferencesRequest) returns (ShoppingPreferences);

message ShoppingPreferences {
  int32 range_km = 1;
  int32 walk_between_m = 2;
  int32 shops_per_trip = 3;
  coolinary.common.v1.Money min_extra_shop_saving = 4;
  string updated_at = 5;            // RFC 3339
}
message GetShoppingPreferencesRequest { string household_id = 1; string country_code = 2; }
message SetShoppingPreferencesRequest { string household_id = 1; ShoppingPreferences preferences = 2; }
```

## Data

`offers/migrations/000NN_shopping_preferences.sql`

```sql
CREATE TABLE shopping_preferences (
    household_id   text PRIMARY KEY,
    range_km       smallint NOT NULL CHECK (range_km IN (1, 3, 10)),
    walk_between_m smallint NOT NULL CHECK (walk_between_m BETWEEN 0 AND 1500 AND walk_between_m % 100 = 0),
    shops_per_trip smallint NOT NULL CHECK (shops_per_trip BETWEEN 1 AND 3),
    min_saving_minor int NOT NULL CHECK (min_saving_minor BETWEEN 0 AND 5000),
    currency       text NOT NULL,
    updated_at     timestamptz NOT NULL DEFAULT now()
);
```

## Errors

| Case | gRPC status |
|---|---|
| Value outside the allowed set | `INVALID_ARGUMENT` naming the field |
| Currency not the household's | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] A new household reads `range_km = 3`, `walk_between_m = 500`, `shops_per_trip = 2`, EUR 2.00 (DE) or PLN 8.00
      (PL).
- [ ] Setting `walk_between_m = 0` with `shops_per_trip = 3` stores `shops_per_trip = 1`.
- [ ] `walk_between_m = 250` is refused.
