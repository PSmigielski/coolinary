package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"

	"coolinary/kitchen/internal/kitchen/store"
)

// Stream is the Redis stream (on redis-state) the relay publishes this service's events to.
const Stream = "events." + Producer

const (
	// batchSize is the most rows one pass publishes.
	batchSize = 100
	// interval is the longest wait between passes when no NOTIFY arrives.
	interval = time.Second
	// maxBackoff caps the wait after failed passes (1 s, 2 s, 4 s ...).
	maxBackoff = 30 * time.Second
	// retention is how long stream entries and published outbox rows are kept.
	retention = 7 * 24 * time.Hour
	// cleanupEvery is how often the loop deletes published rows older than retention.
	cleanupEvery = time.Hour
	// passTimeout bounds one pass, which finishes even when shutdown starts.
	passTimeout = 5 * time.Second
	// relistenAfter is the wait before the listener reconnects; passes go on every interval meanwhile.
	relistenAfter = 5 * time.Second
)

// Metric names are the same in every service that publishes events (Q041).
var (
	publishedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "events_published_total",
		Help: "Outbox events added to the stream and marked published.",
	}, []string{"stream", "type"})
	publishErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "events_publish_errors_total",
		Help: "Failed XADDs; each one stops its relay pass.",
	}, []string{"stream"})
	outboxPending = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "events_outbox_pending",
		Help: "Outbox rows not yet published.",
	})
	outboxOldest = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "events_outbox_oldest_seconds",
		Help: "Age of the oldest unpublished outbox row in seconds; 0 when none.",
	})
)

// RelayOptions change the relay's surroundings; the zero value is what the service runs with.
type RelayOptions struct {
	// Stream to publish to; Stream when empty.
	Stream string
	// Now reads the clock for the stream trim, the cleanup cutoff and schedule, and the backlog age; time.Now when nil.
	Now func() time.Time
	// After waits between passes like time.After (the default); tests pass a fake to step the loop.
	After func(time.Duration) <-chan time.Time
}

// Relay publishes committed outbox rows to the stream in id order, at least once (KIT-01-03). Run it as one
// g.Go(...) loop in main.go. A pass locks up to 100 unpublished rows, XADDs them in id order and marks them
// published, all in one transaction; a crash between XADD and commit publishes those rows again, which consumers
// drop by envelope id. Rows that commit while a pass runs go in a later pass, so the stream follows commit order:
// rows committed before a pass appear in id order, and the events one caller writes in sequence stay in order.
type Relay struct {
	pool   *pgxpool.Pool
	q      *store.Queries
	rdb    *redis.Client
	log    *slog.Logger
	stream string
	now    func() time.Time
	after  func(time.Duration) <-chan time.Time
}

// NewRelay builds the relay over the service's pool and its redis-state client.
func NewRelay(pool *pgxpool.Pool, rdb *redis.Client, log *slog.Logger, opts RelayOptions) *Relay {
	r := &Relay{pool: pool, q: store.New(pool), rdb: rdb, log: log, stream: opts.Stream, now: opts.Now, after: opts.After}
	if r.stream == "" {
		r.stream = Stream
	}
	if r.now == nil {
		r.now = time.Now
	}
	if r.after == nil {
		r.after = time.After
	}
	return r
}

// Run publishes until ctx is done, then returns nil after finishing the pass in progress. A pass starts on each
// NOTIFY from the outbox trigger and at the latest every second; right away again after a full batch; and after a
// backoff of 1 s, 2 s, 4 s ... (at most 30 s, notifications ignored) when it failed. Once an hour it also deletes
// published rows older than 7 days. Failures are logged and counted, never returned.
func (r *Relay) Run(ctx context.Context) error {
	ctx, stop := context.WithCancel(ctx)
	var wg sync.WaitGroup
	defer wg.Wait()
	defer stop()
	wake := make(chan struct{}, 1)
	conn, err := r.listen(ctx)
	if err != nil {
		r.log.Warn("outbox listener unavailable, passing every second", "error", err)
	}
	wg.Go(func() { r.awaitNotifications(ctx, conn, wake) })

	var lastCleanup time.Time
	failures := 0
	for ctx.Err() == nil {
		n, err := r.pass(ctx, &lastCleanup)
		wait, notified := interval, wake
		switch {
		case err != nil:
			failures++
			wait, notified = backoff(failures), nil
			r.log.Error("outbox relay pass failed", "error", err, "failures", failures, "retry_in", wait.String())
		case n == batchSize:
			failures = 0
			continue // more rows wait: pass again at once
		default:
			failures = 0
		}
		select {
		case <-ctx.Done():
		case <-r.after(wait):
		case <-notified:
		}
	}
	return nil
}

// pass runs one Publish, and the hourly Cleanup when it is due, on a context that shutdown doesn't cancel.
func (r *Relay) pass(ctx context.Context, lastCleanup *time.Time) (int, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), passTimeout)
	defer cancel()
	n, err := r.Publish(ctx)
	if now := r.now(); now.Sub(*lastCleanup) >= cleanupEvery {
		*lastCleanup = now
		if _, cerr := r.Cleanup(ctx); cerr != nil {
			r.log.Error("outbox cleanup failed", "error", cerr)
		}
	}
	return n, err
}

// backoff is the wait after the given number of failed passes in a row: 1 s, 2 s, 4 s ... at most 30 s.
func backoff(failures int) time.Duration {
	if failures > 5 { // 32 s would pass the cap
		return maxBackoff
	}
	return time.Second << (failures - 1)
}

