package eventstest_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/events"
	"coolinary/kitchen/internal/platform/events/eventstest"
)

// recordingT is a testing.TB that records the test failures reported to it.
type recordingT struct {
	testing.TB
	failures []string
}

func (r *recordingT) Helper() {}

func (r *recordingT) Error(args ...any) { r.failures = append(r.failures, fmt.Sprint(args...)) }

func (r *recordingT) Errorf(format string, args ...any) {
	r.failures = append(r.failures, fmt.Sprintf(format, args...))
}

// countingDB is a store.DBTX that counts the statements it gets; every query finds no rows.
type countingDB struct{ execs, queryRows int }

func (d *countingDB) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	d.execs++
	return pgconn.CommandTag{}, nil
}

func (d *countingDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errors.New("countingDB runs no queries")
}

func (d *countingDB) QueryRow(context.Context, string, ...any) pgx.Row {
	d.queryRows++
	return noRow{}
}

type noRow struct{}

func (noRow) Scan(...any) error { return pgx.ErrNoRows }

// countingTx is a transaction that counts its statements.
type countingTx struct {
	pgx.Tx
	execs int
}

func (tx *countingTx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	tx.execs++
	return pgconn.CommandTag{}, nil
}

var sample = events.Event{
	Type: events.TypeIngredientMerged,
	Payload: events.IngredientMerged{
		FromIngredientID: "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a01",
		IntoIngredientID: "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a02",
	},
}

func TestRequireTx_AddOutsideATransactionFailsTheTest(t *testing.T) {
	rt, db := &recordingT{}, &countingDB{}
	q := store.New(eventstest.RequireTx(rt, db))

	err := events.Add(t.Context(), q, sample)

	want := "events.Add outside a transaction: use q.WithTx(tx)"
	if len(rt.failures) != 1 || rt.failures[0] != want {
		t.Errorf("test failures = %q, want [%q]", rt.failures, want)
	}
	if err == nil {
		t.Error("Add() error = nil, want one")
	}
	if db.execs != 0 {
		t.Errorf("the outbox insert reached the database (%d statements), want none", db.execs)
	}
}

func TestRequireTx_AddInATransactionPasses(t *testing.T) {
	rt, db, tx := &recordingT{}, &countingDB{}, &countingTx{}
	q := store.New(eventstest.RequireTx(rt, db))

	err := events.Add(t.Context(), q.WithTx(tx), sample)

	if err != nil || len(rt.failures) != 0 {
		t.Errorf("Add() error = %v, test failures = %q; want neither", err, rt.failures)
	}
	if tx.execs != 1 || db.execs != 0 {
		t.Errorf("statements in the transaction, outside it = %d, %d; want 1, 0", tx.execs, db.execs)
	}
}

func TestRequireTx_OtherStatementsPassThrough(t *testing.T) {
	rt, db := &recordingT{}, &countingDB{}
	guarded := eventstest.RequireTx(rt, db)

	_, getErr := store.New(guarded).GetProductByEAN(t.Context(), "9520000000000")
	_, execErr := guarded.Exec(t.Context(), "UPDATE products SET confirmed = true WHERE ean = $1", "9520000000000")

	if !errors.Is(getErr, pgx.ErrNoRows) || execErr != nil || len(rt.failures) != 0 {
		t.Errorf("errors = %v, %v, test failures = %q; want no rows, none, none", getErr, execErr, rt.failures)
	}
	if db.queryRows != 1 || db.execs != 1 {
		t.Errorf("statements reaching the database = %d queries, %d execs; want 1, 1", db.queryRows, db.execs)
	}
}
