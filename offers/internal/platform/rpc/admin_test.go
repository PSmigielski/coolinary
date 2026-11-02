package rpc

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"testing"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const procAdmin = "/coolinary.test.v1.TestAdminService/WhoAmI"

// startAdmin starts a stand-in admin RPC mounted with AdminOptions, as main.go mounts <Svc>AdminService. It answers
// with the admin user id it sees.
func startAdmin(t *testing.T) string {
	t.Helper()
	h := connect.NewUnaryHandler(procAdmin,
		func(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[wrapperspb.StringValue], error) {
			return connect.NewResponse(wrapperspb.String(AdminUserID(ctx))), nil
		}, AdminOptions(slog.New(slog.DiscardHandler))...)
	return startH2C(t, procAdmin, h)
}

func callAdmin(t *testing.T, url string, headers map[string]string) (*connect.Response[wrapperspb.StringValue], error) {
	t.Helper()
	client := connect.NewClient[emptypb.Empty, wrapperspb.StringValue](NewHTTPClient(), url+procAdmin, connect.WithGRPC())
	req := connect.NewRequest(&emptypb.Empty{})
	for k, v := range headers {
		req.Header().Set(k, v)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	return client.CallUnary(ctx, req)
}

func TestAdminOptions_RefuseCallsWithoutAdminUserID(t *testing.T) {
	emptyAdmin := maps.Clone(userRequest)
	emptyAdmin["x-admin-user-id"] = ""
	tests := []struct {
		name    string
		headers map[string]string
	}{
		{name: "no metadata at all", headers: map[string]string{}},
		{name: "a user request through the gateway (no admin id)", headers: userRequest},
		{name: "an empty admin id", headers: emptyAdmin},
	}
	url := startAdmin(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := callAdmin(t, url, tt.headers)

			var ce *connect.Error
			if !errors.As(err, &ce) {
				t.Fatalf("call = %v, %v; want a *connect.Error with code %v", res, err, connect.CodePermissionDenied)
			}
			if ce.Code() != connect.CodePermissionDenied || ce.Message() != "admin only" {
				t.Errorf("call error = %v %q, want %v %q", ce.Code(), ce.Message(), connect.CodePermissionDenied, "admin only")
			}
		})
	}
}

func TestAdminOptions_PassCallsWithAdminUserID(t *testing.T) {
	url := startAdmin(t)

	res, err := callAdmin(t, url, map[string]string{"x-admin-user-id": "admin_test01"})

	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if got := res.Msg.GetValue(); got != "admin_test01" {
		t.Errorf("the RPC saw admin user id %q, want %q", got, "admin_test01")
	}
}
