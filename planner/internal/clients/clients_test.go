package clients

import (
	"context"
	"errors"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/types/known/emptypb"

	aiv1 "coolinary/planner/gen/coolinary/ai/v1"
	"coolinary/planner/gen/coolinary/ai/v1/aiv1connect"
	identityv1 "coolinary/planner/gen/coolinary/identity/v1"
	"coolinary/planner/gen/coolinary/identity/v1/identityv1connect"
	kitchenv1 "coolinary/planner/gen/coolinary/kitchen/v1"
	"coolinary/planner/gen/coolinary/kitchen/v1/kitchenv1connect"
	notifierv1 "coolinary/planner/gen/coolinary/notifier/v1"
	"coolinary/planner/gen/coolinary/notifier/v1/notifierv1connect"
	offersv1 "coolinary/planner/gen/coolinary/offers/v1"
	"coolinary/planner/gen/coolinary/offers/v1/offersv1connect"
	"coolinary/planner/internal/platform/rpc"
)

// Procedures of RPCs the contracts don't have yet; the policies go by method name.
const (
	procKitchenMatchNames      = "/coolinary.kitchen.v1.KitchenService/MatchNames"
	procKitchenListStock       = "/coolinary.kitchen.v1.KitchenService/ListStock"
	procKitchenAddStockItems   = "/coolinary.kitchen.v1.KitchenService/AddStockItems"
	procKitchenUpdateStockItem = "/coolinary.kitchen.v1.KitchenService/UpdateStockItem"
	procKitchenRemoveStock     = "/coolinary.kitchen.v1.KitchenService/RemoveStockItems"
	procKitchenConsumeStock    = "/coolinary.kitchen.v1.KitchenService/ConsumeStock"
	procOffersGetPrices        = "/coolinary.offers.v1.OffersService/GetPrices"
	procAIComposeDayPlan       = "/coolinary.ai.v1.AiService/ComposeDayPlan"
	procAISubmitBatch          = "/coolinary.ai.v1.AiService/SubmitBatch"
)

// gatewayMetadata is what the gateway sets on a signed-in Plus user's request in German.
var gatewayMetadata = map[string]string{
	"x-request-id":   "req-test-0001",
	"traceparent":    "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	"x-user-id":      "user_test01",
	"x-household-id": "0190a5f0-0000-7000-8000-000000000001",
	"x-plan":         "plus",
	"x-language":     "de",
}

// forwardedKeys are the metadata keys a called service may receive.
var forwardedKeys = []string{
	"x-request-id", "traceparent", "x-user-id", "x-household-id", "x-plan", "x-language", "x-admin-user-id",
}

// incoming is a request context carrying md, as the server interceptor builds it.
func incoming(ctx context.Context, md map[string]string) context.Context {
	return rpc.WithMetadata(ctx, func(k string) string { return md[k] })
}

// attempt is one call a fake service received.
type attempt struct {
	at       time.Time
	metadata map[string]string
}

// fake stands in for a called service on one procedure. Each attempt gets the next of its results (the last one
// repeats; nil = OK). With hang set it answers only when the call's deadline has passed.
type fake struct {
	mu       sync.Mutex
	results  []error
	hang     bool
	attempts []attempt
}

func (f *fake) handle(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	md := map[string]string{}
	for _, k := range forwardedKeys {
		if v := req.Header().Get(k); v != "" {
			md[k] = v
		}
	}
	f.mu.Lock()
	n := len(f.attempts)
	f.attempts = append(f.attempts, attempt{at: time.Now(), metadata: md})
	var err error
	if len(f.results) > 0 {
		err = f.results[min(n, len(f.results)-1)]
	}
	f.mu.Unlock()
	if f.hang {
		<-ctx.Done()
		return nil, connect.NewError(connect.CodeDeadlineExceeded, ctx.Err())
	}
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (f *fake) calls() []attempt {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]attempt(nil), f.attempts...)
}

// startFake serves f on procedure from a local h2c server and returns its host:port.
func startFake(t *testing.T, procedure string, f *fake) string {
	t.Helper()
	mux := http.NewServeMux()
	mux.Handle(procedure, connect.NewUnaryHandler(procedure, f.handle))
	srv := httptest.NewUnstartedServer(mux)
	var protocols http.Protocols
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	srv.Config.Protocols = &protocols
	srv.Start()
	t.Cleanup(srv.Close)
	return strings.TrimPrefix(srv.URL, "http://")
}

