// Package planner implements coolinary.planner.v1.PlannerService.
package planner

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"

	plannerv1 "coolinary/planner/gen/coolinary/planner/v1"
	"coolinary/planner/gen/coolinary/planner/v1/plannerv1connect"
	"coolinary/planner/internal/clients"
	"coolinary/planner/internal/planner/store"
)

type Service struct {
	q     *store.Queries
	calls *clients.Clients // kitchen, offers, ai, notifier, identity
	log   *slog.Logger
}

var _ plannerv1connect.PlannerServiceHandler = (*Service)(nil)

func New(q *store.Queries, calls *clients.Clients, log *slog.Logger) *Service {
	return &Service{q: q, calls: calls, log: log}
}

// SuggestDishes lists library dishes for a meal. Ranking by stock on hand
// (via kitchen) and expiry comes with the planner feature; share_on_hand is 0 for now.
func (s *Service) SuggestDishes(ctx context.Context, req *connect.Request[plannerv1.SuggestDishesRequest]) (*connect.Response[plannerv1.SuggestDishesResponse], error) {
	meal := req.Msg.GetMeal()
	if _, ok := plannerv1.Meal_name[int32(meal)]; !ok {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("unknown meal"))
	}
	limit := req.Msg.GetLimit()
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	//nolint:gosec // G115: meal is a known Meal value (checked above), so it fits the smallint column
	rows, err := s.q.ListDishesForMeal(ctx, store.ListDishesForMealParams{Meal: int16(meal), Limit: limit})
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
