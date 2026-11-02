package events_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/events"
)

const (
	mergedType = "catalogue.ingredient_merged"
	day        = 24 * time.Hour
	week       = 7 * day
)

// errRedisDown is what a stopped redis-state looks like to the relay.
var errRedisDown = errors.New("dial tcp 10.77.0.12:6379: connect: connection refused")

// cleanupNow is the relay's clock in the cleanup tests: years before any real row, so only the test's own rows are
// old enough to delete.
var cleanupNow = time.Date(2020, 1, 15, 4, 0, 0, 0, time.UTC)

// openTestRedis returns a client on redis-state's test database whose XADDs the test can break, or skips when
// TEST_REDIS_STATE_URL isn't set (make test sets it while make dev is up).
func openTestRedis(t *testing.T) (*redis.Client, *xaddFaults) {
	t.Helper()
	url := os.Getenv("TEST_REDIS_STATE_URL")
	if url == "" {
		t.Skip("needs TEST_REDIS_STATE_URL (make dev, then make test s=kitchen)")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	rdb := redis.NewClient(opts)
	t.Cleanup(func() { _ = rdb.Close() })
	faults := &xaddFaults{}
	rdb.AddHook(faults)
	return rdb, faults
}

// xaddFaults breaks XADDs the way a stopped Redis or a dying process would. before runs ahead of each XADD with the
// envelope id it carries, and an error from it fails that XADD without sending it; after runs once an XADD succeeded.
type xaddFaults struct {
	mu     sync.Mutex
	before func(ctx context.Context, envelopeID string) error
	after  func(envelopeID string)
}

func (f *xaddFaults) set(before func(context.Context, string) error, after func(string)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.before, f.after = before, after
}

func (f *xaddFaults) DialHook(next redis.DialHook) redis.DialHook { return next }

func (f *xaddFaults) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

func (f *xaddFaults) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if cmd.Name() != "xadd" {
			return next(ctx, cmd)
		}
		f.mu.Lock()
		before, after := f.before, f.after
		f.mu.Unlock()
		id := xaddEnvelopeID(cmd.Args())
		if before != nil {
			if err := before(ctx, id); err != nil {
				cmd.SetErr(err)
				return err
			}
		}
		if err := next(ctx, cmd); err != nil {
			return err
		}
		if after != nil {
			after(id)
		}
		return nil
	}
}

// xaddEnvelopeID reads the envelope id from an XADD's arguments; "" when it carries no envelope.
func xaddEnvelopeID(args []any) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] != "envelope" {
			continue
		}
		var raw []byte
		switch v := args[i+1].(type) {
		case []byte:
			raw = v
		case string:
			raw = []byte(v)
		}
		var env struct {
			ID string `json:"id"`
		}
		if json.Unmarshal(raw, &env) == nil {
			return env.ID
		}
	}
	return ""
}

func redisDown(context.Context, string) error { return errRedisDown }

// failXAddOf fails the XADD of each envelope id in ids with err; other XADDs go through.
func failXAddOf(ids []string, err error) func(context.Context, string) error {
	return func(_ context.Context, id string) error {
		if slices.Contains(ids, id) {
			return err
		}
		return nil
	}
}

// testStream names a fresh stream for one test and deletes it when the test ends.
func testStream(t *testing.T, rdb *redis.Client) string {
	t.Helper()
	s := "test.events.kitchen." + randomUUID(t)
	t.Cleanup(func() { _ = rdb.Del(context.Background(), s).Err() })
	return s
}

// entry is one stream entry as a consumer reads it.
type entry struct {
	streamID string
	fields   map[string]any
	envelope string // the envelope field: the JSON the relay sent
	id       string // the envelope id
}

