package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"

	"coolinary/kitchen/internal/kitchen/store"
)

const (
	// readCount is the most entries one XREADGROUP or XAUTOCLAIM returns.
	readCount = 50
	// readBlock is the longest one XREADGROUP waits for new entries.
	readBlock = 5 * time.Second
	// handlerTimeout bounds one handler's transaction.
	handlerTimeout = 10 * time.Second
	// claimEvery is how often a stream's loop claims entries left pending; claimMinIdle is how long an entry must
	// have been pending to be claimed.
	claimEvery   = time.Minute
	claimMinIdle = time.Minute
	// maxDeliveries is the delivery whose failure makes an entry a dead letter.
	maxDeliveries = 10
	// processedRetention is how long processed_events rows are kept: twice the streams' 7 days, so an id older
	// than that can't be delivered again.
	processedRetention = 14 * 24 * time.Hour
	processedCleanup   = 24 * time.Hour
	// finishTimeout bounds the rollback, XPENDING and XACK after a handler, which run even when shutdown has started.
	finishTimeout = 2 * time.Second
)

// Results of one delivery, the result label of events_consumed_total.
const (
	resultOK        = "ok"        // the handler applied the event
	resultDuplicate = "duplicate" // the event was applied before; acknowledged without the handler
	resultSkipped   = "skipped"   // no handler for the type; acknowledged unread
	resultError     = "error"     // the handler failed or the entry can't be handled yet; left pending for a retry
	resultDead      = "dead"      // the 10th delivery failed too; acknowledged as a dead letter
)

// Metric names are the same in every service that consumes events (Q041).
var (
	consumedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "events_consumed_total",
		Help: "Stream entries delivered to the consumer, by what became of them.",
	}, []string{"stream", "type", "result"})
	handlerDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "events_handler_duration_seconds",
		Help:    "Time an event handler ran.",
		Buckets: prometheus.DefBuckets,
	}, []string{"type"})
)

// errNewerVersion: the envelope's version is above the one its handler understands.
var errNewerVersion = errors.New("event version newer than its handler")

// Handler applies one event. It runs in the consumer's transaction: q is that transaction's Queries (so events.Add
// works with it too), and returning an error rolls back everything it wrote. ctx ends after 10 seconds. A handler
// that calls another service must be idempotent itself: pass e.ID as the request's idempotency key where the callee
// takes one.
type Handler func(ctx context.Context, q *store.Queries, e Envelope) error

type handler struct {
	version int32
	fn      Handler
}

// ConsumerOptions configure the consumer. Streams, Group and Name are required; the rest is for tests, and their
// zero values are what the service runs with.
type ConsumerOptions struct {
	// Streams to read (EVENT_STREAMS), e.g. events.identity.
	Streams []string
	// Group is the consumer group: the service name.
	Group string
	// Name is this consumer in the group: the container hostname.
	Name string
	// Now reads the clock for the processed_events cutoff; time.Now when nil.
	Now func() time.Time
	// After waits like time.After (the default): for the next claim, the next cleanup and a backoff.
	After func(time.Duration) <-chan time.Time
	// Block is the longest one read waits for new entries; 5 s when zero.
	Block time.Duration
}

// Consumer applies other services' events exactly once in effect, in stream order (KIT-01-04). Register handlers
// with Handle, then run it as one g.Go(...) loop in main.go.
//
// Each stream has one loop. On start it creates the group at the stream's beginning when missing (so a new service
// processes everything still retained), handles this consumer's own pending entries (read before a restart, not yet
// acknowledged), claims entries pending over a minute, then reads new entries 50 at a time and handles them one by
// one in stream order; every minute it claims again. A handled entry's handler runs in one transaction that also
// records the event id in processed_events, and the entry is acknowledged only after the commit; an id already
// recorded means the event was applied before, so the entry is acknowledged without running the handler. A failed
// entry stays pending until a claim retries it; when its 10th delivery fails it is acknowledged as a dead letter.
// Entries of unhandled types are acknowledged without parsing their envelope.
type Consumer struct {
	pool     *pgxpool.Pool
	q        *store.Queries
	rdb      *redis.Client
	log      *slog.Logger
	streams  []string
	group    string
	name     string
	now      func() time.Time
	after    func(time.Duration) <-chan time.Time
	block    time.Duration
	handlers map[string]handler
}

