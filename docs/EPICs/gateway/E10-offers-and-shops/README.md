# GW-10 Shops, offers, flyers and price matrix routes

Releases R2 (shops, my shops, offers, own flyers, price matrix) and R3 (cheapest top-up).

## Goal

The app finds shops around a postcode, saves the household's shops and "how you shop" settings, browses this and
next week's offers, scans a flyer that is not covered (Plus), and gets the price matrix the phone uses to plan the
shopping trip (decision 7, Plus). The server never sees an exact position or a distance (D13): the phone sorts shops
and sends their ids nearest first.

Near-shop reminders (R3) run entirely on the phone from the price matrix and the saved shops; they need no route.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-10-01](GW-10-01-shops-by-postcode.md) | Shops around a postcode; suggest a shop change | R2 | S | GW-01-02, OFF-02, OFF-03 |
| [GW-10-02](GW-10-02-my-shops.md) | My shops and shopping preferences | R2 | S | GW-10-01, OFF-04 |
| [GW-10-03](GW-10-03-offers-browse.md) | Offer browse, search and flyer status | R2 | M | GW-10-02, OFF-09, OFF-05 |
| [GW-10-04](GW-10-04-own-flyer.md) | Scan a flyer that is not covered (async) | R2 | M | GW-01-05, GW-03-03, OFF-10 |
| [GW-10-05](GW-10-05-price-matrix.md) | Trip price matrix | R2 | M | GW-08-01, GW-03-03, PLN-19 |
| [GW-10-06](GW-10-06-cheapest-top-up.md) | Cheapest top-up | R3 | S | GW-10-05, GW-07-11, PLN-20 |
