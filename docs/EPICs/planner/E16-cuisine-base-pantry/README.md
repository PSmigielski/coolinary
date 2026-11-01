# PLN-16 Cuisine base pantry

Release R1. Each cuisine has a "base pantry" checklist (Asian: soy sauce, rice, ginger, sesame oil; Italian: olive
oil, passata, pasta, hard cheese). One tap adds the missing base items to the shopping list (features.md §2
"Cuisine choice", "Shopping list" step 4). No AI.

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [PLN-16-01](PLN-16-01-base-pantry-lists.md) | Base pantry lists and admin edits | R1 | S | PLN-01, KIT-02 (GetIngredients) |
| [PLN-16-02](PLN-16-02-base-pantry-checklist.md) | Base pantry checklist for a household | R1 | S | PLN-16-01, KIT-07 (ListStock), KIT-10 (GetPlanningProfile) |
| [PLN-16-03](PLN-16-03-add-missing-base-items.md) | Add missing base items to the list | R1 | S | PLN-16-02, PLN-15-04, PLN-15-08 |
