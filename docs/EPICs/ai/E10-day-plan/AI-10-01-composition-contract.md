# AI-10-01 Composition contract and shared dish types

| | |
|---|---|
| Service | ai |
| Epic | [AI-10 Day plan composition and cuisine choices](README.md) |
| Release | R1 |
| Size | M |
| Status | todo |
| Depends on | AI-02-01 |
| Unblocks | AI-10-02, AI-09-03, AI-09-05, AI-11, AI-12, AI-14, AI-15, AI-17; PLN-08 |

## Goal

One set of messages for everything that composes or rewrites dishes, built so that the AI can only refer to
ingredients it was given and never needs, or produces, calorie numbers.

## Business rules

1. AI never sees user free text in composition except a wish (AI-11-03); every ingredient is referred to by the
   catalogue id planner sent, with the catalogue name in the request language as its label.
2. AI never receives health data: no age, height, weight or goal. R2 family requests carry only per-eater meal
   targets in kcal (AI-12-01), as features.md §2 says ("AI receives only each person's daily target, never their
   weight").
3. Amounts are `coolinary.common.v1.Quantity` in `GRAM`, `MILLILITRE` or `PIECE` only; `display_hint` (≤ 20 chars,
   e.g. "2 tbsp") is display text and never used for maths. A `stock_uses` amount is always in the unit of its stock
   line, so totals can be checked without unit conversion.
4. Tools, cuisines, diets, allergens and effort tags use the shared enums in `coolinary/common/v1/food.proto`
   (README decision 16: `KitchenTool` KIT-10-03, `Cuisine` KIT-10-02, `Diet` KIT-10-04, `Allergen` KIT-02-01,
   `EffortTag` PLN-02-01); ai defines no copies. In prompts and output schemas ai writes them as lowercase codes
   (`KITCHEN_TOOL_AIR_FRYER` → `air_fryer`) and maps the answer back; an unknown code fails validation (AI-10-03).
5. Input size limits per request: stock ≤ 150 lines; `ingredients` facts ≤ 400; library candidates ≤ 30;
   `buyable_ids` ≤ 80; meals 1–3; portions 1–24.
6. Every composition response is a list of `ComposedMeal`; a meal holds either a library dish id or a `NewDish`.
7. Decided (Q103): stock past its date still counts as at home. Planner sends such lines with `past_date = true` (and
   `days_to_expiry = 0`); the AI may use them, and every `stock_uses` entry on such a line comes back with
   `check_first = true`. ai sets this flag itself after validation (AI-10-03), whatever the model wrote, so the app
   can show "check it first" wherever the item is used.

## Contract

`contracts/coolinary/ai/v1/compose.proto`:

```proto
enum MealSlot { MEAL_SLOT_UNSPECIFIED = 0; MEAL_SLOT_BREAKFAST = 1; MEAL_SLOT_DINNER = 2; MEAL_SLOT_SUPPER = 3; }
// Diet, Allergen, Cuisine, KitchenTool, EffortTag: imported from coolinary/common/v1/food.proto (decision 16).
enum IngredientRole { INGREDIENT_ROLE_UNSPECIFIED = 0; INGREDIENT_ROLE_MAIN = 1; INGREDIENT_ROLE_BASE = 2;
  INGREDIENT_ROLE_TOPPING = 3; INGREDIENT_ROLE_GARNISH = 4; }

// Facts about a catalogue ingredient (from kitchen KIT-02); the only way AI learns what an id is.
message IngredientFact {
  string ingredient_id = 1;
  string name = 2;                 // catalogue name in the request language
  repeated coolinary.common.v1.Allergen allergens = 3;
  bool vegan = 4;
  bool vegetarian = 5;
  bool gluten_free = 6;
  bool lactose_free = 7;
  string category = 8;
}

message StockLine {
  string ingredient_id = 1;
  coolinary.common.v1.Quantity available = 2;
  int32 days_to_expiry = 3;        // 0 = today (also for past-date lines); -1 = unknown
  bool opened = 4;
  bool past_date = 5;              // past its date; usable, shown as "check it first" (Q103)
}

message KitchenRules {
  repeated coolinary.common.v1.KitchenTool tools = 1;     // tools the kitchen has (KIT-10)
  repeated coolinary.common.v1.Diet diets = 2;            // hard
  repeated coolinary.common.v1.Allergen excluded_allergens = 3;  // hard
  repeated string disliked_ingredient_ids = 4;  // hard (features.md §1 "Diet and allergens")
}

message MealRequest {
  MealSlot meal = 1;
  repeated coolinary.common.v1.Cuisine cuisines = 2;  // empty = any
  int32 max_total_minutes = 3;
  bool main_hot_meal = 4;
  bool prefer_no_cooking = 5;
}

message CandidateIngredient {
  string ingredient_id = 1;
  coolinary.common.v1.Quantity amount = 2;   // for the request's portions
  bool optional = 3;
}

message LibraryCandidate {
  string dish_id = 1;
  string name = 2;
  repeated MealSlot meals = 3;
  coolinary.common.v1.Cuisine cuisine = 4;
  repeated coolinary.common.v1.KitchenTool tools = 5;
  int32 total_minutes = 6;
  repeated CandidateIngredient ingredients = 7;
  repeated coolinary.common.v1.EffortTag effort_tags = 8;
  double share_on_hand = 9;        // planner's ranking (PLN-06)
  int32 kcal_per_portion = 10;     // set only in family requests (AI-12); never echoed in text
}

message DishIngredient {
  string ingredient_id = 1;
  coolinary.common.v1.Quantity amount = 2;
  string display_hint = 3;
  bool optional = 4;
  IngredientRole role = 5;
  string note = 6;                 // ≤ 60 chars, e.g. "finely chopped"
}

message Step {
  string text = 1;                 // ≤ 400 chars
  int32 timer_seconds = 2;         // 0 = no timer
}

message NewDish {
  string name = 1;                 // ≤ 80 chars
  coolinary.common.v1.Cuisine cuisine = 2;
  repeated MealSlot meals = 3;
  int32 servings = 4;              // = request portions
  int32 total_minutes = 5;
  int32 active_minutes = 6;
  repeated coolinary.common.v1.KitchenTool tools = 7;
  repeated DishIngredient ingredients = 8;
  repeated Step steps = 9;         // empty only for alternatives (AI-10-04)
  repeated coolinary.common.v1.EffortTag effort_tags = 10;
  int32 keeps_days = 11;           // 0-5, used by week plans (AI-15)
  bool freezes = 12;
}

message StockUse {
  string ingredient_id = 1;
  coolinary.common.v1.Quantity amount = 2;   // unit of the stock line
  bool check_first = 3;                      // set by ai for past-date stock lines (Q103)
}

message BuyItem {
  string ingredient_id = 1;        // from buyable_ids
  coolinary.common.v1.Quantity amount = 2;
}

enum MealReason { MEAL_REASON_UNSPECIFIED = 0; MEAL_REASON_USES_EXPIRING = 1; MEAL_REASON_LIBRARY_DISH = 2;
  MEAL_REASON_CUISINE_MATCH = 3; MEAL_REASON_FAST = 4; MEAL_REASON_NO_COOKING = 5; MEAL_REASON_USES_OFFER = 6; }

message ComposedMeal {
  MealSlot meal = 1;
  oneof dish {
    string library_dish_id = 2;
    NewDish new_dish = 3;
  }
  repeated StockUse stock_uses = 4;
  repeated BuyItem to_buy = 5;
  repeated MealReason reasons = 6;
  CuisineFit cuisine_fit = 7;               // AI-10-04
  repeated MealAlternative alternatives = 8; // AI-10-04
}
```

Shared request part of every composition RPC:

```proto
message CompositionContext {
  string household_id = 1;
  string language = 2;
  string date = 3;                          // YYYY-MM-DD, for season only
  int32 portions = 4;                       // 1-24
  KitchenRules rules = 5;
  repeated StockLine stock = 6;             // ≤ 150
  repeated string staple_ids = 7;           // always at home, never used up
  repeated IngredientFact ingredients = 8;  // facts for every id in stock, staples, candidates, buyable
  repeated LibraryCandidate candidates = 9; // ≤ 30, ranked by planner
  repeated string buyable_ids = 10;         // ≤ 80 (cuisine base pantry and common items)
  int32 max_items_to_buy = 11;              // 0-3 for the day
}
```

`CuisineFit` and `MealAlternative` are defined in AI-10-04; `ComposeDayPlanRequest/Response` in AI-10-02.

## Data

No schema change.

## Errors

| Case | gRPC status |
|---|---|
| Request breaks the limits above | `INVALID_ARGUMENT`, `ai_bad_request` |

## Acceptance criteria

- [ ] `make proto` generates Go code for planner and Python code for ai without lint errors.
- [ ] The pydantic output schema for a composed meal mirrors these messages (test compares field names).
