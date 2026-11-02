package events_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"coolinary/planner/internal/planner/store"
	"coolinary/planner/internal/platform/events"
)

const (
	planChanged = "household.plan_changed"
	memberAdded = "household.member_added"
	group       = "planner"
)

var (
	errHandlerBroke = errors.New("handler broke")
	errProcessDied  = errors.New("process died")
)

// consumerEnv is one consumer test's surroundings: the test database, redis-state with breakable commands, a fresh
// stream standing in for events.identity, and the tag that marks the effects the test's handler writes.
type consumerEnv struct {
	pool   *pgxpool.Pool
	rdb    *redis.Client
	faults *commandFaults
	waits  *fakeWaits
	stream string
	tag    string

	mu       sync.Mutex
	eventIDs []string // every envelope id the test made, for cleanup
}

// newConsumerEnv makes the test stream with group planner at its end, as every start after planner's first finds it.
func newConsumerEnv(t *testing.T) *consumerEnv {
	t.Helper()
	e := newConsumerEnvWithoutGroup(t)
	if err := e.rdb.XGroupCreateMkStream(t.Context(), e.stream, group, "$").Err(); err != nil {
		t.Fatal(err)
	}
	return e
}

// newConsumerEnvWithoutGroup is planner's first start: no group yet, and no stream until the test produces.
func newConsumerEnvWithoutGroup(t *testing.T) *consumerEnv {
	t.Helper()
	pool := openTestDB(t)
	rdb, faults := openTestRedis(t)
	e := &consumerEnv{
		pool:   pool,
		rdb:    rdb,
		faults: faults,
		waits:  newFakeWaits(),
		stream: "test.events.identity." + randomUUID(t),
		tag:    randomUUID(t),
	}
	t.Cleanup(func() {
		ctx := context.Background()
		_ = rdb.Del(ctx, e.stream).Err()
		if _, err := pool.Exec(ctx, "DELETE FROM processed_events WHERE event_id::text = ANY($1) OR stream = $2",
			e.eventIDs, effectStream(e.tag)); err != nil {
			t.Error(err)
		}
		if _, err := pool.Exec(ctx, "DELETE FROM dead_events WHERE stream = $1", e.stream); err != nil {
			t.Error(err)
		}
	})
	return e
}

func (e *consumerEnv) track(id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventIDs = append(e.eventIDs, id)
}

// consumerOpts are a test consumer's settings; zero values are a discarded log, a clock fixed at testNow, no jitter
// and a 10 ms blocking read (so a fired wait is noticed at once).
type consumerOpts struct {
	log          *slog.Logger
	now          func() time.Time
	jitter       func(time.Duration) time.Duration
	defaultBlock bool // block as the service does (5 s) instead of 10 ms
}

// consumer builds a consumer of the test's stream in group planner, with fake waits.
func (e *consumerEnv) consumer(name string, o consumerOpts) *events.Consumer {
	if o.log == nil {
		o.log = slog.New(slog.DiscardHandler)
	}
	if o.now == nil {
		o.now = (&clock{now: testNow}).Now
	}
	if o.jitter == nil {
		o.jitter = func(time.Duration) time.Duration { return 0 }
	}
	block := 10 * time.Millisecond
	if o.defaultBlock {
		block = 0
	}
	return events.NewConsumer(e.pool, e.rdb, o.log, events.ConsumerOptions{
		Streams: []string{e.stream},
		Group:   group,
		Name:    name,
		Now:     o.now,
		After:   e.waits.After,
		Jitter:  o.jitter,
		Block:   block,
	})
}

// envelope builds an envelope as identity's publisher writes it, with a fresh id.
func (e *consumerEnv) envelope(t *testing.T, typ string, version int32, household string) events.Envelope {
	t.Helper()
	id := randomUUID(t)
	e.track(id)
	payload := fmt.Sprintf(`{"household_id":%q,"plan":"plus","previous_plan":"free","source":"manual"}`, household)
	return events.Envelope{
		ID:          id,
		Type:        typ,
		Version:     version,
		OccurredAt:  "2026-11-02T10:00:00.000Z",
		Producer:    "identity",
		HouseholdID: household,
		Payload:     json.RawMessage(payload),
	}
}

