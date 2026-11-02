package events_test

import (
	"bytes"
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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/events"
)

const (
	planChanged = "household.plan_changed"
	memberAdded = "household.member_added"
	group       = "kitchen"
)

var (
	errHandlerBroke = errors.New("handler broke")
	errProcessDied  = errors.New("process died before XACK")
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

func newConsumerEnv(t *testing.T) *consumerEnv {
	t.Helper()
	pool := openTestDB(t)
	rdb, _ := openTestRedis(t)
	e := &consumerEnv{
		pool:   pool,
		rdb:    rdb,
		faults: &commandFaults{},
		waits:  newFakeWaits(),
		stream: "test.events.identity." + randomUUID(t),
		tag:    randomUUID(t),
	}
	rdb.AddHook(e.faults)
	t.Cleanup(func() {
		ctx := context.Background()
		_ = rdb.Del(ctx, e.stream).Err()
		if _, err := pool.Exec(ctx, "DELETE FROM outbox WHERE payload->>'from_ingredient_id' = $1", e.tag); err != nil {
			t.Error(err)
		}
		if _, err := pool.Exec(ctx, "DELETE FROM processed_events WHERE event_id::text = ANY($1)", e.eventIDs); err != nil {
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

// consumer builds a consumer of the test's stream in group kitchen, with fake waits and a short blocking read so a
// fired wait is noticed at once. A nil log discards logs.
func (e *consumerEnv) consumer(name string, log *slog.Logger, now func() time.Time) *events.Consumer {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	return events.NewConsumer(e.pool, e.rdb, log, events.ConsumerOptions{
		Streams: []string{e.stream},
		Group:   group,
		Name:    name,
		Now:     now,
		After:   e.waits.After,
		Block:   10 * time.Millisecond,
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

// produce is the fake producer: it adds env to the stream with the fields type and envelope (Q015) and returns the
// entry's stream id.
func (e *consumerEnv) produce(t *testing.T, env events.Envelope) string {
	t.Helper()
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	return e.produceRaw(t, env.Type, string(raw))
}

func (e *consumerEnv) produceRaw(t *testing.T, typ, envelope string) string {
	t.Helper()
	id, err := e.rdb.XAdd(t.Context(), &redis.XAddArgs{
		Stream: e.stream,
		Values: []any{"type", typ, "envelope", envelope},
	}).Result()
	if err != nil {
		t.Fatal(err)
	}
	return id
}

// effects returns, in the order they were applied, the envelope ids of the events the test handler applied.
func (e *consumerEnv) effects(t *testing.T) []string {
	t.Helper()
	rows := outboxRows(t, e.pool, e.tag)
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, payloadOf(t, r).IntoIngredientID)
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

// fakeWaits stands in for time.After in the consumer. Each wait is queued by its duration (1 s and up: backoff,
// 1 min: the next claim, 24 h: the next cleanup) and fires only when the test fires it. The consumer asks for a wait
// only after the work before it is done, so receiving one is a sync point.
type fakeWaits struct {
	mu     sync.Mutex
	queues map[time.Duration]chan *fakeWait
}

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

// testHandler applies an event by writing one outbox row tagged with tag, carrying the envelope id; fail, when set,
// runs after that write and its error fails the handler (so the write must roll back).
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
	if err := events.Add(ctx, q, merged(h.tag, env.ID)); err != nil {
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

// failOn fails the given calls (1 = first delivery) of the events with these ids.
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

// metricKey is a series as /metrics names it: name{label="value",...}, labels sorted by name.
func metricKey(name string, labels map[string]string) string {
	keys := slices.Sorted(maps.Keys(labels))
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%q", k, labels[k]))
	}
	return name + "{" + strings.Join(parts, ",") + "}"
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

func TestConsumer_CreatesGroupAtStreamStartSoRetainedEntriesAreProcessed(t *testing.T) {
	e := newConsumerEnv(t)
	household := randomUUID(t)
	envs := []events.Envelope{e.envelope(t, planChanged, 1, household), e.envelope(t, planChanged, 1, household)}
	for _, env := range envs {
		e.produce(t, env) // the stream exists, the group doesn't yet
	}
	h := e.handler(nil)
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 1, h.handle)

	runConsumer(t, c)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 2})
	assertIDs(t, "applied events", e.effects(t), idsOf(envs))
	groups, err := e.rdb.XInfoGroups(t.Context(), e.stream).Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 || groups[0].Name != group || groups[0].Pending != 0 {
		t.Errorf("groups = %+v, want one group %q with nothing pending", groups, group)
	}
}

func TestConsumer_CreatesMissingStream(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("kitchen-test", nil, nil)
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

func TestConsumer_AppliesEntriesOneAtATimeInStreamOrder(t *testing.T) {
	e := newConsumerEnv(t)
	households := []string{randomUUID(t), randomUUID(t)}
	var envs []events.Envelope
	for i := range 120 { // more than one read of 50
		env := e.envelope(t, planChanged, 1, households[i%2])
		e.produce(t, env)
		envs = append(envs, env)
	}
	h := e.handler(nil)
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 1, h.handle)

	runConsumer(t, c)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 120})
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

func TestConsumer_AcknowledgesUnhandledTypeWithoutParsing(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("kitchen-test", nil, nil)
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
}

func TestConsumer_HandlesVersionsUpToTheHandlersAndRetriesNewer(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 2, h.handle)
	runConsumer(t, c)
	household := randomUUID(t)
	v1, v2, v3 := e.envelope(t, planChanged, 1, household), e.envelope(t, planChanged, 2, household), e.envelope(t, planChanged, 3, household)

	e.produce(t, v1)
	e.produce(t, v2)
	v3Entry := e.produce(t, v3)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 2, "error": 1})
	assertIDs(t, "applied events", e.effects(t), []string{v1.ID, v2.ID})
	assertIDs(t, "handler calls", h.callList(), []string{v1.ID, v2.ID})
	if got, want := e.pending(t), []pendingEntry{{v3Entry, "kitchen-test", 1}}; !slices.Equal(got, want) {
		t.Errorf("pending = %v, want %v (the newer version waits for a retry)", got, want)
	}
	if got := e.processed(t, v3.ID); len(got) != 0 {
		t.Errorf("processed_events of the newer version = %v, want none", got)
	}
}

func TestConsumer_SameEntryTwiceChangesTheDatabaseOnce(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	env := e.envelope(t, planChanged, 1, randomUUID(t))

	e.produce(t, env)
	e.produce(t, env) // the relay published it twice (crash between XADD and commit)

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
			c := e.consumer("kitchen-test", nil, nil)
			c.Handle(planChanged, 1, h.handle)
			runConsumer(t, c)

			failingEntry := e.produce(t, failing)
			e.produce(t, next)

			waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1, "ok": 1})
			assertIDs(t, "applied events", e.effects(t), []string{next.ID})
			if got := e.processed(t, failing.ID); len(got) != 0 {
				t.Errorf("processed_events of the failed event = %v, want none", got)
			}
			if got, want := e.pending(t), []pendingEntry{{failingEntry, "kitchen-test", 1}}; !slices.Equal(got, want) {
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
	c := e.consumer("kitchen-test", nil, nil)
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
	deadline := got.at
	if !got.ok {
		t.Fatal("the handler's context has no deadline")
	}
	if deadline.Before(before.Add(10*time.Second)) || deadline.After(after.Add(10*time.Second)) {
		t.Errorf("deadline = %s after the event was produced, want 10 s after the handler started", deadline.Sub(before))
	}
}

func TestConsumer_CrashBeforeAckRedeliversAsNoOp(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	first := e.consumer("kitchen-test", nil, nil)
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

	second := e.consumer("kitchen-test", nil, nil)
	second.Handle(planChanged, 1, h.handle)
	runConsumer(t, second)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1, "duplicate": 1})
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
	assertIDs(t, "handler calls", h.callList(), []string{env.ID})
	if p := e.pending(t); len(p) != 0 {
		t.Errorf("pending = %v, want none", p)
	}
}

