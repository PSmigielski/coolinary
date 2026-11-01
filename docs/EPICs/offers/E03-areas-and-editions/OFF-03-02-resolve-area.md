# OFF-03-02 Resolve a postcode or town to an area

| | |
|---|---|
| Service | offers |
| Epic | [OFF-03 Postcode areas and flyer editions](README.md) |
| Release | R2 |
| Size | S |
| Status | todo |
| Depends on | OFF-03-01 |
| Unblocks | OFF-03-03, GW-10 |

## Goal

A user who refuses location, or who wants to check another place, types a postcode or a town and gets an area that
the shop list can use (features.md §3 "User position … a typed postcode or town works instead").

## Business rules

1. Input: country code and a query. A query that is a postcode in the country's format (PL `NN-NNN` or `NNNNN`, DE
   `NNNNN`) is a postcode lookup; anything else is a town search.
2. Postcode lookup returns exactly one area (`kind = postcode`) or `NOT_FOUND`.
3. Town search: compare the query (lower-cased, without diacritics, trimmed, at least 2 characters) with town search
   keys: exact match first, then prefix matches, then trigram similarity ≥ 0.4. Return at most 10 areas
   (`kind = town`), larger towns (more postcodes) first within each group.
4. Area ids are opaque strings: `pc:PL:00-001`, `town:DE:1234`.
5. An area's `centre` is returned, because the phone shows it on a map and uses it as the stand-in position for
   distances when there is no location permission (Verification 10: same ordering from the postcode).
6. The phone that has location permission reverse-geocodes on the device and sends only the postcode (README D13);
   offers never receives coordinates of the user.
7. Queries are not logged (they can reveal where someone lives).

## Contract

`contracts/coolinary/offers/v1/shops.proto`

```proto
rpc ResolveArea(ResolveAreaRequest) returns (ResolveAreaResponse);

message ResolveAreaRequest {
  string country_code = 1;
  string query = 2;
}
message ResolveAreaResponse {
  repeated Area areas = 1;
  Attribution attribution = 2;     // GeoNames credit
}
message Area {
  string id = 1;
  AreaKind kind = 2;
  string label = 3;                // "00-001 Warszawa", "München (Bayern)"
  double centre_latitude = 4;
  double centre_longitude = 5;
}
enum AreaKind { AREA_KIND_UNSPECIFIED = 0; AREA_KIND_POSTCODE = 1; AREA_KIND_TOWN = 2; }
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Country not PL or DE | `INVALID_ARGUMENT` "country not covered yet" |
| Query shorter than 2 characters | `INVALID_ARGUMENT` |
| Postcode unknown | `NOT_FOUND` |

## Acceptance criteria

- [ ] `PL` + `00001` and `PL` + `00-001` return the same area.
- [ ] `DE` + `munch` returns München first.
- [ ] Logs of the call contain the request id and duration but not the query.