func rawOf(t *testing.T, env events.Envelope) string {
	t.Helper()
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

// produce is the fake producer: it adds env to the stream with the fields type and envelope (Q015) and returns the
// entry's stream id.
func (e *consumerEnv) produce(t *testing.T, env events.Envelope) string {
	t.Helper()
	return e.produceRaw(t, env.Type, rawOf(t, env))
}

func (e *consumerEnv) produceRaw(t *testing.T, typ, envelope string) string {
	t.Helper()
	return e.produceFields(t, "type", typ, "envelope", envelope)
}

func (e *consumerEnv) produceFields(t *testing.T, fields ...any) string {
	t.Helper()
	id, err := e.rdb.XAdd(t.Context(), &redis.XAddArgs{Stream: e.stream, Values: fields}).Result()
	if err != nil {
		t.Fatal(err)
	}
	return id
}

// The test handler's effect is a processed_events row of its own: stream "effect.<tag>", type the envelope id, and
// an event id that sorts in the order the handler ran. Planner has no other table a handler can write yet.
func effectStream(tag string) string { return "effect." + tag }

func effectID(tag string, seq int) pgtype.UUID {
	var id pgtype.UUID
	_ = id.Scan(fmt.Sprintf("%s%012x", tag[:24], seq))
	return id
}

// effects returns, in the order they were applied, the envelope ids of the events the test handler applied.
func (e *consumerEnv) effects(t *testing.T) []string {
	t.Helper()
	rows, err := e.pool.Query(t.Context(), "SELECT type FROM processed_events WHERE stream = $1 ORDER BY event_id",
		effectStream(e.tag))
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		out = append(out, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

type processedRow struct{ stream, typ string }

func (e *consumerEnv) processed(t *testing.T, eventID string) []processedRow {
	t.Helper()
	rows, err := e.pool.Query(t.Context(), "SELECT stream, type FROM processed_events WHERE event_id = $1", eventID)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []processedRow
	for rows.Next() {
		var r processedRow
		if err := rows.Scan(&r.stream, &r.typ); err != nil {
			t.Fatal(err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

type deadRow struct {
	entryID, eventID, typ, envelope, lastError string
	deliveries                                 int32
}

// deadRows lists the test stream's dead_events rows in stream order; eventID is "" when NULL.
func (e *consumerEnv) deadRows(t *testing.T) []deadRow {
	t.Helper()
	rows, err := e.pool.Query(t.Context(), `SELECT entry_id, COALESCE(event_id::text, ''), type, envelope, last_error,
		deliveries FROM dead_events WHERE stream = $1 ORDER BY entry_id`, e.stream)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []deadRow
	for rows.Next() {
		var r deadRow
		if err := rows.Scan(&r.entryID, &r.eventID, &r.typ, &r.envelope, &r.lastError, &r.deliveries); err != nil {
			t.Fatal(err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

type pendingEntry struct {
	entryID    string
	consumer   string
	deliveries int64
}

// pending lists the group's unacknowledged entries in stream order.
func (e *consumerEnv) pending(t *testing.T) []pendingEntry {
	t.Helper()
	ps, err := e.rdb.XPendingExt(t.Context(), &redis.XPendingExtArgs{
		Stream: e.stream, Group: group, Start: "-", End: "+", Count: 1000,
	}).Result()
	if err != nil {
		t.Fatal(err)
	}
	out := make([]pendingEntry, 0, len(ps))
	for _, p := range ps {
		out = append(out, pendingEntry{entryID: p.ID, consumer: p.Consumer, deliveries: p.RetryCount})
	}
	return out
}

// agePending makes every pending entry look idle for 61 s and owned by owner, as if its consumer had died a minute
// ago. JUSTID keeps the delivery count; deliveries > 0 sets it instead.
func (e *consumerEnv) agePending(t *testing.T, owner string, deliveries int64) {
	t.Helper()
	for _, p := range e.pending(t) {
		args := []any{"XCLAIM", e.stream, group, owner, 0, p.entryID, "IDLE", 61000}
		if deliveries > 0 {
			args = append(args, "RETRYCOUNT", deliveries)
		}
		args = append(args, "JUSTID")
		if err := e.rdb.Do(t.Context(), args...).Err(); err != nil {
			t.Fatal(err)
		}
	}
}

// testHandler applies an event by writing one effect row carrying the envelope id; fail, when set, runs after that
// write and its error fails the handler (so the write must roll back).
type testHandler struct {
	tag  string
	fail func(ctx context.Context, env events.Envelope, call int) error

	mu          sync.Mutex
	calls       []string
	inFlight    int
	maxInFlight int
}

func (e *consumerEnv) handler(fail func(context.Context, events.Envelope, int) error) *testHandler {
	return &testHandler{tag: e.tag, fail: fail}
}

func (h *testHandler) handle(ctx context.Context, q *store.Queries, env events.Envelope) error {
	h.mu.Lock()
	h.calls = append(h.calls, env.ID)
	seq := len(h.calls)
	call := 0
	for _, id := range h.calls {
		if id == env.ID {
			call++
		}
	}
	h.inFlight++
	h.maxInFlight = max(h.maxInFlight, h.inFlight)
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		h.inFlight--
		h.mu.Unlock()
	}()
	if _, err := q.InsertProcessedEvent(ctx, store.InsertProcessedEventParams{
		EventID: effectID(h.tag, seq), Stream: effectStream(h.tag), Type: env.ID,
	}); err != nil {
		return err
	}
	if h.fail != nil {
		return h.fail(ctx, env, call)
	}
	return nil
}

func (h *testHandler) callList() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return slices.Clone(h.calls)
}

func (h *testHandler) mostAtOnce() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.maxInFlight
}

// failOn fails the given calls (1 = first delivery) of the events with these ids; no calls means every call.
func failOn(ids []string, calls ...int) func(context.Context, events.Envelope, int) error {
	return func(_ context.Context, env events.Envelope, call int) error {
		if slices.Contains(ids, env.ID) && (len(calls) == 0 || slices.Contains(calls, call)) {
			return errHandlerBroke
		}
		return nil
	}
}

// runConsumer starts c.Run. stop shuts it down; done delivers what Run returned. The test's end stops it too.
func runConsumer(t *testing.T, c *events.Consumer) (stop context.CancelFunc, done <-chan error) {
	t.Helper()
	ctx, cancel := context.WithCancel(t.Context())
	result, finished := make(chan error, 1), make(chan struct{})
	go func() {
		result <- c.Run(ctx)
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

func awaitReturn(t *testing.T, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run = %v, want nil", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Run did not return within 10 s of shutdown")
	}
}

// consumed reads events_consumed_total for one stream and type: result to count, zero counts left out.
func consumed(t *testing.T, stream, typ string) map[string]float64 {
	t.Helper()
	m := scrapeMetrics(t)
	out := map[string]float64{}
	for _, r := range []string{"ok", "duplicate", "skipped", "error", "dead"} {
		key := metricKey("events_consumed_total", map[string]string{"stream": stream, "type": typ, "result": r})
		if v := m[key]; v != 0 {
			out[r] = v
		}
	}
	return out
}

func waitConsumed(t *testing.T, stream, typ string, want map[string]float64) {
	t.Helper()
	waitUntil(t, fmt.Sprintf("events_consumed_total %v", want), func() bool {
		return maps.Equal(consumed(t, stream, typ), want)
	})
}

func idsOf(envs []events.Envelope) []string {
	out := make([]string, 0, len(envs))
	for _, env := range envs {
		out = append(out, env.ID)
	}
	return out
}

// logRecords parses the JSON log lines whose msg is msg.
func logRecords(t *testing.T, logs, msg string) []map[string]any {
	t.Helper()
	var out []map[string]any
	for line := range strings.SplitSeq(strings.TrimSpace(logs), "\n") {
		if line == "" {
			continue
		}
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("log line %q: %v", line, err)
		}
		if rec["msg"] == msg {
			out = append(out, rec)
		}
	}
	return out
}

func jsonLogger(buf *syncBuffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestConsumer_CreatesGroupAtStreamEndSoOnlyNewEntriesAreProcessed(t *testing.T) {
	e := newConsumerEnvWithoutGroup(t)
	household := randomUUID(t)
	old := e.envelope(t, planChanged, 1, household)
	e.produce(t, old) // published before planner's group existed
	h := e.handler(nil)
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	e.waits.next(t, time.Minute) // started: the group exists and the first claim ran

	fresh := e.envelope(t, planChanged, 1, household)
	e.produce(t, fresh)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1})
	assertIDs(t, "applied events", e.effects(t), []string{fresh.ID})
	assertIDs(t, "handler calls", h.callList(), []string{fresh.ID})
	groups, err := e.rdb.XInfoGroups(t.Context(), e.stream).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 || groups[0].Name != group || groups[0].Pending != 0 {
		t.Errorf("groups = %+v, want one group %q with nothing pending", groups, group)
	}
}

func TestConsumer_CreatesMissingStream(t *testing.T) {
	e := newConsumerEnvWithoutGroup(t)
	h := e.handler(nil)
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)

	runConsumer(t, c)
	e.waits.next(t, time.Minute) // started: the group exists and the first claim ran

	groups, err := e.rdb.XInfoGroups(t.Context(), e.stream).Result()
	if err != nil {
		t.Fatalf("XINFO GROUPS on the stream the consumer should have made: %v", err)
	}
	if len(groups) != 1 || groups[0].Name != group {
		t.Fatalf("groups = %+v, want one group %q", groups, group)
	}
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	e.produce(t, env)
	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1})
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
}

func TestConsumer_ReadsUpTo100NewEntriesBlockingFiveSeconds(t *testing.T) {
	e := newConsumerEnv(t)
	reads := make(chan string, 100)
	e.faults.set(func(name string, args []any) error {
		if name == "xreadgroup" && slices.Contains(args, any(">")) {
			reads <- fmt.Sprint(args)
			return errRedisDown // so the test doesn't wait out the block
		}
		return nil
	})
	c := e.consumer("planner-test", consumerOpts{defaultBlock: true})
	c.Handle(planChanged, 1, e.handler(nil).handle)

	runConsumer(t, c)

	select {
	case got := <-reads:
		want := fmt.Sprint([]any{"xreadgroup", "group", group, "planner-test", "count", int64(100), "block", int64(5000),
			"streams", e.stream, ">"})
		if got != want {
			t.Errorf("read = %s\nwant   %s", got, want)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("the consumer read no new entries within 10 s")
	}
}

func TestConsumer_AppliesEntriesOneAtATimeInStreamOrder(t *testing.T) {
	e := newConsumerEnv(t)
	households := []string{randomUUID(t), randomUUID(t)}
	var envs []events.Envelope
	for i := range 150 { // more than one read of 100
		env := e.envelope(t, planChanged, 1, households[i%2])
		e.produce(t, env)
		envs = append(envs, env)
	}
	h := e.handler(nil)
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)

	runConsumer(t, c)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 150})
	assertIDs(t, "applied events", e.effects(t), idsOf(envs))
	if n := h.mostAtOnce(); n != 1 {
		t.Errorf("up to %d handlers ran at once, want 1", n)
	}
	for _, env := range envs {
		if got, want := e.processed(t, env.ID), []processedRow{{e.stream, planChanged}}; !slices.Equal(got, want) {
			t.Fatalf("processed_events of %s = %v, want %v", env.ID, got, want)
		}
	}
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none", p)
	}
}