func readStream(t *testing.T, rdb *redis.Client, stream string) []entry {
	t.Helper()
	msgs, err := rdb.XRange(context.Background(), stream, "-", "+").Result()
	if err != nil {
		t.Fatal(err)
	}
	out := make([]entry, 0, len(msgs))
	for _, m := range msgs {
		env, _ := m.Values["envelope"].(string)
		var e struct {
			ID string `json:"id"`
		}
		_ = json.Unmarshal([]byte(env), &e)
		out = append(out, entry{streamID: m.ID, fields: m.Values, envelope: env, id: e.ID})
	}
	return out
}

// entriesWith keeps the entries whose envelope id is one of ids, in stream order.
func entriesWith(entries []entry, ids []string) []entry {
	var out []entry
	for _, e := range entries {
		if slices.Contains(ids, e.id) {
			out = append(out, e)
		}
	}
	return out
}

func envelopeIDs(entries []entry) []string {
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.id)
	}
	return out
}

// streamIDsOf returns, in stream order, the envelope ids of the entries carrying one of ids.
func streamIDsOf(t *testing.T, rdb *redis.Client, stream string, ids []string) []string {
	t.Helper()
	return envelopeIDs(entriesWith(readStream(t, rdb, stream), ids))
}

// waitForEntries waits until the stream holds n entries carrying one of ids and returns them; it fails after 10 s.
func waitForEntries(t *testing.T, rdb *redis.Client, stream string, ids []string, n int) []entry {
	t.Helper()
	timeout := time.After(10 * time.Second)
	tick := time.NewTicker(5 * time.Millisecond)
	defer tick.Stop()
	for {
		if got := entriesWith(readStream(t, rdb, stream), ids); len(got) >= n {
			return got
		}
		select {
		case <-timeout:
			t.Fatalf("the stream got fewer than %d of the test's events within 10 s", n)
		case <-tick.C:
		}
	}
}

// assertIDs compares two id lists in order and reports the first difference.
func assertIDs(t *testing.T, what string, got, want []string) {
	t.Helper()
	if slices.Equal(got, want) {
		return
	}
	i := 0
	for i < len(got) && i < len(want) && got[i] == want[i] {
		i++
	}
	t.Errorf("%s: got %d ids, want %d; first difference at position %d", what, len(got), len(want), i)
}

// commitEvents commits one transaction with a catalogue.ingredient_merged event from `from` into each of intos.
func commitEvents(ctx context.Context, pool *pgxpool.Pool, q *store.Queries, from string, intos []string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(context.WithoutCancel(ctx)) }()
	for _, into := range intos {
		if err := events.Add(ctx, q.WithTx(tx), merged(from, into)); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func randomUUIDs(t *testing.T, n int) []string {
	t.Helper()
	out := make([]string, n)
	for i := range out {
		out[i] = randomUUID(t)
	}
	return out
}

// writeEvents commits n events tagged with from in one transaction.
func writeEvents(t *testing.T, pool *pgxpool.Pool, q *store.Queries, from string, n int) {
	t.Helper()
	if err := commitEvents(t.Context(), pool, q, from, randomUUIDs(t, n)); err != nil {
		t.Fatal(err)
	}
}

// outboxIDs returns the ids of the outbox rows tagged with from, in id order.
func outboxIDs(t *testing.T, pool *pgxpool.Pool, from string) []string {
	t.Helper()
	rows := outboxRows(t, pool, from)
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.ID.String())
	}
	return out
}

// insertOutboxRow stores a row tagged with from directly, with the given times (publishedAt nil = unpublished).
func insertOutboxRow(t *testing.T, pool *pgxpool.Pool, from string, occurredAt time.Time, publishedAt *time.Time) string {
	t.Helper()
	id := randomUUID(t)
	if _, err := pool.Exec(t.Context(), `
		INSERT INTO outbox (id, type, payload, occurred_at, published_at)
		VALUES ($1, 'catalogue.ingredient_merged', jsonb_build_object('from_ingredient_id', $2::text), $3, $4)`,
		id, from, occurredAt, publishedAt); err != nil {
		t.Fatal(err)
	}
	return id
}

