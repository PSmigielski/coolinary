package events

import (
	"encoding/json"

	"coolinary/kitchen/internal/kitchen/store"
)

// Envelope is the JSON every event travels in: the stream entry's envelope field (Q015). Field names are snake_case.
type Envelope struct {
	ID          string          `json:"id"`           // UUIDv7, the consumers' idempotency key
	Type        string          `json:"type"`         // e.g. catalogue.ingredient_merged
	Version     int32           `json:"version"`      // starts at 1, raised only for breaking payload changes
	OccurredAt  string          `json:"occurred_at"`  // RFC 3339, UTC, milliseconds
	Producer    string          `json:"producer"`     // kitchen
	HouseholdID string          `json:"household_id"` // "" for catalogue-wide events
	Payload     json.RawMessage `json:"payload"`
}

// occurredAtLayout is RFC 3339 with milliseconds; on a UTC time it ends in "Z".
const occurredAtLayout = "2006-01-02T15:04:05.000Z07:00"

// EnvelopeOf builds the envelope of an outbox row; the relay publishes its JSON.
func EnvelopeOf(row store.Outbox) Envelope {
	return Envelope{
		ID:          row.ID.String(),
		Type:        row.Type,
		Version:     row.Version,
		OccurredAt:  row.OccurredAt.Time.UTC().Format(occurredAtLayout),
		Producer:    Producer,
		HouseholdID: row.HouseholdID.String(), // "" when NULL
		Payload:     row.Payload,
	}
}
