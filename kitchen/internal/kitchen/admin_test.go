package kitchen

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	kitchenv1 "coolinary/kitchen/gen/coolinary/kitchen/v1"
	"coolinary/kitchen/gen/coolinary/kitchen/v1/kitchenv1connect"
	"coolinary/kitchen/internal/platform/db"
	"coolinary/kitchen/internal/platform/rpc"
	"coolinary/kitchen/migrations"
)

const testAdmin = "admin_test01"

// adminNow is the clock the AdminService tests run with.
var adminNow = time.Date(2026, 11, 2, 9, 30, 15, 123_000_000, time.UTC)

// newTestAdmin returns an AdminService on the migrated test database with its clock at adminNow, or skips when
// TEST_DATABASE_URL isn't set (make test sets it while make dev is up).
func newTestAdmin(t *testing.T) (*AdminService, *pgxpool.Pool) {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("needs TEST_DATABASE_URL (make dev, then make test s=kitchen)")
	}
	log := slog.New(slog.DiscardHandler)
	pool, err := db.Open(t.Context(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := db.Migrate(t.Context(), pool, migrations.FS, log); err != nil {
		t.Fatal(err)
	}
	s := NewAdmin(pool, log)
	s.now = func() time.Time { return adminNow }
	return s, pool
}

// adminCtx is the context of a call from the admin service, with x-admin-user-id set to adminID ("" sends none).
func adminCtx(t *testing.T, adminID string) context.Context {
	t.Helper()
	h := http.Header{}
	if adminID != "" {
		h.Set("x-admin-user-id", adminID)
	}
	return rpc.WithMetadata(t.Context(), h.Get)
}

// testEntity returns a random entity id no other test uses, and deletes its admin_changes rows after the test.
func testEntity(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatal(err)
	}
	id := "test-" + hex.EncodeToString(b[:])
	t.Cleanup(func() {
		if _, err := pool.Exec(context.Background(), "DELETE FROM admin_changes WHERE entity_id = $1", id); err != nil {
			t.Error(err)
		}
	})
	return id
}

// record writes one change in its own transaction, as an admin RPC does, and commits it.
func record(t *testing.T, s *AdminService, c adminChange) {
	t.Helper()
	tx, err := s.pool.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(t.Context()) }()
	if err := recordAdminChange(adminCtx(t, testAdmin), s.q.WithTx(tx), c); err != nil {
		t.Fatalf("recordAdminChange: %v", err)
	}
	if err := tx.Commit(t.Context()); err != nil {
		t.Fatal(err)
	}
}

// recordMany writes n changes to one entity in one statement, with actions "Change1" ... "Change<n>" in id order.
func recordMany(t *testing.T, pool *pgxpool.Pool, entityID string, n int) {
	t.Helper()
	_, err := pool.Exec(t.Context(), `INSERT INTO admin_changes (admin_user_id, action, entity_type, entity_id)
		SELECT $1, 'Change' || i, 'ingredient', $2 FROM generate_series(1, $3) AS i ORDER BY i`, testAdmin, entityID, n)
	if err != nil {
		t.Fatal(err)
	}
}

// actions lists the changes' actions in the order returned.
func actions(changes []*kitchenv1.AdminChange) []string {
	out := []string{}
	for _, c := range changes {
		out = append(out, c.GetAction())
	}
	return out
}

// changeActions lists "Change<from>" down to "Change<to>".
func changeActions(from, to int) []string {
	out := []string{}
	for i := from; i >= to; i-- {
		out = append(out, "Change"+strconv.Itoa(i))
	}
	return out
}

// requireCode fails the test unless err is a *connect.Error with this code and message.
func requireCode(t *testing.T, err error, code connect.Code, msg string) {
	t.Helper()
	var ce *connect.Error
	if !errors.As(err, &ce) {
		t.Fatalf("error = %v, want a *connect.Error %v %q", err, code, msg)
	}
	if ce.Code() != code || ce.Message() != msg {
		t.Errorf("error = %v %q, want %v %q", ce.Code(), ce.Message(), code, msg)
	}
}

// jsonEqual reports whether two JSON texts hold the same value ("" only equals "").
func jsonEqual(t *testing.T, a, b string) bool {
	t.Helper()
	if a == "" || b == "" {
		return a == b
	}
	var va, vb any
	if err := json.Unmarshal([]byte(a), &va); err != nil {
		t.Fatalf("not JSON: %q", a)
	}
	if err := json.Unmarshal([]byte(b), &vb); err != nil {
		t.Fatalf("not JSON: %q", b)
	}
	return reflect.DeepEqual(va, vb)
}

type testIngredient struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