func setPublishedAt(t *testing.T, pool *pgxpool.Pool, from string, at *time.Time) {
	t.Helper()
	if _, err := pool.Exec(t.Context(), "UPDATE outbox SET published_at = $2 WHERE payload->>'from_ingredient_id' = $1", from, at); err != nil {
		t.Fatal(err)
	}
}

// rowState is what the relay changes on an outbox row.
type rowState struct {
	published bool
	attempts  int32
	lastError string
}

func rowStates(t *testing.T, pool *pgxpool.Pool, from string) []rowState {
	t.Helper()
	var out []rowState
	for _, r := range outboxRows(t, pool, from) {
		out = append(out, rowState{published: r.PublishedAt.Valid, attempts: r.Attempts, lastError: r.LastError.String})
	}
	return out
}

func newRelay(pool *pgxpool.Pool, rdb *redis.Client, opts events.RelayOptions) *events.Relay {
	return events.NewRelay(pool, rdb, slog.New(slog.DiscardHandler), opts)
}

// drain runs passes until one publishes nothing.
func drain(t *testing.T, r *events.Relay) {
	t.Helper()
	for range 100 {
		n, err := r.Publish(t.Context())
		if err != nil {
			t.Fatalf("Publish error = %v, want none", err)
		}
		if n == 0 {
			return
		}
	}
	t.Fatal("the outbox didn't drain in 100 passes")
}

// runRelay starts r.Run. stop shuts it down (cancels its context); done delivers what Run returned. The test's end
// stops it too.
func runRelay(t *testing.T, r *events.Relay) (stop context.CancelFunc, done <-chan error) {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	result, finished := make(chan error, 1), make(chan struct{})
	go func() {
		result <- r.Run(ctx)
		close(finished)
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case <-finished:
		case <-time.After(10 * time.Second):
			t.Error("Run did not return within 10 s of shutdown")
		}
	})
	return cancel, result
}

// fakeTimer stands in for time.After in the relay loop. Each wait the loop starts is queued in waits and fires only
// when the test fires it. The loop asks for a wait only after it finished its pass (and cleanup), so receiving a
// wait means that iteration is done.
type fakeTimer struct{ waits chan *fakeWait }

type fakeWait struct {
	d time.Duration
	c chan time.Time
}

func newFakeTimer() *fakeTimer { return &fakeTimer{waits: make(chan *fakeWait, 1000)} }

func (f *fakeTimer) After(d time.Duration) <-chan time.Time {
	w := &fakeWait{d: d, c: make(chan time.Time, 1)}
	f.waits <- w
	return w.c
}

// next returns the next wait the loop asked for; it fails after 10 s.
func (f *fakeTimer) next(t *testing.T) *fakeWait {
	t.Helper()
	select {
	case w := <-f.waits:
		return w
	case <-time.After(10 * time.Second):
		t.Fatal("the relay started no wait within 10 s")
		return nil
	}
}

func (w *fakeWait) fire() { w.c <- time.Time{} }

// clock is a settable relay clock.
type clock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t
}

// scrapeMetrics reads /metrics as the admin server serves it: "name{labels}" to value.
func scrapeMetrics(t *testing.T) map[string]float64 {
	t.Helper()
	rec := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/metrics", http.NoBody))
	out := map[string]float64{}
	for line := range strings.SplitSeq(rec.Body.String(), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		i := strings.LastIndexByte(line, ' ')
		v, err := strconv.ParseFloat(line[i+1:], 64)
		if err != nil {
			t.Fatalf("metrics line %q: %v", line, err)
		}
		out[line[:i]] = v
	}
	return out
}

func TestStream_IsEventsKitchen(t *testing.T) {
	if events.Stream != "events.kitchen" {
		t.Errorf("Stream = %q, want events.kitchen", events.Stream)
	}
}