func TestConsumer_AcknowledgesUnknownTypeWithoutParsingAndCountsItSkipped(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)

	e.produceRaw(t, memberAdded, "{not json")
	handled := e.envelope(t, planChanged, 1, randomUUID(t))
	e.produce(t, handled)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1})
	if got, want := consumed(t, e.stream, memberAdded), map[string]float64{"skipped": 1}; !maps.Equal(got, want) {
		t.Errorf("events_consumed_total of %s = %v, want %v", memberAdded, got, want)
	}
	assertIDs(t, "handler calls", h.callList(), []string{handled.ID})
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none", p)
	}
	if got := e.deadRows(t); len(got) != 0 {
		t.Errorf("dead_events = %v, want none", got)
	}
}

func TestConsumer_HandlesVersionsUpToTheHandlersAndParksNewerAtOnce(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 2, h.handle)
	runConsumer(t, c)
	household := randomUUID(t)
	v1, v2, v3 := e.envelope(t, planChanged, 1, household), e.envelope(t, planChanged, 2, household), e.envelope(t, planChanged, 3, household)

	e.produce(t, v1)
	e.produce(t, v2)
	v3Entry := e.produce(t, v3)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 2, "dead": 1})
	assertIDs(t, "applied events", e.effects(t), []string{v1.ID, v2.ID})
	assertIDs(t, "handler calls", h.callList(), []string{v1.ID, v2.ID})
	want := []deadRow{{v3Entry, v3.ID, planChanged, rawOf(t, v3), "event version newer than its handler: version 3, handler takes up to 2", 1}}
	if got := e.deadRows(t); !slices.Equal(got, want) {
		t.Errorf("dead_events = %+v, want %+v", got, want)
	}
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none: the newer version is acknowledged without retries", p)
	}
	if got := e.processed(t, v3.ID); len(got) != 0 {
		t.Errorf("processed_events of the newer version = %v, want none", got)
	}
}