func TestRecordAdminChange_WritesOneRowWithBeforeAndAfter(t *testing.T) {
	tests := []struct {
		name       string
		before     any
		after      any
		wantBefore string
		wantAfter  string
	}{
		{
			name:       "an edit keeps both versions",
			before:     testIngredient{Name: "Oat drink", Language: "en"},
			after:      testIngredient{Name: "Oat drink 1 l", Language: "en"},
			wantBefore: `{"name":"Oat drink","language":"en"}`,
			wantAfter:  `{"name":"Oat drink 1 l","language":"en"}`,
		},
		{
			name:      "a create has no before",
			before:    nil,
			after:     testIngredient{Name: "Oat drink", Language: "en"},
			wantAfter: `{"name":"Oat drink","language":"en"}`,
		},
		{
			name:       "a delete has no after",
			before:     testIngredient{Name: "Oat drink", Language: "en"},
			after:      (*testIngredient)(nil),
			wantBefore: `{"name":"Oat drink","language":"en"}`,
		},
	}
	s, pool := newTestAdmin(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := testEntity(t, pool)

			record(t, s, adminChange{Action: "UpdateIngredient", EntityType: "ingredient", EntityID: entity, Before: tt.before, After: tt.after})

			res, err := s.listAdminChanges(t.Context(), "ingredient", entity, 0, "")
			if err != nil {
				t.Fatalf("listAdminChanges: %v", err)
			}
			if len(res.GetChanges()) != 1 {
				t.Fatalf("listAdminChanges returned %d changes, want 1: %v", len(res.GetChanges()), res)
			}
			got := res.GetChanges()[0]
			var id int64
			var createdAt time.Time
			if err := pool.QueryRow(t.Context(), "SELECT id, created_at FROM admin_changes WHERE entity_id = $1", entity).Scan(&id, &createdAt); err != nil {
				t.Fatal(err)
			}
			want := &kitchenv1.AdminChange{
				Id: id, AdminUserId: testAdmin, Action: "UpdateIngredient", EntityType: "ingredient", EntityId: entity,
				CreatedAt: createdAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			}
			if !jsonEqual(t, got.GetBeforeJson(), tt.wantBefore) || !jsonEqual(t, got.GetAfterJson(), tt.wantAfter) {
				t.Errorf("before, after = %q, %q; want %q, %q", got.GetBeforeJson(), got.GetAfterJson(), tt.wantBefore, tt.wantAfter)
			}
			got.BeforeJson, got.AfterJson = "", ""
			if !proto.Equal(got, want) {
				t.Errorf("change = %v, want %v", got, want)
			}
		})
	}
}

func TestRecordAdminChange_RolledBackWithItsTransaction(t *testing.T) {
	s, pool := newTestAdmin(t)
	entity := testEntity(t, pool)
	tx, err := pool.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	err = recordAdminChange(adminCtx(t, testAdmin), s.q.WithTx(tx), adminChange{
		Action: "UpdateIngredient", EntityType: "ingredient", EntityID: entity, After: testIngredient{Name: "Oat drink"},
	})
	if err != nil {
		t.Fatalf("recordAdminChange: %v", err)
	}
	if err := tx.Rollback(t.Context()); err != nil {
		t.Fatal(err)
	}

	res, err := s.listAdminChanges(t.Context(), "ingredient", entity, 0, "")
	if err != nil {
		t.Fatalf("listAdminChanges: %v", err)
	}
	if len(res.GetChanges()) != 0 {
		t.Errorf("listAdminChanges after rollback = %v, want no changes", res.GetChanges())
	}
}

func TestRecordAdminChange_RefusesWithoutAdminUserID(t *testing.T) {
	s, pool := newTestAdmin(t)
	entity := testEntity(t, pool)
	tx, err := pool.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(t.Context()) }()

	err = recordAdminChange(adminCtx(t, ""), s.q.WithTx(tx), adminChange{
		Action: "UpdateIngredient", EntityType: "ingredient", EntityID: entity, After: testIngredient{Name: "Oat drink"},
	})

	if !errors.Is(err, errNoAdminUserID) {
		t.Errorf("recordAdminChange without x-admin-user-id = %v, want %v", err, errNoAdminUserID)
	}
}

func TestListAdminChanges_OneEntityNewestFirst(t *testing.T) {
	s, pool := newTestAdmin(t)
	entity, other := testEntity(t, pool), testEntity(t, pool)
	for _, c := range []adminChange{
		{Action: "First", EntityType: "ingredient", EntityID: entity},
		{Action: "OtherEntity", EntityType: "ingredient", EntityID: other},
		{Action: "Second", EntityType: "ingredient", EntityID: entity},
		{Action: "OtherType", EntityType: "product", EntityID: entity},
		{Action: "Third", EntityType: "ingredient", EntityID: entity},
	} {
		c.After = testIngredient{Name: "Oat drink"}
		record(t, s, c)
	}

	res, err := s.listAdminChanges(t.Context(), "ingredient", entity, 0, "")

	if err != nil {
		t.Fatalf("listAdminChanges: %v", err)
	}
	if got, want := actions(res.GetChanges()), []string{"Third", "Second", "First"}; !reflect.DeepEqual(got, want) {
		t.Errorf("actions = %v, want %v", got, want)
	}
	if got := res.GetPage().GetNextPageToken(); got != "" {
		t.Errorf("next_page_token = %q, want none on the last page", got)
	}
}

