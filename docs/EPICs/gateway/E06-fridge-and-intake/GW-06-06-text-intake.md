# GW-06-06 Typed or spoken sentence → proposals

| | |
|---|---|
| Service | gateway |
| Epic | [GW-06 Fridge, barcode and intake routes](README.md) |
| Release | R3 |
| Size | S |
| Status | todo |
| Depends on | GW-03-01, GW-01-03, KIT-08 (`IntakeFromText`) |
| Unblocks | MOB-03 |

## Goal

"Bought eggs, two peppers, some cheese" — typed, or spoken and turned into text by the phone (D7) — becomes items
the user confirms (features.md §2, release 3 "voice input").

## Business rules

1. JSON body `{ "text", "kind": "text" | "voice", "storage"? }`; `text` trimmed, 1–2,000 characters (matches
   `ai`'s limit). `voice` means the phone's speech recognition produced the text (D7).
2. Counts as one `scan` for Free (decided, Q112: every AI reading request is a scan). Ships in R3 with voice
   (decided, Q219). Same order as GW-06-04, idempotency required.
3. Calls kitchen `IntakeFromText(household_id, text, kind, storage, language)`; returns kitchen's `IntakeResponse`
   as in GW-06-04 (`intakeRef`, proposals), plus `rulesFallback`.
4. The text is user free text: never logged, never stored by the gateway (only hashed inside the idempotency
   fingerprint).
5. AI unavailable, over budget or refusing: kitchen answers with its own parser (`ai_used = false`,
   `rules_fallback = true`, KIT-08-03 rule 3). The gateway returns 200 with those proposals and releases the quota.

## Contract

`POST /v1/intake/text` — `auth: true`, tag `intake`, idempotency required, quota `scan`.
Request `{ "text": "kupiłam jajka, dwie papryki i trochę sera", "storage": "fridge" }` → 200 same shape as GW-06-04.

## Data

No storage.

## Errors

| Case | HTTP | `code` |
|---|---|---|
| Empty or > 2,000 chars | 400 | `validation` |
| Free quota used up | 429 | `quota_reached` |
| Kitchen down | 503 | `service_unavailable` |

## Acceptance criteria

- [ ] The Polish sentence above returns proposals for eggs, 2 bell peppers and cheese.
- [ ] Request logs contain no part of the text (log-scanning test).
- [ ] A 2,001-character text gets 400 without calling kitchen.
- [ ] With `ai` stopped, the sentence still returns proposals with `rulesFallback: true` and the Free scan quota
      is unchanged.