// NewConsumer builds the consumer over the service's pool and its redis-state client.
func NewConsumer(pool *pgxpool.Pool, rdb *redis.Client, log *slog.Logger, opts ConsumerOptions) *Consumer {
	c := &Consumer{
		pool: pool, q: store.New(pool), rdb: rdb, log: log,
		streams: opts.Streams, group: opts.Group, name: opts.Name,
		now: opts.Now, after: opts.After, block: opts.Block,
		handlers: map[string]handler{},
	}
	if c.now == nil {
		c.now = time.Now
	}
	if c.after == nil {
		c.after = time.After
	}
	if c.block == 0 {
		c.block = readBlock
	}
	return c
}

// Handle registers fn for events of eventType up to the given major version: consumer.Handle("household.plan_changed",
// 1, fn). Entries with a newer version aren't acknowledged; they are retried and end as dead letters, so the
// incompatibility shows. Call Handle before Run.
func (c *Consumer) Handle(eventType string, version int32, fn Handler) {
	c.handlers[eventType] = handler{version: version, fn: fn}
}

// Run consumes every stream until ctx is done, then returns nil once each loop finished the handler it was running
// (rule 11: stop reading, finish the current handler, return). Entries read but not yet handled stay pending and are
// handled first on the next start. It also deletes processed_events rows older than 14 days, at start and daily.
// Failures are logged and counted, never returned.
func (c *Consumer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, stream := range c.streams {
		wg.Go(func() { c.consume(ctx, stream) })
	}
	wg.Go(func() { c.cleanupDaily(ctx) })
	wg.Wait()
	return nil
}

// streamLoop is where one stream's loop stands.
type streamLoop struct {
	stream string
	// started: the group exists and this consumer's own pending entries were handled.
	started bool
	// claimDue fires when the next claim is due; nil means at once.
	claimDue <-chan time.Time
}

// consume runs one stream's loop until ctx is done. A Redis or Postgres error waits 1 s, 2 s, 4 s ... (at most 30 s)
// before the next round; a missing group (the stream was deleted) is created again.
func (c *Consumer) consume(ctx context.Context, stream string) {
	l := &streamLoop{stream: stream}
	failures := 0
	for ctx.Err() == nil {
		err := c.round(ctx, l)
		if err == nil || ctx.Err() != nil {
			failures = 0
			continue
		}
		if isNoGroup(err) {
			l.started = false
		}
		failures++
		wait := backoff(failures)
		c.log.Error("event stream read failed", "stream", stream, "error", err, "failures", failures, "retry_in", wait.String())
		select {
		case <-ctx.Done():
		case <-c.after(wait):
		}
	}
}

// round starts the loop when needed, claims when a claim is due, then reads and handles new entries.
func (c *Consumer) round(ctx context.Context, l *streamLoop) error {
	if !l.started {
		if err := c.createGroup(ctx, l.stream); err != nil {
			return err
		}
		if err := c.resume(ctx, l.stream); err != nil {
			return err
		}
		l.started, l.claimDue = true, nil
	}
	if claimDue(l.claimDue) {
		l.claimDue = nil // retried in the next round if the claim fails
		if err := c.claim(ctx, l.stream); err != nil {
			return err
		}
		l.claimDue = c.after(claimEvery)
	}
	return c.readNew(ctx, l.stream)
}

func claimDue(due <-chan time.Time) bool {
	if due == nil {
		return true
	}
	select {
	case <-due:
		return true
	default:
		return false
	}
}

// createGroup creates the group at the stream's first entry, and the stream when missing:
// XGROUP CREATE <stream> <group> 0 MKSTREAM. An existing group is kept.
func (c *Consumer) createGroup(ctx context.Context, stream string) error {
	err := c.rdb.XGroupCreateMkStream(ctx, stream, c.group, "0").Err()
	if err != nil && !strings.HasPrefix(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("create consumer group: %w", err)
	}
	return nil
}