// stoppedAddr is an address nothing listens on, like a stopped service.
func stoppedAddr(t *testing.T) string {
	t.Helper()
	var lc net.ListenConfig
	l, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := l.Addr().String()
	if err := l.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	return addr
}

// callWith calls procedure at addr through a client built with the policy's options, as New builds them.
func callWith(ctx context.Context, p policy, addr, procedure string) error {
	client := connect.NewClient[emptypb.Empty, emptypb.Empty](rpc.NewHTTPClient(), rpc.BaseURL(addr)+procedure,
		p.options()...)
	_, err := client.CallUnary(ctx, connect.NewRequest(&emptypb.Empty{}))
	return err
}

// clientRequests reads rpc_client_requests_total for the labels; an empty code sums every code.
func clientRequests(t *testing.T, service, method, code string) float64 {
	t.Helper()
	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	var sum float64
	for _, f := range families {
		if f.GetName() != "rpc_client_requests_total" {
			continue
		}
		for _, m := range f.GetMetric() {
			labels := map[string]string{}
			for _, l := range m.GetLabel() {
				labels[l.GetName()] = l.GetValue()
			}
			if labels["rpc_service"] == service && labels["rpc_method"] == method &&
				(code == "" || labels["rpc_code"] == code) {
				sum += m.GetCounter().GetValue()
			}
		}
	}
	return sum
}

func TestCallDeadline_IsTheSmallerOfCallersDeadlineMinus200msAndThePerCallDefault(t *testing.T) {
	tests := []struct {
		name       string
		p          policy
		method     string
		callerLeft time.Duration // 0: the caller has no deadline
		want       time.Duration
	}{
		{name: "kitchen read 2 s", p: kitchenPolicy, method: "GetPlanningProfile", want: 2 * time.Second},
		{name: "kitchen write 2 s", p: kitchenPolicy, method: "ConsumeStock", want: 2 * time.Second},
		{name: "offers 2 s", p: offersPolicy, method: "GetPrices", want: 2 * time.Second},
		{name: "identity 1 s", p: identityPolicy, method: "GetEntitlements", want: time.Second},
		{name: "notifier 3 s", p: notifierPolicy, method: "SendPush", want: 3 * time.Second},
		{name: "ai text 15 s", p: aiPolicy, method: "ExtractItemsFromText", want: 15 * time.Second},
		{name: "ai match 15 s", p: aiPolicy, method: "PickCatalogueMatches", want: 15 * time.Second},
		{name: "ai compose day plan 45 s", p: aiPolicy, method: "ComposeDayPlan", want: 45 * time.Second},
		{name: "ai modify meal 45 s", p: aiPolicy, method: "ModifyMeal", want: 45 * time.Second},
		{name: "ai rewrite dish 45 s", p: aiPolicy, method: "RewriteDish", want: 45 * time.Second},
		{name: "ai compose family day 45 s", p: aiPolicy, method: "ComposeFamilyDay", want: 45 * time.Second},
		{name: "ai compose from deals 45 s", p: aiPolicy, method: "ComposeFromDeals", want: 45 * time.Second},
		{name: "ai compose week plan 45 s", p: aiPolicy, method: "ComposeWeekPlan", want: 45 * time.Second},
		{name: "ai cheapest top-up 45 s", p: aiPolicy, method: "CheapestTopUp", want: 45 * time.Second},
		{name: "ai batch submit 10 s", p: aiPolicy, method: "SubmitBatch", want: 10 * time.Second},
		{name: "ai batch poll 10 s", p: aiPolicy, method: "GetBatch", want: 10 * time.Second},
		{name: "ai batch cancel 10 s", p: aiPolicy, method: "CancelBatch", want: 10 * time.Second},
		{name: "caller with 1 s left: 800 ms", p: kitchenPolicy, method: "ListStock", callerLeft: time.Second,
			want: 800 * time.Millisecond},
		{name: "caller with 2.2 s left: exactly the default", p: kitchenPolicy, method: "ListStock",
			callerLeft: 2200 * time.Millisecond, want: 2 * time.Second},
		{name: "caller with 10 s left: the default", p: kitchenPolicy, method: "ListStock", callerLeft: 10 * time.Second,
			want: 2 * time.Second},
		{name: "caller with 30 s left: ai compose gets 29.8 s", p: aiPolicy, method: "ComposeDayPlan",
			callerLeft: 30 * time.Second, want: 29800 * time.Millisecond},
		{name: "caller with 150 ms left: already past", p: kitchenPolicy, method: "ListStock",
			callerLeft: 150 * time.Millisecond, want: -50 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			ctx := t.Context()
			if tt.callerLeft > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithDeadline(ctx, now.Add(tt.callerLeft))
				defer cancel()
			}
			var got time.Time
			var set bool
			record := connect.UnaryInterceptorFunc(func(connect.UnaryFunc) connect.UnaryFunc {
				return func(ctx context.Context, _ connect.AnyRequest) (connect.AnyResponse, error) {
					got, set = ctx.Deadline()
					return connect.NewResponse(&emptypb.Empty{}), nil
				}
			})
			client := connect.NewClient[emptypb.Empty, emptypb.Empty](http.DefaultClient,
				"http://target.invalid/coolinary.test.v1.Target/"+tt.method, connect.WithGRPC(),
				connect.WithInterceptors(tt.p.interceptor(func() time.Time { return now }), record))

			if _, err := client.CallUnary(ctx, connect.NewRequest(&emptypb.Empty{})); err != nil {
				t.Fatalf("call: %v", err)
			}

			if !set {
				t.Fatalf("the call has no deadline, want now + %v", tt.want)
			}
			if d := got.Sub(now); d != tt.want {
				t.Errorf("deadline is now + %v, want now + %v", d, tt.want)
			}
		})
	}
}

