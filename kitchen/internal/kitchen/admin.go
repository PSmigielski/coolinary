package kitchen

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	commonv1 "coolinary/kitchen/gen/coolinary/common/v1"
	kitchenv1 "coolinary/kitchen/gen/coolinary/kitchen/v1"
	"coolinary/kitchen/gen/coolinary/kitchen/v1/kitchenv1connect"
	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/rpc"
)

const (
	// Admin list RPCs return 50 rows per page by default and at most 200 (.claude/rules/go.md "Lists").
	adminPageSize    = 50
	adminMaxPageSize = 200
	// adminChangesKept is how long admin_changes rows are kept; adminChangesCleanup is how often older ones are
	// deleted (monthly).
	adminChangesKept    = 3 // years
	adminChangesCleanup = 30 * 24 * time.Hour
	// timeLayout is RFC 3339 with milliseconds, as event envelopes use; on a UTC time it ends in "Z".
	timeLayout = "2006-01-02T15:04:05.000Z07:00"
)

// AdminService implements coolinary.kitchen.v1.KitchenAdminService, called only by the admin service. main.go
// mounts it with rpc.AdminOptions, which refuses calls without x-admin-user-id. RPCs that change data write one
// admin_changes row with recordAdminChange in their transaction; read-only RPCs write nothing.
type AdminService struct {
	pool *pgxpool.Pool
	q    *store.Queries
	log  *slog.Logger
	now  func() time.Time
}

var _ kitchenv1connect.KitchenAdminServiceHandler = (*AdminService)(nil)

func NewAdmin(pool *pgxpool.Pool, log *slog.Logger) *AdminService {
	return &AdminService{pool: pool, q: store.New(pool), log: log, now: time.Now}
}

// ListAdminChanges returns the change history of one entity, newest first, paged by change id.
func (s *AdminService) ListAdminChanges(ctx context.Context, req *connect.Request[kitchenv1.ListAdminChangesRequest]) (*connect.Response[kitchenv1.ListAdminChangesResponse], error) {
	res, err := s.listAdminChanges(ctx, req.Msg.GetEntityType(), req.Msg.GetEntityId(), req.Msg.GetPage().GetPageSize(), req.Msg.GetPage().GetPageToken())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(res), nil
}

