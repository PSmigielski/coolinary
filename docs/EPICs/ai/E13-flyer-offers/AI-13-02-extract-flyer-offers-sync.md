# AI-13-02 ExtractFlyerOffers for user flyers (sync, chunked)

| | |
|---|---|
| Service | ai |
| Epic | [AI-13 Flyer offer extraction](README.md) |
| Release | R2 |
| Size | L |
| Status | todo |
| Depends on | AI-13-01, AI-01-08 |
| Unblocks | OFF-10 (user-supplied flyers) |

## Goal

A user photographs or shares a flyer that Coolinary doesn't cover yet and gets its offers within about a minute, in
any country (features.md §2, "Photograph or share any flyer…"; Plus).

## Business rules

1. Job `extract_flyer_offers`, Haiku 5.5 at effort `medium`. Prompt `extract_flyer_offers/v1` (shared with batch).
2. Input: either 1–40 page images (each ≤ 3 MiB, JPEG/PNG/WebP) or one PDF (≤ 25 MiB, ≤ 40 pages; Decided, Q022: user flyers of up to 40 pages from R2); `FlyerMeta`
   with what the caller knows (country always; chain, currency and dates often empty for user flyers).
3. Chunking: ai splits the pages into chunks of at most 10 pages (PDFs split with `pypdf`, added as a dependency)
   so each Claude call stays under the 100K-token prompt line where Haiku's price rises 5× (AI-02-01). Chunks run
   at most 3 at a time, each with the job timeout (40 s); a chunk that fails is retried once, then reported in
   `failed_pages`.
4. The first chunk also returns `flyer_info`: chain name as printed, store or region as printed, validity dates as
   printed. Offers without their own dates inherit nothing in ai; the caller applies `flyer_info` dates.
5. Merge: offers from all chunks in page order; exact duplicates (same name, price, pack size on the same page)
   removed.
6. Response time target: p95 ≤ 60 s for 20 pages; the caller sets a 120 s deadline. If the deadline would pass, ai
   returns what it has with `partial=true` and the missing pages in `failed_pages`.
7. Images and PDFs are never stored (D6); cache per page by SHA-256 of the page bytes, 30 days, so the same flyer
   shared by two users costs once.
8. ai never fetches a URL; for shared links, offers downloads the flyer and sends the bytes (OFF-10).

## Contract

```proto
rpc ExtractFlyerOffers(ExtractFlyerOffersRequest) returns (ExtractFlyerOffersResponse);

message ExtractFlyerOffersRequest {
  repeated InputImage pages = 1;   // 1-40, or
  bytes pdf = 2;                   // one PDF, ≤ 40 pages
  FlyerMeta meta = 3;
  string household_id = 4;
}

message FlyerInfo {
  string chain_name_printed = 1;
  string store_or_region_printed = 2;
  string valid_from = 3;
  string valid_to = 4;
}

message ExtractFlyerOffersResponse {
  repeated FlyerPageOffers pages = 1;
  FlyerInfo flyer_info = 2;
  repeated int32 failed_pages = 3;
  bool partial = 4;
  CallInfo call = 5;               // summed over chunks
}
```

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Both or neither of `pages` and `pdf`; over 40 pages; size or type invalid | `INVALID_ARGUMENT`, `ai_bad_request` |
| Every chunk failed | the first chunk's error (AI-01-04) |

## Acceptance criteria

- [ ] A 24-page PDF is split into 3 calls of ≤ 10 pages; each call's prompt is under 100K tokens (from usage).
- [ ] A 40-page, 24 MiB PDF is accepted and split into 4 calls of 10 pages; a 41-page PDF or a 26 MiB PDF is
      refused with `INVALID_ARGUMENT` (Q022).
- [ ] The same PDF sent twice costs once.
- [ ] A page that fails twice appears in `failed_pages`; the rest is returned.
- [ ] Cost of a 20-page flyer ≈ $0.02 (report).