// Publish runs one pass: it locks up to 100 unpublished rows in id order, XADDs each to the stream (trimming
// entries older than 7 days), and marks the sent rows published, in one transaction. An XADD that fails stops the
// pass at its row, which gets attempts + 1 and last_error; later rows wait so the order holds. It returns the number
// of rows published and refreshes the backlog gauges.
func (r *Relay) Publish(ctx context.Context) (int, error) {
	defer r.observeBacklog(ctx)
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin outbox pass: %w", err)
	}
	defer func() {
		// Rolls back synchronously even when ctx is cancelled, so the rows' locks go at once.
		rctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), passTimeout)
		defer cancel()
		_ = tx.Rollback(rctx)
	}()
	q := r.q.WithTx(tx)
	rows, err := q.LockUnpublishedOutbox(ctx, batchSize)
	if err != nil {
		return 0, fmt.Errorf("lock outbox rows: %w", err)
	}
	minID := strconv.FormatInt(r.now().Add(-retention).UnixMilli(), 10)
	sent := make([]pgtype.UUID, 0, len(rows))
	var xaddErr error
	for _, row := range rows {
		if err := r.xadd(ctx, row, minID); err != nil {
			publishErrors.WithLabelValues(r.stream).Inc()
			xaddErr = fmt.Errorf("publish outbox row %s: %w", row.ID, err)
			failed := store.MarkOutboxFailedParams{ID: row.ID, LastError: pgtype.Text{String: err.Error(), Valid: true}}
			if err := q.MarkOutboxFailed(ctx, failed); err != nil {
				return 0, errors.Join(xaddErr, fmt.Errorf("record the failure: %w", err))
			}
			break
		}
		sent = append(sent, row.ID)
	}
	if len(sent) > 0 {
		if err := q.MarkOutboxPublished(ctx, sent); err != nil {
			return 0, errors.Join(xaddErr, fmt.Errorf("mark outbox rows published: %w", err))
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, errors.Join(xaddErr, fmt.Errorf("commit outbox pass: %w", err))
	}
	for _, row := range rows[:len(sent)] {
		publishedTotal.WithLabelValues(r.stream, row.Type).Inc()
	}
	return len(sent), xaddErr
}

// xadd adds one row to the stream: XADD <stream> MINID ~ <minID> * type <type> envelope <json> (Q015).
func (r *Relay) xadd(ctx context.Context, row store.Outbox, minID string) error {
	env, err := json.Marshal(EnvelopeOf(row))
	if err != nil {
		return fmt.Errorf("encode envelope: %w", err)
	}
	return r.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		MinID:  minID,
		Approx: true,
		ID:     "*",
		Values: []any{"type", row.Type, "envelope", env},
	}).Err()
}

// observeBacklog sets the pending and oldest-age gauges from the outbox.
func (r *Relay) observeBacklog(ctx context.Context) {
	b, err := r.q.OutboxBacklog(ctx)
	if err != nil {
		r.log.Warn("outbox backlog unavailable", "error", err)
		return
	}
	outboxPending.Set(float64(b.Pending))
	age := 0.0
	if b.Oldest.Valid {
		age = max(r.now().Sub(b.Oldest.Time).Seconds(), 0)
	}
	outboxOldest.Set(age)
}

// Cleanup deletes up to 10,000 published rows older than 7 days and returns how many it deleted. Unpublished rows
// are never deleted.
func (r *Relay) Cleanup(ctx context.Context) (int64, error) {
	n, err := r.q.DeletePublishedOutbox(ctx, pgtype.Timestamptz{Time: r.now().Add(-retention), Valid: true})
	if err != nil {
		return 0, fmt.Errorf("delete published outbox rows: %w", err)
	}
	return n, nil
}

// listen opens the relay's own connection and subscribes it to the outbox trigger's NOTIFY.
func (r *Relay) listen(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.ConnectConfig(ctx, r.pool.Config().ConnConfig)
	if err != nil {
		return nil, fmt.Errorf("connect outbox listener: %w", err)
	}
	if err := store.New(conn).ListenOutbox(ctx); err != nil {
		closeConn(ctx, conn)
		return nil, fmt.Errorf("listen outbox: %w", err)
	}
	return conn, nil
}

// awaitNotifications turns each NOTIFY on conn into a wake-up, reconnecting when the connection fails, until ctx is
// done. conn may be nil when the first listen failed.
func (r *Relay) awaitNotifications(ctx context.Context, conn *pgx.Conn, wake chan<- struct{}) {
	for {
		if conn != nil {
			err := waitLoop(ctx, conn, wake)
			closeConn(ctx, conn)
			if ctx.Err() != nil {
				return
			}
			r.log.Warn("outbox listener lost, passing every second until it reconnects", "error", err)
		}
		retry := time.NewTimer(relistenAfter)
		select {
		case <-ctx.Done():
			retry.Stop()
			return
		case <-retry.C:
		}
		var err error
		if conn, err = r.listen(ctx); err != nil && ctx.Err() == nil {
			r.log.Warn("outbox listener unavailable, passing every second", "error", err)
		}
	}
}

func waitLoop(ctx context.Context, conn *pgx.Conn, wake chan<- struct{}) error {
	for {
		if _, err := conn.WaitForNotification(ctx); err != nil {
			return err
		}
		select {
		case wake <- struct{}{}:
		default: // a wake-up is already pending
		}
	}
}

func closeConn(ctx context.Context, conn *pgx.Conn) {
	cctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Second)
	defer cancel()
	_ = conn.Close(cctx)
}
