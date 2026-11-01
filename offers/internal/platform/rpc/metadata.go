package rpc

import (
	"context"
	"strings"
)

// Metadata keys forwarded between services (docs/service-conventions.md "gRPC").
const (
	keyRequestID   = "x-request-id"
	keyTraceparent = "traceparent"
	keyUserID      = "x-user-id"
	keyHouseholdID = "x-household-id"
	keyPlan        = "x-plan"
	keyLanguage    = "x-language"
	keyAdminUserID = "x-admin-user-id" // set only on calls from the admin service
)

// ForwardedHeaders are copied from an incoming request to outgoing calls.
var ForwardedHeaders = []string{
	keyRequestID, keyTraceparent, keyUserID, keyHouseholdID, keyPlan, keyLanguage, keyAdminUserID,
}

type metadataKey struct{}

// WithMetadata stores the forwarded headers of an incoming request in ctx.
// Values are kept as received; the helpers below read and normalise them.
// Never log them, except x-request-id (docs/EPICs KIT-01-01 rule 7).
func WithMetadata(ctx context.Context, get func(string) string) context.Context {
	md := map[string]string{}
	for _, k := range ForwardedHeaders {
		if v := get(k); v != "" {
			md[k] = v
		}
	}
	return context.WithValue(ctx, metadataKey{}, md)
}

func metadata(ctx context.Context, key string) string {
	md, _ := ctx.Value(metadataKey{}).(map[string]string)
	return md[key]
}

// UserID is the calling user's id (x-user-id), or "" when there is none.
func UserID(ctx context.Context) string { return metadata(ctx, keyUserID) }

// HouseholdID is x-household-id as received, or "". Use CheckHousehold to validate it against the request.
func HouseholdID(ctx context.Context) string { return metadata(ctx, keyHouseholdID) }

// AdminUserID is the admin's user id on calls from the admin service (x-admin-user-id), or "".
func AdminUserID(ctx context.Context) string { return metadata(ctx, keyAdminUserID) }

// PlanCode is the household's plan as the gateway forwards it in x-plan (from the token, so possibly stale).
type PlanCode string

const (
	PlanFree PlanCode = "free"
	PlanPlus PlanCode = "plus"
)

// Plan is x-plan, case-insensitive; free when missing or unknown.
func Plan(ctx context.Context) PlanCode {
	if strings.EqualFold(metadata(ctx, keyPlan), string(PlanPlus)) {
		return PlanPlus
	}
	return PlanFree
}

// Language is the primary subtag of the BCP 47 tag in x-language, lowercased ("pl-PL" → "pl"). Languages other
// than pl, de and en are returned as they are, so language-aware code picks its own fallback. Missing or
// malformed → "en". Well-formed here means subtags of 1–8 letters or digits joined by "-", the first one
// 2–8 letters.
func Language(ctx context.Context) string {
	subtags := strings.Split(metadata(ctx, keyLanguage), "-")
	primary := subtags[0]
	if len(primary) < 2 || !isAlnum(primary, false) {
		return "en"
	}
	for _, s := range subtags[1:] {
		if !isAlnum(s, true) {
			return "en"
		}
	}
	return strings.ToLower(primary)
}

// isAlnum reports whether s is 1–8 ASCII letters (and digits, if allowed).
func isAlnum(s string, digits bool) bool {
	if len(s) < 1 || len(s) > 8 {
		return false
	}
	for _, c := range s {
		letter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
		if !letter && (!digits || c < '0' || c > '9') {
			return false
		}
	}
	return true
}