func TestRelay_PublishesCommittedRowsInIDOrderAsTypeAndEnvelope(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	writeEvents(t, pool, q, from, 3)
	rows := outboxRows(t, pool, from)

	n, err := newRelay(pool, rdb, events.RelayOptions{Stream: stream}).Publish(t.Context())

	if err != nil {
		t.Fatalf("Publish error = %v, want none", err)
	}
	all := readStream(t, rdb, stream)
	if n != len(all) {
		t.Errorf("Publish returned %d, but it added %d stream entries", n, len(all))
	}
	var ids []string
	for _, r := range rows {
		ids = append(ids, r.ID.String())
	}
	mine := entriesWith(all, ids)
	assertIDs(t, "stream entries", envelopeIDs(mine), ids)
	for i, e := range mine {
		env, err := json.Marshal(events.EnvelopeOf(rows[i]))
		if err != nil {
			t.Fatal(err)
		}
		if want := map[string]any{"type": mergedType, "envelope": string(env)}; !maps.Equal(e.fields, want) {
			t.Errorf("entry %d = %v, want %v", i, e.fields, want)
		}
	}
	want := slices.Repeat([]rowState{{published: true}}, 3)
	if got := rowStates(t, pool, from); !slices.Equal(got, want) {
		t.Errorf("rows after the pass = %+v, want %+v", got, want)
	}
}

