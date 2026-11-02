// Package clients builds planner's clients to kitchen, offers, ai, notifier and identity (PLN-01-02): every call
// gets a deadline, carries the request metadata, kitchen and offers reads are retried once, and errors are sorted
// into three kinds (KindOf). Clients never log request or response bodies.
package clients

import (
	"context"
	"strings"
	"time"

	"connectrpc.com/connect"

	"coolinary/planner/gen/coolinary/ai/v1/aiv1connect"
	"coolinary/planner/gen/coolinary/identity/v1/identityv1connect"
	"coolinary/planner/gen/coolinary/kitchen/v1/kitchenv1connect"
	"coolinary/planner/gen/coolinary/notifier/v1/notifierv1connect"
	"coolinary/planner/gen/coolinary/offers/v1/offersv1connect"
	"coolinary/planner/internal/platform/rpc"
)

// Addrs are the host:port addresses of the services planner calls.
type Addrs struct {
	Kitchen, Offers, AI, Notifier, Identity string
}

// Clients holds one generated client per called service. Build it once (main.go) and share it.
type Clients struct {
	Kitchen  kitchenv1connect.KitchenServiceClient
	Offers   offersv1connect.OffersServiceClient
	AI       aiv1connect.AiServiceClient
	Notifier notifierv1connect.NotifierServiceClient
	Identity identityv1connect.IdentityServiceClient
}

// New builds the clients on one h2c HTTP client.
func New(a Addrs) *Clients {
	hc := rpc.NewHTTPClient()
	return &Clients{
		Kitchen:  kitchenv1connect.NewKitchenServiceClient(hc, rpc.BaseURL(a.Kitchen), kitchenPolicy.options()...),
		Offers:   offersv1connect.NewOffersServiceClient(hc, rpc.BaseURL(a.Offers), offersPolicy.options()...),
		AI:       aiv1connect.NewAiServiceClient(hc, rpc.BaseURL(a.AI), aiPolicy.options()...),
		Notifier: notifierv1connect.NewNotifierServiceClient(hc, rpc.BaseURL(a.Notifier), notifierPolicy.options()...),
		Identity: identityv1connect.NewIdentityServiceClient(hc, rpc.BaseURL(a.Identity), identityPolicy.options()...),
	}
}

const (
	// deadlineMargin is kept from the caller's deadline, so planner can still answer after a call times out.
	deadlineMargin = 200 * time.Millisecond
	// retryWait is the pause before a read is retried.
	retryWait = 100 * time.Millisecond
)

// policy is how planner calls one service (PLN-01-02 rules 2 and 5).
type policy struct {
	timeout  time.Duration            // per-call default
	byMethod map[string]time.Duration // methods with another default
	retry    bool                     // retry once after retryWait when Unavailable
	noRetry  map[string]bool          // writes, never retried (they carry an idempotency key)
}

var (
	kitchenPolicy = policy{timeout: 2 * time.Second, retry: true, noRetry: map[string]bool{ // Q238
		"AddStockItems": true, "UpdateStockItem": true, "RemoveStockItems": true, "ConsumeStock": true,
	}}
	offersPolicy   = policy{timeout: 2 * time.Second, retry: true}
	notifierPolicy = policy{timeout: 3 * time.Second}
	identityPolicy = policy{timeout: time.Second}
	// ai: text and match calls 15 s, compose calls 45 s, batch submit and poll 10 s; never retried (each costs money).
	aiPolicy = policy{timeout: 15 * time.Second, byMethod: map[string]time.Duration{
		"ComposeDayPlan": aiCompose, "ModifyMeal": aiCompose, "RewriteDish": aiCompose, "ComposeFamilyDay": aiCompose,
		"ComposeFromDeals": aiCompose, "ComposeWeekPlan": aiCompose, "CheapestTopUp": aiCompose,
		"SubmitBatch": aiBatch, "GetBatch": aiBatch, "CancelBatch": aiBatch,
	}}
)

const (
	aiCompose = 45 * time.Second
	aiBatch   = 10 * time.Second
)

// options are the client options for a service called with this policy. The policy's interceptor runs outside
// rpc's, so each attempt is forwarded the metadata and counted in rpc_client_requests_total.
func (p policy) options() []connect.ClientOption {
	return append([]connect.ClientOption{connect.WithInterceptors(p.interceptor(time.Now))}, rpc.ClientOptions()...)
}

// interceptor gives each call its deadline: the per-call default, or the caller's deadline minus deadlineMargin
// when that comes first. Reads of a retrying policy are tried once more after retryWait when Unavailable, within
// the same deadline (Q239). Unary calls only: a server-streaming RPC needs its own deadline.
func (p policy) interceptor(now func() time.Time) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			method := req.Spec().Procedure[strings.LastIndex(req.Spec().Procedure, "/")+1:]
			deadline := now().Add(p.timeoutOf(method))
			if caller, ok := ctx.Deadline(); ok && caller.Add(-deadlineMargin).Before(deadline) {
				deadline = caller.Add(-deadlineMargin)
			}
			ctx, cancel := context.WithDeadline(ctx, deadline)
			defer cancel()
			res, err := next(ctx, req)
			if !p.retry || p.noRetry[method] || KindOf(err) != Unavailable {
				return res, err
			}
			wait := time.NewTimer(retryWait)
			defer wait.Stop()
			select {
			case <-ctx.Done():
				return res, err
			case <-wait.C:
				return next(ctx, req)
			}
		}
	}
}

func (p policy) timeoutOf(method string) time.Duration {
	if d, ok := p.byMethod[method]; ok {
		return d
	}
	return p.timeout
}