func TestConsumer_ParksMalformedEnvelopesAtOnce(t *testing.T) {
	tests := []struct {
		name      string
		fields    []any
		envelope  string // stored as received
		wantError string
	}{
		{"not JSON", []any{"type", planChanged, "envelope", `{"id":`}, `{"id":`, "malformed envelope: not JSON"},
		{
			"id not a UUID",
			[]any{"type", planChanged, "envelope", `{"id":"42","type":"household.plan_changed","version":1}`},
			`{"id":"42","type":"household.plan_changed","version":1}`,
			"malformed envelope: id is not a UUID",
		},
		{"no envelope field", []any{"type", planChanged}, "", "malformed envelope: not JSON"},
		{
			"field of the wrong type",
			[]any{"type", planChanged, "envelope", `{"id":"0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b","version":"1"}`},
			`{"id":"0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b","version":"1"}`,
			"malformed envelope: version has the wrong type",
		},
		{
			// Postgres text holds neither NUL nor invalid UTF-8: they are dropped and replaced so the entry is parked.
			"bytes text can't hold",
			[]any{"type", planChanged, "envelope", "{\x00\xff"},
			"{\uFFFD",
			"malformed envelope: not JSON",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newConsumerEnv(t)
			h := e.handler(nil)
			c := e.consumer("planner-test", consumerOpts{})
			c.Handle(planChanged, 1, h.handle)
			runConsumer(t, c)
			next := e.envelope(t, planChanged, 1, randomUUID(t))

			entry := e.produceFields(t, tt.fields...)
			e.produce(t, next)

			waitConsumed(t, e.stream, planChanged, map[string]float64{"dead": 1, "ok": 1})
			want := []deadRow{{entry, "", planChanged, tt.envelope, tt.wantError, 1}}
			if got := e.deadRows(t); !slices.Equal(got, want) {
				t.Errorf("dead_events = %+v, want %+v", got, want)
			}
			assertIDs(t, "handler calls", h.callList(), []string{next.ID})
			if p := e.pending(t); len(p) != 0 {
				t.Errorf("pending = %v, want none", p)
			}
		})
	}
}

func TestConsumer_SameEnvelopeTwiceRunsTheHandlerOnce(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	env := e.envelope(t, planChanged, 1, randomUUID(t))

	e.produce(t, env)
	e.produce(t, env) // the producer's relay published it twice (crash between XADD and marking it sent)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1, "duplicate": 1})
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
	assertIDs(t, "handler calls", h.callList(), []string{env.ID})
	if got, want := e.processed(t, env.ID), []processedRow{{e.stream, planChanged}}; !slices.Equal(got, want) {
		t.Errorf("processed_events = %v, want %v", got, want)
	}
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none: both entries are acknowledged", p)
	}
}

func TestConsumer_FailedHandlerLeavesEntryPendingAndRollsBack(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"handler error", errHandlerBroke},
		{"handler timeout", context.DeadlineExceeded},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newConsumerEnv(t)
			household := randomUUID(t)
			failing, next := e.envelope(t, planChanged, 1, household), e.envelope(t, planChanged, 1, household)
			h := e.handler(func(_ context.Context, env events.Envelope, _ int) error {
				if env.ID == failing.ID {
					return tt.err
				}
				return nil
			})
			c := e.consumer("planner-test", consumerOpts{})
			c.Handle(planChanged, 1, h.handle)
			runConsumer(t, c)

			failingEntry := e.produce(t, failing)
			e.produce(t, next)

			waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1, "ok": 1})
			assertIDs(t, "applied events", e.effects(t), []string{next.ID})
			if got := e.processed(t, failing.ID); len(got) != 0 {
				t.Errorf("processed_events of the failed event = %v, want none", got)
			}
			if got, want := e.pending(t), []pendingEntry{{failingEntry, "planner-test", 1}}; !slices.Equal(got, want) {
				t.Errorf("pending = %v, want %v", got, want)
			}
		})
	}
}