func TestRelay_PublishesEventsOfConcurrentTransactionsInIDOrder(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	const writers, txsPerWriter, eventsPerTx = 20, 25, 2 // 1,000 events
	intos := make([][]string, writers*txsPerWriter)
	for i := range intos {
		intos[i] = randomUUIDs(t, eventsPerTx)
	}
	var g errgroup.Group
	for w := range writers {
		g.Go(func() error {
			for i := range txsPerWriter {
				if err := commitEvents(t.Context(), pool, q, from, intos[w*txsPerWriter+i]); err != nil {
					return err
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	ids := outboxIDs(t, pool, from)
	if len(ids) != 1000 {
		t.Fatalf("wrote %d events, want 1000", len(ids))
	}

	drain(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream}))

	assertIDs(t, "stream entries", streamIDsOf(t, rdb, stream, ids), ids)
}

func TestRelay_RunningRelayPublishesConcurrentWritesOnceInEachWritersOrder(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream := testStream(t, rdb)
	const writers, txsPerWriter = 20, 50 // 1,000 events, one per transaction
	froms := randomUUIDs(t, writers)
	q := testStore(t, pool, nil, froms)
	intos := make([][]string, writers)
	for w := range intos {
		intos[w] = randomUUIDs(t, txsPerWriter)
	}
	runRelay(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream}))

	var g errgroup.Group
	for w := range writers {
		g.Go(func() error {
			for _, into := range intos[w] {
				if err := commitEvents(t.Context(), pool, q, froms[w], []string{into}); err != nil {
					return err
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

	perWriter := make([][]string, writers)
	var all []string
	for w := range writers {
		perWriter[w] = outboxIDs(t, pool, froms[w])
		all = append(all, perWriter[w]...)
	}
	got := envelopeIDs(waitForEntries(t, rdb, stream, all, len(all)))
	position := map[string]int{}
	for i, id := range got {
		if _, dup := position[id]; dup {
			t.Errorf("event %s was published twice without a crash", id)
		}
		position[id] = i
	}
	if len(position) != len(all) {
		t.Errorf("the stream holds %d of the %d events", len(position), len(all))
	}
	for w, ids := range perWriter {
		for i := 1; i < len(ids); i++ {
			if position[ids[i]] < position[ids[i-1]] {
				t.Errorf("writer %d: its event %d was published before its event %d", w, i, i-1)
			}
		}
	}
}

func TestRelay_FailedXAddStopsThePassAtThatRow(t *testing.T) {
	pool := openTestDB(t)
	rdb, faults := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	writeEvents(t, pool, q, from, 5)
	ids := outboxIDs(t, pool, from)
	faults.set(failXAddOf(ids[2:3], errRedisDown), nil)

	_, err := newRelay(pool, rdb, events.RelayOptions{Stream: stream}).Publish(t.Context())

	if !errors.Is(err, errRedisDown) {
		t.Fatalf("Publish error = %v, want %v", err, errRedisDown)
	}
	assertIDs(t, "stream entries", streamIDsOf(t, rdb, stream, ids), ids[:2])
	want := []rowState{
		{published: true},
		{published: true},
		{attempts: 1, lastError: errRedisDown.Error()},
		{},
		{},
	}
	if got := rowStates(t, pool, from); !slices.Equal(got, want) {
		t.Errorf("rows after the failed pass = %+v, want %+v", got, want)
	}
}

func TestRelay_EventsWrittenWhileRedisIsDownArePublishedInOrderOnceItReturns(t *testing.T) {
	pool := openTestDB(t)
	rdb, faults := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	relay := newRelay(pool, rdb, events.RelayOptions{Stream: stream})
	faults.set(redisDown, nil)
	for range 2 {
		writeEvents(t, pool, q, from, 3)
		if _, err := relay.Publish(t.Context()); !errors.Is(err, errRedisDown) {
			t.Fatalf("Publish while Redis is down: error = %v, want %v", err, errRedisDown)
		}
	}
	if got := rowStates(t, pool, from)[0]; got != (rowState{attempts: 2, lastError: errRedisDown.Error()}) {
		t.Fatalf("first row after two failed passes = %+v, want 2 attempts and the error", got)
	}
	faults.set(nil, nil)

	drain(t, relay)

	ids := outboxIDs(t, pool, from)
	assertIDs(t, "stream entries", streamIDsOf(t, rdb, stream, ids), ids)
	for i, s := range rowStates(t, pool, from) {
		if !s.published {
			t.Errorf("row %d is still unpublished after Redis returned", i)
		}
	}
}

func TestRelay_CrashBetweenXAddAndCommitPublishesAgainWithoutGap(t *testing.T) {
	pool := openTestDB(t)
	rdb, faults := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	writeEvents(t, pool, q, from, 5)
	ids := outboxIDs(t, pool, from)
	// The process dies right after the third XADD: nothing after it reaches Redis or Postgres.
	ctx, crash := context.WithCancel(t.Context())
	defer crash()
	faults.set(
		func(ctx context.Context, _ string) error { return ctx.Err() },
		func(id string) {
			if id == ids[2] {
				crash()
			}
		})
	if _, err := newRelay(pool, rdb, events.RelayOptions{Stream: stream}).Publish(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Publish of the crashing process: error = %v, want %v", err, context.Canceled)
	}
	if got := rowStates(t, pool, from); !slices.Equal(got, make([]rowState, 5)) {
		t.Fatalf("rows after the crash = %+v, want all unpublished and untouched", got)
	}
	faults.set(nil, nil)

	drain(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream})) // the restarted process

	mine := entriesWith(readStream(t, rdb, stream), ids)
	assertIDs(t, "stream entries", envelopeIDs(mine), append(slices.Clone(ids[:3]), ids...))
	if len(mine) == 8 {
		for i := range 3 {
			if mine[i].envelope != mine[3+i].envelope {
				t.Errorf("the duplicate of event %d differs from the first copy:\n%s\n%s", i, mine[i].envelope, mine[3+i].envelope)
			}
		}
	}
}

func TestRelay_TrimsStreamEntriesOlderThanSevenDays(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	now := time.Now()
	oldMs, recentMs := now.Add(-10*day).UnixMilli(), now.Add(-6*day).UnixMilli()
	recent := []string{fmt.Sprintf("%d-0", recentMs), fmt.Sprintf("%d-1", recentMs)}
	if _, err := rdb.Pipelined(t.Context(), func(p redis.Pipeliner) error {
		for i := range 300 {
			p.XAdd(t.Context(), &redis.XAddArgs{Stream: stream, ID: fmt.Sprintf("%d-%d", oldMs, i), Values: []any{"type", "old", "envelope", "{}"}})
		}
		for _, id := range recent {
			p.XAdd(t.Context(), &redis.XAddArgs{Stream: stream, ID: id, Values: []any{"type", "recent", "envelope", "{}"}})
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	writeEvents(t, pool, q, from, 1)
	ids := outboxIDs(t, pool, from)

	if _, err := newRelay(pool, rdb, events.RelayOptions{Stream: stream}).Publish(t.Context()); err != nil {
		t.Fatalf("Publish error = %v, want none", err)
	}

	old := 0
	var kept []string
	for _, e := range readStream(t, rdb, stream) {
		if strings.HasPrefix(e.streamID, strconv.FormatInt(oldMs, 10)+"-") {
			old++
		} else {
			kept = append(kept, e.streamID)
		}
	}
	// Approximate trimming removes whole nodes only, so a node holding old and recent entries (fewer than 100) stays.
	if old >= 100 {
		t.Errorf("%d of 300 entries older than 7 days are left, want the whole old nodes trimmed", old)
	}
	if len(kept) != 3 || kept[0] != recent[0] || kept[1] != recent[1] {
		t.Errorf("entries newer than 7 days = %v, want %v and the new one", kept, recent)
	}
	assertIDs(t, "new entries", streamIDsOf(t, rdb, stream, ids), ids)
}

func TestRelay_CleanupDeletesPublishedRowsOlderThanSevenDays(t *testing.T) {
	pool := openTestDB(t)
	from := randomUUID(t)
	testStore(t, pool, nil, []string{from})
	ago := func(d time.Duration) *time.Time { at := cleanupNow.Add(-d); return &at }
	insertOutboxRow(t, pool, from, cleanupNow.Add(-9*day), ago(week+time.Millisecond))
	exactlyAWeek := insertOutboxRow(t, pool, from, cleanupNow.Add(-9*day), ago(week))
	recent := insertOutboxRow(t, pool, from, cleanupNow.Add(-9*day), ago(6*day))
	unpublished := insertOutboxRow(t, pool, from, cleanupNow.Add(-30*day), nil)

	n, err := newRelay(pool, nil, events.RelayOptions{Now: func() time.Time { return cleanupNow }}).Cleanup(t.Context())

	if err != nil || n != 1 {
		t.Errorf("Cleanup = %d, %v; want 1 row deleted", n, err)
	}
	want := []string{exactlyAWeek, recent, unpublished}
	slices.Sort(want)
	assertIDs(t, "rows left", outboxIDs(t, pool, from), want)
}

func TestRelay_CleanupDeletesAtMostTenThousandRowsPerRun(t *testing.T) {
	pool := openTestDB(t)
	from := randomUUID(t)
	testStore(t, pool, nil, []string{from})
	if _, err := pool.Exec(t.Context(), `
		INSERT INTO outbox (id, type, payload, occurred_at, published_at)
		SELECT gen_random_uuid(), 'catalogue.ingredient_merged', jsonb_build_object('from_ingredient_id', $1::text), $2, $2
		FROM generate_series(1, 10001)`, from, cleanupNow.Add(-8*day)); err != nil {
		t.Fatal(err)
	}
	relay := newRelay(pool, nil, events.RelayOptions{Now: func() time.Time { return cleanupNow }})

	first, err1 := relay.Cleanup(t.Context())
	second, err2 := relay.Cleanup(t.Context())

	if err1 != nil || err2 != nil || first != 10000 || second != 1 {
		t.Errorf("two cleanups = %d (%v), %d (%v); want 10000, then 1", first, err1, second, err2)
	}
}

func TestRelay_RunCleansUpAtStartAndThenOnceAnHour(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	testStore(t, pool, nil, []string{from})
	ago := func(d time.Duration) *time.Time { at := cleanupNow.Add(-d); return &at }
	insertOutboxRow(t, pool, from, cleanupNow.Add(-9*day), ago(8*day))
	later := insertOutboxRow(t, pool, from, cleanupNow.Add(-9*day), ago(6*day))
	clk, timer := &clock{now: cleanupNow}, newFakeTimer()

	runRelay(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream, Now: clk.Now, After: timer.After}))

	w := timer.next(t)
	assertIDs(t, "rows left after the start", outboxIDs(t, pool, from), []string{later})
	setPublishedAt(t, pool, from, ago(8*day)) // later is old now; an UPDATE sends no notification
	clk.Set(cleanupNow.Add(59 * time.Minute))
	w.fire()
	w = timer.next(t)
	assertIDs(t, "rows left 59 minutes later", outboxIDs(t, pool, from), []string{later})
	clk.Set(cleanupNow.Add(time.Hour))
	w.fire()
	timer.next(t)
	assertIDs(t, "rows left an hour later", outboxIDs(t, pool, from), nil)
}

func TestRelay_WakesOnNotify(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	timer := newFakeTimer() // never fired: only a notification can start a pass
	runRelay(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream, After: timer.After}))
	timer.next(t) // listening, first pass done

	writeEvents(t, pool, q, from, 1)

	ids := outboxIDs(t, pool, from)
	waitForEntries(t, rdb, stream, ids, 1)
}

func TestRelay_PassesEverySecondWithoutNotify(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	// A row that becomes pending without a notification: written as published before the relay starts, then made
	// pending again while the relay waits (an UPDATE sends no notification).
	writeEvents(t, pool, q, from, 1)
	published := time.Now()
	setPublishedAt(t, pool, from, &published)
	timer := newFakeTimer()
	runRelay(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream, After: timer.After}))
	w := timer.next(t)
	if w.d != time.Second {
		t.Errorf("the relay waits %s for a notification, want 1s", w.d)
	}
	setPublishedAt(t, pool, from, nil)

	w.fire()

	timer.next(t)
	ids := outboxIDs(t, pool, from)
	assertIDs(t, "stream entries", streamIDsOf(t, rdb, stream, ids), ids)
}

func TestRelay_RetriesAfterBackoffOf1To30Seconds(t *testing.T) {
	pool := openTestDB(t)
	rdb, faults := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	faults.set(redisDown, nil)
	writeEvents(t, pool, q, from, 1)
	timer := newFakeTimer()
	runRelay(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream, After: timer.After}))

	want := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second, 30 * time.Second, 30 * time.Second}
	var waits []time.Duration
	var w *fakeWait
	for i := range want {
		w = timer.next(t)
		waits = append(waits, w.d)
		if i < len(want)-1 {
			w.fire()
		}
	}
	faults.set(nil, nil) // Redis is back
	w.fire()
	w = timer.next(t)
	waits = append(waits, w.d)

	if want = append(want, time.Second); !slices.Equal(waits, want) {
		t.Errorf("waits = %v, want %v", waits, want)
	}
	ids := outboxIDs(t, pool, from)
	assertIDs(t, "stream entries once Redis is back", streamIDsOf(t, rdb, stream, ids), ids)
	faults.set(redisDown, nil) // the next outage starts again at 1 s
	writeEvents(t, pool, q, from, 1)
	w.fire()
	if w = timer.next(t); w.d != time.Second {
		t.Errorf("first wait of the next outage = %s, want 1s", w.d)
	}
}