// resume handles, in stream order, the entries this consumer read but didn't acknowledge before it stopped, so a
// restart keeps the publish order.
func (c *Consumer) resume(ctx context.Context, stream string) error {
	after := "0"
	for ctx.Err() == nil {
		msgs, err := c.read(ctx, stream, after, -1)
		if err != nil || len(msgs) == 0 {
			return err
		}
		c.handleAll(ctx, stream, msgs)
		after = msgs[len(msgs)-1].ID
	}
	return nil
}

// claim takes over, in stream order, the group's entries pending for over a minute (their consumer failed or died)
// and handles them again.
func (c *Consumer) claim(ctx context.Context, stream string) error {
	start := "0-0"
	for ctx.Err() == nil {
		msgs, next, err := c.rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
			Stream: stream, Group: c.group, Consumer: c.name, MinIdle: claimMinIdle, Start: start, Count: readCount,
		}).Result()
		if err != nil {
			return fmt.Errorf("claim pending entries: %w", err)
		}
		c.handleAll(ctx, stream, msgs)
		if next == "0-0" {
			return nil
		}
		start = next
	}
	return nil
}

// readNew reads up to 50 entries no consumer of the group has seen, waiting up to Block, and handles them.
func (c *Consumer) readNew(ctx context.Context, stream string) error {
	msgs, err := c.read(ctx, stream, ">", c.block)
	if err != nil {
		return err
	}
	c.handleAll(ctx, stream, msgs)
	return nil
}

// read is XREADGROUP GROUP <group> <name> COUNT 50 [BLOCK <ms>] STREAMS <stream> <id>; block < 0 doesn't block.
func (c *Consumer) read(ctx context.Context, stream, id string, block time.Duration) ([]redis.XMessage, error) {
	res, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group: c.group, Consumer: c.name, Streams: []string{stream, id}, Count: readCount, Block: block,
	}).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return nil, nil // nothing new within Block
	case err != nil:
		return nil, fmt.Errorf("read stream: %w", err) // after shutdown started, consume drops it
	case len(res) == 0:
		return nil, nil
	}
	return res[0].Messages, nil
}

// handleAll handles msgs one at a time in order and stops early when shutdown starts.
func (c *Consumer) handleAll(ctx context.Context, stream string, msgs []redis.XMessage) {
	for _, m := range msgs {
		if ctx.Err() != nil {
			return
		}
		c.handle(ctx, stream, m)
	}
}

// handle takes one delivered entry to its result. It finishes even when shutdown starts meanwhile.
func (c *Consumer) handle(ctx context.Context, stream string, m redis.XMessage) {
	ctx = context.WithoutCancel(ctx)
	typ, _ := m.Values["type"].(string)
	h, ok := c.handlers[typ]
	if !ok {
		c.ack(ctx, stream, m.ID, typ, "")
		consumedTotal.WithLabelValues(stream, typ, resultSkipped).Inc()
		return
	}
	raw, _ := m.Values["envelope"].(string)
	eventID, applied, err := c.apply(ctx, stream, typ, h, raw)
	if err != nil {
		c.failed(ctx, stream, m.ID, typ, eventID, err)
		return
	}
	c.ack(ctx, stream, m.ID, typ, eventID)
	result := resultDuplicate
	if applied {
		result = resultOK
	}
	consumedTotal.WithLabelValues(stream, typ, result).Inc()
}