func TestConsumer_HandlerGetsTenSeconds(t *testing.T) {
	e := newConsumerEnv(t)
	type deadlineSeen struct {
		at time.Time
		ok bool
	}
	seen := make(chan deadlineSeen, 1)
	h := e.handler(func(ctx context.Context, _ events.Envelope, _ int) error {
		at, ok := ctx.Deadline()
		seen <- deadlineSeen{at, ok}
		return nil
	})
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	before := time.Now()

	e.produce(t, e.envelope(t, planChanged, 1, randomUUID(t)))

	var got deadlineSeen
	select {
	case got = <-seen:
	case <-time.After(10 * time.Second):
		t.Fatal("the handler didn't run within 10 s")
	}
	after := time.Now()
	if !got.ok {
		t.Fatal("the handler's context has no deadline")
	}
	if got.at.Before(before.Add(10*time.Second)) || got.at.After(after.Add(10*time.Second)) {
		t.Errorf("deadline = %s after the event was produced, want 10 s after the handler started", got.at.Sub(before))
	}
}

func TestConsumer_CrashBeforeAckRedeliversAsNoOp(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	first := e.consumer("planner-test", consumerOpts{})
	first.Handle(planChanged, 1, h.handle)
	stop, done := runConsumer(t, first)
	acks := make(chan struct{})
	var once sync.Once
	e.faults.set(func(name string, args []any) error {
		if name == "xack" && slices.Contains(args, any(e.stream)) {
			once.Do(func() { close(acks) })
			return errProcessDied
		}
		return nil
	})
	e.produce(t, env)
	select {
	case <-acks:
	case <-time.After(10 * time.Second):
		t.Fatal("the consumer sent no XACK within 10 s")
	}
	stop()
	awaitReturn(t, done)
	e.faults.set(nil)

	second := e.consumer("planner-test", consumerOpts{})
	second.Handle(planChanged, 1, h.handle)
	runConsumer(t, second)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1, "duplicate": 1})
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
	assertIDs(t, "handler calls", h.callList(), []string{env.ID})
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none", p)
	}
}

func TestConsumer_KilledMidHandlerAppliesTheEventOnceAfterARestart(t *testing.T) {
	e := newConsumerEnv(t)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	var mu sync.Mutex
	var kill context.CancelFunc
	h := e.handler(func(_ context.Context, _ events.Envelope, call int) error {
		if call > 1 {
			return nil
		}
		mu.Lock()
		k := kill
		mu.Unlock()
		k()
		return errProcessDied // nothing commits and nothing is acknowledged, as when the process dies here
	})
	first := e.consumer("planner-test", consumerOpts{})
	first.Handle(planChanged, 1, h.handle)
	stop, done := runConsumer(t, first)
	mu.Lock()
	kill = stop
	mu.Unlock()
	e.produce(t, env)
	awaitReturn(t, done)

	second := e.consumer("planner-test", consumerOpts{})
	second.Handle(planChanged, 1, h.handle)
	runConsumer(t, second)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1, "ok": 1})
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
	assertIDs(t, "handler calls", h.callList(), []string{env.ID, env.ID})
	if got, want := e.processed(t, env.ID), []processedRow{{e.stream, planChanged}}; !slices.Equal(got, want) {
		t.Errorf("processed_events = %v, want %v", got, want)
	}
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none", p)
	}
}

func TestConsumer_ClaimsEntriesPendingOverAMinute(t *testing.T) {
	e := newConsumerEnv(t)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	h := e.handler(failOn([]string{env.ID}, 1))
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	claim := e.waits.next(t, time.Minute)
	e.produce(t, env)
	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})
	e.agePending(t, "planner-gone", 0) // its consumer died a minute ago

	claim.fire()
	e.waits.next(t, time.Minute) // the claim round is done and the next one is a minute away

	if got, want := consumed(t, e.stream, planChanged), map[string]float64{"error": 1, "ok": 1}; !maps.Equal(got, want) {
		t.Errorf("events_consumed_total = %v, want %v", got, want)
	}
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none", p)
	}
}

func TestConsumer_LeavesEntriesPendingUnderAMinute(t *testing.T) {
	e := newConsumerEnv(t)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	h := e.handler(failOn([]string{env.ID}, 1))
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	claim := e.waits.next(t, time.Minute)
	entry := e.produce(t, env)
	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})

	claim.fire()
	e.waits.next(t, time.Minute)

	assertIDs(t, "handler calls", h.callList(), []string{env.ID})
	if got, want := e.pending(t), []pendingEntry{{entry, "planner-test", 1}}; !slices.Equal(got, want) {
		t.Errorf("pending = %v, want %v", got, want)
	}
}

func TestConsumer_HandlerFailingThreeTimesThenSucceedingAppliesOnce(t *testing.T) {
	e := newConsumerEnv(t)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	h := e.handler(failOn([]string{env.ID}, 1, 2, 3))
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	claim := e.waits.next(t, time.Minute)
	e.produce(t, env)
	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})

	for range 3 {
		e.agePending(t, "planner-test", 0)
		claim.fire()
		claim = e.waits.next(t, time.Minute)
	}

	if got, want := consumed(t, e.stream, planChanged), map[string]float64{"error": 3, "ok": 1}; !maps.Equal(got, want) {
		t.Errorf("events_consumed_total = %v, want %v", got, want)
	}
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
	if got, want := e.processed(t, env.ID), []processedRow{{e.stream, planChanged}}; !slices.Equal(got, want) {
		t.Errorf("processed_events = %v, want %v", got, want)
	}
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none", p)
	}
}

