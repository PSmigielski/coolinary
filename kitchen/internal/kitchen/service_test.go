package kitchen

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"testing"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	kitchenv1 "coolinary/kitchen/gen/coolinary/kitchen/v1"
	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/cache"
	"coolinary/kitchen/internal/platform/db"
	"coolinary/kitchen/internal/platform/rpc"
	"coolinary/kitchen/migrations"
)

const (
	householdA = "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b"
	householdB = "0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5c"
)

// newTestService returns a Service on the test database and Redis, or skips when they aren't set
// (make test sets them while make dev is up).
func newTestService(t *testing.T) *Service {
	t.Helper()
	dbURL, cacheURL := os.Getenv("TEST_DATABASE_URL"), os.Getenv("TEST_REDIS_CACHE_URL")
	if dbURL == "" || cacheURL == "" {
		t.Skip("needs TEST_DATABASE_URL and TEST_REDIS_CACHE_URL (make dev, then make test s=kitchen)")
	}
	log := slog.New(slog.DiscardHandler)
	pool, err := db.Open(t.Context(), dbURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := db.Migrate(t.Context(), pool, migrations.FS, log); err != nil {
		t.Fatal(err)
	}
	redisCache, err := cache.Open(t.Context(), cacheURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = redisCache.Close() })
	return New(store.New(pool), redisCache, log)
}

// unknownEAN returns a random 13-digit code in the GS1 example range that no test or migration stores.
func unknownEAN(t *testing.T) string {
	t.Helper()
	n, err := rand.Int(rand.Reader, big.NewInt(1e10))
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("952%010d", n.Int64())
}

func TestResolveBarcode_HouseholdGuard(t *testing.T) {
	s := newTestService(t)
	tests := []struct {
		name     string
		request  string // ResolveBarcodeRequest.household_id
		metadata string // x-household-id; "" sends none
		wantCode connect.Code
		wantMsg  string
	}{
		{
			name: "household in request and metadata differ is refused", request: householdA, metadata: householdB,
			wantCode: connect.CodePermissionDenied, wantMsg: "household mismatch",
		},
		{
			name: "request household that is not a uuid is refused", request: "household-1", metadata: "",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
		{
			name: "metadata household that is not a uuid is refused", request: "", metadata: "household-1",
			wantCode: connect.CodeInvalidArgument, wantMsg: "household_id must be a UUID",
		},
		{name: "same household in request and metadata is looked up", request: householdA, metadata: householdA},
		{name: "household only in the request is looked up (service-to-service)", request: householdA, metadata: ""},
		{name: "household only in metadata is looked up", request: "", metadata: householdA},
		{name: "no household is looked up in the shared catalogue (public route)", request: "", metadata: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			if tt.metadata != "" {
				h.Set("x-household-id", tt.metadata)
			}
			ctx := rpc.WithMetadata(t.Context(), h.Get)

			res, err := s.resolveBarcode(ctx, unknownEAN(t), tt.request)

			if tt.wantCode == 0 {
				if err != nil {
					t.Fatalf("resolveBarcode() error = %v, want none", err)
				}
				if want := (&kitchenv1.ResolveBarcodeResponse{Found: false}); !proto.Equal(res, want) {
					t.Errorf("resolveBarcode() = %v, want %v", res, want)
				}
				return
			}
			var ce *connect.Error
			if !errors.As(err, &ce) {
				t.Fatalf("resolveBarcode() = %v, %v; want a *connect.Error with code %v", res, err, tt.wantCode)
			}
			if ce.Code() != tt.wantCode || ce.Message() != tt.wantMsg {
				t.Errorf("resolveBarcode() error = %v %q, want %v %q", ce.Code(), ce.Message(), tt.wantCode, tt.wantMsg)
			}
		})
	}
}
