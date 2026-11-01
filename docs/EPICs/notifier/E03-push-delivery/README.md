# NTF-03 Push delivery (Expo Push API)

**Goal.** kitchen, planner and identity send pushes (expiry reminders in R1; offer ends, household changes and the
savings report in R2) with one call. notifier applies preferences, quiet hours and deduplication, sends through
Expo in batches, retries, and cleans up dead tokens and message texts.

**Flow.** `SendPush(users | household, category, template, variables, path)` → per user: preferences, dedup,
devices → one queue row per device (rendered in the device language, deferred by quiet hours) → worker sends ≤ 100
per Expo request → tickets → receipts after 15 minutes → final state, texts cleared.

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [NTF-03-01](NTF-03-01-push-templates.md) | Push templates and categories | R1 (R2 templates) | M | NTF-02-01 |
| [NTF-03-02](NTF-03-02-send-push.md) | SendPush and the push queue | R1 | L | NTF-01-02, NTF-02-02, NTF-03-01, NTF-04-01, NTF-04-02 |
| [NTF-03-03](NTF-03-03-expo-delivery-worker.md) | Expo delivery worker | R1 | M | NTF-03-02 |
| [NTF-03-04](NTF-03-04-push-receipts.md) | Push receipts and cleanup | R1 | S | NTF-03-03 |