func TestConsumer_ParksEntryAfterFiveFailedDeliveries(t *testing.T) {
	tests := []struct {
		name        string
		delivered   int64 // deliveries before the claim; the claim is one more
		wantPending bool
		want        map[string]float64
	}{
		{"4th failed delivery stays pending", 3, true, map[string]float64{"error": 2}},
		{"5th failed delivery is parked", 4, false, map[string]float64{"error": 1, "dead": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newConsumerEnv(t)
			env := e.envelope(t, planChanged, 1, randomUUID(t))
			h := e.handler(failOn([]string{env.ID}))
			c := e.consumer("planner-test", consumerOpts{})
			c.Handle(planChanged, 1, h.handle)
			runConsumer(t, c)
			claim := e.waits.next(t, time.Minute)
			entry := e.produce(t, env)
			waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})
			e.agePending(t, "planner-test", tt.delivered)

			claim.fire()
			e.waits.next(t, time.Minute)

			if got := consumed(t, e.stream, planChanged); !maps.Equal(got, tt.want) {
				t.Errorf("events_consumed_total = %v, want %v", got, tt.want)
			}
			var wantPending []pendingEntry
			var wantDead []deadRow
			if tt.wantPending {
				wantPending = []pendingEntry{{entry, "planner-test", tt.delivered + 1}}
			} else {
				wantDead = []deadRow{{entry, env.ID, planChanged, rawOf(t, env), errHandlerBroke.Error(), 5}}
			}
			if got := e.pending(t); !slices.Equal(got, wantPending) {
				t.Errorf("pending = %v, want %v", got, wantPending)
			}
			if got := e.deadRows(t); !slices.Equal(got, wantDead) {
				t.Errorf("dead_events = %+v, want %+v", got, wantDead)
			}
			if got := e.effects(t); len(got) != 0 {
				t.Errorf("applied events = %v, want none", got)
			}
		})
	}
}

func TestConsumer_CleanupWaitsForTheNext0330UTC(t *testing.T) {
	tests := []struct {
		name string
		now  time.Time
		want time.Duration
	}{
		{"after 03:30", time.Date(2020, 1, 15, 4, 0, 0, 0, time.UTC), 23*time.Hour + 30*time.Minute},
		{"before 03:30", time.Date(2020, 1, 15, 1, 0, 0, 0, time.UTC), 2*time.Hour + 30*time.Minute},
		{"at 03:30", time.Date(2020, 1, 15, 3, 30, 0, 0, time.UTC), day},
		{"clock in another zone", time.Date(2020, 1, 15, 4, 0, 0, 0, time.FixedZone("CET", 3600)), 30 * time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newConsumerEnv(t)
			c := e.consumer("planner-test", consumerOpts{now: (&clock{now: tt.now}).Now})

			runConsumer(t, c)

			e.waits.next(t, tt.want)
		})
	}
}

func TestConsumer_DeletesProcessedEventsOlderThan14DaysAt0330UTC(t *testing.T) {
	e := newConsumerEnv(t)
	clk := &clock{now: testNow}
	cleanupAt := time.Date(2020, 1, 16, 3, 30, 0, 0, time.UTC)
	insert := func(age time.Duration) string {
		id := randomUUID(t)
		e.track(id)
		if _, err := e.pool.Exec(t.Context(),
			"INSERT INTO processed_events (event_id, stream, type, processed_at) VALUES ($1, $2, $3, $4)",
			id, e.stream, planChanged, cleanupAt.Add(-age)); err != nil {
			t.Fatal(err)
		}
		return id
	}
	older, exactly, younger := insert(14*day+time.Second), insert(14*day), insert(13*day)
	left := func() []string {
		var out []string
		for _, id := range []string{older, exactly, younger} {
			if len(e.processed(t, id)) == 1 {
				out = append(out, id)
			}
		}
		return out
	}
	c := e.consumer("planner-test", consumerOpts{now: clk.Now})
	runConsumer(t, c)
	cleanup := e.waits.next(t, 23*time.Hour+30*time.Minute)
	assertIDs(t, "rows left before 03:30", left(), []string{older, exactly, younger})

	clk.Set(cleanupAt)
	cleanup.fire()
	e.waits.next(t, day) // the cleanup ran; the next one is a day away

	assertIDs(t, "rows left after the 03:30 cleanup", left(), []string{exactly, younger})
}

func TestConsumer_ShutdownFinishesTheCurrentHandlerThenReturns(t *testing.T) {
	e := newConsumerEnv(t)
	household := randomUUID(t)
	envs := []events.Envelope{e.envelope(t, planChanged, 1, household), e.envelope(t, planChanged, 1, household), e.envelope(t, planChanged, 1, household)}
	started, release := make(chan struct{}), make(chan struct{})
	h := e.handler(func(ctx context.Context, env events.Envelope, _ int) error {
		if env.ID == envs[0].ID {
			close(started)
			<-release
			return ctx.Err() // a cancelled handler would roll back
		}
		return nil
	})
	var entries []string
	for _, env := range envs { // before the start, so the first read takes all three
		entries = append(entries, e.produce(t, env))
	}
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	stop, done := runConsumer(t, c)
	select {
	case <-started:
	case <-time.After(10 * time.Second):
		t.Fatal("the first handler didn't start within 10 s")
	}

	stop()
	close(release)
	awaitReturn(t, done)

	assertIDs(t, "applied events", e.effects(t), []string{envs[0].ID})
	assertIDs(t, "handler calls", h.callList(), []string{envs[0].ID})
	want := []pendingEntry{{entries[1], "planner-test", 1}, {entries[2], "planner-test", 1}}
	if got := e.pending(t); !slices.Equal(got, want) {
		t.Errorf("pending = %v, want %v (read, not yet handled)", got, want)
	}
}

