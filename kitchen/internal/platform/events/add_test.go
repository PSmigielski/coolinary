package events_test

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/events"
)

const (
	ingredientA = "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a01"
	ingredientB = "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a02"
	householdA  = "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b"
)

// countingDB is a store.DBTX that counts the statements it gets and runs none.
type countingDB struct{ execs int }

func (d *countingDB) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	d.execs++
	return pgconn.CommandTag{}, nil
}

func (d *countingDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errors.New("countingDB runs no queries")
}

func (d *countingDB) QueryRow(context.Context, string, ...any) pgx.Row { return nil }

func merged(from, into string) events.Event {
	return events.Event{
		Type:    events.TypeIngredientMerged,
		Payload: events.IngredientMerged{FromIngredientID: from, IntoIngredientID: into},
	}
}

// mergedOfSize returns a catalogue.ingredient_merged event whose payload JSON is exactly n bytes.
func mergedOfSize(t *testing.T, n int) events.Event {
	t.Helper()
	empty, err := json.Marshal(events.IngredientMerged{})
	if err != nil {
		t.Fatal(err)
	}
	p := events.IngredientMerged{FromIngredientID: strings.Repeat("0", n-len(empty))}
	if b, err := json.Marshal(p); err != nil || len(b) != n {
		t.Fatalf("payload is %d bytes (%v), want %d", len(b), err, n)
	}
	return events.Event{Type: events.TypeIngredientMerged, Payload: p}
}

func TestAdd_RejectsBeforeInsert(t *testing.T) {
	withType := func(typ events.Type) events.Event {
		e := merged(ingredientA, ingredientB)
		e.Type = typ
		return e
	}
	withPayload := func(p any) events.Event {
		e := merged(ingredientA, ingredientB)
		e.Payload = p
		return e
	}
	withHousehold := func(h string) events.Event {
		e := merged(ingredientA, ingredientB)
		e.HouseholdID = h
		return e
	}
	tests := []struct {
		name  string
		event events.Event
		want  error
	}{
		{name: "empty type", event: withType(""), want: events.ErrInvalidType},
		{name: "type without a dot", event: withType("catalogue"), want: events.ErrInvalidType},
		{name: "type ending in a dot", event: withType("catalogue."), want: events.ErrInvalidType},
		{name: "type with a capital letter", event: withType("Catalogue.ingredient_merged"), want: events.ErrInvalidType},
		{name: "type with a hyphen", event: withType("catalogue.ingredient-merged"), want: events.ErrInvalidType},
		{name: "type with a digit", event: withType("catalogue.ingredient_merged2"), want: events.ErrInvalidType},
		{name: "well-formed type that is not registered", event: withType("catalogue.ingredient_renamed"), want: events.ErrUnknownType},
		{name: "payload of another struct", event: withPayload(struct{ Name string }{"someone@example.test"}), want: events.ErrWrongPayload},
		{name: "payload as a map", event: withPayload(map[string]any{"from_ingredient_id": "someone@example.test"}), want: events.ErrWrongPayload},
		{name: "payload as a pointer to the struct", event: withPayload(&events.IngredientMerged{}), want: events.ErrWrongPayload},
		{name: "no payload", event: withPayload(nil), want: events.ErrWrongPayload},
		{name: "household id that is not a uuid", event: withHousehold("someone@example.test"), want: events.ErrInvalidHousehold},
		{name: "household id with x for its dashes", event: withHousehold("0190a1b2xc3d4x7e5fx8a9bx0c1d2e3f4a5b"), want: events.ErrInvalidHousehold},
		{name: "household id one character short", event: withHousehold(householdA[:35]), want: events.ErrInvalidHousehold},
		{name: "payload one byte over 64 KB", event: mergedOfSize(t, 64*1024+1), want: events.ErrPayloadTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &countingDB{}

			err := events.Add(t.Context(), store.New(db), tt.event)

			if !errors.Is(err, tt.want) {
				t.Errorf("Add() error = %v, want %v", err, tt.want)
			}
			if err != nil && strings.Contains(err.Error(), "example.test") {
				t.Errorf("Add() error %q repeats the event's data", err)
			}
			if db.execs != 0 {
				t.Errorf("Add() ran %d statements, want none", db.execs)
			}
		})
	}
}

func TestAdd_InsertsOneRow(t *testing.T) {
	withHousehold := merged(ingredientA, ingredientB)
	withHousehold.HouseholdID = householdA
	upperHousehold := merged(ingredientA, ingredientB)
	upperHousehold.HouseholdID = strings.ToUpper(householdA)
	tests := []struct {
		name  string
		event events.Event
	}{
		{name: "catalogue-wide event", event: merged(ingredientA, ingredientB)},
		{name: "household event", event: withHousehold},
		{name: "household id in upper case", event: upperHousehold},
		{name: "payload of exactly 64 KB", event: mergedOfSize(t, 64*1024)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &countingDB{}

			err := events.Add(t.Context(), store.New(db), tt.event)

			if err != nil {
				t.Fatalf("Add() error = %v, want none", err)
			}
			if db.execs != 1 {
				t.Errorf("Add() ran %d statements, want 1", db.execs)
			}
		})
	}
}

func TestTypes_AreWellFormedAndVersioned(t *testing.T) {
	pattern := regexp.MustCompile(`^[a-z]+(\.[a-z_]+)+$`)
	for typ, def := range events.Types {
		if !pattern.MatchString(string(typ)) {
			t.Errorf("registered type %q does not match %s", typ, pattern)
		}
		if def.Version < 1 {
			t.Errorf("registered type %q has version %d, want 1 or more", typ, def.Version)
		}
	}
}
