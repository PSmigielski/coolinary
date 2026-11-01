package rpc

import (
	"errors"
	"testing"

	"connectrpc.com/connect"
)

const (
	householdA      = "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b"
	householdAUpper = "0190A1B2-C3D4-7E5F-8A9B-0C1D2E3F4A5B"
	householdB      = "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5c"
)

func TestCheckHousehold(t *testing.T) {
	tests := []struct {
		name     string
		request  string // the request's household_id field
		metadata string // x-household-id; "" sends none
		want     string
		wantCode connect.Code // 0: no error
		wantMsg  string
	}{
		{name: "request and metadata agree", request: householdA, metadata: householdA, want: householdA},
		{name: "only metadata is set", request: "", metadata: householdA, want: householdA},
		{name: "only request is set (service-to-service)", request: householdA, metadata: "", want: householdA},
		{name: "upper case request is lowered", request: householdAUpper, metadata: "", want: householdA},
		{name: "upper case metadata is lowered", request: "", metadata: householdAUpper, want: householdA},
		{name: "same household in different case agrees", request: householdA, metadata: householdAUpper, want: householdA},
		{
			name: "request and metadata differ", request: householdA, metadata: householdB,
			wantCode: connect.CodePermissionDenied, wantMsg: "household mismatch",
		},
		{
			name: "neither is set", request: "", metadata: "",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id is required",
		},
		{
			name: "request is not a uuid", request: "household-1", metadata: "",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
		{
			name: "metadata is not a uuid", request: "", metadata: "household-1",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
		{
			name: "malformed request is refused even when metadata is valid", request: "household-1", metadata: householdA,
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
		{
			name: "uuid without hyphens is malformed", request: "0190a1b2c3d47e5f8a9b0c1d2e3f4a5b", metadata: "",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
		{
			name: "uuid in braces is malformed", request: "{" + householdA + "}", metadata: "",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
		{
			name: "non-hex character is malformed", request: "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5g", metadata: "",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
		{
			name: "hyphens in the wrong place are malformed", request: "0190a1b2c-3d4-7e5f-8a9b-0c1d2e3f4a5b", metadata: "",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := map[string]string{}
			if tt.metadata != "" {
				headers["x-household-id"] = tt.metadata
			}
			got, err := CheckHousehold(incoming(t, headers), tt.request)
			if tt.wantCode == 0 {
				if err != nil {
					t.Fatalf("CheckHousehold() error = %v, want none", err)
				}
				if got != tt.want {
					t.Errorf("CheckHousehold() = %q, want %q", got, tt.want)
				}
				return
			}
			var ce *connect.Error
			if !errors.As(err, &ce) {
				t.Fatalf("CheckHousehold() error = %v, want a *connect.Error with code %v", err, tt.wantCode)
			}
			if ce.Code() != tt.wantCode || ce.Message() != tt.wantMsg {
				t.Errorf("CheckHousehold() error = %v %q, want %v %q", ce.Code(), ce.Message(), tt.wantCode, tt.wantMsg)
			}
			if got != "" {
				t.Errorf("CheckHousehold() = %q on error, want empty", got)
			}
		})
	}
}

func TestHouseholdHash(t *testing.T) {
	// printf '%s' 0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b | sha256sum | cut -c1-8
	if got, want := HouseholdHash(householdA), "3bdb3cc3"; got != want {
		t.Errorf("HouseholdHash() = %q, want %q", got, want)
	}
}
