package rpc

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metric names are shared by all services (docs/service-conventions.md).
var (
	requests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_server_requests_total",
		Help: "RPCs handled, by service, method and status code.",
	}, []string{"rpc_service", "rpc_method", "rpc_code"})
	duration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "rpc_server_duration_seconds",
		Help:    "RPC handling time.",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30},
	}, []string{"rpc_service", "rpc_method"})
	clientRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rpc_client_requests_total",
		Help: "Outgoing RPCs, by target service, method and status code.",
	}, []string{"rpc_service", "rpc_method", "rpc_code"})
)

func splitProcedure(p string) (service, method string) {
	p = strings.TrimPrefix(p, "/")
	if i := strings.LastIndex(p, "/"); i >= 0 {
		return p[:i], p[i+1:]
	}
	return p, ""
}

func code(err error) string {
	if err == nil {
		return "ok"
	}
	var ce *connect.Error
	if errors.As(err, &ce) {
		return ce.Code().String()
	}
	return connect.CodeUnknown.String()
}

// NewServerInterceptor records metrics and logs every unary RPC.
func NewServerInterceptor(log *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			svc, method := splitProcedure(req.Spec().Procedure)
			ctx = WithMetadata(ctx, req.Header().Get)
			res, err := next(ctx, req)
			c := code(err)
			requests.WithLabelValues(svc, method, c).Inc()
			duration.WithLabelValues(svc, method).Observe(time.Since(start).Seconds())
			lvl := slog.LevelDebug
			if err != nil {
				lvl = slog.LevelWarn
			}
			log.Log(ctx, lvl, "rpc", "method", req.Spec().Procedure, "code", c,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", req.Header().Get("x-request-id"), "error", err)
			return res, err
		}
	}
}

// NewClientInterceptor forwards request metadata and counts outgoing calls.
func NewClientInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if md, ok := ctx.Value(metadataKey{}).(map[string]string); ok {
				for k, v := range md {
					req.Header().Set(k, v)
				}
			}
			res, err := next(ctx, req)
			svc, method := splitProcedure(req.Spec().Procedure)
			clientRequests.WithLabelValues(svc, method, code(err)).Inc()
			return res, err
		}
	}
}

type metadataKey struct{}

// ForwardedHeaders are copied from an incoming request to outgoing calls.
var ForwardedHeaders = []string{"x-user-id", "x-household-id", "x-plan", "x-request-id", "traceparent"}

// WithMetadata stores the forwarded headers of an incoming request in ctx.
func WithMetadata(ctx context.Context, get func(string) string) context.Context {
	md := map[string]string{}
	for _, k := range ForwardedHeaders {
		if v := get(k); v != "" {
			md[k] = v
		}
	}
	return context.WithValue(ctx, metadataKey{}, md)
}
