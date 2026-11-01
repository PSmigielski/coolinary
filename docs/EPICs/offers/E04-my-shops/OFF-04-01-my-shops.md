# OFF-04-01 My shops: set and list

| | |
|---|---|
| Service | offers |
| Epic | [OFF-04 My shops and shopping preferences](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-03-03, OFF-03-04 |
| Unblocks | OFF-04-03, OFF-09-01, OFF-11-04, PLN-17, PLN-18, PLN-19, GW-10 |

## Goal

A household confirms the shops it uses; offers keeps them in the order the phone gave (nearest first) and serves
them to the app and to planner. Free households get offers from 2 shops, Plus from all their shops (features.md §6).

## Business rules

1. `SetMyShops` replaces the household's active list with the given shop ids, in the given order. Position 1 is the
   nearest shop as computed on the phone (README D13). Duplicate ids keep their first position.
2. Each id must be an `active` shop, or a `pending` shop added by this household. Otherwise `INVALID_ARGUMENT` naming
   the id.
3. All shops must be in one country, because prices are compared in one currency (PLN or EUR). Mixed countries:
   `INVALID_ARGUMENT`.
4. Limits by plan (plan from identity through `PlanFor`, OFF-01-03, Decided Q004): Free at most 2 shops; Plus at
   most 40 (Decided, Q161).
5. Free with more than 2 ids: `FAILED_PRECONDITION` with reason `plus_required` and limit `2` in the error details;
   nothing is stored.
6. Entries that became `inactive` after a downgrade (OFF-04-04) are kept untouched by a Free `SetMyShops`, unless the
   same shop is in the new active list.
7. Any household member with an account may change the list (all members share one fridge and list).
8. `ListMyShops` returns active entries in order, then inactive and closed entries, each with the shop, chain,
   coverage, editions, `state` (`active`, `inactive`, `closed`) and the position. planner calls it with the
   household id to know which shops to use and in what order.
9. An empty list is allowed (offers then shows nothing until shops are chosen).
10. Every successful change recomputes edition demand (OFF-04-03) in the same transaction.

## Contract

`contracts/coolinary/offers/v1/shops.proto`

```proto
rpc SetMyShops(SetMyShopsRequest) returns (SetMyShopsResponse);
rpc ListMyShops(ListMyShopsRequest) returns (ListMyShopsResponse);

message SetMyShopsRequest {
  string household_id = 1;
  repeated string shop_ids = 2;      // nearest first
}
message SetMyShopsResponse { ListMyShopsResponse my_shops = 1; }

message ListMyShopsRequest { string household_id = 1; }
message ListMyShopsResponse {
  repeated MyShop shops = 1;
  int32 limit = 2;                   // 2 on Free, 40 on Plus
  Attribution attribution = 3;
}
message MyShop {
  ShopWithEditions shop = 1;
  ChainCoverage coverage = 2;
  int32 position = 3;
  MyShopState state = 4;
}
enum MyShopState { MY_SHOP_STATE_UNSPECIFIED = 0; MY_SHOP_STATE_ACTIVE = 1; MY_SHOP_STATE_INACTIVE = 2; MY_SHOP_STATE_CLOSED = 3; }
```

## Data

`offers/migrations/000NN_my_shops.sql`

```sql
CREATE TABLE my_shops (
    household_id text NOT NULL,
    shop_id      text NOT NULL REFERENCES shops (id),
    position     int NOT NULL,
    state        text NOT NULL DEFAULT 'active' CHECK (state IN ('active', 'inactive', 'closed')),
    inactive_since timestamptz,
    updated_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (household_id, shop_id)
);
CREATE INDEX my_shops_shop ON my_shops (shop_id) WHERE state = 'active';
```

## Errors

Reason codes (in backticks) travel as trailing gRPC metadata `x-error-code` (README decision 18); version
conflicts use `ABORTED`.

| Case | gRPC status |
|---|---|
| Unknown or not visible shop id | `INVALID_ARGUMENT` |
| Shops in two countries | `INVALID_ARGUMENT` |
| Free and more than 2 shops | `FAILED_PRECONDITION` (`plus_required`) |
| Plus and more than 40 shops | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] A Free household sets [Lidl, Biedronka]; `ListMyShops` returns them in that order with `limit = 2`.
- [ ] The same household sending 3 shops gets `plus_required` and the stored list is unchanged.
- [ ] A Plus household stores 12 shops in order.
- [ ] Two household members setting lists one after the other: the later list wins.
