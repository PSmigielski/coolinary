package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"

	"coolinary/planner/internal/planner/store"
)

const (
	// readCount is the most entries one XREADGROUP or XAUTOCLAIM returns.
	readCount = 100
	// readBlock is the longest one XREADGROUP waits for new entries.
	readBlock = 5 * time.Second
	// handlerTimeout bounds one handler's transaction.
	handlerTimeout = 10 * time.Second
	// claimEvery is how often a stream's loop claims entries left pending (and refreshes the lag gauge);
	// claimMinIdle is how long an entry must have been pending to be claimed.
	claimEvery   = time.Minute
	claimMinIdle = time.Minute
	// maxDeliveries is the delivery whose failure parks an entry in dead_events.
	maxDeliveries = 5
	// After a failed round (Redis unreachable) the loop waits retryWait plus up to retryJitter, so planner's
	// containers don't retry in step, and logs the failure at most once per logEvery.
	retryWait   = 2 * time.Second
	retryJitter = 500 * time.Millisecond
	logEvery    = time.Minute
	// processedRetention is how long processed_events rows are kept: twice the streams' 7 days, so an id older
	// than that can't be delivered again. The cleanup runs daily at cleanupHour:cleanupMinute UTC.
	processedRetention = 14 * 24 * time.Hour
	cleanupHour        = 3
	cleanupMinute      = 30
	// finishTimeout bounds the rollback, XPENDING, dead_events insert and XACK after a handler, which run even when
	// shutdown has started.
	finishTimeout = 2 * time.Second
)

// Results of one delivery, the result label of events_consumed_total (kitchen's labels, Q235).
const (
	resultOK        = "ok"        // the handler applied the event
	resultDuplicate = "duplicate" // the event was applied before; acknowledged without the handler
	resultSkipped   = "skipped"   // no handler for the type; acknowledged unread
	resultError     = "error"     // the handler failed or the entry can't be handled yet; left pending for a retry
	resultDead      = "dead"      // given up on: parked in dead_events and acknowledged
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
	consumerLag = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "events_consumer_lag",
		Help: "Stream entries not yet delivered to the consumer group (XINFO GROUPS lag), refreshed every minute.",
	}, []string{"stream", "group"})
)

var (
	// errNewerVersion: the envelope's version is above the one its handler understands.
	errNewerVersion = errors.New("event version newer than its handler")
	// errMalformed: the envelope can't be decoded or has no UUID id, so it can be neither handled nor recorded.
	errMalformed = errors.New("malformed envelope")
)

// Handler applies one event. It runs in the consumer's transaction: q is that transaction's Queries, and returning
// an error rolls back everything it wrote. ctx ends after 10 seconds. A handler must not call another service: when
// an event needs one, it enqueues the work (e.g. a row in a job table) and returns.
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
	// Now reads the clock for the processed_events cleanup and the error log's pace; time.Now when nil.
	Now func() time.Time
	// After waits like time.After (the default): for the next claim, the next cleanup and a retry.
	After func(time.Duration) <-chan time.Time
	// Block is the longest one read waits for new entries; 5 s when zero.
	Block time.Duration
	// Jitter returns a random duration in [0, max), added to the wait before a retry; math/rand/v2's N when nil.
	Jitter func(max time.Duration) time.Duration
}

// Consumer applies other services' events exactly once in effect, in stream order (PLN-01-01). Register handlers
// with Handle, then run it as one g.Go(...) loop in main.go.
//
// Each stream has one loop. On start it creates the group at the stream's end when missing (a new planner starts
// from the events published after it), handles this consumer's own pending entries (read before a restart, not yet
// acknowledged), claims entries pending over a minute, then reads new entries 100 at a time and handles them one by
// one in stream order; every minute it claims again and refreshes events_consumer_lag. A handled entry's handler
// runs in one transaction that also records the event id in processed_events, and the entry is acknowledged only
// after the commit; an id already recorded means the event was applied before, so the entry is acknowledged without
// running the handler. A failed entry stays pending until a claim retries it; when its 5th delivery fails it is
// parked in dead_events and acknowledged. An envelope that is malformed or newer than its handler is parked at once.
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
	jitter   func(time.Duration) time.Duration
	handlers map[string]handler
}

