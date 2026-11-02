// Package events writes kitchen's events to the transactional outbox and defines the envelope they travel in
// (docs/EPICs/README.md "Events", docs/service-conventions.md "Events").
//
// A handler records an event in the same transaction as the change it describes, so the event is published if and
// only if the change commits:
//
//	tx, err := pool.Begin(ctx)
//	...
//	qtx := q.WithTx(tx)
//	... business writes with qtx ...
//	if err := events.Add(ctx, qtx, events.Event{Type: events.TypeIngredientMerged, Payload: events.IngredientMerged{...}}); err != nil {
//		return nil, connect.NewError(connect.CodeInternal, err) // the deferred Rollback drops the change too
//	}
//	return ..., tx.Commit(ctx)
//
// Calling Add with Queries that aren't a transaction's is a programming error; tests catch it with
// eventstest.RequireTx. The relay (KIT-01-03) publishes the rows in id order.
//
// The Consumer (KIT-01-04) applies other services' events: each registered Handler runs in one transaction with
// the record that makes a second delivery of the same event a no-op.
//
//	consumer.Handle("household.plan_changed", 1, func(ctx context.Context, q *store.Queries, e events.Envelope) error {
//		... decode e.Payload, write with q ...
//	})
//
// Planner and offers copy this package (PLN-01-01, OFF-01-01) and keep it identical, except the store import path,
// Producer and types.go, which lists each service's own event types.
package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"coolinary/kitchen/internal/kitchen/store"
)

// Producer is the envelope's producer field.
const Producer = "kitchen"

// MaxPayloadBytes is the largest payload JSON one event may carry (64 KB). Callers split larger fan-outs into
// several events.
const MaxPayloadBytes = 64 << 10

var (
	// ErrInvalidType: the type doesn't match ^[a-z]+(\.[a-z_]+)+$.
	ErrInvalidType = errors.New("invalid event type")
	// ErrUnknownType: the type isn't in Types.
	ErrUnknownType = errors.New("unknown event type")
	// ErrWrongPayload: the payload isn't the payload struct Types names for the type.
	ErrWrongPayload = errors.New("payload is not the event type's struct")
	// ErrInvalidHousehold: the household id is neither empty nor a UUID.
	ErrInvalidHousehold = errors.New("event household_id must be a UUID or empty")
	// ErrPayloadTooLarge: the payload JSON is over MaxPayloadBytes.
	ErrPayloadTooLarge = errors.New("event payload over 64 KB")
)

// Event is one change for other services to learn about.
type Event struct {
	Type Type
	// HouseholdID is the household the event is about, or "" for catalogue-wide events.
	HouseholdID string
	// Payload is the type's payload struct (see Types), by value. Ids, codes, dates and counts only.
	Payload any
}

// Add writes e as one outbox row through q, which must be a transaction's Queries (q.WithTx(tx)). It generates the
// envelope id (UUIDv7) and occurred_at, and fails before any insert when the type is malformed or unregistered, the
// payload isn't the type's struct or its JSON is over 64 KB, or the household id isn't a UUID. The caller turns an
// error into INTERNAL and rolls the transaction back.
func Add(ctx context.Context, q *store.Queries, e Event) error {
	row, err := newRow(e, time.Now())
	if err != nil {
		return err
	}
	if err := q.InsertOutboxEvent(ctx, row); err != nil {
		return fmt.Errorf("insert %s event: %w", e.Type, err)
	}
	return nil
}

var typePattern = regexp.MustCompile(`^[a-z]+(\.[a-z_]+)+$`)

// newRow checks e and builds its outbox row. Errors name the type, never the event's data.
func newRow(e Event, now time.Time) (store.InsertOutboxEventParams, error) {
	if !typePattern.MatchString(string(e.Type)) {
		return store.InsertOutboxEventParams{}, fmt.Errorf("%w %q: want %s", ErrInvalidType, e.Type, typePattern)
	}
	def, ok := Types[e.Type]
	if !ok {
		return store.InsertOutboxEventParams{}, fmt.Errorf("%w %q: add it to events.Types", ErrUnknownType, e.Type)
	}
	if reflect.TypeOf(e.Payload) != def.payload {
		return store.InsertOutboxEventParams{}, fmt.Errorf("%w: %s takes %s, got %T", ErrWrongPayload, e.Type, def.payload, e.Payload)
	}
	household, ok := parseHousehold(e.HouseholdID)
	if !ok {
		return store.InsertOutboxEventParams{}, fmt.Errorf("%w (%s)", ErrInvalidHousehold, e.Type)
	}
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return store.InsertOutboxEventParams{}, fmt.Errorf("encode %s payload: %w", e.Type, err)
	}
	if len(payload) > MaxPayloadBytes {
		return store.InsertOutboxEventParams{}, fmt.Errorf("%w: %s payload is %d bytes, split it into several events", ErrPayloadTooLarge, e.Type, len(payload))
	}
	at := now.UTC().Truncate(time.Millisecond)
	return store.InsertOutboxEventParams{
		ID:          pgtype.UUID{Bytes: ids.next(at), Valid: true},
		Type:        string(e.Type),
		Version:     def.Version,
		HouseholdID: household,
		Payload:     payload,
		OccurredAt:  pgtype.Timestamptz{Time: at, Valid: true},
	}, nil
}

// parseHousehold reads a household id in the 8-4-4-4-12 form; "" is NULL, a catalogue-wide event.
func parseHousehold(s string) (pgtype.UUID, bool) {
	if s == "" {
		return pgtype.UUID{}, true
	}
	if len(s) != 36 || s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return pgtype.UUID{}, false
	}
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		return pgtype.UUID{}, false
	}
	return u, true
}
