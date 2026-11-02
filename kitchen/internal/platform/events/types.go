package events

import "reflect"

// This file is kitchen's own (the rest of the package is shared with planner and offers): every event type kitchen
// writes, its version and its payload struct. Payloads hold ids, codes, dates and counts only, never names,
// e-mail addresses, free text or health data; field names are snake_case and stable. Adding an optional field keeps
// the version; removing, renaming or changing the meaning of one raises it (docs/service-conventions.md "Events").

// Type is an event type, "<area>.<what_happened>" (^[a-z]+(\.[a-z_]+)+$).
type Type string

const (
	// TypeIngredientMerged: an admin merged one ingredient into another (KIT-02-09). Catalogue-wide.
	TypeIngredientMerged Type = "catalogue.ingredient_merged"
)

// IngredientMerged is the payload of catalogue.ingredient_merged.
type IngredientMerged struct {
	FromIngredientID string `json:"from_ingredient_id"`
	IntoIngredientID string `json:"into_ingredient_id"`
}

// Definition is a registered event type: its envelope version and its payload struct.
type Definition struct {
	Version int32
	payload reflect.Type
}

func define[P any](version int32) Definition {
	return Definition{Version: version, payload: reflect.TypeFor[P]()}
}

// Types lists every event type kitchen writes. Add refuses any other.
var Types = map[Type]Definition{
	TypeIngredientMerged: define[IngredientMerged](1),
}
