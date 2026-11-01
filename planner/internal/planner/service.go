// Package planner implements coolinary.planner.v1.PlannerService.
package planner

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"

	plannerv1 "coolinary/planner/gen/coolinary/planner/v1"
	"coolinary/planner/gen/coolinary/planner/v1/plannerv1connect"
	"coolinary/planner/internal/planner/store"
)

type Service struct {
	q   *store.Queries
	log *slog.Logger
}

var _ plannerv1connect.PlannerServiceHandler = (*Service)(nil)

func New(q *store.Queries, log *slog.Logger) *Service { return &Service{q: q, log: log} }

// SuggestDishes lists library dishes for a meal. Ranking by stock on hand
// (via kitchen) and expiry comes with the planner feature; share_on_hand is 0 for now.
func (s *Service) SuggestDishes(ctx context.Context, req *connect.Request[plannerv1.SuggestDishesRequest]) (*connect.Response[plannerv1.SuggestDishesResponse], error) {
	limit := req.Msg.GetLimit()
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	rows, err := s.q.ListDishesForMeal(ctx, store.ListDishesForMealParams{Meal: int16(req.Msg.GetMeal()), Limit: limit})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := &plannerv1.SuggestDishesResponse{}
	for _, r := range rows {
		res.Suggestions = append(res.Suggestions, &plannerv1.DishSuggestion{
			DishId: r.ID, Name: r.Name, KcalPerPortion: r.KcalPerPortion,
		})
	}
	return connect.NewResponse(res), nil
}