func TestKitchenAndOffersReads_RetryOnceAfter100msWhenUnavailable(t *testing.T) {
	unavailable := connect.NewError(connect.CodeUnavailable, errors.New("down"))
	slow := connect.NewError(connect.CodeDeadlineExceeded, errors.New("its own call timed out"))
	tests := []struct {
		name      string
		p         policy
		procedure string
		first     error
	}{
		{name: "kitchen ResolveBarcode", p: kitchenPolicy, procedure: kitchenv1connect.KitchenServiceResolveBarcodeProcedure,
			first: unavailable},
		{name: "kitchen MatchNames is a read (Q238)", p: kitchenPolicy, procedure: procKitchenMatchNames, first: unavailable},
		{name: "kitchen answering DEADLINE_EXCEEDED itself", p: kitchenPolicy, procedure: procKitchenListStock, first: slow},
		{name: "offers GetPrices", p: offersPolicy, procedure: procOffersGetPrices, first: unavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fake{results: []error{tt.first, nil}}
			addr := startFake(t, tt.procedure, f)

			err := callWith(t.Context(), tt.p, addr, tt.procedure)

			if err != nil {
				t.Fatalf("call: %v, want the retry's success", err)
			}
			calls := f.calls()
			if len(calls) != 2 {
				t.Fatalf("%d attempts, want 2", len(calls))
			}
			if gap := calls[1].at.Sub(calls[0].at); gap < 100*time.Millisecond {
				t.Errorf("retried after %v, want at least 100 ms", gap)
			}
		})
	}
}

func TestKitchenRead_RetriesOnlyOnce(t *testing.T) {
	f := &fake{results: []error{connect.NewError(connect.CodeUnavailable, errors.New("down"))}}
	addr := startFake(t, kitchenv1connect.KitchenServiceResolveBarcodeProcedure, f)

	err := callWith(t.Context(), kitchenPolicy, addr, kitchenv1connect.KitchenServiceResolveBarcodeProcedure)

	if KindOf(err) != Unavailable {
		t.Errorf("kind %v (%v), want Unavailable", KindOf(err), err)
	}
	if n := len(f.calls()); n != 2 {
		t.Errorf("%d attempts, want 2", n)
	}
}

