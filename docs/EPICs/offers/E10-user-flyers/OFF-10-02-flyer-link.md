# OFF-10-02 Flyer from a link

| | |
|---|---|
| Service | offers |
| Epic | [OFF-10 User-supplied flyers](README.md) |
| Release | R2 |
| Size | M |
| Depends on | OFF-10-01 |
| Unblocks | GW-10 |

## Goal

A user can paste or share a link to a flyer file instead of photographing it (features.md §3 "photo, PDF, or link").

## Business rules

1. Same plan rule and limits as OFF-10-01 (a link counts as one upload).
2. Decided (Q172): in R2 only direct links to a PDF or an image are accepted; a link to an HTML page
   or an online flipbook is refused with `link_not_a_file` (the app suggests sharing the PDF or taking photos).
   Alternative: adapters for common flyer viewers.
3. Safe fetching (the server fetches a URL chosen by a user):
   - only `https`;
   - the host name is resolved and every resulting address must be public: loopback, private (10/8, 172.16/12,
     192.168/16), link-local, carrier-grade NAT, multicast and IPv6 equivalents are refused; the connection goes to the
     checked address (no second lookup);
   - at most 3 redirects, each checked the same way;
   - timeout 30 s; at most 20 MB; `Content-Type` must be `application/pdf`, `image/jpeg`, `image/png` or `image/webp`,
     and the first bytes must match the type;
   - the polite fetcher's `User-Agent` is sent and the host's `robots.txt` is honoured.
4. After download the job continues exactly as an upload (OFF-10-03). The URL itself is not stored (it can contain
   personal tokens); only the host name is kept, for abuse checks.

## Contract

```proto
rpc SubmitUserFlyerLink(SubmitUserFlyerLinkRequest) returns (UserFlyer);
message SubmitUserFlyerLinkRequest {
  string household_id = 1;
  string user_id = 2;
  string shop_id = 3;
  string url = 4;
}
```

## Data

`ALTER TABLE user_flyers ADD COLUMN link_host text;`

## Errors

| Case | Result |
|---|---|
| `http://`, private address, too many redirects | `INVALID_ARGUMENT` (`link_not_allowed`) |
| HTML page | job `failed` (`link_not_a_file`) |
| Download fails or too large | job `failed` (`link_failed`, `too_large`) |

## Acceptance criteria

- [ ] `https://127.0.0.1/x.pdf`, `https://10.77.0.30/`, and a public host that redirects to a private address are all
      refused without any request reaching the private address.
- [ ] A PDF link to a test server (address check pointed at a test resolver) produces a draft.
- [ ] The stored row contains the host but not the path or query.