// apply decodes the envelope and, unless its id is recorded already, runs the handler and records the id in one
// transaction. It returns the envelope id ("" until it is known to be a UUID) and whether the handler applied it.
func (c *Consumer) apply(ctx context.Context, stream, typ string, h handler, raw string) (string, bool, error) {
	var env Envelope
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		return "", false, fmt.Errorf("decode envelope: %w", err)
	}
	var id pgtype.UUID
	if err := id.Scan(env.ID); err != nil {
		return "", false, errors.New("envelope id is not a UUID")
	}
	eventID := id.String()
	if env.Version > h.version {
		return eventID, false, fmt.Errorf("%w: version %d, handler takes up to %d", errNewerVersion, env.Version, h.version)
	}
	ctx, cancel := context.WithTimeout(ctx, handlerTimeout)
	defer cancel()
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return eventID, false, fmt.Errorf("begin handler transaction: %w", err)
	}
	defer func() {
		// Rolls back even after a timeout, so the row locks go at once; a no-op after Commit.
		rctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), finishTimeout)
		defer cancel()
		_ = tx.Rollback(rctx)
	}()
	q := c.q.WithTx(tx)
	n, err := q.InsertProcessedEvent(ctx, store.InsertProcessedEventParams{EventID: id, Stream: stream, Type: typ})
	if err != nil {
		return eventID, false, fmt.Errorf("record processed event: %w", err)
	}
	if n == 0 {
		return eventID, false, nil // applied before: roll back, acknowledge
	}
	started := time.Now()
	err = h.fn(ctx, q, env)
	handlerDuration.WithLabelValues(typ).Observe(time.Since(started).Seconds())
	if err != nil {
		return eventID, false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return eventID, false, fmt.Errorf("commit handler transaction: %w", err)
	}
	return eventID, true, nil
}

// failed leaves a failed entry pending for a retry, or acknowledges it as a dead letter when this was its 10th
// delivery. The log carries ids, the type and the error, never the payload.
func (c *Consumer) failed(ctx context.Context, stream, entryID, typ, eventID string, err error) {
	deliveries, perr := c.deliveries(ctx, stream, entryID)
	if perr != nil {
		c.log.Warn("event delivery count unavailable, retrying later", "stream", stream, "entry_id", entryID, "error", perr)
	}
	if perr == nil && deliveries >= maxDeliveries {
		c.log.Error("event handler failed for the last time; dropped as a dead letter", "stream", stream,
			"entry_id", entryID, "event_id", eventID, "type", typ, "error", err, "deliveries", deliveries)
		c.ack(ctx, stream, entryID, typ, eventID)
		consumedTotal.WithLabelValues(stream, typ, resultDead).Inc()
		return
	}
	c.log.Warn("event handler failed; retrying later", "stream", stream, "entry_id", entryID, "event_id", eventID,
		"type", typ, "error", err, "deliveries", deliveries)
	consumedTotal.WithLabelValues(stream, typ, resultError).Inc()
}

// deliveries is how often the entry was delivered, counting this delivery (from XPENDING).
func (c *Consumer) deliveries(ctx context.Context, stream, entryID string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, finishTimeout)
	defer cancel()
	p, err := c.rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: stream, Group: c.group, Start: entryID, End: entryID, Count: 1,
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("read delivery count: %w", err)
	}
	if len(p) == 0 {
		return 0, errors.New("entry no longer pending")
	}
	return p[0].RetryCount, nil
}

// ack acknowledges an entry. When XACK fails the entry stays pending; its retry finds the event recorded and only
// acknowledges it.
func (c *Consumer) ack(ctx context.Context, stream, entryID, typ, eventID string) {
	ctx, cancel := context.WithTimeout(ctx, finishTimeout)
	defer cancel()
	if err := c.rdb.XAck(ctx, stream, c.group, entryID).Err(); err != nil {
		c.log.Warn("event acknowledgement failed; it will be delivered again", "stream", stream,
			"entry_id", entryID, "event_id", eventID, "type", typ, "error", err)
	}
}

// cleanupDaily deletes processed_events rows older than 14 days, at start and then every 24 hours, until ctx is done.
func (c *Consumer) cleanupDaily(ctx context.Context) {
	for {
		cutoff := pgtype.Timestamptz{Time: c.now().Add(-processedRetention), Valid: true}
		if _, err := c.q.DeleteProcessedEvents(ctx, cutoff); err != nil && ctx.Err() == nil {
			c.log.Error("processed events cleanup failed", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-c.after(processedCleanup):
		}
	}
}

// isNoGroup reports Redis's NOGROUP error: the group (or the whole stream) is gone.
func isNoGroup(err error) bool {
	var rerr redis.Error
	return errors.As(err, &rerr) && strings.HasPrefix(rerr.Error(), "NOGROUP")
}