func TestConsumer_KeepsPublishOrderAcrossARestart(t *testing.T) {
	tests := []struct {
		name string
		// interrupt stops the first consumer while it handles the second event.
		interrupt func(t *testing.T, e *consumerEnv, second events.Envelope, stop context.CancelFunc) func(context.Context, events.Envelope, int) error
	}{
		{"shutdown during a read", func(_ *testing.T, _ *consumerEnv, second events.Envelope, stop context.CancelFunc) func(context.Context, events.Envelope, int) error {
			return func(_ context.Context, env events.Envelope, _ int) error {
				if env.ID == second.ID {
					stop()
				}
				return nil
			}
		}},
		{"crash before XACK", func(t *testing.T, e *consumerEnv, second events.Envelope, stop context.CancelFunc) func(context.Context, events.Envelope, int) error {
			e.faults.set(func(name string, args []any) error {
				if name != "xack" {
					return nil
				}
				for _, entry := range e.entriesOf(t, second.ID) {
					if slices.Contains(args, any(entry)) {
						stop()
						return errProcessDied
					}
				}
				return nil
			})
			return nil
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newConsumerEnv(t)
			household := randomUUID(t)
			var envs []events.Envelope
			for range 7 {
				envs = append(envs, e.envelope(t, planChanged, 1, household))
			}
			var mu sync.Mutex
			var fail func(context.Context, events.Envelope, int) error
			setFail := func(f func(context.Context, events.Envelope, int) error) {
				mu.Lock()
				defer mu.Unlock()
				fail = f
			}
			h := e.handler(func(ctx context.Context, env events.Envelope, call int) error {
				mu.Lock()
				f := fail
				mu.Unlock()
				if f != nil {
					return f(ctx, env, call)
				}
				return nil
			})
			first := e.consumer("planner-test", consumerOpts{})
			first.Handle(planChanged, 1, h.handle)
			stop, done := runConsumer(t, first)
			setFail(tt.interrupt(t, e, envs[1], stop))
			for _, env := range envs[:5] {
				e.produce(t, env)
			}
			awaitReturn(t, done)
			e.faults.set(nil)
			setFail(nil)
			for _, env := range envs[5:] {
				e.produce(t, env)
			}

			second := e.consumer("planner-test", consumerOpts{})
			second.Handle(planChanged, 1, h.handle)
			runConsumer(t, second)

			waitUntil(t, "all 7 events applied", func() bool { return len(e.effects(t)) >= 7 })
			assertIDs(t, "applied events", e.effects(t), idsOf(envs))
		})
	}
}

// entriesOf returns the stream ids of the entries carrying envelope id.
func (e *consumerEnv) entriesOf(t *testing.T, id string) []string {
	t.Helper()
	msgs, err := e.rdb.XRange(context.Background(), e.stream, "-", "+").Result()
	if err != nil {
		t.Error(err)
		return nil
	}
	var out []string
	for _, m := range msgs {
		raw, _ := m.Values["envelope"].(string)
		var env events.Envelope
		if json.Unmarshal([]byte(raw), &env) == nil && env.ID == id {
			out = append(out, m.ID)
		}
	}
	return out
}

func TestConsumer_LogsHandlerErrorsWithEventIDTypeAndErrorOnly(t *testing.T) {
	e := newConsumerEnv(t)
	var buf syncBuffer
	marker := "payload-marker-" + randomUUID(t)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	env.Payload = json.RawMessage(fmt.Sprintf(`{"note":%q}`, marker))
	h := e.handler(failOn([]string{env.ID}))
	c := e.consumer("planner-test", consumerOpts{log: jsonLogger(&buf)})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)

	e.produce(t, env)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})
	logs := buf.String()
	if strings.Contains(logs, marker) {
		t.Errorf("logs carry payload content: %s", logs)
	}
	recs := logRecords(t, logs, "event handler failed; retrying later")
	if len(recs) != 1 || recs[0]["event_id"] != env.ID || recs[0]["type"] != planChanged || recs[0]["error"] != errHandlerBroke.Error() {
		t.Errorf("handler failure log records = %v, want one with event_id %s, type %q and error %q", recs, env.ID, planChanged, errHandlerBroke)
	}
}

func TestConsumer_LogsParkedEntriesAsErrorsWithoutPayload(t *testing.T) {
	e := newConsumerEnv(t)
	var buf syncBuffer
	marker := "payload-marker-" + randomUUID(t)
	newer := e.envelope(t, planChanged, 2, randomUUID(t))
	newer.Payload = json.RawMessage(fmt.Sprintf(`{"note":%q}`, marker))
	c := e.consumer("planner-test", consumerOpts{log: jsonLogger(&buf)})
	c.Handle(planChanged, 1, e.handler(nil).handle)
	runConsumer(t, c)

	e.produce(t, newer)
	malformedEntry := e.produceRaw(t, planChanged, fmt.Sprintf(`{"note":%q`, marker))

	waitConsumed(t, e.stream, planChanged, map[string]float64{"dead": 2})
	logs := buf.String()
	if strings.Contains(logs, marker) {
		t.Errorf("logs carry payload content: %s", logs)
	}
	recs := logRecords(t, logs, "event parked in dead_events")
	if len(recs) != 2 {
		t.Fatalf("parked log records = %v, want 2", recs)
	}
	for i, want := range []map[string]any{
		{"event_id": newer.ID, "type": planChanged},
		{"event_id": "", "entry_id": malformedEntry, "type": planChanged},
	} {
		if recs[i]["level"] != "ERROR" {
			t.Errorf("record %d level = %v, want ERROR", i, recs[i]["level"])
		}
		for k, v := range want {
			if recs[i][k] != v {
				t.Errorf("record %d %s = %v, want %v", i, k, recs[i][k], v)
			}
		}
	}
}

