# KIT-10-03 Kitchen tools and presets

| | |
|---|---|
| Service | kitchen |
| Epic | [KIT-10 Kitchen profile](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | KIT-10-01 |
| Unblocks | KIT-10-07, KIT-10-08, PLN-07 (tool check), PLN-05 (tool filter) |

## Goal

The household ticks the tools it owns, or picks a preset in one tap, and no dish ever needs an unticked tool
(features.md §1 "Kitchen tools": hard limit).

## Business rules

1. Tool codes (shared enum `coolinary.common.v1.KitchenTool`, also used for dish tags by PLN-02-01): hob, oven,
   microwave, air fryer, blender, hand blender, food processor, multicooker / pressure cooker, cooking robot,
   toaster / sandwich maker, grill, steamer, rice cooker, mixer, wok, baking tins, kettle.
2. "Kettle" is not in the features.md checklist but the "student room" preset names it, so it is a tool too.
   Decided (Q061): the kettle is a tool that dishes can require.
3. Presets (returned by KIT-10-07, applied by the app as a full tool list):
   1. `student_room`: microwave, kettle;
   2. `basic_kitchen`: hob, oven, kettle, toaster, hand blender, baking tins (Decided, Q060; kitchen owns this
      list and planner's coverage targets use it, PLN-03-09);
   3. `fully_equipped`: all 17 tools.
4. The stored value is always the full list; a preset is not remembered (changing one tool after a preset is
   normal). 0–17 tools, no duplicates.
5. Default before the first save: `basic_kitchen`.
6. An empty list is allowed: planner then only offers "no cooking" dishes.
7. Tools are a hard limit everywhere (planner); the one-off "air fryer only tonight" filter is planner's request
   parameter and never changes the profile.

## Contract

```proto
// coolinary/common/v1/food.proto
enum KitchenTool {
  KITCHEN_TOOL_UNSPECIFIED = 0;
  KITCHEN_TOOL_HOB = 1; KITCHEN_TOOL_OVEN = 2; KITCHEN_TOOL_MICROWAVE = 3; KITCHEN_TOOL_AIR_FRYER = 4;
  KITCHEN_TOOL_BLENDER = 5; KITCHEN_TOOL_HAND_BLENDER = 6; KITCHEN_TOOL_FOOD_PROCESSOR = 7;
  KITCHEN_TOOL_MULTICOOKER = 8; KITCHEN_TOOL_COOKING_ROBOT = 9; KITCHEN_TOOL_TOASTER = 10; KITCHEN_TOOL_GRILL = 11;
  KITCHEN_TOOL_STEAMER = 12; KITCHEN_TOOL_RICE_COOKER = 13; KITCHEN_TOOL_MIXER = 14; KITCHEN_TOOL_WOK = 15;
  KITCHEN_TOOL_BAKING_TINS = 16; KITCHEN_TOOL_KETTLE = 17;
}
```

Section `tools` of `UpdateProfile` (KIT-10-01).

## Data

`kitchen/migrations/000NN_profile_tools.sql`

```sql
-- +goose Up
CREATE TABLE profile_tools (
    household_id uuid NOT NULL REFERENCES kitchen_profiles (household_id) ON DELETE CASCADE,
    tool         smallint NOT NULL CHECK (tool BETWEEN 1 AND 17),
    PRIMARY KEY (household_id, tool)
);

-- +goose Down
DROP TABLE profile_tools;
```

## Errors

| Case | gRPC status | Reason |
|---|---|---|
| Unknown code or duplicate | `INVALID_ARGUMENT` | `invalid_value` |

## Acceptance criteria

- [ ] Saving the `student_room` list stores exactly microwave and kettle; features.md Verification 2 (50 composed
      dishes need no other tool) uses this profile.
- [ ] A new household reads the `basic_kitchen` tools with `is_default = true`.
- [ ] Saving an empty list is accepted.
