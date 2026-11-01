package rpc

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	procEcho    = "/coolinary.test.v1.EchoService/Headers"
	procForward = "/coolinary.test.v1.ForwardService/Forward"
)

// allKeys is every metadata key a service forwards.
var allKeys = []string{
	"x-request-id", "traceparent", "x-user-id", "x-household-id", "x-plan", "x-language", "x-admin-user-id",
}

// userRequest is the metadata the gateway sets on a signed-in user's request: six keys, no admin id.
var userRequest = map[string]string{
	"x-request-id":   "req-test-0001",
	"traceparent":    "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	"x-user-id":      "user_test01",
	"x-household-id": householdA,
	"x-plan":         "plus",
	"x-language":     "pl-PL",
}

// syncBuffer is a bytes.Buffer the server goroutine writes and the test reads.
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

// startH2C serves h on a local h2c server (gRPC needs HTTP/2) and returns its base URL.
func startH2C(t *testing.T, path string, h http.Handler) string {
	t.Helper()
	mux := http.NewServeMux()
	mux.Handle(path, h)
	srv := httptest.NewUnstartedServer(mux)
	var protocols http.Protocols
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	srv.Config.Protocols = &protocols
	srv.Start()
	t.Cleanup(srv.Close)
	return srv.URL
}

// startEcho starts a stand-in for a called service (ai) that answers with the metadata it received.
func startEcho(t *testing.T) string {
	t.Helper()
	path, h := procEcho, connect.NewUnaryHandler(procEcho,
		func(_ context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[structpb.Struct], error) {
			got := map[string]any{}
			for _, k := range allKeys {
				if v := req.Header().Get(k); v != "" {
					got[k] = v
				}
			}
			s, err := structpb.NewStruct(got)
			if err != nil {
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			return connect.NewResponse(s), nil
		})
	return startH2C(t, path, h)
}

// startForwarder starts a stand-in for kitchen: a handler behind the server interceptor that calls the echo
// service through a client built with ClientOptions, as kitchen calls ai, and returns the echo's answer.
func startForwarder(t *testing.T, echoURL string, log *slog.Logger, fail error) string {
	t.Helper()
	echo := connect.NewClient[emptypb.Empty, structpb.Struct](NewHTTPClient(), echoURL+procEcho, ClientOptions()...)
	h := connect.NewUnaryHandler(procForward,
		func(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[structpb.Struct], error) {
			if fail != nil {
				return nil, fail
			}
			ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			res, err := echo.CallUnary(ctx, connect.NewRequest(&emptypb.Empty{}))
			if err != nil {
				return nil, err
			}
			return connect.NewResponse(res.Msg), nil
		}, Options(log)...)
	return startH2C(t, procForward, h)
}

// call sends one gRPC request with these headers to the forwarder.
func call(t *testing.T, url string, headers map[string]string) (*connect.Response[structpb.Struct], error) {
	t.Helper()
	client := connect.NewClient[emptypb.Empty, structpb.Struct](NewHTTPClient(), url+procForward, connect.WithGRPC())
	req := connect.NewRequest(&emptypb.Empty{})
	for k, v := range headers {
		req.Header().Set(k, v)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	return client.CallUnary(ctx, req)
}

func TestOutgoingCallsForwardIncomingMetadata(t *testing.T) {
	withAdmin := maps.Clone(userRequest)
	withAdmin["x-admin-user-id"] = "admin_test01"
	tests := []struct {
		name    string
		headers map[string]string
	}{
		{name: "a user request carries all six keys on", headers: userRequest},
		{name: "an admin request also carries x-admin-user-id on", headers: withAdmin},
		{name: "keys missing on the incoming request are not sent", headers: map[string]string{"x-request-id": "req-test-0002"}},
		{name: "a request without metadata sends none", headers: map[string]string{}},
	}
	url := startForwarder(t, startEcho(t), slog.New(slog.DiscardHandler), nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := call(t, url, tt.headers)
			if err != nil {
				t.Fatalf("call: %v", err)
			}
			got := map[string]string{}
			for k, v := range res.Msg.AsMap() {
				got[k], _ = v.(string)
			}
			if !maps.Equal(got, tt.headers) {
				t.Errorf("the called service received %v, want %v", got, tt.headers)
			}
		})
	}
}

func TestServerInterceptorLogsOnlyRequestIDOfMetadata(t *testing.T) {
	headers := maps.Clone(userRequest)
	headers["x-admin-user-id"] = "admin_test01"
	tests := []struct {
		name string
		fail error
	}{
		{name: "successful call (debug line)", fail: nil},
		{name: "failed call (warn line)", fail: connect.NewError(connect.CodeInternal, errors.New("boom"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf syncBuffer
			log := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
			url := startForwarder(t, startEcho(t), log, tt.fail)

			_, _ = call(t, url, headers)

			out := buf.String()
			if !strings.Contains(out, `"request_id":"req-test-0001"`) {
				t.Errorf("log has no request_id: %s", out)
			}
			for _, k := range allKeys {
				if k == "x-request-id" {
					continue
				}
				if strings.Contains(out, headers[k]) {
					t.Errorf("log contains the %s value %q: %s", k, headers[k], out)
				}
			}
		})
	}
}
