# OFF-02-03 OSM shop import job

| | |
|---|---|
| Service | offers |
| Epic | [OFF-02 Chains and shops](README.md) |
| Release | R2 |
| Size | L |
| Status | todo |
| Depends on | OFF-02-02, OFF-01-05 |
| Unblocks | OFF-02-04, OFF-03-03, OFF-03-05 |

## Goal

Every shop of every seeded chain in Poland and Germany is known with its position, refreshed monthly from
OpenStreetMap, so nothing per country has to be negotiated (features.md §3 "Location").

## Business rules

1. Source (Decided, Q144): the Overpass API, one query per chain per country, monthly, selecting nodes and ways with
   the chain's `brand:wikidata` (or, if none matched, `brand` name) and one of its `osm_shop_values`, inside the
   country's area. Ways are reduced to their centre point (`out center`).
2. Overpass etiquette: one query at a time, at least 30 s between queries, `User-Agent` and `Referer` identifying
   Coolinary with a contact address, timeout 180 s, honour HTTP 429 by waiting the `Retry-After` time. Base URL is a
   setting (`OVERPASS_URL`) so a paid or self-hosted instance can replace the public one.
3. Schedule: monthly, on the 3rd at 03:00 `Europe/Berlin`; manual run through `TriggerJob`.
4. For each returned element, upsert by `(osm_type, osm_id)`: chain, name, street (`addr:street` + `addr:housenumber`),
   postcode (`addr:postcode`), city (`addr:city`), latitude, longitude. Set `missing_refreshes = 0`.
5. An OSM shop already in the table but absent from this run's result for its chain gets `missing_refreshes + 1`.
   At 2 it becomes `closed`. A closed shop that reappears becomes `active` again.
6. A chain query that fails (HTTP error, timeout, empty result where the last run had more than 20 shops) changes
   nothing for that chain: no shop is counted as missing. A founder alert is sent (`SendAlert`, key
   `offers.osm.<chain_id>`).
7. A run that would close more than 10% of a chain's active shops at once stops for that chain and alerts instead
   (protects against a broken query).
8. User-added shops (`source = 'user'`) are never closed or changed by the import.
9. Shops whose status is `merged` or `rejected` are not reactivated by the import.
10. After the run, OFF-02-04 merges duplicates and OFF-03-05 links new shops to chain store numbers.
11. Metrics: `osm_import_shops{country,chain,status}` gauge, `osm_import_runs_total{chain,result}`.

## Contract

No public RPC. Admin: the job appears in `ListJobs`; `TriggerJob` with name `osm_import` (all chains) or
`osm_import:<chain_id>`.

## Data

`offers/migrations/000NN_osm_import.sql`

```sql
CREATE TABLE osm_import_runs (
    id          bigserial PRIMARY KEY,
    chain_id    text NOT NULL REFERENCES chains (id),
    started_at  timestamptz NOT NULL DEFAULT now(),
    finished_at timestamptz,
    result      text CHECK (result IN ('ok', 'failed', 'stopped_mass_close')),
    seen        int NOT NULL DEFAULT 0,
    added       int NOT NULL DEFAULT 0,
    closed      int NOT NULL DEFAULT 0,
    error       text
);
```

## Errors

| Case | Behaviour |
|---|---|
| Overpass 429 / timeout | Retry the chain twice (wait `Retry-After` or 5 min); then `failed`, alert |
| Mass close (> 10%) | Chain run `stopped_mass_close`, alert, no changes |

## Acceptance criteria

- [ ] With a recorded Overpass response for Lidl DE (fixture), the import creates one shop per element with
      positions to 6 decimals and the right postcodes.
- [ ] Running the same fixture twice changes nothing; a fixture with one shop removed marks it missing once, and a
      second run closes it.
- [ ] A fixture with 30% of shops removed stops with `stopped_mass_close` and sends one alert (Mailpit).
- [ ] Test positions in Hamburg and Munich (Verification 8) have Aldi Nord and Aldi Süd shops respectively within
      5 km after a real import in dev.

## Notes

The "© OpenStreetMap contributors" credit and the share-alike question are handled in OFF-02-08.
