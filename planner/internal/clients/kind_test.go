package clients

import (
	"errors"
	"fmt"
	"testing"

	"connectrpc.com/connect"
)

func TestKindOf_SortsErrorsIntoUnavailableRejectedAndBug(t *testing.T) {
	coded := func(c connect.Code) error { return connect.NewError(c, errors.New("from the called service")) }
	tests := []struct {
		name string
		err  error
		want Kind
	}{
		{name: "no error", err: nil, want: OK},
		{name: "unavailable", err: coded(connect.CodeUnavailable), want: Unavailable},
		{name: "deadline exceeded", err: coded(connect.CodeDeadlineExceeded), want: Unavailable},
		{name: "invalid argument", err: coded(connect.CodeInvalidArgument), want: Rejected},
		{name: "failed precondition", err: coded(connect.CodeFailedPrecondition), want: Rejected},
		{name: "not found", err: coded(connect.CodeNotFound), want: Rejected},
		{name: "permission denied", err: coded(connect.CodePermissionDenied), want: Rejected},
		{name: "resource exhausted (ai spend limit)", err: coded(connect.CodeResourceExhausted), want: Rejected},
		{name: "aborted (Q237)", err: coded(connect.CodeAborted), want: Rejected},
		{name: "already exists (Q237)", err: coded(connect.CodeAlreadyExists), want: Rejected},
		{name: "internal", err: coded(connect.CodeInternal), want: Bug},
		{name: "unknown", err: coded(connect.CodeUnknown), want: Bug},
		{name: "canceled (Q237)", err: coded(connect.CodeCanceled), want: Bug},
		{name: "unimplemented", err: coded(connect.CodeUnimplemented), want: Bug},
		{name: "unauthenticated", err: coded(connect.CodeUnauthenticated), want: Bug},
		{name: "out of range", err: coded(connect.CodeOutOfRange), want: Bug},
		{name: "data loss", err: coded(connect.CodeDataLoss), want: Bug},
		{name: "an error without a code", err: errors.New("plain"), want: Bug},
		{name: "a wrapped unavailable", err: fmt.Errorf("load profile: %w", coded(connect.CodeUnavailable)),
			want: Unavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KindOf(tt.err); got != tt.want {
				t.Errorf("KindOf(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