func TestCallsThatMustNotRetry_AreAttemptedOnce(t *testing.T) {
	unavailable := connect.NewError(connect.CodeUnavailable, errors.New("down"))
	tests := []struct {
		name      string
		p         policy
		procedure string
		err       *connect.Error
	}{
		{name: "kitchen write AddStockItems", p: kitchenPolicy, procedure: procKitchenAddStockItems, err: unavailable},
		{name: "kitchen write UpdateStockItem", p: kitchenPolicy, procedure: procKitchenUpdateStockItem, err: unavailable},
		{name: "kitchen write RemoveStockItems", p: kitchenPolicy, procedure: procKitchenRemoveStock, err: unavailable},
		{name: "kitchen write ConsumeStock", p: kitchenPolicy, procedure: procKitchenConsumeStock, err: unavailable},
		{name: "ai compose", p: aiPolicy, procedure: procAIComposeDayPlan, err: unavailable},
		{name: "ai batch submit", p: aiPolicy, procedure: procAISubmitBatch, err: unavailable},
		{name: "identity", p: identityPolicy, procedure: identityv1connect.IdentityServiceGetHouseholdProcedure,
			err: unavailable},
		{name: "notifier", p: notifierPolicy, procedure: notifierv1connect.NotifierServiceSendEmailProcedure,
			err: unavailable},
		{name: "kitchen read rejected", p: kitchenPolicy, procedure: procKitchenListStock,
			err: connect.NewError(connect.CodeNotFound, errors.New("household not found"))},
		{name: "kitchen read failing with a bug", p: kitchenPolicy, procedure: procKitchenListStock,
			err: connect.NewError(connect.CodeInternal, errors.New("boom"))},
		{name: "offers read rejected", p: offersPolicy, procedure: procOffersGetPrices,
			err: connect.NewError(connect.CodeInvalidArgument, errors.New("too many shops"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fake{results: []error{tt.err, nil}}
			addr := startFake(t, tt.procedure, f)

			err := callWith(t.Context(), tt.p, addr, tt.procedure)

			if got := connect.CodeOf(err); got != tt.err.Code() {
				t.Errorf("code %v (%v), want %v", got, err, tt.err.Code())
			}
			if n := len(f.calls()); n != 1 {
				t.Errorf("%d attempts, want 1", n)
			}
		})
	}
}

func TestKitchenRead_NotRetriedWhenTheFirstAttemptUsedUpTheDeadline(t *testing.T) {
	f := &fake{hang: true}
	addr := startFake(t, procKitchenListStock, f)
	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond) // the call gets 300 ms
	defer cancel()
	const svc, method = kitchenv1connect.KitchenServiceName, "ListStock"
	before, beforeTimeouts := clientRequests(t, svc, method, ""), clientRequests(t, svc, method, "deadline_exceeded")

	err := callWith(ctx, kitchenPolicy, addr, procKitchenListStock)

	if KindOf(err) != Unavailable || connect.CodeOf(err) != connect.CodeDeadlineExceeded {
		t.Errorf("kind %v, code %v (%v), want Unavailable with deadline_exceeded", KindOf(err), connect.CodeOf(err), err)
	}
	if n := clientRequests(t, svc, method, "") - before; n != 1 {
		t.Errorf("%v attempts counted, want 1", n)
	}
	if n := clientRequests(t, svc, method, "deadline_exceeded") - beforeTimeouts; n != 1 {
		t.Errorf("rpc_client_requests_total{rpc_code=deadline_exceeded} went up by %v, want 1", n)
	}
}

func TestKitchenStopped_ReadReturnsUnavailableWithin2300msAfterOneRetry(t *testing.T) {
	c := New(Addrs{Kitchen: stoppedAddr(t)})
	const svc, method = kitchenv1connect.KitchenServiceName, "ResolveBarcode"
	before := clientRequests(t, svc, method, "unavailable")
	start := time.Now()

	_, err := c.Kitchen.ResolveBarcode(t.Context(), connect.NewRequest(&kitchenv1.ResolveBarcodeRequest{Ean: "4006381333931"}))

	if took := time.Since(start); took > 2300*time.Millisecond {
		t.Errorf("took %v, want at most 2.3 s", took)
	}
	if KindOf(err) != Unavailable {
		t.Errorf("kind %v (%v), want Unavailable", KindOf(err), err)
	}
	if n := clientRequests(t, svc, method, "unavailable") - before; n != 2 {
		t.Errorf("%v attempts counted, want 2 (one retry)", n)
	}
}

func TestAICall_ReceivesPlanAndLanguageWithTheOtherRequestMetadata(t *testing.T) {
	f := &fake{}
	c := New(Addrs{AI: startFake(t, aiv1connect.AiServiceExtractItemsFromTextProcedure, f)})
	ctx := incoming(t.Context(), gatewayMetadata)

	_, err := c.AI.ExtractItemsFromText(ctx, connect.NewRequest(&aiv1.ExtractItemsFromTextRequest{}))

	if err != nil {
		t.Fatalf("call: %v", err)
	}
	calls := f.calls()
	if len(calls) != 1 {
		t.Fatalf("%d attempts, want 1", len(calls))
	}
	if got := calls[0].metadata; !maps.Equal(got, gatewayMetadata) {
		t.Errorf("ai received %v, want %v", got, gatewayMetadata)
	}
}

