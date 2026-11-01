# GW-08 Shopping list routes

Release R1 (offer badges appear on lines from R2, PLN-17; live sharing is GW-09).

## Goal

The household's one shopping list (features.md §2 "Shopping list"): proposed from the plan with what is at home
subtracted, reviewed before it joins the list, edited by hand, ticked in the shop (ticked items move into the
fridge), plus cuisine base-pantry items. Planner owns the list; every write bumps the list `version`, which the live
socket (GW-09) and the offline queue (MOB-04) rely on.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [GW-08-01](GW-08-01-list-lines.md) | Read the list; add, edit and delete lines | R1 | M | GW-01-03, PLN-15 |
| [GW-08-02](GW-08-02-proposals.md) | Propose lines from the plan and accept them | R1 | M | GW-08-01, GW-07-05, PLN-15 |
| [GW-08-03](GW-08-03-tick.md) | Tick in the shop → into the fridge | R1 | S | GW-08-01, PLN-15 |
| [GW-08-04](GW-08-04-base-pantry.md) | Cuisine base pantry | R1 | S | GW-08-01, PLN-16 |
