package events_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"regexp"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/db"
	"coolinary/kitchen/internal/platform/events"
	"coolinary/kitchen/internal/platform/events/eventstest"
	"coolinary/kitchen/migrations"
)

var uuidV7 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// openTestDB returns a pool on the migrated test database, or skips when TEST_DATABASE_URL isn't set
// (make test sets it while make dev is up).
func openTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("needs TEST_DATABASE_URL (make dev, then make test s=kitchen)")
	}
	pool, err := db.Open(t.Context(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := db.Migrate(t.Context(), pool, migrations.FS, slog.New(slog.DiscardHandler)); err != nil {
		t.Fatal(err)
	}
	return pool
}

func randomUUID(t *testing.T) string {
	t.Helper()
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatal(err)
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return pgtype.UUID{Bytes: b, Valid: true}.String()
}

// unknownEAN returns a random 13-digit code in the GS1 example range that no test or migration stores.
func unknownEAN(t *testing.T) string {
	t.Helper()
	n, err := rand.Int(rand.Reader, big.NewInt(1e10))
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("952%010d", n.Int64())
}

var errHandlerFailed = errors.New("handler failed after writing")

// writeProductAndEvent stands in for an RPC handler: in one transaction it writes a business row (a product) and an
// event, then commits, or rolls back when fail is set or events.Add fails.
func writeProductAndEvent(ctx context.Context, pool *pgxpool.Pool, q *store.Queries, ean string, e events.Event, fail bool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "INSERT INTO products (ean, name) VALUES ($1, 'Test product')", ean); err != nil {
		return err
	}
	if err := events.Add(ctx, q.WithTx(tx), e); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	if fail {
		return errHandlerFailed
	}
	return tx.Commit(ctx)
}

// testStore returns Queries on the test database that fail the test when events.Add is called on them outside a
// transaction, and removes the products and outbox rows the test names when it ends.
func testStore(t *testing.T, pool *pgxpool.Pool, eans, fromIngredientIDs []string) *store.Queries {
	t.Helper()
	t.Cleanup(func() {
		ctx := context.Background()
		if _, err := pool.Exec(ctx, "DELETE FROM products WHERE ean = ANY($1)", eans); err != nil {
			t.Error(err)
		}
		if _, err := pool.Exec(ctx, "DELETE FROM outbox WHERE payload->>'from_ingredient_id' = ANY($1)", fromIngredientIDs); err != nil {
			t.Error(err)
		}
	})
	return store.New(eventstest.RequireTx(t, pool))
}