func TestListAdminChanges_PagesFollowOnFromTheToken(t *testing.T) {
	s, pool := newTestAdmin(t)
	entity := testEntity(t, pool)
	recordMany(t, pool, entity, 5)

	first, err := s.listAdminChanges(t.Context(), "ingredient", entity, 2, "")
	if err != nil {
		t.Fatalf("first page: %v", err)
	}
	second, err := s.listAdminChanges(t.Context(), "ingredient", entity, 2, first.GetPage().GetNextPageToken())
	if err != nil {
		t.Fatalf("second page: %v", err)
	}
	last, err := s.listAdminChanges(t.Context(), "ingredient", entity, 2, second.GetPage().GetNextPageToken())
	if err != nil {
		t.Fatalf("last page: %v", err)
	}

	pages := [][]string{actions(first.GetChanges()), actions(second.GetChanges()), actions(last.GetChanges())}
	want := [][]string{{"Change5", "Change4"}, {"Change3", "Change2"}, {"Change1"}}
	if !reflect.DeepEqual(pages, want) {
		t.Errorf("pages = %v, want %v", pages, want)
	}
	if first.GetPage().GetNextPageToken() == "" || second.GetPage().GetNextPageToken() == "" {
		t.Errorf("next_page_token missing on a full page: %q, %q", first.GetPage().GetNextPageToken(), second.GetPage().GetNextPageToken())
	}
	if got := last.GetPage().GetNextPageToken(); got != "" {
		t.Errorf("next_page_token = %q on the last page, want none", got)
	}
}

func TestListAdminChanges_PageSize(t *testing.T) {
	tests := []struct {
		name      string
		rows      int
		pageSize  int32
		wantCount int
		wantMore  bool
	}{
		{name: "50 per page when unset", rows: 51, pageSize: 0, wantCount: 50, wantMore: true},
		{name: "50 per page when negative", rows: 51, pageSize: -1, wantCount: 50, wantMore: true},
		{name: "a smaller page is kept", rows: 11, pageSize: 10, wantCount: 10, wantMore: true},
		{name: "200 is allowed", rows: 201, pageSize: 200, wantCount: 200, wantMore: true},
		{name: "over 200 is capped at 200", rows: 201, pageSize: 500, wantCount: 200, wantMore: true},
		{name: "a page that holds every row has no token", rows: 200, pageSize: 200, wantCount: 200, wantMore: false},
	}
	s, pool := newTestAdmin(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := testEntity(t, pool)
			recordMany(t, pool, entity, tt.rows)

			res, err := s.listAdminChanges(t.Context(), "ingredient", entity, tt.pageSize, "")

			if err != nil {
				t.Fatalf("listAdminChanges: %v", err)
			}
			if got, want := actions(res.GetChanges()), changeActions(tt.rows, tt.rows-tt.wantCount+1); !reflect.DeepEqual(got, want) {
				t.Errorf("got %d changes %v, want %d newest", len(got), got, tt.wantCount)
			}
			if got := res.GetPage().GetNextPageToken() != ""; got != tt.wantMore {
				t.Errorf("next_page_token = %q, want one: %v", res.GetPage().GetNextPageToken(), tt.wantMore)
			}
		})
	}
}

func TestListAdminChanges_RefusesBadInput(t *testing.T) {
	tests := []struct {
		name       string
		entityType string
		entityID   string
		pageToken  string
		wantMsg    string
	}{
		{name: "no entity type", entityType: "", entityID: "test-1", wantMsg: "entity_type and entity_id are required"},
		{name: "no entity id", entityType: "ingredient", entityID: "", wantMsg: "entity_type and entity_id are required"},
		{name: "a token that isn't a number", entityType: "ingredient", entityID: "test-1", pageToken: "abc", wantMsg: "page_token is invalid"},
		{name: "a zero token", entityType: "ingredient", entityID: "test-1", pageToken: "0", wantMsg: "page_token is invalid"},
		{name: "a negative token", entityType: "ingredient", entityID: "test-1", pageToken: "-5", wantMsg: "page_token is invalid"},
	}
	s := NewAdmin(nil, slog.New(slog.DiscardHandler))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := s.listAdminChanges(t.Context(), tt.entityType, tt.entityID, 0, tt.pageToken)

			if res != nil {
				t.Errorf("listAdminChanges = %v, want nil", res)
			}
			requireCode(t, err, connect.CodeInvalidArgument, tt.wantMsg)
		})
	}
}

