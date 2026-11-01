# MOB-02-05 Entitlements and AI quota state

| | |
|---|---|
| Service | mobile |
| Epic | [MOB-02 API client and error model](README.md) |
| Release | R1 |
| Size | M |
| Depends on | MOB-01-04, MOB-02-02, GW-02 (`GET /v1/entitlements`), GW-03 (quota), IDN-05 |
| Unblocks | later paywall and quota UI |

## Goal

The app knows, before the user taps, whether a feature needs Plus and how many free AI day plans and photo scans
are left in the last 7 days, so it can show the Plus offer instead of a failed request.

## Business rules

1. `useEntitlements()` queries `GET /v1/entitlements` (GW-02-04), `staleTime` 5 min, persisted (it holds
   no personal data). It refetches on app foreground and whenever the JWT `plan` claim changes (MOB-01-04 rule 8).
2. The gate rules are the README table "Free and Plus"; the server is the source of truth and the app only mirrors
   them. `usePlusGate(feature)` returns `{ allowed: boolean; reason?: "plus_required" | "quota_reached"; remaining?: number; nextFreeAt?: string }`.
3. Quota counts update from three sources, newest wins:
   1. the entitlements response;
   2. the `quota` object and headers `X-Quota-Limit`, `X-Quota-Remaining`, `X-Quota-Next-Free-At` on successful AI
      responses (GW-03-02);
   3. the `quotaReached` error body (MOB-02-02), which sets `remaining = 0`.
4. **Rolling window (decided, Q111).** Free allows 2 AI day plans and 2 scans in any 7 days; each use frees up again
   7 days after it was made. When `remaining = 0`, the app shows "next free plan on <nextFreeAt>" (local date and time)
   and re-reads entitlements once `nextFreeAt` has passed, instead of waiting for a weekly reset.
5. Counts are never decremented on the phone before the server answers; a failed or fallback (non-AI) plan does not
   use a free plan (GW-03 counts a day plan only when planner reports `quota_consumed`).
6. `GET /v1/entitlements` returns `features` (named limits such as `aiDayPlansPerWeek`, `shopsWithOffers`,
   `tripPlan`) and `quota.dayPlan` / `quota.scan` (GW-02-04); the hook maps them into the `Entitlements` shape below.
7. `plusRequired` from any route invalidates the entitlements query once, in case the cached plan is stale.
8. After a plan change on the server (admin switch now, store purchase later), the app calls
   `tokenProvider.forceRefresh()` and refetches entitlements; Plus-gated queries (price matrix, per-person profiles,
   shared household) are invalidated. The JWT can be up to 15 min old otherwise.

## Contract

Assumed `GET /v1/entitlements` (align with GW-02/IDN-05):

```ts
interface Entitlements {
  plan: "free" | "plus";
  anonymous: boolean;
  features: Record<string, boolean>;   // e.g. "expiry_reminders", "price_matrix", "shared_household"
  quotas: { feature: "day_plan" | "photo_scan"; limit: number | null; used: number; nextFreeAt: string | null }[];   // null when nothing is used
  limits: { shops: number | null; accountsPerHousehold: number };
}
```

## Data

Query cache (persisted): entitlements.

## Errors

| Case | Behaviour |
|---|---|
| Entitlements unavailable offline | Use the persisted copy; if none, treat AI features as allowed and let the server decide |

## Acceptance criteria

- [ ] A Free household after two AI day plans shows `remaining = 0` without another request, and the third tap is
      blocked on the phone.
- [ ] Switching the household to Plus in admin and foregrounding the app unlocks Plus features within one refresh.
- [ ] A non-AI fallback plan leaves the remaining count unchanged.
- [ ] With both free plans used on Monday 10:00 and Wednesday 09:00, the app shows the next free plan on the
      following Monday at 10:00 and allows one plan again after that time without restarting.