// NewConsumer builds the consumer over the service's pool and its redis-state client.
func NewConsumer(pool *pgxpool.Pool, rdb *redis.Client, log *slog.Logger, opts ConsumerOptions) *Consumer {
	c := &Consumer{
		pool: pool, q: store.New(pool), rdb: rdb, log: log,
		streams: opts.Streams, group: opts.Group, name: opts.Name,
		now: opts.Now, after: opts.After, block: opts.Block, jitter: opts.Jitter,
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
	if c.jitter == nil {
		c.jitter = rand.N[time.Duration]
	}
	return c
}

// Handle registers fn for events of eventType up to the given major version: consumer.Handle("household.plan_changed",
// 1, fn). Entries with a newer version are parked in dead_events at once, so the incompatibility shows. Call Handle
// before Run.
func (c *Consumer) Handle(eventType string, version int32, fn Handler) {
	c.handlers[eventType] = handler{version: version, fn: fn}
}

// Run consumes every stream until ctx is done, then returns nil once each loop finished the handler it was running
// (rule 8: finish the current message, take no new ones). Entries read but not yet handled stay pending and are
// handled first on the next start. It also deletes processed_events rows older than 14 days, daily at 03:30 UTC.
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
	// groupStart is where a missing group is created: "$" on start, "0" once the group vanished while running.
	groupStart string
	// started: the group exists and this consumer's own pending entries were handled.
	started bool
	// claimDue fires when the next claim is due; nil means at once.
	claimDue <-chan time.Time
}

// consume runs one stream's loop until ctx is done. A failed round (Redis unreachable) is retried after 2 s plus up
// to 0.5 s of jitter and logged at most once a minute. A missing group means the stream was deleted while planner
// ran: the group is created again at the stream's beginning, since everything in the new stream is new to planner.
func (c *Consumer) consume(ctx context.Context, stream string) {
	l := &streamLoop{stream: stream, groupStart: "$"}
	failures := 0
	var loggedAt time.Time
	for ctx.Err() == nil {
		err := c.round(ctx, l)
		if err == nil || ctx.Err() != nil {
			failures = 0
			continue
		}
		if isNoGroup(err) {
			l.started, l.groupStart = false, "0"
		}
		failures++
		if now := c.now(); loggedAt.IsZero() || now.Sub(loggedAt) >= logEvery {
			loggedAt = now
			c.log.Error("event stream read failed", "stream", stream, "error", err, "failures", failures)
		}
		select {
		case <-ctx.Done():
		case <-c.after(retryWait + c.jitter(retryJitter)):
		}
	}
}