func (s *AdminService) listAdminChanges(ctx context.Context, entityType, entityID string, pageSize int32, pageToken string) (*kitchenv1.ListAdminChangesResponse, error) {
	if entityType == "" || entityID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("entity_type and entity_id are required"))
	}
	var afterID int64
	if pageToken != "" {
		id, err := strconv.ParseInt(pageToken, 10, 64)
		if err != nil || id <= 0 {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("page_token is invalid"))
		}
		afterID = id
	}
	size := pageSize
	if size <= 0 {
		size = adminPageSize
	}
	size = min(size, adminMaxPageSize)

	rows, err := s.q.ListAdminChanges(ctx, store.ListAdminChangesParams{
		EntityType: entityType, EntityID: entityID, AfterID: afterID, RowLimit: size + 1,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := &kitchenv1.ListAdminChangesResponse{Page: &commonv1.PageResponse{}}
	for i, r := range rows {
		if int32(i) == size {
			res.Page.NextPageToken = strconv.FormatInt(rows[i-1].ID, 10)
			break
		}
		res.Changes = append(res.Changes, &kitchenv1.AdminChange{
			Id: r.ID, AdminUserId: r.AdminUserID, Action: r.Action, EntityType: r.EntityType, EntityId: r.EntityID,
			BeforeJson: string(r.Before), AfterJson: string(r.After),
			CreatedAt: r.CreatedAt.Time.UTC().Format(timeLayout),
		})
	}
	return res, nil
}

// GetQueueCounts returns the number of open items in each work queue, for the admin overview.
func (s *AdminService) GetQueueCounts(ctx context.Context, _ *connect.Request[kitchenv1.GetQueueCountsRequest]) (*connect.Response[kitchenv1.GetQueueCountsResponse], error) {
	res, err := s.getQueueCounts(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(res), nil
}

// getQueueCounts counts each queue in adminQueues. Each count is one indexed count(*).
func (s *AdminService) getQueueCounts(ctx context.Context) (*kitchenv1.GetQueueCountsResponse, error) {
	counts := map[string]int64{}
	for key, count := range s.adminQueues() {
		n, err := count(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("count %s: %w", key, err))
		}
		counts[key] = n
	}
	return &kitchenv1.GetQueueCountsResponse{Counts: counts, ComputedAt: s.now().UTC().Format(timeLayout)}, nil
}

// adminQueues maps each work queue's stable key to its count. A queue is added here by the task that builds it:
// match_reviews (KIT-03-06), nutrition_suggestions (KIT-06-06), unreviewed_names (KIT-02-10), product_reviews
// (KIT-05-05), parked_items (KIT-02-08, KIT-06-05), open_batches (KIT-01-08), dead_letters (KIT-01-05). Queues not
// built yet are left out of GetQueueCounts.
func (s *AdminService) adminQueues() map[string]func(context.Context) (int64, error) {
	return map[string]func(context.Context) (int64, error){}
}

// adminChange is one change an admin RPC records with recordAdminChange.
type adminChange struct {
	Action     string // the RPC name, e.g. "MergeIngredients"
	EntityType string
	EntityID   string
	// Before and After are the entity's catalogue fields (never household data), stored as JSON. Before is nil
	// when the change created the entity, After is nil when it deleted it.
	Before any
	After  any
}

// errNoAdminUserID: recordAdminChange was called on a context without x-admin-user-id.
var errNoAdminUserID = errors.New("admin change without x-admin-user-id")

// recordAdminChange writes c as one admin_changes row, with the admin user id from ctx, through q, which must be the
// transaction's Queries of the change it records (q.WithTx(tx)), so the record commits or rolls back with it. The
// caller turns an error into INTERNAL and rolls back:
//
//	if err := recordAdminChange(ctx, qtx, adminChange{Action: "MergeIngredients", ...}); err != nil {
//		return nil, connect.NewError(connect.CodeInternal, err)
//	}
func recordAdminChange(ctx context.Context, q *store.Queries, c adminChange) error {
	admin := rpc.AdminUserID(ctx)
	if admin == "" {
		return errNoAdminUserID
	}
	before, err := entityJSON(c.Before)
	if err != nil {
		return fmt.Errorf("encode %s %s before: %w", c.Action, c.EntityType, err)
	}
	after, err := entityJSON(c.After)
	if err != nil {
		return fmt.Errorf("encode %s %s after: %w", c.Action, c.EntityType, err)
	}
	err = q.InsertAdminChange(ctx, store.InsertAdminChangeParams{
		AdminUserID: admin, Action: c.Action, EntityType: c.EntityType, EntityID: c.EntityID, Before: before, After: after,
	})
	if err != nil {
		return fmt.Errorf("insert %s admin change: %w", c.Action, err)
	}
	return nil
}

// entityJSON encodes v; nil (also a nil pointer) is no JSON at all, stored as NULL.
func entityJSON(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil || string(b) == "null" {
		return nil, err
	}
	return b, nil
}

// DeleteExpiredAdminChanges deletes admin_changes rows older than 3 years and returns how many it deleted.
func (s *AdminService) DeleteExpiredAdminChanges(ctx context.Context) (int64, error) {
	cutoff := s.now().AddDate(-adminChangesKept, 0, 0)
	n, err := s.q.DeleteAdminChangesBefore(ctx, pgtype.Timestamptz{Time: cutoff, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("delete expired admin changes: %w", err)
	}
	return n, nil
}

// RunCleanup deletes expired admin_changes rows at start and then monthly, until ctx is done. Failures are logged.
func (s *AdminService) RunCleanup(ctx context.Context) {
	for {
		if _, err := s.DeleteExpiredAdminChanges(ctx); err != nil && ctx.Err() == nil {
			s.log.Error("admin changes cleanup failed", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(adminChangesCleanup):
		}
	}
}
