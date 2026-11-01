# OFF-04 My shops and shopping preferences

| | |
|---|---|
| Service | offers |
| Release | R2 |
| Goal | Store each household's confirmed shops (nearest first, as the phone sends them) and its "How you shop" settings, enforce the Free limit of 2 shops, and tell flyer sync which store editions someone actually needs. |

## Flow

1. The phone suggests shops nearest first from `ListShopsNearArea`; the user confirms; the phone sends the ordered list.
2. offers stores it, checks the plan (Free 2, Plus up to 40), and recomputes which editions are wanted.
3. "How you shop" settings are stored so every household device and the phone's trip plan use the same values.
4. Plan changes and shop closures adjust the list without losing the user's choices.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [OFF-04-01](OFF-04-01-my-shops.md) | My shops: set and list | R2 | M | OFF-03-03, OFF-03-04 |
| [OFF-04-02](OFF-04-02-shopping-preferences.md) | Shopping preferences ("How you shop") | R2 | S | — |
| [OFF-04-03](OFF-04-03-edition-demand.md) | Store-edition demand | R2 | S | OFF-04-01, OFF-03-04 |
| [OFF-04-04](OFF-04-04-plan-change-rules.md) | My shops on plan changes | R2 | S | OFF-04-01, OFF-01-03, OFF-01-05 |
| [OFF-04-05](OFF-04-05-closed-merged-shops.md) | Closed and merged shops in my shops | R2 | S | OFF-04-01, OFF-02-04 |