func TestConsumer_RetriesRedisErrorsEvery2sPlusJitter(t *testing.T) {
	tests := []struct {
		name   string
		jitter time.Duration
		want   time.Duration
	}{
		{"no jitter", 0, 2 * time.Second},
		{"most jitter", 499 * time.Millisecond, 2*time.Second + 499*time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newConsumerEnv(t)
			var mu sync.Mutex
			reads := 0
			var jitterMax []time.Duration
			e.faults.set(func(name string, _ []any) error {
				mu.Lock()
				defer mu.Unlock()
				if name == "xreadgroup" && reads < 2 {
					reads++
					return errRedisDown
				}
				return nil
			})
			h := e.handler(nil)
			c := e.consumer("planner-test", consumerOpts{jitter: func(max time.Duration) time.Duration {
				mu.Lock()
				defer mu.Unlock()
				jitterMax = append(jitterMax, max)
				return tt.jitter
			}})
			c.Handle(planChanged, 1, h.handle)
			runConsumer(t, c)

			e.waits.next(t, tt.want).fire()
			e.waits.next(t, tt.want).fire()
			env := e.envelope(t, planChanged, 1, randomUUID(t))
			e.produce(t, env)

			waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1})
			assertIDs(t, "applied events", e.effects(t), []string{env.ID})
			mu.Lock()
			defer mu.Unlock()
			if want := []time.Duration{500 * time.Millisecond, 500 * time.Millisecond}; !slices.Equal(jitterMax, want) {
				t.Errorf("jitter asked for up to %v, want %v", jitterMax, want)
			}
		})
	}
}

func TestConsumer_LogsRedisErrorsOncePerMinute(t *testing.T) {
	e := newConsumerEnv(t)
	var buf syncBuffer
	clk := &clock{now: testNow}
	e.faults.set(func(name string, _ []any) error {
		if name == "xreadgroup" {
			return errRedisDown
		}
		return nil
	})
	h := e.handler(nil)
	c := e.consumer("planner-test", consumerOpts{log: jsonLogger(&buf), now: clk.Now})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)

	retry := e.waits.next(t, 2*time.Second) // 1st failure
	clk.Set(testNow.Add(59 * time.Second))
	retry.fire()
	retry = e.waits.next(t, 2*time.Second) // 2nd, within the minute
	clk.Set(testNow.Add(time.Minute))
	retry.fire()
	retry = e.waits.next(t, 2*time.Second) // 3rd, a minute after the first line

	recs := logRecords(t, buf.String(), "event stream read failed")
	var failures []float64
	for _, r := range recs {
		f, _ := r["failures"].(float64)
		failures = append(failures, f)
		if want := "read stream: " + errRedisDown.Error(); r["stream"] != e.stream || r["error"] != want {
			t.Errorf("log record = %v, want stream %s and error %q", r, e.stream, want)
		}
	}
	if want := []float64{1, 3}; !slices.Equal(failures, want) {
		t.Errorf("logged failures = %v, want %v (one line a minute)", failures, want)
	}
	e.faults.set(nil)
	retry.fire()
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	e.produce(t, env)
	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1})
}

func TestConsumer_RecreatesItsGroupWhenTheStreamIsDeleted(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	e.waits.next(t, time.Minute)

	if err := e.rdb.Del(t.Context(), e.stream).Err(); err != nil {
		t.Fatal(err)
	}
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	e.produce(t, env) // a new stream without the group: everything in it is new to planner
	e.waits.next(t, 2*time.Second).fire()

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1})
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
}

func TestConsumer_RecordsHandlerDuration(t *testing.T) {
	e := newConsumerEnv(t)
	key := metricKey("events_handler_duration_seconds_count", map[string]string{"type": planChanged})
	before := scrapeMetrics(t)[key]
	failing := e.envelope(t, planChanged, 1, randomUUID(t))
	h := e.handler(failOn([]string{failing.ID}))
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)

	e.produceRaw(t, memberAdded, "{}")
	e.produce(t, failing)
	e.produce(t, e.envelope(t, planChanged, 1, randomUUID(t)))

	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1, "ok": 1})
	if got := scrapeMetrics(t)[key] - before; got != 2 {
		t.Errorf("handler durations recorded = %v, want 2 (one per handler run)", got)
	}
}

func TestConsumer_ExportsGroupLagInEntriesEveryMinute(t *testing.T) {
	e := newConsumerEnv(t)
	for range 3 { // published while planner was down
		e.produce(t, e.envelope(t, planChanged, 1, randomUUID(t)))
	}
	c := e.consumer("planner-test", consumerOpts{})
	c.Handle(planChanged, 1, e.handler(nil).handle)
	key := metricKey("events_consumer_lag", map[string]string{"stream": e.stream, "group": group})
	lag := func() (float64, bool) {
		v, ok := scrapeMetrics(t)[key]
		return v, ok
	}

	runConsumer(t, c)
	claim := e.waits.next(t, time.Minute) // started: the first claim ran

	if v, ok := lag(); !ok || v != 3 {
		t.Errorf("%s = %v (exported: %v), want 3 at start", key, v, ok)
	}
	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 3})
	claim.fire()
	e.waits.next(t, time.Minute)
	if v, ok := lag(); !ok || v != 0 {
		t.Errorf("%s = %v (exported: %v), want 0 a minute later", key, v, ok)
	}
}
