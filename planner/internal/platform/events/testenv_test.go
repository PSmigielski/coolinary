package events_test

import (
	"bytes"
	"context"
	"crypto/rand"
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

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"coolinary/planner/internal/platform/db"
	"coolinary/planner/migrations"
)

const day = 24 * time.Hour

// errRedisDown is what a stopped redis-state looks like to the consumer.
var errRedisDown = errors.New("dial tcp 10.77.0.12:6379: connect: connection refused")

// testNow is the consumer's clock in the tests: years before any real row, so only a test's own rows are old enough
// for the processed_events cleanup, and half an hour after the daily cleanup time.
var testNow = time.Date(2020, 1, 15, 4, 0, 0, 0, time.UTC)

// openTestDB returns a pool on the migrated test database, or skips when TEST_DATABASE_URL isn't set
// (make test sets it while make dev is up).
func openTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("needs TEST_DATABASE_URL (make dev, then make test s=planner)")
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

// openTestRedis returns a client on redis-state's test database whose commands the test can break, or skips when
// TEST_REDIS_STATE_URL isn't set (make test sets it while make dev is up).
func openTestRedis(t *testing.T) (*redis.Client, *commandFaults) {
	t.Helper()
	url := os.Getenv("TEST_REDIS_STATE_URL")
	if url == "" {
		t.Skip("needs TEST_REDIS_STATE_URL (make dev, then make test s=planner)")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	rdb := redis.NewClient(opts)
	t.Cleanup(func() { _ = rdb.Close() })
	faults := &commandFaults{}
	rdb.AddHook(faults)
	return rdb, faults
}

// commandFaults breaks the Redis commands fail returns an error for, before they are sent.
type commandFaults struct {
	mu   sync.Mutex
	fail func(name string, args []any) error
}

func (f *commandFaults) set(fail func(name string, args []any) error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fail = fail
}

func (f *commandFaults) DialHook(next redis.DialHook) redis.DialHook { return next }

func (f *commandFaults) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

func (f *commandFaults) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		f.mu.Lock()
		fail := f.fail
		f.mu.Unlock()
		if fail != nil {
			if err := fail(cmd.Name(), cmd.Args()); err != nil {
				cmd.SetErr(err)
				return err
			}
		}
		return next(ctx, cmd)
	}
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

// fakeWaits stands in for time.After in the consumer. Each wait is queued by its duration (2 s and a little: a retry
// after a Redis error, 1 min: the next claim, hours: the next cleanup) and fires only when the test fires it. The
// consumer asks for a wait only after the work before it is done, so receiving one is a sync point.
type fakeWaits struct {
	mu     sync.Mutex
	queues map[time.Duration]chan *fakeWait
}

type fakeWait struct {
	d time.Duration
	c chan time.Time
}

func (w *fakeWait) fire() { w.c <- time.Time{} }

func newFakeWaits() *fakeWaits { return &fakeWaits{queues: map[time.Duration]chan *fakeWait{}} }

func (f *fakeWaits) queue(d time.Duration) chan *fakeWait {
	f.mu.Lock()
	defer f.mu.Unlock()
	q, ok := f.queues[d]
	if !ok {
		q = make(chan *fakeWait, 1000)
		f.queues[d] = q
	}
	return q
}

func (f *fakeWaits) After(d time.Duration) <-chan time.Time {
	w := &fakeWait{d: d, c: make(chan time.Time, 1)}
	f.queue(d) <- w
	return w.c
}

// next returns the next wait of duration d the consumer asked for; it fails after 10 s.
func (f *fakeWaits) next(t *testing.T, d time.Duration) *fakeWait {
	t.Helper()
	select {
	case w := <-f.queue(d):
		return w
	case <-time.After(10 * time.Second):
		t.Fatalf("the consumer started no %s wait within 10 s", d)
		return nil
	}
}

// clock is a settable consumer clock.
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

// waitUntil polls cond until it holds; it fails after 10 s.
func waitUntil(t *testing.T, what string, cond func() bool) {
	t.Helper()
	timeout := time.After(10 * time.Second)
	tick := time.NewTicker(5 * time.Millisecond)
	defer tick.Stop()
	for !cond() {
		select {
		case <-timeout:
			t.Fatalf("%s: not within 10 s", what)
		case <-tick.C:
		}
	}
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

// metricKey is a series as /metrics names it: name{label="value",...}, labels sorted by name.
func metricKey(name string, labels map[string]string) string {
	keys := slices.Sorted(maps.Keys(labels))
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%q", k, labels[k]))
	}
	return name + "{" + strings.Join(parts, ",") + "}"
}

// syncBuffer is a bytes.Buffer the consumer's goroutine writes and the test reads.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
