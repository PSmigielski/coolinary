package events

import "encoding/json"

// Envelope is the JSON every event travels in: the stream entry's envelope field (Q015). Field names are snake_case.
type Envelope struct {
	ID          string          `json:"id"`           // UUIDv7, the consumers' idempotency key
	Type        string          `json:"type"`         // e.g. catalogue.ingredient_merged
	Version     int32           `json:"version"`      // starts at 1, raised only for breaking payload changes
	OccurredAt  string          `json:"occurred_at"`  // RFC 3339, UTC, milliseconds
	Producer    string          `json:"producer"`     // the service that published it
	HouseholdID string          `json:"household_id"` // "" for catalogue-wide events
	Payload     json.RawMessage `json:"payload"`
}