// outboxRows returns the stored catalogue.ingredient_merged rows with this from_ingredient_id, in id order.
func outboxRows(t *testing.T, pool *pgxpool.Pool, fromIngredientID string) []store.Outbox {
	t.Helper()
	rows, err := pool.Query(t.Context(), `
		SELECT id, type, version, household_id, payload, occurred_at, published_at, attempts, last_error
		FROM outbox WHERE payload->>'from_ingredient_id' = $1 ORDER BY id`, fromIngredientID)
	if err != nil {
		t.Fatal(err)
	}
	out, err := pgx.CollectRows(rows, pgx.RowToStructByPos[store.Outbox])
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func productExists(t *testing.T, pool *pgxpool.Pool, ean string) bool {
	t.Helper()
	var n int
	if err := pool.QueryRow(t.Context(), "SELECT count(*) FROM products WHERE ean = $1", ean).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n == 1
}

func payloadOf(t *testing.T, row store.Outbox) events.IngredientMerged {
	t.Helper()
	var p events.IngredientMerged
	if err := json.Unmarshal(row.Payload, &p); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestAdd_RolledBackChangeLeavesNeitherRow(t *testing.T) {
	pool := openTestDB(t)
	ean, from := unknownEAN(t), randomUUID(t)
	q := testStore(t, pool, []string{ean}, []string{from})

	err := writeProductAndEvent(t.Context(), pool, q, ean, merged(from, randomUUID(t)), true)

	if !errors.Is(err, errHandlerFailed) {
		t.Fatalf("handler error = %v, want %v", err, errHandlerFailed)
	}
	if productExists(t, pool, ean) {
		t.Error("the product of the rolled-back transaction exists")
	}
	if rows := outboxRows(t, pool, from); len(rows) != 0 {
		t.Errorf("the rolled-back transaction left %d outbox rows, want none", len(rows))
	}
}

func TestAdd_CommittedChangeKeepsBothRows(t *testing.T) {
	pool := openTestDB(t)
	ean, from, into := unknownEAN(t), randomUUID(t), randomUUID(t)
	q := testStore(t, pool, []string{ean}, []string{from})
	before := time.Now()

	err := writeProductAndEvent(t.Context(), pool, q, ean, merged(from, into), false)

	after := time.Now()
	if err != nil {
		t.Fatalf("handler error = %v, want none", err)
	}
	if !productExists(t, pool, ean) {
		t.Error("the product of the committed transaction is missing")
	}
	rows := outboxRows(t, pool, from)
	if len(rows) != 1 {
		t.Fatalf("the committed transaction left %d outbox rows, want 1", len(rows))
	}
	row := rows[0]
	if !uuidV7.MatchString(row.ID.String()) {
		t.Errorf("id = %s, want a UUIDv7", row.ID)
	}
	if row.Type != "catalogue.ingredient_merged" || row.Version != 1 || row.HouseholdID.Valid {
		t.Errorf("type, version, household_id = %q, %d, %v; want catalogue.ingredient_merged, 1, NULL",
			row.Type, row.Version, row.HouseholdID)
	}
	if got, want := payloadOf(t, row), (events.IngredientMerged{FromIngredientID: from, IntoIngredientID: into}); got != want {
		t.Errorf("payload = %+v, want %+v", got, want)
	}
	if row.PublishedAt.Valid || row.Attempts != 0 || row.LastError.Valid {
		t.Errorf("published_at, attempts, last_error = %v, %d, %v; want NULL, 0, NULL", row.PublishedAt, row.Attempts, row.LastError)
	}
	occurred := row.OccurredAt.Time
	if !occurred.Equal(occurred.Truncate(time.Millisecond)) {
		t.Errorf("occurred_at = %s, want whole milliseconds", occurred.Format(time.RFC3339Nano))
	}
	if occurred.Before(before.Truncate(time.Millisecond)) || occurred.After(after) {
		t.Errorf("occurred_at = %s, want between %s and %s", occurred, before, after)
	}
}

func TestAdd_FailingEventRollsBackTheChange(t *testing.T) {
	pool := openTestDB(t)
	ean, from := unknownEAN(t), randomUUID(t)
	q := testStore(t, pool, []string{ean}, []string{from})
	e := merged(from, randomUUID(t))
	e.Type = "catalogue.ingredient_renamed"

	err := writeProductAndEvent(t.Context(), pool, q, ean, e, false)

	if connect.CodeOf(err) != connect.CodeInternal || !errors.Is(err, events.ErrUnknownType) {
		t.Fatalf("handler error = %v, want INTERNAL wrapping %v", err, events.ErrUnknownType)
	}
	if productExists(t, pool, ean) {
		t.Error("the product of the failed transaction exists")
	}
	if rows := outboxRows(t, pool, from); len(rows) != 0 {
		t.Errorf("the failed transaction left %d outbox rows, want none", len(rows))
	}
}

func TestAdd_HouseholdEventStoresItsHousehold(t *testing.T) {
	pool := openTestDB(t)
	ean, from, household := unknownEAN(t), randomUUID(t), randomUUID(t)
	q := testStore(t, pool, []string{ean}, []string{from})
	e := merged(from, randomUUID(t))
	e.HouseholdID = household

	if err := writeProductAndEvent(t.Context(), pool, q, ean, e, false); err != nil {
		t.Fatalf("handler error = %v, want none", err)
	}

	rows := outboxRows(t, pool, from)
	if len(rows) != 1 || rows[0].HouseholdID.String() != household {
		t.Errorf("stored rows = %+v, want one with household_id %s", rows, household)
	}
}

func TestAdd_EventsOfOneTransactionSortInCallOrder(t *testing.T) {
	pool := openTestDB(t)
	from := randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	var want []string
	tx, err := pool.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()
	for range 50 {
		into := randomUUID(t)
		want = append(want, into)
		if err := events.Add(t.Context(), q.WithTx(tx), merged(from, into)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tx.Commit(t.Context()); err != nil {
		t.Fatal(err)
	}

	rows := outboxRows(t, pool, from)

	if len(rows) != len(want) {
		t.Fatalf("stored %d rows, want %d", len(rows), len(want))
	}
	for i, row := range rows {
		if got := payloadOf(t, row).IntoIngredientID; got != want[i] {
			t.Fatalf("row %d in id order is event %s, want %s (call order)", i, got, want[i])
		}
	}
}
