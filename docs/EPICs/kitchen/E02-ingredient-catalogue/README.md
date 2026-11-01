# KIT-02 Ingredient catalogue and translations

| | |
|---|---|
| Service | kitchen |
| Release | R1 |
| Unblocks | KIT-03 matching, KIT-04/05 barcodes, KIT-06 nutrition, KIT-07 stock, PLN-02 recipes, GW-05 autocomplete |

## Goal

One shared list of ingredients that every part of Coolinary points to: recipes (planner), fridge items, scanned
products, flyer offers and nutrition values. Each ingredient knows its category, the 14 EU allergens it contains,
diet flags, how much a piece or spoon weighs, where it keeps and for how long, typical pack sizes, and its names in
every launch language (Polish, German, English). Names come from a reviewed seed plus an AI translation batch
(once for everyone, features.md §2 "Fridge and pantry, With AI").

## Flow

```
seed CSV ──seed command──> ingredients, names, measures, packs
translation batch (AI-09 translate_ingredients) ──> machine names (reviewed later in ADM-03)
app typing ──SearchIngredients──> autocomplete      planner/offers ──GetIngredients──> full records
admin ──KitchenAdminService──> edits, merges ──> catalogue.ingredient_merged / ingredient_flags_changed
```

## Tasks

| ID | Title | Rel. | Size | Depends on |
|---|---|---|---|---|
| [KIT-02-01](KIT-02-01-ingredient-core-schema.md) | Ingredient core: categories, allergens, diet flags | R1 | M | — |
| [KIT-02-02](KIT-02-02-measures-and-units.md) | Typical weights, density and unit conversion | R1 | M | KIT-02-01 |
| [KIT-02-03](KIT-02-03-storage-shelf-life-packs.md) | Storage place, shelf life and pack sizes | R1 | S | KIT-02-01 |
| [KIT-02-04](KIT-02-04-names-and-synonyms.md) | Names, synonyms and plurals per language | R1 | M | KIT-02-01, KIT-03-01 |
| [KIT-02-05](KIT-02-05-search-ingredients.md) | SearchIngredients (autocomplete) | R1 | M | KIT-02-04 |
| [KIT-02-06](KIT-02-06-get-ingredients.md) | GetIngredients and catalogue version | R1 | M | KIT-02-02, KIT-02-03, KIT-02-04 |
| [KIT-02-07](KIT-02-07-starter-catalogue-seed.md) | Starter catalogue seed and seed command | R1 | M | KIT-02-01 … KIT-02-04 |
| [KIT-02-08](KIT-02-08-ai-translation-batch.md) | AI translation batch job | R1 | M | KIT-02-04, KIT-01-08, AI-09 |
| [KIT-02-09](KIT-02-09-merge-ingredients.md) | Merge duplicate ingredients | R1 | M | KIT-02-04, KIT-01-02, KIT-01-07, KIT-03-02 |
| [KIT-02-10](KIT-02-10-admin-catalogue-rpcs.md) | Admin catalogue editing RPCs | R1 | M | KIT-01-07, KIT-02-01 … KIT-02-04 |

## Decisions (taken 1 November 2026)

- Q050 (KIT-02-02): No default; the dish shows "calories incomplete" until a density is added.
- Q051 (KIT-02-02): Pinch = 0.5 g; "to taste" = 0 g, shown as "not counted".
- Q053 (KIT-02-07): About 700 ingredients, AI-drafted, checked by the founder and a native speaker per language.
- Q054 (KIT-02-10): Yes: new event catalogue.ingredient_flags_changed.
- Q055 (KIT-02-10): Yes: hidden from search and new matches; existing references keep working.
- Q056 (KIT-02-09): No undo; create a new ingredient if a merge was wrong.
- Q066 (KIT-02-08): Yes, live at once; the native-speaker check per language follows.
