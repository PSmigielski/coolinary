# NTF-05-02 Founder alert e-mail and recipients

| | |
|---|---|
| Service | notifier |
| Epic | [NTF-05 Founder alerts](README.md) |
| Release | R1 |
| Size | S |
| Depends on | NTF-02-01 |
| Unblocks | NTF-05-01 |

## Goal

Alert e-mails are short, scannable on a phone, and say what broke, since when, and where to look.

## Business rules

1. Template id `founder_alert`, channel `email`, English only (internal; the en text is used for any language).
   The "every language" test of NTF-02-01 rule 8 exempts templates flagged `internal`.
2. Variables: `severity`, `source`, `summary`, `details`, `link`, `first_at`, `count`, `state` (`open`, `repeat`,
   `resolved`).
3. Subject: `[<SEVERITY>] <source>: <summary>`, `[RESOLVED] …` for resolved, cut to 120 characters.
4. Body (text only; HTML is not needed): summary, details, "First seen <first_at> UTC", "Seen <count> times", the
   link when present, and "Acknowledge in admin: <ADMIN_PUBLIC_URL>/alerts" when `ADMIN_PUBLIC_URL` is set.
5. Config: `FOUNDER_ALERT_EMAILS` (comma list, ≤ 5), `ADMIN_PUBLIC_URL` (optional). Dev: one test address that lands
   in Mailpit.

## Contract

No gRPC change (used by NTF-05-01).

## Data

New config in `notifier/.env.example`, `.env.dev`, `.env.prod`: `FOUNDER_ALERT_EMAILS`, `ADMIN_PUBLIC_URL`.

## Errors

As NTF-02-01.

## Acceptance criteria

- [ ] Snapshot for `open`, `repeat` and `resolved`.
- [ ] A 300-character summary gives a 120-character subject.