func TestRelay_FinishesTheCurrentBatchOnShutdown(t *testing.T) {
	pool := openTestDB(t)
	rdb, faults := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	writeEvents(t, pool, q, from, 3)
	ids := outboxIDs(t, pool, from)
	reached, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	faults.set(func(_ context.Context, id string) error {
		if id == ids[1] {
			once.Do(func() { close(reached) })
			<-release
		}
		return nil
	}, nil)
	stop, done := runRelay(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream, After: newFakeTimer().After}))
	select {
	case <-reached:
	case <-time.After(10 * time.Second):
		t.Fatal("the relay didn't reach the second event within 10 s")
	}

	stop()
	close(release)

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run returned %v on shutdown, want nil", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Run did not return within 10 s of shutdown")
	}
	assertIDs(t, "stream entries", streamIDsOf(t, rdb, stream, ids), ids)
	if got, want := rowStates(t, pool, from), slices.Repeat([]rowState{{published: true}}, 3); !slices.Equal(got, want) {
		t.Errorf("rows after shutdown = %+v, want all published", got)
	}
}

func TestRelay_PassesAgainAtOnceAfterAFullBatch(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	writeEvents(t, pool, q, from, 150)
	ids := outboxIDs(t, pool, from)
	timer := newFakeTimer()

	runRelay(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream, After: timer.After}))

	timer.next(t) // the first wait comes only once the backlog is gone
	assertIDs(t, "stream entries before the first wait", streamIDsOf(t, rdb, stream, ids), ids)
}

