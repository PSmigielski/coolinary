# ADM-03-03 Ingredient editor and translation review

| | |
|---|---|
| Service | admin |
| Epic | [ADM-03 Catalogue, matching and nutrition review](README.md) |
| Release | R1 |
| Size | M |
| Depends on | ADM-02-03; KIT-02-05/06 (`SearchIngredients`, `GetIngredients`), KIT-02-10 (`ListIngredients`, `CreateIngredient`, `UpdateIngredient`, `SetIngredientNames`, `ApproveMachineNames`, `SetMeasures`, `SetPackSizes`, `RetireIngredient`) |
| Unblocks | ADM-03-01 (`new_ingredient`) |

## Goal

The founder can create and correct catalogue ingredients (names per language, typical weights, shelf life, allergen
and diet flags, pack sizes) and approve AI translations, so matching, calories, expiry dates and allergen filters
rest on correct data.

## Business rules

1. Search by name in any language (`KitchenService.SearchIngredients`, 20 results), or browse with
   `KitchenAdminService.ListIngredients` and its filters: text, category, status, `missing_language`,
   `missing_measures`, `missing_nutrition`, `unreviewed_names` (KIT-02-10 rule 7).
2. The editor loads one ingredient with `KitchenService.GetIngredients([id])` and shows every field kitchen's model
   has (KIT-02 owns the list). Saving calls the matching admin RPC:
   - category, base unit, allergens, diet flags, staple candidate, density, storage, shelf lives →
     `UpdateIngredient` with a field mask of the changed fields;
   - names of one language → `SetIngredientNames` (full replacement of that language's names);
   - typical weights → `SetMeasures`; pack sizes → `SetPackSizes`;
   - retire → `RetireIngredient` (hidden from search and new matches, references keep working; KIT-02-10 rule 6).
3. Two admins editing the same ingredient: the later save wins per field; kitchen's `admin_changes` (shown with
   `ListAdminChanges`, KIT-01-07) and the admin audit log keep both.
4. Create (`CreateIngredient`) needs the fields KIT-02-01 … 04 require; kitchen answers `INVALID_ARGUMENT` naming
   the field. After saving, admin shows a checklist from the same checks as `ListIngredients`' filters: no primary
   name in pl, de or en; no piece weight for a `PIECE` ingredient; no density for a `MILLILITRE` ingredient; no
   approved nutrition link (ADM-03-02).
5. Allergen flags (the 14 EU allergens) and diet flags are hard limits for users (features.md §1), so changing them
   asks for confirmation in the UI and the audit row lists old and new values. Kitchen tells planner through
   `catalogue.ingredient_flags_changed` at once (KIT-02-10 rule 5, Q054).
6. Translation review: machine translations from the AI batch (KIT-02-08) are listed with
   `ListIngredients(unreviewed_names = true, language)`. Actions: `accept` → `ApproveMachineNames(name ids)`;
   `edit` → `SetIngredientNames` with the corrected text (marked reviewed, KIT-02-10 rule 3); `reject` →
   `SetIngredientNames` without the machine name (the language stays untranslated and the daily job may try again,
   KIT-02-08). Bulk accept up to 100 name ids in one `ApproveMachineNames` call.
7. Audited: `kitchen.ingredient.create`, `kitchen.ingredient.update` (changed field names with new values),
   `kitchen.ingredient.retire`, `kitchen.translation.resolve`.
8. Role `admin` only.

## Contract

| Method | Path | Body / response |
|---|---|---|
| GET | `/admin/v1/kitchen/ingredients?q=&language=&category=&status=&missing=&page_token=` | `{ items, next_page_token }` |
| GET | `/admin/v1/kitchen/ingredients/{id}` | ingredient with its checklist |
| POST | `/admin/v1/kitchen/ingredients` | ingredient fields → created ingredient, `checklist` |
| PATCH | `/admin/v1/kitchen/ingredients/{id}` | `{ ...changed fields }` → ingredient (one kitchen call per changed group, rule 2) |
| PUT | `/admin/v1/kitchen/ingredients/{id}/names/{language}` | `{ names[] }` → ingredient |
| POST | `/admin/v1/kitchen/ingredients/{id}/retire` | `{ reason }` → ingredient |
| GET | `/admin/v1/kitchen/translation-reviews?language=&page_token=` | `{ items: [{ ingredient_id, language, name_id, text, source_text }] }` |
| POST | `/admin/v1/kitchen/translation-reviews/resolve` | `{ items: [{ ingredient_id, name_id, action, text? }] }` (≤ 100) |

## Data

None in admin.

## Errors

| Case | HTTP |
|---|---|
| Invalid field value (kitchen's validation) | 400 with field name |
| Unknown ingredient | 404 |

## Acceptance criteria

- [ ] Changing only the density sends one `UpdateIngredient` with mask `density`; kitchen's `admin_changes` shows it.
- [ ] Changing an allergen flag writes old and new value to the audit log.
- [ ] Accepting 100 translations in one request works; 101 → 400.
