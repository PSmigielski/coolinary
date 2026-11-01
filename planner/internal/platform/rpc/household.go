package rpc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"connectrpc.com/connect"
)

// CheckHousehold returns the household a user-scoped RPC acts on: the request's household_id field and
// x-household-id must agree when both are set (PERMISSION_DENIED otherwise); when only one is set, that one is
// used, so service-to-service calls may send the field alone. Neither set, or a value that isn't a UUID, is
// INVALID_ARGUMENT. The id is returned lowercased.
func CheckHousehold(ctx context.Context, requestHouseholdID string) (string, error) {
	fromRequest, err := normalizeUUID(requestHouseholdID)
	if err != nil {
		return "", err
	}
	fromMetadata, err := normalizeUUID(HouseholdID(ctx))
	if err != nil {
		return "", err
	}
	switch {
	case fromRequest == "" && fromMetadata == "":
		return "", connect.NewError(connect.CodeInvalidArgument, errors.New("household_id is required"))
	case fromRequest == "":
		return fromMetadata, nil
	case fromMetadata != "" && fromMetadata != fromRequest:
		return "", connect.NewError(connect.CodePermissionDenied, errors.New("household mismatch"))
	}
	return fromRequest, nil
}

// normalizeUUID lowercases a UUID in its canonical 8-4-4-4-12 form. "" stays "".
func normalizeUUID(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	invalid := connect.NewError(connect.CodeInvalidArgument, errors.New("household_id must be a UUID"))
	if len(s) != 36 {
		return "", invalid
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return "", invalid
			}
		default:
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
				return "", invalid
			}
		}
	}
	return strings.ToLower(s), nil
}

// HouseholdHash is the only form in which a household id may be logged: the first 8 hex characters of its SHA-256.
func HouseholdHash(householdID string) string {
	sum := sha256.Sum256([]byte(householdID))
	return hex.EncodeToString(sum[:4])
}