func TestConsumer_ClaimsEntriesPendingOverAMinute(t *testing.T) {
	e := newConsumerEnv(t)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	h := e.handler(failOn([]string{env.ID}, 1))
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	claim := e.waits.next(t, time.Minute)
	e.produce(t, env)
	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})
	e.agePending(t, "kitchen-gone", 0) // its consumer died a minute ago

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
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	claim := e.waits.next(t, time.Minute)
	entry := e.produce(t, env)
	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})

	claim.fire()
	e.waits.next(t, time.Minute)

	assertIDs(t, "handler calls", h.callList(), []string{env.ID})
	if got, want := e.pending(t), []pendingEntry{{entry, "kitchen-test", 1}}; !slices.Equal(got, want) {
		t.Errorf("pending = %v, want %v", got, want)
	}
}

func TestConsumer_HandlerFailingThreeTimesThenSucceedingAppliesOnce(t *testing.T) {
	e := newConsumerEnv(t)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	h := e.handler(failOn([]string{env.ID}, 1, 2, 3))
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	claim := e.waits.next(t, time.Minute)
	e.produce(t, env)
	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})

	for range 3 {
		e.agePending(t, "kitchen-test", 0)
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

func TestConsumer_DeadLetterAfterTenDeliveries(t *testing.T) {
	tests := []struct {
		name        string
		delivered   int64 // deliveries before the claim; the claim is one more
		wantPending bool
		want        map[string]float64
	}{
		{"9th failed delivery stays pending", 8, true, map[string]float64{"error": 2}},
		{"10th failed delivery is a dead letter", 9, false, map[string]float64{"error": 1, "dead": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newConsumerEnv(t)
			env := e.envelope(t, planChanged, 1, randomUUID(t))
			h := e.handler(failOn([]string{env.ID}))
			c := e.consumer("kitchen-test", nil, nil)
			c.Handle(planChanged, 1, h.handle)
			runConsumer(t, c)
			claim := e.waits.next(t, time.Minute)
			entry := e.produce(t, env)
			waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})
			e.agePending(t, "kitchen-test", tt.delivered)

			claim.fire()
			e.waits.next(t, time.Minute)

			if got := consumed(t, e.stream, planChanged); !maps.Equal(got, tt.want) {
				t.Errorf("events_consumed_total = %v, want %v", got, tt.want)
			}
			var want []pendingEntry
			if tt.wantPending {
				want = []pendingEntry{{entry, "kitchen-test", tt.delivered + 1}}
			}
			if got := e.pending(t); !slices.Equal(got, want) {
				t.Errorf("pending = %v, want %v", got, want)
			}
			if got := e.effects(t); len(got) != 0 {
				t.Errorf("applied events = %v, want none", got)
			}
		})
	}
}

