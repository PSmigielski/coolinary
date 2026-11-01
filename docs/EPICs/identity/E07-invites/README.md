# IDN-07 Invites (QR, code, link)

features.md "Linking people into a household": an invite can be handed over as a QR code on screen, a 6-character
code, or a link sent through any messaging app. It works for 48 hours and for one person. Waiting invites are
listed and can be cancelled. The person joining sees whose household it is before confirming and chooses whether
to bring their fridge items. Sharing needs Plus, held by the owner; a household has at most 5 accounts.

One invite = one code. The QR code and the link both carry that code.

## Flow

1. A member creates an invite → code `KX7-M4P`, link `https://<invite host>/i/KX7M4P`, QR of the link.
2. The other person opens the link, scans the QR or types the code → preview: household and owner name, whether
   they can bring their fridge.
3. They confirm → they join; their old personal household is deleted (its fridge moved first if they chose so).

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| IDN-07-01 | Invite table and code generator | R2 | S | IDN-04-01 |
| IDN-07-02 | Create, list and cancel invites | R2 | M | IDN-07-01, IDN-05-01, IDN-04-05 |
| IDN-07-03 | Invite preview | R2 | S | IDN-07-02 |
| IDN-07-04 | Accept an invite | R2 | L | IDN-07-03, IDN-01-01, IDN-04-02 |
| IDN-07-05 | Invite expiry and cleanup | R2 | S | IDN-07-01 |
