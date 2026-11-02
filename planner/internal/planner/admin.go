package planner

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

	commonv1 "coolinary/planner/gen/coolinary/common/v1"
	plannerv1 "coolinary/planner/gen/coolinary/planner/v1"
	"coolinary/planner/gen/coolinary/planner/v1/plannerv1connect"
	"coolinary/planner/internal/planner/store"
	"coolinary/planner/internal/platform/rpc"
)

const (
	// Admin list RPCs return 50 rows per page by default and at most 200 (.claude/rules/go.md "Lists").
	adminPageSize    = 50
	adminMaxPageSize = 200
	// adminChangesKept is how long admin_changes rows are kept (Q016); adminChangesCleanup is how often older ones
	// are deleted (monthly).
	adminChangesKept    = 1 // year
	adminChangesCleanup = 30 * 24 * time.Hour
	// timeLayout is RFC 3339 with milliseconds; on a UTC time it ends in "Z".
	timeLayout = "2006-01-02T15:04:05.000Z07:00"
)

// AdminService implements coolinary.planner.v1.PlannerAdminService, called only by the admin service. main.go
// mounts it with rpc.AdminOptions, which refuses calls without x-admin-user-id. RPCs that change data write one
// admin_changes row with recordAdminChange in their transaction; read-only RPCs write nothing.
type AdminService struct {
	pool *pgxpool.Pool
	q    *store.Queries
	log  *slog.Logger
	now  func() time.Time
}

var _ plannerv1connect.PlannerAdminServiceHandler = (*AdminService)(nil)

func NewAdmin(pool *pgxpool.Pool, log *slog.Logger) *AdminService {
	return &AdminService{pool: pool, q: store.New(pool), log: log, now: time.Now}
}

// ListAdminChanges returns the change history of one target (e.g. one dish), newest first, paged by change id.
func (s *AdminService) ListAdminChanges(ctx context.Context, req *connect.Request[plannerv1.ListAdminChangesRequest]) (*connect.Response[plannerv1.ListAdminChangesResponse], error) {
	res, err := s.listAdminChanges(ctx, req.Msg.GetTarget(), req.Msg.GetPage().GetPageSize(), req.Msg.GetPage().GetPageToken())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(res), nil
}

func (s *AdminService) listAdminChanges(ctx context.Context, target string, pageSize int32, pageToken string) (*plannerv1.ListAdminChangesResponse, error) {
	if target == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("target is required"))
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

	rows, err := s.q.ListAdminChanges(ctx, store.ListAdminChangesParams{Target: target, AfterID: afterID, RowLimit: size + 1})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := &plannerv1.ListAdminChangesResponse{Page: &commonv1.PageResponse{}}
	for i, r := range rows {
		if int32(i) == size {
			res.Page.NextPageToken = strconv.FormatInt(rows[i-1].ID, 10)
			break
		}
		res.Changes = append(res.Changes, &plannerv1.AdminChange{
			AdminUserId: r.AdminUserID, Rpc: r.Rpc, Target: r.Target, SummaryJson: string(r.Summary),
			CreatedAt: r.CreatedAt.Time.UTC().Format(timeLayout),
		})
	}
	return res, nil
}

// GetQueueCounts returns the number of open items in each work queue, for the admin overview.
func (s *AdminService) GetQueueCounts(ctx context.Context, _ *connect.Request[plannerv1.GetQueueCountsRequest]) (*connect.Response[plannerv1.GetQueueCountsResponse], error) {
	res, err := s.getQueueCounts(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(res), nil
}

// getQueueCounts counts each queue in adminQueues. Each count is one indexed count(*).
func (s *AdminService) getQueueCounts(ctx context.Context) (*plannerv1.GetQueueCountsResponse, error) {
	counts := map[string]int64{}
	for key, count := range s.adminQueues() {
		n, err := count(ctx)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("count %s: %w", key, err))
		}
		counts[key] = n
	}
	return &plannerv1.GetQueueCountsResponse{Counts: counts, ComputedAt: s.now().UTC().Format(timeLayout)}, nil
}

// adminQueues maps each work queue's stable key to its count. A queue is added here by the task that builds it:
// review_queue (PLN-03-07), dietitian_pending (PLN-03-08), effort_tags_pending (PLN-03-06), promotion_candidates
// (PLN-11-03), open_batches (planner_batches not yet processed, PLN-03-04). Queues not built yet are left out of
// GetQueueCounts.
func (s *AdminService) adminQueues() map[string]func(context.Context) (int64, error) {
	return map[string]func(context.Context) (int64, error){}
}

// adminChange is one change an admin RPC records with recordAdminChange.
type adminChange struct {
	RPC    string // the RPC name, e.g. "ApproveDishVersion"
	Target string // what was changed: "dish:<id>", "import:<id>", ...
	// Summary is stored as JSON: ids and counts only, never free text from users. nil (also a nil pointer) is
	// stored as {}.
	Summary any
}

// errNoAdminUserID: recordAdminChange was called on a context without x-admin-user-id.
var errNoAdminUserID = errors.New("admin change without x-admin-user-id")

// recordAdminChange writes c as one admin_changes row, with the admin user id from ctx, through q, which must be the
// transaction's Queries of the change it records (q.WithTx(tx)), so the record commits or rolls back with it. The
// caller turns an error into INTERNAL and rolls back:
//
//	if err := recordAdminChange(ctx, qtx, adminChange{RPC: "ApproveDishVersion", Target: "dish:" + id, ...}); err != nil {
//		return nil, connect.NewError(connect.CodeInternal, err)
//	}
func recordAdminChange(ctx context.Context, q *store.Queries, c adminChange) error {
	admin := rpc.AdminUserID(ctx)
	if admin == "" {
		return errNoAdminUserID
	}
	summary, err := json.Marshal(c.Summary)
	if err != nil {
		return fmt.Errorf("encode %s summary: %w", c.RPC, err)
	}
	if string(summary) == "null" {
		summary = []byte("{}")
	}
	err = q.InsertAdminChange(ctx, store.InsertAdminChangeParams{AdminUserID: admin, Rpc: c.RPC, Target: c.Target, Summary: summary})
	if err != nil {
		return fmt.Errorf("insert %s admin change: %w", c.RPC, err)
	}
	return nil
}

// DeleteExpiredAdminChanges deletes admin_changes rows older than 1 year and returns how many it deleted.
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
