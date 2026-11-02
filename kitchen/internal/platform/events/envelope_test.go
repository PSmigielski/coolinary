package events_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/events"
)

func mustUUID(t *testing.T, s string) pgtype.UUID {
	t.Helper()
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		t.Fatal(err)
	}
	return u
}

// sampleRow is the outbox row of a catalogue.ingredient_merged event, as the relay reads it.
func sampleRow(t *testing.T) store.Outbox {
	t.Helper()
	payload, err := json.Marshal(events.IngredientMerged{FromIngredientID: ingredientA, IntoIngredientID: ingredientB})
	if err != nil {
		t.Fatal(err)
	}
	return store.Outbox{
		ID:      mustUUID(t, "019a4c3e-5b2a-7c1d-8e9f-0a1b2c3d4e5f"),
		Type:    string(events.TypeIngredientMerged),
		Version: 1,
		Payload: payload,
		// Sub-millisecond digits and a non-UTC zone, both normalised away in the envelope.
		OccurredAt: pgtype.Timestamptz{Time: time.Date(2026, 11, 2, 11, 30, 15, 123456789, time.FixedZone("CET", 3600)), Valid: true},
		Attempts:   2,
		LastError:  pgtype.Text{String: "redis unavailable", Valid: true},
	}
}

func TestEnvelope_MatchesGoldenFile(t *testing.T) {
	row := sampleRow(t)

	got, err := json.Marshal(events.EnvelopeOf(row))

	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "envelope.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got)+"\n" != string(want) {
		t.Errorf("envelope JSON\n got %s\nwant %s", got, want)
	}
}

func TestEnvelope_HouseholdEventCarriesItsHousehold(t *testing.T) {
	row := sampleRow(t)
	row.HouseholdID = mustUUID(t, householdA)

	got := events.EnvelopeOf(row)

	if got.HouseholdID != householdA {
		t.Errorf("envelope household_id = %q, want %q", got.HouseholdID, householdA)
	}
}
