# OFF-03-03 Shops near an area

| | |
|---|---|
| Service | offers |
| Epic | [OFF-03 Postcode areas and flyer editions](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | OFF-03-02, OFF-02-03, OFF-02-04, OFF-03-04 |
| Unblocks | OFF-04-01, GW-10, MOB-03 |

## Goal

The phone gets every shop it might need for the user's travel range, with exact shop positions, and computes all
distances itself (features.md §3 "Distance is calculated on the phone", README D13).

## Business rules

1. Input: an area id (OFF-03-02) or a country plus postcode.
2. offers selects active shops whose position lies in a box of ±15 km (north-south and east-west) around the area's
   centre. 15 km = the largest travel range (10 km, "by car") plus 5 km for the gap between the user and their
   postcode's centre. The box is computed from the centre in degrees (15 / 111.32 for latitude, 15 / (111.32 ×
   cos(latitude)) for longitude); no distance to the user is ever calculated on the server.
3. Also included: the caller household's own pending shops anywhere in the box (OFF-02-05), flagged
   `household_only`. Shops with `reported_closed = true` are left out.
4. Every shop is returned with its chain, chain coverage, position, address, and its edition ids (OFF-03-04).
5. Order: by chain id, then shop id (stable; the phone sorts by distance).
6. At most 1,000 shops. If the box holds more, the response sets `truncated = true` and keeps shops of chains with
   coverage `live` or `feed` first, then the rest by id. Decided (Q149): the cap is 1,000 shops (central Berlin holds
   about 400 shops of the seeded chains).
7. Response is cached per (country, postcode or town id, OSM import version) in `redis-cache` for 24 h; the
   household's pending shops are added after the cache read.
8. Chains in the response carry coverage, so the app can say "offers not covered yet, scan the flyer yourself".
9. The response carries both attributions (OSM, OFF-02-08; GeoNames, OFF-03-01).

## Contract

```proto
rpc ListShopsNearArea(ListShopsNearAreaRequest) returns (ListShopsNearAreaResponse);

message ListShopsNearAreaRequest {
  string household_id = 1;
  string area_id = 2;          // or country_code + postcode
  string country_code = 3;
  string postcode = 4;
}
message ListShopsNearAreaResponse {
  repeated ShopWithEditions shops = 1;
  repeated Chain chains = 2;
  bool truncated = 3;
  repeated Attribution attributions = 4;
}
message ShopWithEditions {
  Shop shop = 1;
  repeated string edition_ids = 2;
}
```

## Data

No schema change (index `shops_position` from OFF-02-02).

## Errors

| Case | gRPC status |
|---|---|
| Unknown area | `NOT_FOUND` |
| Neither area nor postcode | `INVALID_ARGUMENT` |

## Acceptance criteria

- [ ] For a fixed test postcode, the list contains every seeded-chain shop within 10 km of any point of that postcode
      area (checked against the fixture).
- [ ] The response has no field derived from the user's position; the request has none either.
- [ ] A test position in Hamburg returns Aldi Nord and no Aldi Süd shops; Munich the reverse (Verification 8).
- [ ] Second call for the same postcode is served from cache (metric `cache_hits_total{cache="shops_near"}`).