func TestRelay_ExportsMetrics(t *testing.T) {
	pool := openTestDB(t)
	rdb, faults := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	writeEvents(t, pool, q, from, 3)
	ids := outboxIDs(t, pool, from)
	now := time.Now().UTC().Truncate(time.Millisecond)
	if _, err := pool.Exec(t.Context(), "UPDATE outbox SET occurred_at = $2 WHERE payload->>'from_ingredient_id' = $1", from, now.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	faults.set(failXAddOf(ids[2:], errRedisDown), nil)
	relay := newRelay(pool, rdb, events.RelayOptions{Stream: stream, Now: func() time.Time { return now }})

	_, _ = relay.Publish(t.Context())

	var pending int
	var oldest time.Time
	if err := pool.QueryRow(t.Context(), "SELECT count(*), min(occurred_at) FROM outbox WHERE published_at IS NULL").Scan(&pending, &oldest); err != nil {
		t.Fatal(err)
	}
	if pending < 1 || now.Sub(oldest) < time.Hour {
		t.Fatalf("backlog = %d rows, oldest %s; want the test's failed row in it", pending, oldest)
	}
	want := map[string]float64{
		fmt.Sprintf("events_published_total{stream=%q,type=%q}", stream, mergedType): 2,
		fmt.Sprintf("events_publish_errors_total{stream=%q}", stream):                1,
		"events_outbox_pending":        float64(pending),
		"events_outbox_oldest_seconds": now.Sub(oldest).Seconds(),
	}
	got := scrapeMetrics(t)
	for name, v := range want {
		if g, ok := got[name]; !ok || g != v {
			t.Errorf("/metrics %s = %v (present %v), want %v", name, g, ok, v)
		}
	}
	faults.set(nil, nil)
	drain(t, relay)
	got = scrapeMetrics(t)
	if got["events_outbox_pending"] != 0 || got["events_outbox_oldest_seconds"] != 0 {
		t.Errorf("after draining: pending %v, oldest %v s; want 0 and 0", got["events_outbox_pending"], got["events_outbox_oldest_seconds"])
	}
}

func TestRelay_MedianLatencyFromCommitToXAddIsUnderOneSecond(t *testing.T) {
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	stream, from := testStream(t, rdb), randomUUID(t)
	q := testStore(t, pool, nil, []string{from})
	timer := newFakeTimer() // never fired: every pass is started by a notification
	runRelay(t, newRelay(pool, rdb, events.RelayOptions{Stream: stream, After: timer.After}))
	timer.next(t)

	var latencies []time.Duration
	for range 21 {
		writeEvents(t, pool, q, from, 1)
		committed := time.Now()
		ids := outboxIDs(t, pool, from)
		e := waitForEntries(t, rdb, stream, ids[len(ids)-1:], 1)[0]
		ms, err := strconv.ParseInt(strings.SplitN(e.streamID, "-", 2)[0], 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		latencies = append(latencies, time.UnixMilli(ms).Sub(committed))
	}

	slices.Sort(latencies)
	if median := latencies[len(latencies)/2]; median >= time.Second {
		t.Errorf("median latency from commit to XADD = %s, want under 1s", median)
	}
}