func TestConsumer_DeletesProcessedEventsOlderThan14DaysDaily(t *testing.T) {
	e := newConsumerEnv(t)
	clk := &clock{now: cleanupNow}
	insert := func(age time.Duration) string {
		id := randomUUID(t)
		e.track(id)
		if _, err := e.pool.Exec(t.Context(),
			"INSERT INTO processed_events (event_id, stream, type, processed_at) VALUES ($1, $2, $3, $4)",
			id, e.stream, planChanged, cleanupNow.Add(-age)); err != nil {
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
	c := e.consumer("kitchen-test", nil, clk.Now)

	runConsumer(t, c)
	nextCleanup := e.waits.next(t, day)

	assertIDs(t, "rows left after the start-up cleanup", left(), []string{exactly, younger})
	clk.Set(cleanupNow.Add(day))
	nextCleanup.fire()
	e.waits.next(t, day)
	assertIDs(t, "rows left a day later", left(), []string{younger})
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
	c := e.consumer("kitchen-test", nil, nil)
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
	want := []pendingEntry{{entries[1], "kitchen-test", 1}, {entries[2], "kitchen-test", 1}}
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
			first := e.consumer("kitchen-test", nil, nil)
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

			second := e.consumer("kitchen-test", nil, nil)
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
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	marker := "payload-marker-" + randomUUID(t)
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	env.Payload = json.RawMessage(fmt.Sprintf(`{"note":%q}`, marker))
	h := e.handler(failOn([]string{env.ID}))
	c := e.consumer("kitchen-test", logger, nil)
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)

	e.produce(t, env)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"error": 1})
	logs := buf.String()
	if strings.Contains(logs, marker) {
		t.Errorf("logs carry payload content: %s", logs)
	}
	var found bool
	for line := range strings.SplitSeq(strings.TrimSpace(logs), "\n") {
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("log line %q: %v", line, err)
		}
		if rec["event_id"] != env.ID {
			continue
		}
		found = true
		if rec["type"] != planChanged || rec["error"] != errHandlerBroke.Error() {
			t.Errorf("log record = %v, want type %q and error %q", rec, planChanged, errHandlerBroke)
		}
	}
	if !found {
		t.Errorf("no log record with event_id %s in %s", env.ID, logs)
	}
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

func TestConsumer_BacksOffAfterRedisErrorsAndResumes(t *testing.T) {
	e := newConsumerEnv(t)
	var mu sync.Mutex
	reads := 0
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
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)

	e.waits.next(t, time.Second).fire()
	e.waits.next(t, 2*time.Second).fire()
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	e.produce(t, env)

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1})
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
}

func TestConsumer_RecreatesItsGroupWhenTheStreamIsDeleted(t *testing.T) {
	e := newConsumerEnv(t)
	h := e.handler(nil)
	c := e.consumer("kitchen-test", nil, nil)
	c.Handle(planChanged, 1, h.handle)
	runConsumer(t, c)
	e.waits.next(t, time.Minute)

	if err := e.rdb.Del(t.Context(), e.stream).Err(); err != nil {
		t.Fatal(err)
	}
	env := e.envelope(t, planChanged, 1, randomUUID(t))
	e.produce(t, env) // a new stream without the group
	e.waits.next(t, time.Second).fire()

	waitConsumed(t, e.stream, planChanged, map[string]float64{"ok": 1})
	assertIDs(t, "applied events", e.effects(t), []string{env.ID})
}

func TestConsumer_RecordsHandlerDuration(t *testing.T) {
	e := newConsumerEnv(t)
	key := metricKey("events_handler_duration_seconds_count", map[string]string{"type": planChanged})
	before := scrapeMetrics(t)[key]
	failing := e.envelope(t, planChanged, 1, randomUUID(t))
	h := e.handler(failOn([]string{failing.ID}))
	c := e.consumer("kitchen-test", nil, nil)
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
