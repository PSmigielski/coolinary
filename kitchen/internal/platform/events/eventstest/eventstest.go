// Package eventstest catches events.Add called outside a transaction, in tests.
package eventstest

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"coolinary/kitchen/internal/kitchen/store"
)

// RequireTx wraps db, the connection a test builds its *store.Queries on (store.New(eventstest.RequireTx(t, pool))),
// so that events.Add through those Queries fails the test: outbox rows must be written with q.WithTx(tx), in the
// transaction of the change they describe. Every other statement passes through to db, and calls made through
// q.WithTx(tx) never reach the wrapper.
func RequireTx(t testing.TB, db store.DBTX) store.DBTX {
	return &guard{t: t, db: db}
}

var errOutsideTx = errors.New("events.Add outside a transaction: use q.WithTx(tx)")

// outboxInsert starts the SQL that sqlc generates for the query events.Add runs.
const outboxInsert = "-- name: InsertOutboxEvent "

type guard struct {
	t  testing.TB
	db store.DBTX
}

func (g *guard) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if strings.HasPrefix(sql, outboxInsert) {
		g.t.Helper()
		g.t.Error(errOutsideTx.Error())
		return pgconn.CommandTag{}, errOutsideTx
	}
	return g.db.Exec(ctx, sql, args...)
}

func (g *guard) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return g.db.Query(ctx, sql, args...)
}

func (g *guard) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return g.db.QueryRow(ctx, sql, args...)
}