func TestDeleteExpiredAdminChanges_KeepsThreeYears(t *testing.T) {
	s, pool := newTestAdmin(t)
	entity := testEntity(t, pool)
	cutoff := adminNow.AddDate(-3, 0, 0)
	for action, at := range map[string]time.Time{
		"JustOverThreeYears":  cutoff.Add(-time.Second),
		"ExactlyThreeYears":   cutoff,
		"JustUnderThreeYears": cutoff.Add(time.Second),
		"Recent":              adminNow.AddDate(0, -1, 0),
	} {
		if _, err := pool.Exec(t.Context(), `INSERT INTO admin_changes (admin_user_id, action, entity_type, entity_id, created_at)
			VALUES ($1, $2, 'ingredient', $3, $4)`, testAdmin, action, entity, at); err != nil {
			t.Fatal(err)
		}
	}

	n, err := s.DeleteExpiredAdminChanges(t.Context())

	if err != nil || n != 1 {
		t.Fatalf("DeleteExpiredAdminChanges = %d, %v; want the 1 expired row deleted", n, err)
	}
	res, err := s.listAdminChanges(t.Context(), "ingredient", entity, 0, "")
	if err != nil {
		t.Fatalf("listAdminChanges: %v", err)
	}
	got := map[string]bool{}
	for _, c := range res.GetChanges() {
		got[c.GetAction()] = true
	}
	if want := map[string]bool{"ExactlyThreeYears": true, "JustUnderThreeYears": true, "Recent": true}; !reflect.DeepEqual(got, want) {
		t.Errorf("rows left = %v, want %v", got, want)
	}
}

func TestGetQueueCounts_LeavesOutQueuesNotBuiltYet(t *testing.T) {
	s := NewAdmin(nil, slog.New(slog.DiscardHandler))
	s.now = func() time.Time { return adminNow }

	res, err := s.getQueueCounts(adminCtx(t, testAdmin))

	if err != nil {
		t.Fatalf("getQueueCounts: %v", err)
	}
	want := &kitchenv1.GetQueueCountsResponse{Counts: map[string]int64{}, ComputedAt: "2026-11-02T09:30:15.123Z"}
	if !proto.Equal(res, want) {
		t.Errorf("getQueueCounts = %v, want %v", res, want)
	}
}

// startAdminService serves KitchenAdminService as main.go mounts it, on a local h2c server.
func startAdminService(t *testing.T, s *AdminService) string {
	t.Helper()
	path, h := kitchenv1connect.NewKitchenAdminServiceHandler(s, rpc.AdminOptions(slog.New(slog.DiscardHandler))...)
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

func TestKitchenAdminService_EveryRPCRefusesCallsWithoutAdminUserID(t *testing.T) {
	url := startAdminService(t, NewAdmin(nil, slog.New(slog.DiscardHandler)))
	methods := kitchenv1.File_coolinary_kitchen_v1_admin_proto.Services().ByName("KitchenAdminService").Methods()
	if methods.Len() == 0 {
		t.Fatal("KitchenAdminService has no RPCs")
	}
	for i := range methods.Len() {
		name := string(methods.Get(i).Name())
		t.Run(name, func(t *testing.T) {
			// An empty message decodes as any request type, so one client shape reaches every RPC.
			client := connect.NewClient[emptypb.Empty, emptypb.Empty](rpc.NewHTTPClient(),
				url+"/"+kitchenv1connect.KitchenAdminServiceName+"/"+name, connect.WithGRPC())
			req := connect.NewRequest(&emptypb.Empty{})
			req.Header().Set("x-user-id", "user_test01")
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()

			_, err := client.CallUnary(ctx, req)

			requireCode(t, err, connect.CodePermissionDenied, "admin only")
		})
	}
}

func TestKitchenAdminService_ServesCallsWithAdminUserID(t *testing.T) {
	url := startAdminService(t, NewAdmin(nil, slog.New(slog.DiscardHandler)))
	client := kitchenv1connect.NewKitchenAdminServiceClient(rpc.NewHTTPClient(), url, connect.WithGRPC())
	req := connect.NewRequest(&kitchenv1.GetQueueCountsRequest{})
	req.Header().Set("x-admin-user-id", testAdmin)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	res, err := client.GetQueueCounts(ctx, req)

	if err != nil {
		t.Fatalf("GetQueueCounts with x-admin-user-id: %v", err)
	}
	if len(res.Msg.GetCounts()) != 0 {
		t.Errorf("counts = %v, want none before any queue is built", res.Msg.GetCounts())
	}
}
