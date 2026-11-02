package rpc

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
)

// AdminOptions returns the interceptors for a <Svc>AdminService handler: those of Options, then a check that refuses
// any call without x-admin-user-id with PERMISSION_DENIED "admin only". Only the admin service calls these RPCs and
// it checks the admin role with identity first (ADM-02); services don't check the role themselves.
func AdminOptions(log *slog.Logger) []connect.HandlerOption {
	return append(Options(log), connect.WithInterceptors(adminOnly()))
}

func adminOnly() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if AdminUserID(ctx) == "" {
				return nil, connect.NewError(connect.CodePermissionDenied, errors.New("admin only"))
			}
			return next(ctx, req)
		}
	}
}
