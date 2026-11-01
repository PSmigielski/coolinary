package rpc

import (
	"context"
	"net/http"
	"testing"
)

// incoming returns the ctx a handler sees for a request carrying these headers.
func incoming(t *testing.T, headers map[string]string) context.Context {
	t.Helper()
	h := http.Header{}
	for k, v := range headers {
		h.Set(k, v)
	}
	return WithMetadata(t.Context(), h.Get)
}

func TestUserID(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    string
	}{
		{name: "missing is empty", headers: nil, want: ""},
		{name: "set is returned", headers: map[string]string{"x-user-id": "user_test01"}, want: "user_test01"},
		{name: "mixed case is kept (ids are case-sensitive)", headers: map[string]string{"x-user-id": "uSeR_Test-01"}, want: "uSeR_Test-01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserID(incoming(t, tt.headers)); got != tt.want {
				t.Errorf("UserID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAdminUserID(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    string
	}{
		{name: "missing is empty", headers: nil, want: ""},
		{name: "set is returned", headers: map[string]string{"x-admin-user-id": "admin_test01"}, want: "admin_test01"},
		{name: "mixed case is kept (ids are case-sensitive)", headers: map[string]string{"x-admin-user-id": "AdMin_Test-01"}, want: "AdMin_Test-01"},
		{name: "user id is not an admin id", headers: map[string]string{"x-user-id": "user_test01"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AdminUserID(incoming(t, tt.headers)); got != tt.want {
				t.Errorf("AdminUserID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHouseholdID(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    string
	}{
		{name: "missing is empty", headers: nil, want: ""},
		{name: "set is returned", headers: map[string]string{"x-household-id": householdA}, want: householdA},
		{name: "mixed case is returned as received", headers: map[string]string{"x-household-id": householdAUpper}, want: householdAUpper},
		{name: "malformed is returned as received (CheckHousehold rejects it)", headers: map[string]string{"x-household-id": "household-1"}, want: "household-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HouseholdID(incoming(t, tt.headers)); got != tt.want {
				t.Errorf("HouseholdID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPlan(t *testing.T) {
	tests := []struct {
		name string
		plan string // "" sends no x-plan
		want PlanCode
	}{
		{name: "missing is free", plan: "", want: PlanFree},
		{name: "free", plan: "free", want: PlanFree},
		{name: "plus", plan: "plus", want: PlanPlus},
		{name: "mixed case plus is plus", plan: "Plus", want: PlanPlus},
		{name: "upper case plus is plus", plan: "PLUS", want: PlanPlus},
		{name: "upper case free is free", plan: "FREE", want: PlanFree},
		{name: "unknown is free", plan: "premium", want: PlanFree},
		{name: "malformed is free", plan: "plus;q=1", want: PlanFree},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := map[string]string{}
			if tt.plan != "" {
				headers["x-plan"] = tt.plan
			}
			if got := Plan(incoming(t, headers)); got != tt.want {
				t.Errorf("Plan() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLanguage(t *testing.T) {
	tests := []struct {
		name     string
		language string // "" sends no x-language
		want     string
	}{
		{name: "missing is en", language: "", want: "en"},
		{name: "primary subtag only", language: "pl", want: "pl"},
		{name: "region is dropped", language: "pl-PL", want: "pl"},
		{name: "upper case is lowered", language: "PL", want: "pl"},
		{name: "mixed case with region", language: "De-at", want: "de"},
		{name: "english with region", language: "en-GB", want: "en"},
		{name: "script and region are dropped", language: "zh-Hant-TW", want: "zh"},
		{name: "non-launch language is passed through", language: "fr-FR", want: "fr"},
		{name: "three-letter primary subtag is kept", language: "fil-PH", want: "fil"},
		{name: "digits in a later subtag are allowed", language: "de-DE-1996", want: "de"},
		{name: "underscore separator is malformed", language: "pl_PL", want: "en"},
		{name: "empty trailing subtag is malformed", language: "pl-", want: "en"},
		{name: "later subtag over 8 characters is malformed", language: "en-abcdefghi", want: "en"},
		{name: "digits in the primary subtag are malformed", language: "p1-PL", want: "en"},
		{name: "digits are malformed", language: "123", want: "en"},
		{name: "wildcard is malformed", language: "*", want: "en"},
		{name: "accept-language list is malformed", language: "pl-PL,pl;q=0.9", want: "en"},
		{name: "single-letter primary subtag is malformed", language: "x-private", want: "en"},
		{name: "empty primary subtag is malformed", language: "-PL", want: "en"},
		{name: "primary subtag over 8 letters is malformed", language: "abcdefghi", want: "en"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := map[string]string{}
			if tt.language != "" {
				headers["x-language"] = tt.language
			}
			if got := Language(incoming(t, headers)); got != tt.want {
				t.Errorf("Language() = %q, want %q", got, tt.want)
			}
		})
	}
}
