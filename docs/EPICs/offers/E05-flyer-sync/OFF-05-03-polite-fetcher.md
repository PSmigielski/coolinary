# OFF-05-03 Polite fetcher

| | |
|---|---|
| Service | offers |
| Epic | [OFF-05 Flyer source framework and sync scheduler](README.md) |
| Release | R2 |
| Size | M |
| Status | todo |
| Depends on | — |
| Unblocks | OFF-05-04, OFF-06-*, OFF-02-03 (Overpass uses it too) |

## Goal

Every request offers makes to a retailer's site is slow, identified, within the site's robot rules and limited to the
hosts that source allows (features.md §3 step 3: "downloaded once, slowly, with the app identified in the request and
the site's robot rules respected").

## Business rules

1. `User-Agent: CoolinaryBot/1.0 (+https://<DOMAIN>/bot; bot@<DOMAIN>)`, with `DOMAIN` from the root `.env`. The
   `/bot` page (what the bot reads, how to contact us, how to opt out) is a frontend/web task; until it exists the
   URL is the landing page.
2. `robots.txt` of each host is read before the first request and cached for 24 h. A disallowed path is never
   requested; the call fails with `ErrRobots`, which marks the source failing with reason `robots` (OFF-05-06).
   A missing `robots.txt` (404) allows everything; a 5xx on `robots.txt` blocks the host for that run.
3. Rate: at most one request at a time per host and at least 3 s between requests to the same host, or the
   `Crawl-delay` in `robots.txt` if larger (Decided, Q151: 3 s).
4. At most 4 hosts fetched in parallel across the service.
5. Only `https`. Redirects: at most 5, and only to hosts in the source's `allowed_hosts`.
6. Timeouts: 10 s to connect, 60 s per response. Size caps: 5 MB for HTML/JSON, 10 MB per image, 100 MB per PDF.
7. Retries: network errors and 5xx retried twice after 30 s and 120 s; 429 waits `Retry-After` (at most 1 h, else
   failure); 401/403 are not retried and fail with `ErrBlocked`.
8. List pages are requested with `If-None-Match` / `If-Modified-Since` when the last response had `ETag` /
   `Last-Modified`; a 304 counts as success with "no change".
9. No log-in, no stored cookies across flyers, no JavaScript execution. A chain whose flyers can only be read by
   running scripts stays `user_flyers_only`.
10. Metrics: `fetch_requests_total{host,code}`, `fetch_bytes_total{host}`, `fetch_wait_seconds` (time spent waiting
    for the rate limit).

## Contract

No gRPC change. Go interface `Fetcher` with `Get(ctx, url, opts) (Response, error)`, scoped to a source (its allowed
hosts).

## Data

`offers/migrations/000NN_fetch_cache.sql`

```sql
CREATE TABLE fetch_validators (
    url           text PRIMARY KEY,
    etag          text,
    last_modified text,
    checked_at    timestamptz NOT NULL DEFAULT now()
);
```

## Errors

| Case | Behaviour |
|---|---|
| Disallowed by `robots.txt` | `ErrRobots` |
| 401 / 403 | `ErrBlocked` |
| 429 with `Retry-After` over 1 h | `ErrRateLimited` |
| Redirect to a host not allowed | `ErrHostNotAllowed` |
| Body over the size cap | `ErrTooLarge` |

## Acceptance criteria

- [ ] Against a local test server: two requests to the same host are at least 3 s apart; `Crawl-delay: 10` makes it
      10 s.
- [ ] A path disallowed in `robots.txt` is never requested (server log shows no hit) and returns `ErrRobots`.
- [ ] A redirect to a host not in `allowed_hosts` is not followed.
- [ ] The `User-Agent` header is present on every request, including `robots.txt`.