// round starts the loop when needed, claims and refreshes the lag when due, then reads and handles new entries.
func (c *Consumer) round(ctx context.Context, l *streamLoop) error {
	if !l.started {
		if err := c.createGroup(ctx, l.stream, l.groupStart); err != nil {
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
		if err := c.observeLag(ctx, l.stream); err != nil {
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

// createGroup creates the group at start ("$" = the stream's end, "0" = its beginning), and the stream when missing:
// XGROUP CREATE <stream> <group> <start> MKSTREAM. An existing group is kept.
func (c *Consumer) createGroup(ctx context.Context, stream, start string) error {
	err := c.rdb.XGroupCreateMkStream(ctx, stream, c.group, start).Err()
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

// observeLag sets events_consumer_lag to the entries not yet delivered to the group (XINFO GROUPS). When Redis can't
// tell (lag -1, e.g. after entries were deleted from the middle of the stream) the gauge keeps its value.
func (c *Consumer) observeLag(ctx context.Context, stream string) error {
	groups, err := c.rdb.XInfoGroups(ctx, stream).Result()
	if err != nil {
		return fmt.Errorf("read consumer group lag: %w", err)
	}
	for _, g := range groups {
		if g.Name == c.group && g.Lag >= 0 {
			consumerLag.WithLabelValues(stream, c.group).Set(float64(g.Lag))
		}
	}
	return nil
}

// readNew reads up to 100 entries no consumer of the group has seen, waiting up to Block, and handles them.
func (c *Consumer) readNew(ctx context.Context, stream string) error {
	msgs, err := c.read(ctx, stream, ">", c.block)
	if err != nil {
		return err
	}
	c.handleAll(ctx, stream, msgs)
	return nil
}

// read is XREADGROUP GROUP <group> <name> COUNT 100 [BLOCK <ms>] STREAMS <stream> <id>; block < 0 doesn't block.
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
		c.failed(ctx, stream, m.ID, typ, eventID, raw, err)
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
// Decode errors name the field at most, never the envelope's content.
func (c *Consumer) apply(ctx context.Context, stream, typ string, h handler, raw string) (string, bool, error) {
	var env Envelope
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		if te, ok := errors.AsType[*json.UnmarshalTypeError](err); ok && te.Field != "" {
			return "", false, fmt.Errorf("%w: %s has the wrong type", errMalformed, te.Field)
		}
		return "", false, fmt.Errorf("%w: not JSON", errMalformed)
	}
	var id pgtype.UUID
	if err := id.Scan(env.ID); err != nil {
		return "", false, fmt.Errorf("%w: id is not a UUID", errMalformed)
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

// failed leaves a failed entry pending for a retry, or parks it in dead_events when it can never be handled
// (malformed, newer version) or this was its 5th delivery. The log carries ids, the type and the error, never the
// payload.
func (c *Consumer) failed(ctx context.Context, stream, entryID, typ, eventID, raw string, err error) {
	deliveries, perr := c.deliveries(ctx, stream, entryID)
	if perr != nil {
		c.log.Warn("event delivery count unavailable, retrying later", "stream", stream, "entry_id", entryID, "error", perr)
	} else if errors.Is(err, errMalformed) || errors.Is(err, errNewerVersion) || deliveries >= maxDeliveries {
		c.park(ctx, stream, entryID, typ, eventID, raw, deliveries, err)
		return
	}
	c.log.Warn("event handler failed; retrying later", "stream", stream, "entry_id", entryID, "event_id", eventID,
		"type", typ, "error", err, "deliveries", deliveries)
	consumedTotal.WithLabelValues(stream, typ, resultError).Inc()
}

// park copies an entry the consumer gives up on to dead_events with the error that ended it, then acknowledges it.
// When the row can't be written the entry stays pending, and a claim brings it back. The log line carries ids, the
// type and the error, never the envelope.
func (c *Consumer) park(ctx context.Context, stream, entryID, typ, eventID, raw string, deliveries int64, cause error) {
	var id pgtype.UUID
	if eventID != "" {
		_ = id.Scan(eventID) // apply set it from a parsed UUID
	}
	ictx, cancel := context.WithTimeout(ctx, finishTimeout)
	defer cancel()
	err := c.q.InsertDeadEvent(ictx, store.InsertDeadEventParams{
		Stream: stream, EntryID: entryID, EventID: id, Type: textOf(typ), Envelope: textOf(raw),
		LastError: textOf(cause.Error()), Deliveries: int32(min(deliveries, math.MaxInt32)), //nolint:gosec // G115: capped at MaxInt32 just before
	})
	if err != nil {
		c.log.Warn("event could not be parked in dead_events; retrying later", "stream", stream, "entry_id", entryID,
			"event_id", eventID, "type", typ, "error", err)
		consumedTotal.WithLabelValues(stream, typ, resultError).Inc()
		return
	}
	c.log.Error("event parked in dead_events", "stream", stream, "entry_id", entryID, "event_id", eventID,
		"type", typ, "error", cause, "deliveries", deliveries)
	c.ack(ctx, stream, entryID, typ, eventID)
	consumedTotal.WithLabelValues(stream, typ, resultDead).Inc()
}

// textOf makes s storable as Postgres text, which holds neither NUL bytes nor invalid UTF-8: NULs are dropped and
// invalid bytes become U+FFFD.
func textOf(s string) string {
	return strings.ToValidUTF8(strings.ReplaceAll(s, "\x00", ""), "�")
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

// ack acknowledges an entry. When XACK fails the entry stays pending; its retry finds the event recorded (or the
// entry parked) and only acknowledges it.
func (c *Consumer) ack(ctx context.Context, stream, entryID, typ, eventID string) {
	ctx, cancel := context.WithTimeout(ctx, finishTimeout)
	defer cancel()
	if err := c.rdb.XAck(ctx, stream, c.group, entryID).Err(); err != nil {
		c.log.Warn("event acknowledgement failed; it will be delivered again", "stream", stream,
			"entry_id", entryID, "event_id", eventID, "type", typ, "error", err)
	}
}

// cleanupDaily deletes processed_events rows older than 14 days every day at 03:30 UTC, until ctx is done.
func (c *Consumer) cleanupDaily(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.after(untilCleanup(c.now())):
		}
		cutoff := pgtype.Timestamptz{Time: c.now().Add(-processedRetention), Valid: true}
		if _, err := c.q.DeleteProcessedEvents(ctx, cutoff); err != nil && ctx.Err() == nil {
			c.log.Error("processed events cleanup failed", "error", err)
		}
	}
}

// untilCleanup is the time from now to the next 03:30 UTC; a full day when it is 03:30 now.
func untilCleanup(now time.Time) time.Duration {
	now = now.UTC()
	next := time.Date(now.Year(), now.Month(), now.Day(), cleanupHour, cleanupMinute, 0, 0, time.UTC)
	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	return next.Sub(now)
}

// isNoGroup reports Redis's NOGROUP error: the group (or the whole stream) is gone.
func isNoGroup(err error) bool {
	var rerr redis.Error
	return errors.As(err, &rerr) && strings.HasPrefix(rerr.Error(), "NOGROUP")
}