func TestAICall_AttemptedExactlyOnceWhenUnavailable(t *testing.T) {
	f := &fake{results: []error{connect.NewError(connect.CodeUnavailable, errors.New("down")), nil}}
	c := New(Addrs{AI: startFake(t, aiv1connect.AiServiceExtractItemsFromTextProcedure, f)})

	_, err := c.AI.ExtractItemsFromText(t.Context(), connect.NewRequest(&aiv1.ExtractItemsFromTextRequest{}))

	if KindOf(err) != Unavailable {
		t.Errorf("kind %v (%v), want Unavailable", KindOf(err), err)
	}
	if n := len(f.calls()); n != 1 {
		t.Errorf("ai was called %d times, want once", n)
	}
}

func TestAICall_SpendLimitIsRejectedWithItsCodeKept(t *testing.T) {
	f := &fake{results: []error{connect.NewError(connect.CodeResourceExhausted, errors.New("daily spend limit reached"))}}
	c := New(Addrs{AI: startFake(t, aiv1connect.AiServiceExtractItemsFromTextProcedure, f)})

	_, err := c.AI.ExtractItemsFromText(t.Context(), connect.NewRequest(&aiv1.ExtractItemsFromTextRequest{}))

	if KindOf(err) != Rejected || connect.CodeOf(err) != connect.CodeResourceExhausted {
		t.Errorf("kind %v, code %v (%v), want Rejected with resource_exhausted", KindOf(err), connect.CodeOf(err), err)
	}
}

func TestClientRequestsTotal_IsLabelledWithTargetServiceMethodAndCode(t *testing.T) {
	notFound := connect.NewError(connect.CodeNotFound, errors.New("not found"))
	tests := []struct {
		name    string
		service string
		method  string
		call    func(context.Context, string) error // calls the method on a client to this address
	}{
		{name: "kitchen", service: kitchenv1connect.KitchenServiceName, method: "ResolveBarcode",
			call: func(ctx context.Context, addr string) error {
				_, err := New(Addrs{Kitchen: addr}).Kitchen.ResolveBarcode(ctx,
					connect.NewRequest(&kitchenv1.ResolveBarcodeRequest{}))
				return err
			}},
		{name: "offers", service: offersv1connect.OffersServiceName, method: "ListOffers",
			call: func(ctx context.Context, addr string) error {
				_, err := New(Addrs{Offers: addr}).Offers.ListOffers(ctx, connect.NewRequest(&offersv1.ListOffersRequest{}))
				return err
			}},
		{name: "ai", service: aiv1connect.AiServiceName, method: "ExtractItemsFromText",
			call: func(ctx context.Context, addr string) error {
				_, err := New(Addrs{AI: addr}).AI.ExtractItemsFromText(ctx,
					connect.NewRequest(&aiv1.ExtractItemsFromTextRequest{}))
				return err
			}},
		{name: "notifier", service: notifierv1connect.NotifierServiceName, method: "SendEmail",
			call: func(ctx context.Context, addr string) error {
				_, err := New(Addrs{Notifier: addr}).Notifier.SendEmail(ctx,
					connect.NewRequest(&notifierv1.SendEmailRequest{}))
				return err
			}},
		{name: "identity", service: identityv1connect.IdentityServiceName, method: "GetHousehold",
			call: func(ctx context.Context, addr string) error {
				_, err := New(Addrs{Identity: addr}).Identity.GetHousehold(ctx,
					connect.NewRequest(&identityv1.GetHouseholdRequest{}))
				return err
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := startFake(t, "/"+tt.service+"/"+tt.method, &fake{results: []error{notFound}})
			before := clientRequests(t, tt.service, tt.method, "not_found")

			err := tt.call(t.Context(), addr)

			if connect.CodeOf(err) != connect.CodeNotFound {
				t.Fatalf("code %v (%v), want not_found", connect.CodeOf(err), err)
			}
			if n := clientRequests(t, tt.service, tt.method, "not_found") - before; n != 1 {
				t.Errorf("rpc_client_requests_total{rpc_service=%q,rpc_method=%q,rpc_code=\"not_found\"} went up by %v, want 1",
					tt.service, tt.method, n)
			}
		})
	}
}
